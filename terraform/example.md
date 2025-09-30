Un **Terraform provider custom** te permet de créer tes propres ressources et data sources pour gérer des services ou infrastructures qui ne sont pas encore couverts par les providers officiels.

Voici une explication complète, étape par étape, incluant la mise en place des **tests unitaires** et **tests d’acceptation**.

---

## **1. Concept de Terraform Provider**

Un provider Terraform est un plugin écrit en Go qui :

* **Expose des ressources** (`resource`) → création, mise à jour, suppression d'objets gérés par Terraform.
* **Expose des data sources** (`data source`) → lecture seule de données externes.
* **Interagit avec une API** ou système externe (HTTP, gRPC, SDK, etc.).

Terraform utilise le SDK officiel **terraform-plugin-framework** (nouveau) ou **terraform-plugin-sdk** (ancien).
Le flux typique est :

```
Terraform CLI <-> Provider <-> API externe
```

---

## **2. Structure d’un projet provider**

Voici un squelette classique pour un provider en Go :

```
terraform-provider-example/
├── go.mod
├── main.go                 # Point d'entrée
├── internal/
│   ├── provider.go         # Initialisation du provider
│   ├── resource_foo.go     # Exemple de ressource
│   └── data_source_bar.go  # Exemple de data source
└── tests/
    ├── provider_test.go    # Tests unitaires
    └── acceptance_test.go  # Tests d'acceptation
```

---

## **3. Initialisation avec le SDK Terraform**

Ajout du SDK dans `go.mod` :

```bash
go mod init github.com/username/terraform-provider-example
go get github.com/hashicorp/terraform-plugin-framework
go get github.com/hashicorp/terraform-plugin-testing
```

---

## **4. Exemple de provider minimal**

**`main.go`**

```go
package main

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-example/internal"
)

func main() {
	providerserver.Serve(context.Background(), internal.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/username/example",
	})
}
```

**`internal/provider.go`**

```go
package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

// Provider struct principal
type ExampleProvider struct{}

func New() provider.Provider {
	return &ExampleProvider{}
}

// Définition du schéma du provider
func (p *ExampleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (p *ExampleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "example"
}
```

---

## **5. Exemple de ressource simple**

Créons une ressource `example_resource` qui stocke une valeur dans un service externe.

**`internal/resource_foo.go`**

```go
package internal

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type FooResource struct{}

func NewFooResource() resource.Resource {
	return &FooResource{}
}

func (r *FooResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "example_resource"
}

func (r *FooResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"value": schema.StringAttribute{
				Required: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *FooResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		Value string `tfsdk:"value"`
	}
	req.Plan.Get(ctx, &plan)

	// Ici on simule un appel API
	resp.State.Set(ctx, map[string]string{
		"id":    "resource-123",
		"value": plan.Value,
	})
}

func (r *FooResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Lecture de l'état actuel (normalement on appelle l'API)
	resp.State.Set(ctx, map[string]string{
		"id":    "resource-123",
		"value": "current-value",
	})
}

func (r *FooResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan struct {
		Value string `tfsdk:"value"`
	}
	req.Plan.Get(ctx, &plan)

	resp.State.Set(ctx, map[string]string{
		"id":    "resource-123",
		"value": plan.Value,
	})
}

func (r *FooResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Suppression simulée
	resp.State.RemoveResource(ctx)
}
```

---

## **6. Tests unitaires**

Les tests unitaires valident la logique interne du provider **sans Terraform CLI**.
On utilise `terraform-plugin-testing/helper/resource` pour les tests.

**`tests/provider_test.go`**

```go
package tests

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/stretchr/testify/assert"
	"terraform-provider-example/internal"
)

func TestProviderSchema(t *testing.T) {
	p := internal.New()
	resp := provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	assert.Contains(t, resp.Schema.Attributes, "api_token", "Le provider doit contenir un champ api_token")
}
```

---

## **7. Tests d'acceptation**

Les **tests d'acceptation** exécutent réellement Terraform avec ton provider et valident le comportement global.

**`tests/acceptance_test.go`**

```go
package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleResource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (interface{}, error){
			"example": func() (interface{}, error) {
				return nil, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "example" {
					  api_token = "fake-token"
					}

					resource "example_resource" "test" {
					  value = "Hello World"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("example_resource.test", "value", "Hello World"),
				),
			},
		},
	})
}
```

---

## **8. Lancer les tests**

* **Tests unitaires** :

```bash
go test ./tests -v -run TestProviderSchema
```

* **Tests d'acceptation** (nécessite Terraform CLI installé) :

```bash
TF_ACC=1 go test ./tests -v -run TestAccExampleResource_basic
```

---

## **9. Bonnes pratiques**

* **Séparer la logique API** dans un client interne (ex : `internal/client/`).
* Écrire **des mocks** pour les appels API (utiliser `mockgen` ou `testify/mock`).
* Suivre la convention `TF_ACC=1` pour activer les tests d'acceptation.
* Valider que chaque ressource gère bien les étapes `Create`, `Read`, `Update`, `Delete`.

---

## **Résumé**

| Étape                  | Outil / Fichier               |
| ---------------------- | ----------------------------- |
| Init projet            | `go mod init` + SDK Terraform |
| Définir provider       | `internal/provider.go`        |
| Créer une ressource    | `internal/resource_foo.go`    |
| Tests unitaires        | `tests/provider_test.go`      |
| Tests d'acceptation    | `tests/acceptance_test.go`    |
| Lancer tests unitaires | `go test ./tests`             |
| Lancer tests ACC       | `TF_ACC=1 go test ./tests`    |

Souhaites-tu que je te prépare un **exemple complet packagé** que tu puisses lancer directement avec `terraform apply` ?

Je vais te préparer un **exemple complet de Terraform provider custom**, que tu pourras lancer **localement**, sans passer par le registre Terraform.

Ce guide couvre :

1. La structure du projet complète
2. Le code du provider et d'une ressource simple
3. Le test unitaire et le test d'acceptation
4. Le lancement local avec `terraform apply` via un binaire compilé

---

# **1. Structure du projet**

Voici la structure finale du projet :

```
terraform-provider-example/
├── go.mod
├── go.sum
├── main.go
├── internal/
│   ├── provider.go
│   ├── resource_foo.go
│   └── client/
│       └── fake_client.go
└── tests/
    ├── provider_test.go
    └── acceptance_test.go
```

---

# **2. Initialisation du projet**

Dans un dossier vide :

```bash
mkdir terraform-provider-example
cd terraform-provider-example
go mod init github.com/username/terraform-provider-example
go get github.com/hashicorp/terraform-plugin-framework
go get github.com/hashicorp/terraform-plugin-testing
go get github.com/stretchr/testify
```

---

# **3. Implémentation du provider**

## `main.go`

Point d'entrée du provider.

```go
package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-example/internal"
)

func main() {
	err := providerserver.Serve(context.Background(), internal.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/username/example",
	})

	if err != nil {
		log.Fatalf("Error starting provider: %s", err)
	}
}
```

---

## `internal/provider.go`

Définition du provider et enregistrement de la ressource.

```go
package internal

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Provider principal
type ExampleProvider struct{}

func New() provider.Provider {
	return &ExampleProvider{}
}

func (p *ExampleProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "example"
}

func (p *ExampleProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_token": schema.StringAttribute{
				Required:    true,
				Description: "Token pour l'API fake",
			},
		},
	}
}

func (p *ExampleProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFooResource,
	}
}

func (p *ExampleProvider) DataSources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{}
}
```

---

## `internal/client/fake_client.go`

Un petit client fake qui simule une API externe.

```go
package client

import "fmt"

type FakeClient struct {
	APIToken string
	Store    map[string]string
}

func NewFakeClient(apiToken string) *FakeClient {
	return &FakeClient{
		APIToken: apiToken,
		Store:    make(map[string]string),
	}
}

func (c *FakeClient) Create(value string) string {
	id := fmt.Sprintf("fake-%d", len(c.Store)+1)
	c.Store[id] = value
	return id
}

func (c *FakeClient) Read(id string) (string, bool) {
	v, ok := c.Store[id]
	return v, ok
}

func (c *FakeClient) Update(id, value string) {
	c.Store[id] = value
}

func (c *FakeClient) Delete(id string) {
	delete(c.Store, id)
}
```

---

## `internal/resource_foo.go`

Notre ressource gérée par Terraform : `example_foo`.

```go
package internal

import (
	"context"
	"terraform-provider-example/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type FooResource struct {
	client *client.FakeClient
}

func NewFooResource() resource.Resource {
	return &FooResource{}
}

func (r *FooResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "example_foo"
}

func (r *FooResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "ID généré par la fake API",
			},
			"value": schema.StringAttribute{
				Required:    true,
				Description: "Valeur à stocker",
			},
		},
	}
}

func (r *FooResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	apiToken := req.Config.GetAttribute(ctx, "api_token")
	if apiToken.IsNull() {
		return
	}
	r.client = client.NewFakeClient(apiToken.String())
}

func (r *FooResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		Value string `tfsdk:"value"`
	}
	req.Plan.Get(ctx, &plan)

	id := r.client.Create(plan.Value)
	resp.State.Set(ctx, map[string]string{
		"id":    id,
		"value": plan.Value,
	})
}

func (r *FooResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state struct {
		ID string `tfsdk:"id"`
	}
	req.State.Get(ctx, &state)

	value, exists := r.client.Read(state.ID)
	if !exists {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.State.Set(ctx, map[string]string{
		"id":    state.ID,
		"value": value,
	})
}

func (r *FooResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan struct {
		ID    string `tfsdk:"id"`
		Value string `tfsdk:"value"`
	}
	req.Plan.Get(ctx, &plan)

	r.client.Update(plan.ID, plan.Value)
	resp.State.Set(ctx, map[string]string{
		"id":    plan.ID,
		"value": plan.Value,
	})
}

func (r *FooResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state struct {
		ID string `tfsdk:"id"`
	}
	req.State.Get(ctx, &state)

	r.client.Delete(state.ID)
	resp.State.RemoveResource(ctx)
}
```

---

# **4. Tests**

## `tests/provider_test.go` (Unit Test)

Teste la configuration du provider uniquement.

```go
package tests

import (
	"context"
	"testing"
	"terraform-provider-example/internal"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/stretchr/testify/assert"
)

func TestProviderSchema(t *testing.T) {
	p := internal.New()
	resp := provider.SchemaResponse{}
	p.Schema(context.Background(), provider.SchemaRequest{}, &resp)

	assert.Contains(t, resp.Schema.Attributes, "api_token")
}
```

---

## `tests/acceptance_test.go` (Acceptance Test)

Teste un `terraform apply` réel.

```go
package tests

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccExampleFooResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: map[string]func() (interface{}, error){
			"example": func() (interface{}, error) {
				return nil, nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
					provider "example" {
					  api_token = "my-token"
					}

					resource "example_foo" "test" {
					  value = "Hello World"
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("example_foo.test", "value", "Hello World"),
				),
			},
		},
	})
}
```

Pour exécuter :

```bash
TF_ACC=1 go test ./tests -v -run TestAccExampleFooResource
```

---

# **5. Lancer en local avec Terraform**

## **5.1. Compiler le provider**

```bash
go build -o terraform-provider-example
```

Le binaire `terraform-provider-example` est maintenant prêt.

---

## **5.2. Configurer Terraform pour un provider local**

Crée un dossier pour ton test local :

```
terraform-test/
└── main.tf
```

**main.tf** :

```hcl
terraform {
  required_providers {
    example = {
      source = "registry.terraform.io/username/example"
      version = "0.0.1"
    }
  }
}

provider "example" {
  api_token = "my-local-token"
}

resource "example_foo" "local_test" {
  value = "Bonjour Terraform"
}
```

---

## **5.3. Utiliser le provider local**

Terraform cherche le provider dans `~/.terraform.d/plugins/`.

Crée le chemin attendu :

```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/darwin_amd64
```

Puis déplace le binaire :

```bash
mv terraform-provider-example ~/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/darwin_amd64/
```

*(Remplace `darwin_amd64` par `linux_amd64` si tu es sur Linux)*

---

## **5.4. Lancer Terraform**

Dans `terraform-test/` :

```bash
terraform init
terraform apply
```

Exemple de sortie attendue :

```
Terraform will perform the following actions:

  # example_foo.local_test will be created
  + resource "example_foo" "local_test" {
      + id    = (known after apply)
      + value = "Bonjour Terraform"
    }

Plan: 1 to add, 0 to change, 0 to destroy.
```

---

# **6. Résumé**

| Étape              | Commande / Action                                                                          |
| ------------------ | ------------------------------------------------------------------------------------------ |
| Compiler           | `go build -o terraform-provider-example`                                                   |
| Préparer plugins   | `mkdir -p ~/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/linux_amd64` |
| Déplacer plugin    | `mv terraform-provider-example ~/.terraform.d/plugins/...`                                 |
| Init Terraform     | `terraform init`                                                                           |
| Apply Terraform    | `terraform apply`                                                                          |
| Test unitaire      | `go test ./tests -run TestProviderSchema`                                                  |
| Test d'acceptation | `TF_ACC=1 go test ./tests -v -run TestAccExampleFooResource`                               |

---
