Très bonne question 👍

En Go, le **refresh token HTTP** est une pratique courante pour renouveler un **access token** (souvent JWT) lorsqu’il expire, sans redemander à l’utilisateur de se reconnecter.
Voyons comment l’implémenter **proprement** — côté **serveur** et côté **client**.

---

## 🧩 1. Principe du refresh token

* Le **client** reçoit deux jetons après connexion :

  * `access_token` (court, ex: 15 min)
  * `refresh_token` (long, ex: 7 jours)
* Quand l’`access_token` expire, le client envoie le `refresh_token` à `/auth/refresh`
* Le **serveur** vérifie le refresh token, et renvoie un **nouvel access_token**

---

## 🏗️ 2. Exemple en Go (avec Gin)

### 📁 Arborescence simple

```
main.go
auth/
 ├── handler.go
 ├── service.go
 └── token.go
```

---

### 🔐 `auth/token.go`

Gestion des JWT.

```go
package auth

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var (
    accessSecret  = []byte("ACCESS_SECRET_KEY")
    refreshSecret = []byte("REFRESH_SECRET_KEY")
)

type Claims struct {
    UserID string `json:"user_id"`
    jwt.RegisteredClaims
}

func GenerateTokens(userID string) (accessToken, refreshToken string, err error) {
    accessClaims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
        },
    }

    refreshClaims := Claims{
        UserID: userID,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
        },
    }

    at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
    rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

    accessToken, err = at.SignedString(accessSecret)
    if err != nil {
        return "", "", err
    }
    refreshToken, err = rt.SignedString(refreshSecret)
    return
}

func ValidateToken(tokenStr string, isRefresh bool) (*Claims, error) {
    key := accessSecret
    if isRefresh {
        key = refreshSecret
    }

    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return key, nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, jwt.ErrTokenInvalidClaims
    }
    return claims, nil
}
```

---

### 🧠 `auth/service.go`

Service pour gérer le refresh.

```go
package auth

import "errors"

var ErrInvalidRefresh = errors.New("invalid refresh token")

type AuthService struct{}

func (s *AuthService) Refresh(refreshToken string) (string, error) {
    claims, err := ValidateToken(refreshToken, true)
    if err != nil {
        return "", ErrInvalidRefresh
    }

    newAccess, _, err := GenerateTokens(claims.UserID)
    if err != nil {
        return "", err
    }

    return newAccess, nil
}
```

---

### 🌐 `auth/handler.go`

Handler HTTP (avec Gin).

```go
package auth

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

type RefreshRequest struct {
    RefreshToken string `json:"refresh_token"`
}

func RefreshHandler(c *gin.Context) {
    var req RefreshRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
        return
    }

    service := AuthService{}
    newAccess, err := service.Refresh(req.RefreshToken)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"access_token": newAccess})
}
```

---

### 🚀 `main.go`

Pour lancer le serveur :

```go
package main

import (
    "github.com/gin-gonic/gin"
    "myapp/auth"
)

func main() {
    r := gin.Default()

    r.POST("/auth/refresh", auth.RefreshHandler)

    r.Run(":8080")
}
```

---

## 💡 Côté client (exemple Go)

Lorsqu’un appel retourne `401 Unauthorized`, on rafraîchit le token :

```go
func RefreshAccessToken(refreshToken string) (string, error) {
    reqBody := map[string]string{"refresh_token": refreshToken}
    body, _ := json.Marshal(reqBody)

    resp, err := http.Post("http://localhost:8080/auth/refresh", "application/json", bytes.NewBuffer(body))
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("refresh failed: %s", resp.Status)
    }

    var data map[string]string
    json.NewDecoder(resp.Body).Decode(&data)
    return data["access_token"], nil
}
```

---

## 🧱 Bonnes pratiques

✅ Stocker les refresh tokens en base (ou Redis) pour pouvoir les **révoquer**
✅ Utiliser des secrets distincts pour access et refresh
✅ Toujours **invalider les anciens refresh tokens** après un refresh
✅ Transmettre les tokens via **HTTPOnly cookies** pour éviter les attaques XSS

---

## unit Tests

Excellent réflexe 👏 — ajouter des **tests unitaires** à ton système d’auth avec `Gin` et les JWT est essentiel.

Je vais te montrer une version **complète et testable** du code précédent, avec des tests unitaires pour :

* la génération et validation des tokens (`token_test.go`),
* le service de refresh (`service_test.go`),
* et le handler HTTP (`handler_test.go`) à l’aide de `httptest`.

---

## 🧩 Structure du projet

```
auth/
 ├── token.go
 ├── token_test.go
 ├── service.go
 ├── service_test.go
 ├── handler.go
 ├── handler_test.go
main.go
go.mod
```

---

## 🔐 `auth/token.go` (identique à avant)

*(pas de changement ici — contient `GenerateTokens()` et `ValidateToken()`)*

---

## ✅ `auth/token_test.go`

Tests unitaires pour la logique JWT.

```go
package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateTokens(t *testing.T) {
	access, refresh, err := GenerateTokens("user123")
	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	// Validate access token
	claims, err := ValidateToken(access, false)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims.UserID)

	// Validate refresh token
	claims2, err := ValidateToken(refresh, true)
	assert.NoError(t, err)
	assert.Equal(t, "user123", claims2.UserID)
}

func TestExpiredToken(t *testing.T) {
	// Génération d'un token déjà expiré
	expiredClaims := Claims{
		UserID: "u1",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	tokenStr, _ := tok.SignedString(accessSecret)

	_, err := ValidateToken(tokenStr, false)
	assert.Error(t, err, "should fail on expired token")
}
```

---

## 🧠 `auth/service.go` (comme avant)

```go
package auth

import "errors"

var ErrInvalidRefresh = errors.New("invalid refresh token")

type AuthService struct{}

func (s *AuthService) Refresh(refreshToken string) (string, error) {
	claims, err := ValidateToken(refreshToken, true)
	if err != nil {
		return "", ErrInvalidRefresh
	}

	newAccess, _, err := GenerateTokens(claims.UserID)
	if err != nil {
		return "", err
	}

	return newAccess, nil
}
```

---

## 🧪 `auth/service_test.go`

```go
package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAuthService_Refresh(t *testing.T) {
	svc := &AuthService{}

	// token valide
	_, refresh, err := GenerateTokens("user1")
	assert.NoError(t, err)

	newAccess, err := svc.Refresh(refresh)
	assert.NoError(t, err)
	assert.NotEmpty(t, newAccess)

	// token invalide
	_, err = svc.Refresh("invalid.token.string")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidRefresh, err)
}
```

---

## 🌐 `auth/handler.go` (identique à avant)

```go
package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func RefreshHandler(c *gin.Context) {
	var req RefreshRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	service := AuthService{}
	newAccess, err := service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccess})
}
```

---

## 🧪 `auth/handler_test.go`

Tests du handler HTTP via `httptest`.

```go
package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestRefreshHandler_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/refresh", RefreshHandler)

	_, refresh, _ := GenerateTokens("user123")
	body, _ := json.Marshal(map[string]string{"refresh_token": refresh})

	req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NotEmpty(t, resp["access_token"])
}

func TestRefreshHandler_InvalidToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/refresh", RefreshHandler)

	body, _ := json.Marshal(map[string]string{"refresh_token": "invalid"})
	req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshHandler_BadRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/auth/refresh", RefreshHandler)

	req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer([]byte(`{bad json`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

---

## 🧰 Commandes de test

```bash
go test ./auth -v
```

ou pour la couverture :

```bash
go test ./auth -cover
```

Exemple de sortie :

```
=== RUN   TestGenerateAndValidateTokens
--- PASS: TestGenerateAndValidateTokens (0.00s)
=== RUN   TestAuthService_Refresh
--- PASS: TestAuthService_Refresh (0.00s)
=== RUN   TestRefreshHandler_Success
--- PASS: TestRefreshHandler_Success (0.01s)
PASS
ok  	myapp/auth	0.04s	coverage: 94.5% of statements
```

---

## 🔒 Bonus (si tu veux aller plus loin)

On peut ajouter :

* un **mock du service** via `mockgen` (`AuthService` mocké dans les tests HTTP),
* des tests d’expiration réels,
* un test d’intégration complet avec un vrai serveur.

---

Souhaites-tu que je t’ajoute **les mocks générés avec `mockgen`** pour tester le handler sans dépendre du vrai service (`AuthService`) ?

Excellent 👏 — ajoutons **des mocks générés avec `mockgen`** pour tester ton `AuthHandler` sans dépendre de la vraie logique du `AuthService`.

---

## 🧩 Objectif

On va :

1. **Extraire une interface** du `AuthService`
2. **Générer un mock** avec `mockgen`
3. **Utiliser ce mock dans les tests HTTP**
4. **Tester le handler isolément**, sans vraiment valider de JWT.

---

## ⚙️ Étape 1 : Créer l’interface

### 🔄 `auth/service.go`

On modifie légèrement le code pour exposer une interface `AuthRefresher` :

```go
package auth

import "errors"

var ErrInvalidRefresh = errors.New("invalid refresh token")

// Interface du service — pour faciliter le mock
type AuthRefresher interface {
	Refresh(refreshToken string) (string, error)
}

type AuthService struct{}

func (s *AuthService) Refresh(refreshToken string) (string, error) {
	claims, err := ValidateToken(refreshToken, true)
	if err != nil {
		return "", ErrInvalidRefresh
	}

	newAccess, _, err := GenerateTokens(claims.UserID)
	if err != nil {
		return "", err
	}

	return newAccess, nil
}
```

---

## ⚙️ Étape 2 : Modifier le handler pour accepter l’interface

### 🌐 `auth/handler.go`

On injecte le service plutôt que de le créer dans la fonction.

```go
package auth

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type AuthHandler struct {
	Service AuthRefresher
}

func (h *AuthHandler) RefreshHandler(c *gin.Context) {
	var req RefreshRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	newAccess, err := h.Service.Refresh(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": newAccess})
}
```

👉 Cela permet d’injecter un mock dans les tests.

---

## ⚙️ Étape 3 : Générer le mock avec `mockgen`

Ajoute cette directive en haut du fichier `service.go` :

```go
//go:generate mockgen -destination=mock_auth.go -package=auth . AuthRefresher
```

Ensuite, lance :

```bash
go install github.com/golang/mock/mockgen@latest
go generate ./auth
```

Cela va créer le fichier :

### 🧪 `auth/mock_auth.go`

*(généré automatiquement — ne pas modifier)*

Il contiendra quelque chose comme :

```go
// Code generated by MockGen. DO NOT EDIT.
package auth

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

type MockAuthRefresher struct {
	ctrl     *gomock.Controller
	recorder *MockAuthRefresherMockRecorder
}

type MockAuthRefresherMockRecorder struct {
	mock *MockAuthRefresher
}

func NewMockAuthRefresher(ctrl *gomock.Controller) *MockAuthRefresher {
	mock := &MockAuthRefresher{ctrl: ctrl}
	mock.recorder = &MockAuthRefresherMockRecorder{mock}
	return mock
}

func (m *MockAuthRefresher) EXPECT() *MockAuthRefresherMockRecorder {
	return m.recorder
}

func (m *MockAuthRefresher) Refresh(refreshToken string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh", refreshToken)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockAuthRefresherMockRecorder) Refresh(refreshToken interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockAuthRefresher)(nil).Refresh), refreshToken)
}
```

---

## ⚙️ Étape 4 : Écrire les tests avec le mock

### 🧪 `auth/handler_test.go`

```go
package auth

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRefreshHandler_WithMock_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockAuthRefresher(ctrl)
	handler := &AuthHandler{Service: mockSvc}

	router := gin.Default()
	router.POST("/auth/refresh", handler.RefreshHandler)

	body, _ := json.Marshal(map[string]string{"refresh_token": "valid_token"})
	req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	mockSvc.EXPECT().Refresh("valid_token").Return("new_access_token", nil)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var resp map[string]string
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, "new_access_token", resp["access_token"])
}

func TestRefreshHandler_WithMock_Unauthorized(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockAuthRefresher(ctrl)
	handler := &AuthHandler{Service: mockSvc}

	router := gin.Default()
	router.POST("/auth/refresh", handler.RefreshHandler)

	body, _ := json.Marshal(map[string]string{"refresh_token": "invalid"})
	req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	mockSvc.EXPECT().Refresh("invalid").Return("", ErrInvalidRefresh)

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestRefreshHandler_WithMock_BadRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := NewMockAuthRefresher(ctrl)
	handler := &AuthHandler{Service: mockSvc}

	router := gin.Default()
	router.POST("/auth/refresh", handler.RefreshHandler)

	req, _ := http.NewRequest("POST", "/auth/refresh", bytes.NewBuffer([]byte(`{invalid_json`)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
```

---

## 🧰 Étape 5 : Lancer les tests

```bash
go test ./auth -v
```

✅ Tu verras que les handlers sont maintenant **testés sans dépendre du vrai JWT service**, grâce aux **mocks générés par `mockgen`**.

---

## 🚀 Résumé

* `AuthRefresher` est une interface
* `AuthHandler` dépend de cette interface
* `mockgen` crée `MockAuthRefresher`
* Les tests injectent le mock et vérifient les réponses HTTP

---


```go
type TokenResponse struct {
    AccessToken string `json:"access_token"`
    ExpiredIn string `json:"expired_in"`
    RefreshExpiresIn string `json:"refresh_expires_in"`
    TokenType string `json:"token_type"`
    NotBeforePolicy string `json:"not_before_policy"`
    Scope string `json:"scope"`
}
```
