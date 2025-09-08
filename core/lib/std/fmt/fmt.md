Le package `fmt` en Go (Golang) est l’un des packages standard les plus utilisés. Il fournit des fonctions de **formatage d'entrée/sortie** (I/O formatting), similaires à ce que fait `printf` en C. Le nom `fmt` est une abréviation de “format”.

### Fonctions principales du package `fmt`

Voici une vue d’ensemble des fonctions les plus courantes classées par usage :

---

### 🔹 **Affichage standard (console)**

* **`fmt.Print()`** : imprime les arguments dans la sortie standard sans retour à la ligne.

  ```go
  fmt.Print("Bonjour", " ", "le monde") // Bonjour le monde
  ```

* **`fmt.Println()`** : comme `Print` mais ajoute un retour à la ligne à la fin.

  ```go
  fmt.Println("Bonjour le monde") // Bonjour le monde\n
  ```

* **`fmt.Printf()`** : formatage avec des verbes de format, comme en C.

  ```go
  fmt.Printf("Nom: %s, Âge: %d\n", "Alice", 30)
  ```

---

### 🔹 **Formatage en chaîne (sans affichage)**

* **`fmt.Sprintf()`** : retourne une chaîne formatée sans l'afficher.

  ```go
  s := fmt.Sprintf("Pi vaut environ %.2f", 3.14159)
  fmt.Println(s) // Pi vaut environ 3.14
  ```

* **`fmt.Sprint()`**, **`fmt.Sprintln()`** : comme `Print()` et `Println()`, mais retourne une chaîne.

---

### 🔹 **Lecture (entrée standard ou chaînes)**

* **`fmt.Scan()`** : lit depuis l’entrée standard.

  ```go
  var name string
  fmt.Print("Votre nom ? ")
  fmt.Scan(&name)
  fmt.Println("Bonjour", name)
  ```

* **`fmt.Sscanf()`** : lit à partir d’une chaîne.

  ```go
  var a int
  fmt.Sscanf("123", "%d", &a)
  fmt.Println(a) // 123
  ```

---

### 🔹 **Verbes de format courants**

| Verbe | Signification                             |
| ----- | ----------------------------------------- |
| `%d`  | entier (base 10)                          |
| `%f`  | flottant (décimal)                        |
| `%s`  | chaîne de caractères                      |
| `%t`  | booléen                                   |
| `%v`  | valeur par défaut d'une variable          |
| `%+v` | affiche les champs nommés d'une struct    |
| `%#v` | affiche la représentation Go de la valeur |
| `%T`  | type de la variable                       |
| `%%`  | le caractère `%`                          |

---

### Exemples complets

```go
package main

import "fmt"

type Person struct {
    Name string
    Age  int
}

func main() {
    // Impression simple
    fmt.Println("Hello, world")

    // Formatage avec Printf
    fmt.Printf("Nombre: %d, Pi: %.2f\n", 42, 3.14159)

    // Struct avec %+v
    p := Person{"Alice", 30}
    fmt.Printf("Struct: %+v\n", p)

    // Lecture utilisateur
    var name string
    fmt.Print("Entrez votre nom : ")
    fmt.Scan(&name)
    fmt.Println("Bonjour", name)
}
```

---
