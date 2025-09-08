En Go, **le struct embedding** (ou "inclusion de structure") est une fonctionnalité qui permet à une structure (`struct`) d’en inclure une autre, ce qui permet de **composer** des comportements ou des données sans recourir à l’héritage (Go n’a pas d’héritage classique comme dans les langages orientés objet).

### Concept de base

Quand tu **embed** une structure dans une autre, tu peux accéder directement aux champs et méthodes de la structure embarquée **comme s’ils faisaient partie de la structure principale**.

### Exemple simple

```go
type Person struct {
    Name string
    Age  int
}

type Employee struct {
    Person // struct embedded ici
    Company string
}
```

Tu peux maintenant faire :

```go
func main() {
    e := Employee{
        Person: Person{
            Name: "Alice",
            Age:  30,
        },
        Company: "Acme Inc.",
    }

    fmt.Println(e.Name)   // accès direct au champ Name de Person
    fmt.Println(e.Age)    // idem
    fmt.Println(e.Company)
}
```

### Détails importants

* L’**embedding** n’est pas de l’héritage, mais plutôt une **composition**.
* Tu peux aussi **embed** plusieurs structs.
* Si deux structs embarquées ont des champs/méthodes avec le même nom, il y a un conflit, et tu dois spécifier explicitement à laquelle tu fais référence.

### Exemple avec méthode

```go
type Logger struct{}

func (l Logger) Log(msg string) {
    fmt.Println("LOG:", msg)
}

type Service struct {
    Logger // embed du logger
}

func main() {
    s := Service{}
    s.Log("Service started") // Appelle Logger.Log
}
```

### Pourquoi utiliser struct embedding ?

* **Réutilisation** de code (par exemple des fonctions utilitaires).
* Pour **simuler une forme d’héritage** comportemental.
* Très utile pour les **patterns de design** comme les middlewares ou les décorateurs.
* Utilisé dans la **standard library** (ex : `http.Server` embed `http.Handler` parfois).
