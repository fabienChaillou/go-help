Créer un **provider Terraform en Go** consiste à implémenter un plugin Terraform qui permet à Terraform de gérer des ressources spécifiques (ex : API custom, infrastructure non supportée officiellement, etc.). Terraform utilise le SDK Terraform Plugin Framework pour cela. Voici les grandes étapes :

---

## 🧱 1. Prérequis

* Go installé (v1.18+ recommandé)
* Terraform CLI installé
* Connaissance de base de Go et Terraform
* `git`, `make`, etc. (facultatif mais utile)

tout le code d'un provider terraform repose sur [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework/tree/main)


---

## 📦 2. Initialisation du projet

Utilise l’outil de scaffolding de HashiCorp :

```bash
git clone https://github.com/hashicorp/terraform-provider-scaffolding.git terraform-provider-myprovider
cd terraform-provider-myprovider
```

Puis renomme les références (`terraform-provider-scaffolding`) dans tous les fichiers (`go.mod`, `main.go`, etc.) en `terraform-provider-myprovider`.

---

## 📚 3. Structure générale d’un provider

Un provider Terraform en Go expose deux types principaux :

* **Provider** : contient la config globale, les ressources et les datasources.
* **Resources / DataSources** : définissent les actions CRUD.

Structure typique :

```
terraform-provider-myprovider/
├── main.go
├── internal/
│   └── provider/
│       ├── provider.go
│       └── resource_*.go
```

---

## 🛠️ 4. Exemple de provider minimal

### `main.go`

```go
package main

import (
	"context"
	"terraform-provider-myprovider/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/mycorp/myprovider",
	})
}
```

---

### `provider.go`

```go
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type myProvider struct{}

func New() provider.Provider {
	return &myProvider{}
}

func (p *myProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "myprovider"
}

func (p *myProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *myProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Lecture des attributs de config ici (ex: token)
}
```

---

### `resource_example.go`

```go
package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type exampleResource struct{}

func (r *exampleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Logique pour créer la ressource ici
}

func (r *exampleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Logique pour lire l'état réel
}

func (r *exampleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Logique pour mise à jour
}

func (r *exampleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Logique pour suppression
}
```

---

## 🧪 5. Compilation & utilisation

```bash
go mod tidy
go build -o terraform-provider-myprovider
```

Place le binaire dans :

```
~/.terraform.d/plugins/mycorp/myprovider/1.0.0/linux_amd64/
```

et crée un fichier `versions.tf` :

```hcl
terraform {
  required_providers {
    myprovider = {
      source = "mycorp/myprovider"
      version = "1.0.0"
    }
  }
}
```

---

## 🧰 6. Test avec un `main.tf`

```hcl
provider "myprovider" {
  token = "xyz"
}

resource "myprovider_example" "test" {
  name = "demo"
}
```

---

## 🚀 Aller plus loin

* [Terraform Plugin Framework](https://developer.hashicorp.com/terraform/plugin/framework)
* Test avec `terraform-plugin-testing`
* Génération de documentation automatique
* Publication sur le [Terraform Registry](https://registry.terraform.io/)

---

Souhaites-tu un exemple plus complet avec une ressource CRUD simulant une API REST ou autre ?

Top, on va faire un **exemple complet** d’un provider Terraform en Go qui gère une ressource simple : un objet `item` avec un `name` et un `description` en mémoire (pas d’API réelle, pour rester simple).

---

# Exemple complet d’un provider Terraform basique en Go

## 1. Structure minimale

```
terraform-provider-simple/
├── main.go
└── internal/
    └── provider/
        ├── provider.go
        └── resource_item.go
```

---

## 2. main.go

```go
package main

import (
	"context"
	"terraform-provider-simple/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/example/simple",
	})
}
```

---

## 3. internal/provider/provider.go

```go
package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

type simpleProvider struct{}

// New retourne une instance du provider
func New() provider.Provider {
	return &simpleProvider{}
}

// Metadata définit le nom du provider dans Terraform
func (p *simpleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "simple"
}

// Schema définit la config du provider (ici aucune config)
func (p *simpleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure est appelé pour configurer le provider (pas nécessaire ici)
func (p *simpleProvider) Configure(context.Context, provider.ConfigureRequest, *provider.ConfigureResponse) {}

// Resources liste les ressources offertes
func (p *simpleProvider) Resources(_ context.Context) []func() provider.Resource {
	return []func() provider.Resource{
		NewItemResource,
	}
}

// DataSources liste les datasources offertes (aucune ici)
func (p *simpleProvider) DataSources(_ context.Context) []func() provider.DataSource {
	return nil
}
```

---

## 4. internal/provider/resource\_item.go

```go
package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ItemResource gère la ressource "simple_item"
type ItemResource struct{}

// NewItemResource crée une nouvelle instance de la ressource
func NewItemResource() resource.Resource {
	return &ItemResource{}
}

// ItemModel représente l’état Terraform de la ressource
type ItemModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
}

// Metadata définit le nom de la ressource dans Terraform
func (r *ItemResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "simple_item"
}

// Schema définit le schéma de la ressource
func (r *ItemResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifiant unique généré automatiquement",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Nom de l'item",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "Description de l'item",
			},
		},
	}
}

// Create crée la ressource
func (r *ItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ItemModel

	// Récupère les valeurs d’entrée
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Génère un ID simple (dans la vraie vie, appelle une API ici)
	data.ID = types.StringValue(fmt.Sprintf("item-%s", data.Name.ValueString()))

	// Enregistre dans l’état Terraform
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Read lit la ressource (ici on fait rien, tout est en mémoire)
func (r *ItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Pas d’API distante, donc rien à faire
}

// Update modifie la ressource
func (r *ItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ItemModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	// On pourrait appeler une API pour mettre à jour, ici on fait juste set l’état
	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
}

// Delete supprime la ressource (ici on ne fait rien de spécial)
func (r *ItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Suppression locale, Terraform supprimera l’état
}
```

---

## 5. Utilisation

Dans un dossier Terraform, crée `main.tf` :

```hcl
terraform {
  required_providers {
    simple = {
      source  = "example/simple"
      version = "0.1.0"
    }
  }
}

provider "simple" {}

resource "simple_item" "foo" {
  name        = "monitem"
  description = "un exemple d'item"
}
```

---

## 6. Build & tester

Compile le provider :

```bash
go mod tidy
go build -o terraform-provider-simple
```

Puis place le binaire :

```bash
mkdir -p ~/.terraform.d/plugins/example/simple/0.1.0/darwin_amd64
mv terraform-provider-simple ~/.terraform.d/plugins/example/simple/0.1.0/darwin_amd64/
```

(adapte `darwin_amd64` selon ton OS, ex: `linux_amd64`)

---

## 7. Lance Terraform

```bash
terraform init
terraform apply
```

Tu verras que Terraform crée ta ressource `simple_item` et gère son état via ton provider.

---

---
