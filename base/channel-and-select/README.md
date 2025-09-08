Voici un exemple simple en Go qui utilise un **channel** et une **instruction `select`**. Cela illustre comment on peut écouter plusieurs canaux en même temps et réagir à celui qui reçoit une donnée en premier.

---

### 🎯 Objectif

Créer deux goroutines qui envoient des messages après un certain délai, et utiliser `select` pour recevoir celui qui arrive en premier.

---

### ✅ Exemple complet

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	// Création de deux canaux
	ch1 := make(chan string)
	ch2 := make(chan string)

	// Goroutine qui envoie un message après 2 secondes
	go func() {
		time.Sleep(2 * time.Second)
		ch1 <- "Message de ch1"
	}()

	// Goroutine qui envoie un message après 1 seconde
	go func() {
		time.Sleep(1 * time.Second)
		ch2 <- "Message de ch2"
	}()

	// On écoute avec select
	select {
	case msg1 := <-ch1:
		fmt.Println("Reçu :", msg1)
	case msg2 := <-ch2:
		fmt.Println("Reçu :", msg2)
	default:
		fmt.Println("Default msg")
	}
}
```

---

### 📌 Explication

* `make(chan string)` crée un canal pour échanger des chaînes.
* Deux goroutines envoient un message après un certain temps.
* L’instruction `select` attend qu’un des deux canaux reçoive une valeur.
* Le premier canal à envoyer un message déclenche le `case` correspondant.

---

### 🧠 Résultat attendu (à peu près) :

```
Reçu : Message de ch2
```

> Parce que `ch2` attend 1 seconde alors que `ch1` attend 2 secondes.

Souhaites-tu un exemple où on boucle avec `select` ou avec un `default` pour gérer les canaux vides ?

Parfait ! Voici un **exemple avec une boucle `for` et un `select`** qui permet d’attendre plusieurs messages dans le temps, **avec un `default`** pour gérer les cas où aucun canal n'est prêt (i.e., non bloquant).

---

### 🔁 Exemple : boucle avec `select` + `default`

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// Goroutine 1 : envoie toutes les 2 secondes
	go func() {
		for {
			time.Sleep(2 * time.Second)
			ch1 <- "Message de ch1"
		}
	}()

	// Goroutine 2 : envoie toutes les 3 secondes
	go func() {
		for {
			time.Sleep(3 * time.Second)
			ch2 <- "Message de ch2"
		}
	}()

	// Boucle infinie avec select
	for {
		select {
		case msg1 := <-ch1:
			fmt.Println("Reçu de ch1 :", msg1)
		case msg2 := <-ch2:
			fmt.Println("Reçu de ch2 :", msg2)
		default:
			// Aucun canal prêt : ne pas bloquer
			fmt.Println("Aucun message, on continue...")
			time.Sleep(500 * time.Millisecond) // pour éviter un loop trop rapide
		}
	}
}
```

---

### 🧠 Explication

* `select` surveille les deux canaux.
* Si aucun n'est prêt, on entre dans le `default`, ce qui **évite le blocage**.
* `time.Sleep(500ms)` évite de surcharger la CPU pendant l’attente.

---

### 🔍 Exemple de sortie (extrait) :

```
Aucun message, on continue...
Aucun message, on continue...
Reçu de ch1 : Message de ch1
Aucun message, on continue...
Reçu de ch2 : Message de ch2
...
```

Souhaites-tu voir comment on peut **fermer un canal** et le gérer dans un `select` aussi ?

Très bien ! Voici un exemple simple qui montre **comment fermer un canal** et **gérer sa fermeture dans un `select`**.

---

### 🔐 Objectif :

* Envoyer quelques messages dans un canal.
* Fermer le canal une fois terminé.
* Utiliser `select` pour lire les messages et détecter quand le canal est fermé.

---

### ✅ Exemple avec fermeture de canal

```go
package main

import (
	"fmt"
)

func main() {
	ch := make(chan string)

	// Goroutine qui envoie 3 messages, puis ferme le canal
	go func() {
		ch <- "Message 1"
		ch <- "Message 2"
		ch <- "Message 3"
		close(ch) // très important
	}()

	// Lire depuis le canal jusqu'à ce qu'il soit fermé
	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				fmt.Println("Canal fermé, on quitte.")
				return
			}
			fmt.Println("Reçu :", msg)
		}
	}
}
```

---

### 🔍 Explication

* `msg, ok := <-ch` :

  * `msg` contient la valeur lue,
  * `ok == false` signifie que le **canal est fermé**.
* Quand `ok` est `false`, on peut arrêter la boucle (`return` ou `break`).

---

### 🧠 Sortie attendue :

```
Reçu : Message 1
Reçu : Message 2
Reçu : Message 3
Canal fermé, on quitte.
```

---

Tu veux voir un exemple avec **plusieurs canaux où certains se ferment pendant que d'autres continuent** ?

Super choix ! Voici un exemple un peu plus réaliste avec **plusieurs canaux**, dont **l’un se ferme pendant que l’autre continue**.

---

### 🎯 Objectif :

* `ch1` envoie 3 messages puis se ferme.
* `ch2` continue à envoyer des messages toutes les 2 secondes.
* On utilise `select` pour lire les deux, et on détecte la fermeture de `ch1`.

---

### ✅ Exemple : plusieurs canaux, un qui se ferme

```go
package main

import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan string)
	ch2 := make(chan string)

	// ch1 : envoie 3 messages, puis se ferme
	go func() {
		for i := 1; i <= 3; i++ {
			ch1 <- fmt.Sprintf("ch1 message %d", i)
			time.Sleep(1 * time.Second)
		}
		close(ch1)
	}()

	// ch2 : continue à envoyer
	go func() {
		for i := 1; ; i++ {
			ch2 <- fmt.Sprintf("ch2 message %d", i)
			time.Sleep(2 * time.Second)
		}
	}()

	// boucle principale
	for {
		select {
		case msg, ok := <-ch1:
			if !ok {
				fmt.Println("✅ ch1 est fermé")
				ch1 = nil // évite de re-tester un canal fermé
			} else {
				fmt.Println("📥 Reçu de ch1 :", msg)
			}
		case msg := <-ch2:
			fmt.Println("📥 Reçu de ch2 :", msg)
		}
	}
}
```

---

### 🔍 Explication importante :

* On met `ch1 = nil` après fermeture, car un `select` avec un canal fermé reste toujours prêt, ce qui **bloquerait l'accès à `ch2`** si on ne le fait pas.
* `ch2` tourne indéfiniment, simulant un flux continu.
* On voit clairement que `ch1` s'arrête, mais `ch2` continue à fonctionner.

---

### 🧠 Exemple de sortie :

```
📥 Reçu de ch1 : ch1 message 1
📥 Reçu de ch1 : ch1 message 2
📥 Reçu de ch2 : ch2 message 1
📥 Reçu de ch1 : ch1 message 3
✅ ch1 est fermé
📥 Reçu de ch2 : ch2 message 2
📥 Reçu de ch2 : ch2 message 3
...
```
