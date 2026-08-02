Oui — tu peux construire un endpoint avec Huma qui exécute des commandes Terraform **tout en filtrant/sécurisant les arguments**. Il faut être très strict sur la validation pour éviter toute injection de commande.

Voici un exemple simple et sécurisé 👇

---

## 🔧 Exemple d’endpoint Huma pour exécuter Terraform

### 1. Installation

```bash
go get github.com/danielgtaylor/huma/v2
```

---

### 2. Code Go (endpoint sécurisé)

```go
package main

import (
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humago"
)

type TerraformInput struct {
	Action    string   `json:"action" doc:"Terraform action (plan/apply/destroy)"`
	Directory string   `json:"directory" doc:"Working directory"`
	Vars      []string `json:"vars" doc:"List of -var key=value"`
}

type TerraformOutput struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
}

// whitelist des actions autorisées
var allowedActions = map[string]bool{
	"plan":    true,
	"apply":   true,
	"destroy": true,
}

// validation stricte des variables
func validateVars(vars []string) ([]string, error) {
	var safeVars []string

	for _, v := range vars {
		// simple validation: key=value sans caractères dangereux
		if !strings.Contains(v, "=") {
			return nil, fmt.Errorf("invalid var format: %s", v)
		}

		if strings.ContainsAny(v, "&|;$`") {
			return nil, fmt.Errorf("forbidden characters in var: %s", v)
		}

		safeVars = append(safeVars, "-var", v)
	}

	return safeVars, nil
}

func main() {
	app := humago.New()

	huma.Register(app, huma.Operation{
		OperationID: "run-terraform",
		Method:      "POST",
		Path:        "/terraform",
	}, func(ctx context.Context, input *TerraformInput) (*TerraformOutput, error) {

		// validation action
		if !allowedActions[input.Action] {
			return nil, fmt.Errorf("action not allowed")
		}

		// validation vars
		argsVars, err := validateVars(input.Vars)
		if err != nil {
			return nil, err
		}

		// construction des args
		args := []string{input.Action}
		args = append(args, argsVars...)

		// exec terraform (sans shell = plus safe)
		cmd := exec.CommandContext(ctx, "terraform", args...)
		cmd.Dir = input.Directory

		out, err := cmd.CombinedOutput()

		resp := &TerraformOutput{
			Output: string(out),
		}

		if err != nil {
			resp.Error = err.Error()
		}

		return resp, nil
	})

	humago.ListenAndServe(app, ":8888")
}
```

---

## 🔐 Points de sécurité importants

Ne prends surtout pas ça à la légère :

### ✅ À faire absolument

* whitelist des actions (`plan`, `apply`, etc.)
* **NE JAMAIS utiliser `sh -c`**
* validation stricte des arguments (`-var`)
* isoler les répertoires autorisés
* utiliser un user système restreint

### ❌ À éviter

* passer des arguments bruts depuis l’utilisateur
* autoriser `-var-file` sans contrôle
* exposer directement `terraform apply` sans validation

---

## 🚀 Améliorations possibles

* Ajouter une whitelist de dossiers Terraform
* Support de `-var-file` avec chemins contrôlés
* Timeout avec `context.WithTimeout`
* Logs structurés
* Mode dry-run obligatoire avant apply

---

## 💡 Variante plus safe (recommandée en prod)

Au lieu d’exécuter Terraform directement :

* queue (Redis / worker)
* sandbox (Docker)
* ou utiliser Terraform Cloud API

---

Si tu veux, je peux te faire :

* une version avec **Docker sandbox**
* une version avec **auth + RBAC**
* ou une API complète prête pour prod 👍

___________________________________________

est-ce que tu peux m'ajouter un test unitaire de la route ?

___________________________________________

Oui — et c’est une excellente idée, surtout pour ce type d’endpoint qui exécute des commandes système 👍

On va tester ton endpoint Huma **sans exécuter réellement Terraform** (sinon ce serait lent + dangereux).
👉 On va injecter une fonction mock pour remplacer l’appel à Terraform.

---

## 🔧 1. Petite refacto pour rendre testable

Ajoute une variable pour pouvoir mocker l’exécution :

```go
var execCommand = exec.CommandContext
```

Puis remplace :

```go
cmd := exec.CommandContext(ctx, "terraform", args...)
```

par :

```go
cmd := execCommand(ctx, "terraform", args...)
```

---

## 🧪 2. Test unitaire du endpoint

```go
package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/humago"
)

// mock exec
func fakeExecCommand(ctx context.Context, name string, args ...string) *exec.Cmd {
	// on appelle un helper test à la place
	cs := []string{"-test.run=TestHelperProcess", "--", name}
	cs = append(cs, args...)

	cmd := exec.CommandContext(ctx, os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

// helper pour simuler terraform
func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}

	// simule une sortie terraform
	println("Terraform plan successful")
	os.Exit(0)
}

func setupTestApp() *humago.App {
	app := humago.New()

	huma.Register(app, huma.Operation{
		OperationID: "run-terraform",
		Method:      "POST",
		Path:        "/terraform",
	}, terraformHandler) // factorise ton handler

	return app
}

func TestTerraformEndpoint(t *testing.T) {
	// remplace execCommand par mock
	execCommand = fakeExecCommand
	defer func() { execCommand = exec.CommandContext }()

	app := setupTestApp()

	// payload
	body := map[string]interface{}{
		"action":    "plan",
		"directory": ".",
		"vars":      []string{"env=dev"},
	}

	jsonBody, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/terraform", bytes.NewReader(jsonBody))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	if !bytes.Contains(w.Body.Bytes(), []byte("Terraform plan successful")) {
		t.Fatalf("unexpected response body: %s", w.Body.String())
	}
}
```

---

## 🧱 3. (Option recommandé) Factoriser le handler

Pour que le test marche proprement :

```go
func terraformHandler(ctx context.Context, input *TerraformInput) (*TerraformOutput, error) {
	// 👉 ton code existant ici
}
```

---

## ✅ Ce que teste ce code

* appel HTTP POST `/terraform`
* validation du JSON
* passage dans ton handler
* exécution simulée de Terraform
* vérification de la réponse

---

## 🔒 Bonus : tests à ajouter

Tu peux renforcer avec :

* ❌ action interdite (`applyyyy`)
* ❌ var invalide (`env;rm -rf`)
* ❌ directory interdite
* ⏱ timeout

---

Si tu veux, je peux te faire une version avec :

* tests table-driven (idiomatique Go)
* coverage complète
* ou mock plus propre avec interface (clean architecture)


___________________________________________

est-ce que tu peux m'ajouter une whitelist sur les actions ?

___________________________________________

Oui — et tu es sur la bonne piste : **la whitelist des actions est indispensable** pour éviter qu’un utilisateur fasse passer des commandes imprévues à Terraform.

Tu avais déjà une base, mais on peut la rendre plus robuste et explicite 👇

---

## 🔐 1. Whitelist stricte des actions

On définit une whitelist claire + normalisation :

```go
var allowedActions = map[string]struct{}{
	"plan":    {},
	"apply":   {},
	"destroy": {},
}
```

---

## 🧼 2. Fonction de validation dédiée

```go
func validateAction(action string) (string, error) {
	action = strings.ToLower(strings.TrimSpace(action))

	if _, ok := allowedActions[action]; !ok {
		return "", fmt.Errorf("action '%s' is not allowed", action)
	}

	return action, nil
}
```

---

## 🔧 3. Utilisation dans ton handler Huma

Remplace la validation actuelle par :

```go
action, err := validateAction(input.Action)
if err != nil {
	return nil, err
}
```

Puis :

```go
args := []string{action}
```

---

## 🚫 4. Option encore plus sécurisée (recommandée)

Tu peux **restreindre `apply` et `destroy`** :

```go
var allowedActions = map[string]struct{}{
	"plan": {}, // seul autorisé par défaut
}
```

Ou bien exiger une confirmation :

```go
type TerraformInput struct {
	Action    string   `json:"action"`
	Confirm   bool     `json:"confirm"`
	Directory string   `json:"directory"`
	Vars      []string `json:"vars"`
}
```

```go
if (action == "apply" || action == "destroy") && !input.Confirm {
	return nil, fmt.Errorf("confirmation required for %s", action)
}
```

---

## 🧠 5. Bonus : empêcher les flags dangereux

Même avec une whitelist d’actions, certains flags sont risqués :

* `-auto-approve`
* `-target`
* `-lock=false`

👉 Tu peux aussi filtrer explicitement :

```go
var forbiddenFlags = []string{
	"-auto-approve",
	"-lock=false",
}

func containsForbidden(args []string) error {
	for _, arg := range args {
		for _, forbidden := range forbiddenFlags {
			if arg == forbidden {
				return fmt.Errorf("forbidden flag: %s", arg)
			}
		}
	}
	return nil
}
```

---

## ✅ Résultat

Avec ça :

* seules les actions autorisées passent
* format nettoyé (trim + lowercase)
* possibilité d’ajouter des règles métier (confirmations, RBAC…)

---

Si tu veux aller plus loin, je peux te proposer :

* une whitelist **par rôle utilisateur (RBAC)**
* ou un système où `apply` passe uniquement via une **queue + validation humaine** (très utilisé en prod)

____________________________________________

je veux bien un système où apply passe uniquement via une queue + validation humaine (très utilisé en prod)

____________________________________________

Très bonne idée — c’est exactement ce qui est fait en prod pour sécuriser Terraform 👍
On va transformer ton endpoint Huma en système **asynchrone avec validation humaine**.

---

# 🧠 Architecture (simple et robuste)

### 🔁 Flux

1. `POST /terraform` → crée une **job**
2. Job stockée en base (status = `pending`)
3. `plan` exécuté automatiquement
4. Un humain valide via `/approve/:id`
5. Worker exécute `apply`

---

# 🧱 1. Modèle de Job

```go id="r8a6kq"
type JobStatus string

const (
	StatusPending  JobStatus = "pending"
	StatusPlanned  JobStatus = "planned"
	StatusApproved JobStatus = "approved"
	StatusRunning  JobStatus = "running"
	StatusDone     JobStatus = "done"
	StatusFailed   JobStatus = "failed"
)

type TerraformJob struct {
	ID        string
	Action    string
	Directory string
	Vars      []string
	Status    JobStatus
	PlanOut   string
	ApplyOut  string
	CreatedAt time.Time
}
```

---

# 📨 2. Endpoint → création de job

👉 `apply` ne lance plus rien directement

```go id="u3c5eh"
func createJob(ctx context.Context, input *TerraformInput) (*TerraformJob, error) {
	action, err := validateAction(input.Action)
	if err != nil {
		return nil, err
	}

	job := &TerraformJob{
		ID:        uuid.NewString(),
		Action:    action,
		Directory: input.Directory,
		Vars:      input.Vars,
		Status:    StatusPending,
		CreatedAt: time.Now(),
	}

	saveJob(job)

	// enqueue
	jobQueue <- job.ID

	return job, nil
}
```

---

# ⚙️ 3. Worker (plan auto)

```go id="9y9g6p"
func worker() {
	for jobID := range jobQueue {
		job := getJob(jobID)

		if job.Action == "plan" || job.Action == "apply" {
			out, err := runTerraform("plan", job)

			if err != nil {
				job.Status = StatusFailed
				job.PlanOut = err.Error()
			} else {
				job.Status = StatusPlanned
				job.PlanOut = out
			}

			saveJob(job)
		}
	}
}
```

---

# ✅ 4. Endpoint de validation humaine

```go id="3u8n3r"
func approveJob(ctx context.Context, input *struct {
	ID string `path:"id"`
}) (*TerraformJob, error) {

	job := getJob(input.ID)

	if job.Status != StatusPlanned {
		return nil, fmt.Errorf("job not ready for approval")
	}

	job.Status = StatusApproved
	saveJob(job)

	// enqueue pour apply
	jobQueue <- job.ID

	return job, nil
}
```

---

# 🚀 5. Worker (apply après validation)

```go id="v4tfp2"
func worker() {
	for jobID := range jobQueue {
		job := getJob(jobID)

		switch job.Status {

		case StatusPending:
			out, err := runTerraform("plan", job)
			if err != nil {
				job.Status = StatusFailed
			} else {
				job.Status = StatusPlanned
				job.PlanOut = out
			}

		case StatusApproved:
			job.Status = StatusRunning

			out, err := runTerraform("apply", job)
			if err != nil {
				job.Status = StatusFailed
				job.ApplyOut = err.Error()
			} else {
				job.Status = StatusDone
				job.ApplyOut = out
			}
		}

		saveJob(job)
	}
}
```

---

# 🔐 6. Sécurité critique

### 🔒 Obligatoire

* RBAC sur `/approve` (admin uniquement)
* logs complets (audit trail)
* stockage immutable du plan
* timeout sur exec

### 💡 Très recommandé

* envoyer une notif (Slack / email)
* afficher le `plan` avant validation
* signer les jobs (anti-tampering)

---

# 🧰 7. Version production

Au lieu d’un simple channel Go :

### Queue réelle

* Redis + workers
* RabbitMQ
* Apache Kafka

### Stockage

* PostgreSQL (jobs)
* S3 (plans volumineux)

---

# 🔥 Bonus : UX réaliste

Flow typique :

1. Dev appelle API → job créé
2. Bot Slack poste :

   > "Terraform plan ready for approval"
3. Reviewer clique :

   * ✅ Approve
   * ❌ Reject
4. Apply lancé automatiquement

---

# 🧠 Alternative encore plus clean

👉 déléguer à :

* Terraform Cloud
* ou Atlantis (GitOps)

---

# ✅ Résultat

Avec ce système :

* aucun `apply` direct ❌
* validation humaine obligatoire ✅
* traçabilité complète ✅
* modèle scalable (queue + workers) ✅

---

Si tu veux, je peux te générer :

* une version **PostgreSQL + migrations**
* une intégration **Slack approval**
* ou une version **Docker sandbox ultra sécurisée**

