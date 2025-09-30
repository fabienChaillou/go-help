Oui, absolument !
Je vais te montrer **un exemple complet d’un provider Terraform custom en Go** qui se connecte à **LocalStack** pour gérer des **buckets S3**, sans dépendre d’AWS directement.

LocalStack est parfait pour tester un provider en local, car il expose une API AWS compatible, et tu pourras tout tester sans coût ni dépendance cloud.

---

# **1. Objectif du provider**

Notre provider `localstack_s3` va :

* Se connecter à **LocalStack** (service S3 mocké).
* Gérer un bucket S3 avec les opérations :

  * **Create** → Créer un bucket.
  * **Read** → Vérifier si le bucket existe.
  * **Update** → (optionnel, ici on changera juste un tag).
  * **Delete** → Supprimer le bucket.

---

# **2. Architecture du projet**

```
terraform-provider-localstack-s3/
├── go.mod
├── go.sum
├── main.go
├── provider.go
├── client.go
├── resource_bucket.go
└── examples/
    └── main.tf
```

---

# **3. Dépendances Go**

On utilise le SDK AWS v2 pour communiquer avec LocalStack.

```bash
go mod init github.com/ton-org/terraform-provider-localstack-s3
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/s3
go get github.com/hashicorp/terraform-plugin-framework
go get github.com/hashicorp/terraform-plugin-log
```

Le fichier `go.mod` ressemblera à ceci :

```go
module github.com/ton-org/terraform-provider-localstack-s3

go 1.22

require (
    github.com/aws/aws-sdk-go-v2 v1.30.0
    github.com/aws/aws-sdk-go-v2/config v1.27.0
    github.com/aws/aws-sdk-go-v2/service/s3 v1.54.0
    github.com/hashicorp/terraform-plugin-framework v1.10.0
    github.com/hashicorp/terraform-plugin-log v0.19.0
)
```

---

# **4. Client S3 (`client.go`)**

Ce fichier gère la connexion à LocalStack en utilisant le SDK AWS.

```go
package localstack_s3

import (
    "context"
    "fmt"
    "net/http"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
    Client *s3.Client
}

func NewS3Client(ctx context.Context, endpoint, region string) (*S3Client, error) {
    if endpoint == "" {
        endpoint = "http://localhost:4566" // endpoint LocalStack par défaut
    }
    if region == "" {
        region = "us-east-1"
    }

    // Config AWS spécifique pour LocalStack
    cfg, err := config.LoadDefaultConfig(ctx,
        config.WithRegion(region),
        config.WithHTTPClient(http.DefaultClient),
        config.WithEndpointResolverWithOptions(aws.EndpointResolverWithOptionsFunc(
            func(service, region string, options ...interface{}) (aws.Endpoint, error) {
                return aws.Endpoint{
                    URL:           endpoint,
                    SigningRegion: region,
                }, nil
            },
        )),
    )
    if err != nil {
        return nil, fmt.Errorf("erreur configuration AWS: %w", err)
    }

    return &S3Client{
        Client: s3.NewFromConfig(cfg),
    }, nil
}
```

---

# **5. Provider (`provider.go`)**

Le provider expose les ressources et configure le client S3.

```go
package localstack_s3

import (
    "context"

    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/provider"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
)

type LocalStackS3Provider struct {
    s3Client *S3Client
}

var _ provider.Provider = &LocalStackS3Provider{}

func New() provider.Provider {
    return &LocalStackS3Provider{}
}

func (p *LocalStackS3Provider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
    resp.TypeName = "localstack_s3"
}

func (p *LocalStackS3Provider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "endpoint": schema.StringAttribute{
                Optional:    true,
                Description: "Endpoint LocalStack (par défaut: http://localhost:4566)",
            },
            "region": schema.StringAttribute{
                Optional:    true,
                Description: "Région AWS simulée (par défaut: us-east-1)",
            },
        },
    }
}

func (p *LocalStackS3Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
    var config struct {
        Endpoint string `tfsdk:"endpoint"`
        Region   string `tfsdk:"region"`
    }

    resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
    if resp.Diagnostics.HasError() {
        return
    }

    client, err := NewS3Client(ctx, config.Endpoint, config.Region)
    if err != nil {
        resp.Diagnostics.AddError("Erreur de connexion", err.Error())
        return
    }

    p.s3Client = client
    resp.DataSourceData = client
    resp.ResourceData = client
}

func (p *LocalStackS3Provider) Resources(_ context.Context) []func() resource.Resource {
    return []func() resource.Resource{
        NewBucketResource,
    }
}

func (p *LocalStackS3Provider) DataSources(_ context.Context) []func() datasource.DataSource {
    return nil
}
```

---

# **6. Ressource Bucket S3 (`resource_bucket.go`)**

Implémentation CRUD complète pour un bucket S3.

```go
package localstack_s3

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
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
    ID   types.String `tfsdk:"id"`
    Name types.String `tfsdk:"name"`
}

func (r *BucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_bucket"
}

func (r *BucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "Nom du bucket",
            },
            "name": schema.StringAttribute{
                Required:    true,
                Description: "Nom du bucket",
            },
        },
    }
}

func (r *BucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData == nil {
        return
    }
    r.client = req.ProviderData.(*S3Client)
}

// ===== CREATE =====
func (r *BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan BucketModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    _, err := r.client.Client.CreateBucket(ctx, &s3.CreateBucketInput{
        Bucket: &plan.Name.ValueString(),
    })
    if err != nil {
        resp.Diagnostics.AddError("Erreur création", fmt.Sprintf("Impossible de créer le bucket: %s", err))
        return
    }

    plan.ID = types.StringValue(plan.Name.ValueString())
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

// ===== READ =====
func (r *BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state BucketModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    _, err := r.client.Client.HeadBucket(ctx, &s3.HeadBucketInput{
        Bucket: &state.Name.ValueString(),
    })
    if err != nil {
        resp.State.RemoveResource(ctx) // Le bucket n'existe plus
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

// ===== UPDATE =====
func (r *BucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    // Pour S3, le nom du bucket ne peut pas être changé.
    // Ici, on pourrait gérer des tags ou des ACL à la place.
    resp.Diagnostics.Append(resp.State.Set(ctx, req.Plan)...)
}

// ===== DELETE =====
func (r *BucketResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
    var state BucketModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
    if resp.Diagnostics.HasError() {
        return
    }

    _, err := r.client.Client.DeleteBucket(ctx, &s3.DeleteBucketInput{
        Bucket: &state.Name.ValueString(),
    })
    if err != nil {
        resp.Diagnostics.AddError("Erreur suppression", fmt.Sprintf("Impossible de supprimer le bucket: %s", err))
    }
}
```

---

# **7. Main (`main.go`)**

Point d’entrée du provider.

```go
package main

import (
    "context"
    "log"

    "github.com/hashicorp/terraform-plugin-framework/providerserver"
    "github.com/ton-org/terraform-provider-localstack-s3/localstack_s3"
)

func main() {
    err := providerserver.Serve(context.Background(), localstack_s3.New, providerserver.ServeOpts{
        Address: "local/demo/localstack_s3",
    })

    if err != nil {
        log.Fatal(err.Error())
    }
}
```

---

# **8. Exemple Terraform (`examples/main.tf`)**

```hcl
terraform {
  required_providers {
    localstack_s3 = {
      source  = "local/demo/localstack_s3"
      version = "0.1.0"
    }
  }
}

provider "localstack_s3" {
  endpoint = "http://localhost:4566" # endpoint LocalStack
  region   = "us-east-1"
}

resource "localstack_s3_bucket" "test" {
  name = "my-test-bucket"
}

output "bucket_name" {
  value = localstack_s3_bucket.test.name
}
```

---

# **9. Démarrer LocalStack**

Lancer LocalStack avec Docker :

```bash
docker run --rm -it -p 4566:4566 localstack/localstack
```

---

# **10. Installation et test**

### Compilation du provider :

```bash
go build -o terraform-provider-localstack_s3
```

### Placement dans le dossier plugins :

```bash
mkdir -p ~/.terraform.d/plugins/local/demo/localstack_s3/0.1.0/linux_amd64
mv terraform-provider-localstack_s3 ~/.terraform.d/plugins/local/demo/localstack_s3/0.1.0/linux_amd64/
```

---

### Initialisation Terraform :

```bash
cd examples
terraform init
```

Sortie attendue :

```
Initializing provider plugins...
- Finding local/demo/localstack_s3 versions matching "0.1.0"...
- Installing local/demo/localstack_s3 v0.1.0...
- Installed local/demo/localstack_s3 v0.1.0 (unauthenticated)
```

---

### Création du bucket :

```bash
terraform apply -auto-approve
```

Sortie attendue :

```
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

Outputs:

bucket_name = "my-test-bucket"
```

---

# **11. Vérification dans LocalStack**

Lister les buckets :

```bash
aws --endpoint-url=http://localhost:4566 s3 ls
```

Résultat attendu :

```
2025-09-17 12:34:56 my-test-bucket
```

---

# **12. Résumé des étapes**

| Étape                 | Commande                                                                                                 |
| --------------------- | -------------------------------------------------------------------------------------------------------- |
| Lancer LocalStack     | `docker run --rm -it -p 4566:4566 localstack/localstack`                                                 |
| Compiler provider     | `go build -o terraform-provider-localstack_s3`                                                           |
| Placer binaire        | `mv terraform-provider-localstack_s3 ~/.terraform.d/plugins/local/demo/localstack_s3/0.1.0/linux_amd64/` |
| Init Terraform        | `terraform init`                                                                                         |
| Créer bucket          | `terraform apply -auto-approve`                                                                          |
| Vérifier avec AWS CLI | `aws --endpoint-url=http://localhost:4566 s3 ls`                                                         |

---

# **13. Améliorations possibles**

* Gérer **les tags** des buckets (`PUT Bucket tagging`).
* Supporter la configuration du **versioning**.
* Ajouter un **data source** pour lister les buckets.
* Intégrer des tests d’intégration avec `terraform-plugin-testing`.
* Automatiser le build et le test via un **Makefile** et GitHub Actions.

Parfait ! On va aller plus loin et transformer l’exemple précédent en un **mini-projet complet et professionnel**.

On va y inclure :

1. **Provider Terraform complet** pour LocalStack S3
2. **Gestion des tags et du versioning des buckets**
3. **Data Source** pour lister les buckets existants
4. **Tests d'intégration** avec `terraform-plugin-testing`
5. **Makefile** pour automatiser le build et les tests
6. **Workflow GitHub Actions** pour CI/CD
7. **Docker Compose** pour démarrer LocalStack facilement
8. **Un exemple Terraform clair**

---

# **1. Arborescence complète**

Voici la structure finale :

```
terraform-provider-localstack-s3/
├── Makefile
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
├── main.go
├── client.go
├── provider.go
├── resource_bucket.go
├── datasource_buckets.go
├── internal/
│   └── test/
│       ├── provider_test.go
│       └── terraform.tfvars
├── .github/
│   └── workflows/
│       └── ci.yml
└── examples/
    ├── main.tf
    └── variables.tf
```

---

# **2. Configuration Go (`go.mod`)**

```go
module github.com/ton-org/terraform-provider-localstack-s3

go 1.22

require (
    github.com/aws/aws-sdk-go-v2 v1.30.0
    github.com/aws/aws-sdk-go-v2/config v1.27.0
    github.com/aws/aws-sdk-go-v2/service/s3 v1.54.0
    github.com/hashicorp/terraform-plugin-framework v1.10.0
    github.com/hashicorp/terraform-plugin-log v0.19.0
    github.com/hashicorp/terraform-plugin-testing v1.7.0
)
```

---

# **3. Docker Compose pour LocalStack**

`docker-compose.yml` pour démarrer LocalStack facilement :

```yaml
version: "3.9"

services:
  localstack:
    image: localstack/localstack:latest
    container_name: localstack
    ports:
      - "4566:4566"
      - "4571:4571"
    environment:
      - SERVICES=s3
      - DEBUG=1
      - DATA_DIR=/tmp/localstack/data
    volumes:
      - "./.localstack:/tmp/localstack"
```

Lancer LocalStack :

```bash
docker-compose up -d
```

---

# **4. Makefile pour automatiser**

```makefile
BINARY_NAME=terraform-provider-localstack_s3
PLUGIN_DIR=$(HOME)/.terraform.d/plugins/local/demo/localstack_s3/0.1.0/linux_amd64

all: build install

build:
	go build -o $(BINARY_NAME)

install: build
	mkdir -p $(PLUGIN_DIR)
	cp $(BINARY_NAME) $(PLUGIN_DIR)/

test:
	go test ./internal/test/... -v

clean:
	rm -f $(BINARY_NAME)
```

---

# **5. Gestion des Tags et Versioning**

On étend notre ressource S3 pour :

* Créer des tags
* Activer ou désactiver le versioning

```go
package localstack_s3

import (
    "context"
    "fmt"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
    "github.com/hashicorp/terraform-plugin-framework/resource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
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
    ID        types.String            `tfsdk:"id"`
    Name      types.String            `tfsdk:"name"`
    Versioned types.Bool              `tfsdk:"versioned"`
    Tags      map[string]types.String `tfsdk:"tags"`
}

func (r *BucketResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_bucket"
}

func (r *BucketResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "id": schema.StringAttribute{
                Computed:    true,
                Description: "ID du bucket",
            },
            "name": schema.StringAttribute{
                Required:    true,
                Description: "Nom du bucket",
            },
            "versioned": schema.BoolAttribute{
                Optional:    true,
                Description: "Activer le versioning du bucket",
            },
            "tags": schema.MapAttribute{
                Optional:    true,
                ElementType: types.StringType,
                Description: "Tags associés au bucket",
            },
        },
    }
}

func (r *BucketResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
    if req.ProviderData != nil {
        r.client = req.ProviderData.(*S3Client)
    }
}

func (r *BucketResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
    var plan BucketModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
    if resp.Diagnostics.HasError() {
        return
    }

    // Création du bucket
    _, err := r.client.Client.CreateBucket(ctx, &s3.CreateBucketInput{
        Bucket: aws.String(plan.Name.ValueString()),
    })
    if err != nil {
        resp.Diagnostics.AddError("Erreur création bucket", err.Error())
        return
    }

    // Versioning
    if plan.Versioned.ValueBool() {
        _, err = r.client.Client.PutBucketVersioning(ctx, &s3.PutBucketVersioningInput{
            Bucket: aws.String(plan.Name.ValueString()),
            VersioningConfiguration: &s3types.VersioningConfiguration{
                Status: s3types.BucketVersioningStatusEnabled,
            },
        })
        if err != nil {
            resp.Diagnostics.AddError("Erreur configuration versioning", err.Error())
            return
        }
    }

    // Tags
    if len(plan.Tags) > 0 {
        var tagSet []s3types.Tag
        for k, v := range plan.Tags {
            tagSet = append(tagSet, s3types.Tag{Key: aws.String(k), Value: aws.String(v.ValueString())})
        }

        _, err = r.client.Client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
            Bucket: aws.String(plan.Name.ValueString()),
            Tagging: &s3types.Tagging{
                TagSet: tagSet,
            },
        })
        if err != nil {
            resp.Diagnostics.AddError("Erreur configuration tags", err.Error())
            return
        }
    }

    plan.ID = types.StringValue(plan.Name.ValueString())
    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *BucketResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
    var state BucketModel
    resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

    _, err := r.client.Client.HeadBucket(ctx, &s3.HeadBucketInput{
        Bucket: aws.String(state.Name.ValueString()),
    })
    if err != nil {
        resp.State.RemoveResource(ctx)
        return
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *BucketResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
    var plan BucketModel
    resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

    if len(plan.Tags) > 0 {
        var tagSet []s3types.Tag
        for k, v := range plan.Tags {
            tagSet = append(tagSet, s3types.Tag{Key: aws.String(k), Value: aws.String(v.ValueString())})
        }

        _, err := r.client.Client.PutBucketTagging(ctx, &s3.PutBucketTaggingInput{
            Bucket: aws.String(plan.Name.ValueString()),
            Tagging: &s3types.Tagging{
                TagSet: tagSet,
            },
        })
        if err != nil {
            resp.Diagnostics.AddError("Erreur update tags", err.Error())
            return
        }
    }

    resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

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

# **6. Data Source : Lister les Buckets**

```go
package localstack_s3

import (
    "context"

    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/hashicorp/terraform-plugin-framework/datasource"
    "github.com/hashicorp/terraform-plugin-framework/schema"
    "github.com/hashicorp/terraform-plugin-framework/types"
)

type BucketsDataSource struct {
    client *S3Client
}

var _ datasource.DataSource = &BucketsDataSource{}

func NewBucketsDataSource() datasource.DataSource {
    return &BucketsDataSource{}
}

type BucketItem struct {
    Name types.String `tfsdk:"name"`
}

type BucketsModel struct {
    Buckets []BucketItem `tfsdk:"buckets"`
}

func (d *BucketsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
    resp.TypeName = req.ProviderTypeName + "_buckets"
}

func (d *BucketsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
    resp.Schema = schema.Schema{
        Attributes: map[string]schema.Attribute{
            "buckets": schema.ListNestedAttribute{
                Computed: true,
                NestedObject: schema.NestedAttributeObject{
                    Attributes: map[string]schema.Attribute{
                        "name": schema.StringAttribute{Computed: true},
                    },
                },
            },
        },
    }
}

func (d *BucketsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
    if req.ProviderData != nil {
        d.client = req.ProviderData.(*S3Client)
    }
}

func (d *BucketsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
    output, err := d.client.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
    if err != nil {
        resp.Diagnostics.AddError("Erreur list buckets", err.Error())
        return
    }

    var buckets []BucketItem
    for _, b := range output.Buckets {
        buckets = append(buckets, BucketItem{Name: types.StringValue(*b.Name)})
    }

    result := BucketsModel{Buckets: buckets}
    resp.Diagnostics.Append(resp.State.Set(ctx, &result)...)
}
```

---

# **7. Exemple Terraform (`examples/main.tf`)**

```hcl
terraform {
  required_providers {
    localstack_s3 = {
      source  = "local/demo/localstack_s3"
      version = "0.1.0"
    }
  }
}

provider "localstack_s3" {
  endpoint = "http://localhost:4566"
  region   = "us-east-1"
}

resource "localstack_s3_bucket" "test" {
  name      = "bucket-demo"
  versioned = true
  tags = {
    environment = "dev"
    owner       = "team1"
  }
}

data "localstack_s3_buckets" "all" {}

output "bucket_names" {
  value = [for b in data.localstack_s3_buckets.all.buckets : b.name]
}
```

---

# **8. Tests avec `terraform-plugin-testing`**

`internal/test/provider_test.go` :

```go
package test

import (
    "testing"

    "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccBucketResource(t *testing.T) {
    resource.Test(t, resource.TestCase{
        PreCheck: func() {},
        Steps: []resource.TestStep{
            {
                Config: `
provider "localstack_s3" {
  endpoint = "http://localhost:4566"
  region   = "us-east-1"
}

resource "localstack_s3_bucket" "test" {
  name = "tf-test-bucket"
}
`,
                Check: resource.ComposeAggregateTestCheckFunc(
                    resource.TestCheckResourceAttr("localstack_s3_bucket.test", "name", "tf-test-bucket"),
                ),
            },
        },
    })
}
```

Lancer les tests :

```bash
make test
```

---

# **9. CI GitHub Actions**

`.github/workflows/ci.yml` :

```yaml
name: CI

on:
  push:
    branches: [ "main" ]
  pull_request:
    branches: [ "main" ]

jobs:
  build-test:
    runs-on: ubuntu-latest
    services:
      localstack:
        image: localstack/localstack
        ports:
          - 4566:4566
        env:
          SERVICES: s3
    steps:
      - uses: actions/checkout@v3
      - name: Set up Go
        uses: actions/setup-go@v4
        with:
          go-version: 1.22
      - run: make test
```

---

# **10. Workflow complet**

| Étape                 | Commande                                         |
| --------------------- | ------------------------------------------------ |
| Démarrer LocalStack   | `docker-compose up -d`                           |
| Compiler provider     | `make build`                                     |
| Installer provider    | `make install`                                   |
| Initialiser Terraform | `terraform init`                                 |
| Créer un bucket       | `terraform apply -auto-approve`                  |
| Vérifier avec AWS CLI | `aws --endpoint-url=http://localhost:4566 s3 ls` |
| Lancer les tests      | `make test`                                      |

---

# **11. Résultat final**

Ce projet te donne :

* Un provider Terraform complet pour LocalStack S3
* Support du versioning et des tags
* Data source pour lister les buckets
* Tests d’intégration automatisés
* CI/CD GitHub Actions
* Workflow Docker Compose pour lancer LocalStack facilement

Avec ça, tu as une **base professionnelle** pour développer et publier un provider Terraform fiable.
