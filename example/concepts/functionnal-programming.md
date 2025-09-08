La **programmation fonctionnelle** est un **paradigme de programmation** qui met l'accent sur l’utilisation de **fonctions pures**, l’**immuabilité**, et l’**absence d’effets de bord**.

Bien que Go ne soit **pas un langage fonctionnel** à la base (il est impérativement structuré), **il permet d’adopter plusieurs concepts fonctionnels** grâce à son support des fonctions de premier ordre, des closures et de l’anonymat fonctionnel.

---

## 🔹 Concepts fonctionnels en Go

### 1. **Fonctions de premier ordre**

> Les fonctions sont des **valeurs** : elles peuvent être stockées dans des variables, passées en argument, ou retournées.

```go
func doubler(x int) int {
    return x * 2
}

func appliquer(f func(int) int, v int) int {
    return f(v)
}

func main() {
    fmt.Println(appliquer(doubler, 5)) // 10
}
```

---

### 2. **Fonctions anonymes et closures**

```go
carre := func(x int) int {
    return x * x
}
fmt.Println(carre(4)) // 16
```

Les **closures** permettent de capturer un état tout en restant dans un style fonctionnel.

---

### 3. **Fonctions retournant des fonctions**

```go
func multiplier(facteur int) func(int) int {
    return func(x int) int {
        return x * facteur
    }
}

fois3 := multiplier(3)
fmt.Println(fois3(4)) // 12
```

---

### 4. **Immuabilité**

> Go n’impose pas l’immuabilité, mais tu peux l'appliquer par discipline.

Exemple :

```go
func ajouter(x int, liste []int) []int {
    nouvelleListe := append([]int{}, liste...) // copie
    nouvelleListe = append(nouvelleListe, x)
    return nouvelleListe
}
```

Ici, on **ne modifie pas l'entrée** `liste`, on crée une **nouvelle version**.

---

### 5. **Pas d’effets de bord**

Une fonction pure :

* Ne modifie pas l’état externe.
* Ne dépend pas de l’état externe.
* Retourne toujours le même résultat pour les mêmes entrées.

```go
func somme(a, b int) int {
    return a + b // fonction pure
}
```

---

## 🔹 Limites de la programmation fonctionnelle en Go

* Pas de **map/filter/reduce** intégrés comme en Python ou JavaScript (mais on peut les écrire).
* Pas de **currying automatique**, **pattern matching**, ou **monades**.
* Go préfère la simplicité et la lisibilité sur les abstractions complexes.

---

## 🔹 Exemple : `filter` fonctionnel

```go
func filter(liste []int, predicat func(int) bool) []int {
    resultat := []int{}
    for _, v := range liste {
        if predicat(v) {
            resultat = append(resultat, v)
        }
    }
    return resultat
}

func main() {
    nombres := []int{1, 2, 3, 4, 5}
    pairs := filter(nombres, func(x int) bool { return x%2 == 0 })
    fmt.Println(pairs) // [2 4]
}
```

---

## 🔹 En résumé

| Concept                | Supporté en Go ?    | Exemple                      |
| ---------------------- | ------------------- | ---------------------------- |
| Fonctions de 1er ordre | ✅                   | `appliquer(f, x)`            |
| Fonctions anonymes     | ✅                   | `func(x int) int { ... }`    |
| Closures               | ✅                   | `capturer une variable`      |
| Immuabilité            | 🟡 (manuelle)       | `append([]T{}, original...)` |
| Effets de bord évités  | 🟡 (par discipline) | fonctions pures              |
| map/filter/reduce      | 🟡 (implémentables) | `filter(slice, func)`        |

---

Souhaites-tu une mini-lib ou un exemple complet de programmation fonctionnelle simulée (avec `map`, `reduce`, etc.) en Go ?
