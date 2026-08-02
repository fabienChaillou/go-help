Parfait ! On va créer un exemple **complet avec `go-oidc` et OAuth2** qui :

1. Récupère un **Access Token** et **ID Token** via **Resource Owner Password Credentials (ROPC)**.
2. Vérifie l’ID Token.
3. Utilise le **Refresh Token** pour obtenir un nouveau Access Token quand l’ancien expire.

---

## Exemple Go avec Refresh Token

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	issuer := "https://keycloak.example.com/auth/realms/myrealm"
	clientID := "my-client"
	clientSecret := "my-secret" // si client confidentiel
	username := "alice"
	password := "s3cr3t"

	// 1) Récupération du provider OIDC
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	// 2) Config OAuth2
	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// 3) Récupérer un token initial (ROPC)
	token, err := conf.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	fmt.Println("=== Initial Tokens ===")
	printToken(token)

	// 4) Vérifier l'ID Token
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	idToken, err := verifier.Verify(ctx, token.Extra("id_token").(string))
	if err != nil {
		log.Fatalf("Failed to verify ID token: %v", err)
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		log.Fatalf("Failed to parse claims: %v", err)
	}
	fmt.Printf("ID Token claims: %+v\n", claims)

	// 5) Simuler expiration et rafraîchissement
	fmt.Println("\n=== Simulating token refresh in 3 seconds ===")
	time.Sleep(3 * time.Second)

	ts := conf.TokenSource(ctx, token)
	newToken, err := ts.Token() // rafraîchit automatiquement si nécessaire
	if err != nil {
		log.Fatalf("Failed to refresh token: %v", err)
	}

	fmt.Println("=== Refreshed Tokens ===")
	printToken(newToken)
}

// Helper pour afficher le token
func printToken(t *oauth2.Token) {
	fmt.Printf("Access Token: %s\n", t.AccessToken)
	fmt.Printf("Refresh Token: %s\n", t.RefreshToken)
	fmt.Printf("Expiry: %v\n", t.Expiry)
}
```

---

### Comment ça marche :

1. `PasswordCredentialsToken` → récupère **access token**, **refresh token** et **id token**.
2. `go-oidc` vérifie la signature de l’ID Token et les claims.
3. `oauth2.Config.TokenSource` gère automatiquement le **rafraîchissement** via le Refresh Token.
4. `ts.Token()` va renvoyer le même token si encore valide ou rafraîchir si expiré.

---

💡 Points à retenir :

* Tu peux utiliser le même pattern pour **service back-end ou scripts automatisés**.
* Ne jamais exposer `clientSecret` dans un front-end public.
* Pour un usage web SPA, préfère **Authorization Code + PKCE** plutôt que ROPC.

---

## other example

Parfait ! On peut faire un exemple simple avec **`go-oidc` + `golang.org/x/oauth2`** pour récupérer un **access token** via le **Resource Owner Password Credentials (ROPC)** ou via **Authorization Code Flow**.

⚠️ Important : avec `go-oidc`, on ne “récupère pas de mots de passe Keycloak en clair” — on récupère un **token JWT**.

Je te montre un exemple **Récupération de token avec ROPC** (simple pour script / service back-end).

---

## Exemple Go : récupérer un token Keycloak avec `go-oidc`

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	issuer := "https://keycloak.example.com/auth/realms/myrealm"
	clientID := "my-client"
	clientSecret := "my-secret" // si client confidentiel
	username := "alice"
	password := "s3cr3t"

	// Récupérer le provider OIDC
	provider, err := oidc.NewProvider(ctx, issuer)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	// Config OAuth2
	conf := oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// ROPC (Resource Owner Password Credentials)
	token, err := conf.PasswordCredentialsToken(ctx, username, password)
	if err != nil {
		log.Fatalf("Failed to get token: %v", err)
	}

	fmt.Printf("Access token: %s\n", token.AccessToken)
	fmt.Printf("Refresh token: %s\n", token.RefreshToken)
	fmt.Printf("Expiry: %v\n", token.Expiry)

	// Vérification du token ID (optionnel)
	verifier := provider.Verifier(&oidc.Config{ClientID: clientID})
	idToken, err := verifier.Verify(ctx, token.Extra("id_token").(string))
	if err != nil {
		log.Fatalf("Failed to verify ID token: %v", err)
	}

	var claims map[string]interface{}
	if err := idToken.Claims(&claims); err != nil {
		log.Fatalf("Failed to parse claims: %v", err)
	}

	fmt.Printf("ID Token claims: %+v\n", claims)
}
```

### Explications :

* `oidc.NewProvider(ctx, issuer)` : récupère la configuration OIDC depuis Keycloak (`.well-known/openid-configuration`).
* `oauth2.Config` : gère le flow OAuth2.
* `PasswordCredentialsToken` : récupère un **access token** et un **ID token** via login/password.
* `verifier.Verify(...)` : vérifie que l’ID token est bien signé par Keycloak et valide pour ton client.
* Tu peux maintenant utiliser `token.AccessToken` pour appeler des API protégées.

---

💡 Remarques importantes :

1. ROPC est **pratique pour les scripts ou services back-end**, mais **à éviter pour une appli web ou mobile publique**.
2. Pour des applications web ou SPA, il faut préférer **Authorization Code + PKCE**.
3. Tu peux ensuite rafraîchir le token via `oauth2.TokenSource`.

---

