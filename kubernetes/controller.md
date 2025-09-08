Un **controller Kubernetes** (ou **contrôleur**) est un composant essentiel de Kubernetes qui surveille l'état de ton cluster et agit pour le faire correspondre à l'état désiré défini dans les objets de configuration (comme les fichiers YAML).

---

### 🔧 Définition simple

Un **controller** est un **programme en boucle** qui :

1. **Observe** l'état actuel du cluster.
2. **Compare** cet état à l'état désiré (défini dans des objets Kubernetes comme `Deployment`, `Service`, etc.).
3. **Agit** pour faire converger l'état actuel vers l'état désiré (par exemple, créer ou supprimer des pods, ajuster des ressources...).

---

### 📦 Exemples de controllers intégrés

Voici quelques controllers intégrés à Kubernetes :

| Controller                 | Rôle                                                                                                   |
| -------------------------- | ------------------------------------------------------------------------------------------------------ |
| **Deployment Controller**  | S’assure qu’un nombre donné de pods sont en cours d’exécution selon un `Deployment`.                   |
| **ReplicaSet Controller**  | Maintient le nombre de copies d’un pod (`replicas`) défini dans un `ReplicaSet`.                       |
| **Job Controller**         | Gère l'exécution de tâches ponctuelles (`Job`).                                                        |
| **DaemonSet Controller**   | S’assure qu’un pod spécifique s’exécute sur **chaque** nœud (utile pour les agents système).           |
| **StatefulSet Controller** | Gère des applications avec des identités stables et du stockage persistant (bases de données par ex.). |

---

### 🔁 Exemple concret

Tu crées un objet `Deployment` avec 3 réplicas. Le controller `Deployment` :

* Vérifie s'il y a bien 3 pods en cours d'exécution.
* S’il n’y en a que 2, il en crée un nouveau.
* S’il y en a 4, il en supprime un.
* Si un pod crashe, il en redéploie un.

---

### 🔨 Controllers personnalisés (Custom Controllers)

Kubernetes te permet aussi de créer des **controllers personnalisés**, souvent associés à des **Custom Resource Definitions (CRD)**. Cela te permet d’étendre Kubernetes pour gérer des ressources spécifiques à ton application.

Ex : Un controller custom qui déploie automatiquement une base PostgreSQL quand tu crées une ressource `MyDatabase`.

---

### 📌 Résumé

Un **controller Kubernetes** :

* Surveille l’état d’une ressource.
* Interagit avec l’API Kubernetes.
* Prend des décisions pour rapprocher l’état actuel de l’état désiré.
* Peut être intégré ou personnalisé (via CRD + opérateur).

Parfait ! Voici une **explication claire et complète d’un controller custom Kubernetes**, étape par étape, avec un petit exemple pour t'aider à visualiser.

---

## 🎯 Objectif

Créer un **controller custom** qui va surveiller une **Custom Resource (CR)** nommée `Website`, et à chaque fois qu’un `Website` est créé, le controller va **créer un Pod** qui sert un site statique.

---

## 📘 Étape 1 : Définir une **Custom Resource Definition (CRD)**

Un CRD permet à Kubernetes de reconnaître un **nouveau type de ressource**. Par exemple, on ajoute une ressource nommée `Website`.

```yaml
# website-crd.yaml
apiVersion: apiextensions.k8s.io/v1
kind: CustomResourceDefinition
metadata:
  name: websites.example.com
spec:
  group: example.com
  versions:
    - name: v1
      served: true
      storage: true
      schema:
        openAPIV3Schema:
          type: object
          properties:
            spec:
              type: object
              properties:
                domain:
                  type: string
  scope: Namespaced
  names:
    plural: websites
    singular: website
    kind: Website
    shortNames:
      - web
```

---

## 📘 Étape 2 : Créer une instance `Website`

```yaml
# website.yaml
apiVersion: example.com/v1
kind: Website
metadata:
  name: mysite
spec:
  domain: "mysite.local"
```

---

## ⚙️ Étape 3 : Écrire un controller custom en Go (simplifié)

Tu peux utiliser [Kubebuilder](https://book.kubebuilder.io/) ou [Operator SDK](https://sdk.operatorframework.io/) pour simplifier, mais voici une version manuelle et simplifiée.

### ✏️ Exemple en Go (simplifié)

```go
package main

import (
    "context"
    "fmt"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/tools/clientcmd"
    "k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
    config, _ := clientcmd.BuildConfigFromFlags("", "~/.kube/config")
    dynClient, _ := dynamic.NewForConfig(config)

    websiteGVR := schema.GroupVersionResource{
        Group:    "example.com",
        Version:  "v1",
        Resource: "websites",
    }

    for {
        list, _ := dynClient.Resource(websiteGVR).Namespace("default").List(context.TODO(), metav1.ListOptions{})
        for _, w := range list.Items {
            domain := w.Object["spec"].(map[string]interface{})["domain"].(string)
            fmt.Println("Detected Website:", w.GetName(), "->", domain)

            // Ici tu pourrais créer un Pod, Deployment ou autre
        }
        time.Sleep(10 * time.Second)
    }
}
```

---

## 📦 Ce que ferait un vrai controller

Dans un projet réel, le controller :

* Utilise des **informer** ou un **controller-runtime** pour surveiller les événements.
* Réagit automatiquement aux créations/modifications/suppressions.
* Applique de la **reconciliation** (comme les controllers natifs de Kubernetes).
* Gère les erreurs, les conflits, les statuts, etc.

---

## 🧰 Outils recommandés

* [`kubebuilder`](https://book.kubebuilder.io/): pour générer la structure de projet.
* [`controller-runtime`](https://pkg.go.dev/sigs.k8s.io/controller-runtime): la base pour écrire un controller facilement.
* [`Operator SDK`](https://sdk.operatorframework.io/): pour créer des opérateurs prêts pour la production.

---

## 📌 Résumé

| Élément           | Rôle                                                       |
| ----------------- | ---------------------------------------------------------- |
| CRD               | Définit une nouvelle ressource Kubernetes (`Website`).     |
| Custom Resource   | Une instance de ce nouveau type.                           |
| Custom Controller | Code qui observe et agit en fonction des objets `Website`. |

---

## controller specifique
Très bien ! Prenons un **controller spécifique intégré à Kubernetes** et examinons-le **en profondeur**. Un bon exemple pédagogique est le **Deployment Controller**, car il est largement utilisé et illustre très bien le modèle de **reconciliation**.

---

## 🎯 Objectif du Deployment Controller

Le **Deployment Controller** a pour mission de :

* S'assurer qu'un nombre désiré de **pods** tournent avec la bonne version de l'application.
* Mettre à jour les pods progressivement (rolling update).
* Recréer automatiquement les pods si nécessaire.

---

## 🧠 Comment ça marche

Quand tu appliques un objet `Deployment`, voici ce que le controller fait :

1. **Observation** : il regarde l'état du `Deployment` (dans l'API Server).
2. **Comparaison** :

   * Il compare le nombre réel de pods (via ReplicaSet) avec le nombre désiré (`.spec.replicas`).
   * Il regarde si l'image ou la configuration a changé.
3. **Action** :

   * Il crée ou met à jour un **ReplicaSet** qui gère les pods.
   * Il supprime progressivement les anciens pods, et en crée des nouveaux (rolling update).
   * Il suit les progrès de la mise à jour dans `.status`.

---

## 📦 Exemple concret

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
        - name: nginx
          image: nginx:1.25
```

Quand tu appliques ce manifeste :

* Le Deployment Controller crée un **ReplicaSet** pour gérer ces 3 pods.
* Ce ReplicaSet va lui-même créer les 3 pods avec l’image `nginx:1.25`.

Si tu modifies l’image (par ex. `nginx:1.26`), le controller :

* Crée un **nouveau ReplicaSet** avec la nouvelle image.
* Fait un **rolling update** : il supprime les pods de l’ancienne version un par un et en crée de nouveaux.
* Si un pod échoue, il peut revenir en arrière (`rollback` possible).

---

## 🔁 Cycle de Reconciliation

Le Deployment Controller suit ce cycle :

```
[User] --> [Deployment object] --> [Controller]
                                   |
                                   |--> Compare state
                                   |--> Update or create ReplicaSet
                                   |--> Create/delete Pods
```

Et il continue à surveiller et corriger les écarts si des pods sont supprimés manuellement, ou plantent.

---

## 🧰 Code source du Deployment Controller (si tu veux creuser)

Le code du Deployment Controller se trouve ici :

* 📁 [https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/deployment](https://github.com/kubernetes/kubernetes/tree/master/pkg/controller/deployment)

Il est écrit en Go, et utilise un **reconcile loop** basée sur des **informers**, une **workqueue**, et des **listers** pour l’efficacité.

---

## 📌 Résumé

| Élément               | Rôle                                                 |
| --------------------- | ---------------------------------------------------- |
| Deployment            | Spécifie l’état désiré d’un groupe de pods.          |
| Deployment Controller | Fait en sorte que cet état soit atteint et maintenu. |
| ReplicaSet            | Géré par le controller, il crée les pods.            |
| Rolling update        | Remplacement progressif des anciens pods.            |

---

Parfait, plongeons dans le fonctionnement du **Horizontal Pod Autoscaler (HPA)**, un controller Kubernetes très puissant et souvent utilisé en production.

---

## 🌡️ Horizontal Pod Autoscaler (HPA) – Vue d’ensemble

Le **HPA** est un **controller natif** de Kubernetes qui **ajuste automatiquement le nombre de pods** d’un déploiement (ou autre objet scalable comme un `StatefulSet`) en fonction de **l'utilisation réelle des ressources** (CPU, mémoire, ou custom metrics).

---

## 🎯 Objectif

Maintenir **le bon nombre de pods** pour répondre à la charge :
➡️ **Scaler horizontalement** (ajouter ou retirer des pods)

---

## 📦 Exemple concret

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: nginx-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: nginx-deployment
  minReplicas: 2
  maxReplicas: 10
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 50
```

### 🧠 Explication :

* **minReplicas** : au minimum 2 pods
* **maxReplicas** : au maximum 10 pods
* **metrics** : si la moyenne d’utilisation CPU dépasse 50%, le HPA augmente le nombre de pods

---

## 🔁 Cycle de fonctionnement du HPA Controller

1. **Observe** l’état de la cible (`Deployment`, `ReplicaSet`, etc.)
2. **Collecte les métriques** via le **Metrics Server** (CPU, mémoire, ou Prometheus avec custom metrics)
3. **Compare** l'utilisation mesurée avec la cible (`averageUtilization`)
4. **Décide** si on doit scale up/down
5. **Met à jour** le champ `.spec.replicas` du `Deployment` cible

---

## 🔧 Prérequis

* Le **Metrics Server** doit être déployé dans ton cluster :

  ```bash
  kubectl top pods
  ```

* Si cette commande ne fonctionne pas, tu dois installer `metrics-server` :

  ```bash
  kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
  ```

---

## 📈 Formule de scaling (simplifiée)

Le nombre de réplicas est ajusté avec cette formule :

```
desiredReplicas = ceil(currentReplicas * currentMetric / targetMetric)
```

Exemple :

* 3 pods utilisent en moyenne 80% de CPU
* Target = 50%
  → `ceil(3 * 80 / 50) = 5` → le HPA scale à 5 pods

---

## 🔍 Que fait le HPA controller concrètement ?

* Il lit le champ `.status` du `Deployment` cible
* Il consulte les métriques d’utilisation via `metrics.k8s.io` API
* Il met à jour le champ `.spec.replicas` de l’objet cible
* Il répète ce processus toutes les **15 secondes par défaut**

---

## 🧠 Architecture interne

```text
+------------------------+
| HorizontalPodAutoscaler|
+------------------------+
           |
           v
  +---------------------+     +---------------------+
  | Metrics API (CPU)   | --> | Custom Metrics API  |
  +---------------------+     +---------------------+
           |
           v
     +--------------+
     | Target Scale |
     | (Deployment) |
     +--------------+
```

---

## 🧪 Pour tester

Tu peux simuler de la charge avec `stress` :

```bash
kubectl run -i --tty load-generator --image=busybox /bin/sh
> while true; do wget -q -O- http://nginx-service; done
```

Puis observer :

```bash
kubectl get hpa
kubectl describe hpa nginx-hpa
```

---

## 📌 Résumé

| Élément        | Rôle                                                |
| -------------- | --------------------------------------------------- |
| HPA Controller | Surveille les métriques et ajuste le nombre de pods |
| Metrics Server | Fournit les métriques CPU/mémoire                   |
| ScaleTargetRef | L’objet que le HPA gère (Deployment, etc.)          |
| Scaling        | Horizontal : variation du nombre de pods            |

---

## example de controller natif kubernetes sur la base d'un deployment
Bien sûr ! Voici un **exemple simplifié en Go** basé sur le **code d’un controller natif Kubernetes**, inspiré du fonctionnement du `DeploymentController`.

Ce code reflète les **principes réels** utilisés dans les controllers de Kubernetes : informer, queue, reconcile, etc.
💡 Kubernetes lui-même est écrit en Go, et ses controllers suivent ce modèle.

---

## 🏗️ Objectif du code

Créer un controller qui :

* Surveille les **Pods** dans le cluster
* Affiche un message dans les logs quand un Pod est **ajouté**

---

## 📦 Exemple minimal inspiré du code source Kubernetes

### 📝 Fichier `main.go`

```go
package main

import (
    "context"
    "fmt"
    "time"

    corev1 "k8s.io/api/core/v1"
    "k8s.io/apimachinery/pkg/fields"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/client-go/informers"
    "k8s.io/client-go/kubernetes"
    "k8s.io/client-go/rest"
    "k8s.io/client-go/tools/cache"
)

func main() {
    // Configuration in-cluster (ou locale via ~/.kube/config avec clientcmd)
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err)
    }

    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err)
    }

    factory := informers.NewSharedInformerFactory(clientset, time.Minute)

    // Création de l’informer Pod
    podInformer := factory.Core().V1().Pods().Informer()

    // Gestion des événements
    podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
        AddFunc: func(obj interface{}) {
            pod := obj.(*corev1.Pod)
            fmt.Printf("🔍 Nouveau Pod détecté : %s/%s\n", pod.Namespace, pod.Name)
        },
        UpdateFunc: func(oldObj, newObj interface{}) {
            // Optionnel : handle update
        },
        DeleteFunc: func(obj interface{}) {
            // Optionnel : handle delete
        },
    })

    // Démarrage de l’informer
    stopCh := make(chan struct{})
    defer close(stopCh)

    fmt.Println("▶️ Lancement du controller...")
    factory.Start(stopCh)

    // Attendre que le cache soit synchronisé
    if ok := cache.WaitForCacheSync(stopCh, podInformer.HasSynced); !ok {
        panic("❌ Impossible de synchroniser le cache")
    }

    // Bloquer pour laisser le controller tourner
    <-stopCh
}
```

---

## 🧠 Ce que fait ce controller

* Connecte-toi à l’API Kubernetes (`client-go`)
* Utilise un **informer** (très utilisé dans les controllers natifs) pour suivre les événements sur les Pods
* Exécute une fonction quand un Pod est **ajouté** (via `AddFunc`)
* Démarre un **cache local** qui suit les objets Kubernetes

---

## 📦 À propos des controllers natifs

Les controllers Kubernetes suivent tous ce **modèle d’architecture** :

```text
[Informer] --> [Queue] --> [Reconcile Function]
```

1. **Informer** : surveille les événements (add/update/delete)
2. **Queue** : envoie les événements à traiter
3. **Reconciler** : applique la logique pour corriger ou adapter l’état du cluster

Tu retrouveras cette architecture dans tous les controllers natifs comme :

* Deployment Controller : `pkg/controller/deployment`
* ReplicaSet Controller : `pkg/controller/replicaset`
* Job Controller : `pkg/controller/job`

---

## 🔎 Tu veux aller plus loin ?

* Voir le code source du [Deployment Controller sur GitHub](https://github.com/kubernetes/kubernetes/blob/master/pkg/controller/deployment/deployment_controller.go)
* Ou je peux t’écrire un **controller de type natif** plus avancé avec `workqueue`, `rate limiting`, et gestion d’erreurs — comme dans le vrai code Kubernetes.

