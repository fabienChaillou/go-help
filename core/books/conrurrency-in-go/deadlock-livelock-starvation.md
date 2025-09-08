Voici une explication claire et contextualisée de **deadlock**, **livelock** et **starvation** en Go, avec des exemples simples :

---

## 🧱 1. **Deadlock (Interblocage)**

Un **deadlock** survient lorsque deux (ou plus) goroutines attendent mutuellement une ressource détenue par l'autre, et aucune ne peut continuer.

### 🔧 Exemple en Go :

```go
var mu1, mu2 sync.Mutex

func main() {
    go func() {
        mu1.Lock()
        time.Sleep(1 * time.Second)
        mu2.Lock() // bloqué car l'autre goroutine a mu2
        mu2.Unlock()
        mu1.Unlock()
    }()

    go func() {
        mu2.Lock()
        time.Sleep(1 * time.Second)
        mu1.Lock() // bloqué car la première a mu1
        mu1.Unlock()
        mu2.Unlock()
    }()

    time.Sleep(5 * time.Second)
}
```

### 🔍 Résultat :

Les deux goroutines se bloquent mutuellement — **deadlock**.

---

## 🔄 2. **Livelock (Interblocage actif)**

Un **livelock** est similaire au deadlock, mais ici les goroutines **continuent de s'exécuter**, en réessayant en boucle, sans progresser.

### 🔧 Exemple en Go :

```go
type Philosopher struct {
    name   string
    trying bool
}

func (p *Philosopher) tryEat(other *Philosopher) {
    for p.trying {
        if !other.trying {
            fmt.Println(p.name, "commence à manger")
            return
        }
        fmt.Println(p.name, "attend que", other.name, "finisse")
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    p1 := &Philosopher{name: "A", trying: true}
    p2 := &Philosopher{name: "B", trying: true}

    go p1.tryEat(p2)
    go p2.tryEat(p1)

    time.Sleep(3 * time.Second)
}
```

### 🔍 Résultat :

Les deux essaient de manger, attendent l'autre en boucle sans jamais réussir — **livelock**.

---

## ⌛ 3. **Starvation (Famine)**

La **starvation** survient lorsqu’une goroutine **n’obtient jamais l'accès à une ressource** car d'autres goroutines monopolisent cette ressource en continu.

### 🔧 Exemple en Go :

```go
var mu sync.Mutex

func main() {
    // Goroutine "gloutonne"
    go func() {
        for {
            mu.Lock()
            time.Sleep(100 * time.Millisecond)
            mu.Unlock()
            time.Sleep(10 * time.Millisecond)
        }
    }()

    // Goroutine affamée
    go func() {
        for {
            mu.Lock()
            fmt.Println("Goroutine affamée a obtenu le verrou")
            mu.Unlock()
            time.Sleep(100 * time.Millisecond)
        }
    }()

    time.Sleep(5 * time.Second)
}
```

### 🔍 Résultat :

La goroutine affamée **a peu de chances d'accéder au mutex** car l'autre la devance tout le temps — **starvation**.

---

Souhaites-tu que je t’aide à écrire un détecteur ou simulateur de ces cas en Go pour expérimenter ?

Parfait ! Voici un petit **simulateur interactif** en Go qui te permet de tester des cas de **deadlock**, **livelock** et **starvation**. Chaque exemple est activable via une constante.

---

### 🧪 Code complet du simulateur :

```go
package main

import (
	"fmt"
	"sync"
	"time"
)

// Active une des trois situations :
const (
	testDeadlock   = false
	testLivelock   = false
	testStarvation = true
)

func main() {
	switch {
	case testDeadlock:
		runDeadlock()
	case testLivelock:
		runLivelock()
	case testStarvation:
		runStarvation()
	default:
		fmt.Println("Active une des constantes pour lancer un test.")
	}
}

// ======= DEADLOCK =======
func runDeadlock() {
	var mu1, mu2 sync.Mutex

	go func() {
		mu1.Lock()
		fmt.Println("Goroutine 1 a verrouillé mu1")
		time.Sleep(1 * time.Second)
		mu2.Lock()
		fmt.Println("Goroutine 1 a verrouillé mu2")
		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		mu2.Lock()
		fmt.Println("Goroutine 2 a verrouillé mu2")
		time.Sleep(1 * time.Second)
		mu1.Lock()
		fmt.Println("Goroutine 2 a verrouillé mu1")
		mu1.Unlock()
		mu2.Unlock()
	}()

	time.Sleep(5 * time.Second)
	fmt.Println("Fin du test Deadlock (tu devrais voir un blocage)")
}

// ======= LIVELOCK =======
type Philosopher struct {
	name    string
	trying  bool
}

func (p *Philosopher) tryEat(other *Philosopher) {
	for p.trying {
		if !other.trying {
			fmt.Println(p.name, "commence à manger")
			return
		}
		fmt.Println(p.name, "attend que", other.name, "finisse")
		time.Sleep(200 * time.Millisecond)
	}
}

func runLivelock() {
	p1 := &Philosopher{name: "A", trying: true}
	p2 := &Philosopher{name: "B", trying: true}

	go p1.tryEat(p2)
	go p2.tryEat(p1)

	time.Sleep(5 * time.Second)
	fmt.Println("Fin du test Livelock (tu devrais voir les philosophes attendre indéfiniment)")
}

// ======= STARVATION =======
func runStarvation() {
	var mu sync.Mutex

	// Glouton : prend le verrou très souvent
	go func() {
		for {
			mu.Lock()
			time.Sleep(100 * time.Millisecond)
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	// Affamé : essaie périodiquement d'obtenir le verrou
	go func() {
		for {
			time.Sleep(150 * time.Millisecond)
			mu.Lock()
			fmt.Println("Goroutine affamée a obtenu le verrou")
			mu.Unlock()
		}
	}()

	time.Sleep(10 * time.Second)
	fmt.Println("Fin du test Starvation")
}
```

---

### ✅ Instructions :

1. Choisis le test à exécuter en changeant la valeur des constantes (`testDeadlock`, `testLivelock`, `testStarvation`).
2. Lance le programme :

   ```bash
   go run main.go
   ```

---

Souhaites-tu une version avec une interface Web simple (en Go + HTML avec `net/http`) pour contrôler les tests via un navigateur ?

Voici une **version web** simple du simulateur (utilisant `net/http` et du HTML) qui te permet de lancer les tests **Deadlock**, **Livelock**, ou **Starvation** depuis ton navigateur.

---

### 🌐 Fonctionnalités :

* Interface web avec 3 boutons
* Chaque bouton déclenche un test côté serveur (via des `goroutines`)
* Affichage du résultat côté terminal (log)

---

### 📁 Fichiers requis

#### 1. `main.go`

```go
package main

import (
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"
)

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/test", testHandler)
	fmt.Println("Serveur sur http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := `
<!DOCTYPE html>
<html>
<head><title>Simulateur Concurrence Go</title></head>
<body style="font-family:sans-serif;text-align:center;margin-top:50px">
	<h1>Simulateur Deadlock / Livelock / Starvation</h1>
	<button onclick="location.href='/test?mode=deadlock'">🧱 Deadlock</button>
	<button onclick="location.href='/test?mode=livelock'">🔄 Livelock</button>
	<button onclick="location.href='/test?mode=starvation'">⌛ Starvation</button>
	<p>Regarde les logs dans le terminal !</p>
</body>
</html>
`
	t := template.Must(template.New("web").Parse(tmpl))
	t.Execute(w, nil)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	mode := r.URL.Query().Get("mode")
	switch mode {
	case "deadlock":
		fmt.Println("=== DEADLOCK ===")
		runDeadlock()
	case "livelock":
		fmt.Println("=== LIVELOCK ===")
		runLivelock()
	case "starvation":
		fmt.Println("=== STARVATION ===")
		runStarvation()
	default:
		fmt.Println("Mode inconnu")
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// ========== DEADLOCK ==========
func runDeadlock() {
	var mu1, mu2 sync.Mutex

	go func() {
		mu1.Lock()
		fmt.Println("Goroutine 1: mu1 locked")
		time.Sleep(1 * time.Second)
		mu2.Lock()
		fmt.Println("Goroutine 1: mu2 locked")
		mu2.Unlock()
		mu1.Unlock()
	}()

	go func() {
		mu2.Lock()
		fmt.Println("Goroutine 2: mu2 locked")
		time.Sleep(1 * time.Second)
		mu1.Lock()
		fmt.Println("Goroutine 2: mu1 locked")
		mu1.Unlock()
		mu2.Unlock()
	}()
}

// ========== LIVELOCK ==========
type Philosopher struct {
	name   string
	trying bool
}

func (p *Philosopher) tryEat(other *Philosopher) {
	for p.trying {
		if !other.trying {
			fmt.Println(p.name, "mange.")
			return
		}
		fmt.Println(p.name, "attend que", other.name, "finisse.")
		time.Sleep(200 * time.Millisecond)
	}
}

func runLivelock() {
	p1 := &Philosopher{name: "A", trying: true}
	p2 := &Philosopher{name: "B", trying: true}

	go p1.tryEat(p2)
	go p2.tryEat(p1)
}

// ========== STARVATION ==========
func runStarvation() {
	var mu sync.Mutex

	go func() {
		for {
			mu.Lock()
			time.Sleep(100 * time.Millisecond)
			mu.Unlock()
			time.Sleep(10 * time.Millisecond)
		}
	}()

	go func() {
		for {
			time.Sleep(150 * time.Millisecond)
			mu.Lock()
			fmt.Println("Goroutine affamée a obtenu le verrou")
			mu.Unlock()
		}
	}()
}
```

---

### ▶️ Instructions :

1. Crée un fichier `main.go` avec ce code.
2. Lance le serveur :

   ```bash
   go run main.go
   ```
3. Ouvre ton navigateur : [http://localhost:8080](http://localhost:8080)
4. Clique sur un bouton, puis observe les **logs du terminal** pour voir le comportement.

---

Souhaites-tu que je t’ajoute un **terminal dans la page web** pour voir les logs en direct via WebSocket ?
