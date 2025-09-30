```go
// internal/provider/custom_plan_modifiers.go
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

// Plan modifier personnalisé pour remplacer null par une valeur par défaut
type defaultValueIfNullModifier struct {
	defaultValue attr.Value
}

func DefaultValueIfNull(defaultValue attr.Value) planmodifier.String {
	return &defaultValueIfNullModifier{
		defaultValue: defaultValue,
	}
}

func (m *defaultValueIfNullModifier) Description(ctx context.Context) string {
	return "Set default value if null"
}

func (m *defaultValueIfNullModifier) MarkdownDescription(ctx context.Context) string {
	return "Set default value if null"
}

func (m *defaultValueIfNullModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// Si la valeur planifiée est null et qu'il n'y a pas de valeur dans l'état
	if req.PlanValue.IsNull() && req.StateValue.IsNull() {
		resp.PlanValue = m.defaultValue.(types.String)
	}
}

// Plan modifier pour générer une valeur calculée si null
type generateIfNullModifier struct {
	generator func(ctx context.Context, config attr.Value, state attr.Value) (types.String, error)
}

func GenerateIfNull(generator func(ctx context.Context, config attr.Value, state attr.Value) (types.String, error)) planmodifier.String {
	return &generateIfNullModifier{
		generator: generator,
	}
}

func (m *generateIfNullModifier) Description(ctx context.Context) string {
	return "Generate value if null"
}

func (m *generateIfNullModifier) MarkdownDescription(ctx context.Context) string {
	return "Generate value if null"
}

func (m *generateIfNullModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() {
		generatedValue, err := m.generator(ctx, req.ConfigValue, req.StateValue)
		if err != nil {
			resp.Diagnostics.AddError("Generation Error", fmt.Sprintf("Failed to generate value: %s", err))
			return
		}
		resp.PlanValue = generatedValue
	}
}

// internal/provider/server_resource.go
package provider

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &ServerResource{}

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

type ServerResource struct{}

type ServerResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Environment types.String `tfsdk:"environment"`
	Hostname    types.String `tfsdk:"hostname"`
	IPAddress   types.String `tfsdk:"ip_address"`
	Password    types.String `tfsdk:"password"`
	Tags        types.Map    `tfsdk:"tags"`
}

func (r *ServerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Server resource with null-to-value transformations",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"environment": schema.StringAttribute{
				Optional: true,
				// Si null, utilise "development" par défaut
				PlanModifiers: []planmodifier.String{
					DefaultValueIfNull(types.StringValue("development")),
				},
			},
			"hostname": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// Génère un hostname si null basé sur le nom
				PlanModifiers: []planmodifier.String{
					GenerateIfNull(func(ctx context.Context, config attr.Value, state attr.Value) (types.String, error) {
						// Accéder au nom depuis la configuration
						// Note: Dans un cas réel, vous auriez accès au plan complet
						return types.StringValue("server-" + fmt.Sprintf("%d", time.Now().Unix())), nil
					}),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
				// IP générée pendant l'apply si null
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Computed:  true,
				Sensitive: true,
				// Mot de passe généré si null
				PlanModifiers: []planmodifier.String{
					GenerateIfNull(func(ctx context.Context, config attr.Value, state attr.Value) (types.String, error) {
						return generatePassword(16)
					}),
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"tags": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				// Tags par défaut si null
			},
		},
	}
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Génération de l'ID
	data.ID = types.StringValue(fmt.Sprintf("srv-%d", time.Now().Unix()))

	// Gestion des valeurs null -> computed pendant l'apply
	if data.IPAddress.IsNull() {
		// Simulation d'allocation d'IP par l'API
		data.IPAddress = types.StringValue(fmt.Sprintf("10.0.%d.%d", 
			randInt(1, 255), randInt(1, 254)))
	}

	// Gestion des tags par défaut si null
	if data.Tags.IsNull() {
		defaultTags := map[string]attr.Value{
			"managed_by":  types.StringValue("terraform"),
			"created_at":  types.StringValue(time.Now().Format(time.RFC3339)),
			"environment": data.Environment,
		}
		tagsValue, diags := types.MapValue(types.StringType, defaultTags)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		data.Tags = tagsValue
	}

	// Simulation d'appel API
	time.Sleep(100 * time.Millisecond)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Simulation: certaines valeurs peuvent être mises à jour par l'API externe
	// mais on garde l'état existant grâce aux plan modifiers

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Même logique que Create pour les nouvelles valeurs null
	if data.IPAddress.IsNull() {
		data.IPAddress = types.StringValue(fmt.Sprintf("10.0.%d.%d", 
			randInt(1, 255), randInt(1, 254)))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Simulation de suppression
	time.Sleep(50 * time.Millisecond)
}

// internal/provider/advanced_plan_modifier.go
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// Plan modifier avancé qui peut accéder à tout le plan
type conditionalValueModifier struct {
	condition func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) bool
	value     func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (types.String, error)
}

func ConditionalValue(
	condition func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) bool,
	value func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (types.String, error),
) planmodifier.String {
	return &conditionalValueModifier{
		condition: condition,
		value:     value,
	}
}

func (m *conditionalValueModifier) Description(ctx context.Context) string {
	return "Set value based on condition"
}

func (m *conditionalValueModifier) MarkdownDescription(ctx context.Context) string {
	return "Set value based on condition"
}

func (m *conditionalValueModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.PlanValue.IsNull() && m.condition(ctx, req.Plan, req.State) {
		newValue, err := m.value(ctx, req.Plan, req.State)
		if err != nil {
			resp.Diagnostics.AddError("Conditional Value Error", err.Error())
			return
		}
		resp.PlanValue = newValue
	}
}

// internal/provider/utils.go
package provider

import (
	"crypto/rand"
	"math/big"
	"strings"
	
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func generatePassword(length int) (types.String, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	
	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return types.StringNull(), err
		}
		password[i] = charset[num.Int64()]
	}
	
	return types.StringValue(string(password)), nil
}

func randInt(min, max int) int {
	num, _ := rand.Int(rand.Reader, big.NewInt(int64(max-min+1)))
	return min + int(num.Int64())
}

// internal/provider/server_resource_test.go
package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccServerResource_NullToValueTransformation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test avec valeurs null qui deviennent des valeurs par défaut
			{
				Config: testAccServerResourceConfigMinimal("test-server"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("custom_server.test", "name", "test-server"),
					// Environment null -> "development"
					resource.TestCheckResourceAttr("custom_server.test", "environment", "development"),
					// Hostname généré
					resource.TestMatchResourceAttr("custom_server.test", "hostname", regexp.MustCompile(`^server-\d+$`)),
					// IP générée pendant apply
					resource.TestMatchResourceAttr("custom_server.test", "ip_address", regexp.MustCompile(`^10\.0\.\d+\.\d+$`)),
					// Password généré
					resource.TestCheckResourceAttrSet("custom_server.test", "password"),
					// Tags par défaut créés
					resource.TestCheckResourceAttr("custom_server.test", "tags.managed_by", "terraform"),
					resource.TestCheckResourceAttrSet("custom_server.test", "tags.created_at"),
				),
			},
			// Test avec valeurs explicites (ne doivent pas être overridées)
			{
				Config: testAccServerResourceConfigComplete("test-server", "production", "custom-host", "192.168.1.100"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("custom_server.test", "environment", "production"),
					resource.TestCheckResourceAttr("custom_server.test", "hostname", "custom-host"),
					resource.TestCheckResourceAttr("custom_server.test", "ip_address", "192.168.1.100"),
				),
			},
		},
	})
}

func TestAccServerResource_PlanStability(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccServerResourceConfigMinimal("test-server"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("custom_server.test", "hostname"),
					resource.TestCheckResourceAttrSet("custom_server.test", "password"),
				),
			},
			// Re-apply - les valeurs générées doivent rester stables
			{
				Config:   testAccServerResourceConfigMinimal("test-server"),
				PlanOnly: true, // Aucun changement ne devrait être planifié
			},
		},
	})
}

func testAccServerResourceConfigMinimal(name string) string {
	return fmt.Sprintf(`
resource "custom_server" "test" {
  name = %[1]q
  # environment, hostname, ip_address, password et tags sont null
  # et seront transformés par les plan modifiers
}
`, name)
}

func testAccServerResourceConfigComplete(name, environment, hostname, ipAddress string) string {
	return fmt.Sprintf(`
resource "custom_server" "test" {
  name        = %[1]q
  environment = %[2]q
  hostname    = %[3]q
  ip_address  = %[4]q
  
  tags = {
    custom = "value"
  }
}
`, name, environment, hostname, ipAddress)
}
```

```hcl
// Exemple d'utilisation dans un fichier Terraform
// examples/main.tf

terraform {
  required_providers {
    custom = {
      source = "registry.terraform.io/example/custom"
    }
  }
}

provider "custom" {}

# Exemple avec transformation null -> valeurs par défaut
resource "custom_server" "example1" {
  name = "my-server"
  # environment sera "development"
  # hostname sera généré
  # ip_address sera généré pendant apply
  # password sera généré
  # tags seront créés par défaut
}

# Exemple avec valeurs explicites
resource "custom_server" "example2" {
  name        = "my-custom-server"
  environment = "production"
  hostname    = "prod-web-01"
  ip_address  = "10.0.1.50"
  
  tags = {
    role = "web-server"
    team = "platform"
  }
}

# Exemple de sortie pour voir les valeurs générées
output "generated_password" {
  value     = custom_server.example1.password
  sensitive = true
}

output "generated_ip" {
  value = custom_server.example1.ip_address
}
```
