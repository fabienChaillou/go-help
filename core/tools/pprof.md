`pprof` est un outil de **profilage** en Go, intégré à la bibliothèque standard, qui permet d’analyser les performances d’une application Go (CPU, mémoire, goroutines, etc.). Il est très utile pour diagnostiquer des lenteurs, des fuites mémoire, ou simplement comprendre le comportement d’un programme.

---

### 🔧 Fonctionnement général

`pprof` collecte des **profils d’exécution** à différents niveaux :

* **CPU** : quelles fonctions consomment le plus de temps CPU.
* **Mémoire (heap)** : quelles allocations de mémoire sont faites.
* **Goroutines** : quelles sont les goroutines actives.
* **Blocages (blocking)** : où le programme est bloqué (mutex, canal…).

---

### 📦 Importer pprof

```go
import _ "net/http/pprof"
```

> L'import underscore (`_`) sert à exécuter le `init()` du package, qui enregistre automatiquement des endpoints HTTP utiles.

---

### 🌐 Exposer pprof via HTTP

Souvent, on lance un serveur HTTP pour permettre l’accès aux endpoints de profilage :

```go
package main

import (
    _ "net/http/pprof"
    "log"
    "net/http"
)

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // Le reste de ton programme ici
}
```

---

### 📊 Endpoints disponibles

Une fois lancé, tu peux accéder aux profils via :

* [http://localhost:6060/debug/pprof/](http://localhost:6060/debug/pprof/)
* [http://localhost:6060/debug/pprof/profile?seconds=30](http://localhost:6060/debug/pprof/profile?seconds=30) (CPU profile de 30s)
* [http://localhost:6060/debug/pprof/heap](http://localhost:6060/debug/pprof/heap) (profil mémoire)
* [http://localhost:6060/debug/pprof/goroutine](http://localhost:6060/debug/pprof/goroutine) (goroutines actives)

---

### 🧪 Utilisation avec `go tool pprof`

```bash
go tool pprof http://localhost:6060/debug/pprof/profile
```

Puis, dans l’interface interactive :

* `top` : montre les fonctions les plus coûteuses
* `list maFonction` : montre le détail dans une fonction
* `web` : génère un graphe interactif (nécessite Graphviz)

---

### 🧠 Exemple simple

```go
func main() {
    go func() {
        log.Println(http.ListenAndServe(":6060", nil))
    }()

    for {
        doSomething()
    }
}

func doSomething() {
    time.Sleep(10 * time.Millisecond)
}
```

Tu peux ensuite lancer :

```bash
go run main.go
go tool pprof http://localhost:6060/debug/pprof/profile
```

---

### 📌 Astuces

* Utilise `net/http/pprof` en dev seulement (pas en prod sans restriction d'accès !).
* Tu peux aussi utiliser `runtime/pprof` manuellement si tu veux contrôler précisément quand et quoi profiler.
* Compatible avec des outils comme **GoLand**, **Visual Studio Code**, ou **Perfetto** pour analyse visuelle.

---
