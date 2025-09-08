# Le concept de Reflect en Go (Golang)

La bibliothèque `reflect` est un package puissant du langage Go qui permet d'examiner et de manipuler des objets à l'exécution. C'est l'implémentation de la réflexion en Go, une fonctionnalité qui permet d'inspecter et de modifier des structures, des types et des valeurs pendant l'exécution du programme.

## Principes fondamentaux

En Go, le package `reflect` repose sur deux concepts principaux :

1. **Type** : représente le type d'une variable (via `reflect.Type`)
2. **Value** : représente la valeur d'une variable (via `reflect.Value`)

Ces deux éléments permettent d'examiner et de manipuler les structures de données dynamiquement.

## Utilisations courantes

### 1. Inspection des types à l'exécution

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    x := 3.14
    fmt.Println("Type:", reflect.TypeOf(x))
    fmt.Println("Value:", reflect.ValueOf(x))
    
    // Récupérer des informations supplémentaires sur le type
    v := reflect.ValueOf(x)
    fmt.Println("Type est float64?", v.Kind() == reflect.Float64)
    fmt.Println("Valeur:", v.Float())
}
```

### 2. Accès et modification de structures

```go
package main

import (
    "fmt"
    "reflect"
)

type Person struct {
    Name string
    Age  int
}

func main() {
    p := Person{"Alice", 30}
    
    // Obtenir la valeur de reflect
    v := reflect.ValueOf(&p).Elem()
    
    // Accéder aux champs
    nameField := v.FieldByName("Name")
    fmt.Println("Name:", nameField.String())
    
    // Modification d'un champ
    if nameField.CanSet() {
        nameField.SetString("Bob")
    }
    
    fmt.Println("Après modification:", p)
}
```

### 3. Parcourir les champs d'une structure

```go
package main

import (
    "fmt"
    "reflect"
)

type User struct {
    ID      int
    Name    string
    Email   string
    IsAdmin bool
}

func main() {
    user := User{1, "Jean", "jean@example.com", false}
    
    val := reflect.ValueOf(user)
    typ := val.Type()
    
    fmt.Println("Structure avec", val.NumField(), "champs:")
    
    for i := 0; i < val.NumField(); i++ {
        field := typ.Field(i)
        value := val.Field(i)
        fmt.Printf("%d: %s (%s) = %v\n", i, field.Name, field.Type, value.Interface())
    }
}
```

### 4. Appel dynamique de méthodes

```go
package main

import (
    "fmt"
    "reflect"
)

type Calculator struct{}

func (c *Calculator) Add(a, b int) int {
    return a + b
}

func main() {
    calc := &Calculator{}
    
    // Obtenir la méthode par réflexion
    method := reflect.ValueOf(calc).MethodByName("Add")
    
    // Préparer les arguments
    args := []reflect.Value{
        reflect.ValueOf(5),
        reflect.ValueOf(3),
    }
    
    // Appeler la méthode
    results := method.Call(args)
    
    // Afficher le résultat
    fmt.Println("Résultat:", results[0].Interface())
}
```

### 5. Création dynamique d'objets

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    // Créer un type de tableau dynamiquement
    sliceType := reflect.SliceOf(reflect.TypeOf(""))
    
    // Créer une instance de ce type
    slice := reflect.MakeSlice(sliceType, 0, 10)
    
    // Ajouter des éléments
    slice = reflect.Append(slice, reflect.ValueOf("hello"))
    slice = reflect.Append(slice, reflect.ValueOf("world"))
    
    // Convertir en type concret
    s := slice.Interface().([]string)
    fmt.Println(s) // [hello world]
}
```

## Avantages et inconvénients

### Avantages
- Permet d'écrire du code générique qui fonctionne avec différents types
- Facilite la création de frameworks et bibliothèques
- Utile pour la sérialisation/désérialisation (comme dans les packages `encoding/json`)
- Permet l'introspection des structures et des types

### Inconvénients
- Perte des vérifications de type à la compilation
- Code moins lisible et plus difficile à maintenir
- Performances réduites (2-3 fois plus lent que le code typé statiquement)
- Risque d'erreurs d'exécution (panics)

## Bonnes pratiques

1. Éviter `reflect` quand des solutions typées statiquement sont possibles
2. Documenter attentivement le code utilisant la réflexion
3. Utiliser des tests unitaires complets pour le code de réflexion
4. Capturer les panics potentielles avec des defer/recover
5. Ne pas exposer la réflexion dans les API publiques

La réflexion en Go est un outil puissant mais qui doit être utilisé avec parcimonie, pour des cas spécifiques où la généricité et l'inspection dynamique sont vraiment nécessaires.
