Voici un exemple concret de déploiement d'un cluster OpenShift sur AWS à l'aide de Terraform, en utilisant le fournisseur officiel Red Hat Cloud Services (RHCS). Ce modèle est adapté à un environnement de production et peut être intégré dans vos pipelines CI/CD.

---

### 🛠️ Prérequis

Avant de commencer, assurez-vous d'avoir :

* Un compte AWS avec les permissions nécessaires.
* Le [fournisseur Terraform Red Hat Cloud Services (RHCS)](https://registry.terraform.io/providers/terraform-redhat/rhcs/latest/docs) installé.
* Le [fournisseur Terraform AWS](https://registry.terraform.io/providers/hashicorp/aws/latest/docs) configuré.
* Un jeton d'authentification ROSA (Red Hat OpenShift Service on AWS).

---

### 📁 Structure du projet

Voici la structure des fichiers Terraform pour déployer un cluster ROSA :

```shell
terraform-cluster/
├── main.tf
├── variables.tf
├── outputs.tf
├── provider.tf
└── rosa.tf
```

---

### 🧱 Contenu des fichiers

#### `provider.tf`

Ce fichier configure les fournisseurs nécessaires :

```hcl
provider "aws" {
  region = var.aws_region
}

provider "rhcs" {
  token = var.rosa_token
}
```

#### `variables.tf`

Définissez les variables utilisées dans votre configuration :

```hcl
variable "aws_region" {
  type    = string
  default = "us-east-2"
}

variable "rosa_token" {
  type    = string
  default = "<votre_token_ROSA>"
}

variable "cluster_name" {
  type    = string
  default = "mon-cluster-rosa"
}
```

#### `rosa.tf`

Ce fichier contient la configuration spécifique à ROSA :

```hcl
resource "rhcs_rosa_cluster" "mon_cluster" {
  name     = var.cluster_name
  region   = var.aws_region
  version  = "4.14"
  channel  = "stable-4.14"

  control_plane {
    instance_type = "m5.xlarge"
    size          = 3
  }

  compute {
    instance_type = "m5.xlarge"
    size          = 3
  }

  networking {
    vpc_id = "<votre_vpc_id>"
  }

  identity_provider {
    type = "oidc"
    config {
      client_id     = "<votre_client_id>"
      client_secret = "<votre_client_secret>"
      issuer_url    = "<votre_issuer_url>"
    }
  }
}
```

#### `outputs.tf`

Affichez les informations importantes après le déploiement :

```hcl
output "cluster_id" {
  value = rhcs_rosa_cluster.mon_cluster.id
}
```

---

### 🚀 Déploiement

1. Initialisez le projet Terraform :

   ```bash
   terraform init
   ```

2. Appliquez la configuration :

   ```bash
   terraform apply
   ```

3. Suivez les instructions pour valider le plan et lancer le déploiement.

---

### 📚 Ressources supplémentaires

* Documentation officielle de ROSA avec Terraform : ([docs.redhat.com][1])
* Blog AWS sur la création de clusters ROSA avec Terraform : ([Amazon Web Services, Inc.][2])
* Exemple de déploiement d'OpenShift sur AWS avec Terraform : ([GitHub][3])

---

Si vous souhaitez un exemple pour une autre plateforme cloud ou un cas d'utilisation spécifique, n'hésitez pas à me le faire savoir.

[1]: https://docs.redhat.com/en/documentation/red_hat_openshift_service_on_aws/4/html/install_clusters/creating-a-rosa-cluster-using-terraform?utm_source=chatgpt.com "Chapter 3. Creating a ROSA cluster using Terraform"
[2]: https://aws.amazon.com/blogs/containers/build-rosa-clusters-with-terraform/?utm_source=chatgpt.com "Build ROSA Clusters with Terraform | Containers"
[3]: https://github.com/dwmkerr/terraform-aws-openshift?utm_source=chatgpt.com "dwmkerr/terraform-aws-openshift: Create infrastructure with ..."

Pour déployer un cluster OpenShift en utilisant un *OpenShift Manager* avec Terraform, vous pouvez exploiter le fournisseur Terraform officiel de Red Hat Cloud Services (RHCS). Ce fournisseur permet de gérer des clusters OpenShift Service on AWS (ROSA) de manière déclarative, y compris la création de clusters, la gestion des pools de machines et l'intégration avec des fournisseurs d'identité externes.

---

### 🧱 Prérequis

Avant de commencer, assurez-vous d'avoir :

* Un compte AWS avec les permissions nécessaires.
* Le [fournisseur Terraform Red Hat Cloud Services (RHCS)](https://registry.terraform.io/providers/terraform-redhat/rhcs/latest/docs) installé.
* Le [fournisseur Terraform AWS](https://registry.terraform.io/providers/hashicorp/aws/latest/docs) configuré.
* Un jeton d'authentification ROSA (Red Hat OpenShift Service on AWS).
* Les rôles IAM nécessaires créés dans votre compte AWS.

---

### 📁 Structure du projet

Voici la structure des fichiers Terraform pour déployer un cluster ROSA :

```shell
terraform-cluster/
├── main.tf
├── variables.tf
├── outputs.tf
├── provider.tf
└── rosa.tf
```

---

### 🛠️ Contenu des fichiers

#### `provider.tf`

Ce fichier configure les fournisseurs nécessaires :

```hcl
provider "aws" {
  region = var.aws_region
}

provider "rhcs" {
  token = var.rosa_token
}
```

#### `variables.tf`

Définissez les variables utilisées dans votre configuration :

```hcl
variable "aws_region" {
  type    = string
  default = "us-east-2"
}

variable "rosa_token" {
  type    = string
  default = "<votre_token_ROSA>"
}

variable "cluster_name" {
  type    = string
  default = "mon-cluster-rosa"
}
```

#### `rosa.tf`

Ce fichier contient la configuration spécifique à ROSA :

```hcl
resource "rhcs_rosa_cluster" "mon_cluster" {
  name     = var.cluster_name
  region   = var.aws_region
  version  = "4.14"
  channel  = "stable-4.14"

  control_plane {
    instance_type = "m5.xlarge"
    size          = 3
  }

  compute {
    instance_type = "m5.xlarge"
    size          = 3
  }

  networking {
    vpc_id = "<votre_vpc_id>"
  }

  identity_provider {
    type = "oidc"
    config {
      client_id     = "<votre_client_id>"
      client_secret = "<votre_client_secret>"
      issuer_url    = "<votre_issuer_url>"
    }
  }
}
```

#### `outputs.tf`

Affichez les informations importantes après le déploiement :

```hcl
output "cluster_id" {
  value = rhcs_rosa_cluster.mon_cluster.id
}
```

---

### 🚀 Déploiement

1. Initialisez le projet Terraform :

   ```bash
   terraform init
   ```

2. Appliquez la configuration :

   ```bash
   terraform apply
   ```

3. Suivez les instructions pour valider le plan et lancer le déploiement.

---

### 📚 Ressources supplémentaires

* Documentation officielle de ROSA avec Terraform : ([docs.redhat.com][1])
* Blog AWS sur la création de clusters ROSA avec Terraform : ([Red Hat Developer][2])
* Exemple de déploiement d'OpenShift sur AWS avec Terraform : ([GitHub][3])

---

[1]: https://docs.redhat.com/en/documentation/red_hat_openshift_service_on_aws/4/html/install_clusters/creating-a-rosa-cluster-using-terraform?utm_source=chatgpt.com "Chapter 3. Creating a ROSA cluster using Terraform"
[2]: https://developers.redhat.com/blog/2024/11/27/performance-terraform-provider-manage-openshift-fleets?utm_source=chatgpt.com "Performance of Terraform provider to manage Openshift ..."
[3]: https://github.com/RedHatGov/terraform-packet-openshift?utm_source=chatgpt.com "RedHatGov/terraform-packet-openshift"
