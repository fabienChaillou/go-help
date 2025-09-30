Voici un exemple complet qui illustre :

1. **Un provider Terraform custom** en Go.
2. L'utilisation de `PlanModifier` pour modifier un champ lors de la phase `plan`.
3. **Des tests d'acceptance** pour valider le comportement du provider avec `terraform-plugin-testing`.

Cet exemple est basé sur le SDK Terraform Plugin Framework (`terraform-plugin-framework`).

---

## **1. Structure du projet**

```
terraform-provider-example/
├── go.mod
├── go.sum
├── main.go
├── internal/
│   └── provider/
│       ├── provider.go
│       ├── resource_example.go
│       └── plan_modifier.go
└── internal/
    └── acceptance/
        └── resource_example_test.go
```

---

## **2. Implémentation du PlanModifier**

Dans Terraform, un `PlanModifier` te permet d'ajuster la valeur d'un champ **au moment du plan** avant l'application.

Par exemple, ici nous allons automatiquement générer un `slug` à partir du nom.

### `internal/provider/plan_modifier.go`

```go
package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// CustomSlugPlanModifier - génère un slug à partir du champ "name"
type CustomSlugPlanModifier struct{}

func (m CustomSlugPlanModifier) Description(ctx context.Context) string {
	return "Génère un slug basé sur le champ 'name'."
}

func (m CustomSlugPlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Génère un slug basé sur le champ 'name'."
}

func (m CustomSlugPlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	// Génère un slug simple
	name := req.ConfigValue.ValueString()
	slug := strings.ToLower(strings.ReplaceAll(name, " ", "-"))
	resp.PlanValue = types.StringValue(fmt.Sprintf("%s-slug", slug))
}

// Factory pour l'utiliser facilement dans le schéma
func SlugPlanModifier() planmodifier.String {
	return CustomSlugPlanModifier{}
}
```

---

## **3. Ressource Terraform**

Une ressource qui prend un `name` en entrée et génère un `slug` automatiquement.

### `internal/provider/resource_example.go`

```go
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ExampleResource struct{}

func NewExampleResource() resource.Resource {
	return &ExampleResource{}
}

type ExampleResourceModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
	Slug types.String `tfsdk:"slug"`
}

func (r *ExampleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_example"
}

func (r *ExampleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"slug": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					SlugPlanModifier(),
				},
			},
		},
	}
}

func (r *ExampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan ExampleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.ID = types.StringValue("example-id-123")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ExampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state ExampleResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.Set(ctx, state)
}

func (r *ExampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan ExampleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.Set(ctx, plan)
}

func (r *ExampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
```

---

## **4. Provider**

### `internal/provider/provider.go`

```go
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

type ExampleProvider struct{}

func New() provider.Provider {
	return &ExampleProvider{}
}

func (p *ExampleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "example"
}

func (p *ExampleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *ExampleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
	}
}

func (p *ExampleProvider) DataSources(_ context.Context) []func() resource.DataSource {
	return nil
}
```

---

## **5. main.go**

```go
package main

import (
	"context"
	"log"

	"terraform-provider-example/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/example/provider",
	})
	if err != nil {
		log.Fatal(err)
	}
}
```

---

## **6. Test d'acceptance**

On utilise `terraform-plugin-testing` pour lancer des vrais tests avec Terraform.

### `internal/acceptance/resource_example_test.go`

```go
package acceptance

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() {},
		ProtoV6ProviderFactories: map[string]func() (interface{}, error){
			"example": func() (interface{}, error) {
				return provider.New(), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
resource "example_example" "test" {
  name = "My Test Name"
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("example_example.test", "name", "My Test Name"),
					resource.TestCheckResourceAttr("example_example.test", "slug", "my-test-name-slug"),
				),
			},
		},
	})
}
```

---

## **7. Lancer les tests**

```bash
go test ./internal/acceptance -v
```

Terraform va être lancé en mode test et vérifier le comportement du provider.

---

## **Résumé**

* **PlanModifier** est défini dans `plan_modifier.go` et permet de calculer automatiquement le champ `slug`.
* **Ressource** avec `name` et `slug`, le slug est généré au `plan`.
* **Tests d'acceptance** validant que le comportement est correct via `terraform-plugin-testing`.

Ce setup est une base propre pour créer un provider custom avec génération dynamique de valeurs pendant la phase de plan.

Voici comment **ajouter des tests avec Terratest** pour ton provider Terraform.

Nous allons :

1. Ajouter une configuration Terraform de test.
2. Créer un test en Go qui utilise **Terratest** (`github.com/gruntwork-io/terratest`).
3. Exécuter Terraform via Terratest et vérifier que le champ `slug` est bien généré par le `PlanModifier`.

---

## **1. Structure mise à jour**

On ajoute un répertoire `test/terratest` :

```
terraform-provider-example/
├── internal/
│   ├── provider/
│   │   ├── provider.go
│   │   ├── resource_example.go
│   │   └── plan_modifier.go
│   └── acceptance/
│       └── resource_example_test.go
├── test/
│   └── terratest/
│       ├── main.tf
│       └── terratest_test.go
├── main.go
├── go.mod
└── go.sum
```

---

## **2. Configuration Terraform pour le test**

Crée un fichier `test/terratest/main.tf` qui utilise le provider local.

```hcl
terraform {
  required_providers {
    example = {
      source = "registry.terraform.io/example/provider"
    }
  }
}

provider "example" {}

resource "example_example" "test" {
  name = "Terratest Name"
}

output "slug" {
  value = example_example.test.slug
}
```

> **Remarque** :
> Assure-toi que ton provider est compilé et disponible en local dans le répertoire Terraform plugin cache :
>
> ```
> mkdir -p ~/.terraform.d/plugins/registry.terraform.io/example/provider/1.0.0/darwin_amd64
> go build -o ~/.terraform.d/plugins/registry.terraform.io/example/provider/1.0.0/darwin_amd64/terraform-provider-example
> ```

Adapte `darwin_amd64` selon ton OS (`linux_amd64` ou `windows_amd64`).

---

## **3. Test Terratest**

Voici un test Terratest en Go qui :

* Lance Terraform,
* Applique la config,
* Vérifie que la sortie `slug` est correcte,
* Puis détruit l'infrastructure.

### `test/terratest/terratest_test.go`

```go
package terratest

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformExampleResource(t *testing.T) {
	t.Parallel()

	// Chemin vers le dossier contenant la configuration Terraform
	terraformDir, err := filepath.Abs("./")
	if err != nil {
		t.Fatalf("Erreur pour trouver le chemin absolu : %v", err)
	}

	// Options Terraform pour Terratest
	terraformOptions := &terraform.Options{
		TerraformDir: terraformDir,
		NoColor:      true,
	}

	// Initialisation et application de Terraform
	defer terraform.Destroy(t, terraformOptions)
	terraform.InitAndApply(t, terraformOptions)

	// Récupération de la sortie 'slug'
	slug := terraform.Output(t, terraformOptions, "slug")

	fmt.Printf("Slug généré : %s\n", slug)

	// Vérifie que le slug correspond bien au format attendu
	assert.True(t, strings.HasSuffix(slug, "-slug"), "Le slug doit se terminer par '-slug'")
	assert.Equal(t, "terratest-name-slug", slug, "Le slug doit être correctement généré")
}
```

---

## **4. Installer Terratest et Testify**

Ajoute les dépendances au `go.mod` :

```bash
go get github.com/gruntwork-io/terratest/modules/terraform
go get github.com/stretchr/testify
```

---

## **5. Lancer le test Terratest**

```bash
cd test/terratest
go test -v
```

Si tout fonctionne, la sortie devrait être similaire à :

```
=== RUN   TestTerraformExampleResource
Slug généré : terratest-name-slug
--- PASS: TestTerraformExampleResource (6.54s)
PASS
ok      terraform-provider-example/test/terratest    7.023s
```

---

## **6. Résumé des tests disponibles**

| Type de test               | Objectif                                                                                 | Commande                           |
| -------------------------- | ---------------------------------------------------------------------------------------- | ---------------------------------- |
| **Acceptance Tests**       | Teste le provider via `terraform-plugin-testing`, vérifie la logique interne du provider | `go test ./internal/acceptance -v` |
| **Terratest (End-to-End)** | Teste la configuration Terraform avec le provider compilé, comme dans un vrai usage      | `cd test/terratest && go test -v`  |

---

## **Points clés**

* **Acceptance Tests** → vérifient le provider **sans Terraform complet**, en simulant un plan et apply.
* **Terratest** → vérifie **de bout en bout**, comme un utilisateur final, en exécutant Terraform avec la CLI.
* Le `PlanModifier` est vérifié dans les deux cas :

  * Dans les acceptance tests, via `TestCheckResourceAttr`.
  * Dans Terratest, via la sortie `terraform.Output`.

Ainsi, tu as un pipeline de tests complet couvrant la logique interne et le comportement réel du provider.
