En Go, **les closures** (ou **fermetures**) sont des fonctions **définies à l’intérieur d’une autre fonction** et qui **capturent les variables de leur environnement**. Cela signifie qu'elles peuvent utiliser et **se souvenir des variables locales** de la fonction dans laquelle elles ont été créées, même après que cette fonction ait terminé son exécution.

---

### 🧠 Concept de base

Quand une fonction interne **référence** une variable déclarée dans la fonction externe, elle la **capture**. Cette variable reste **vivante** tant que la closure est utilisée, même si la fonction parente a fini de s’exécuter.

---

### 📌 Exemple simple

```go
package main

import "fmt"

func adder() func(int) int {
    sum := 0
    return func(x int) int {
        sum += x
        return sum
    }
}

func main() {
    f := adder()
    fmt.Println(f(10)) // 10
    fmt.Println(f(5))  // 15
    fmt.Println(f(3))  // 18
}
```

#### 🔍 Explication :

* La fonction `adder` retourne une **closure**.
* `sum` est une variable locale à `adder`, mais la fonction retournée continue d'y accéder.
* À chaque appel de `f`, la variable `sum` est mise à jour.

---

### ✅ Pourquoi utiliser des closures ?

1. **Encapsulation** de l’état sans utiliser de structures.
2. Créer des **fonctions génératrices** ou des **fonctions configurables**.
3. Simplifier du code en évitant des variables globales.

---

### ⚠️ Attention

* Les variables capturées **sont partagées** par toutes les closures créées dans la même fonction. Cela peut causer des effets de bord si mal géré.
* Cela peut parfois compliquer le **débogage** si on ne comprend pas bien la portée des variables capturées.

---
