La **récursivité** en Go fonctionne comme dans la plupart des langages : une fonction s'appelle elle-même pour **résoudre un problème en le décomposant**.

---

## 🔹 Définition simple

> Une **fonction récursive** est une fonction qui s’appelle elle-même avec un ou plusieurs arguments modifiés, jusqu’à atteindre une **condition d’arrêt**.

---

## 🔹 Exemple classique : Factorielle

```go
func factorielle(n int) int {
    if n == 0 {
        return 1 // condition d’arrêt
    }
    return n * factorielle(n-1) // appel récursif
}

func main() {
    fmt.Println(factorielle(5)) // Output: 120
}
```

🔸 La fonction `factorielle(5)` appelle `factorielle(4)`, puis `factorielle(3)`...
🔸 Lorsqu’on atteint `n == 0`, la récursion **s’arrête**.

---

## 🔹 Autre exemple : Suite de Fibonacci

```go
func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return fibonacci(n-1) + fibonacci(n-2)
}
```

⚠️ **Peu efficace** car recalcul répétitif → complexité exponentielle.
➡️ Solution : **mémoïsation** ou approche **itérative**.

---

## 🔹 Points importants

* **Toujours avoir une condition d’arrêt** (`if`) pour éviter une récursion infinie.
* Chaque appel crée une nouvelle **frame** dans la pile d’appel (`stack`).
* Trop de récursion → risque de **stack overflow** (dépassement de pile).
* Go n'a **pas d'optimisation de récursion terminale** (contrairement à certains langages fonctionnels).

---

## 🔹 Quand utiliser la récursivité ?

* Arbres (parcours, structure hiérarchique)
* Algorithmes diviser pour régner (merge sort, quick sort)
* Problèmes mathématiques définis récursivement

---

Souhaites-tu un exemple d'algorithme récursif plus avancé (comme un parcours d’arbre binaire ou un tri récursif) ?
