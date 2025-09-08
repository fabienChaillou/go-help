En Go, pour comparer des types (ou vérifier si deux variables sont du même type), tu peux utiliser la fonction `reflect.TypeOf` du package `reflect`.

### Exemple de fonction pour comparer deux types en Go

```go
import (
	"fmt"
	"reflect"
)

func CompareType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}

func main() {
	var x int = 10
	var y int64 = 10
	var z int = 20

	fmt.Println(CompareType(x, y)) // false
	fmt.Println(CompareType(x, z)) // true
}
```

### Explication :

* `reflect.TypeOf(variable)` retourne le type de la variable (`reflect.Type`).
* Tu peux comparer deux types directement avec `==`.

---

### Bonus : afficher les types pour déboguer

```go
fmt.Println(reflect.TypeOf(x)) // int
fmt.Println(reflect.TypeOf(y)) // int64
```

Merci pour la précision ! Voici différents cas courants de comparaison de types en Go, selon le contexte :

---

### 🔹 1. **Comparer deux variables de types quelconques**

```go
func SameType(a, b interface{}) bool {
	return reflect.TypeOf(a) == reflect.TypeOf(b)
}
```

---

### 🔹 2. **Vérifier si une variable est d’un type spécifique**

Exemple : vérifier si `v` est de type `string` ou `*MyStruct` :

```go
if reflect.TypeOf(v) == reflect.TypeOf("") {
	fmt.Println("v est une string")
}

if reflect.TypeOf(v) == reflect.TypeOf(&MyStruct{}) {
	fmt.Println("v est un pointeur vers MyStruct")
}
```

---

### 🔹 3. **Vérifier si une variable implémente une interface**

Par exemple :

```go
type MyInterface interface {
	DoSomething()
}

type MyType struct{}

func (m MyType) DoSomething() {}

func ImplementsInterface(val interface{}) bool {
	_, ok := val.(MyInterface)
	return ok
}
```

---

### 🔹 4. **Vérifier le type dynamique d’une interface vide (`interface{}`)**

Utilise un `switch type` :

```go
func PrintType(v interface{}) {
	switch val := v.(type) {
	case int:
		fmt.Println("int", val)
	case string:
		fmt.Println("string", val)
	case MyStruct:
		fmt.Println("MyStruct", val)
	default:
		fmt.Println("Unknown type", reflect.TypeOf(v))
	}
}
```

---
