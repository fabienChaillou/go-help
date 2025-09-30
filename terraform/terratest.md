Ok, parfait ! On va mettre en place **des tests d’intégration avec Terratest** pour ton provider Terraform local.
Terratest est idéal pour exécuter un `terraform init` / `terraform apply` / `terraform destroy` complet dans un test Go, et vérifier que tout se comporte comme prévu.

---

# **1. Présentation rapide de Terratest**

Terratest permet :

* De **lancer des commandes Terraform** (`init`, `apply`, `destroy`) depuis un test Go.
* De **vérifier l’état final des ressources**.
* De gérer le **cycle de vie complet** : déploiement puis suppression.

C'est un outil complémentaire aux tests d’acceptation Terraform.

---

# **2. Installation**

Dans ton projet, installe la lib :

```bash
go get github.com/gruntwork-io/terratest/modules/terraform
```

Vérifie que tu as Terraform CLI installé localement (`terraform version`).

---

# **3. Structure du projet avec Terratest**

On ajoute un dossier `test/terratest/` :

```
terraform-provider-example/
├── go.mod
├── internal/
│   ├── provider.go
│   ├── resource_foo.go
│   └── client/
│       └── fake_client.go
├── tests/                 # tests unitaires & acceptation
│   ├── provider_test.go
│   └── acceptance_test.go
└── terratest/              # tests avec Terratest
    ├── main.tf             # config Terraform de test
    └── terratest_test.go   # test Go avec Terratest
```

---

# **4. Code Terraform pour le test local**

Dans `terratest/main.tf` :

```hcl
terraform {
  required_providers {
    example = {
      source  = "registry.terraform.io/username/example"
      version = "0.0.1"
    }
  }
}

provider "example" {
  api_token = "terratest-token"
}

resource "example_foo" "test" {
  value = "Hello from Terratest"
}
```

> ⚠️ Le provider doit être installé localement dans :
>
> ```
> ~/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/linux_amd64/
> ```
>
> Comme vu précédemment.

---

# **5. Test avec Terratest**

Dans `terratest/terratest_test.go` :

```go
package terratest

import (
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformExampleFoo(t *testing.T) {
	t.Parallel()

	// Config pour Terratest
	options := &terraform.Options{
		TerraformDir: ".", // dossier contenant main.tf
		NoColor:      true,
	}

	// Nettoyage automatique après le test
	defer terraform.Destroy(t, options)

	// Lancement init + apply
	terraform.InitAndApply(t, options)

	// Récupère la valeur finale
	value := terraform.Output(t, options, "example_foo_value")

	assert.Equal(t, "Hello from Terratest", value)
}
```

---

# **6. Ajouter une sortie dans le provider Terraform**

Pour que Terratest récupère la valeur, on ajoute un `output` dans `terratest/main.tf` :

```hcl
output "example_foo_value" {
  value = example_foo.test.value
}
```

---

# **7. Exécuter le test**

Dans le dossier racine du projet :

```bash
cd terratest
go test -v
```

Exemple de sortie attendue :

```
=== RUN   TestTerraformExampleFoo
Terraform init and apply
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

--- PASS: TestTerraformExampleFoo (5.89s)
PASS
ok      terraform-provider-example/terratest   6.123s
```

---

# **8. Scénarios supplémentaires**

Tu peux enrichir tes tests Terratest avec :

### 8.1 Vérification que l’état est bien ce que tu veux

```go
terraform.AssertOutput(t, options, "example_foo_value", "Hello from Terratest")
```

---

### 8.2 Test de mise à jour (`terraform apply` deux fois)

Tu peux simuler un changement dans la configuration :

```go
func TestTerraformExampleFoo_Update(t *testing.T) {
	t.Parallel()

	options := &terraform.Options{
		TerraformDir: ".",
	}

	defer terraform.Destroy(t, options)

	// Premier déploiement
	terraform.InitAndApply(t, options)

	// Mettre à jour la variable `value`
	options.Vars = map[string]interface{}{
		"value": "Updated value",
	}

	// Second apply
	terraform.Apply(t, options)

	value := terraform.Output(t, options, "example_foo_value")
	assert.Equal(t, "Updated value", value)
}
```

---

# **9. Makefile pour automatiser**

On peut ajouter un Makefile pour automatiser les étapes :

```makefile
PLUGIN_DIR=~/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/linux_amd64
BINARY=terraform-provider-example

build:
	go build -o $(BINARY)

install: build
	mkdir -p $(PLUGIN_DIR)
	mv $(BINARY) $(PLUGIN_DIR)/

test:
	cd terratest && go test -v

clean:
	rm -rf .terraform terraform.tfstate*
```

---

# **10. Flux complet**

| Étape                       | Commande                            |
| --------------------------- | ----------------------------------- |
| Compiler le provider        | `make build`                        |
| Installer le provider local | `make install`                      |
| Init + Apply Terraform      | `terraform init && terraform apply` |
| Lancer les tests Terratest  | `make test`                         |
| Nettoyer                    | `make clean`                        |

---

# **11. Résumé**

Avec ce setup, tu as :

* Un **provider custom** compilé et installé localement.
* Des **tests unitaires** (`tests/provider_test.go`).
* Des **tests d’acceptation Terraform** (`tests/acceptance_test.go`).
* Des **tests d’intégration avec Terratest** (`terratest/terratest_test.go`).
* Un **Makefile** pour automatiser la compilation, l’installation et le test.

---

Souhaites-tu que je te prépare une version avec **Docker** pour isoler totalement les tests Terraform et Terratest ? 🚀

Parfait, on va encapsuler tout le workflow dans **Docker** pour isoler ton provider et tes tests.
Avec cette approche, tu n'auras plus besoin d’installer Terraform ou Go directement sur ta machine, tout sera géré dans des conteneurs.

---

# **1. Objectifs de la version Docker**

Avec cette solution, tu pourras :

* Compiler et exécuter ton **provider Terraform custom** dans un conteneur.
* Installer le provider dans l’arborescence locale attendue par Terraform (`~/.terraform.d/plugins/`).
* Lancer tes **tests unitaires**, **tests d’acceptation**, et **tests Terratest** depuis Docker.
* Déployer un environnement Terraform en local **sans rien installer**.

---

# **2. Structure finale du projet**

Voici l’arborescence finale :

```
terraform-provider-example/
├── Dockerfile
├── docker-compose.yml
├── Makefile
├── go.mod
├── internal/
│   ├── provider.go
│   ├── resource_foo.go
│   └── client/
│       └── fake_client.go
├── main.go
├── terratest/
│   ├── main.tf
│   └── terratest_test.go
└── tests/
    ├── provider_test.go
    └── acceptance_test.go
```

---

# **3. Dockerfile**

Ce Dockerfile crée une image avec Go, Terraform et Terratest.

**Dockerfile**

```Dockerfile
FROM golang:1.22 AS builder

# Installer Terraform CLI
ENV TERRAFORM_VERSION=1.8.5
RUN apt-get update && apt-get install -y wget unzip git make && \
    wget https://releases.hashicorp.com/terraform/${TERRAFORM_VERSION}/terraform_${TERRAFORM_VERSION}_linux_amd64.zip && \
    unzip terraform_${TERRAFORM_VERSION}_linux_amd64.zip -d /usr/local/bin && \
    terraform version

# Préparer l'environnement de build
WORKDIR /app

# Copier les fichiers du projet
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compiler le provider
RUN go build -o terraform-provider-example .

# Créer le répertoire du plugin
RUN mkdir -p /root/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/linux_amd64 \
    && mv terraform-provider-example /root/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/linux_amd64/

# Exposer terraform et le provider
ENV PATH="/usr/local/bin:/root/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/linux_amd64:${PATH}"

CMD ["bash"]
```

---

# **4. docker-compose.yml**

Fichier pour orchestrer facilement le conteneur de test.

**docker-compose.yml**

```yaml
version: "3.9"
services:
  terraform-test:
    build: .
    container_name: terraform-provider-example
    working_dir: /app
    volumes:
      - .:/app
    environment:
      TF_ACC: "1"   # Active les tests d'acceptation
    command: bash
```

---

# **5. main.tf pour Terratest**

Dans `terratest/main.tf` :

```hcl
terraform {
  required_providers {
    example = {
      source  = "registry.terraform.io/username/example"
      version = "0.0.1"
    }
  }
}

provider "example" {
  api_token = "docker-token"
}

resource "example_foo" "test" {
  value = "Hello from Docker Terratest"
}

output "example_foo_value" {
  value = example_foo.test.value
}
```

---

# **6. Test Terratest**

Dans `terratest/terratest_test.go` :

```go
package terratest

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestTerraformExampleFoo(t *testing.T) {
	t.Parallel()

	// Configuration Terratest
	options := &terraform.Options{
		TerraformDir: ".", // Dossier avec main.tf
		NoColor:      true,
	}

	// Nettoyage après test
	defer terraform.Destroy(t, options)

	// Init et Apply
	terraform.InitAndApply(t, options)

	// Vérification de la valeur retournée
	value := terraform.Output(t, options, "example_foo_value")
	assert.Equal(t, "Hello from Docker Terratest", value)
}
```

---

# **7. Makefile**

Pour simplifier la compilation, l'installation et les tests.

**Makefile**

```makefile
PLUGIN_DIR=/root/.terraform.d/plugins/registry.terraform.io/username/example/0.0.1/linux_amd64
BINARY=terraform-provider-example

# Build du provider local
build:
	docker-compose run --rm terraform-test go build -o $(BINARY) .

# Installer le provider dans le conteneur
install: build
	docker-compose run --rm terraform-test sh -c "mkdir -p $(PLUGIN_DIR) && mv $(BINARY) $(PLUGIN_DIR)/"

# Lancer les tests unitaires et d'acceptation
test:
	docker-compose run --rm terraform-test go test ./tests -v

# Lancer Terratest
terratest:
	docker-compose run --rm terraform-test sh -c "cd terratest && go test -v"

# Terraform apply dans le conteneur
apply:
	docker-compose run --rm terraform-test sh -c "cd terratest && terraform init && terraform apply -auto-approve"

# Terraform destroy dans le conteneur
destroy:
	docker-compose run --rm terraform-test sh -c "cd terratest && terraform destroy -auto-approve"

# Nettoyer les fichiers Terraform
clean:
	rm -rf terratest/.terraform terratest/terraform.tfstate*
```

---

# **8. Lancement pas à pas**

## 8.1. Builder l'image Docker

```bash
docker-compose build
```

---

## 8.2. Compiler et installer le provider

```bash
make install
```

Cela :

* Compile le binaire `terraform-provider-example`.
* L’installe dans `/root/.terraform.d/plugins/...` du conteneur.

---

## 8.3. Lancer un `terraform apply` dans le conteneur

```bash
make apply
```

Exemple de sortie :

```
Terraform will perform the following actions:

  # example_foo.test will be created
  + resource "example_foo" "test" {
      + id    = (known after apply)
      + value = "Hello from Docker Terratest"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

---

## 8.4. Lancer Terratest

```bash
make terratest
```

Sortie attendue :

```
=== RUN   TestTerraformExampleFoo
Terraform init and apply
Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

--- PASS: TestTerraformExampleFoo (6.23s)
PASS
ok      terraform-provider-example/terratest   6.248s
```

---

## 8.5. Nettoyer

```bash
make destroy
make clean
```

---

# **9. Flux complet avec Docker**

| Étape                  | Commande               |
| ---------------------- | ---------------------- |
| Builder l'image Docker | `docker-compose build` |
| Compiler provider      | `make build`           |
| Installer provider     | `make install`         |
| Terraform apply        | `make apply`           |
| Lancer Terratest       | `make terratest`       |
| Nettoyer les fichiers  | `make clean`           |

---

# **10. Points clés**

* **Isolation complète** : aucun outil installé sur ta machine hôte.
* **Cycle complet automatisé** avec Makefile.
* **Tests Terratest** directement intégrés dans le conteneur.
* Compatible **CI/CD** (GitHub Actions, GitLab CI, etc.).
* Facile à déployer sur d'autres machines grâce au `docker-compose.yml`.

---

Super — voici une présentation claire et pratique (en français) de **Terratest** quand on veut tester des configurations Terraform qui utilisent **un provider personnalisé local**. Je donne la logique, les alternatives pour “placer” le provider local, un exemple complet de test Go avec Terratest (incluant comment pointer Terraform vers le provider local via un fichier de config CLI), et des conseils/précautions.

# 1) Rappels rapides

* **Terratest** est une bibliothèque Go (Gruntwork) pour écrire des tests d’infrastructure (terraform, k8s, packer, etc.). Elle lance `terraform init/apply/destroy` et permet d’écrire assertions en Go. ([terratest.gruntwork.io][1])
* Terraform charge les providers comme des plugins. Pour utiliser un provider personnalisé **local**, il faut fournir la structure attendue (emplacement/version/OS-arch) ou configurer Terraform pour qu’il installe les providers à partir d’un mirror local / dossier. ([HashiCorp Developer][2])

# 2) Stratégies pour utiliser un provider local dans des tests

Trois approches courantes :

A. **Installer le binaire du provider dans `~/.terraform.d/plugins/...`** — méthode simple pour dev local (Terraform cherchera là avant d’aller sur le registry). Utile pour tests locaux rapides. ([HashiCorp Discuss][3])

B. **`terraform providers mirror`** — créer un “miroir” sur le filesystem (layout attendu) et pointer Terraform vers ce dossier via la config CLI (`provider_installation` avec `filesystem_mirror`). Robuste, reproductible (bon pour CI/air-gapped). ([HashiCorp Developer][2])

C. **Placer les binaires dans un `plugin-dir` / `-plugin-dir`** (historiquement possible) ou utiliser l’option `TF_CLI_CONFIG_FILE` pour donner un fichier CLI configuré pour chercher dans un dossier local. (La méthode recommandée actuelle : mirror + provider\_installation / \~/.terraform.d/plugins.) ([Stack Overflow][4])

# 3) Exemple pratique : workflow recommandé (mirror + Terratest)

Principe :

1. Compiler votre provider personnalisé en binaire (ex : `terraform-provider-mycorp_v0.1.0_linux_amd64`).
2. Créer un mirror local avec l’arborescence attendue (ou utiliser `terraform providers mirror` pour l’automatiser).
3. Écrire un fichier CLI Terraform (ex : `.terraformrc` ou un fichier temporaire) qui pointe `provider_installation` vers le dossier mirror.
4. Dans votre test Terratest, exporter `TF_CLI_CONFIG_FILE` vers ce fichier CLI, puis lancer `terraform init` / `apply` via Terratest.

### Exemple de layout (mirror)

```
/tmp/provider-mirror/
└── registry.terraform.io
    └── mycorp
        └── myprovider
            └── 0.1.0
                └── linux_amd64
                    └── terraform-provider-myprovider_v0.1.0
```

(Ce layout est celui attendu par `terraform providers mirror` / Terraform.) ([HashiCorp Developer][2])

### Exemple de fichier CLI (provider\_installation)

Nommer ce fichier `tf-cli-config.hcl` (vous pouvez le générer dans le test) :

```hcl
provider_installation {
  filesystem_mirror {
    path    = "/tmp/provider-mirror"
  }
  direct {
    # fallback to registry if not found in mirror
  }
}
```

### Exemple de test Go + Terratest

(assume `go.mod` et dépendances `github.com/gruntwork-io/terratest/modules/terraform`)

```go
package test

import (
  "os"
  "path/filepath"
  "testing"
  "github.com/gruntwork-io/terratest/modules/terraform"
  "github.com/stretchr/testify/assert"
)

func TestMyProviderModule(t *testing.T) {
  t.Parallel()

  // 1) Préparer un fichier CLI temporaire qui pointe vers le mirror local
  cliConfig := `provider_installation {
    filesystem_mirror {
      path = "/tmp/provider-mirror"
    }
    direct {}
  }`
  cliPath := filepath.Join(os.TempDir(), "tf-cli-config.hcl")
  if err := os.WriteFile(cliPath, []byte(cliConfig), 0600); err != nil {
    t.Fatalf("failed to write tf cli config: %v", err)
  }
  // s'assurer de supprimer après
  defer os.Remove(cliPath)

  // 2) Options terraform pour Terratest
  terraformOptions := &terraform.Options{
    // dossier contenant votre example Terraform (main.tf qui requerra mycorp/myprovider)
    TerraformDir: "../examples/myprovider",
    EnvVars: map[string]string{
      // dire à Terraform d'utiliser ce fichier CLI
      "TF_CLI_CONFIG_FILE": cliPath,
      // variables éventuelles pour votre provider/module
      // "MYPROVIDER_ENDPOINT": "http://localhost:1234",
    },
  }

  // 3) Init & Apply via Terratest (et toujours cleanup)
  defer terraform.Destroy(t, terraformOptions)
  terraform.InitAndApply(t, terraformOptions)

  // 4) Assertions : par ex. lire outputs
  out := terraform.Output(t, terraformOptions, "example_output")
  assert.Equal(t, "expected-value", out)
}
```

Remarques :

* `TF_CLI_CONFIG_FILE` est pratique car vous pouvez générer ce fichier **temporaire** depuis le test, ce qui évite de modifier la config globale du runner CI.
* Si vous préférez, au lieu d’un mirror vous pouvez pré-installer le binaire dans `~/.terraform.d/plugins/...` sur l’image runner, mais c’est moins reproductible pour CI sans préparation. ([HashiCorp Discuss][3])

# 4) Alternatives pour le développement & tests rapides

* **Placer le binaire dans `~/.terraform.d/plugins`** (chemin exact dépend OS) : rapide pour dev local mais moins propre pour CI. ([registry.terraform.io][5])
* **Utiliser Localstack / simulateurs** si votre provider parle à un API cloud : utile pour tests isolés (ex. Localstack pour AWS) — on combine souvent Localstack + Terratest. ([Medium][6])

# 5) Pièges / conseils pratiques

* **Versions et nommage** : Terraform s’attend à version dans le required\_providers et au nommage/metadata du provider. Pour du dev local, utilisez une version explicite (ex `0.1.0` ou un numéro élevé `99.0.0` pour dev). ([HashiCorp Discuss][7])
* **OS/arch** : placez le binaire pour la bonne plateforme (linux\_amd64, darwin\_amd64, etc.).
* **Permissions** : assurez-vous que le binaire est exécutable.
* **Nettoyage** : Terratest supprime automatiquement si vous utilisez `defer terraform.Destroy(...)`, mais testez aussi la destruction si votre provider gère états côté remote.
* **Debug** : si `terraform init` ignore votre binaire local, vérifier : fichier CLI utilisé (`TF_CLI_CONFIG_FILE`), arborescence du mirror, et que `required_providers` mentionne la même source (ex `mycorp/myprovider`) que l’arborescence. ([HashiCorp Developer][2])

# 6) Références utiles (lecture / copy-paste)

* Terratest — quick start & examples (Gruntwork). Exemples de tests Terraform : guides et snippets. ([terratest.gruntwork.io][1])
* Terraform — gestion des plugins & `providers mirror` / `provider_installation`. (page officielle `developer.hashicorp.com/terraform/cli/plugins`.) ([HashiCorp Developer][2])
* Discussions HashiCorp + guides sur comment faire du dev local de providers et layout `~/.terraform.d/plugins`. ([HashiCorp Discuss][3])

---

Si tu veux, je peux tout de suite :

* te fournir un **exemple complet et prêt** (repo-style) : `examples/myprovider/main.tf`, le script `build-provider.sh` qui construit le binaire et prépare le mirror, **et** le fichier de test Go complet (prêt à lancer).
* ou générer un **tf-cli-config.hcl** adapté à ton chemin mirror et un fragment de `required_providers` pour `main.tf`.

Dis-moi ce que tu préfères (exemple complet prêt-à-lancer ou juste les snippets).

[1]: https://terratest.gruntwork.io/docs/getting-started/quick-start/?utm_source=chatgpt.com "Quick start - Terratest - Gruntwork"
[2]: https://developer.hashicorp.com/terraform/cli/plugins?utm_source=chatgpt.com "Manage plugins overview - Terraform - HashiCorp Developer"
[3]: https://discuss.hashicorp.com/t/local-provider-installation-not-working-in-terraform-13/21586?utm_source=chatgpt.com "Local Provider Installation not working in Terraform 13"
[4]: https://stackoverflow.com/questions/50944395/use-pre-installed-terraform-plugins-instead-of-downloading-them-with-terraform-i?utm_source=chatgpt.com "Use pre-installed Terraform plugins instead of downloading them ..."
[5]: https://registry.terraform.io/providers/wearespindle/proxmox/latest/docs/guides/installation?utm_source=chatgpt.com "installation | Guides | wearespindle/proxmox - Terraform Registry"
[6]: https://medium.com/%40ledevedeccorentin/terraform-iac-testing-localstack-terratest-9946dafe98b6?utm_source=chatgpt.com "Terraform (IaC) testing: Localstack + Terratest. - Medium"
[7]: https://discuss.hashicorp.com/t/easiest-way-to-use-a-local-custom-provider-with-terraform-0-13/12691?utm_source=chatgpt.com "Easiest way to use a local custom provider with Terraform 0.13"

Parfait ! Voici un **exemple complet et prêt à l’emploi** pour tester un **provider Terraform personnalisé local** avec **Terratest**.
On va construire une petite arborescence avec :

* **Un provider custom minimal** (en Go)
* **Un module Terraform exemple** qui utilise ce provider
* **Un test Terratest en Go** qui :

  * génère un fichier CLI temporaire (`TF_CLI_CONFIG_FILE`)
  * utilise un **mirror local** pour le provider
  * exécute `terraform init` + `apply`
  * fait des assertions sur les outputs
* **Un script `Makefile` ou `build.sh`** pour compiler et préparer le provider.

---

# 1) Arborescence finale

```
my-terraform-provider-test/
├── provider/                     # Code du provider custom
│   ├── main.go
│   └── go.mod
├── terraform/                    # Module Terraform qui l'utilise
│   ├── main.tf
│   └── outputs.tf
├── test/                          # Tests Terratest
│   ├── provider_test.go
│   └── go.mod
└── build.sh                       # Script build + mirror
```

---

# 2) Provider minimal en Go (`provider/main.go`)

Ce provider expose une **ressource simple** `mycorp_hello` qui prend une variable `name` et retourne un message `greeting`.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func main() {
	err := providerserver.Serve(context.Background(), NewProvider, providerserver.ServeOpts{
		Address: "registry.terraform.io/mycorp/myprovider",
	})
	if err != nil {
		log.Fatal(err)
	}
}

func NewProvider() provider.Provider {
	return &helloProvider{}
}

type helloProvider struct{}

func (p *helloProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "mycorp_myprovider"
}

func (p *helloProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

func (p *helloProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {}

func (p *helloProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewHelloResource,
	}
}

func (p *helloProvider) DataSources(_ context.Context) []func() resource.Resource {
	return nil
}

// --------- Resource hello ---------

type helloResource struct{}

func NewHelloResource() resource.Resource {
	return &helloResource{}
}

func (r *helloResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "mycorp_hello"
}

func (r *helloResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"greeting": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *helloResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan struct {
		Name types.String `tfsdk:"name"`
	}
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	greeting := fmt.Sprintf("Hello, %s!", plan.Name.ValueString())

	resp.State.Set(ctx, map[string]any{
		"name":     plan.Name.ValueString(),
		"greeting": greeting,
	})
}

func (r *helloResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Nothing to do
}

func (r *helloResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	r.Create(ctx, resource.CreateRequest{Plan: req.Plan}, resp)
}

func (r *helloResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}
```

### `provider/go.mod`

```go
module mycorp/myprovider

go 1.21

require (
	github.com/hashicorp/terraform-plugin-framework v1.8.0
)
```

---

# 3) Module Terraform qui utilise ce provider (`terraform/main.tf`)

```hcl
terraform {
  required_providers {
    myprovider = {
      source = "mycorp/myprovider"
      version = "0.1.0"
    }
  }
}

provider "myprovider" {}

resource "mycorp_hello" "example" {
  name = "Terratest"
}

output "example_output" {
  value = mycorp_hello.example.greeting
}
```

### `terraform/outputs.tf`

```hcl
output "greeting" {
  value = mycorp_hello.example.greeting
}
```

---

# 4) Test Terratest (`test/provider_test.go`)

```go
package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

func TestMyProvider(t *testing.T) {
	t.Parallel()

	// 1. Générer un fichier CLI Terraform temporaire
	cliConfig := `provider_installation {
	  filesystem_mirror {
	    path = "./.mirror"
	  }
	  direct {}
	}`

	tmpFile := filepath.Join(os.TempDir(), "tf-cli-config.hcl")
	if err := os.WriteFile(tmpFile, []byte(cliConfig), 0600); err != nil {
		t.Fatalf("failed to write tf cli config: %v", err)
	}
	defer os.Remove(tmpFile)

	// 2. Options Terraform
	terraformOptions := &terraform.Options{
		TerraformDir: "../terraform",
		EnvVars: map[string]string{
			"TF_CLI_CONFIG_FILE": tmpFile,
		},
	}

	// Cleanup à la fin
	defer terraform.Destroy(t, terraformOptions)

	// 3. Init & Apply
	terraform.InitAndApply(t, terraformOptions)

	// 4. Assertions
	output := terraform.Output(t, terraformOptions, "example_output")
	assert.Equal(t, "Hello, Terratest!", output)
}
```

### `test/go.mod`

```go
module mycorp/terratest-example

go 1.21

require (
	github.com/gruntwork-io/terratest v0.46.0
	github.com/stretchr/testify v1.8.4
)
```

---

# 5) Script de build et préparation (`build.sh`)

Ce script :

1. Compile le provider
2. Crée l’arborescence du mirror
3. Copie le provider dans le bon dossier

```bash
#!/usr/bin/env bash
set -euo pipefail

VERSION="0.1.0"
OS="linux"
ARCH="amd64"

# Dossier de sortie
MIRROR_DIR="./.mirror/registry.terraform.io/mycorp/myprovider/${VERSION}/${OS}_${ARCH}"
mkdir -p "$MIRROR_DIR"

echo "==> Building provider..."
cd provider
go mod tidy
go build -o "../terraform-provider-myprovider_v${VERSION}"
cd ..

echo "==> Moving binary to mirror..."
mv "terraform-provider-myprovider_v${VERSION}" "$MIRROR_DIR/terraform-provider-myprovider_v${VERSION}"

echo "==> Done. Mirror ready at ./.mirror"
```

---

# 6) Lancement des tests

```bash
# Compiler le provider et préparer le mirror
./build.sh

# Lancer le test Terratest
cd test
go test -v
```

Sortie attendue :

```
=== RUN   TestMyProvider
--- PASS: TestMyProvider (10.32s)
PASS
ok  	mycorp/terratest-example/test	10.345s
```

---

# 7) Explications clés

* **Mirror local (`.mirror`)**
  Terraform trouve le provider grâce au fichier CLI `provider_installation` pointant sur `./.mirror`.

* **Structure `registry.terraform.io/mycorp/myprovider/VERSION/OS_ARCH`**
  Obligatoire pour que Terraform reconnaisse le provider.

* **`TF_CLI_CONFIG_FILE`**
  Le test Go génère ce fichier à la volée pour ne pas dépendre d'une config globale.

* **Isolation totale**
  Pas besoin de publier le provider, tout reste local.

---

# 8) Bonus : intégration CI (GitHub Actions)

Extrait simple pour CI :

```yaml
name: Test Provider
on: [push]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Install Go
        uses: actions/setup-go@v5
        with:
          go-version: 1.21
      - name: Build provider
        run: ./build.sh
      - name: Run tests
        run: cd test && go test -v
```

---

# 9) Résumé

* **`build.sh`** gère le provider local et le mirror.
* **Terraform module** consomme le provider via un `source = "mycorp/myprovider"`.
* **Terratest** pilote `init` / `apply` via un fichier CLI temporaire.
* **Structure reproductible**, parfaite pour CI/CD et développement local.
