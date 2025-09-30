Créer un **provider Terraform custom** en Go implique de comprendre la structure interne d’un provider et le flux d’exécution entre **Terraform CLI**, ton **plugin**, et éventuellement une **API externe**.

Je vais t’expliquer tout le processus étape par étape :

1. **Architecture et concepts**
2. **Initialisation d’un projet Go**
3. **Structure d’un provider Terraform**
4. **Implémentation d’un provider minimal**
5. **Création d’une resource et d’un data source**
6. **Test en local et avec Terraform**
7. **Tests unitaires et d’intégration**

---

## **1. Architecture d’un provider Terraform**

Terraform suit un modèle **plugin**.
Quand tu exécutes `terraform apply` :

```
Terraform CLI <----> Provider <----> API externe / Services
```

* **Terraform CLI** lit tes fichiers `.tf` et construit un plan.
* Il démarre ton **provider** comme un processus externe (un binaire Go).
* La communication se fait via RPC (plugin gRPC).
* Ton provider exécute des opérations sur une API ou un système (ex: Kubernetes, AWS, base de données).

Un provider est donc essentiellement :

* Une **déclaration du schéma** des ressources et data sources.
* Du **code Go** qui sait gérer le CRUD : `Create`, `Read`, `Update`, `Delete`.

---

## **2. Initialiser le projet Go**

Créons un nouveau projet :

```bash
mkdir terraform-provider-demo
cd terraform-provider-demo
go mod init github.com/ton-org/terraform-provider-demo
go mod tidy
```

Ajoutons la dépendance officielle Terraform Plugin Framework :

```bash
go get github.com/hashicorp/terraform-plugin-framework
```

Le framework `terraform-plugin-framework` est la version moderne recommandée par HashiCorp (à utiliser plutôt que `terraform-plugin-sdk/v2`).

---

## **3. Structure d’un provider Terraform**

Un provider minimal comporte :

```
terraform-provider-demo/
├── go.mod
├── main.go             # Point d'entrée
├── provider.go         # Définition du provider
├── resource_item.go    # Une ressource exemple
└── data_source_item.go # Un data source exemple
```

* **`main.go`** : démarre le plugin et l’enregistre auprès de Terraform.
* **`provider.go`** : décrit le provider et ses configurations globales (`endpoint`, `token`, etc.).
* **`resource_item.go`** : implémente une ressource gérée par Terraform (`Create`, `Read`, `Update`, `Delete`).
* **`data_source_item.go`** : implémente un data source en lecture seule.

---

## **4. Implémentation du provider minimal**

### **main.go**

```go
package main

import (
    "context"
    "github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/ton-org/terraform-provider-demo/demo"
)

func main() {
    providerserver.Serve(context.Background(), demo.New, providerserver.ServeOpts{
        Address: "local/demo/demo",
    })
}
```

* `providerserver.Serve` démarre le provider et écoute les appels RPC de Terraform.
* `Address` définit l’identifiant du provider (`local/demo/demo` pour un provider local).

---

### **provider.go**

```go
package demo

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
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
                Description: "URL de l'API externe",
            },
            "token": schema.StringAttribute{
                Optional:    true,
                Description: "Clé API pour authentification",
                Sensitive:   true,
            },
        },
    }
}

func (p *DemoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    // Ici, tu pourrais créer un client HTTP ou gRPC vers ton service
}

func (p *DemoProvider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewItemResource, // Ajout de la ressource
    }
}

func (p *DemoProvider) DataSources(_ context.Context) []func() datasource.DataSource {
    return []func() datasource.DataSource{
        NewItemDataSource, // Ajout du data source
    }
}
```

---

## **5. Créer une Resource Terraform**

Une ressource doit implémenter les 4 opérations CRUD :

* **Create** → Terraform crée l'objet dans ton système externe.
* **Read** → Terraform vérifie l'état réel.
* **Update** → Terraform modifie l'objet.
* **Delete** → Terraform supprime l'objet.

### **resource\_item.go**

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
                Computed: true,
            },
            "name": schema.StringAttribute{
                Required: true,
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

    // Génération d'un ID fictif
    plan.ID = types.StringValue(fmt.Sprintf("item-%s", plan.Name.ValueString()))

    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *ItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state ItemModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

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
    // Terraform retirera simplement l'objet du state
}
```

---

## **6. Créer un Data Source**

Un data source est une source **en lecture seule**, comme une recherche ou une info provenant d'une API.

### **data\_source\_item.go**

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
                Required: true,
            },
            "name": schema.StringAttribute{
                Computed: true,
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

    // Simule un retour API
    data.Name = types.StringValue(fmt.Sprintf("Nom de %s", data.ID.ValueString()))
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

---

## **7. Test du provider en local**

### **Compilation**

```bash
go build -o terraform-provider-demo
```

### **Installation**

Terraform cherche les providers dans `~/.terraform.d/plugins` :

```bash
mkdir -p ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64
mv terraform-provider-demo ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64/
```

> Remplace `linux_amd64` par `darwin_amd64` (Mac) ou `windows_amd64` (Windows).

---

## **8. Exemple Terraform (`main.tf`)**

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
  endpoint = "http://localhost:8080"
}

resource "demo_item" "example" {
  name = "my-item"
}

data "demo_item_info" "info" {
  id = demo_item.example.id
}

output "item_id" {
  value = demo_item.example.id
}

output "item_name" {
  value = data.demo_item_info.info.name
}
```

---

### **Initialiser et appliquer**

```bash
terraform init
terraform apply -auto-approve
```

Résultat attendu :

```
Outputs:

item_id = "item-my-item"
item_name = "Nom de item-my-item"
```

---

## **9. Tests**

### **Tests unitaires Go**

```bash
go test ./... -cover
```

### **Tests d’intégration Terraform**

En utilisant `terraform-plugin-testing` :

```bash
TF_ACC=1 go test ./internal/tests -v
```

---

## **Résumé des étapes**

| Étape | Action                                            |
| ----- | ------------------------------------------------- |
| 1     | Créer un projet Go avec `go mod init`             |
| 2     | Ajouter `terraform-plugin-framework`              |
| 3     | Écrire `main.go` pour lancer le provider          |
| 4     | Définir le provider et son schéma (`provider.go`) |
| 5     | Implémenter des ressources CRUD                   |
| 6     | Implémenter des data sources                      |
| 7     | Compiler et installer le provider                 |
| 8     | Tester avec `terraform init` et `terraform apply` |
| 9     | Ajouter des tests unitaires et d’intégration      |

---

Souhaites-tu que je t’ajoute une **explication détaillée de la gestion de l’état Terraform** côté provider ? Cela te permettrait de bien comprendre comment Terraform suit les ressources.

La **gestion de l’état (state) dans un provider Terraform** est essentielle pour que Terraform puisse suivre les ressources qu’il crée et éviter des incohérences entre ce qui est **déclaré dans le code** et ce qui **existe réellement dans le système externe**.

Voyons en détail :

---

## **1. Qu'est-ce que le state dans Terraform**

Terraform stocke l'état de l'infrastructure dans un fichier appelé `terraform.tfstate`.

Ce fichier contient :

* L'identifiant des ressources réelles dans le système externe.
* Les valeurs actuelles des attributs connus de Terraform.
* Les métadonnées nécessaires pour faire la comparaison lors d'un `terraform plan`.

### Exemple de `terraform.tfstate` pour notre `demo_item` :

```json
{
  "resources": [
    {
      "type": "demo_item",
      "name": "example",
      "instances": [
        {
          "attributes": {
            "id": "item-my-item",
            "name": "my-item"
          }
        }
      ]
    }
  ]
}
```

> 💡 **But principal** :
> L'état sert de **source de vérité locale** pour que Terraform sache si une ressource doit être créée, mise à jour ou supprimée.

---

## **2. Cycle de vie d'une ressource Terraform**

Pour une ressource, Terraform suit un flux CRUD classique :

| Étape      | Méthode Go dans le provider | Rôle                                                                          |
| ---------- | --------------------------- | ----------------------------------------------------------------------------- |
| **Create** | `Create`                    | Crée la ressource dans le système externe et enregistre l'état dans Terraform |
| **Read**   | `Read`                      | Vérifie que la ressource existe et met à jour l'état                          |
| **Update** | `Update`                    | Met à jour la ressource et l'état                                             |
| **Delete** | `Delete`                    | Supprime la ressource du système externe et du state                          |

> **Note :**
> Après chaque action, Terraform **synchronise** ce que tu renvoies dans le `State` avec le fichier `.tfstate`.

---

## **3. Comment l'état est géré côté provider**

Le provider interagit avec Terraform via le `State`, qui est une **représentation en mémoire** de la ressource.

Prenons notre ressource `demo_item`.
Voici les points clés :

---

### **3.1. Schéma et modèle Go**

Chaque attribut de la ressource doit être défini dans deux endroits :

* **Dans le schéma Terraform** → ce que l'utilisateur peut définir dans son `.tf`.
* **Dans une struct Go** → qui sert à mapper les valeurs entre Terraform et ton code.

Exemple dans `resource_item.go` :

```go
type ItemModel struct {
    ID   types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
}
```

> Le tag `tfsdk:"name"` relie l'attribut Terraform `name` à la variable Go `Name`.

---

### **3.2. Create → Insérer dans le state**

Lors de la création, Terraform appelle `Create`.
Le provider doit :

1. Lire le plan proposé (`req.Plan.Get`).
2. Créer la ressource dans le système externe (API, DB, etc.).
3. Écrire la ressource dans le state (`resp.State.Set`).

```go
func (r *ItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan ItemModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Simuler un appel API pour créer la ressource
    plan.ID = types.StringValue(fmt.Sprintf("item-%s", plan.Name.ValueString()))

    // Enregistrer la ressource dans le state
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
```

---

### **3.3. Read → Synchroniser l'état réel**

La méthode `Read` est essentielle pour **détecter les dérives** entre Terraform et la réalité.

Objectifs :

* Vérifier que la ressource existe encore dans le système externe.
* Mettre à jour les attributs dans Terraform.
* Supprimer la ressource du state si elle n'existe plus.

```go
func (r *ItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state ItemModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Simuler un appel API pour vérifier que l'item existe
    exists := true // remplacer par un appel réel
    if !exists {
        // La ressource a été supprimée en dehors de Terraform
        resp.State.RemoveResource(ctx)
        return
    }

    // Mettre à jour le state avec les valeurs actuelles
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
```

> 💡 **Bonnes pratiques :**
>
> * Si l'API retourne une erreur "Not Found", **supprimer la ressource du state**.
> * Toujours remettre à jour les champs calculés (`Computed`).

---

### **3.4. Update → Modifier et synchroniser**

Lorsque l'utilisateur change un attribut dans son fichier `.tf`, Terraform appelle `Update`.

```go
func (r *ItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan ItemModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Simuler une mise à jour via API
    plan.ID = types.StringValue(fmt.Sprintf("item-%s", plan.Name.ValueString()))

    // Mise à jour du state
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}
```

---

### **3.5. Delete → Nettoyer le state**

Lors de la suppression :

* Supprimer la ressource dans le système externe.
* Terraform supprime automatiquement l'entrée du `.tfstate`.

```go
func (r *ItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state ItemModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Suppression côté API
    // Exemple fictif : deleteAPI(state.ID.ValueString())

    // Rien à faire côté state, Terraform s'en occupe
}
```

---

## **4. Synchronisation du state avec Terraform**

| Action dans le `.tf`        | Méthode Go appelée | Ce qui se passe                                      |
| --------------------------- | ------------------ | ---------------------------------------------------- |
| Ajout d'une ressource       | `Create` + `Read`  | Création dans le système et ajout dans le `.tfstate` |
| Modification d'un champ     | `Update` + `Read`  | Mise à jour dans le système et `.tfstate`            |
| Suppression d'une ressource | `Delete`           | Suppression dans le système et `.tfstate`            |
| `terraform refresh`         | `Read`             | Synchronisation entre le `.tfstate` et l'état réel   |

---

## **5. Exemple complet du cycle**

### **main.tf**

```hcl
resource "demo_item" "test" {
  name = "my-item"
}
```

### **Workflow Terraform**

```bash
terraform apply
```

1. Terraform détecte qu'il faut **créer** la ressource.
2. Appelle `Create`.
3. Le provider retourne `id = "item-my-item"`.
4. Terraform ajoute ceci au `.tfstate`.

```bash
terraform state show demo_item.test
```

Résultat :

```
# demo_item.test:
resource "demo_item" "test" {
    id   = "item-my-item"
    name = "my-item"
}
```

---

## **6. Détection des dérives (Drift Detection)**

Supposons que quelqu'un supprime l'objet **hors Terraform** :

* Terraform lance `Read`.
* Le provider voit que l'objet n'existe plus.
* Il appelle `resp.State.RemoveResource(ctx)`.
* Terraform propose de le recréer lors du prochain `terraform plan`.

```go
if !exists {
    resp.State.RemoveResource(ctx)
    return
}
```

---

## **7. Résumé visuel**

```
+-----------------------+
| .tf files (config)    |
+-----------------------+
          |
          v
+-----------------------+
| Terraform CLI         |
+-----------------------+
          |
          v
+-----------------------+
| Provider (Go plugin)  |
| - Create/Read/Update  |
| - CRUD Operations     |
+-----------------------+
          |
          v
+-----------------------+
| API / System externe  |
+-----------------------+

State synchronisé entre
Terraform et le Provider
```

---

## **8. Bonnes pratiques pour le state**

1. **Toujours gérer l'état réel dans `Read`**

   * Si une ressource a disparu, la retirer du state.
   * Mettre à jour les champs calculés.

2. **Ne jamais stocker d'informations sensibles en clair**

   * Marquer les attributs comme `Sensitive: true` dans le schéma.

3. **Garder le state cohérent**

   * Après chaque modification (`Create`/`Update`), toujours écrire dans `resp.State.Set`.

4. **Valider les entrées utilisateur**

   * Utiliser des `Validators` dans le schéma pour éviter des incohérences.

---

## **Récapitulatif**

| Étape                       | But                                        | Méthode  |
| --------------------------- | ------------------------------------------ | -------- |
| Définir le schéma           | Décrire les attributs                      | `Schema` |
| Créer une ressource         | Générer un ID et enregistrer dans le state | `Create` |
| Lire l'état réel            | Synchroniser Terraform et l'API            | `Read`   |
| Mettre à jour une ressource | Appliquer des changements et resync        | `Update` |
| Supprimer une ressource     | Nettoyer l'API et le state                 | `Delete` |

---

Souhaites-tu que je t'ajoute un **exemple de gestion d'état avancée** avec une vraie API, par exemple en appelant une API REST fictive ?

Parfait ! Je vais te montrer un **exemple avancé de provider Terraform** en Go qui interagit avec une **API REST externe**.
Cet exemple met en avant une gestion d’état réaliste avec un vrai flux CRUD (Create, Read, Update, Delete).

On va créer un provider `demo` qui gère des **items** via une API REST fictive exposant des endpoints comme :

| Action         | Méthode | Endpoint      | Corps requis             |
| -------------- | ------- | ------------- | ------------------------ |
| Créer item     | POST    | `/items`      | `{ "name": "my-item" }`  |
| Lire item      | GET     | `/items/{id}` | -                        |
| Mettre à jour  | PUT     | `/items/{id}` | `{ "name": "new-name" }` |
| Supprimer item | DELETE  | `/items/{id}` | -                        |

---

# **1. Structure du projet**

```
terraform-provider-demo/
├── go.mod
├── main.go
├── provider.go
├── client.go            # Client HTTP pour l'API REST
├── resource_item.go     # Resource Terraform "demo_item"
├── data_source_item.go  # Data source pour récupérer un item
└── examples/
    └── main.tf
```

---

# **2. Dépendances et initialisation**

### `go.mod`

```go
module github.com/ton-org/terraform-provider-demo

go 1.22

require (
    github.com/hashicorp/terraform-plugin-framework v1.10.0
    github.com/hashicorp/terraform-plugin-log v0.19.0
)
```

Initialiser le projet :

```bash
go mod tidy
```

---

# **3. Création du client HTTP (`client.go`)**

Ce client encapsule les appels à l'API externe et sera partagé dans le provider.

```go
package demo

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type APIClient struct {
    endpoint string
    token    string
    client   *http.Client
}

func NewAPIClient(endpoint, token string) *APIClient {
    return &APIClient{
        endpoint: endpoint,
        token:    token,
        client:   &http.Client{},
    }
}

type Item struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

// ===== CREATE =====
func (c *APIClient) CreateItem(name string) (*Item, error) {
    body, _ := json.Marshal(map[string]string{"name": name})
    req, err := http.NewRequest("POST", fmt.Sprintf("%s/items", c.endpoint), bytes.NewBuffer(body))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusCreated {
        data, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error: %s", string(data))
    }

    var item Item
    json.NewDecoder(resp.Body).Decode(&item)
    return &item, nil
}

// ===== READ =====
func (c *APIClient) GetItem(id string) (*Item, error) {
    req, _ := http.NewRequest("GET", fmt.Sprintf("%s/items/%s", c.endpoint, id), nil)
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode == http.StatusNotFound {
        return nil, nil // Not found
    }

    if resp.StatusCode != http.StatusOK {
        data, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error: %s", string(data))
    }

    var item Item
    json.NewDecoder(resp.Body).Decode(&item)
    return &item, nil
}

// ===== UPDATE =====
func (c *APIClient) UpdateItem(id, name string) (*Item, error) {
    body, _ := json.Marshal(map[string]string{"name": name})
    req, _ := http.NewRequest("PUT", fmt.Sprintf("%s/items/%s", c.endpoint, id), bytes.NewBuffer(body))
    req.Header.Set("Content-Type", "application/json")
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        data, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API error: %s", string(data))
    }

    var item Item
    json.NewDecoder(resp.Body).Decode(&item)
    return &item, nil
}

// ===== DELETE =====
func (c *APIClient) DeleteItem(id string) error {
    req, _ := http.NewRequest("DELETE", fmt.Sprintf("%s/items/%s", c.endpoint, id), nil)
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }

    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
        data, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("API error: %s", string(data))
    }
    return nil
}
```

---

# **4. Provider (`provider.go`)**

Ici, on configure notre provider pour initialiser le `APIClient`.

```go
package demo

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
)

type DemoProvider struct {
    client *APIClient
}

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
                Required:    true,
                Description: "Endpoint de l'API REST",
            },
            "token": schema.StringAttribute{
                Optional:    true,
                Sensitive:   true,
                Description: "Token d'authentification",
            },
        },
    }
}

func (p *DemoProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    var config struct {
        Endpoint string `tfsdk:"endpoint"`
        Token    string `tfsdk:"token"`
    }

    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }

    if config.Endpoint == "" {
        resp.Diagnostics.AddError(
            "Configuration manquante",
            "Le champ 'endpoint' est requis pour se connecter à l'API.",
        )
        return
    }

    p.client = NewAPIClient(config.Endpoint, config.Token)
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

# **5. Ressource avec gestion d'état (`resource_item.go`)**

Cette ressource gère un objet dans l'API, en synchronisant le state Terraform avec l'état réel côté API.

```go
package demo

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type ItemResource struct {
    client *APIClient
}

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
                Description: "Identifiant unique de l'item",
            },
            "name": schema.StringAttribute{
                Required:    true,
                Description: "Nom de l'item",
            },
        },
    }
}

func (r *ItemResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }
    r.client = req.ProviderData.(*APIClient)
}

// ===== CREATE =====
func (r *ItemResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan ItemModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    item, err := r.client.CreateItem(plan.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError("Erreur création", fmt.Sprintf("Impossible de créer l'item: %s", err))
        return
    }

    plan.ID = types.StringValue(item.ID)
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// ===== READ =====
func (r *ItemResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state ItemModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    item, err := r.client.GetItem(state.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError("Erreur lecture", fmt.Sprintf("Impossible de lire l'item: %s", err))
        return
    }

    if item == nil {
        // L'item n'existe plus côté API, on le supprime du state
        resp.State.RemoveResource(ctx)
        return
    }

    state.Name = types.StringValue(item.Name)
    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// ===== UPDATE =====
func (r *ItemResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan ItemModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    item, err := r.client.UpdateItem(plan.ID.ValueString(), plan.Name.ValueString())
    if err != nil {
        resp.Diagnostics.AddError("Erreur mise à jour", fmt.Sprintf("Impossible de mettre à jour l'item: %s", err))
        return
    }

    plan.Name = types.StringValue(item.Name)
    resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

// ===== DELETE =====
func (r *ItemResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state ItemModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    if err := r.client.DeleteItem(state.ID.ValueString()); err != nil {
        resp.Diagnostics.AddError("Erreur suppression", fmt.Sprintf("Impossible de supprimer l'item: %s", err))
        return
    }
}
```

---

# **6. Data Source pour lecture (`data_source_item.go`)**

Ce data source permet de récupérer un item par son `id`.

```go
package demo

import (
    "context"
    "fmt"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type ItemDataSource struct {
    client *APIClient
}

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
                Description: "ID de l'item à récupérer",
            },
            "name": schema.StringAttribute{
                Computed:    true,
                Description: "Nom de l'item",
            },
        },
    }
}

func (d *ItemDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }
    d.client = req.ProviderData.(*APIClient)
}

func (d *ItemDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    var data ItemDataSourceModel
    resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
    if resp.Diagnostics.HasError() {
        return
    }

    item, err := d.client.GetItem(data.ID.ValueString())
    if err != nil {
        resp.Diagnostics.AddError("Erreur data source", fmt.Sprintf("Impossible de récupérer l'item: %s", err))
        return
    }

    if item == nil {
        resp.Diagnostics.AddError("Item introuvable", "L'ID fourni ne correspond à aucun item existant")
        return
    }

    data.Name = types.StringValue(item.Name)
    resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
```

---

# **7. Exemple Terraform (`examples/main.tf`)**

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
  endpoint = "http://localhost:8080"
  token    = "my-secret-token"
}

resource "demo_item" "my_item" {
  name = "mon-item"
}

data "demo_item_info" "fetched" {
  id = demo_item.my_item.id
}

output "created_item_id" {
  value = demo_item.my_item.id
}

output "fetched_item_name" {
  value = data.demo_item_info.fetched.name
}
```

---

# **8. Compilation et installation locale**

### Compilation

```bash
go build -o terraform-provider-demo
```

### Installation dans le dossier des plugins

```bash
mkdir -p ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64
mv terraform-provider-demo ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64/
```

---

# **9. Test avec Terraform**

```bash
cd examples
terraform init
terraform apply -auto-approve
```

Sortie attendue :

```
Outputs:

created_item_id = "12345"
fetched_item_name = "mon-item"
```

---

# **10. Points clés sur la gestion d'état**

| Étape  | Gestion dans le provider                                                                     |
| ------ | -------------------------------------------------------------------------------------------- |
| Create | Enregistrer l'ID de l'API dans `resp.State.Set`                                              |
| Read   | Vérifier l'existence de l'objet → supprimer du state si absent (`resp.State.RemoveResource`) |
| Update | Mettre à jour dans l'API puis synchroniser avec `resp.State.Set`                             |
| Delete | Supprimer côté API, Terraform nettoie automatiquement le `.tfstate`                          |

---

# **11. Améliorations possibles**

* Ajouter des **validators** pour vérifier les entrées utilisateur (`Validators`).
* Implémenter la pagination si l’API renvoie plusieurs résultats.
* Gérer les **timeouts** et **retries** pour plus de résilience.
* Ajouter des **tests d’intégration** avec `terraform-plugin-testing`.
* Signer et publier le provider sur le **Terraform Registry**.

---

# **Résumé**

Nous avons maintenant un provider complet qui :

* Configure une connexion à une API REST.
* Implémente une ressource avec CRUD complet.
* Met à jour le state Terraform automatiquement.
* Fournit un data source en lecture seule.
* Peut être testé en local avec `terraform apply`.

Pour **tester ton provider Terraform en local avec `terraform apply`**, il faut suivre une série d'étapes précises afin que Terraform puisse reconnaître ton plugin et l'exécuter correctement.

Voici la procédure complète :

---

## **1. Compiler le provider**

Dans le répertoire racine du projet :

```bash
go build -o terraform-provider-demo
```

Cela produit un binaire nommé `terraform-provider-demo` correspondant à notre provider.

> 💡 **Astuce :** tu peux ajouter `GOOS` et `GOARCH` pour cross-compiler, par exemple :
>
> ```bash
> GOOS=linux GOARCH=amd64 go build -o terraform-provider-demo
> ```

---

## **2. Créer le dossier plugin attendu par Terraform**

Terraform s'attend à trouver le provider dans un chemin très spécifique basé sur :

```
~/.terraform.d/plugins/<namespace>/<provider_name>/<version>/<os>_<arch>/
```

Dans notre exemple :

| Champ         | Valeur                                      |
| ------------- | ------------------------------------------- |
| Namespace     | `local`                                     |
| Provider name | `demo`                                      |
| Version       | `0.1.0`                                     |
| OS/Arch       | `linux_amd64` (ou `darwin_amd64` sur macOS) |

Crée le dossier complet :

```bash
mkdir -p ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64
```

---

## **3. Déplacer le binaire dans ce dossier**

Déplace ton provider compilé dans le répertoire créé :

```bash
mv terraform-provider-demo ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64/
```

Vérifie :

```bash
ls ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64/
```

Tu dois voir :

```
terraform-provider-demo
```

---

## **4. Créer la configuration Terraform d'exemple**

Dans le dossier `examples/`, crée un fichier `main.tf` :

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
  endpoint = "http://localhost:8080"
  token    = "my-secret-token"
}

resource "demo_item" "my_item" {
  name = "mon-item"
}

data "demo_item_info" "fetched" {
  id = demo_item.my_item.id
}

output "created_item_id" {
  value = demo_item.my_item.id
}

output "fetched_item_name" {
  value = data.demo_item_info.fetched.name
}
```

---

## **5. Lancer une API REST factice pour tester**

Si tu n'as pas encore de vraie API REST, tu peux utiliser un **serveur mock** avec Python ou Go.

### Option A : Serveur Python ultra-simple

```bash
python3 -m http.server 8080
```

⚠️ Ça ne gérera que des réponses statiques.
Pour simuler une vraie API CRUD, mieux vaut utiliser [json-server](https://github.com/typicode/json-server) ou un petit serveur Go.

---

### Option B : Serveur Go minimaliste

Crée un fichier `mock_api.go` :

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type Item struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}

var store = make(map[string]Item)

func main() {
    http.HandleFunc("/items", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            var data map[string]string
            json.NewDecoder(r.Body).Decode(&data)
            id := fmt.Sprintf("item-%s", data["name"])
            item := Item{ID: id, Name: data["name"]}
            store[id] = item
            w.WriteHeader(http.StatusCreated)
            json.NewEncoder(w).Encode(item)
            return
        }
        w.WriteHeader(http.StatusMethodNotAllowed)
    })

    http.HandleFunc("/items/", func(w http.ResponseWriter, r *http.Request) {
        id := r.URL.Path[len("/items/"):]
        switch r.Method {
        case http.MethodGet:
            item, ok := store[id]
            if !ok {
                http.NotFound(w, r)
                return
            }
            json.NewEncoder(w).Encode(item)
        case http.MethodPut:
            var data map[string]string
            json.NewDecoder(r.Body).Decode(&data)
            item, ok := store[id]
            if !ok {
                http.NotFound(w, r)
                return
            }
            item.Name = data["name"]
            store[id] = item
            json.NewEncoder(w).Encode(item)
        case http.MethodDelete:
            delete(store, id)
            w.WriteHeader(http.StatusNoContent)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })

    log.Println("API Mock sur http://localhost:8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Lance l'API :

```bash
go run mock_api.go
```

---

## **6. Initialiser Terraform**

Dans le dossier `examples/` :

```bash
terraform init
```

Terraform doit reconnaître ton provider :

```
Initializing provider plugins...
- Finding local/demo/demo versions matching "0.1.0"...
- Installing local/demo/demo v0.1.0...
- Installed local/demo/demo v0.1.0 (unauthenticated)

Terraform has been successfully initialized!
```

---

## **7. Tester avec `terraform apply`**

Toujours dans `examples/` :

```bash
terraform apply -auto-approve
```

Sortie attendue :

```
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

created_item_id = "item-mon-item"
fetched_item_name = "mon-item"
```

---

## **8. Vérifier l'état**

Montre ce que Terraform a enregistré dans le state :

```bash
terraform state show demo_item.my_item
```

Résultat :

```
# demo_item.my_item:
resource "demo_item" "my_item" {
    id   = "item-mon-item"
    name = "mon-item"
}
```

---

## **9. Scénarios de test**

Maintenant que tout est en place, tu peux tester différents cas :

| Action                           | Commande                                                   | Comportement attendu                                 |
| -------------------------------- | ---------------------------------------------------------- | ---------------------------------------------------- |
| Création                         | `terraform apply`                                          | Crée un nouvel item                                  |
| Modification du `name`           | `terraform apply` après avoir changé `name` dans `main.tf` | Appelle `Update` et met à jour côté API              |
| Suppression                      | `terraform destroy`                                        | Appelle `Delete` et supprime l'item                  |
| Ressource supprimée manuellement | Supprime l'item via API puis `terraform plan`              | Terraform propose de le recréer (détection de drift) |

---

## **10. Récapitulatif**

| Étape                 | Commande                                                                               |
| --------------------- | -------------------------------------------------------------------------------------- |
| Compiler le provider  | `go build -o terraform-provider-demo`                                                  |
| Créer dossier plugin  | `mkdir -p ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64`                    |
| Déplacer binaire      | `mv terraform-provider-demo ~/.terraform.d/plugins/local/demo/demo/0.1.0/linux_amd64/` |
| Initialiser Terraform | `terraform init`                                                                       |
| Appliquer la config   | `terraform apply`                                                                      |
| Vérifier l'état       | `terraform state show demo_item.my_item`                                               |

---

Avec ce setup, tu peux tester **tout le cycle Terraform** localement, sans avoir besoin de publier ton provider sur le Terraform Registry.
C'est la méthode standard pour le développement et la validation d'un provider custom.
