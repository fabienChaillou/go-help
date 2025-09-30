Pour tester un **provider Terraform custom** en local, tu as plusieurs approches, en fonction du niveau de test que tu veux réaliser :

* **Tests unitaires** → Vérifier la logique interne de ton provider avec Go (`go test`).
* **Tests d'intégration** → Lancer Terraform en local avec ton provider compilé.
* **Tests automatisés** → Utiliser le framework officiel `terraform-plugin-testing`.

Voici un guide complet étape par étape.

---

## **1. Prérequis**

Assure-toi d’avoir :

* **Go** (≥ 1.21)
* **Terraform CLI** (≥ 1.6)
* Le SDK Terraform (`terraform-plugin-sdk/v2` ou `terraform-plugin-framework`)
* Ton provider dans un repo Go (par ex. `github.com/ton-org/terraform-provider-demo`).

---

## **2. Compiler ton provider en local**

Dans le répertoire du provider :

```bash
go mod tidy
go build -o terraform-provider-demo
```

Vérifie que l'exécutable fonctionne :

```bash
./terraform-provider-demo
```

Il devrait afficher une erreur du type :

```
This binary is a Terraform plugin
```

---

## **3. Déclarer ton provider dans Terraform**

Terraform recherche les providers dans un répertoire spécifique.
Sur Linux/Mac :

```bash
~/.terraform.d/plugins/<namespace>/<type>/<version>/<os>_<arch>/terraform-provider-<name>_v<version>
```

Exemple pour un provider `demo` en version `0.1.0` :

```bash
mkdir -p ~/.terraform.d/plugins/local/demo/demo/0.1.0/darwin_amd64
mv terraform-provider-demo ~/.terraform.d/plugins/local/demo/demo/0.1.0/darwin_amd64/
```

Sur Windows :

```
%APPDATA%\terraform.d\plugins\local\demo\demo\0.1.0\windows_amd64\
```

---

## **4. Configurer Terraform pour utiliser le provider local**

Dans un projet Terraform :

**`main.tf`**

```hcl
terraform {
  required_providers {
    demo = {
      source  = "local/demo/demo"
      version = "0.1.0"
    }
  }
}

provider "demo" {
  # Configuration éventuelle de ton provider
  endpoint = "http://localhost:8080"
}
```

Puis :

```bash
terraform init
```

Tu devrais voir :

```
- Finding local/demo/demo versions matching "0.1.0"...
- Installing local/demo/demo v0.1.0...
- Installed local/demo/demo v0.1.0 (unauthenticated)
```

---

## **5. Tester ton provider avec Terraform**

Maintenant tu peux tester ton provider comme un utilisateur final :

```bash
terraform plan
terraform apply
```

Vérifie que ton provider réagit correctement.
C'est la première étape d'un test manuel.

---

## **6. Tests unitaires en Go**

Pour tester la logique interne (sans Terraform) :

**Exemple :**

```go
package demo

import "testing"

func TestHello(t *testing.T) {
    result := Hello("world")
    if result != "Hello world" {
        t.Errorf("Expected 'Hello world', got '%s'", result)
    }
}
```

Lancer les tests :

```bash
go test ./...
```

Avec la couverture :

```bash
go test ./... -cover
```

---

## **7. Tests d’intégration avec terraform-plugin-testing**

Terraform fournit un framework pour tester un provider comme si Terraform l'exécutait.

Ajoute la dépendance :

```bash
go get github.com/hashicorp/terraform-plugin-testing/helper/resource
```

**Exemple de test d'intégration :**

```go
package demo_test

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDemoResource_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Steps: []resource.TestStep{
            {
                Config: `
                provider "demo" {}
                resource "demo_item" "test" {
                  name = "example"
                }
                `,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("demo_item.test", "name", "example"),
                ),
            },
        },
    })
}
```

Lance les tests :

```bash
TF_ACC=1 go test ./... -v
```

> **Note :**
> `TF_ACC=1` active les tests d'intégration car ils font appel à Terraform en réel.

---

## **8. Déboguer ton provider**

Pour voir les logs détaillés de Terraform :

```bash
TF_LOG=DEBUG terraform apply
```

Pour déboguer avec Go :

```bash
dlv exec ./terraform-provider-demo
```

Ou lancer Terraform avec le débogueur attaché à ton provider.

---

## **Récapitulatif**

| Étape                        | Objectif                              | Commande clé                                                               |
| ---------------------------- | ------------------------------------- | -------------------------------------------------------------------------- |
| **Build provider**           | Compiler en Go                        | `go build -o terraform-provider-demo`                                      |
| **Installer provider local** | Copier dans `~/.terraform.d/plugins/` | `mv terraform-provider-demo ~/.terraform.d/plugins/local/demo/demo/0.1.0/` |
| **Init Terraform**           | Détecter le provider                  | `terraform init`                                                           |
| **Tester manuellement**      | Valider avec plan/apply               | `terraform apply`                                                          |
| **Unit tests**               | Tester la logique Go                  | `go test ./... -cover`                                                     |
| **Integration tests**        | Tester Terraform + provider           | `TF_ACC=1 go test ./...`                                                   |
| **Debug logs**               | Voir les logs Terraform               | `TF_LOG=DEBUG terraform apply`                                             |

---

Veux-tu que je te prépare un **exemple complet de provider minimal** avec un `resource` et un `data source`, prêt à être testé ?

Voici un **exemple complet** d’un provider Terraform minimal écrit en Go, incluant :

* Un **provider** nommé `demo`.
* Une **resource** `demo_item`.
* Un **data source** `demo_item_info`.
* Un setup pour le tester **en local** avec `terraform init` et `terraform apply`.
* Un **test d'intégration** avec `terraform-plugin-testing`.

Je vais te donner la structure du projet, le code, puis les commandes pour l'exécuter.

---

## **1. Structure du projet**

```
terraform-provider-demo/
├── go.mod
├── go.sum
├── main.go
├── provider.go
├── resource_item.go
├── data_source_item.go
└── internal/
    └── tests/
        └── provider_test.go
```

---

## **2. Fichier `go.mod`**

```go
module github.com/ton-org/terraform-provider-demo

go 1.22

require (
    github.com/hashicorp/terraform-plugin-framework v1.10.0
    github.com/hashicorp/terraform-plugin-testing v1.8.0
)
```

Puis :

```bash
go mod tidy
```

---

## **3. Fichier `main.go`**

C'est le point d'entrée du provider.

```go
package main

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/ton-org/terraform-provider-demo"
)

func main() {
    providerserver.Serve(context.Background(), demo.New, providerserver.ServeOpts{
        Address: "local/demo/demo",
    })
}
```

---

## **4. Fichier `provider.go`**

Définit le provider et ses paramètres globaux.

```go
package demo

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
)

type DemoProvider struct{}

var _ provider.Provider = &DemoProvider{}

func New() provider.Provider {
    return &DemoProvider{}
}

func (p *DemoProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "demo"
}

func (p *DemoProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "endpoint": schema.StringAttribute{
                Optional:    true,
                Description: "Endpoint de l'API externe",
            },
        },
    }
}

func (p *DemoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    // Ici tu pourrais initialiser une connexion à une API
}

func (p *DemoProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewItemResource,
    }
}

func (p *DemoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        NewItemDataSource,
    }
}
```

---

## **5. Fichier `resource_item.go`**

Une ressource simple qui stocke un `name` et retourne un `id` généré.

```go
package demo

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type ItemResource struct{}

var _ resource.Resource = &ItemResource{}

func NewItemResource() resource.Resource {
    return &ItemResource{}
}

type ItemModel struct {
    ID   types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
}

func (r *ItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_item"
}

func (r *ItemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "ID de l'item généré automatiquement",
            },
            "name": schema.StringAttribute{
                Required:    true,
                Description: "Nom de l'item",
            },
        },
    }
}

func (r *ItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan ItemModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    plan.ID = types.StringValue(fmt.Sprintf("item-%s", plan.Name.ValueString()))

    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state ItemModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Ici, on suppose que l'item existe toujours.
    resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (r *ItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan ItemModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    plan.ID = types.StringValue(fmt.Sprintf("item-%s", plan.Name.ValueString()))
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ItemResource) Delete(_ context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
    // Rien à faire, juste supprimer du state
}
```

---

## **6. Fichier `data_source_item.go`**

Un data source qui récupère des infos sur un item fictif.

```go
package demo

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type ItemDataSource struct{}

var _ datasource.DataSource = &ItemDataSource{}

func NewItemDataSource() datasource.DataSource {
    return &ItemDataSource{}
}

type ItemDataSourceModel struct {
    ID   types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
}

func (d *ItemDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_item_info"
}

func (d *ItemDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Required:    true,
                Description: "ID de l'item",
            },
            "name": schema.StringAttribute{
                Computed:    true,
                Description: "Nom de l'item récupéré",
            },
        },
    }
}

func (d *ItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ItemDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    data.Name = types.StringValue(fmt.Sprintf("Nom de %s", data.ID.ValueString()))
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

---

## **7. Test d'intégration (`internal/tests/provider_test.go`)**

```go
package tests

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccItemResource_basic(t *testing.T) {
    resource.Test(t, resource.TestCase{
        Steps: []resource.TestStep{
            {
                Config: `
                provider "demo" {}

                resource "demo_item" "test" {
                  name = "example"
                }
                `,
                Check: resource.ComposeTestCheckFunc(
                    resource.TestCheckResourceAttr("demo_item.test", "name", "example"),
                ),
            },
        },
    })
}
```

---

## **8. Compilation et installation du provider local**

```bash
go build -o terraform-provider-demo
mkdir -p ~/.terraform.d/plugins/local/demo/demo/0.1.0/darwin_amd64
mv terraform-provider-demo ~/.terraform.d/plugins/local/demo/demo/0.1.0/darwin_amd64/
```

*(Remplace `darwin_amd64` par `linux_amd64` ou `windows_amd64` selon ton OS.)*

---

## **9. Exemple Terraform (`main.tf`)**

```hcl
terraform {
  required_providers {
    demo = {
      source  = "local/demo/demo"
      version = "0.1.0"
    }
  }
}

provider "demo" {}

resource "demo_item" "test" {
  name = "mon-item"
}

data "demo_item_info" "test_info" {
  id = demo_item.test.id
}

output "item_id" {
  value = demo_item.test.id
}

output "item_name" {
  value = data.demo_item_info.test_info.name
}
```

---

## **10. Lancer le test en local**

Initialiser Terraform :

```bash
terraform init
```

Vérifier le plan :

```bash
terraform plan
```

Appliquer :

```bash
terraform apply -auto-approve
```

Tu devrais voir :

```
Outputs:

item_id = "item-mon-item"
item_name = "Nom de item-mon-item"
```

---

## **11. Tests unitaires et d'intégration**

Unit tests :

```bash
go test ./... -cover
```

Test Terraform d'intégration :

```bash
TF_ACC=1 go test ./internal/tests -v
```

---

## **Récapitulatif**

| Étape                  | Commande clé                                            |
| ---------------------- | ------------------------------------------------------- |
| Compiler le provider   | `go build -o terraform-provider-demo`                   |
| L'installer localement | `mv terraform-provider-demo ~/.terraform.d/plugins/...` |
| Initialiser Terraform  | `terraform init`                                        |
| Tester manuellement    | `terraform apply`                                       |
| Tests unitaires Go     | `go test ./... -cover`                                  |
| Tests Terraform        | `TF_ACC=1 go test ./internal/tests -v`                  |

---

Veux-tu que je prépare un **Makefile** pour automatiser la compilation, l’installation et les tests ?

Voici un **Makefile complet** qui automatise :

* La **compilation** du provider Terraform.
* Son **installation dans le répertoire local des plugins**.
* L’exécution des **tests unitaires** et **tests d’intégration Terraform**.
* Le **nettoyage** des artefacts.

---

## **Makefile**

```makefile
# Nom et version du provider
PROVIDER_NAME=demo
PROVIDER_NAMESPACE=local
PROVIDER_VERSION=0.1.0

# GOOS et GOARCH automatiques
OS=$(shell uname | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')

# Dossier d'installation pour Terraform
PLUGIN_DIR=$(HOME)/.terraform.d/plugins/$(PROVIDER_NAMESPACE)/$(PROVIDER_NAME)/$(PROVIDER_NAME)/$(PROVIDER_VERSION)/$(OS)_$(ARCH)

# Nom de l'exécutable généré
BINARY=terraform-provider-$(PROVIDER_NAME)

# Options Go
GO_FILES=$(shell find . -type f -name '*.go')
GO_FLAGS=-ldflags="-s -w"

# Fichiers Terraform pour tests locaux
TF_DIR=examples

# Commande Terraform
TF=terraform

# ========================================
# Cibles principales
# ========================================

.PHONY: all build install init apply test clean test-acc

all: build install

## Build du provider
build:
	@echo "==> Compilation du provider pour $(OS)/$(ARCH)..."
	GOOS=$(OS) GOARCH=$(ARCH) go build $(GO_FLAGS) -o $(BINARY) .

## Installation locale du provider
install: build
	@echo "==> Installation du provider dans $(PLUGIN_DIR)"
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY) $(PLUGIN_DIR)/$(BINARY)

## Initialiser un projet Terraform de test
init: install
	@echo "==> Initialisation de Terraform"
	cd $(TF_DIR) && $(TF) init

## Appliquer la configuration Terraform
apply: init
	@echo "==> Application de la configuration Terraform"
	cd $(TF_DIR) && $(TF) apply -auto-approve

## Nettoyage des fichiers générés
clean:
	@echo "==> Nettoyage..."
	rm -f $(BINARY)
	rm -rf $(TF_DIR)/.terraform $(TF_DIR)/terraform.tfstate*

## Lancer les tests unitaires Go
test:
	@echo "==> Lancement des tests unitaires Go"
	go test ./... -cover

## Lancer les tests d'intégration Terraform (acceptance tests)
test-acc: install
	@echo "==> Lancement des tests d'intégration Terraform"
	TF_ACC=1 go test ./internal/tests -v -timeout 10m
```

---

## **2. Dossier d'exemple Terraform (`examples/main.tf`)**

Ajoute un dossier `examples` à la racine avec un fichier `main.tf` :

```hcl
terraform {
  required_providers {
    demo = {
      source  = "local/demo/demo"
      version = "0.1.0"
    }
  }
}

provider "demo" {}

resource "demo_item" "test" {
  name = "mon-item"
}

data "demo_item_info" "test_info" {
  id = demo_item.test.id
}

output "item_id" {
  value = demo_item.test.id
}

output "item_name" {
  value = data.demo_item_info.test_info.name
}
```

---

## **3. Utilisation du Makefile**

### **Compilation et installation**

```bash
make install
```

Cela va :

1. Compiler le binaire pour ton OS/architecture.
2. L'installer dans `~/.terraform.d/plugins/local/demo/demo/0.1.0/<os>_<arch>/`.

---

### **Initialiser et appliquer Terraform**

```bash
make apply
```

Sortie attendue :

```
Outputs:

item_id = "item-mon-item"
item_name = "Nom de item-mon-item"
```

---

### **Tests unitaires**

```bash
make test
```

---

### **Tests d'intégration Terraform**

```bash
make test-acc
```

> ⚠️ N'oublie pas de définir `TF_ACC=1` dans la cible, ce qui est déjà géré par le Makefile.

---

### **Nettoyage complet**

```bash
make clean
```

---

## **4. Récapitulatif des commandes**

| Commande        | Action                                   |
| --------------- | ---------------------------------------- |
| `make build`    | Compile uniquement le provider           |
| `make install`  | Compile + installe le provider local     |
| `make init`     | `terraform init` dans `examples/`        |
| `make apply`    | `terraform apply` dans `examples/`       |
| `make test`     | Lancer les tests unitaires Go            |
| `make test-acc` | Lancer les tests d'intégration Terraform |
| `make clean`    | Supprimer les fichiers générés           |

---

## **5. Améliorations possibles**

* Ajouter la **cross-compilation** pour Linux, Windows et macOS :

  ```bash
  make build OS=linux ARCH=amd64
  ```
* Générer une archive `.zip` pour distribution.
* Intégrer ce Makefile dans une **GitHub Action CI/CD** pour build + tests automatiques.
* Publier le provider sur **Terraform Registry** (changement du namespace + signature des releases).

