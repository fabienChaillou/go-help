Le package `reflect` en Go permet d'inspecter et de manipuler dynamiquement des variables, leurs types, et leurs valeurs à l’exécution. C’est ce qu’on appelle **la réflexion**.

Voici les concepts de base à connaître sur `reflect` :

---

### 🔧 1. Importation

```go
import "reflect"
```

---

### 🔍 2. Obtenir des informations sur une variable

```go
var x int = 42

t := reflect.TypeOf(x)   // donne le type (Type)
v := reflect.ValueOf(x)  // donne la valeur (Value)

fmt.Println(t.Name()) // "int"
fmt.Println(v.Int())  // 42
```

---

### 🧱 3. Manipuler dynamiquement une structure

```go
type User struct {
	Name string
	Age  int
}

u := User{"Alice", 30}
v := reflect.ValueOf(u)

for i := 0; i < v.NumField(); i++ {
	fmt.Printf("Champ %d: %v\n", i, v.Field(i))
}
```

> ⚠️ Si vous voulez **modifier** des valeurs, vous devez passer un **pointeur** à `reflect.ValueOf` **et** utiliser `.Elem()` :

```go
v := reflect.ValueOf(&u).Elem()
v.FieldByName("Age").SetInt(35)
```

---

### 🔁 4. Appeler dynamiquement une méthode

```go
type Greeter struct{}

func (g Greeter) Greet(name string) {
	fmt.Println("Hello,", name)
}

g := Greeter{}
v := reflect.ValueOf(g)
method := v.MethodByName("Greet")
args := []reflect.Value{reflect.ValueOf("Bob")}
method.Call(args) // imprime: Hello, Bob
```

---

### ⚠️ Attention

* `reflect` est **lent** par rapport au code statique : évite de l'utiliser dans des chemins critiques de performance.
* Le code `reflect` est **plus difficile à lire, à maintenir et à tester**.
* Mais c’est très utile pour :

  * Généricité (avant Go 1.18)
  * Sérialisation/désérialisation (ex : `encoding/json`)
  * Validation automatique (comme avec `go-playground/validator`)

---

Souhaites-tu un exemple pratique d’utilisation, par exemple pour parcourir une struct et valider ses champs ?
