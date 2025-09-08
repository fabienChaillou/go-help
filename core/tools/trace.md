Pour tracer les **goroutines** en Go, tu as plusieurs options selon ce que tu veux faire (debug, profiling, visualisation, etc.). Voici un aperçu des méthodes les plus courantes :

---

### 🔹 1. **Utiliser le package `runtime`**

Tu peux lister les goroutines et obtenir leur stack trace :

```go
import (
    "runtime"
    "os"
)

func dumpGoroutines() {
    buf := make([]byte, 1<<20) // 1 MB buffer
    stacklen := runtime.Stack(buf, true)
    os.Stdout.Write(buf[:stacklen])
}
```

👉 Appelle `dumpGoroutines()` à un moment stratégique pour voir ce qui tourne.

---

### 🔹 2. **Profiling avec `net/http/pprof`**

Ajoute ceci à ton code :

```go
import (
    _ "net/http/pprof"
    "log"
    "net/http"
)

func init() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
}
```

Puis, dans ton navigateur ou avec `curl` :

* Affiche les goroutines : [http://localhost:6060/debug/pprof/goroutine?debug=2](http://localhost:6060/debug/pprof/goroutine?debug=2)
* Génère un profil : `go tool pprof http://localhost:6060/debug/pprof/goroutine`

---

### 🔹 3. **Outils de tracing visuel**

* **Go tool trace** :

  1. Ajoute ce code pour générer un fichier de trace :

     ```go
     import (
         "os"
         "runtime/trace"
     )

     func main() {
         f, _ := os.Create("trace.out")
         defer f.Close()
         trace.Start(f)
         defer trace.Stop()

         // Ton code ici
     }
     ```
  2. Lance le trace viewer :

     ```
     go tool trace trace.out
     ```

  Cela ouvre une interface web avec le détail des goroutines.

---

### 🔹 4. **Delve (le débogueur Go)**

Tu peux exécuter ton programme avec [Delve](https://github.com/go-delve/delve) :

```bash
dlv debug
```

Dans le shell Delve :

```bash
goroutines
```

Cela affiche toutes les goroutines actives.

---

### 🔹 5. **Bibliothèques tierces**

* [Gops](https://github.com/google/gops) — inspecte les goroutines en temps réel.
* [go-trace](https://github.com/nikolaydubina/go-trace) — générateur de visualisation d’événements d'exécution.

---

Souhaites-tu un exemple complet avec un serveur qui déclenche plusieurs goroutines et trace leur activité ?

Parfait ! Voici un **exemple complet** d’un serveur HTTP en Go qui :

1. Lance plusieurs goroutines.
2. Expose une interface de diagnostic via `net/http/pprof`.
3. Peut générer un `trace.out` à analyser avec `go tool trace`.

---

### ✅ Exemple Complet

```go
package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/trace"
	"time"
)

func worker(id int) {
	for {
		fmt.Printf("Worker %d is working...\n", id)
		time.Sleep(2 * time.Second)
	}
}

func main() {
	// Démarrer le trace
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalf("Erreur lors de la création du fichier de trace: %v", err)
	}
	defer f.Close()
	if err := trace.Start(f); err != nil {
		log.Fatalf("Erreur lors du démarrage du trace: %v", err)
	}
	defer trace.Stop()

	// Lancer quelques goroutines
	for i := 1; i <= 3; i++ {
		go worker(i)
	}

	// Lancer le serveur pprof sur localhost:6060
	go func() {
		log.Println("Serveur pprof sur http://localhost:6060/debug/pprof/")
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	// Garder le serveur principal actif
	select {}
}
```

---

### 🧪 Tester l'exemple

1. **Lancer le programme** :

   ```bash
   go run main.go
   ```

2. **Accéder à pprof dans le navigateur** :

   * [http://localhost:6060/debug/pprof/](http://localhost:6060/debug/pprof/)
   * Voir toutes les goroutines :
     [http://localhost:6060/debug/pprof/goroutine?debug=2](http://localhost:6060/debug/pprof/goroutine?debug=2)

3. **Analyser le fichier de trace** :
   Après avoir arrêté le programme, exécute :

   ```bash
   go tool trace trace.out
   ```

   Cela ouvrira une interface Web avec des vues comme :

   * **Goroutines**
   * **Scheduler latency**
   * **Network blocking profile**
   * etc.

---
