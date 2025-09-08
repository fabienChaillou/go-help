En Go (Golang), une **fonction anonyme** est une fonction **sans nom** que tu peux définir **directement dans une expression**. Ces fonctions sont souvent utilisées pour :

* Les affecter à une variable
* Les passer en argument à une autre fonction
* Les exécuter immédiatement (IIFE : Immediately Invoked Function Expression)

---

### 1. **Déclaration et affectation à une variable**

```go
package main

import "fmt"

func main() {
    add := func(a int, b int) int {
        return a + b
    }

    fmt.Println(add(3, 4)) // Affiche 7
}
```

---

### 2. **Fonction anonyme exécutée immédiatement (IIFE)**

```go
package main

import "fmt"

func main() {
    result := func(a int, b int) int {
        return a * b
    }(3, 5)

    fmt.Println(result) // Affiche 15
}
```

---

### 3. **Passée comme argument**

```go
package main

import "fmt"

func operate(a, b int, f func(int, int) int) int {
    return f(a, b)
}

func main() {
    result := operate(10, 2, func(x, y int) int {
        return x - y
    })

    fmt.Println(result) // Affiche 8
}
```

---

### 4. **Utilité des fonctions anonymes**

* **Encapsulation temporaire** de logique
* **Callbacks**, surtout pour les routines ou événements
* **Fermetures (closures)** : elles peuvent capturer les variables de l’environnement

---

### 5. **Exemple de closure**

```go
package main

import "fmt"

func main() {
    counter := 0

    increment := func() int {
        counter++
        return counter
    }

    fmt.Println(increment()) // 1
    fmt.Println(increment()) // 2
}
```

Ici, la fonction anonyme **capture** la variable `counter` même après la fin de sa portée locale.

---
