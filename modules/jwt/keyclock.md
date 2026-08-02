Très bien 👍
On va faire un **petit module Go complet** pour interagir avec Keycloak :

* récupération d’un **access_token** et d’un **refresh_token**,
* **rafraîchissement** du token,
* avec des **tests unitaires (mockés)**, sans requête réseau réelle.

---

## 🗂️ Structure du module

```
keycloak/
├── keycloak.go          # code du client
├── keycloak_test.go     # tests unitaires
├── go.mod
```

---

## 📦 `keycloak/keycloak.go`

```go
package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// TokenResponse représente la réponse de Keycloak
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// Client permet d’interagir avec Keycloak
type Client struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client
}

// NewClient crée un client Keycloak
func NewClient(baseURL, realm, clientID, clientSecret string) *Client {
	return &Client{
		BaseURL:      baseURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) tokenEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.BaseURL, c.Realm)
}

// GetToken récupère un access_token et un refresh_token
func (c *Client) GetToken(ctx context.Context) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak: status %d", resp.StatusCode)
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}

	if tr.AccessToken == "" {
		return nil, errors.New("missing access_token")
	}
	return &tr, nil
}

// RefreshToken permet de rafraîchir un token expiré
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint(), bytes.NewBufferString(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak: status %d", resp.StatusCode)
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}

	if tr.AccessToken == "" {
		return nil, errors.New("missing access_token")
	}
	return &tr, nil
}
```

---

## 🧪 `keycloak/keycloak_test.go`

Ces tests utilisent un **serveur HTTP mocké** (`httptest.Server`) pour simuler les réponses Keycloak.

```go
package keycloak

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// crée une fausse réponse JSON
func mockTokenResponse() TokenResponse {
	return TokenResponse{
		AccessToken:  "access123",
		RefreshToken: "refresh456",
		ExpiresIn:    3600,
		TokenType:    "bearer",
	}
}

func TestGetToken_Success(t *testing.T) {
	mockResp := mockTokenResponse()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/realms/test/protocol/openid-connect/token" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	token, err := client.GetToken(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if token.AccessToken != mockResp.AccessToken {
		t.Errorf("expected %s, got %s", mockResp.AccessToken, token.AccessToken)
	}
	if token.RefreshToken != mockResp.RefreshToken {
		t.Errorf("expected %s, got %s", mockResp.RefreshToken, token.RefreshToken)
	}
}

func TestRefreshToken_Success(t *testing.T) {
	mockResp := mockTokenResponse()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			t.Fatalf("ParseForm: %v", err)
		}
		if r.FormValue("grant_type") != "refresh_token" {
			t.Errorf("expected grant_type=refresh_token, got %s", r.FormValue("grant_type"))
		}
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	token, err := client.RefreshToken(context.Background(), "refresh456")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token.AccessToken != mockResp.AccessToken {
		t.Errorf("expected %s, got %s", mockResp.AccessToken, token.AccessToken)
	}
}

func TestGetToken_Failure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	_, err := client.GetToken(context.Background())
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
```

---

## 🚀 Lancer les tests

```bash
go test ./keycloak -v
```

Sortie attendue :

```
=== RUN   TestGetToken_Success
--- PASS: TestGetToken_Success (0.00s)
=== RUN   TestRefreshToken_Success
--- PASS: TestRefreshToken_Success (0.00s)
=== RUN   TestGetToken_Failure
--- PASS: TestGetToken_Failure (0.00s)
PASS
ok      keycloak    0.005s
```

---

Souhaites-tu que je t’ajoute :

* 🔐 la gestion automatique du **renouvellement** de token (si expiré), ou

Excellent 👍
On va enrichir le module précédent avec une **gestion automatique du renouvellement du token** dès qu’il est expiré — c’est-à-dire :

✅ stockage du token courant et de son expiration,
✅ vérification avant chaque utilisation,
✅ appel automatique à `RefreshToken()` si nécessaire,
✅ sinon appel à `GetToken()` si aucun token ou si le refresh échoue.

---

## 🧩 Nouvelle version de `keycloak/keycloak.go`

```go
package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	ExpiresAt    time.Time
}

type Client struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client

	mu    sync.Mutex
	token *TokenResponse
}

// NewClient crée un client Keycloak
func NewClient(baseURL, realm, clientID, clientSecret string) *Client {
	return &Client{
		BaseURL:      baseURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) tokenEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.BaseURL, c.Realm)
}

// internal helper pour effectuer la requête vers Keycloak
func (c *Client) doTokenRequest(ctx context.Context, form url.Values) (*TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint(), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak: status %d", resp.StatusCode)
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	if tr.AccessToken == "" {
		return nil, errors.New("missing access_token")
	}
	tr.ExpiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return &tr, nil
}

// GetToken récupère un access_token et un refresh_token
func (c *Client) GetToken(ctx context.Context) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)

	return c.doTokenRequest(ctx, data)
}

// RefreshToken permet de rafraîchir un token expiré
func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)

	return c.doTokenRequest(ctx, data)
}

// GetValidToken renvoie toujours un token valide, en le rafraîchissant si nécessaire.
func (c *Client) GetValidToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Aucun token -> on en récupère un
	if c.token == nil {
		tr, err := c.GetToken(ctx)
		if err != nil {
			return "", err
		}
		c.token = tr
		return tr.AccessToken, nil
	}

	// Token encore valide
	if time.Now().Before(c.token.ExpiresAt.Add(-10 * time.Second)) { // marge 10s
		return c.token.AccessToken, nil
	}

	// Token expiré → tentative de refresh
	tr, err := c.RefreshToken(ctx, c.token.RefreshToken)
	if err != nil {
		// fallback : re-auth complète
		tr, err = c.GetToken(ctx)
		if err != nil {
			return "", err
		}
	}
	c.token = tr
	return tr.AccessToken, nil
}
```

---

## 🧪 Tests : `keycloak/keycloak_test.go`

On ajoute un test de **renouvellement automatique**.

```go
package keycloak

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func mockTokenResponse() TokenResponse {
	return TokenResponse{
		AccessToken:  "access123",
		RefreshToken: "refresh456",
		ExpiresIn:    1, // court pour simuler expiration
		TokenType:    "bearer",
	}
}

func TestGetValidToken_Renewal(t *testing.T) {
	var callCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		var resp TokenResponse
		if r.FormValue("grant_type") == "refresh_token" {
			resp = TokenResponse{
				AccessToken:  "new_access",
				RefreshToken: "new_refresh",
				ExpiresIn:    3600,
			}
		} else {
			resp = mockTokenResponse()
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	ctx := context.Background()
	token1, err := client.GetValidToken(ctx)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if token1 != "access123" {
		t.Fatalf("expected access123, got %s", token1)
	}

	// Attendre expiration
	time.Sleep(2 * time.Second)

	token2, err := client.GetValidToken(ctx)
	if err != nil {
		t.Fatalf("err: %v", err)
	}
	if token2 != "new_access" {
		t.Fatalf("expected new_access after refresh, got %s", token2)
	}

	if callCount < 2 {
		t.Fatalf("expected at least 2 calls (get+refresh), got %d", callCount)
	}
}
```

---

## ✅ Résumé du comportement

| Situation           | Action                       |
| ------------------- | ---------------------------- |
| Aucun token stocké  | `GetToken()` appelé          |
| Token encore valide | Renvoi direct                |
| Token expiré        | Tentative `RefreshToken()`   |
| Refresh échoue      | Reconnexion via `GetToken()` |

---

## 🚀 Utilisation

```go
ctx := context.Background()
kc := keycloak.NewClient("https://keycloak.example.com", "myrealm", "my-client", "secret")

token, err := kc.GetValidToken(ctx)
if err != nil {
    log.Fatalf("Keycloak error: %v", err)
}

// token toujours valide ici
req, _ := http.NewRequest("GET", "https://api.example.com", nil)
req.Header.Set("Authorization", "Bearer "+token)
```

---

Souhaites-tu que je t’ajoute maintenant :

* 🔁 un **auto-refresh goroutine** qui renouvelle le token périodiquement en arrière-plan ?

Parfait 👌
On va compléter le module Keycloak avec un **auto-refresh en arrière-plan** — un **goroutine** qui :

* vérifie périodiquement la validité du token,
* appelle `RefreshToken()` si nécessaire,
* se relance automatiquement,
* et s’arrête proprement avec un `context.CancelFunc`.

---

## 🧩 Nouvelle version : `keycloak/keycloak.go`

Voici la version complète avec l’auto-refresh intégré :

```go
package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// TokenResponse représente la réponse JSON de Keycloak
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	ExpiresAt    time.Time
}

// Client encapsule la configuration Keycloak et le token courant
type Client struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client

	mu        sync.Mutex
	token     *TokenResponse
	stopCh    chan struct{}
	refreshWG sync.WaitGroup
}

// NewClient crée un client Keycloak
func NewClient(baseURL, realm, clientID, clientSecret string) *Client {
	return &Client{
		BaseURL:      baseURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) tokenEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.BaseURL, c.Realm)
}

func (c *Client) doTokenRequest(ctx context.Context, form url.Values) (*TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint(), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak: status %d", resp.StatusCode)
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	if tr.AccessToken == "" {
		return nil, errors.New("missing access_token")
	}
	tr.ExpiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return &tr, nil
}

func (c *Client) GetToken(ctx context.Context) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	return c.doTokenRequest(ctx, data)
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	return c.doTokenRequest(ctx, data)
}

// GetValidToken garantit qu'on obtient toujours un token valide
func (c *Client) GetValidToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token == nil {
		tr, err := c.GetToken(ctx)
		if err != nil {
			return "", err
		}
		c.token = tr
		return tr.AccessToken, nil
	}

	if time.Now().Before(c.token.ExpiresAt.Add(-10 * time.Second)) {
		return c.token.AccessToken, nil
	}

	tr, err := c.RefreshToken(ctx, c.token.RefreshToken)
	if err != nil {
		tr, err = c.GetToken(ctx)
		if err != nil {
			return "", err
		}
	}
	c.token = tr
	return tr.AccessToken, nil
}

// StartAutoRefresh démarre une goroutine qui renouvelle le token automatiquement
func (c *Client) StartAutoRefresh(ctx context.Context, interval time.Duration) {
	c.stopCh = make(chan struct{})
	c.refreshWG.Add(1)
	go func() {
		defer c.refreshWG.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				return
			case <-ticker.C:
				c.refreshIfNeeded(ctx)
			}
		}
	}()
}

// StopAutoRefresh arrête la goroutine proprement
func (c *Client) StopAutoRefresh() {
	if c.stopCh != nil {
		close(c.stopCh)
		c.refreshWG.Wait()
	}
}

// refreshIfNeeded est interne : tente un refresh si le token est expiré
func (c *Client) refreshIfNeeded(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token == nil {
		tr, err := c.GetToken(ctx)
		if err == nil {
			c.token = tr
		}
		return
	}

	if time.Now().After(c.token.ExpiresAt.Add(-15 * time.Second)) {
		tr, err := c.RefreshToken(ctx, c.token.RefreshToken)
		if err == nil {
			c.token = tr
		}
	}
}
```

---

## 🧪 Exemple d’utilisation

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"yourmodule/keycloak"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kc := keycloak.NewClient(
		"https://keycloak.example.com",
		"myrealm",
		"my-client",
		"secret",
	)

	// Récupère le premier token valide
	token, err := kc.GetValidToken(ctx)
	if err != nil {
		log.Fatalf("failed to get token: %v", err)
	}
	fmt.Println("Initial token:", token)

	// Démarre le rafraîchissement auto toutes les 30 secondes
	kc.StartAutoRefresh(ctx, 30*time.Second)
	defer kc.StopAutoRefresh()

	// Simulation de service qui tourne
	for i := 0; i < 5; i++ {
		time.Sleep(10 * time.Second)
		tok, _ := kc.GetValidToken(ctx)
		fmt.Printf("[%d] Current token: %.10s...\n", i, tok)
	}
}
```

---

## ⚙️ Fonctionnement

| Élément              | Rôle                                             |
| -------------------- | ------------------------------------------------ |
| `StartAutoRefresh()` | Lance une goroutine qui tourne en fond           |
| `interval`           | fréquence à laquelle on vérifie la validité      |
| `refreshIfNeeded()`  | vérifie si le token expire bientôt (marge 15s)   |
| `StopAutoRefresh()`  | arrête proprement la goroutine                   |
| `sync.Mutex`         | empêche les race conditions sur le token partagé |

---

## 🧠 Astuce production

* `interval` peut être mis à `time.Minute` ou `time.Duration(token.ExpiresIn/2)` pour adapter la fréquence.
* Tu peux ajouter un callback (`OnTokenRefresh func(token string)`) pour notifier un autre composant (par ex. un client HTTP Gin).

---

Souhaites-tu que je t’ajoute une **version avec callback (hook)**, qui est appelée à chaque rafraîchissement réussi du token (utile si tu veux le propager dans des requêtes sortantes, des middlewares, etc.) ?

Parfait 💪
On va donc enrichir notre client Keycloak avec un **callback automatique** (`OnTokenRefresh`) — appelé **chaque fois qu’un token est renouvelé avec succès** (que ce soit un refresh ou un nouveau token).

Ce hook est très pratique pour :

* Mettre à jour un header Authorization dans un client HTTP,
* Stocker le nouveau token dans un cache,
* Notifier un autre service (par ex. via channel ou log).

---

## 🧩 Version finale du module : `keycloak/keycloak.go`

```go
package keycloak

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

// TokenResponse représente la réponse JSON de Keycloak
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	ExpiresAt    time.Time
}

// Client encapsule la configuration Keycloak et le token courant
type Client struct {
	BaseURL      string
	Realm        string
	ClientID     string
	ClientSecret string
	HTTPClient   *http.Client

	// Hook appelé à chaque rafraîchissement réussi
	OnTokenRefresh func(newToken *TokenResponse)

	mu        sync.Mutex
	token     *TokenResponse
	stopCh    chan struct{}
	refreshWG sync.WaitGroup
}

// NewClient crée un client Keycloak
func NewClient(baseURL, realm, clientID, clientSecret string) *Client {
	return &Client{
		BaseURL:      baseURL,
		Realm:        realm,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *Client) tokenEndpoint() string {
	return fmt.Sprintf("%s/realms/%s/protocol/openid-connect/token", c.BaseURL, c.Realm)
}

func (c *Client) doTokenRequest(ctx context.Context, form url.Values) (*TokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.tokenEndpoint(), bytes.NewBufferString(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("keycloak: status %d", resp.StatusCode)
	}

	var tr TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return nil, err
	}
	if tr.AccessToken == "" {
		return nil, errors.New("missing access_token")
	}
	tr.ExpiresAt = time.Now().Add(time.Duration(tr.ExpiresIn) * time.Second)
	return &tr, nil
}

func (c *Client) GetToken(ctx context.Context) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	tr, err := c.doTokenRequest(ctx, data)
	if err == nil && c.OnTokenRefresh != nil {
		c.OnTokenRefresh(tr)
	}
	return tr, err
}

func (c *Client) RefreshToken(ctx context.Context, refreshToken string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.ClientID)
	data.Set("client_secret", c.ClientSecret)
	tr, err := c.doTokenRequest(ctx, data)
	if err == nil && c.OnTokenRefresh != nil {
		c.OnTokenRefresh(tr)
	}
	return tr, err
}

// GetValidToken garantit qu'on obtient toujours un token valide
func (c *Client) GetValidToken(ctx context.Context) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token == nil {
		tr, err := c.GetToken(ctx)
		if err != nil {
			return "", err
		}
		c.token = tr
		return tr.AccessToken, nil
	}

	if time.Now().Before(c.token.ExpiresAt.Add(-10 * time.Second)) {
		return c.token.AccessToken, nil
	}

	tr, err := c.RefreshToken(ctx, c.token.RefreshToken)
	if err != nil {
		tr, err = c.GetToken(ctx)
		if err != nil {
			return "", err
		}
	}
	c.token = tr
	return tr.AccessToken, nil
}

// StartAutoRefresh démarre une goroutine qui renouvelle le token automatiquement
func (c *Client) StartAutoRefresh(ctx context.Context, interval time.Duration) {
	c.stopCh = make(chan struct{})
	c.refreshWG.Add(1)
	go func() {
		defer c.refreshWG.Done()
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-c.stopCh:
				return
			case <-ticker.C:
				c.refreshIfNeeded(ctx)
			}
		}
	}()
}

// StopAutoRefresh arrête la goroutine proprement
func (c *Client) StopAutoRefresh() {
	if c.stopCh != nil {
		close(c.stopCh)
		c.refreshWG.Wait()
	}
}

func (c *Client) refreshIfNeeded(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.token == nil {
		tr, err := c.GetToken(ctx)
		if err == nil {
			c.token = tr
		}
		return
	}

	if time.Now().After(c.token.ExpiresAt.Add(-15 * time.Second)) {
		tr, err := c.RefreshToken(ctx, c.token.RefreshToken)
		if err == nil {
			c.token = tr
		}
	}
}
```

---

## 🧠 Exemple d’utilisation avec callback

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"yourmodule/keycloak"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	kc := keycloak.NewClient(
		"https://keycloak.example.com",
		"myrealm",
		"my-client",
		"secret",
	)

	// 🔔 Callback : appelé à chaque fois que le token est rafraîchi
	kc.OnTokenRefresh = func(newToken *keycloak.TokenResponse) {
		fmt.Println("✅ Nouveau token reçu ! Expires at:", newToken.ExpiresAt)
		// Exemple : mettre à jour un client HTTP global
		// myHttpClient.SetAuthToken(newToken.AccessToken)
	}

	token, err := kc.GetValidToken(ctx)
	if err != nil {
		log.Fatalf("Erreur initiale : %v", err)
	}
	fmt.Println("Premier token :", token[:10], "...")

	// Lancer le rafraîchissement automatique toutes les 30s
	kc.StartAutoRefresh(ctx, 30*time.Second)
	defer kc.StopAutoRefresh()

	// Simulation : ton service tourne pendant 2 minutes
	for i := 0; i < 12; i++ {
		time.Sleep(10 * time.Second)
		tok, _ := kc.GetValidToken(ctx)
		fmt.Printf("[%d] Token actif: %.10s...\n", i, tok)
	}
}
```

---

## ⚙️ Fonctionnement du callback

| Moment                                                      | Ce qui se passe                              |
| ----------------------------------------------------------- | -------------------------------------------- |
| `GetToken()` ou `RefreshToken()` réussit                    | Appel de `OnTokenRefresh`                    |
| Token inchangé (encore valide)                              | Pas d’appel                                  |
| Token récupéré automatiquement (via auto-refresh goroutine) | Appel du callback                            |
| Erreur Keycloak                                             | Pas d’appel, le dernier token reste inchangé |

---

## 🔐 Avantages

* 🔄 Rafraîchissement automatique et transparent
* 🔔 Notification immédiate sur chaque mise à jour du token
* 🧵 Sécurisé via mutex et arrêt propre de goroutine
* 🧰 Facile à intégrer dans un service long-running (API, worker, etc.)

---

## unit tests

Excellent 💪
On va compléter ton module Keycloak avec **des tests unitaires complets**, y compris :

✅ tests de `GetToken` et `RefreshToken`,
✅ test du **renouvellement automatique** (`GetValidToken` + expiration simulée),
✅ test du **callback `OnTokenRefresh`**,
✅ test du **goroutine d’auto-refresh** (`StartAutoRefresh` + `StopAutoRefresh`).

---

## 🧪 `keycloak/keycloak_test.go`

```go
package keycloak

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// crée une fausse réponse JSON de Keycloak
func mockTokenResponse(access, refresh string, expiresIn int) TokenResponse {
	return TokenResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		ExpiresIn:    expiresIn,
		TokenType:    "bearer",
	}
}

func TestGetToken_Success(t *testing.T) {
	mockResp := mockTokenResponse("access123", "refresh456", 3600)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	token, err := client.GetToken(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token.AccessToken != "access123" {
		t.Errorf("expected access123, got %s", token.AccessToken)
	}
	if token.RefreshToken != "refresh456" {
		t.Errorf("expected refresh456, got %s", token.RefreshToken)
	}
}

func TestRefreshToken_Success(t *testing.T) {
	mockResp := mockTokenResponse("newAccess", "newRefresh", 3600)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.FormValue("grant_type") != "refresh_token" {
			t.Fatalf("expected grant_type=refresh_token, got %s", r.FormValue("grant_type"))
		}
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	token, err := client.RefreshToken(context.Background(), "oldRefresh")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if token.AccessToken != "newAccess" {
		t.Errorf("expected newAccess, got %s", token.AccessToken)
	}
}

func TestGetValidToken_RefreshFlow(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		var resp TokenResponse
		switch atomic.LoadInt32(&callCount) {
		case 1:
			resp = mockTokenResponse("firstAccess", "firstRefresh", 1) // expire vite
		default:
			resp = mockTokenResponse("secondAccess", "secondRefresh", 3600)
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	ctx := context.Background()
	token1, err := client.GetValidToken(ctx)
	if err != nil {
		t.Fatalf("GetValidToken error: %v", err)
	}
	if token1 != "firstAccess" {
		t.Fatalf("expected firstAccess, got %s", token1)
	}

	time.Sleep(2 * time.Second) // attendre expiration

	token2, err := client.GetValidToken(ctx)
	if err != nil {
		t.Fatalf("GetValidToken refresh error: %v", err)
	}
	if token2 != "secondAccess" {
		t.Fatalf("expected secondAccess after refresh, got %s", token2)
	}
}

func TestOnTokenRefresh_CallbackCalled(t *testing.T) {
	mockResp := mockTokenResponse("cbAccess", "cbRefresh", 1)
	called := int32(0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	client.OnTokenRefresh = func(tr *TokenResponse) {
		if tr.AccessToken != "cbAccess" {
			t.Errorf("expected cbAccess, got %s", tr.AccessToken)
		}
		atomic.StoreInt32(&called, 1)
	}

	_, err := client.GetToken(context.Background())
	if err != nil {
		t.Fatalf("GetToken error: %v", err)
	}

	if atomic.LoadInt32(&called) == 0 {
		t.Fatal("expected OnTokenRefresh to be called")
	}
}

func TestAutoRefresh_Goroutine(t *testing.T) {
	var refreshCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&refreshCount, 1)
		resp := mockTokenResponse("autoToken", "autoRefresh", 1)
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	client.OnTokenRefresh = func(tr *TokenResponse) {
		t.Logf("token refreshed: %s", tr.AccessToken)
	}

	client.StartAutoRefresh(ctx, 500*time.Millisecond)
	defer client.StopAutoRefresh()

	// attendre que la goroutine fasse au moins 2 refreshs
	time.Sleep(1600 * time.Millisecond)

	if atomic.LoadInt32(&refreshCount) < 2 {
		t.Fatalf("expected at least 2 refresh calls, got %d", refreshCount)
	}
}
```

---

## 🧩 Points testés

| Test                                | Description                                                  |
| ----------------------------------- | ------------------------------------------------------------ |
| `TestGetToken_Success`              | Vérifie récupération de token simple                         |
| `TestRefreshToken_Success`          | Vérifie le flow `grant_type=refresh_token`                   |
| `TestGetValidToken_RefreshFlow`     | Simule expiration et rafraîchissement automatique            |
| `TestOnTokenRefresh_CallbackCalled` | Vérifie que le hook `OnTokenRefresh` est bien appelé         |
| `TestAutoRefresh_Goroutine`         | Lance le goroutine auto-refresh et vérifie qu’il tourne bien |

---

## 🧪 Exécution

```bash
go test ./keycloak -v
```

Sortie attendue :

```
=== RUN   TestGetToken_Success
--- PASS: TestGetToken_Success (0.00s)
=== RUN   TestRefreshToken_Success
--- PASS: TestRefreshToken_Success (0.00s)
=== RUN   TestGetValidToken_RefreshFlow
--- PASS: TestGetValidToken_RefreshFlow (2.01s)
=== RUN   TestOnTokenRefresh_CallbackCalled
--- PASS: TestOnTokenRefresh_CallbackCalled (0.00s)
=== RUN   TestAutoRefresh_Goroutine
    keycloak_test.go:153: token refreshed: autoToken
    keycloak_test.go:153: token refreshed: autoToken
--- PASS: TestAutoRefresh_Goroutine (1.60s)
PASS
ok      keycloak    3.6s
```

---

Souhaites-tu que je t’ajoute aussi :

* ✅ des **tests de concurrence** (plusieurs appels simultanés à `GetValidToken` pour vérifier les verrous),

Parfait ⚙️ — on va donc ajouter des **tests de concurrence** à ton module Keycloak.
Ces tests vont vérifier que :

✅ les accès simultanés à `GetValidToken()` ne provoquent **aucune condition de course**,
✅ le **mutex (`mu`)** protège correctement le token,
✅ et que **un seul rafraîchissement** est réellement exécuté, même si plusieurs goroutines appellent `GetValidToken()` en même temps.

---

## 🧪 Nouveau test : `TestGetValidToken_ConcurrentAccess`

Ajoute ceci à la fin de ton fichier `keycloak_test.go` :

```go
func TestGetValidToken_ConcurrentAccess(t *testing.T) {
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)

		// Le premier appel renvoie un token très court
		if atomic.LoadInt32(&callCount) == 1 {
			_ = json.NewEncoder(w).Encode(mockTokenResponse("firstAccess", "firstRefresh", 1))
		} else {
			// Les suivants simulent un rafraîchissement
			_ = json.NewEncoder(w).Encode(mockTokenResponse("refreshedAccess", "refreshedRefresh", 3600))
		}
	}))
	defer server.Close()

	client := NewClient(server.URL, "test", "client-id", "secret")
	client.HTTPClient = server.Client()

	ctx := context.Background()
	_, err := client.GetValidToken(ctx)
	if err != nil {
		t.Fatalf("initial GetValidToken failed: %v", err)
	}

	// On attend que le token expire
	time.Sleep(1200 * time.Millisecond)

	const goroutines = 20
	results := make(chan string, goroutines)

	// Lancement de plusieurs goroutines concurrentes
	for i := 0; i < goroutines; i++ {
		go func() {
			token, err := client.GetValidToken(ctx)
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}
			results <- token
		}()
	}

	// Collecter les résultats
	tokens := make(map[string]int)
	for i := 0; i < goroutines; i++ {
		tok := <-results
		tokens[tok]++
	}

	if len(tokens) != 1 {
		t.Errorf("expected only one final token shared, got multiple: %+v", tokens)
	}

	if atomic.LoadInt32(&callCount) > 3 {
		t.Errorf("expected at most 2-3 calls to token endpoint, got %d", callCount)
	}
}
```

---

## 🧩 Explication

| Étape | Description                                                                                 |
| ----- | ------------------------------------------------------------------------------------------- |
| 1️⃣   | On simule un premier token à durée de vie très courte (`expires_in = 1s`).                  |
| 2️⃣   | Après expiration, on lance 20 goroutines qui appellent **simultanément** `GetValidToken()`. |
| 3️⃣   | Grâce au **mutex interne (`c.mu`)**, un seul thread rafraîchit réellement le token.         |
| 4️⃣   | Tous les threads obtiennent le **même token final** (même valeur d’`AccessToken`).          |
| 5️⃣   | On vérifie qu’il n’y a **pas plus de 2–3 appels HTTP** au total vers le serveur mocké.      |

---

## ✅ Exemple de sortie

```bash
=== RUN   TestGetValidToken_ConcurrentAccess
--- PASS: TestGetValidToken_ConcurrentAccess (1.25s)
PASS
ok      keycloak    1.3s
```

---

## 🧠 Que garantit ce test ?

* Aucun **data race** (`go test -race ./keycloak` doit être propre).
* Un seul **rafraîchissement simultané** a lieu, même sous forte charge.
* Le cache interne (`client.token`) est partagé en lecture en toute sécurité.

---

Souhaites-tu que je t’ajoute un **test de concurrence avancé** où :

* le refresh échoue une fois (simulateur de `401` ou `500`),
* plusieurs goroutines déclenchent un **fallback automatique** vers `GetToken()` ?

C’est très utile pour tester la **tolérance aux erreurs réseau / Keycloak down**.
