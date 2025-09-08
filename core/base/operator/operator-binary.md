En Go, l’expression `1 << 2` utilise l’opérateur **de décalage binaire vers la gauche** (`<<`).

### **Explication pas à pas**

* `1` en binaire :

  ```
  0001
  ```
* `<< 2` signifie *décaler les bits de deux positions vers la gauche*.

Après le décalage :

```
0001 << 2  →  0100
```

En décimal, `0100` correspond à **4**.

Donc :

```go
fmt.Println(1 << 2) // Affiche 4
```

---

### **Formule générale**

Le décalage à gauche de `n` positions revient à **multiplier par 2ⁿ** :

$$
1 << n = 1 \times 2^n
$$

Exemples :

| Expression | Résultat binaire | Résultat décimal |
| ---------- | ---------------- | ---------------- |
| `1 << 0`   | `0001`           | 1                |
| `1 << 1`   | `0010`           | 2                |
| `1 << 2`   | `0100`           | 4                |
| `1 << 3`   | `1000`           | 8                |

---

### **Utilisation courante**

Cet opérateur est souvent utilisé pour :

1. **Manipuler des flags ou des masques binaires** :

   ```go
   const (
       Read  = 1 << 0 // 0001 → 1
       Write = 1 << 1 // 0010 → 2
       Exec  = 1 << 2 // 0100 → 4
   )

   perms := Read | Exec // Combine Read et Exec → 0101
   fmt.Println(perms)   // 5
   ```

2. **Calculer des puissances de deux rapidement** :

   ```go
   for i := 0; i < 5; i++ {
       fmt.Println(1 << i)
   }
   // Affiche : 1, 2, 4, 8, 16
   ```

---

En résumé :

* `<<` décale les bits vers la gauche.
* `1 << 2` équivaut à $1 \times 2^2 = 4$.
* C'est très utile pour gérer des permissions, flags, et puissances de deux.
