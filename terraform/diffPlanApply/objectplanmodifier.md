Je vais vous créer un exemple avec un `ObjectPlanModifier` personnalisé qui transforme les valeurs null dans des objets complexes.Voici un exemple complet avec des `ObjectPlanModifier` personnalisés qui transforment les objets null ou partiels. L'exemple utilise une ressource de base de données avec trois objets complexes qui sont transformés différemment.

## ObjectPlanModifiers personnalisés créés :

### 1. **NullToDefaultObject**
```go
NullToDefaultObject(func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (types.Object, error) {
    // Génère un objet par défaut complet si l'objet est null
})
```

### 2. **FillNullFields** 
```go
FillNullFields(func(ctx context.Context, obj types.Object, plan tfsdk.Plan, state tfsdk.State) (types.Object, error) {
    // Remplit uniquement les champs null dans un objet existant
})
```

### 3. **ConditionalObjectTransform**
```go
ConditionalObjectTransform(
    condition, // fonction qui détermine si la transformation doit avoir lieu
    transform  // fonction de transformation
)
```

## Exemples de transformations :

### **Config object
... when the response to Claude.


```go
// internal/provider/custom_object_plan_modifiers.go
package provider

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

// ObjectPlanModifier personnalisé pour transformer les objets null
type nullToDefaultObjectModifier struct {
	defaultObject func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (types.Object, error)
}

func NullToDefaultObject(defaultObject func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (types.Object, error)) planmodifier.Object {
	return &nullToDefaultObjectModifier{
		defaultObject: defaultObject,
	}
}

func (m *nullToDefaultObjectModifier) Description(ctx context.Context) string {
	return "Set default object if null"
}

func (m *nullToDefaultObjectModifier) MarkdownDescription(ctx context.Context) string {
	return "Set default object if null"
}

func (m *nullToDefaultObjectModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Si l'objet planifié est null
	if req.PlanValue.IsNull() {
		defaultObj, err := m.defaultObject(ctx, req.Plan, req.State)
		if err != nil {
			resp.Diagnostics.AddError("Default Object Generation Error", err.Error())
			return
		}
		resp.PlanValue = defaultObj
	}
}

// ObjectPlanModifier pour remplir les champs null dans un objet existant
type fillNullFieldsModifier struct {
	fillFields func(ctx context.Context, obj types.Object, plan tfsdk.Plan, state tfsdk.State) (types.Object, error)
}

func FillNullFields(fillFields func(ctx context.Context, obj types.Object, plan tfsdk.Plan, state tfsdk.State) (types.Object, error)) planmodifier.Object {
	return &fillNullFieldsModifier{
		fillFields: fillFields,
	}
}

func (m *fillNullFieldsModifier) Description(ctx context.Context) string {
	return "Fill null fields in object"
}

func (m *fillNullFieldsModifier) MarkdownDescription(ctx context.Context) string {
	return "Fill null fields in object"
}

func (m *fillNullFieldsModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	// Si l'objet n'est pas null mais peut contenir des champs null
	if !req.PlanValue.IsNull() {
		filledObj, err := m.fillFields(ctx, req.PlanValue, req.Plan, req.State)
		if err != nil {
			resp.Diagnostics.AddError("Fill Fields Error", err.Error())
			return
		}
		resp.PlanValue = filledObj
	}
}

// ObjectPlanModifier conditionnel avancé
type conditionalObjectModifier struct {
	condition func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, current types.Object) bool
	transform func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, current types.Object) (types.Object, error)
}

func ConditionalObjectTransform(
	condition func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, current types.Object) bool,
	transform func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, current types.Object) (types.Object, error),
) planmodifier.Object {
	return &conditionalObjectModifier{
		condition: condition,
		transform: transform,
	}
}

func (m *conditionalObjectModifier) Description(ctx context.Context) string {
	return "Conditionally transform object"
}

func (m *conditionalObjectModifier) MarkdownDescription(ctx context.Context) string {
	return "Conditionally transform object"
}

func (m *conditionalObjectModifier) PlanModifyObject(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
	if m.condition(ctx, req.Plan, req.State, req.PlanValue) {
		transformedObj, err := m.transform(ctx, req.Plan, req.State, req.PlanValue)
		if err != nil {
			resp.Diagnostics.AddError("Object Transform Error", err.Error())
			return
		}
		resp.PlanValue = transformedObj
	}
}

// internal/provider/database_resource.go
package provider

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &DatabaseResource{}

func NewDatabaseResource() resource.Resource {
	return &DatabaseResource{}
}

type DatabaseResource struct{}

type DatabaseResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Name     types.String `tfsdk:"name"`
	Engine   types.String `tfsdk:"engine"`
	Config   types.Object `tfsdk:"config"`
	Security types.Object `tfsdk:"security"`
	Backup   types.Object `tfsdk:"backup"`
}

type DatabaseConfigModel struct {
	Host         types.String `tfsdk:"host"`
	Port         types.Int64  `tfsdk:"port"`
	Database     types.String `tfsdk:"database"`
	MaxConns     types.Int64  `tfsdk:"max_connections"`
	Timeout      types.Int64  `tfsdk:"timeout"`
	SSLMode      types.String `tfsdk:"ssl_mode"`
}

type SecurityConfigModel struct {
	Username        types.String `tfsdk:"username"`
	Password        types.String `tfsdk:"password"`
	EncryptionKey   types.String `tfsdk:"encryption_key"`
	AllowedIPs      types.List   `tfsdk:"allowed_ips"`
	RequireSSL      types.Bool   `tfsdk:"require_ssl"`
}

type BackupConfigModel struct {
	Enabled         types.Bool   `tfsdk:"enabled"`
	Schedule        types.String `tfsdk:"schedule"`
	RetentionDays   types.Int64  `tfsdk:"retention_days"`
	StorageLocation types.String `tfsdk:"storage_location"`
	EncryptBackups  types.Bool   `tfsdk:"encrypt_backups"`
}

func (r *DatabaseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database"
}

func (r *DatabaseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Database resource with complex object transformations",

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
			"engine": schema.StringAttribute{
				Required: true,
			},
			"config": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"host": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"port": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"database": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"max_connections": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"timeout": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"ssl_mode": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					// Si config est null, crée un objet par défaut
					NullToDefaultObject(func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State) (types.Object, error) {
						var data DatabaseResourceModel
						plan.Get(ctx, &data)
						
						return r.generateDefaultConfig(ctx, data.Name.ValueString(), data.Engine.ValueString())
					}),
					// Remplit les champs null dans un objet existant
					FillNullFields(func(ctx context.Context, obj types.Object, plan tfsdk.Plan, state tfsdk.State) (types.Object, error) {
						var data DatabaseResourceModel
						plan.Get(ctx, &data)
						
						return r.fillConfigNullFields(ctx, obj, data.Name.ValueString(), data.Engine.ValueString())
					}),
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"security": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"username": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"password": schema.StringAttribute{
						Optional:  true,
						Computed:  true,
						Sensitive: true,
					},
					"encryption_key": schema.StringAttribute{
						Optional:  true,
						Computed:  true,
						Sensitive: true,
					},
					"allowed_ips": schema.ListAttribute{
						Optional:    true,
						Computed:    true,
						ElementType: types.StringType,
					},
					"require_ssl": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					// Transformation conditionnelle basée sur l'engine
					ConditionalObjectTransform(
						func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, current types.Object) bool {
							var data DatabaseResourceModel
							plan.Get(ctx, &data)
							// Génère security uniquement pour postgres et mysql
							engine := data.Engine.ValueString()
							return current.IsNull() && (engine == "postgresql" || engine == "mysql")
						},
						func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, current types.Object) (types.Object, error) {
							var data DatabaseResourceModel
							plan.Get(ctx, &data)
							
							return r.generateSecurityConfig(ctx, data.Name.ValueString(), data.Engine.ValueString())
						},
					),
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"backup": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
					"schedule": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"retention_days": schema.Int64Attribute{
						Optional: true,
						Computed: true,
					},
					"storage_location": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"encrypt_backups": schema.BoolAttribute{
						Optional: true,
						Computed: true,
					},
				},
				PlanModifiers: []planmodifier.Object{
					// Backup par défaut si null et en production
					ConditionalObjectTransform(
						func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, current types.Object) bool {
							// Logique: créer backup par défaut si null et nom contient "prod"
							var data DatabaseResourceModel
							plan.Get(ctx, &data)
							return current.IsNull() && strings.Contains(strings.ToLower(data.Name.ValueString()), "prod")
						},
						func(ctx context.Context, plan tfsdk.Plan, state tfsdk.State, current types.Object) (types.Object, error) {
							return r.generateProductionBackupConfig(ctx)
						},
					),
					objectplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

// Méthodes helper pour la génération d'objets
func (r *DatabaseResource) generateDefaultConfig(ctx context.Context, name, engine string) (types.Object, error) {
	var port int64 = 5432 // PostgreSQL par défaut
	var sslMode string = "require"
	
	switch engine {
	case "mysql":
		port = 3306
		sslMode = "REQUIRED"
	case "mongodb":
		port = 27017
		sslMode = "requireSSL"
	}

	config := DatabaseConfigModel{
		Host:        types.StringValue(fmt.Sprintf("%s-db.internal", name)),
		Port:        types.Int64Value(port),
		Database:    types.StringValue(name),
		MaxConns:    types.Int64Value(100),
		Timeout:     types.Int64Value(30),
		SSLMode:     types.StringValue(sslMode),
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"host":            types.StringType,
		"port":            types.Int64Type,
		"database":        types.StringType,
		"max_connections": types.Int64Type,
		"timeout":         types.Int64Type,
		"ssl_mode":        types.StringType,
	}, config)
}

func (r *DatabaseResource) fillConfigNullFields(ctx context.Context, obj types.Object, name, engine string) (types.Object, error) {
	var config DatabaseConfigModel
	obj.As(ctx, &config, types.ObjectAsOptions{})

	// Remplit uniquement les champs null
	if config.Host.IsNull() {
		config.Host = types.StringValue(fmt.Sprintf("%s-db.internal", name))
	}
	if config.Port.IsNull() {
		var port int64 = 5432
		if engine == "mysql" {
			port = 3306
		} else if engine == "mongodb" {
			port = 27017
		}
		config.Port = types.Int64Value(port)
	}
	if config.Database.IsNull() {
		config.Database = types.StringValue(name)
	}
	if config.MaxConns.IsNull() {
		config.MaxConns = types.Int64Value(100)
	}
	if config.Timeout.IsNull() {
		config.Timeout = types.Int64Value(30)
	}
	if config.SSLMode.IsNull() {
		sslMode := "require"
		if engine == "mysql" {
			sslMode = "REQUIRED"
		} else if engine == "mongodb" {
			sslMode = "requireSSL"
		}
		config.SSLMode = types.StringValue(sslMode)
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"host":            types.StringType,
		"port":            types.Int64Type,
		"database":        types.StringType,
		"max_connections": types.Int64Type,
		"timeout":         types.Int64Type,
		"ssl_mode":        types.StringType,
	}, config)
}

func (r *DatabaseResource) generateSecurityConfig(ctx context.Context, name, engine string) (types.Object, error) {
	// Génération de mot de passe
	password, err := generateSecurePassword(32)
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}

	// Génération de clé de chiffrement
	encKey, err := generateEncryptionKey(32)
	if err != nil {
		return types.ObjectNull(map[string]attr.Type{}), err
	}

	// IPs par défaut
	defaultIPs, _ := types.ListValueFrom(ctx, types.StringType, []string{"10.0.0.0/8", "172.16.0.0/12"})

	security := SecurityConfigModel{
		Username:      types.StringValue(fmt.Sprintf("%s_user", name)),
		Password:      types.StringValue(password),
		EncryptionKey: types.StringValue(encKey),
		AllowedIPs:    defaultIPs,
		RequireSSL:    types.BoolValue(true),
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"username":       types.StringType,
		"password":       types.StringType,
		"encryption_key": types.StringType,
		"allowed_ips":    types.ListType{ElemType: types.StringType},
		"require_ssl":    types.BoolType,
	}, security)
}

func (r *DatabaseResource) generateProductionBackupConfig(ctx context.Context) (types.Object, error) {
	backup := BackupConfigModel{
		Enabled:         types.BoolValue(true),
		Schedule:        types.StringValue("0 2 * * *"), // 2h du matin
		RetentionDays:   types.Int64Value(30),
		StorageLocation: types.StringValue("s3://prod-db-backups/"),
		EncryptBackups:  types.BoolValue(true),
	}

	return types.ObjectValueFrom(ctx, map[string]attr.Type{
		"enabled":          types.BoolType,
		"schedule":         types.StringType,
		"retention_days":   types.Int64Type,
		"storage_location": types.StringType,
		"encrypt_backups":  types.BoolType,
	}, backup)
}

func (r *DatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Génération de l'ID
	data.ID = types.StringValue(fmt.Sprintf("db-%d", time.Now().Unix()))

	// Les ObjectPlanModifiers ont déjà fait le travail de transformation
	// Ici on simule juste la création dans l'API

	time.Sleep(200 * time.Millisecond) // Simulation d'appel API

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *DatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	time.Sleep(100 * time.Millisecond) // Simulation de suppression
}

// Fonctions utilitaires
func generateSecurePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	
	password := make([]byte, length)
	for i := range password {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[num.Int64()]
	}
	
	return string(password), nil
}

func generateEncryptionKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%x", bytes), nil
}

// internal/provider/database_resource_test.go
package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDatabaseResource_ObjectNullTransformation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test avec objets null qui deviennent des objets par défaut
			{
				Config: testAccDatabaseResourceConfigMinimal("test-db", "postgresql"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("custom_database.test", "name", "test-db"),
					resource.TestCheckResourceAttr("custom_database.test", "engine", "postgresql"),
					
					// Config généré par NullToDefaultObject
					resource.TestCheckResourceAttr("custom_database.test", "config.host", "test-db-db.internal"),
					resource.TestCheckResourceAttr("custom_database.test", "config.port", "5432"),
					resource.TestCheckResourceAttr("custom_database.test", "config.database", "test-db"),
					resource.TestCheckResourceAttr("custom_database.test", "config.max_connections", "100"),
					resource.TestCheckResourceAttr("custom_database.test", "config.ssl_mode", "require"),
					
					// Security généré par ConditionalObjectTransform
					resource.TestCheckResourceAttr("custom_database.test", "security.username", "test-db_user"),
					resource.TestCheckResourceAttrSet("custom_database.test", "security.password"),
					resource.TestCheckResourceAttr("custom_database.test", "security.require_ssl", "true"),
					
					// Backup null car pas "prod" dans le nom
					resource.TestCheckResourceAttr("custom_database.test", "backup.%", "0"),
				),
			},
		},
	})
}

func TestAccDatabaseResource_PartialObjectFill(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test avec config partielle - les champs null sont remplis
			{
				Config: testAccDatabaseResourceConfigPartial("my-app", "mysql"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Valeur explicite conservée
					resource.TestCheckResourceAttr("custom_database.test", "config.host", "custom-host.example.com"),
					// Valeurs null remplies par FillNullFields
					resource.TestCheckResourceAttr("custom_database.test", "config.port", "3306"),
					resource.TestCheckResourceAttr("custom_database.test", "config.database", "my-app"),
					resource.TestCheckResourceAttr("custom_database.test", "config.ssl_mode", "REQUIRED"),
				),
			},
		},
	})
}

func TestAccDatabaseResource_ProductionBackup(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			// Test que "prod" dans le nom génère backup automatiquement
			{
				Config: testAccDatabaseResourceConfigMinimal("prod-api-db", "postgresql"),
				Check: resource.ComposeAggregateTestCheckFunc(
					// Backup généré car "prod" dans le nom
					resource.TestCheckResourceAttr("custom_database.test", "backup.enabled", "true"),
					resource.TestCheckResourceAttr("custom_database.test", "backup.schedule", "0 2 * * *"),
					resource.TestCheckResourceAttr("custom_database.test", "backup.retention_days", "30"),
					resource.TestCheckResourceAttr("custom_database.test", "backup.encrypt_backups", "true"),
				),
			},
		},
	})
}

func TestAccDatabaseResource_PlanStability(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDatabaseResourceConfigMinimal("test-db", "postgresql"),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttrSet("custom_database.test", "security.password"),
				),
			},
			// Re-apply - les objets générés doivent rester stables
			{
				Config:   testAccDatabaseResourceConfigMinimal("test-db", "postgresql"),
				PlanOnly: true, // Aucun changement planifié
			},
		},
	})
}

func testAccDatabaseResourceConfigMinimal(name, engine string) string {
	return fmt.Sprintf(`
resource "custom_database" "test" {
  name   = %[1]q
  engine = %[2]q
  # config, security, backup sont null
  # et seront transformés par les ObjectPlanModifiers
}
`, name, engine)
}

func testAccDatabaseResourceConfigPartial(name, engine string) string {
	return fmt.Sprintf(`
resource "custom_database" "test" {
  name   = %[1]q
  engine = %[2]q
  
  config = {
    host = "custom-host.example.com"
    # port, database, max_connections, timeout, ssl_mode sont null
    # et seront remplis par FillNullFields
  }
}
`, name, engine)
}
```

```hcl
// Exemple d'utilisation Terraform
// examples/database.tf

terraform {
  required_providers {
    custom = {
      source = "registry.terraform.io/example/custom"
    }
  }
}

provider "custom" {}

# Base de données minimale - tout sera généré
resource "custom_database" "app_db" {
  name   = "my-app"
  engine = "postgresql"
  # config = null → objet par défaut généré
  # security = null → objet de sécurité généré (car postgresql)
  # backup = null → reste null (pas "prod" dans le nom)
}

# Base de données MySQL avec config partielle
resource "custom_database" "mysql_db" {
  name   = "analytics"
  engine = "mysql"
  
  config = {
    host = "mysql.internal.company.com"
    max_connections = 200
    # port, database, timeout, ssl_mode → générés automatiquement
  }
}

# Base de données de production - backup automatique
resource "custom_database" "prod_api" {
  name   = "prod-api-db"  # "prod" déclenche backup automatique
  engine = "postgresql"
  
  # Tout généré avec des valeurs de production appropriées
}

# Base de données avec configuration complète explicite
resource "custom_database" "custom_db" {
  name   = "special-db"
  engine = "mongodb"
  
  config = {
    host            = "mongo-cluster.example.com"
    port            = 27017
    database        = "special_data"
    max_connections = 500
    timeout         = 60
    ssl_mode        = "requireSSL"
  }
  
  security = {
    username    = "admin"
    require_ssl = true
    allowed_ips = ["10.0.0.0/16", "172.31.0.0/16"]
    # password, encryption_key → générés même si objet explicite
  }
  
  backup = {
    enabled          = true
    schedule         = "0 3 * * 0"  # Dimanche 3h
    retention_days   = 90
    storage_location = "s3://my-custom-backups/"
    encrypt_backups  = true
  }
}

# Outputs pour voir les valeurs générées
output "app_db_config" {
  value = custom_database.app_db.config
}

output "mysql_generated_values" {
  value = {
    port     = custom_database.mysql_db.config.port
    database = custom_database.mysql_db.config.database
    ssl_mode = custom_database.mysql_db.config.ssl_mode
  }
}

output "prod_backup_config" {
  value     = custom_database.prod_api.backup
  sensitive = false
}
```