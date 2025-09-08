En Go, le **package `context`** est utilisé pour **transporter des délais d'expiration, des annulations et d'autres valeurs spécifiques à la requête** à travers les appels de fonctions et les goroutines. Il est essentiel pour gérer les **opérations asynchrones**, les **requêtes HTTP**, et les **bases de données** de manière contrôlée.

---

### 📦 Pourquoi utiliser `context` ?

* ⏱️ **Timeouts** : interrompre une opération qui prend trop de temps.
* ❌ **Annulation** : stopper une opération quand elle n’est plus nécessaire.
* 🔐 **Propagation de valeurs** : comme un identifiant de requête, un token utilisateur, etc.

---

### 🧱 Structure de base

```go
ctx := context.Background()
```

ou

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
```

---

### 🔁 Exemples pratiques

#### 1. **Timeout**

```go
ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
defer cancel()

select {
case <-time.After(3 * time.Second):
    fmt.Println("Fini")
case <-ctx.Done():
    fmt.Println("Contexte terminé :", ctx.Err()) // timeout
}
```

#### 2. **Annulation manuelle**

```go
ctx, cancel := context.WithCancel(context.Background())

go func() {
    time.Sleep(1 * time.Second)
    cancel() // Annule le contexte
}()

<-ctx.Done()
fmt.Println("Opération annulée :", ctx.Err())
```

#### 3. **Passage de contexte à une fonction**

```go
func fetchData(ctx context.Context) error {
    select {
    case <-time.After(5 * time.Second):
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}
```

---

### 💡 Bonnes pratiques

* Toujours appeler `cancel()` après `WithCancel`, `WithTimeout` ou `WithDeadline` pour libérer les ressources.
* Ne pas stocker un `context.Context` dans une structure (comme un struct global).
* Utiliser `context` comme **premier argument** d'une fonction :

  ```go
  func doSomething(ctx context.Context, ...) {}
  ```

---

Souhaites-tu un exemple intégré dans un handler HTTP ou une base de données avec `sqlx` ?

### avec un handler HTTP

Parfait ! Voici un **exemple complet d'utilisation de `context` dans un handler HTTP** avec le framework **Gin**. On y applique un **timeout** pour éviter que les requêtes longues ne bloquent le serveur.

---

### 🛠️ Exemple : Handler avec timeout et contexte

```go
package main

import (
    "context"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()
    r.GET("/long-task", longTaskHandler)
    r.Run(":8080")
}

func longTaskHandler(c *gin.Context) {
    // Crée un contexte avec timeout de 2 secondes
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    resultChan := make(chan string, 1)

    // Lance une tâche longue dans une goroutine
    go func() {
        time.Sleep(3 * time.Second) // Simule une tâche lente
        resultChan <- "Tâche terminée avec succès"
    }()

    select {
    case result := <-resultChan:
        c.JSON(http.StatusOK, gin.H{"message": result})
    case <-ctx.Done():
        // Timeout ou annulation
        c.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout de la requête"})
    }
}
```

---

### 📌 Points clés

* `context.WithTimeout` crée un contexte avec délai d’expiration.
* `ctx.Done()` est un canal fermé lorsque le contexte est annulé ou expire.
* Le `select` permet de réagir soit au résultat, soit à l'expiration.

---

### 🔄 Variante avec une fonction séparée

Si tu veux une structure plus propre :

```go
func longTask(ctx context.Context) (string, error) {
    select {
    case <-time.After(3 * time.Second):
        return "Terminé", nil
    case <-ctx.Done():
        return "", ctx.Err()
    }
}
```

Et dans le handler :

```go
func longTaskHandler(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    result, err := longTask(ctx)
    if err != nil {
        c.JSON(http.StatusRequestTimeout, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": result})
}
```

---

### plus

Voici un fichier `context.md` qui explique le fonctionnement du package `context` en Go, avec exemples et bonnes pratiques, adapté à une lecture claire pour documentation interne ou onboarding :

---

## 📘 `context` en Go

Le package `context` de Go permet de contrôler l’exécution des fonctions longues ou concurrentes à l’aide d’un **timeout**, d’une **annulation**, ou pour **passer des métadonnées** à travers des appels imbriqués.

---

### 🔧 Pourquoi utiliser `context` ?

* **Limiter la durée d'une requête** (timeout)
* **Annuler une opération en cascade**
* **Transmettre des informations liées à la requête** (ex : userID, traceID)
* **Prévenir les fuites de goroutines**

---

### 🧱 Création de base

```go
ctx := context.Background() // racine du contexte, souvent utilisé dans main()
ctx := context.TODO()       // utilisé comme placeholder

// Avec annulation
ctx, cancel := context.WithCancel(parentCtx)
defer cancel()

// Avec timeout
ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
defer cancel()

// Avec deadline (temps absolu)
ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(2*time.Second))
defer cancel()
```

---

### 🧵 Exemple avec `select` et `ctx.Done()`

```go
select {
case res := <-someChan:
    // opération réussie
case <-ctx.Done():
    // opération annulée ou expirée
    return ctx.Err()
}
```

---

### 🧪 Exemple dans un handler HTTP avec Gin

```go
func longTaskHandler(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    resultChan := make(chan string, 1)

    go func() {
        time.Sleep(3 * time.Second) // tâche simulée
        resultChan <- "Terminé"
    }()

    select {
    case result := <-resultChan:
        c.JSON(http.StatusOK, gin.H{"message": result})
    case <-ctx.Done():
        c.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout atteint"})
    }
}
```

---

### ✅ Bonnes pratiques

* Toujours appeler `defer cancel()` après `WithCancel`, `WithTimeout` ou `WithDeadline`.
* Passer le `context.Context` **en premier argument** d'une fonction.
* Ne **pas stocker** un contexte dans une structure.
* Préférer `context.Background()` comme racine sauf si dans un handler HTTP, où l’on utilise `c.Request.Context()`.

---

### 🧬 Passage de valeurs dans le contexte

```go
type contextKey string

const userIDKey contextKey = "userID"

ctx := context.WithValue(context.Background(), userIDKey, "123")

func getUserID(ctx context.Context) string {
    val := ctx.Value(userIDKey)
    if val == nil {
        return ""
    }
    return val.(string)
}
```

🔴 **Attention** : `WithValue` est à utiliser avec parcimonie. Préférer des paramètres explicites dans la majorité des cas.

---

Souhaites-tu que je te génère aussi une version `.md` prête à être enregistrée comme fichier ?

Voici le contenu prêt à être enregistré dans un fichier nommé `context.md` :

---

````markdown
# 📘 Utilisation du `context` en Go

Le package `context` de Go permet de contrôler l’exécution des fonctions longues ou concurrentes à l’aide d’un **timeout**, d’une **annulation**, ou pour **passer des métadonnées** à travers des appels imbriqués.

---

## 🔧 Pourquoi utiliser `context` ?

- ✅ Limiter la durée d'une requête (timeout)
- ❌ Annuler une opération en cascade
- 📦 Transmettre des informations liées à la requête (ex : userID, traceID)
- 🔐 Prévenir les fuites de goroutines

---

## 🧱 Création de base

```go
ctx := context.Background() // racine du contexte, souvent utilisé dans main()
ctx := context.TODO()       // utilisé comme placeholder

// Avec annulation
ctx, cancel := context.WithCancel(parentCtx)
defer cancel()

// Avec timeout
ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
defer cancel()

// Avec deadline (temps absolu)
ctx, cancel := context.WithDeadline(parentCtx, time.Now().Add(2*time.Second))
defer cancel()
````

---

## 🧵 Exemple avec `select` et `ctx.Done()`

```go
select {
case res := <-someChan:
    // opération réussie
case <-ctx.Done():
    // opération annulée ou expirée
    return ctx.Err()
}
```

---

## 🧪 Exemple dans un handler HTTP avec Gin

```go
func longTaskHandler(c *gin.Context) {
    ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
    defer cancel()

    resultChan := make(chan string, 1)

    go func() {
        time.Sleep(3 * time.Second) // tâche simulée
        resultChan <- "Terminé"
    }()

    select {
    case result := <-resultChan:
        c.JSON(http.StatusOK, gin.H{"message": result})
    case <-ctx.Done():
        c.JSON(http.StatusRequestTimeout, gin.H{"error": "Timeout atteint"})
    }
}
```

---

## ✅ Bonnes pratiques

* Toujours appeler `defer cancel()` après `WithCancel`, `WithTimeout` ou `WithDeadline`.
* Passer le `context.Context` **en premier argument** d'une fonction.
* Ne **pas stocker** un contexte dans une structure.
* Préférer `context.Background()` comme racine sauf dans un handler HTTP, où l’on utilise `c.Request.Context()`.

---

## 🧬 Passage de valeurs dans le contexte

```go
type contextKey string

const userIDKey contextKey = "userID"

ctx := context.WithValue(context.Background(), userIDKey, "123")

func getUserID(ctx context.Context) string {
    val := ctx.Value(userIDKey)
    if val == nil {
        return ""
    }
    return val.(string)
}
```

> ⚠️ `WithValue` est à utiliser avec parcimonie. Préférer des paramètres explicites quand c'est possible.

---

```yaml

---

Souhaites-tu que je le génère en fichier `.md` téléchargeable ?
```
