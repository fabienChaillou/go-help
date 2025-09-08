Le **runtime** en Go fait référence à l’**environnement d'exécution** qui gère les **comportements bas niveau** du langage pendant que ton programme tourne. C’est une partie **intégrée du langage Go**, pas une VM comme en Java, mais un **ensemble de fonctions et mécanismes** natifs.

---

## 🔹 Que fait le runtime en Go ?

Voici les principales responsabilités du **runtime Go** :

### 1. **Goroutines et planification**

Le runtime :

* Lance et gère les **goroutines**
* Effectue le **scheduler** (répartition automatique des goroutines sur les threads du système)
* Gère le **M\:N threading model** : plusieurs goroutines sont multiplexées sur un nombre réduit de threads.

📌 Le **planificateur** repose sur le modèle **GMP** :

* **G** = Goroutine
* **M** = Thread du système (Machine)
* **P** = Processeur logique (scheduler local)

---

### 2. **Garbage collector (GC)**

Go intègre un **ramasse-miettes automatique** :

* Libère la mémoire des objets non utilisés.
* Fonctionne **en parallèle** avec les goroutines.
* Optimisé pour avoir une **latence faible** (surtout après Go 1.5+).

---

### 3. **Gestion de la mémoire**

* Allocation/désallocation de mémoire.
* Gestion du **heap** et du **stack**.
* Le runtime adapte la taille des stacks des goroutines dynamiquement (ils commencent très petits, quelques ko).

---

### 4. **Reflection & Type Info**

* Permet l’inspection des types à l’exécution via le package `reflect`.
* Utilisé en interne pour le typage dynamique, par exemple avec les interfaces.

---

### 5. **Panics et recover**

* Le runtime gère les **panic**, un équivalent de crash contrôlé.
* Tu peux intercepter un panic avec `recover()` dans un `defer`.

```go
func main() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Erreur rattrapée :", r)
        }
    }()
    panic("Quelque chose s'est mal passé")
}
```

---

### 6. **Timers, select, channels**

* Le runtime synchronise et gère les **channels**, les opérations `select`, les timers et délais (`time.After`, `time.Sleep`, etc.).

---

## 🔹 Accès au runtime depuis ton code

Tu peux accéder à certaines fonctions internes via le package `runtime` :

```go
import "runtime"

func main() {
    fmt.Println(runtime.NumCPU())     // Nombre de CPU
    fmt.Println(runtime.NumGoroutine()) // Nombre de goroutines en cours
    runtime.GC()                     // Force un GC
}
```

---

## 🔹 Exemple : voir le runtime en action

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func tache() {
    for i := 0; i < 5; i++ {
        fmt.Println("Tâche", i)
        time.Sleep(100 * time.Millisecond)
    }
}

func main() {
    go tache()
    fmt.Println("Goroutines actives :", runtime.NumGoroutine())
    time.Sleep(600 * time.Millisecond)
}
```

---

## 🔹 En résumé

| Fonction du runtime | Rôle principal                         |
| ------------------- | -------------------------------------- |
| Goroutine scheduler | Exécution concurrente efficace         |
| Garbage collector   | Nettoyage automatique de la mémoire    |
| Memory management   | Allocation et gestion des stacks/heaps |
| Panic/recover       | Gestion des erreurs fatales            |
| Reflection          | Inspection des types à l'exécution     |
| Timers/channels     | Outils de synchronisation              |

---

Souhaites-tu un schéma du modèle **GMP (Goroutine-Machine-Processor)** ou un exemple plus poussé de profiling avec le runtime ?
