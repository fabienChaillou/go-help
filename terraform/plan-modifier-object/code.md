```go
// main.go
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/example/terraform-provider-custom/internal/provider"
)

var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/example/custom",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}

// internal/provider/provider.go
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type customProvider struct {
	version string
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &customProvider{
			version: version,
		}
	}
}

func (p *customProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "custom"
	resp.Version = p.version
}

func (p *customProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Optional: true,
			},
		},
	}
}

func (p *customProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Configuration du provider (client API, etc.)
}

func (p *customProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewApplicationResource,
	}
}

func (p *customProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

// internal/provider/application_resource.go
package provider

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ApplicationResource{}

func NewApplicationResource() resource.Resource {
	return &ApplicationResource{}
}

type ApplicationResource struct{}

type ApplicationResourceModel struct {
	ID     types.String `tfsdk:"id"`
	Name   types.String `tfsdk:"name"`
	Config types.Object `tfsdk:"config"`
}

type ConfigModel struct {
	Environment types.String `tfsdk:"environment"`
	Endpoint    types.String `tfsdk:"endpoint"`
	Port        types.Int64  `tfsdk:"port"`
	SSL         types.Object `tfsdk:"ssl"`
}

type SSLModel struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	CACert  types.String `tfsdk:"ca_cert"`
}

func (r *ApplicationResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_application"
}

func (r *ApplicationResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Application resource avec configuration complexe",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Application identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Application name",
				Required:            true,
			},
			"config": schema.SingleNestedAttribute{
				MarkdownDescription: "Application configuration",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"environment": schema.StringAttribute{
						MarkdownDescription: "Environment name",
						Required:            true,
					},
					"endpoint": schema.StringAttribute{
						MarkdownDescription: "Generated endpoint URL",
						Computed:            true,
					},
					"port": schema.Int64Attribute{
						MarkdownDescription: "Generated port number",
						Computed:            true,
					},
					"ssl": schema.SingleNestedAttribute{
						MarkdownDescription: "SSL configuration",
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "SSL enabled status",
								Computed:            true,
							},
							"ca_cert": schema.StringAttribute{
								MarkdownDescription: "CA certificate",
								Computed:            true,
							},
						},
					},
				},
				// Utilisation de UseStateForUnknown pour éviter les diffs inutiles
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *ApplicationResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Configuration du client si nécessaire
}

func (r *ApplicationResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Simulation de la création avec des valeurs calculées
	data.ID = types.StringValue(fmt.Sprintf("app-%d", rand.Int()))
	
	// Création des valeurs calculées pour la configuration
	configModel := ConfigModel{}
	resp.Diagnostics.Append(data.Config.As(ctx, &configModel, types.ObjectAsOptions{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Génération des valeurs calculées
	configModel.Endpoint = types.StringValue(fmt.Sprintf("https://%s-api.example.com", configModel.Environment.ValueString()))
	configModel.Port = types.Int64Value(8080 + rand.Int63n(1000))
	
	// Configuration SSL calculée
	sslModel := SSLModel{
		Enabled: types.BoolValue(true),
		CACert:  types.StringValue("-----BEGIN CERTIFICATE-----\nMIIBIjANBgkqhki..."),
	}
	
	sslObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"enabled": types.BoolType,
		"ca_cert": types.StringType,
	}, sslModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	configModel.SSL = sslObj

	// Reconversion en Object
	configObj, diags := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"environment": types.StringType,
		"endpoint":    types.StringType,
		"port":        types.Int64Type,
		"ssl":         types.ObjectType{AttrTypes: map[string]attr.Type{
			"enabled": types.BoolType,
			"ca_cert": types.StringType,
		}},
	}, configModel)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	
	data.Config = configObj

	// Simulation d'un délai API
	time.Sleep(100 * time.Millisecond)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Simulation de la lecture - ici on garde l'état tel quel
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Logique de mise à jour similaire au Create
	// Les valeurs calculées peuvent changer selon la logique métier

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ApplicationResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ApplicationResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Simulation de la suppression
	time.Sleep(50 * time.Millisecond)
}

// internal/provider/application_resource_test.go
package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccApplicationResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Create and Read testing
			{
				Config: testAccApplicationResourceConfig("test-app", "production"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("custom_application.test", "name", "test-app"),
					resource.TestCheckResourceAttr("custom_application.test", "config.environment", "production"),
					resource.TestCheckResourceAttrSet("custom_application.test", "config.endpoint"),
					resource.TestCheckResourceAttrSet("custom_application.test", "config.port"),
					resource.TestCheckResourceAttr("custom_application.test", "config.ssl.enabled", "true"),
					resource.TestCheckResourceAttrSet("custom_application.test", "config.ssl.ca_cert"),
				),
			},
			// ImportState testing
			{
				ResourceName:      "custom_application.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update and Read testing
			{
				Config: testAccApplicationResourceConfig("test-app-updated", "staging"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("custom_application.test", "name", "test-app-updated"),
					resource.TestCheckResourceAttr("custom_application.test", "config.environment", "staging"),
					// Vérifier que les valeurs calculées sont toujours présentes
					resource.TestCheckResourceAttrSet("custom_application.test", "config.endpoint"),
					resource.TestCheckResourceAttrSet("custom_application.test", "config.port"),
				),
			},
		},
	})
}

func TestAccApplicationResource_UseStateForUnknown(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccApplicationResourceConfig("test-app", "production"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("custom_application.test", "config.endpoint"),
					resource.TestCheckResourceAttrSet("custom_application.test", "config.port"),
				),
			},
			// Test que les valeurs calculées ne changent pas lors d'un re-apply
			{
				Config: testAccApplicationResourceConfig("test-app", "production"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Les valeurs doivent rester les mêmes grâce à UseStateForUnknown
					resource.TestCheckResourceAttrSet("custom_application.test", "config.endpoint"),
					resource.TestCheckResourceAttrSet("custom_application.test", "config.port"),
				),
				PlanOnly: true, // Vérifier qu'il n'y a pas de changement planifié
			},
		},
	})
}

func testAccApplicationResourceConfig(name, environment string) string {
	return fmt.Sprintf(`
resource "custom_application" "test" {
  name = %[1]q
  config = {
    environment = %[2]q
  }
}
`, name, environment)
}

// internal/provider/provider_test.go
package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

const (
	providerConfig = `
provider "custom" {}
`
)

var (
	testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
		"custom": providerserver.NewProtocol6WithError(New("test")()),
	}
)

func testAccPreCheck(t *testing.T) {
	// Ajouter des vérifications préalables si nécessaire
}

func TestAccProvider(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig,
			},
		},
	})
}

// go.mod
module github.com/example/terraform-provider-custom

go 1.21

require (
	github.com/hashicorp/terraform-plugin-framework v1.4.2
	github.com/hashicorp/terraform-plugin-go v0.19.1
	github.com/hashicorp/terraform-plugin-testing v1.5.1
)
```

```makefile
// Makefile
default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Build provider
.PHONY: build
build:
	go build -o terraform-provider-custom

# Install provider locally
.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/registry.terraform.io/example/custom/0.1.0/linux_amd64
	mv terraform-provider-custom ~/.terraform.d/plugins/registry.terraform.io/example/custom/0.1.0/linux_amd64

# Run tests
.PHONY: test
test:
	go test -v ./...
```
