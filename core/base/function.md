En Go (Golang), une **fonction** est un bloc de code que l’on peut appeler pour effectuer une tâche spécifique. Les fonctions sont essentielles en Go, car elles permettent de structurer, organiser et réutiliser du code.

---

### 🧩 Syntaxe d'une fonction

```go
func nomDeLaFonction(param1 type1, param2 type2) retourType {
    // corps de la fonction
    return valeur
}
```

### 🔹 Exemple simple

```go
func ajouter(a int, b int) int {
    return a + b
}
```

Ici, `ajouter` prend deux entiers `a` et `b`, et retourne leur somme (`int`).

---

### ✅ Appel d'une fonction

```go
resultat := ajouter(3, 5)
fmt.Println(resultat) // Affiche 8
```

---

### 🔄 Retour multiple

Go permet de **retourner plusieurs valeurs** :

```go
func diviser(a, b int) (int, int) {
    return a / b, a % b
}
```

```go
quotient, reste := diviser(10, 3)
fmt.Println(quotient, reste) // Affiche 3 1
```

---

### 🧪 Paramètres nommés et types groupés

Tu peux regrouper les paramètres qui ont le même type :

```go
func multiplier(a, b int) int {
    return a * b
}
```

---

### 🔁 Fonctions anonymes

Tu peux déclarer des fonctions sans nom (utiles pour des callbacks, goroutines, etc.) :

```go
f := func(x int) int {
    return x * x
}
fmt.Println(f(4)) // Affiche 16
```

---

### 🔧 Fonctions comme arguments

Les fonctions peuvent être passées en paramètre :

```go
func appliquer(f func(int) int, valeur int) int {
    return f(valeur)
}
```

---

### 🧠 Fonction récursive

Une fonction peut s’appeler elle-même :

```go
func factorielle(n int) int {
    if n <= 1 {
        return 1
    }
    return n * factorielle(n-1)
}
```

---
