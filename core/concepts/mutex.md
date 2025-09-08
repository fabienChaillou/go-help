Un **mutex** (abréviation de *mutual exclusion*) en Go est un outil de synchronisation qui permet de **protéger l'accès à une ressource partagée** (comme une variable ou une structure) contre les accès concurrents par plusieurs goroutines.

---

### 🔹 Pourquoi un mutex ?

Quand plusieurs goroutines modifient une même variable **simultanément**, cela peut provoquer un **comportement non déterministe** (data race).

Le mutex permet de garantir que **seule une goroutine à la fois** accède à une section critique du code.

---

### 🔹 Utilisation basique

Go fournit un `sync.Mutex` dans le package `sync`.

```go
import (
    "fmt"
    "sync"
)

var compteur int
var mutex sync.Mutex

func incrementer(wg *sync.WaitGroup) {
    mutex.Lock()   // verrouille le mutex
    compteur++
    mutex.Unlock() // libère le mutex
    wg.Done()
}

func main() {
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go incrementer(&wg)
    }

    wg.Wait()
    fmt.Println("Valeur finale du compteur :", compteur)
}
```

🔸 Sans le mutex, le `compteur++` pourrait être **corrompu** à cause d'accès concurrents.
🔸 Avec `mutex.Lock()` et `Unlock()`, chaque goroutine attend son tour.

---

### 🔹 Points importants

* `Lock()` bloque la goroutine jusqu’à ce qu’elle puisse entrer dans la section critique.
* `Unlock()` libère le verrou.
* Toujours appeler `Unlock()` **même si une erreur survient** → souvent utilisé avec `defer` :

```go
mutex.Lock()
defer mutex.Unlock()
// section critique
```

---

### 🔹 Exemple : mauvaise utilisation sans mutex

```go
var x int

func incrementer() {
    for i := 0; i < 1000; i++ {
        x = x + 1
    }
}
```

Si tu lances cette fonction avec 10 goroutines, tu **ne garantis pas** que `x == 10000` à la fin, sauf si tu utilises un mutex.

---

### 🔹 Alternatives et compléments

* `sync.RWMutex` : un mutex avec **verrou lecture/écriture**, pour permettre plusieurs lectures concurrentes mais une seule écriture.
* `atomic` : pour des opérations simples sans mutex, comme `atomic.AddInt32`.

---
