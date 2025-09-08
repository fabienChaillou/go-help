✅ Voici le squelette complet du projet Kubebuilder avec :

### 📁 Structure

* `main.go` : point d'entrée du manager
* `api/v1/website_types.go` : définition de la ressource `Website`
* `controllers/website_controller.go` : logique de reconciliation

### 🚀 Ce que fait ce controller

Dès qu’une ressource `Website` est créée :

* Un Pod nginx est lancé dans le namespace
* Le nom du Pod est stocké dans `.status.createdPod`

Souhaites-tu maintenant :

* un exemple de fichier `config/samples/website.yaml` ?
* ou les manifests de déploiement (`kustomize`) ?

```go
// main.go - Point d'entrée du manager
package main

import (
	"flag"
	"os"

	webv1 "github.com/example/website-operator/api/v1"
	"github.com/example/website-operator/controllers"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(webv1.AddToScheme(scheme))
}

func main() {
	var metricsAddr string
	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme: scheme,
		MetricsBindAddress: metricsAddr,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.WebsiteReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Website")
		os.Exit(1)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// api/v1/website_types.go
package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/scheme"
)

// WebsiteSpec définit les champs désirés
type WebsiteSpec struct {
	Domain string `json:"domain"`
}

// WebsiteStatus définit les champs observés
type WebsiteStatus struct {
	CreatedPod string `json:"createdPod,omitempty"`
}

// Website est le Schema API principal
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Website struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WebsiteSpec   `json:"spec,omitempty"`
	Status WebsiteStatus `json:"status,omitempty"`
}

// WebsiteList contient une liste de Website
// +kubebuilder:object:root=true
type WebsiteList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Website `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Website{}, &WebsiteList{})
}

// controllers/website_controller.go
package controllers

import (
	context "context"
	corev1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	webv1 "github.com/example/website-operator/api/v1"
)

type WebsiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *WebsiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var website webv1.Website
	if err := r.Get(ctx, req.NamespacedName, &website); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "website-" + website.Name,
			Namespace: website.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  "nginx",
				Image: "nginx",
				Ports: []corev1.ContainerPort{{ContainerPort: 80}},
			}},
		},
	}

	if err := r.Create(ctx, pod); err != nil && !apierrors.IsAlreadyExists(err) {
		return ctrl.Result{}, err
	}

	website.Status.CreatedPod = pod.Name
	if err := r.Status().Update(ctx, &website); err != nil {
		logger.Error(err, "unable to update Website status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *WebsiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&webv1.Website{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
```
