Les **channels** en Go sont des **tuyaux de communication** entre les goroutines. Ils permettent de **transmettre des données de manière sûre** et **synchronisée**, sans avoir à utiliser de mutex dans la plupart des cas.

---

### 🔹 Définition simple

Un **channel** permet à une goroutine d’envoyer une valeur, et à une autre de la recevoir.
C’est un **mécanisme de synchronisation intégré** : l'envoi bloque jusqu'à ce que la réception ait lieu, et vice versa (dans un channel non bufferisé).

---

### 🔹 Syntaxe de base

#### Création d’un channel :

```go
ch := make(chan int) // channel de type int
```

#### Envoi et réception :

```go
ch <- 42      // envoyer 42 dans le channel
x := <-ch     // recevoir depuis le channel
```

---

### 🔹 Exemple simple

```go
package main

import (
    "fmt"
)

func saluer(ch chan string) {
    ch <- "Bonjour depuis la goroutine" // envoie une valeur
}

func main() {
    ch := make(chan string)

    go saluer(ch)           // lance la goroutine
    message := <-ch         // attend et lit la valeur
    fmt.Println(message)
}
```

Sortie :

```
Bonjour depuis la goroutine
```

---

### 🔹 Channel bufferisé

Un **channel bufferisé** permet de stocker plusieurs messages sans attendre leur réception immédiate.

```go
ch := make(chan int, 3) // channel avec un buffer de 3

ch <- 1
ch <- 2
ch <- 3
```

🔸 Si tu essaies d’envoyer une 4ᵉ valeur sans qu’aucune n’ait été lue, **le programme bloque**.

---

### 🔹 Fermeture d’un channel

Tu peux fermer un channel pour signaler qu’il n’y aura plus d’envoi :

```go
close(ch)
```

Et pour lire jusqu’à ce qu’il soit fermé :

```go
for val := range ch {
    fmt.Println(val)
}
```

---

### 🔹 Exemple complet avec boucle

```go
func generateur(ch chan int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}

func main() {
    ch := make(chan int)
    go generateur(ch)

    for val := range ch {
        fmt.Println(val)
    }
}
```

---

### 🔹 Sélection avec `select`

Quand tu veux écouter plusieurs channels en même temps :

```go
select {
case msg1 := <-ch1:
    fmt.Println("Reçu de ch1:", msg1)
case msg2 := <-ch2:
    fmt.Println("Reçu de ch2:", msg2)
default:
    fmt.Println("Aucune donnée reçue")
}
```

---

### 🔹 Avantages des channels

* Synchronisation naturelle entre goroutines.
* Moins besoin de mutex.
* Bonne lisibilité du code concurrent.

---

Souhaites-tu que je te montre un exemple **channel vs mutex**, ou un cas plus avancé avec `select` et `timeout` ?
