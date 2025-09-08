En Go (ou Golang), une **interface** est un type qui définit un ensemble de **méthodes** mais **ne fournit pas d’implémentation**. Une interface permet donc de **décrire un comportement** sans spécifier comment ce comportement est réalisé. Elle joue un rôle fondamental dans l'**abstraction**, le **polymorphisme**, et les **tests**.

---

### 💡 Définition simple

Une interface dit : "tout type qui implémente ces méthodes est de ce type d’interface".

```go
type Animal interface {
    Parler() string
}
```

Ici, tout type qui a une méthode `Parler() string` **implémente l’interface `Animal`**.

---

### ✅ Implémentation implicite

En Go, il **n’y a pas besoin de déclarer** qu’un type implémente une interface. C’est fait **implicitement** si le type a toutes les méthodes requises.

```go
type Chien struct{}

func (c Chien) Parler() string {
    return "Wouf"
}

type Chat struct{}

func (c Chat) Parler() string {
    return "Miaou"
}

func FaireParler(a Animal) {
    fmt.Println(a.Parler())
}

func main() {
    var c Chien
    var ch Chat
    FaireParler(c)   // Wouf
    FaireParler(ch)  // Miaou
}
```

---

### 🎯 Utilité des interfaces

1. **Abstraction** : on code sur des comportements, pas des implémentations.
2. **Découplage** : facilite la maintenance et les tests.
3. **Tests unitaires** : on peut facilement utiliser des mocks.
4. **Polymorphisme** : des objets différents peuvent être utilisés de la même manière.

---

### 🔄 Interface intégrée : `io.Reader`

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

Beaucoup de types standard (fichiers, réseau, buffer, etc.) implémentent cette interface.

---

### 🧪 Exemple avec interface personnalisée et test

```go
type Repository interface {
    GetUser(id int) (User, error)
}

// En test
type FakeRepo struct{}

func (f FakeRepo) GetUser(id int) (User, error) {
    return User{Name: "Fake"}, nil
}
```

---
