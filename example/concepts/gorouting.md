Le système de **goroutines** en **Go (Golang)** est un mécanisme léger et puissant pour gérer la **concurrence**. Voici une explication claire et progressive :

---

### 🔹 Qu'est-ce qu'une goroutine ?

Une **goroutine** est une fonction qui s’exécute **concurrentement** avec d’autres goroutines dans le même programme.
Elles sont gérées par le **runtime de Go**, et non directement par le système d’exploitation.

---

### 🔹 Syntaxe de base

Pour lancer une goroutine :

```go
go maFonction()
```

Cela lance `maFonction()` **en parallèle** du reste du programme.

---

### 🔹 Caractéristiques principales

* **Très légères** : chaque goroutine utilise quelques kilo-octets de mémoire (beaucoup moins qu’un thread classique).
* **Multiplexées** : le runtime de Go répartit automatiquement les goroutines sur les threads disponibles du système.
* **Pas de gestion explicite des threads** : tout est abstrait, ce qui simplifie énormément la programmation concurrente.

---

### 🔹 Exemple simple

```go
package main

import (
    "fmt"
    "time"
)

func direBonjour() {
    for i := 0; i < 5; i++ {
        fmt.Println("Bonjour")
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    go direBonjour() // lance direBonjour en tant que goroutine
    for i := 0; i < 5; i++ {
        fmt.Println("Salut")
        time.Sleep(100 * time.Millisecond)
    }
}
```

**Sortie possible :**

```shell
Salut
Bonjour
Salut
Bonjour
...
```

Les lignes "Salut" et "Bonjour" s'entremêlent, car les deux fonctions tournent **en parallèle**.

---

### 🔹 Synchronisation

Comme les goroutines peuvent s’exécuter en même temps, tu dois parfois **synchroniser** leur exécution.

Go propose plusieurs outils :

* **Channels** : pour communiquer entre goroutines.
* **WaitGroups** : pour attendre la fin de plusieurs goroutines.

**Exemple avec `sync.WaitGroup`** :

```go
import (
    "fmt"
    "sync"
)

func travail(wg *sync.WaitGroup) {
    fmt.Println("Travail en cours...")
    wg.Done() // signale que la goroutine est terminée
}

func main() {
    var wg sync.WaitGroup
    wg.Add(1) // on attend une goroutine
    go travail(&wg)
    wg.Wait() // on attend la fin
    fmt.Println("Fini")
}
```

---

### 🔹 Avantages des goroutines

* Haute performance : peut en lancer des **milliers** sans gros coût mémoire.
* Simplicité de syntaxe.
* Intégration fluide avec les channels pour la communication.

---

En Go, les appels HTTP sont **asynchrones par nature** dans le sens où ils tournent dans des goroutines si tu veux, mais si tu n’utilises **pas de goroutine**, l’appel est **bloquant** : ton programme **attend la réponse HTTP avant de continuer**.

Voici un exemple clair 👇

---

### 🧩 Exemple simple — Attendre la réponse avant d’appeler une fonction

```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    // 1️⃣ On fait un appel HTTP bloquant
    resp, err := http.Get("https://api.github.com")
    if err != nil {
        fmt.Println("Erreur HTTP :", err)
        return
    }
    defer resp.Body.Close()

    // 2️⃣ Lecture du corps de la réponse
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Erreur lecture body :", err)
        return
    }

    fmt.Println("Réponse reçue !")

    // 3️⃣ Appel de ta fonction une fois la réponse obtenue
    processResponse(body)
}

func processResponse(data []byte) {
    fmt.Println("Traitement de la réponse...")
    fmt.Println(string(data[:200])) // affiche les 200 premiers caractères
}
```

➡️ Ici, le programme **attend** que `http.Get()` et `ioutil.ReadAll()` soient terminés avant de passer à `processResponse`.

---

### 🧵 Exemple avec goroutine (si tu veux attendre explicitement)

Si tu veux faire **plusieurs appels HTTP en parallèle** et **attendre qu’ils soient tous terminés**, tu peux utiliser un **`sync.WaitGroup`** :

```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
)

func fetch(url string, wg *sync.WaitGroup) {
    defer wg.Done() // signale la fin
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Erreur:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    fmt.Println("Réponse de", url, ":", len(body), "octets")
}

func main() {
    var wg sync.WaitGroup

    urls := []string{
        "https://api.github.com",
        "https://golang.org",
        "https://example.com",
    }

    // Lance les appels en parallèle
    for _, url := range urls {
        wg.Add(1)
        go fetch(url, &wg)
    }

    // Attend que tous les appels soient terminés
    wg.Wait()

    fmt.Println("Tous les appels HTTP sont terminés.")
    finalFunction()
}

func finalFunction() {
    fmt.Println("Appel de la fonction finale après les requêtes HTTP.")
}
```

---

### 🧠 En résumé

* Si tu veux **attendre une seule requête** → fais simplement un `http.Get()` (il est bloquant).
* Si tu veux **attendre plusieurs requêtes parallèles** → utilise un `sync.WaitGroup`.
* Si tu veux attendre une seule goroutine → utilise un **canal (`chan`)** ou **WaitGroup** aussi.

---

Parfait 👍
Voici un exemple clair où **l’appel HTTP est lancé dans une goroutine**, et **le programme attend la réponse avant d’appeler une autre fonction** — grâce à un **canal (`chan`)** ou un **`sync.WaitGroup`**.

---

## 🧵 Exemple 1 — Avec un canal (`chan`)

Ce modèle est très idiomatique en Go :
la goroutine fait la requête, envoie le résultat sur un canal, et le main attend la réponse avant de continuer.

```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
)

// fetchURL lance un appel HTTP et renvoie le résultat via un canal
func fetchURL(url string, ch chan<- []byte) {
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Erreur HTTP:", err)
        close(ch)
        return
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println("Erreur lecture body:", err)
        close(ch)
        return
    }

    // Envoi de la réponse sur le canal
    ch <- body
    close(ch)
}

func processResponse(data []byte) {
    fmt.Println("Traitement de la réponse...")
    fmt.Println(string(data[:200])) // affiche les 200 premiers caractères
}

func main() {
    ch := make(chan []byte)

    fmt.Println("Lancement de la requête HTTP...")
    go fetchURL("https://api.github.com", ch)

    // On attend la réponse sur le canal (bloquant)
    body := <-ch

    // Appel de la fonction une fois la réponse reçue
    if body != nil {
        processResponse(body)
    } else {
        fmt.Println("Erreur lors de la récupération des données.")
    }

    fmt.Println("Programme terminé.")
}
```

### 🔍 Ce qu’il fait :

* `fetchURL()` tourne dans une goroutine (asynchrone).
* `main()` attend la réponse via `<-ch`, ce qui **bloque jusqu’à réception**.
* Dès que la réponse est reçue, on appelle `processResponse()`.

---

## 🪄 Exemple 2 — Avec `sync.WaitGroup` + variable partagée

Une autre façon (souvent utilisée quand tu veux récupérer plusieurs résultats) :

```go
package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
)

func fetch(url string, wg *sync.WaitGroup, result *string) {
    defer wg.Done()
    resp, err := http.Get(url)
    if err != nil {
        fmt.Println("Erreur HTTP:", err)
        return
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)
    *result = string(body)
}

func processResponse(data string) {
    fmt.Println("Traitement terminé :")
    fmt.Println(data[:200]) // affiche les 200 premiers caractères
}

func main() {
    var wg sync.WaitGroup
    var response string

    wg.Add(1)
    go fetch("https://api.github.com", &wg, &response)

    // Attend que la goroutine se termine
    wg.Wait()

    processResponse(response)
}
```

---

### ⚖️ Différences entre les deux approches

| Méthode                         | Avantage                                                           | Inconvénient                                                 |
| ------------------------------- | ------------------------------------------------------------------ | ------------------------------------------------------------ |
| **Canal (`chan`)**              | Idiomatique, simple à composer, pas besoin de variable globale     | Moins pratique pour récupérer plusieurs valeurs nommées      |
| **WaitGroup + pointeur/struct** | Facile à combiner avec plusieurs goroutines et résultats complexes | Nécessite attention à la concurrence (éviter les data races) |

---

Bien sûr ✅
On va reprendre l’exemple avec **canal (`chan`)**, car il est clair et idiomatique, et on va y ajouter **un test unitaire** propre en Go.

---

## 🧩 Code source : `main.go`

Voici la version légèrement refactorisée pour être testable :

```go
package main

import (
	"fmt"
	"io"
	"net/http"
)

// fetchURL lance un appel HTTP et renvoie la réponse (ou une erreur)
func fetchURL(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("erreur HTTP: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("erreur lecture body: %w", err)
	}

	return body, nil
}

// processResponse fait quelque chose avec les données (ex : affichage)
func processResponse(data []byte) {
	fmt.Println("Traitement de la réponse...")
	fmt.Println(string(data[:100])) // affiche 100 premiers caractères
}

func main() {
	fmt.Println("Lancement de la requête HTTP...")
	body, err := fetchURL("https://api.github.com")
	if err != nil {
		fmt.Println("Erreur:", err)
		return
	}

	processResponse(body)
}
```

---

## 🧪 Test unitaire : `main_test.go`

On va tester `fetchURL()` en **mockant un serveur HTTP local** avec `httptest` (lib standard de Go).

```go
package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestFetchURLOK vérifie que fetchURL lit correctement la réponse d’un serveur
func TestFetchURLOK(t *testing.T) {
	// Crée un faux serveur HTTP
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, `{"message":"Hello, world!"}`)
	}))
	defer mockServer.Close()

	// Appel de la fonction testée
	body, err := fetchURL(mockServer.URL)
	if err != nil {
		t.Fatalf("fetchURL a échoué : %v", err)
	}

	got := string(body)
	want := `{"message":"Hello, world!"}`

	if got != want {
		t.Errorf("Résultat inattendu :\n got  = %q\n want = %q", got, want)
	}
}

// TestFetchURLError vérifie la gestion d’erreurs
func TestFetchURLError(t *testing.T) {
	// URL invalide pour forcer une erreur
	_, err := fetchURL("http://adresse-invalide")
	if err == nil {
		t.Error("fetchURL aurait dû retourner une erreur pour une URL invalide")
	}
}
```

---

## ▶️ Lancer le test

```bash
go test -v
```

### Sortie attendue :

```
=== RUN   TestFetchURLOK
--- PASS: TestFetchURLOK (0.00s)
=== RUN   TestFetchURLError
--- PASS: TestFetchURLError (0.00s)
PASS
ok      example/test   0.003s
```

---

## 💡 Ce qu’on teste ici

* ✅ Que `fetchURL()` reçoit bien les données d’un serveur HTTP valide.
* ✅ Que `fetchURL()` renvoie une erreur si l’URL est invalide.

---

