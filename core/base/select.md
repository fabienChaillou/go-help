En Go, `select` est une instruction utilisée pour écouter plusieurs canaux (`chan`) en même temps. Elle permet de gérer la concurrence de manière élégante en attendant qu’un ou plusieurs canaux soient prêts pour la communication (lecture ou écriture). Elle est souvent utilisée avec les goroutines pour synchroniser les opérations concurrentes.

---

### Syntaxe de base

```go
select {
case msg1 := <-chan1:
    // Code exécuté si on reçoit un message de chan1
case chan2 <- msg2:
    // Code exécuté si on peut envoyer msg2 dans chan2
default:
    // (optionnel) Code exécuté si aucun canal n'est prêt
}
```

---

### Exemple simple

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "message de ch1"
    }()

    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "message de ch2"
    }()

    select {
    case msg := <-ch1:
        fmt.Println("Reçu:", msg)
    case msg := <-ch2:
        fmt.Println("Reçu:", msg)
    case <-time.After(3 * time.Second):
        fmt.Println("Timeout après 3 secondes")
    }
}
```

### Comportement :

* `select` attend que l’un des canaux soit prêt.
* Si plusieurs cas sont prêts en même temps, Go en choisit un **au hasard**.
* Si aucun cas n’est prêt et qu’il n’y a pas de `default`, `select` **bloque**.
* Si `default` est présent, il est exécuté immédiatement si rien n’est prêt.

---

### Cas d’usage typiques

* Lire sur plusieurs canaux en parallèle.
* Mettre en œuvre des timeout (`time.After`).
* Multiplexer des résultats venant de plusieurs goroutines.
* Implémenter une fermeture de service avec un canal de signal.

Souhaites-tu un exemple plus complet avec gestion d'annulation ou multiplexage ?

Voici un exemple complet en Go qui montre comment utiliser `select` pour gérer :

1. **Multiplexage** de messages depuis deux goroutines.
2. **Timeout** avec `time.After`.
3. **Annulation propre** avec un canal de signal (`done`).

---

### 🔧 Code d'exemple

```go
package main

import (
    "fmt"
    "time"
)

func worker(name string, out chan<- string, done <-chan struct{}) {
    for i := 1; i <= 5; i++ {
        select {
        case <-done:
            fmt.Println(name, "annulé.")
            return
        case out <- fmt.Sprintf("%s envoie %d", name, i):
            time.Sleep(time.Duration(500+100*i) * time.Millisecond)
        }
    }
}

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)
    done := make(chan struct{})

    go worker("Worker1", ch1, done)
    go worker("Worker2", ch2, done)

    timeout := time.After(3 * time.Second)

    for {
        select {
        case msg := <-ch1:
            fmt.Println("Reçu de ch1:", msg)
        case msg := <-ch2:
            fmt.Println("Reçu de ch2:", msg)
        case <-timeout:
            fmt.Println("⏰ Timeout atteint ! On annule les workers.")
            close(done)
            time.Sleep(1 * time.Second) // Laisse le temps aux goroutines de terminer
            return
        }
    }
}
```

---

### 🧠 Ce que fait ce programme :

* Deux "workers" envoient des messages dans deux canaux différents.
* La fonction `main()` lit ces messages via un `select`.
* Après 3 secondes, un timeout se déclenche et envoie un signal `done` aux goroutines pour les arrêter proprement.

---
