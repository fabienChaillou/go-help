Très bonne question ! En Go, les **fonctions et types `builtin`** (intégrés) sont ceux que tu peux utiliser **sans avoir besoin de les importer** depuis un package. Ce sont des **éléments de base fournis directement par le langage**.

---

## 🧰 Qu’est-ce que `builtin` en Go ?

> Le mot `builtin` fait référence à tous les **types, fonctions et constantes** disponibles **par défaut** dans Go.

Ils sont déclarés dans un package spécial interne appelé `builtin`, que tu **n’as pas besoin d’importer**.

---

## 📦 1. Types built-in

Ce sont les **types de base** que tu connais déjà :

* **Numériques** : `int`, `int8`, `float64`, `complex128`, etc.
* **Booléen** : `bool`
* **Texte** : `string`, `rune`, `byte`
* **Pointeurs** : `*T`
* **Agrégats** : `array`, `slice`, `map`, `struct`, `interface`, `chan`

---

## ⚙️ 2. Fonctions built-in les plus utiles

Voici les principales :

| Fonction                  | Description                                    |
| ------------------------- | ---------------------------------------------- |
| `len(x)`                  | Longueur (slice, array, map, string, chan)     |
| `cap(x)`                  | Capacité (slice, array, chan)                  |
| `append(slice, elems...)` | Ajoute des éléments à un slice                 |
| `copy(dst, src)`          | Copie le contenu d’un slice dans un autre      |
| `close(ch)`               | Ferme un channel                               |
| `delete(m, key)`          | Supprime une entrée dans une map               |
| `complex(r, i)`           | Crée un nombre complexe                        |
| `real(c)` / `imag(c)`     | Parties réelle et imaginaire                   |
| `new(T)`                  | Alloue un pointeur vers un type T              |
| `make(T)`                 | Alloue une slice, map ou channel               |
| `panic(v)`                | Provoque une erreur fatale                     |
| `recover()`               | Intercepte une panique (dans un `defer`)       |
| `print()`, `println()`    | Version de débogage, à éviter en prod          |
| `unsafe.Sizeof()`         | Obtenir la taille mémoire (nécessite `unsafe`) |

---

## 🧪 Exemple

```go
package main

import "fmt"

func main() {
    s := []int{1, 2, 3}
    fmt.Println("len:", len(s))     // 3
    fmt.Println("cap:", cap(s))     // dépend

    s = append(s, 4)
    fmt.Println(s)                  // [1 2 3 4]

    m := map[string]int{"a": 1, "b": 2}
    delete(m, "a")
    fmt.Println(m)                  // map[b:2]
}
```

---

## 🧠 Différence entre `make` et `new`

| Fonction | Utilisée pour               | Retourne                      |
| -------- | --------------------------- | ----------------------------- |
| `new(T)` | Allouer n’importe quel type | `*T` (pointeur)               |
| `make()` | Spécifique à slice/map/chan | valeur directe (pas pointeur) |

```go
p := new(int)     // *int, pointe vers 0
s := make([]int, 5) // slice avec 5 éléments
```

---

## 🧩 Où voir tout ça ?

Tu peux voir tous les éléments `builtin` ici (officiel) :
👉 [https://pkg.go.dev/builtin](https://pkg.go.dev/builtin)

---

--------------------------------------------------------------------------------------------------------

Souhaites-tu un exemple montrant plusieurs fonctions builtin dans un mini projet Go ?

Super ! Voici un **mini projet Go** qui utilise plusieurs **fonctions et types `builtin`** en même temps, pour bien illustrer leur utilité dans un contexte réaliste.

---

## 🎯 Objectif :

Créer un programme qui :

1. Gère une liste de personnes (struct slice)
2. Utilise `append`, `len`, `make`, `delete`, `copy`
3. Affiche les infos avec `fmt.Println`

---

## 🧱 Code complet avec explication

```go
package main

import (
    "fmt"
)

// Struct représentant une personne
type Person struct {
    Name string
    Age  int
}

func main() {
    // make : création d’un slice de 0 élément, capacité 3
    people := make([]Person, 0, 3)

    // append : ajout d’éléments
    people = append(people, Person{Name: "Alice", Age: 30})
    people = append(people, Person{Name: "Bob", Age: 25})
    people = append(people, Person{Name: "Charlie", Age: 35})

    fmt.Println("Liste initiale (len =", len(people), ", cap =", cap(people), "):")
    for i, p := range people {
        fmt.Printf("%d. %s (%d ans)\n", i+1, p.Name, p.Age)
    }

    // Supposons qu’on veuille "supprimer" Bob (index 1)
    indexToRemove := 1
    people = append(people[:indexToRemove], people[indexToRemove+1:]...) // suppression

    fmt.Println("\nAprès suppression de Bob:")
    for _, p := range people {
        fmt.Println("-", p.Name)
    }

    // copy : copier la liste dans un nouveau slice
    copied := make([]Person, len(people))
    copy(copied, people)

    // Modifions l’original pour montrer que copied est indépendant
    people[0].Name = "Alicia"

    fmt.Println("\nCopie indépendante :")
    for _, p := range copied {
        fmt.Println("-", p.Name)
    }

    // map + delete
    ages := map[string]int{
        "Alice":   30,
        "Charlie": 35,
    }
    delete(ages, "Alice") // suppression d’une clé

    fmt.Println("\nMap des âges restants :", ages)
}
```

---

## 🔍 Fonctions `builtin` utilisées ici

| Fonction | Utilisation dans le code       |
| -------- | ------------------------------ |
| `make`   | Créer un slice ou une map      |
| `append` | Ajouter un élément à un slice  |
| `len`    | Obtenir la taille d’un slice   |
| `cap`    | Obtenir la capacité d’un slice |
| `copy`   | Copier un slice                |
| `delete` | Supprimer une clé dans une map |

---

## ✅ Sortie attendue (exemple) :

```
Liste initiale (len = 3 , cap = 3 ):
1. Alice (30 ans)
2. Bob (25 ans)
3. Charlie (35 ans)

Après suppression de Bob:
- Alice
- Charlie

Copie indépendante :
- Alice
- Charlie

Map des âges restants : map[Charlie:35]
```

---

--------------------------------------------------------------------------------------------------------

Souhaites-tu un exemple similaire mais avec des goroutines et channels intégrés ?

Parfait ! Voici un **mini projet Go** qui combine les **fonctions `builtin`**, les **structs**, **goroutines**, et **channels**. C'est un excellent exemple de **concurrence contrôlée et typée**.

---

## 🎯 Objectif :

* Traiter une liste de personnes (structs)
* Lancer une goroutine par personne
* Communiquer les résultats via un channel
* Utiliser `make`, `len`, `cap`, `append`, `close`, `range`, `go`, `chan`

---

## 🧱 Code complet

```go
package main

import (
    "fmt"
    "time"
)

type Person struct {
    Name string
    Age  int
}

// simulateWork simule une tâche qui prend du temps (par ex. traitement de données)
func simulateWork(p Person, ch chan string) {
    time.Sleep(time.Millisecond * 100) // Simule un traitement
    msg := fmt.Sprintf("✔️  %s traité (âge %d)", p.Name, p.Age)
    ch <- msg
}

func main() {
    // Slice de personnes
    people := []Person{
        {Name: "Alice", Age: 30},
        {Name: "Bob", Age: 25},
        {Name: "Charlie", Age: 35},
    }

    // Channel pour recevoir les résultats
    results := make(chan string, len(people)) // buffered pour ne pas bloquer

    // Lancement des goroutines
    for _, p := range people {
        go simulateWork(p, results)
    }

    // Collecte des résultats
    for i := 0; i < len(people); i++ {
        fmt.Println(<-results) // Lecture bloquante tant que pas reçu
    }

    close(results) // bonne pratique même si pas indispensable ici
}
```

---

## 🧪 Exemple de sortie :

```
✔️  Bob traité (âge 25)
✔️  Alice traité (âge 30)
✔️  Charlie traité (âge 35)
```

*(L’ordre peut varier selon la planification des goroutines)*

---

## ✅ Éléments `builtin` utilisés ici :

| Élément      | Rôle dans le code                          |
| ------------ | ------------------------------------------ |
| `make`       | Créer un `chan string` avec capacité       |
| `len`        | Taille du slice pour contrôler les boucles |
| `chan`, `<-` | Communication entre goroutines             |
| `go`         | Lancer les fonctions concurrentes          |
| `close`      | Fermer le channel (bonne pratique)         |

---

--------------------------------------------------------------------------------------------------------

Souhaites-tu qu’on ajoute un `select` ou un `timeout` pour contrôler les goroutines plus finement ?

Parfait ! Ajoutons maintenant un **`select` avec un `timeout`** pour illustrer comment :

* gérer plusieurs `channels` (ou événements),
* éviter de bloquer indéfiniment si une goroutine est trop lente ou plante.

---

## 🎯 Objectif modifié :

* Attendre un résultat depuis le channel,
* Mais **passer à autre chose si ça prend trop de temps** (timeout).

---

## 🧱 Code avec `select` + `timeout`

```go
package main

import (
    "fmt"
    "time"
)

type Person struct {
    Name string
    Age  int
}

func simulateWork(p Person, ch chan string) {
    delay := time.Duration(100+int(p.Age)%3*300) * time.Millisecond // randomisé
    time.Sleep(delay)
    ch <- fmt.Sprintf("✔️  %s traité en %v", p.Name, delay)
}

func main() {
    people := []Person{
        {Name: "Alice", Age: 30},
        {Name: "Bob", Age: 25},
        {Name: "Charlie", Age: 35},
    }

    results := make(chan string)

    for _, p := range people {
        go simulateWork(p, results)

        select {
        case msg := <-results:
            fmt.Println(msg)
        case <-time.After(200 * time.Millisecond): // timeout
            fmt.Printf("⚠️  %s a mis trop de temps, on passe !\n", p.Name)
        }
    }
}
```

---

## ✅ Fonctionnement

* On lance une goroutine pour traiter chaque personne.
* On utilise `select` pour attendre soit :

  * une réponse dans `results`,
  * soit un **timeout** après `200ms`.
* Si une personne est lente, on affiche un avertissement.

---

## 🧪 Exemple de sortie :

```
✔️  Alice traité en 100ms
⚠️  Bob a mis trop de temps, on passe !
✔️  Charlie traité en 100ms
```

*(Les délais sont simulés, donc peuvent varier)*

---

## 🧠 Ce que tu as appris ici :

| Élément      | Utilité                                        |
| ------------ | ---------------------------------------------- |
| `select`     | Attendre sur plusieurs channels                |
| `time.After` | Créer un canal qui se déclenche après un délai |
| `chan`       | Transfert de message                           |
| `go`         | Exécution parallèle                            |

---

Souhaites-tu aller plus loin avec un **worker pool** (plusieurs workers traitant une file de tâches concurrentes avec un canal) ?
