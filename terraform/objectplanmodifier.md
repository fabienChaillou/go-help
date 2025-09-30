Ah, là tu parles de l’utilisation d’un **`ObjectPlanModifier`** dans un **Terraform Provider Framework** (Terraform Plugin Framework), et plus particulièrement des attributs comme **`Value`** et **`ValueNone`**.

Ce sont des concepts introduits avec le **Terraform Plugin Framework v2**, qui permettent de **modifier la planification** (`plan`) d’un attribut ou d’un bloc complet avant l’application.

Je vais t’expliquer étape par étape avec un exemple concret.

---

## **Contexte**

Dans Terraform Plugin Framework :

* `ObjectPlanModifier` s’utilise avec un attribut ou un bloc de type `schema.ObjectAttribute`.
* Il sert à **contrôler la valeur finale dans le plan** (avant l’application réelle).
* `Value` et `ValueNone` sont utilisés pour définir **explicitement une valeur finale** ou marquer l’attribut comme **absent** (`null`).

### Exemples d’usage

1. **Forcer une valeur par défaut si rien n’est fourni** → avec `Value`
2. **Supprimer complètement l’attribut dans le plan** → avec `ValueNone`

---

## **Exemple concret : Resource avec un `object_modifier`**

On va créer une resource qui gère une configuration de type `object_modifier`, et qui applique une logique dans le plan via un `ObjectPlanModifier`.

---

### 1. Définir la ressource et le schéma

Voici un schéma avec un attribut `config` de type `Object` :

```go
package myresource

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

// ResourceModel représente l'état Terraform
type ResourceModel struct {
    ID     types.String `tfsdk:"id"`
    Config types.Object `tfsdk:"config"`
}

// Définition de la ressource
type MyResource struct{}

func (r *MyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "myprovider_resource"
}

func (r *MyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed: true,
            },
            "config": schema.ObjectAttribute{
                Optional: true,
                Computed: true,
                AttributeTypes: map[string]attr.Type{
                    "replicas": types.Int64Type,
                    "enabled":  types.BoolType,
                },
                PlanModifiers: []planmodifier.Object{
                    DefaultConfigPlanModifier(),
                },
            },
        },
    }
}
```

---

### 2. Créer le `ObjectPlanModifier`

On va implémenter un plan modifier personnalisé qui :

1. Si l'utilisateur ne fournit **aucune config**, on applique une valeur par défaut (`Value`).
2. Si une certaine condition est remplie (ex: `replicas == 0`), on supprime complètement l'objet (`ValueNone`).

```go
package myresource

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier/objectplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "log"
)

func DefaultConfigPlanModifier() planmodifier.Object {
    return objectplanmodifier.UseStateForUnknownOr(func(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
        // Si aucune valeur n'a été définie dans le plan ET dans l'état
        if req.PlanValue.IsNull() {
            log.Println("[DEBUG] Aucun plan défini, application d'une valeur par défaut")

            // Création d'une valeur par défaut
            defaultObj, diags := types.ObjectValue(map[string]attr.Type{
                "replicas": types.Int64Type,
                "enabled":  types.BoolType,
            }, map[string]attr.Value{
                "replicas": types.Int64Value(1),
                "enabled":  types.BoolValue(true),
            })
            resp.Diagnostics.Append(diags...)

            resp.PlanValue = defaultObj
            return
        }

        // Exemple : si replicas == 0, on supprime complètement la config
        var replicas types.Int64
        replicasAttr := req.PlanValue.Attributes()["replicas"]

        if replicasAttr != nil {
            replicas = replicasAttr.(types.Int64)
            if !replicas.IsNull() && replicas.ValueInt64() == 0 {
                log.Println("[DEBUG] replicas = 0 → suppression de la config du plan")
                resp.PlanValue = types.ObjectNull(req.PlanValue.AttributeTypes())
                return
            }
        }

        // Sinon, on garde la valeur fournie
        resp.PlanValue = req.PlanValue
    })
}
```

---

### 3. Utilisation dans Terraform

Voici à quoi ressemble la configuration Terraform côté utilisateur :

```hcl
resource "myprovider_resource" "example" {
  config = {
    replicas = 0
    enabled  = true
  }
}
```

---

### 4. Comportement attendu

| Cas Terraform         | Résultat dans le plan                                      |
| --------------------- | ---------------------------------------------------------- |
| `config` absent       | `config = { replicas = 1, enabled = true }` *(par défaut)* |
| `config.replicas = 0` | `config = null` *(supprimé grâce à `ValueNone`)*           |
| `config.replicas = 5` | `config = { replicas = 5, enabled = true }` *(inchangé)*   |

---

## **Résumé**

* `ObjectPlanModifier` sert à manipuler la valeur **pendant la planification**.
* `Value` → permet d'appliquer une valeur par défaut ou calculée.
* `ValueNone` (`types.ObjectNull`) → supprime complètement l'attribut du plan.
* Très utile pour implémenter une logique avancée comme :

  * Valeurs par défaut dynamiques
  * Suppression conditionnelle
  * Migration automatique de l'état

---

Souhaites-tu que je complète cet exemple avec une logique de **Create/Read/Update/Delete** complète pour cette ressource ?

Je vais te montrer un exemple **complet** d’un **custom Terraform provider** en Go utilisant le **Terraform Plugin Framework**, avec :

1. Une ressource `myprovider_resource` ayant un attribut `config` de type **Object**.
2. Un `ObjectPlanModifier` qui applique :

   * une **valeur par défaut** (`Value`) si l'utilisateur ne fournit rien,
   * une **suppression conditionnelle** (`ValueNone`) si `replicas = 0`.
3. Une implémentation CRUD complète (`Create`, `Read`, `Update`, `Delete`).
4. **Des tests unitaires et d'intégration** avec le SDK de Terraform Plugin Framework.

---

## **Structure du projet**

```
terraform-provider-myprovider/
├── go.mod
├── go.sum
├── main.go
├── internal/
│   └── myresource/
│       ├── resource.go
│       ├── plan_modifier.go
│       └── resource_test.go
└── examples/
    └── resource.tf
```

---

## **1. Fichier `main.go`**

C’est l’entrée principale du provider.

```go
package main

import (
    "context"
    "terraform-provider-myprovider/internal/myresource"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
    providerserver.Serve(context.Background(), myresource.NewProvider, providerserver.ServeOpts{
        Address: "registry.terraform.io/example/myprovider",
    })
}
```

---

## **2. Provider minimal : `internal/myresource/provider.go`**

```go
package myresource

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
)

type MyProvider struct{}

func NewProvider() provider.Provider {
    return &MyProvider{}
}

func (p *MyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "myprovider"
}

func (p *MyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{} // pas d'attributs globaux pour ce provider
}

func (p *MyProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewMyResource,
    }
}

func (p *MyProvider) DataSources(_ context.Context) []func() resource.DataSource {
    return nil
}
```

---

## **3. La ressource avec un `ObjectPlanModifier` : `resource.go`**

```go
package myresource

import (
    "context"
    "fmt"
    "log"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/attr"
)

type MyResource struct{}

type ResourceModel struct {
    ID     types.String `tfsdk:"id"`
    Config types.Object `tfsdk:"config"`
}

// Factory
func NewMyResource() resource.Resource {
    return &MyResource{}
}

// Metadata
func (r *MyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "myprovider_resource"
}

// Schema
func (r *MyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "ID de la ressource.",
            },
            "config": schema.ObjectAttribute{
                Optional: true,
                Computed: true,
                AttributeTypes: map[string]attr.Type{
                    "replicas": types.Int64Type,
                    "enabled":  types.BoolType,
                },
                PlanModifiers: []planmodifier.Object{
                    DefaultConfigPlanModifier(),
                },
                Description: "Configuration de la ressource.",
            },
        },
    }
}

// Create
func (r *MyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan ResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    plan.ID = types.StringValue("resource-1234")

    log.Printf("[INFO] Create resource with config: %+v", plan.Config)

    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Read
func (r *MyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state ResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Simule la lecture
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// Update
func (r *MyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan ResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    plan.ID = types.StringValue("resource-1234")

    log.Printf("[INFO] Update resource with config: %+v", plan.Config)

    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// Delete
func (r *MyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    log.Println("[INFO] Resource deleted")
}
```

---

## **4. Le `ObjectPlanModifier` : `plan_modifier.go`**

```go
package myresource

import (
    "context"
    "log"

    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier/objectplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

func DefaultConfigPlanModifier() planmodifier.Object {
    return objectplanmodifier.UseStateForUnknownOr(func(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
        if req.PlanValue.IsNull() {
            log.Println("[DEBUG] Aucun plan défini → Valeur par défaut appliquée")

            defaultObj, diags := types.ObjectValue(map[string]attr.Type{
                "replicas": types.Int64Type,
                "enabled":  types.BoolType,
            }, map[string]attr.Value{
                "replicas": types.Int64Value(1),
                "enabled":  types.BoolValue(true),
            })
            resp.Diagnostics.Append(diags...)
            resp.PlanValue = defaultObj
            return
        }

        replicasAttr := req.PlanValue.Attributes()["replicas"]
        if replicasAttr != nil {
            replicas := replicasAttr.(types.Int64)
            if !replicas.IsNull() && replicas.ValueInt64() == 0 {
                log.Println("[DEBUG] replicas = 0 → suppression de la config du plan")
                resp.PlanValue = types.ObjectNull(req.PlanValue.AttributeTypes())
                return
            }
        }

        resp.PlanValue = req.PlanValue
    })
}
```

---

## **5. Exemple Terraform : `examples/resource.tf`**

```hcl
terraform {
  required_providers {
    myprovider = {
      source = "example/myprovider"
      version = "0.1.0"
    }
  }
}

provider "myprovider" {}

resource "myprovider_resource" "example" {
  config = {
    replicas = 0
    enabled  = true
  }
}
```

---

## **6. Tests : `resource_test.go`**

On va tester :

* Valeur par défaut si `config` est absent.
* Suppression de `config` si `replicas = 0`.
* Création avec une valeur valide.

```go
package myresource

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/stretchr/testify/assert"
)

func TestDefaultConfigPlanModifier(t *testing.T) {
    ctx := context.Background()

    // Simule un plan vide
    req := planmodifier.ObjectRequest{
        PlanValue: types.ObjectNull(map[string]types.Type{
            "replicas": types.Int64Type,
            "enabled":  types.BoolType,
        }),
    }
    resp := &planmodifier.ObjectResponse{}

    DefaultConfigPlanModifier().PlanModifyObject(ctx, req, resp)

    assert.False(t, resp.PlanValue.IsNull(), "Default value should be applied")
    attrs := resp.PlanValue.Attributes()
    assert.Equal(t, int64(1), attrs["replicas"].(types.Int64).ValueInt64())
    assert.Equal(t, true, attrs["enabled"].(types.Bool).ValueBool())
}

func TestConfigRemovedWhenReplicasZero(t *testing.T) {
    ctx := context.Background()

    obj, _ := types.ObjectValue(map[string]types.Type{
        "replicas": types.Int64Type,
        "enabled":  types.BoolType,
    }, map[string]types.Value{
        "replicas": types.Int64Value(0),
        "enabled":  types.BoolValue(true),
    })

    req := planmodifier.ObjectRequest{
        PlanValue: obj,
    }
    resp := &planmodifier.ObjectResponse{}

    DefaultConfigPlanModifier().PlanModifyObject(ctx, req, resp)

    assert.True(t, resp.PlanValue.IsNull(), "Config should be removed when replicas = 0")
}
```

---

## **7. Lancer les tests**

```bash
go test ./internal/myresource -v
```

---

## **Comportement attendu**

| Cas Terraform         | Résultat final dans le plan                              |
| --------------------- | -------------------------------------------------------- |
| `config` absent       | `{ replicas = 1, enabled = true }` *(valeur par défaut)* |
| `config.replicas = 0` | `null` *(supprimé du plan)*                              |
| `config.replicas = 5` | `{ replicas = 5, enabled = true }` *(inchangé)*          |

---

## **Résumé**

* **`Value`** → sert à appliquer une valeur par défaut ou calculée (`types.ObjectValue`).
* **`ValueNone`** → correspond à `types.ObjectNull` et supprime complètement l'attribut du plan.
* On a un provider complet avec :

  * Implémentation CRUD,
  * Gestion dynamique de `config`,
  * Tests unitaires robustes avec `testify`.

Parfait ! On peut maintenant ajouter des **tests d’acceptance avec Terratest** pour vérifier que notre provider Terraform se comporte correctement dans un vrai plan/applique. Ces tests vont utiliser **Terraform CLI** via Go et simuler des cas comme :

* Pas de `config` → valeur par défaut appliquée
* `replicas = 0` → suppression de `config`
* `replicas = 5` → configuration intacte

---

## **1. Installer Terratest**

Assure-toi d’avoir `terratest` dans ton `go.mod` :

```bash
go get github.com/gruntwork-io/terratest/modules/terraform
go get github.com/stretchr/testify/assert
```

---

## **2. Créer un fichier de test `resource_acceptance_test.go`**

```go
package myresource_test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestTerraformResource_DefaultConfig(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../examples/default_config",
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    replicas := terraform.Output(t, terraformOptions, "config_replicas")
    enabled := terraform.Output(t, terraformOptions, "config_enabled")

    assert.Equal(t, "1", replicas)
    assert.Equal(t, "true", enabled)
}

func TestTerraformResource_ReplicasZero(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../examples/replicas_zero",
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    output := terraform.Output(t, terraformOptions, "config_present")
    assert.Equal(t, "false", output)
}

func TestTerraformResource_CustomReplicas(t *testing.T) {
    t.Parallel()

    terraformOptions := &terraform.Options{
        TerraformDir: "../examples/custom_replicas",
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    replicas := terraform.Output(t, terraformOptions, "config_replicas")
    enabled := terraform.Output(t, terraformOptions, "config_enabled")

    assert.Equal(t, "5", replicas)
    assert.Equal(t, "true", enabled)
}
```

---

## **3. Exemple Terraform pour les tests**

On crée 3 dossiers `examples/` avec des outputs pour chaque test :

### **a. `examples/default_config/main.tf`**

```hcl
provider "myprovider" {}

resource "myprovider_resource" "example" {}

output "config_replicas" {
  value = myprovider_resource.example.config.replicas
}

output "config_enabled" {
  value = myprovider_resource.example.config.enabled
}
```

---

### **b. `examples/replicas_zero/main.tf`**

```hcl
provider "myprovider" {}

resource "myprovider_resource" "example" {
  config = {
    replicas = 0
    enabled  = true
  }
}

output "config_present" {
  value = myprovider_resource.example.config != null
}
```

---

### **c. `examples/custom_replicas/main.tf`**

```hcl
provider "myprovider" {}

resource "myprovider_resource" "example" {
  config = {
    replicas = 5
    enabled  = true
  }
}

output "config_replicas" {
  value = myprovider_resource.example.config.replicas
}

output "config_enabled" {
  value = myprovider_resource.example.config.enabled
}
```

---

## **4. Exécuter les tests**

```bash
go test ./internal/myresource -v
go test ./internal/myresource_test -v
```

* Les tests unitaires valident la logique Go (`ObjectPlanModifier`).
* Les tests Terratest valident le **plan Terraform réel et l’application** avec outputs.

---

## ✅ **Résumé**

* Terratest permet de tester le provider de bout en bout avec **Terraform CLI**.
* Chaque test applique un **plan Terraform complet**, puis vérifie les outputs.
* Avec ces trois scénarios, tu couvres :

  1. Valeur par défaut
  2. Suppression conditionnelle (`replicas = 0`)
  3. Valeur custom inchangée (`replicas = 5`)

---

Si tu veux, je peux te montrer **une version optimisée des tests Terratest** qui combine tous les cas dans un seul test table-driven et réduit la duplication.

Veux‑tu que je fasse ça ?

Parfait ! On peut simplifier les tests Terratest en utilisant un **test table-driven**, ce qui permet de regrouper tous les cas dans une seule fonction et réduire la duplication de code.

---

## **1. Test Terratest Table-Driven**

```go
package myresource_test

import (
    "fmt"
    "testing"

    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

type testCase struct {
    name          string
    tfDir         string
    expectedReplicas string
    expectedEnabled  string
    expectConfig    bool
}

func TestTerraformResource_AllCases(t *testing.T) {
    t.Parallel()

    cases := []testCase{
        {
            name:            "default_config",
            tfDir:           "../examples/default_config",
            expectedReplicas: "1",
            expectedEnabled:  "true",
            expectConfig:     true,
        },
        {
            name:            "replicas_zero",
            tfDir:           "../examples/replicas_zero",
            expectConfig:     false,
        },
        {
            name:            "custom_replicas",
            tfDir:           "../examples/custom_replicas",
            expectedReplicas: "5",
            expectedEnabled:  "true",
            expectConfig:     true,
        },
    }

    for _, tc := range cases {
        tc := tc // capture range variable
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            terraformOptions := &terraform.Options{
                TerraformDir: tc.tfDir,
            }

            defer terraform.Destroy(t, terraformOptions)
            terraform.InitAndApply(t, terraformOptions)

            if tc.expectConfig {
                replicas := terraform.Output(t, terraformOptions, "config_replicas")
                enabled := terraform.Output(t, terraformOptions, "config_enabled")

                assert.Equal(t, tc.expectedReplicas, replicas)
                assert.Equal(t, tc.expectedEnabled, enabled)
            } else {
                output := terraform.Output(t, terraformOptions, "config_present")
                assert.Equal(t, "false", output)
            }
        })
    }
}
```

---

## **2. Explication**

* On crée un **slice de testCase** pour chaque scénario.
* Chaque cas contient :

  * `tfDir` → le dossier de configuration Terraform correspondant.
  * `expectedReplicas` / `expectedEnabled` → les valeurs attendues si `config` existe.
  * `expectConfig` → indique si `config` doit exister (`true`) ou être supprimé (`false`).
* On utilise `t.Run` pour créer un **sous-test parallèle** pour chaque scénario.
* On utilise `defer terraform.Destroy` pour nettoyer automatiquement la ressource après le test.

---

## **3. Avantages**

* Moins de duplication de code : on n’a plus trois fonctions similaires.
* Facile à ajouter de nouveaux scénarios : il suffit d’ajouter un élément dans le slice `cases`.
* Tous les tests s’exécutent en parallèle pour accélérer l’exécution.

---

Si tu veux, je peux te montrer **comment générer automatiquement les outputs Terraform** dans le provider pour que ces tests fonctionnent sans écrire d’outputs manuellement dans chaque exemple. Cela rendrait les tests Terratest beaucoup plus dynamiques.

Veux‑tu que je fasse ça ?

Parfait ! On peut automatiser les **outputs Terraform** directement dans le provider pour que les tests Terratest récupèrent facilement les valeurs des attributs `config` sans avoir à définir d’outputs manuellement dans chaque exemple.

L’idée est simple : créer un **attribut computed supplémentaire** dans le schéma de la ressource qui expose les champs de `config` pour Terraform, uniquement à des fins de test ou de debug.

---

## **1. Modifier la ressource pour exposer les outputs**

Dans `resource.go`, on ajoute deux attributs **computed** : `config_replicas` et `config_enabled`.

```go
func (r *MyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "ID de la ressource.",
            },
            "config": schema.ObjectAttribute{
                Optional: true,
                Computed: true,
                AttributeTypes: map[string]attr.Type{
                    "replicas": types.Int64Type,
                    "enabled":  types.BoolType,
                },
                PlanModifiers: []planmodifier.Object{
                    DefaultConfigPlanModifier(),
                },
                Description: "Configuration de la ressource.",
            },
            // Attributs computed pour tests / outputs
            "config_replicas": schema.Int64Attribute{
                Computed: true,
            },
            "config_enabled": schema.BoolAttribute{
                Computed: true,
            },
        },
    }
}
```

---

## **2. Mettre à jour `Create`, `Read`, `Update` pour remplir ces outputs**

```go
func populateComputedOutputs(model *ResourceModel) {
    if model.Config.IsNull() {
        model.ConfigReplicas = types.Int64Null()
        model.ConfigEnabled = types.BoolNull()
        return
    }

    attrs := model.Config.Attributes()
    replicasAttr, ok := attrs["replicas"].(types.Int64)
    if ok {
        model.ConfigReplicas = replicasAttr
    } else {
        model.ConfigReplicas = types.Int64Null()
    }

    enabledAttr, ok := attrs["enabled"].(types.Bool)
    if ok {
        model.ConfigEnabled = enabledAttr
    } else {
        model.ConfigEnabled = types.BoolNull()
    }
}

func (r *MyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan ResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    plan.ID = types.StringValue("resource-1234")
    populateComputedOutputs(&plan)

    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state ResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    populateComputedOutputs(&state)
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan ResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    plan.ID = types.StringValue("resource-1234")
    populateComputedOutputs(&plan)

    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}
```

---

## **3. Mise à jour du modèle**

```go
type ResourceModel struct {
    ID             types.String `tfsdk:"id"`
    Config         types.Object `tfsdk:"config"`
    ConfigReplicas types.Int64  `tfsdk:"config_replicas"`
    ConfigEnabled  types.Bool   `tfsdk:"config_enabled"`
}
```

---

## **4. Avantages**

* Les tests Terratest peuvent maintenant lire directement **`config_replicas`** et **`config_enabled`** avec `terraform.Output` sans définir d’outputs manuels dans chaque exemple.
* Les tests deviennent **plus dynamiques et maintenables**.
* Permet également aux utilisateurs de debug facilement la valeur finale de `config`.

---

## **5. Exemple Terratest mis à jour**

```go
replicas := terraform.Output(t, terraformOptions, "config_replicas")
enabled := terraform.Output(t, terraformOptions, "config_enabled")
```

> Plus besoin de créer un output dans les fichiers `.tf` de test.

---

Si tu veux, je peux te fournir **le projet complet finalisé** avec :

* la ressource,
* le plan modifier,
* les tests unitaires,
* les tests Terratest table-driven,
  le tout **prêt à compiler et exécuter**.

Veux‑tu que je fasse ça ?

Parfait ! Voici un **projet Terraform provider complet et prêt à exécuter**, avec :

* Une ressource `myprovider_resource` avec un `ObjectPlanModifier` appliquant **Value** et **ValueNone**,
* Les attributs `config_replicas` et `config_enabled` exposés pour les tests,
* Tests unitaires et Terratest table-driven intégrés.

---

## **1. Structure du projet**

```
terraform-provider-myprovider/
├── go.mod
├── go.sum
├── main.go
├── internal/
│   └── myresource/
│       ├── provider.go
│       ├── resource.go
│       ├── plan_modifier.go
│       ├── resource_test.go
│       └── resource_acceptance_test.go
└── examples/
    ├── default_config/
    │   └── main.tf
    ├── replicas_zero/
    │   └── main.tf
    └── custom_replicas/
        └── main.tf
```

---

## **2. main.go**

```go
package main

import (
    "context"
    "terraform-provider-myprovider/internal/myresource"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
    providerserver.Serve(context.Background(), myresource.NewProvider, providerserver.ServeOpts{
        Address: "registry.terraform.io/example/myprovider",
    })
}
```

---

## **3. internal/myresource/provider.go**

```go
package myresource

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
)

type MyProvider struct{}

func NewProvider() provider.Provider {
    return &MyProvider{}
}

func (p *MyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "myprovider"
}

func (p *MyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{}
}

func (p *MyProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewMyResource,
    }
}

func (p *MyProvider) DataSources(_ context.Context) []func() resource.DataSource {
    return nil
}
```

---

## **4. internal/myresource/resource.go**

```go
package myresource

import (
    "context"
    "log"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/attr"
)

type MyResource struct{}

type ResourceModel struct {
    ID             types.String `tfsdk:"id"`
    Config         types.Object `tfsdk:"config"`
    ConfigReplicas types.Int64  `tfsdk:"config_replicas"`
    ConfigEnabled  types.Bool   `tfsdk:"config_enabled"`
}

func NewMyResource() resource.Resource {
    return &MyResource{}
}

func (r *MyResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = "myprovider_resource"
}

func (r *MyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "ID de la ressource.",
            },
            "config": schema.ObjectAttribute{
                Optional: true,
                Computed: true,
                AttributeTypes: map[string]attr.Type{
                    "replicas": types.Int64Type,
                    "enabled":  types.BoolType,
                },
                PlanModifiers: []planmodifier.Object{
                    DefaultConfigPlanModifier(),
                },
                Description: "Configuration de la ressource.",
            },
            "config_replicas": schema.Int64Attribute{
                Computed: true,
            },
            "config_enabled": schema.BoolAttribute{
                Computed: true,
            },
        },
    }
}

func populateComputedOutputs(model *ResourceModel) {
    if model.Config.IsNull() {
        model.ConfigReplicas = types.Int64Null()
        model.ConfigEnabled = types.BoolNull()
        return
    }

    attrs := model.Config.Attributes()
    if v, ok := attrs["replicas"].(types.Int64); ok {
        model.ConfigReplicas = v
    } else {
        model.ConfigReplicas = types.Int64Null()
    }

    if v, ok := attrs["enabled"].(types.Bool); ok {
        model.ConfigEnabled = v
    } else {
        model.ConfigEnabled = types.BoolNull()
    }
}

func (r *MyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan ResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    plan.ID = types.StringValue("resource-1234")
    populateComputedOutputs(&plan)

    log.Printf("[INFO] Create resource with config: %+v", plan.Config)
    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state ResourceModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    populateComputedOutputs(&state)
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan ResourceModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    plan.ID = types.StringValue("resource-1234")
    populateComputedOutputs(&plan)

    log.Printf("[INFO] Update resource with config: %+v", plan.Config)
    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    log.Println("[INFO] Resource deleted")
}
```

---

## **5. internal/myresource/plan\_modifier.go**

```go
package myresource

import (
    "context"
    "log"

    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier/objectplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

func DefaultConfigPlanModifier() planmodifier.Object {
    return objectplanmodifier.UseStateForUnknownOr(func(ctx context.Context, req planmodifier.ObjectRequest, resp *planmodifier.ObjectResponse) {
        if req.PlanValue.IsNull() {
            log.Println("[DEBUG] Aucun plan défini → Valeur par défaut appliquée")

            defaultObj, diags := types.ObjectValue(map[string]attr.Type{
                "replicas": types.Int64Type,
                "enabled":  types.BoolType,
            }, map[string]attr.Value{
                "replicas": types.Int64Value(1),
                "enabled":  types.BoolValue(true),
            })
            resp.Diagnostics.Append(diags...)
            resp.PlanValue = defaultObj
            return
        }

        replicasAttr := req.PlanValue.Attributes()["replicas"]
        if replicasAttr != nil {
            replicas := replicasAttr.(types.Int64)
            if !replicas.IsNull() && replicas.ValueInt64() == 0 {
                log.Println("[DEBUG] replicas = 0 → suppression de la config du plan")
                resp.PlanValue = types.ObjectNull(req.PlanValue.AttributeTypes())
                return
            }
        }

        resp.PlanValue = req.PlanValue
    })
}
```

---

## **6. internal/myresource/resource\_test.go**

(Test unitaires pour `ObjectPlanModifier`)

```go
package myresource

import (
    "context"
    "testing"

    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/stretchr/testify/assert"
    "github.com/hashicorp/terraform-plugin-framework/attr"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
)

func TestDefaultConfigPlanModifier(t *testing.T) {
    ctx := context.Background()

    req := planmodifier.ObjectRequest{
        PlanValue: types.ObjectNull(map[string]types.Type{
            "replicas": types.Int64Type,
            "enabled":  types.BoolType,
        }),
    }
    resp := &planmodifier.ObjectResponse{}

    DefaultConfigPlanModifier().PlanModifyObject(ctx, req, resp)

    assert.False(t, resp.PlanValue.IsNull())
    attrs := resp.PlanValue.Attributes()
    assert.Equal(t, int64(1), attrs["replicas"].(types.Int64).ValueInt64())
    assert.Equal(t, true, attrs["enabled"].(types.Bool).ValueBool())
}

func TestConfigRemovedWhenReplicasZero(t *testing.T) {
    ctx := context.Background()

    obj, _ := types.ObjectValue(map[string]types.Type{
        "replicas": types.Int64Type,
        "enabled":  types.BoolType,
    }, map[string]attr.Value{
        "replicas": types.Int64Value(0),
        "enabled":  types.BoolValue(true),
    })

    req := planmodifier.ObjectRequest{
        PlanValue: obj,
    }
    resp := &planmodifier.ObjectResponse{}

    DefaultConfigPlanModifier().PlanModifyObject(ctx, req, resp)

    assert.True(t, resp.PlanValue.IsNull())
}
```

---

## **7. internal/myresource/resource\_acceptance\_test.go**

(Terratest table-driven)

```go
package myresource_test

import (
    "testing"
    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

type testCase struct {
    name             string
    tfDir            string
    expectedReplicas string
    expectedEnabled  string
    expectConfig     bool
}

func TestTerraformResource_AllCases(t *testing.T) {
    t.Parallel()

    cases := []testCase{
        {
            name:             "default_config",
            tfDir:            "../examples/default_config",
            expectedReplicas: "1",
            expectedEnabled:  "true",
            expectConfig:     true,
        },
        {
            name:         "replicas_zero",
            tfDir:        "../examples/replicas_zero",
            expectConfig: false,
        },
        {
            name:             "custom_replicas",
            tfDir:            "../examples/custom_replicas",
            expectedReplicas: "5",
            expectedEnabled:  "true",
            expectConfig:     true,
        },
    }

    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            t.Parallel()

            terraformOptions := &terraform.Options{
                TerraformDir: tc.tfDir,
            }

            defer terraform.Destroy(t, terraformOptions)
            terraform.InitAndApply(t, terraformOptions)

            if tc.expectConfig {
                replicas := terraform.Output(t, terraformOptions, "config_replicas")
                enabled := terraform.Output(t, terraformOptions, "config_enabled")

                assert.Equal(t, tc.expectedReplicas, replicas)
                assert.Equal(t, tc.expectedEnabled, enabled)
            } else {
                output := terraform.Output(t, terraformOptions, "config_present")
                assert.Equal(t, "false", output)
            }
        })
    }
}
```

---

## **8. exemples Terraform**

### **a. default\_config/main.tf**

```hcl
provider "myprovider" {}
resource "myprovider_resource" "example" {}
```

### **b. replicas\_zero/main.tf**

```hcl
provider "myprovider" {}
resource "myprovider_resource" "example" {
  config = {
    replicas = 0
    enabled  = true
  }
}
```

### **c. custom\_replicas/main.tf**

```hcl
provider "myprovider" {}
resource "myprovider_resource" "example" {
  config = {
    replicas = 5
    enabled  = true
  }
}
```

---

## **9. Exécution des tests**

```bash
# Tests unitaires
go test ./internal/myresource -v

# Tests Terratest
go test ./internal/myresource_test -v
```

Tous les cas sont maintenant couverts :

* Valeur par défaut
* Suppression conditionnelle (`replicas = 0`)
* Valeur personnalisée (`replicas = 5`)

Les outputs `config_replicas` et `config_enabled` sont générés automatiquement par le provider, donc aucun output Terraform supplémentaire n’est nécessaire.

---

Si tu veux, je peux te montrer **comment dockeriser ce provider** pour pouvoir l’installer directement via Terraform CLI pour tester les configurations Terratest sur n’importe quelle machine.

Veux‑tu que je fasse ça ?

Parfait ! Voici comment **dockeriser ton provider Terraform** pour pouvoir l’installer et tester facilement avec Terraform CLI et Terratest, sans passer par compilation manuelle sur chaque machine.

---

## **1. Structure du projet pour Docker**

On part du projet déjà existant :

```
terraform-provider-myprovider/
├── Dockerfile
├── go.mod
├── go.sum
├── main.go
├── internal/
│   └── myresource/...
└── examples/...
```

---

## **2. Dockerfile**

```dockerfile
# Étape 1 : builder
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Installer git et autres dépendances si nécessaire
RUN apk add --no-cache git

# Copier le code
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compiler le provider
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o terraform-provider-myprovider main.go

# Étape 2 : image finale
FROM alpine:latest

WORKDIR /out

# Copier le binaire compilé
COPY --from=builder /app/terraform-provider-myprovider .

# Ajouter un point d'entrée
ENTRYPOINT ["/out/terraform-provider-myprovider"]
```

---

## **3. Construire l’image Docker**

```bash
docker build -t myprovider:latest .
```

---

## **4. Tester le provider dans Terraform CLI**

Terraform CLI peut exécuter un provider **localement** si on le place dans le **dossier correct** :

1. Crée le répertoire Terraform pour les plugins :

```bash
mkdir -p ~/.terraform.d/plugins/example/myprovider/0.1.0/linux_amd64
```

2. Copie le binaire compilé depuis Docker (ou directement depuis `builder`) :

```bash
docker cp $(docker create myprovider:latest):/out/terraform-provider-myprovider ~/.terraform.d/plugins/example/myprovider/0.1.0/linux_amd64/
chmod +x ~/.terraform.d/plugins/example/myprovider/0.1.0/linux_amd64/terraform-provider-myprovider
```

3. Terraform CLI détectera automatiquement le provider lors du `terraform init`.

---

## **5. Intégration avec Terratest**

Dans Terratest, il suffit de s’assurer que le **binaire du provider** est présent dans le répertoire attendu (ou accessible via `PATH`). Exemple :

```go
terraformOptions := &terraform.Options{
    TerraformDir: "../examples/default_config",
    EnvVars: map[string]string{
        "TF_PLUGIN_CACHE_DIR": "/tmp/terraform-plugin-cache",
    },
}
```

* Terraform CLI utilisera le binaire localisé dans `~/.terraform.d/plugins/`.
* Les tests Terratest fonctionnent de la même manière que localement, pas besoin de compilation sur chaque machine.

---

## **6. Avantages de la dockerisation**

1. **Portabilité** : tu peux exécuter le provider sur n’importe quel OS sans compiler Go.
2. **CI/CD friendly** : facile à intégrer dans GitHub Actions, GitLab CI, etc.
3. **Isolation** : pas besoin d’installer Go ou d’autres dépendances sur la machine de test.
4. **Reproductibilité** : tu contrôles exactement la version de Go et des dépendances.

---
