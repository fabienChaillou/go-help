En Go, il n'existe pas de fonction intégrée universelle pour vérifier si une valeur existe, car cela dépend du **type de donnée** que tu utilises (slice, map, etc.). Voici quelques cas fréquents :

---

### ✅ 1. **Dans une map (recommandé)**

```go
m := map[string]int{"a": 1, "b": 2}

if val, exists := m["a"]; exists {
    fmt.Println("Existe avec la valeur :", val)
} else {
    fmt.Println("N'existe pas")
}
```

---

### ✅ 2. **Dans un slice**

Il faut parcourir manuellement le slice :

```go
arr := []string{"a", "b", "c"}
valToFind := "b"
found := false

for _, v := range arr {
    if v == valToFind {
        found = true
        break
    }
}
fmt.Println("Existe :", found)
```

Tu peux aussi en faire une fonction générique si tu utilises Go 1.18+ (avec les **génériques**) :

```go
func Contains[T comparable](slice []T, val T) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}
```

Utilisation :

```go
found := Contains([]int{1, 2, 3}, 2) // true
```

---
