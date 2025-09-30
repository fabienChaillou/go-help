Oui ! Comprendre la méthode `Schema` d’une ressource Terraform est **fondamental** pour bien structurer un provider custom.

Dans Terraform, le schéma (`Schema`) définit **l’interface entre l’utilisateur et le provider**, c’est-à-dire :

* Ce que l’utilisateur peut définir dans le fichier `.tf`
* Ce qui est **obligatoire**, **optionnel** ou **calculé**
* Le **type de donnée** attendu
* La façon dont Terraform doit gérer les changements (force recréation, mise à jour, etc.)

On va détailler **chaque aspect avec des exemples concrets**.

---

# **1. Structure générale de `Schema`**

Une ressource Terraform implémente la méthode `Schema` qui retourne une structure `schema.Schema` :

```go
func (r *MyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "Identifiant unique de la ressource",
            },
            "name": schema.StringAttribute{
                Required:    true,
                Description: "Nom de la ressource",
            },
            "tags": schema.MapAttribute{
                Optional:    true,
                ElementType: types.StringType,
                Description: "Liste de tags sous forme clé/valeur",
            },
        },
    }
}
```

Ici :

* `Attributes` est un dictionnaire avec la clé = nom de l'attribut (`id`, `name`, `tags`)
* La valeur est une configuration spécifique à chaque type (`StringAttribute`, `MapAttribute`, etc.)

---

# **2. Types d'attributs**

Terraform supporte plusieurs types de données, qui correspondent aux types `types.*` du **Terraform Plugin Framework**.

| Type Terraform | Go struct utilisé            | Exemple                          |
| -------------- | ---------------------------- | -------------------------------- |
| `string`       | `schema.StringAttribute`     | Nom, ID                          |
| `bool`         | `schema.BoolAttribute`       | Activer ou non le versioning     |
| `number`       | `schema.NumberAttribute`     | Stockage max en Go, int ou float |
| `map(string)`  | `schema.MapAttribute`        | Tags clé/valeur                  |
| `list(object)` | `schema.ListNestedAttribute` | Liste d’objets complexes         |
| `set(object)`  | `schema.SetNestedAttribute`  | Liste unique d’objets complexes  |

---

## **Exemple simple avec différents types**

```go
Attributes: map[string]schema.Attribute{
    "id": schema.StringAttribute{
        Computed:    true,
        Description: "Identifiant unique généré par l'API",
    },

    "enabled": schema.BoolAttribute{
        Optional:    true,
        Description: "Activer ou désactiver la ressource",
    },

    "retries": schema.NumberAttribute{
        Optional:    true,
        Description: "Nombre de tentatives en cas d'erreur",
    },

    "tags": schema.MapAttribute{
        Optional:    true,
        ElementType: types.StringType,
        Description: "Tags associés à la ressource",
    },
}
```

---

# **3. Propriétés des attributs**

Chaque attribut possède plusieurs propriétés importantes qui définissent son **comportement**.

| Propriété            | Description                                                                     |
| -------------------- | ------------------------------------------------------------------------------- |
| `Required`           | L'attribut **doit** être défini dans la configuration `.tf`.                    |
| `Optional`           | L'attribut est facultatif, l'utilisateur peut le définir ou non.                |
| `Computed`           | La valeur est **calculée par le provider** (ex: un ID généré par l'API).        |
| `Description`        | Texte descriptif visible dans `terraform-docs` et `terraform providers schema`. |
| `DeprecationMessage` | Message affiché si l'attribut est obsolète.                                     |
| `Sensitive`          | Masque la valeur dans les logs et l'état Terraform (ex: mots de passe).         |
| `PlanModifiers`      | Contrôle le comportement lors des changements (ex: forcer la recréation).       |

---

## **Exemple : ID généré automatiquement**

L'attribut `id` est **Computed**, car c'est l'API qui le fournit :

```go
"id": schema.StringAttribute{
    Computed:    true,
    Description: "ID généré par l'API lors de la création",
}
```

Dans Terraform :

```hcl
resource "demo_item" "test" {
  name = "bucket1"
}

output "item_id" {
  value = demo_item.test.id
}
```

---

## **Exemple : champ obligatoire**

Un champ `Required` doit être fourni par l'utilisateur :

```go
"name": schema.StringAttribute{
    Required:    true,
    Description: "Nom unique du bucket",
}
```

Dans Terraform, l'utilisateur **doit** l'écrire :

```hcl
resource "demo_item" "test" {
  name = "bucket1" # Obligatoire
}
```

---

## **Exemple : champ optionnel avec valeur par défaut**

Terraform ne gère pas les valeurs par défaut côté provider directement,
mais tu peux les simuler en Go dans la fonction `Create`.

```go
"enabled": schema.BoolAttribute{
    Optional:    true,
    Description: "Activer la ressource (par défaut: true)",
}
```

Dans la fonction `Create` :

```go
if plan.Enabled.IsNull() {
    plan.Enabled = types.BoolValue(true)
}
```

---

## **Exemple : champ sensible (mot de passe)**

```go
"password": schema.StringAttribute{
    Required:    true,
    Sensitive:   true,
    Description: "Mot de passe de l'utilisateur",
}
```

Dans `terraform apply`, la valeur sera masquée :

```
+ password = (sensitive value)
```

---

# **4. Attributs complexes (Map, List, Set, Object)**

## **Map de strings**

```go
"tags": schema.MapAttribute{
    Optional:    true,
    ElementType: types.StringType,
    Description: "Tags du bucket en clé/valeur",
}
```

Terraform :

```hcl
tags = {
  environment = "dev"
  owner       = "team1"
}
```

---

## **Liste d'objets**

Pour des données structurées complexes :

```go
"users": schema.ListNestedAttribute{
    Optional: true,
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{Required: true},
            "email": schema.StringAttribute{Required: true},
        },
    },
},
```

Terraform :

```hcl
users = [
  {
    name  = "Alice"
    email = "alice@example.com"
  },
  {
    name  = "Bob"
    email = "bob@example.com"
  }
]
```

---

## **Set d'objets (liste sans doublons)**

```go
"buckets": schema.SetNestedAttribute{
    Optional: true,
    NestedObject: schema.NestedAttributeObject{
        Attributes: map[string]schema.Attribute{
            "name": schema.StringAttribute{Required: true},
        },
    },
},
```

Terraform :

```hcl
buckets = [
  { name = "bucket1" },
  { name = "bucket2" },
]
```

---

# **5. Exemple complet d'une ressource**

Exemple d'une ressource S3 avec tous les types :

```go
func (r *BucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "ID généré automatiquement",
            },

            "name": schema.StringAttribute{
                Required:    true,
                Description: "Nom du bucket",
            },

            "versioned": schema.BoolAttribute{
                Optional:    true,
                Description: "Activer le versioning",
            },

            "tags": schema.MapAttribute{
                Optional:    true,
                ElementType: types.StringType,
                Description: "Tags du bucket en clé/valeur",
            },

            "users": schema.ListNestedAttribute{
                Optional: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "name": schema.StringAttribute{Required: true},
                        "email": schema.StringAttribute{Required: true},
                    },
                },
                Description: "Liste des utilisateurs autorisés",
            },
        },
    }
}
```

---

# **6. Plan Modifiers (force recreate, etc.)**

Les `PlanModifiers` permettent de contrôler le comportement lors des changements.
Par exemple, forcer la **recréation de la ressource** si le nom change :

```go
import "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
import "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

"name": schema.StringAttribute{
    Required: true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(), // force un delete/create si le nom change
    },
},
```

Terraform affichera alors :

```
-/+ name = "old" -> "new"
```

Ce qui indique une suppression et une création.

---

# **7. Résumé des bonnes pratiques**

| Cas d’usage                                  | Solution                           |
| -------------------------------------------- | ---------------------------------- |
| Valeur définie par l’utilisateur             | `Required: true`                   |
| Valeur facultative définie par l’utilisateur | `Optional: true`                   |
| Valeur générée par le provider (API)         | `Computed: true`                   |
| Masquer une donnée sensible                  | `Sensitive: true`                  |
| Forcer une recréation si changement          | `PlanModifiers: RequiresReplace()` |
| Structure clé/valeur                         | `MapAttribute`                     |
| Liste ordonnée                               | `ListNestedAttribute`              |
| Liste unique sans doublons                   | `SetNestedAttribute`               |

---

# **8. Exemple Terraform**

Avec notre schéma :

```hcl
resource "localstack_s3_bucket" "example" {
  name      = "my-bucket"
  versioned = true

  tags = {
    environment = "dev"
    owner       = "team1"
  }

  users = [
    {
      name  = "Alice"
      email = "alice@example.com"
    },
    {
      name  = "Bob"
      email = "bob@example.com"
    }
  ]
}
```

---

# **9. Résumé visuel**

| Attribut    | Utilisateur fournit ? | Calculé par provider ? | Type de donnée | Forcer recréation       |
| ----------- | --------------------- | ---------------------- | -------------- | ----------------------- |
| `name`      | ✅ Oui                 | ❌ Non                  | `string`       | Oui (`RequiresReplace`) |
| `versioned` | ✅ Oui (optionnel)     | ❌ Non                  | `bool`         | Non                     |
| `id`        | ❌ Non                 | ✅ Oui                  | `string`       | N/A                     |
| `tags`      | ✅ Oui (optionnel)     | ❌ Non                  | `map(string)`  | Non                     |
| `users`     | ✅ Oui (optionnel)     | ❌ Non                  | `list(object)` | Non                     |

---

En résumé, la méthode `Schema` **définit l'API utilisateur** du provider Terraform.
Elle est cruciale pour :

* Structurer les ressources
* Contrôler les validations
* Définir les relations entre Terraform et l'API réelle
* Gérer le cycle de vie complet des ressources (`Create`, `Read`, `Update`, `Delete`).

Oui, on va détailler ça calmement. Dans le **Terraform Plugin Framework**, les **PlanModifiers** et **ApplyModifiers** sont des mécanismes très puissants pour contrôler comment Terraform **planifie** et **applique** les changements d’une ressource.

---

# **1. Contexte**

Quand Terraform exécute un plan (`terraform plan`) :

1. Il lit la configuration `.tf` (desired state).
2. Il lit l’état actuel (`terraform state`).
3. Il calcule un **plan d’actions** (`Create`, `Update`, `Delete`) pour arriver de l’état actuel à l’état souhaité.

Les **modifiers** interviennent dans cette étape pour **affiner ou transformer les valeurs** avant que Terraform ne décide quoi faire.

* **PlanModifier** → agit **au moment du plan**, avant que Terraform calcule le diff.
* **ApplyModifier** → agit **au moment de l’application**, juste avant d’envoyer les données à la fonction `Create` ou `Update`.

---

# **2. PlanModifier**

Le **PlanModifier** change ou influence la valeur **durant le plan**, pour influencer la décision de Terraform.

### **Types principaux**

1. **RequiresReplace** → force Terraform à recréer la ressource si la valeur change.
2. **UseStateForUnknown** → garde la valeur existante si la valeur du plan est inconnue.
3. **Default** → définit une valeur par défaut pour le plan.

---

### **Exemple : forcer la recréation si le nom change**

```go
import "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
import "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"

"name": schema.StringAttribute{
    Required: true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
},
```

**Effet :**

* Si l’utilisateur change `name = "old"` → `name = "new"`, Terraform planifiera un `- / +` (delete + create).
* Terraform affichera :

```
-/+ name = "old" -> "new" (forces replacement)
```

---

### **Exemple : valeur par défaut au plan**

```go
"enabled": schema.BoolAttribute{
    Optional: true,
    PlanModifiers: []planmodifier.Bool{
        planmodifier.UseStateForUnknown(),
    },
},
```

* Si l’utilisateur ne définit pas `enabled`, Terraform prend la valeur actuelle du state ou une valeur par défaut définie dans le provider.

---

# **3. ApplyModifier**

L’**ApplyModifier** agit **juste avant que Terraform n’applique la valeur à la fonction Create ou Update**.

* On peut transformer ou calculer la valeur finale.
* Très utile pour :

  * Hachage d’un mot de passe
  * Conversion de formats (ex: `"yes"` → `true`)
  * Validation supplémentaire

### **Exemple : transformer un mot de passe avant l’application**

```go
"password": schema.StringAttribute{
    Required:  true,
    Sensitive: true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(),
    },
    ApplyModifiers: []planmodifier.String{
        stringplanmodifier.MaskValue(),
    },
},
```

* `MaskValue()` cache la valeur pour qu’elle ne soit pas exposée dans le state ni dans les logs.
* Terraform continue de l’envoyer à la fonction `Create`/`Update`.

---

# **4. Différence clé entre PlanModifier et ApplyModifier**

| Concept       | Quand il agit             | Effet                                                                                                 |
| ------------- | ------------------------- | ----------------------------------------------------------------------------------------------------- |
| PlanModifier  | Pendant `terraform plan`  | Transforme ou influence la valeur dans le plan, peut forcer replace ou remplir des valeurs par défaut |
| ApplyModifier | Pendant `terraform apply` | Transforme ou sécurise la valeur avant de l’envoyer à `Create` ou `Update`                            |

* **PlanModifier** → influence le **plan et le diff**.
* **ApplyModifier** → influence la **valeur envoyée à l’API**.

---

# **5. Exemple complet dans un Bucket S3**

```go
"versioned": schema.BoolAttribute{
    Optional: true,
    PlanModifiers: []planmodifier.Bool{
        // Si versioning change, Terraform ne recrée pas le bucket (il fait juste update)
        planmodifier.UseStateForUnknown(),
    },
},
"name": schema.StringAttribute{
    Required: true,
    PlanModifiers: []planmodifier.String{
        stringplanmodifier.RequiresReplace(), // change → delete/create
    },
},
"password": schema.StringAttribute{
    Optional: true,
    Sensitive: true,
    ApplyModifiers: []planmodifier.String{
        stringplanmodifier.MaskValue(), // cache la valeur lors de apply
    },
},
```

---

# **6. Résumé pratique**

* **PlanModifier**

  * Agit avant `terraform plan` pour influencer la création du plan.
  * Peut **forcer remplacement**, remplir une valeur par défaut ou conserver l’état actuel.

* **ApplyModifier**

  * Agit avant `terraform apply`, juste avant d’envoyer les valeurs à `Create`/`Update`.
  * Sert pour **sécurité, transformation, masquage ou conversion** des valeurs.

---

Si tu veux, je peux te faire un **schéma visuel du flux Terraform avec PlanModifier et ApplyModifier**, ce qui rend tout ça très clair.

Parfait ! Voici un **schéma visuel simplifié du flux Terraform** avec **PlanModifier** et **ApplyModifier**, pour bien comprendre leur rôle dans le cycle d’une ressource :

```
Utilisateur
   |
   v
Configuration .tf
   |
   |  Terraform lit la configuration
   v
+------------------------+
|      Terraform Plan    |
+------------------------+
| - Compare state actuel |
|   avec desired state   |
| - Crée un plan d'action|
|                        |
|  [PlanModifier]        |  <-- Agit ici pour ajuster / transformer / forcer replace
+------------------------+
   |
   v
Diff / Plan final affiché à l'utilisateur
   |
   v
terraform apply
   |
+------------------------+
|    Apply phase         |
|                        |
|  [ApplyModifier]       |  <-- Agit ici avant d'appeler Create/Update
+------------------------+
   |
   v
Provider (Create / Update / Delete)
   |
   v
API distante (ex: S3, LocalStack)
   |
   v
État mis à jour (terraform state)
```

---

### **Exemple concret avec un bucket S3**

| Attribut    | PlanModifier           | ApplyModifier | Effet                                                                     |
| ----------- | ---------------------- | ------------- | ------------------------------------------------------------------------- |
| `name`      | `RequiresReplace()`    | —             | Si le nom change, Terraform supprime et recrée le bucket                  |
| `versioned` | `UseStateForUnknown()` | —             | Si non défini, utilise l’état actuel pour éviter de réinitialiser le plan |
| `password`  | —                      | `MaskValue()` | Masque le mot de passe avant l’envoi à l’API                              |

---

### **Résumé rapide**

1. **PlanModifier** → influence **le plan**, décide si Terraform doit créer, mettre à jour ou supprimer la ressource.
2. **ApplyModifier** → influence **la valeur envoyée à l’API** pendant `apply`, utile pour transformer ou sécuriser les données.

---

Si tu veux, je peux te faire un **mini exemple Go complet** montrant à la fois **PlanModifier** et **ApplyModifier** dans une ressource Terraform S3 custom, pour le tester directement.

# example en go

Parfait ! Voici un **mini exemple Go complet** d’une ressource Terraform custom qui utilise **PlanModifier** et **ApplyModifier** pour un bucket S3 (LocalStack) :

---

# **1. Fichier : `resource_bucket.go`**

```go
package localstack_s3

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
    "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type BucketResource struct {
    client *S3Client
}

var _ resource.Resource = &BucketResource{}

func NewBucketResource() resource.Resource {
    return &BucketResource{}
}

type BucketModel struct {
    ID       types.String `tfsdk:"id"`
    Name     types.String `tfsdk:"name"`
    Password types.String `tfsdk:"password"`
}

// Schema avec PlanModifier et ApplyModifier
func (r *BucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "ID du bucket généré par l'API",
            },
            "name": schema.StringAttribute{
                Required: true,
                PlanModifiers: []planmodifier.String{
                    stringplanmodifier.RequiresReplace(), // Si le nom change -> delete/create
                },
                Description: "Nom du bucket",
            },
            "password": schema.StringAttribute{
                Required: true,
                Sensitive: true,
                ApplyModifiers: []planmodifier.String{
                    stringplanmodifier.MaskValue(), // Masque le mot de passe avant apply
                },
                Description: "Mot de passe associé au bucket (masqué)",
            },
        },
    }
}

// Configure le client S3
func (r *BucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData != nil {
        r.client = req.ProviderData.(*S3Client)
    }
}

// ===== Create =====
func (r *BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan BucketModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    _, err := r.client.Client.CreateBucket(ctx, &s3.CreateBucketInput{
        Bucket: aws.String(plan.Name.ValueString()),
    })
    if err != nil {
        resp.Diagnostics.AddError("Erreur création bucket", err.Error())
        return
    }

    plan.ID = types.StringValue(plan.Name.ValueString())
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// ===== Read =====
func (r *BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state BucketModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    _, err := r.client.Client.HeadBucket(ctx, &s3.HeadBucketInput{
        Bucket: aws.String(state.Name.ValueString()),
    })
    if err != nil {
        resp.State.RemoveResource(ctx)
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// ===== Update =====
func (r *BucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan BucketModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

    // Ici, password est masqué par ApplyModifier avant d’être utilisé
    fmt.Println("Password envoyé à l'API (masqué par Terraform state):", plan.Password.ValueString())

    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// ===== Delete =====
func (r *BucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state BucketModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    _, err := r.client.Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
        Bucket: aws.String(state.Name.ValueString()),
    })
    if err != nil {
        resp.Diagnostics.AddError("Erreur suppression bucket", err.Error())
    }
}
```

---

# **2. Explications clés**

| Attribut   | PlanModifier        | ApplyModifier | Effet                                                        |
| ---------- | ------------------- | ------------- | ------------------------------------------------------------ |
| `name`     | `RequiresReplace()` | —             | Si le nom change, Terraform supprime et recrée le bucket     |
| `password` | —                   | `MaskValue()` | Masque la valeur dans l’état et les logs avant l’application |

* **PlanModifier** agit **avant le plan**, influence la décision Terraform de recréer ou mettre à jour.
* **ApplyModifier** agit **juste avant l’appel de Create/Update**, utile pour sécuriser ou transformer la valeur.

---

# **3. Exemple Terraform pour tester**

```hcl
provider "localstack_s3" {
  endpoint = "http://localhost:4566"
  region   = "us-east-1"
}

resource "localstack_s3_bucket" "demo" {
  name     = "bucket-demo"
  password = "supersecret"
}
```

* `terraform plan` montrera un `+` si le bucket n’existe pas.
* Si tu changes `name = "bucket-demo-2"`, Terraform planifiera un `-/+` (delete + create).
* Le `password` sera masqué dans le state et les logs grâce à `ApplyModifier`.

---

Si tu veux, je peux te faire **un petit diagramme illustrant exactement le rôle de PlanModifier et ApplyModifier avec ce bucket**, ce qui rend le flux très visuel.

Veux‑tu que je fasse ça ?

Parfait ! Voici un **diagramme visuel simplifié** illustrant exactement comment **PlanModifier** et **ApplyModifier** interviennent dans le cycle Terraform pour notre bucket S3 custom :

```
Utilisateur
   |
   v
Fichier Terraform
   name = "bucket-demo"
   password = "supersecret"
   |
   v
Terraform lit la configuration
   |
   v
===============================
        PLAN PHASE
===============================
PlanModifier agit ici
- name : RequiresReplace()
  => Si le nom change, Terraform planifie une suppression + création
- password : aucun
   |
Terraform calcule le diff
- Nouvelle ressource ? + create
- Changement de nom ? -/+ (delete + create)
- Autres champs : update ou pas
   |
   v
Diff affiché à l’utilisateur
===============================

Terraform APPLY
===============================
ApplyModifier agit ici
- password : MaskValue()
  => La valeur est transformée / masquée avant d’être envoyée à l’API
- name : valeur finale envoyée
   |
Appel des fonctions du provider
- Create(ctx, plan) ou Update(ctx, plan)
- API S3 / LocalStack reçoit :
    - name = "bucket-demo"
    - password = "supersecret" (masqué dans state/log)
   |
   v
État mis à jour dans terraform.tfstate
===============================
```

---

### **Résumé clé**

* **PlanModifier**

  * Agit pendant `terraform plan`
  * Influence le **diff** et la décision Terraform (`create`, `update`, `delete`)
  * Exemple : `RequiresReplace()` sur `name` force delete/create si le nom change

* **ApplyModifier**

  * Agit pendant `terraform apply`, juste avant `Create`/`Update`
  * Transforme, masque ou valide la valeur avant de l’envoyer à l’API
  * Exemple : `MaskValue()` sur `password` masque la valeur dans le state/logs

---
