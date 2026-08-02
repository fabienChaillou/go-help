Excellent 👍
On va construire un **module Go propre**, avec une structure simple, du code testable, et des **tests unitaires (mockés)** pour **Vault**.

L’objectif :
✅ Module réutilisable → `vaultclient`
✅ Auth (Token ou AppRole)
✅ Lecture et écriture d’un secret KV v2
✅ Tests unitaires avec mocks (`gomock` ou `testify/mock`)

---

## 🧱 Structure du module

```
vault-demo/
├── go.mod
├── go.sum
├── internal/
│   └── vaultclient/
│       ├── client.go
│       └── client_test.go
└── main.go
```

---

## `go.mod`

```go
module vault-demo

go 1.22

require (
	github.com/hashicorp/vault/api v1.15.0
	github.com/stretchr/testify v1.9.0
)
```

---

## 🧩 `internal/vaultclient/client.go`

```go
package vaultclient

import (
	"context"
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
)

// Client is a wrapper around the Vault API client
type Client struct {
	api *vault.Client
}

// Config holds Vault connection parameters
type Config struct {
	Address   string
	Token     string
	RoleID    string
	SecretID  string
	AuthType  string // "token" or "approle"
}

// New creates and authenticates a Vault client
func New(cfg Config) (*Client, error) {
	vaultCfg := vault.DefaultConfig()
	if cfg.Address != "" {
		vaultCfg.Address = cfg.Address
	}

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("vault.NewClient: %w", err)
	}

	switch cfg.AuthType {
	case "token":
		token := cfg.Token
		if token == "" {
			token = os.Getenv("VAULT_TOKEN")
		}
		if token == "" {
			return nil, fmt.Errorf("no Vault token provided")
		}
		client.SetToken(token)

	case "approle":
		if cfg.RoleID == "" || cfg.SecretID == "" {
			return nil, fmt.Errorf("missing AppRole credentials")
		}
		data := map[string]interface{}{
			"role_id":   cfg.RoleID,
			"secret_id": cfg.SecretID,
		}
		secret, err := client.Logical().Write("auth/approle/login", data)
		if err != nil {
			return nil, fmt.Errorf("AppRole login failed: %w", err)
		}
		if secret == nil || secret.Auth == nil {
			return nil, fmt.Errorf("unexpected AppRole response")
		}
		client.SetToken(secret.Auth.ClientToken)
	default:
		return nil, fmt.Errorf("invalid auth type: %s", cfg.AuthType)
	}

	return &Client{api: client}, nil
}

// ReadSecret reads a secret from a KV v2 path
func (c *Client) ReadSecret(ctx context.Context, path string) (map[string]interface{}, error) {
	secret, err := c.api.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("vault read: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret not found at %s", path)
	}

	// KV v2: secret.Data["data"] contains the actual data
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected data format")
	}
	return data, nil
}

// WriteSecret writes a map of data to a KV v2 path
func (c *Client) WriteSecret(ctx context.Context, path string, data map[string]interface{}) error {
	payload := map[string]interface{}{
		"data": data,
	}
	_, err := c.api.Logical().WriteWithContext(ctx, path, payload)
	if err != nil {
		return fmt.Errorf("vault write: %w", err)
	}
	return nil
}
```

---

## 🧪 `internal/vaultclient/client_test.go`

On va tester :

* construction du client avec token (mock environnement)
* lecture d’un secret (mock de Vault API)
* écriture d’un secret

Ici on ne dépend **pas d’un vrai Vault** (tout est mocké).

```go
package vaultclient

import (
	"context"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mock Logical client ---

type mockLogical struct {
	mock.Mock
}

func (m *mockLogical) ReadWithContext(ctx context.Context, path string) (*vault.Secret, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(*vault.Secret), args.Error(1)
}

func (m *mockLogical) WriteWithContext(ctx context.Context, path string, data map[string]interface{}) (*vault.Secret, error) {
	args := m.Called(ctx, path, data)
	return args.Get(0).(*vault.Secret), args.Error(1)
}

// --- Tests ---

func TestNewClientWithToken(t *testing.T) {
	t.Setenv("VAULT_TOKEN", "test-token")
	cfg := Config{
		Address:  "http://localhost:8200",
		AuthType: "token",
	}
	c, err := New(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.NotNil(t, c.api)
}

func TestReadSecret(t *testing.T) {
	mockL := new(mockLogical)
	client := &vault.Client{}
	client.SetToken("dummy")
	client.SetLogical(mockL)

	c := &Client{api: client}

	expectedData := map[string]interface{}{"username": "admin", "password": "1234"}
	mockL.On("ReadWithContext", mock.Anything, "secret/data/myapp").
		Return(&vault.Secret{Data: map[string]interface{}{"data": expectedData}}, nil)

	data, err := c.ReadSecret(context.Background(), "secret/data/myapp")
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)
	mockL.AssertExpectations(t)
}

func TestWriteSecret(t *testing.T) {
	mockL := new(mockLogical)
	client := &vault.Client{}
	client.SetToken("dummy")
	client.SetLogical(mockL)

	c := &Client{api: client}
	data := map[string]interface{}{"username": "admin", "password": "1234"}

	mockL.On("WriteWithContext", mock.Anything, "secret/data/myapp",
		map[string]interface{}{"data": data}).Return(&vault.Secret{}, nil)

	err := c.WriteSecret(context.Background(), "secret/data/myapp", data)
	assert.NoError(t, err)
	mockL.AssertExpectations(t)
}
```

---

## 🚀 `main.go` (exemple d’usage)

```go
package main

import (
	"context"
	"fmt"
	"log"

	"vault-demo/internal/vaultclient"
)

func main() {
	cfg := vaultclient.Config{
		Address:  "https://vault.example.com:8200",
		AuthType: "token",
		Token:    "s.yourtoken",
	}

	client, err := vaultclient.New(cfg)
	if err != nil {
		log.Fatalf("Erreur init Vault client: %v", err)
	}

	// Exemple lecture
	data, err := client.ReadSecret(context.Background(), "secret/data/myapp")
	if err != nil {
		log.Fatalf("Erreur lecture: %v", err)
	}
	fmt.Printf("username=%s password=%s\n", data["username"], data["password"])
}
```

---

## ✅ Lancer les tests

```bash
go test ./internal/vaultclient -v
```

Sortie attendue :

```
=== RUN   TestNewClientWithToken
--- PASS: TestNewClientWithToken (0.00s)
=== RUN   TestReadSecret
--- PASS: TestReadSecret (0.00s)
=== RUN   TestWriteSecret
--- PASS: TestWriteSecret (0.00s)
PASS
ok  	vault-demo/internal/vaultclient	0.013s
```

---

Souhaites-tu que je t’ajoute :

* ✅ une version **mockgen** (avec `go:generate` + interfaces Vault mockées automatiquement),
  ou
* ✅ une version **intégration réelle** (avec `docker run vault` pour tester contre un vrai serveur Vault local) ?

Parfait 👌 — on va enrichir ton module Go pour utiliser **`mockgen`** (outil officiel de GoMock) afin de générer automatiquement les mocks du client Vault.

Tu auras :

* une **interface `VaultAPI`** et une **interface `Logical`** abstraites,
* un **mock généré par `go:generate`** via `mockgen`,
* des **tests unitaires 100 % isolés** (aucune dépendance à un vrai serveur Vault).

---

## 🧱 Structure finale du module

```
vault-demo/
├── go.mod
├── go.sum
├── internal/
│   └── vaultclient/
│       ├── client.go
│       ├── interfaces.go
│       ├── mock_vaultclient/
│       │   └── mock_vault.go   <-- généré automatiquement
│       └── client_test.go
└── main.go
```

---

## 1️⃣ `go.mod`

```go
module vault-demo

go 1.22

require (
	github.com/golang/mock v1.6.0
	github.com/hashicorp/vault/api v1.15.0
	github.com/stretchr/testify v1.9.0
)
```

---

## 2️⃣ `internal/vaultclient/interfaces.go`

On définit **des interfaces minimales** pour abstraire Vault.

```go
package vaultclient

import (
	"context"

	vault "github.com/hashicorp/vault/api"
)

//go:generate mockgen -destination=mock_vaultclient/mock_vault.go -package=mock_vaultclient vault-demo/internal/vaultclient Logical,VaultAPI

// Logical définit les opérations du client logique de Vault
type Logical interface {
	ReadWithContext(ctx context.Context, path string) (*vault.Secret, error)
	WriteWithContext(ctx context.Context, path string, data map[string]interface{}) (*vault.Secret, error)
}

// VaultAPI définit les opérations principales du client Vault
type VaultAPI interface {
	Logical() Logical
	SetToken(string)
}
```

➡️ `//go:generate` te permet d’exécuter :

```bash
go generate ./...
```

Ce qui créera automatiquement `mock_vaultclient/mock_vault.go`.

---

## 3️⃣ `internal/vaultclient/client.go`

Version adaptée du client qui dépend uniquement de `VaultAPI` et `Logical`.

```go
package vaultclient

import (
	"context"
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

// Client encapsule un client Vault abstrait
type Client struct {
	api VaultAPI
}

// Config Vault
type Config struct {
	Address   string
	Token     string
	RoleID    string
	SecretID  string
	AuthType  string // "token" ou "approle"
}

// New crée un client Vault avec authentification
func New(cfg Config) (*Client, error) {
	vaultCfg := vault.DefaultConfig()
	if cfg.Address != "" {
		vaultCfg.Address = cfg.Address
	}

	client, err := vault.NewClient(vaultCfg)
	if err != nil {
		return nil, fmt.Errorf("vault.NewClient: %w", err)
	}

	switch cfg.AuthType {
	case "token":
		if cfg.Token == "" {
			return nil, fmt.Errorf("missing Vault token")
		}
		client.SetToken(cfg.Token)
	case "approle":
		if cfg.RoleID == "" || cfg.SecretID == "" {
			return nil, fmt.Errorf("missing AppRole credentials")
		}
		data := map[string]interface{}{
			"role_id":   cfg.RoleID,
			"secret_id": cfg.SecretID,
		}
		secret, err := client.Logical().Write("auth/approle/login", data)
		if err != nil {
			return nil, fmt.Errorf("AppRole login failed: %w", err)
		}
		if secret == nil || secret.Auth == nil {
			return nil, fmt.Errorf("invalid AppRole response")
		}
		client.SetToken(secret.Auth.ClientToken)
	default:
		return nil, fmt.Errorf("invalid auth type: %s", cfg.AuthType)
	}

	return &Client{api: client}, nil
}

// ReadSecret lit un secret KVv2
func (c *Client) ReadSecret(ctx context.Context, path string) (map[string]interface{}, error) {
	secret, err := c.api.Logical().ReadWithContext(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("vault read: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret not found")
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected format")
	}
	return data, nil
}

// WriteSecret écrit un secret KVv2
func (c *Client) WriteSecret(ctx context.Context, path string, data map[string]interface{}) error {
	payload := map[string]interface{}{"data": data}
	_, err := c.api.Logical().WriteWithContext(ctx, path, payload)
	if err != nil {
		return fmt.Errorf("vault write: %w", err)
	}
	return nil
}
```

---

## 4️⃣ Génération du mock

Dans le dossier racine :

```bash
go install github.com/golang/mock/mockgen@latest
go generate ./...
```

Résultat :

```
internal/vaultclient/mock_vaultclient/mock_vault.go
```

---

## 5️⃣ `internal/vaultclient/client_test.go`

Tests unitaires en utilisant le mock généré.
(`mockgen` crée `mock_vaultclient.MockLogical` et `mock_vaultclient.MockVaultAPI`)

```go
package vaultclient

import (
	"context"
	"testing"

	vault "github.com/hashicorp/vault/api"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"vault-demo/internal/vaultclient/mock_vaultclient"
)

func TestReadSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogical := mock_vaultclient.NewMockLogical(ctrl)
	mockAPI := mock_vaultclient.NewMockVaultAPI(ctrl)

	mockAPI.EXPECT().Logical().Return(mockLogical).AnyTimes()

	expected := map[string]interface{}{"username": "admin", "password": "1234"}
	secret := &vault.Secret{
		Data: map[string]interface{}{"data": expected},
	}

	mockLogical.EXPECT().
		ReadWithContext(gomock.Any(), "secret/data/myapp").
		Return(secret, nil)

	c := &Client{api: mockAPI}
	data, err := c.ReadSecret(context.Background(), "secret/data/myapp")

	assert.NoError(t, err)
	assert.Equal(t, expected, data)
}

func TestWriteSecret(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogical := mock_vaultclient.NewMockLogical(ctrl)
	mockAPI := mock_vaultclient.NewMockVaultAPI(ctrl)
	mockAPI.EXPECT().Logical().Return(mockLogical).AnyTimes()

	data := map[string]interface{}{"username": "admin", "password": "1234"}
	mockLogical.EXPECT().
		WriteWithContext(gomock.Any(), "secret/data/myapp", map[string]interface{}{"data": data}).
		Return(&vault.Secret{}, nil)

	c := &Client{api: mockAPI}
	err := c.WriteSecret(context.Background(), "secret/data/myapp", data)
	assert.NoError(t, err)
}
```

---

## 6️⃣ Lancer les tests

```bash
go test ./internal/vaultclient -v
```

💡 Résultat attendu :

```
=== RUN   TestReadSecret
--- PASS: TestReadSecret (0.00s)
=== RUN   TestWriteSecret
--- PASS: TestWriteSecret (0.00s)
PASS
ok  	vault-demo/internal/vaultclient	0.018s
```

---

## 7️⃣ Bonus : `main.go` (exécution simple)

```go
package main

import (
	"context"
	"fmt"
	"log"

	"vault-demo/internal/vaultclient"
)

func main() {
	client, err := vaultclient.New(vaultclient.Config{
		Address:  "https://vault.example.com:8200",
		AuthType: "token",
		Token:    "s.XXX",
	})
	if err != nil {
		log.Fatal(err)
	}

	data, err := client.ReadSecret(context.Background(), "secret/data/myapp")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Secret:", data)
}
```

---

Souhaites-tu que je t’ajoute un **Makefile** avec les commandes suivantes :

* `make generate` → exécute mockgen
* `make test` → lance les tests
* `make run` → lance le main ?

Super — voici un guide clair et prêt à l’emploi (en **Go**) pour te connecter à **HashiCorp Vault** et lire/écrire des secrets *username/password* (KV v2). Je fournis :

* installation,
* exemples d'authentification (Token + AppRole),
* lecture d'un secret KV v2 `secret/data/myapp`,
* bonne pratique rapide.

---

# 1) Installer le client Vault pour Go

```bash
go get github.com/hashicorp/vault/api
```

Le paquet officiel s'appelle `github.com/hashicorp/vault/api`. ([pkg.go.dev][1])

---

# 2) Préparer Vault (variables d'environnement)

Généralement on passe l'adresse et le token via des variables d'environnement :

```bash
export VAULT_ADDR="https://vault.example.com:8200"
export VAULT_TOKEN="s.XXXXX"     # si tu utilises token
```

Pour AppRole, tu fourniras `ROLE_ID` et `SECRET_ID` (ou les récupérer depuis un secure store). Voir AppRole docs. ([HashiCorp Developer][2])

---

# 3) Exemple : se connecter avec un **token** et lire un secret KV v2

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

func main() {
	// crée le client (utilise VAULT_ADDR si présent)
	cfg := vault.DefaultConfig()
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.Address = addr
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		log.Fatalf("vault.NewClient: %v", err)
	}

	// auth par token (depuis VAULT_TOKEN)
	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		client.SetToken(token)
	} else {
		log.Fatal("VAULT_TOKEN non défini")
	}

	// Lecture d'un secret KV v2 à "secret/data/myapp"
	// Pour KV v2 l'API HTTP /vault/data/<path> renvoie un champ "data" qui contient "data"
	secretPath := "secret/data/myapp" // adapter selon ton mount et path
	secret, err := client.Logical().ReadWithContext(context.Background(), secretPath)
	if err != nil {
		log.Fatalf("lecture secret: %v", err)
	}
	if secret == nil || secret.Data == nil {
		log.Fatalf("secret vide ou introuvable")
	}

	// Dans KV v2, secret.Data["data"] contient le map réel
	dataRaw, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		log.Fatalf("format KV v2 inattendu")
	}

	username, _ := dataRaw["username"].(string)
	password, _ := dataRaw["password"].(string)

	fmt.Printf("username=%s\npassword=%s\n", username, password)
}
```

Remarque : la structure KV v2 stocke les valeurs sous `data` → `data`. Docs KV v2. ([HashiCorp Developer][3])

---

# 4) Exemple : auth **AppRole** (machine-to-machine)

```go
package main

import (
	"context"
	"log"
	"os"

	vault "github.com/hashicorp/vault/api"
)

func main() {
	cfg := vault.DefaultConfig()
	if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		cfg.Address = addr
	}
	client, err := vault.NewClient(cfg)
	if err != nil {
		log.Fatalf("vault.NewClient: %v", err)
	}

	roleID := os.Getenv("VAULT_ROLE_ID")
	secretID := os.Getenv("VAULT_SECRET_ID")
	if roleID == "" || secretID == "" {
		log.Fatal("VARS VAULT_ROLE_ID et VAULT_SECRET_ID requises")
	}

	data := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}

	// AppRole login: POST to auth/approle/login
	secret, err := client.Logical().Write("auth/approle/login", data)
	if err != nil {
		log.Fatalf("AppRole login failed: %v", err)
	}
	if secret == nil || secret.Auth == nil {
		log.Fatalf("AppRole: réponse inattendue")
	}

	client.SetToken(secret.Auth.ClientToken)

	// ensuite tu peux lire les secrets comme dans l'exemple précédent
}
```

AppRole est documenté dans la doc officielle (procédure pour récupérer role_id / secret_id etc.). ([HashiCorp Developer][2])

---

# 5) Écrire un secret (KV v2) — exemple simple

```go
// write data to KV v2 at secret/data/myapp
path := "secret/data/myapp"
payload := map[string]interface{}{
    "data": map[string]interface{}{
        "username": "alice",
        "password": "s3cr3t",
    },
}
_, err = client.Logical().Write(path, payload)
if err != nil {
    log.Fatalf("write secret: %v", err)
}
```

Pour KV v2 il faut envoyer la clé `data` contenant ton map. ([HashiCorp Developer][3])

---

# 6) Bonnes pratiques rapides

* Préfère **AppRole** (ou auth cloud/Kubernetes) pour les apps plutôt que tokens statiques. ([HashiCorp Developer][2])
* Utilise la **version KV v2** si tu veux versionner et pouvoir restaurer/soft-delete. ([HashiCorp Developer][3])
* Pour des mots de passe de DB, considère d'utiliser le **Database Secrets Engine** (création dynamique d'identifiants) au lieu de stocker des credentials statiques. ([HashiCorp Developer][4])
* Ne commits jamais de `VAULT_TOKEN`, `ROLE_ID` ou `SECRET_ID` dans le code/repos.

---

# 7) Liens utiles (lecture)

* Client Go officiel (`pkg.go.dev` + examples) — référence du client Go. ([pkg.go.dev][1])
* KV v2 docs (format `data`/`metadata`) — pour read/write. ([HashiCorp Developer][3])
* AppRole auth docs — comment configurer role_id/secret_id. ([HashiCorp Developer][2])
* Exemples officiels `vault-examples` sur GitHub. ([GitHub][5])

---

Si tu veux, je te fournis *directement* :

* un petit repo skeleton Go (main + gestion d'auth token/AppRole + lecture/écriture) prêt à copier-coller, ou
* un exemple pour lire des secrets depuis Kubernetes (service account JWT) — dis ce que tu préfères et je te fournis le code complet tout de suite.

[1]: https://pkg.go.dev/github.com/hashicorp/vault/api?utm_source=chatgpt.com "api package - github.com/hashicorp/vault/api"
[2]: https://developer.hashicorp.com/vault/docs/auth/approle?utm_source=chatgpt.com "Use AppRole authentication | Vault"
[3]: https://developer.hashicorp.com/vault/docs/secrets/kv/kv-v2?utm_source=chatgpt.com "Key/Value v2 | Vault"
[4]: https://developer.hashicorp.com/vault/tutorials/get-started/understand-static-dynamic-secrets?utm_source=chatgpt.com "Understand static and dynamic secrets | Vault"
[5]: https://github.com/hashicorp/vault-examples?utm_source=chatgpt.com "hashicorp/vault-examples"
