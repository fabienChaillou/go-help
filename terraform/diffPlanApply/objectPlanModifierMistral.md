Voici un exemple de **provider custom Terraform** (en Go) où tu peux implémenter la logique pour modifier un objet de plan **uniquement s'il est `null`**. Ce pattern est souvent utilisé dans les providers Terraform pour gérer les valeurs par défaut ou les transformations conditionnelles.

### Exemple de code (Go) dans un provider custom

Supposons que tu as un schéma avec un attribut `optional_field` que tu veux modifier dans le plan si sa valeur est `null`.

```go
package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Exemple de schéma pour une ressource
type myResource struct{}

func (r *myResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"optional_field": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
		},
	}
}

// Méthode ModifyPlan pour modifier le plan
func (r *myResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	var plan myResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Si optional_field est null, on le modifie dans le plan
	if plan.OptionalField.IsNull() {
		plan.OptionalField = types.StringValue("valeur_par_defaut")
		resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
	}
}

// Modèle de données pour la ressource
type myResourceModel struct {
	OptionalField types.String `tfsdk:"optional_field"`
}
```

### Explication

- **`ModifyPlan`** est la méthode où tu peux intercepter et modifier le plan avant qu'il ne soit appliqué.
- **`plan.OptionalField.IsNull()`** vérifie si la valeur est `null`.
- Si c'est le cas, tu peux la remplacer par une valeur par défaut ou une autre logique métier.
- Enfin, tu mets à jour le plan avec `resp.Plan.Set(ctx, &plan)`.

---

#### acc test #####
Voici un exemple de **test d’acceptance** pour un provider Terraform custom, utilisant le framework `terraform-plugin-testing` (Go), afin de vérifier que ton `ModifyPlan` modifie bien un champ `null` dans le plan.

On va tester une ressource fictive `my_resource` avec un attribut `optional_field`. Le test vérifie que si `optional_field` est `null` dans la configuration, il est bien remplacé par une valeur par défaut dans le plan.

---

### 1. Structure de base du test

Supposons que ton provider s’appelle `myprovider` et que ta ressource s’appelle `my_resource`.

#### Fichier : `my_resource_test.go`

```go
package myprovider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccMyResource_ModifyPlan(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
				resource "myprovider_my_resource" "test" {
					// optional_field n'est pas défini, donc null
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(
							"myprovider_my_resource.test",
							tfjsonpath.New("optional_field"),
							knownvalue.StringExact("valeur_par_defaut"),
						),
					},
				},
			},
		},
	})
}
```

---

### 2. Explication du test

- **`Config`** : On déclare la ressource sans spécifier `optional_field`, donc il sera `null` dans le plan initial.
- **`ConfigPlanChecks`** : On utilise `PreApply` pour vérifier le plan **avant** l’application.
- **`ExpectKnownValue`** : On vérifie que `optional_field` a bien la valeur `"valeur_par_defaut"` dans le plan, grâce à ton `ModifyPlan`.

---

### 3. Prérequis

- Assure-toi d’avoir importé les packages nécessaires :
  ```go
  import (
      "testing"
      "github.com/hashicorp/terraform-plugin-testing/helper/resource"
      "github.com/hashicorp/terraform-plugin-testing/knownvalue"
      "github.com/hashicorp/terraform-plugin-testing/plancheck"
      "github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
  )
  ```
- `testAccProtoV6ProviderFactories` doit être défini pour instancier ton provider pendant le test.

---

### 4. Exemple complet avec setup

Si tu utilises le framework Terraform Plugin v2+ :

```go
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"myprovider": func() (tfprotov6.ProviderServer, error) {
		return New()(), nil // New() est la fonction qui crée ton provider
	},
}
```

---

### 5. Exécution du test

Lance le test avec :
```bash
go test -v ./... -run TestAccMyResource_ModifyPlan
```

---

### 6. Points clés

- Le test vérifie que le plan est modifié **avant** l’application.
- Si `optional_field` reste `null`, le test échouera.
- Si `ModifyPlan` fonctionne, le test passe.

---
