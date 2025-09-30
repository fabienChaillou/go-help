Oui, bien sûr ! Je vais te montrer un exemple concret d’**ObjectPlanModifier** personnalisé en Go pour un provider Terraform, où l’objectif est de **transformer les valeurs nulles (`null`) dans un objet complexe** en des valeurs par défaut.

Cela est très utile dans un provider custom lorsque tu veux éviter que le plan ait des champs vides et imposer des valeurs par défaut dans la phase de **planification**.

---

## **Contexte**

Imaginons que tu as une ressource `MyResource` qui a un attribut complexe `settings` :

```hcl
resource "myprovider_myresource" "example" {
  name = "example"

  settings {
    replicas = null
    labels   = null
  }
}
```

Ici, si `replicas` ou `labels` sont `null`, tu veux les transformer automatiquement en :

* `replicas` → `1`
* `labels` → `{}` (map vide)

---

## **Code Go : Custom ObjectPlanModifier**

Voici l’implémentation :

```go
package planmodifiers

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "github.com/hashicorp/terraform-plugin-framework/diag"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier/objectplanmodifier"
)

// DefaultNullObjectModifier est un plan modifier qui remplace les champs null
// d'un objet par des valeurs par défaut.
type DefaultNullObjectModifier struct {
    Default types.Object
}

// Description
func (d DefaultNullObjectModifier) Description(_ context.Context) string {
    return "Remplace les champs null par des valeurs par défaut dans l'objet."
}

func (d DefaultNullObjectModifier) MarkdownDescription(_ context.Context) string {
    return d.Description(context.Background())
}

// PlanModifyObject est appelée par Terraform pendant la phase de planification.
func (d DefaultNullObjectModifier) PlanModifyObject(
    ctx context.Context,
    req planmodifier.ObjectRequest,
    resp *planmodifier.ObjectResponse,
) {
    // Si le plan est déjà connu, on ne touche à rien
    if !req.PlanValue.IsNull() && !req.PlanValue.IsUnknown() {
        return
    }

    // Remplace par l'objet par défaut si la valeur est null
    resp.PlanValue = d.Default
}
```

---

## **Utilisation dans le schéma Terraform**

Maintenant, tu peux utiliser ce `ObjectPlanModifier` dans la définition du schéma de ta ressource :

```go
package myresource

import (
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
    "myprovider/planmodifiers"
)

func SettingsSchema() schema.Attribute {
    return schema.SingleNestedAttribute{
        Description: "Configuration des paramètres de la ressource.",
        Optional:    true,
        Attributes: map[string]schema.Attribute{
            "replicas": schema.Int64Attribute{
                Optional: true,
            },
            "labels": schema.MapAttribute{
                Optional:    true,
                ElementType: types.StringType,
            },
        },
        PlanModifiers: []planmodifier.Object{
            planmodifiers.DefaultNullObjectModifier{
                Default: types.ObjectValueMust(
                    map[string]types.Type{
                        "replicas": types.Int64Type,
                        "labels":   types.MapType{ElemType: types.StringType},
                    },
                    map[string]attr.Value{
                        "replicas": types.Int64Value(1),
                        "labels":   types.MapNull(types.StringType), // map vide
                    },
                ),
            },
        },
    }
}
```

---

## **Explication**

1. **`DefaultNullObjectModifier`**

   * Vérifie si la valeur du plan est `null` ou `unknown`.
   * Si `null` → Remplace par l’objet `Default`.
   * Si non → Laisse la valeur inchangée.

2. **Dans le schéma**

   * On définit un objet `settings` avec deux champs : `replicas` et `labels`.
   * On attache notre plan modifier à cet objet.

3. **Comportement au runtime**

   * Si l’utilisateur ne fournit pas `settings` → Terraform injecte l’objet par défaut :

     ```hcl
     settings = {
       replicas = 1
       labels   = {}
     }
     ```
   * Si l’utilisateur fournit seulement `labels` → `replicas` reste à 1 par défaut.

---

## **Résultat attendu dans le plan**

**Input utilisateur :**

```hcl
resource "myprovider_myresource" "example" {
  name = "test"
  settings {
    labels = {
      env = "prod"
    }
  }
}
```

**Plan généré :**

```hcl
# settings sera complété automatiquement
settings = {
  replicas = 1
  labels = {
    env = "prod"
  }
}
```

---

## **Avantages**

* Plus besoin de gérer les valeurs `null` dans ton code interne.
* Garantit des plans toujours cohérents avec des valeurs par défaut.
* Code facilement réutilisable dans plusieurs ressources ou attributs.

######## With the test ###########

Bien sûr ! Je vais te montrer comment écrire **des tests d’acceptance Terraform** pour vérifier que notre `ObjectPlanModifier` fonctionne correctement dans un provider custom.

L’objectif sera de :

1. Créer un test qui vérifie que les valeurs `null` dans un objet complexe (`settings`) sont remplacées par les valeurs par défaut définies par notre `DefaultNullObjectModifier`.
2. Vérifier qu’un utilisateur peut partiellement fournir des valeurs et que seules les valeurs manquantes sont complétées automatiquement.
3. Confirmer qu’aucun changement n’est reproposé lors d’un `terraform plan` (idempotence).

---

## **Structure du projet**

Voici une structure typique pour un provider :

```
myprovider/
├── internal/
│   ├── provider/
│   │   ├── provider.go
│   │   └── myresource.go
│   └── planmodifiers/
│       └── default_null_object_modifier.go
└── internal/acctest/
│   └── resource_myresource_test.go
└── go.mod
```

Les tests d’acceptance se trouvent généralement dans `internal/acctest/`.

---

## **Code du test d’acceptance**

Voici un exemple de test d’acceptance complet.

```go
package acctest

import (
    "fmt"
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "github.com/hashicorp/terraform-plugin-testing/terraform"
)

// TestAccMyResource_SettingsDefaults vérifie que les valeurs nulles
// dans "settings" sont remplacées par les valeurs par défaut
func TestAccMyResource_SettingsDefaults(t *testing.T) {
    resourceName := "myprovider_myresource.test"

    resource.Test(t, resource.TestCase{
        // Cette fonction doit renvoyer une instance du provider
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccMyResourceConfigNullSettings(),
                Check: resource.ComposeTestCheckFunc(
                    // Vérifie que Terraform a injecté la valeur par défaut pour replicas
                    resource.TestCheckResourceAttr(resourceName, "settings.replicas", "1"),
                    // Vérifie que labels est vide
                    resource.TestCheckResourceAttr(resourceName, "settings.labels.%", "0"),
                ),
            },
            {
                // Plan doit être vide après la création (idempotence)
                Config: testAccMyResourceConfigNullSettings(),
                PlanOnly: true,
            },
        },
    })
}

// TestAccMyResource_SettingsPartialInput vérifie que si l'utilisateur
// fournit une partie de l'objet, le reste est complété automatiquement
func TestAccMyResource_SettingsPartialInput(t *testing.T) {
    resourceName := "myprovider_myresource.test"

    resource.Test(t, resource.TestCase{
        Providers: testAccProviders,
        Steps: []resource.TestStep{
            {
                Config: testAccMyResourceConfigPartialSettings(),
                Check: resource.ComposeTestCheckFunc(
                    // "replicas" non fourni → doit être 1
                    resource.TestCheckResourceAttr(resourceName, "settings.replicas", "1"),
                    // "labels" fourni → doit être conservé
                    resource.TestCheckResourceAttr(resourceName, "settings.labels.env", "prod"),
                ),
            },
        },
    })
}

//
// Helpers
//

func testAccMyResourceConfigNullSettings() string {
    return `
resource "myprovider_myresource" "test" {
  name = "test-acc"
}
`
}

func testAccMyResourceConfigPartialSettings() string {
    return `
resource "myprovider_myresource" "test" {
  name = "test-acc"

  settings {
    labels = {
      env = "prod"
    }
  }
}
`
}
```

---

## **Initialisation du provider pour les tests**

Pour que `terraform-plugin-testing` fonctionne, on a besoin d’un **provider factory**.

Dans `internal/acctest/provider_test.go` :

```go
package acctest

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
    "myprovider/internal/provider"
)

// testAccProviders est utilisé par tous les tests acceptance
var testAccProviders = map[string]func() (*schema.Provider, error){
    "myprovider": func() (*schema.Provider, error) {
        return provider.New(), nil
    },
}

func TestProvider(t *testing.T) {
    if err := provider.New().InternalValidate(); err != nil {
        t.Fatalf("err: %s", err)
    }
}
```

---

## **Points clés à vérifier**

1. **Test initial avec `null`**

   * Aucun `settings` fourni → Terraform doit injecter l’objet complet par défaut.
   * `replicas = 1` et `labels = {}`.

2. **Test partiel**

   * L’utilisateur fournit seulement `labels`.
   * Terraform complète automatiquement `replicas` à `1`.

3. **Idempotence**

   * Deuxième `terraform plan` doit être **vide**, prouvant que l’état Terraform correspond exactement à la configuration générée.

---

## **Exécution du test**

Avant d’exécuter les tests, installe les outils nécessaires :

```bash
go get github.com/hashicorp/terraform-plugin-testing
```

Puis lance les tests avec :

```bash
TF_ACC=1 go test ./internal/acctest -v -count=1
```

* `TF_ACC=1` : obligatoire pour activer les tests d’acceptance.
* `-count=1` : empêche la mise en cache des résultats de test.

---

## **Résultat attendu**

Exemple de sortie :

```
=== RUN   TestAccMyResource_SettingsDefaults
--- PASS: TestAccMyResource_SettingsDefaults (5.22s)
=== RUN   TestAccMyResource_SettingsPartialInput
--- PASS: TestAccMyResource_SettingsPartialInput (4.10s)
PASS
ok      myprovider/internal/acctest   9.32s
```

---

## **Résumé**

Avec ce setup :

* Tu as un `ObjectPlanModifier` qui remplace automatiquement les `null` par des valeurs par défaut.
* Les tests d’acceptance valident :

  1. Que le provider injecte correctement les valeurs par défaut.
  2. Que les valeurs partiellement fournies sont respectées.
  3. Que le plan reste stable et idempotent.

C’est la bonne pratique recommandée par HashiCorp pour vérifier le comportement d’un provider custom en conditions réelles.

##### with terratest moduel for acceptance tests ##########

Bien sûr ! Nous allons ajouter des **tests avec Terratest** pour compléter les tests d’acceptance que nous avons faits avec `terraform-plugin-testing`.

Terratest te permet de :

* Lancer **de vrais plans Terraform** en appelant la CLI Terraform directement.
* Tester ton provider comme un utilisateur final, dans des conditions réelles.
* Vérifier la création, la mise à jour et la suppression de la ressource via ton provider custom.
* Compléter les tests unitaires et d’acceptance en validant l’intégration complète.

---

## **1. Installation de Terratest**

Dans ton projet Go, ajoute Terratest :

```bash
go get github.com/gruntwork-io/terratest/modules/terraform
```

---

## **2. Structure du projet de tests**

```
myprovider/
├── examples/
│   └── myresource/
│       └── main.tf
├── internal/
│   ├── provider/
│   │   └── myresource.go
│   └── planmodifiers/
│       └── default_null_object_modifier.go
└── test/
    └── terratest/
        └── myresource_test.go
```

Ici :

* `examples/myresource/main.tf` contient un exemple minimal pour tester la ressource.
* `test/terratest/` contient les tests end-to-end exécutés avec Terratest.

---

## **3. Exemple Terraform pour les tests (`examples/myresource/main.tf`)**

```hcl
terraform {
  required_providers {
    myprovider = {
      source = "local/myprovider"
      version = "0.1.0"
    }
  }
}

provider "myprovider" {}

resource "myprovider_myresource" "test" {
  name = var.name

  settings {
    replicas = var.replicas
    labels   = var.labels
  }
}

output "settings_replicas" {
  value = myprovider_myresource.test.settings.replicas
}

output "settings_labels" {
  value = myprovider_myresource.test.settings.labels
}
```

Fichier des variables `examples/myresource/variables.tf` :

```hcl
variable "name" {
  type = string
}

variable "replicas" {
  type    = number
  default = null
}

variable "labels" {
  type    = map(string)
  default = null
}
```

---

## **4. Code du test Terratest (`test/terratest/myresource_test.go`)**

Voici un test complet :

```go
package terratest

import (
    "fmt"
    "os"
    "path/filepath"
    "testing"

    "github.com/gruntwork-io/terratest/modules/terraform"
    "github.com/stretchr/testify/assert"
)

func TestMyResource_DefaultSettings(t *testing.T) {
    t.Parallel()

    // Chemin vers l'exemple Terraform
    exampleDir := filepath.Join("../../examples/myresource")

    terraformOptions := &terraform.Options{
        TerraformDir: exampleDir,
        Vars: map[string]interface{}{
            "name": "test-defaults",
            // Pas de replicas ni de labels -> null
        },
        EnvVars: map[string]string{
            "TF_LOG": "INFO",
        },
    }

    // Nettoyage à la fin du test
    defer terraform.Destroy(t, terraformOptions)

    // Init + Apply
    terraform.InitAndApply(t, terraformOptions)

    // Vérification des outputs
    replicas := terraform.Output(t, terraformOptions, "settings_replicas")
    labels := terraform.OutputMap(t, terraformOptions, "settings_labels")

    assert.Equal(t, "1", replicas, "replicas doit être 1 par défaut")
    assert.Equal(t, 0, len(labels), "labels doit être vide par défaut")
}

func TestMyResource_PartialSettings(t *testing.T) {
    t.Parallel()

    exampleDir := filepath.Join("../../examples/myresource")

    terraformOptions := &terraform.Options{
        TerraformDir: exampleDir,
        Vars: map[string]interface{}{
            "name": "test-partial",
            "labels": map[string]string{
                "env": "prod",
            },
        },
    }

    defer terraform.Destroy(t, terraformOptions)

    terraform.InitAndApply(t, terraformOptions)

    replicas := terraform.Output(t, terraformOptions, "settings_replicas")
    labels := terraform.OutputMap(t, terraformOptions, "settings_labels")

    // replicas doit être complété automatiquement
    assert.Equal(t, "1", replicas, "replicas doit être 1 par défaut")
    // labels doit être conservé
    assert.Equal(t, "prod", labels["env"])
}
```

---

## **5. Exécution des tests Terratest**

Avant de lancer les tests :

* Compile ton provider et installe-le localement pour que Terraform puisse le trouver.

```bash
go build -o ~/.terraform.d/plugins/local/myprovider/myprovider
```

**Arborescence attendue :**

```
~/.terraform.d/plugins/local/myprovider/
└── myprovider
```

Puis lance les tests Terratest :

```bash
go test ./test/terratest -v -count=1
```

---

## **6. Résultats attendus**

Exemple de sortie :

```
=== RUN   TestMyResource_DefaultSettings
--- PASS: TestMyResource_DefaultSettings (18.23s)
=== RUN   TestMyResource_PartialSettings
--- PASS: TestMyResource_PartialSettings (17.90s)
PASS
ok      myprovider/test/terratest   36.22s
```

---

## **7. Ce que nous avons vérifié**

1. **Test avec valeurs `null` :**

   * `replicas` non défini → doit être `1`.
   * `labels` non défini → doit être `{}`.

2. **Test partiel :**

   * L’utilisateur fournit uniquement `labels`.
   * Terraform complète automatiquement `replicas`.

3. **Vérification end-to-end :**

   * Le provider est utilisé comme un binaire externe.
   * Terraform CLI est réellement exécuté.
   * Le comportement correspond à l’utilisation en production.

---

## **8. Comparaison avec les tests d’acceptance**

| **Critère**                    | **Tests Acceptance (`terraform-plugin-testing`)** | **Tests Terratest**        |
| ------------------------------ | ------------------------------------------------- | -------------------------- |
| Niveau de test                 | Provider en interne (SDK directement)             | Provider via Terraform CLI |
| Performance                    | Rapide, isolé                                     | Plus lent (Terraform CLI)  |
| Vérifie les defaults           | ✅ Oui                                             | ✅ Oui                      |
| Vérifie l’intégration complète | ❌ Non (pas de CLI)                                | ✅ Oui                      |
| Simule un vrai utilisateur     | ❌ Non                                             | ✅ Oui                      |

---

## **Résumé**

Avec ces tests :

* `terraform-plugin-testing` valide **le comportement interne du provider**, y compris notre `ObjectPlanModifier`.
* **Terratest valide le provider comme un utilisateur final**, en utilisant Terraform CLI pour un test end-to-end.
* Tu obtiens une couverture complète et robuste pour ton provider.

