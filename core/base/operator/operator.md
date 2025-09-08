En Go, les opérateurs se divisent en plusieurs catégories selon leur rôle. Voici une présentation complète et bien structurée avec des exemples.

---

## **1. Opérateurs arithmétiques**

Ces opérateurs servent à effectuer des calculs mathématiques.

| Opérateur | Signification  | Exemple | Résultat                          |
| --------- | -------------- | ------- | --------------------------------- |
| `+`       | Addition       | `3 + 2` | `5`                               |
| `-`       | Soustraction   | `3 - 2` | `1`                               |
| `*`       | Multiplication | `3 * 2` | `6`                               |
| `/`       | Division       | `5 / 2` | `2` *(division entière pour int)* |
| `%`       | Modulo (reste) | `5 % 2` | `1`                               |

### Exemple :

```go
a := 10
b := 3
fmt.Println(a + b) // 13
fmt.Println(a - b) // 7
fmt.Println(a * b) // 30
fmt.Println(a / b) // 3
fmt.Println(a % b) // 1
```

---

## **2. Opérateurs de comparaison**

Utilisés dans les conditions (`if`, `for`, etc.).
Ils retournent toujours un **booléen** (`true` ou `false`).

| Opérateur | Signification       | Exemple  | Résultat |
| --------- | ------------------- | -------- | -------- |
| `==`      | Égal à              | `3 == 3` | `true`   |
| `!=`      | Différent de        | `3 != 4` | `true`   |
| `<`       | Inférieur à         | `2 < 3`  | `true`   |
| `>`       | Supérieur à         | `5 > 2`  | `true`   |
| `<=`      | Inférieur ou égal à | `3 <= 3` | `true`   |
| `>=`      | Supérieur ou égal à | `3 >= 4` | `false`  |

### Exemple :

```go
x := 5
y := 8
fmt.Println(x == y) // false
fmt.Println(x != y) // true
fmt.Println(x < y)  // true
```

---

## **3. Opérateurs logiques (booléens)**

Utilisés pour combiner des conditions.

| Opérateur | Signification     | Exemple           | Résultat |
| --------- | ----------------- | ----------------- | -------- |
| `&&`      | ET logique (AND)  | `true && false`   | `false`  |
| `\|\|`    | OU logique (OR)   | `true \|\| false` | `true`   |
| `!`       | NON logique (NOT) | `!true`           | `false`  |

### Exemple :

```go
a := true
b := false

fmt.Println(a && b) // false
fmt.Println(a || b) // true
fmt.Println(!a)     // false
```

---

## **4. Opérateurs d’affectation**

Ils servent à attribuer ou mettre à jour des valeurs.

| Opérateur | Signification                 | Exemple  | Équivalent à |
| --------- | ----------------------------- | -------- | ------------ |
| `=`       | Affectation simple            | `x = 5`  | —            |
| `+=`      | Ajout et affectation          | `x += 2` | `x = x + 2`  |
| `-=`      | Soustraction et affectation   | `x -= 2` | `x = x - 2`  |
| `*=`      | Multiplication et affectation | `x *= 2` | `x = x * 2`  |
| `/=`      | Division et affectation       | `x /= 2` | `x = x / 2`  |
| `%=`      | Modulo et affectation         | `x %= 2` | `x = x % 2`  |

### Exemple :

```go
x := 10
x += 5
fmt.Println(x) // 15
```

---

## **5. Opérateurs de bits (bitwise)**

Ils manipulent les **bits** directement.

| Opérateur | Signification       | Exemple (a=6, b=3) | Résultat               |
| --------- | ------------------- | ------------------ | ---------------------- |
| `&`       | ET binaire (AND)    | `a & b`            | `2` (110 & 011 = 010)  |
| `\|`      | OU binaire (OR)     | `a \| b`           | `7` (110 \| 011 = 111) |
| `^`       | XOR binaire         | `a ^ b`            | `5` (110 ^ 011 = 101)  |
| `&^`      | AND NOT (bit clear) | `a &^ b`           | `4` (110 &^ 011 = 100) |
| `<<`      | Décalage à gauche   | `1 << 2`           | `4` (0001 → 0100)      |
| `>>`      | Décalage à droite   | `8 >> 2`           | `2` (1000 → 0010)      |

### Exemple :

```go
a := 6 // 110
b := 3 // 011

fmt.Println(a & b)  // 2
fmt.Println(a | b)  // 7
fmt.Println(a ^ b)  // 5
fmt.Println(a &^ b) // 4
fmt.Println(1 << 3) // 8
fmt.Println(8 >> 2) // 2
```

---

## **6. Opérateurs divers**

Quelques opérateurs spéciaux :

| Opérateur | Signification                   | Exemple                           |
| --------- | ------------------------------- | --------------------------------- |
| `*`       | Déférencement (pointeur)        | `*ptr`                            |
| `&`       | Adresse d’une variable          | `&x`                              |
| `<-`      | Opérateur channel (Go routines) | `ch <- value` et `value := <- ch` |

### Exemple avec pointeurs :

```go
x := 10
p := &x        // adresse de x
fmt.Println(*p) // 10
```

### Exemple avec channels :

```go
ch := make(chan int)
go func() {
    ch <- 42
}()
value := <-ch
fmt.Println(value) // 42
```

---

## **Résumé visuel**

| Catégorie         | Exemples           |
| ----------------- | ------------------ |
| **Arithmétiques** | `+ - * / %`        |
| **Comparaison**   | `== != < > <= >=`  |
| **Logiques**      | `&& \|\| !`        |
| **Affectation**   | `= += -= *= /= %=` |
| **Bits**          | `& \| ^ &^ << >>`  |
| **Divers**        | `* & <-`           |

---

## **En résumé**

* `<<` et `>>` : décalages binaires (puissances de deux).
* `&` et `|` : manipulation de bits pour flags et masques.
* `&^` : effacer des bits spécifiques.
* `*` et `&` : gestion des pointeurs.
* `<-` : communication entre goroutines via des channels.

Go possède un ensemble d'opérateurs clair et concis, bien adaptés à la programmation système et concurrente.
