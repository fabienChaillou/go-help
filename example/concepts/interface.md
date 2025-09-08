Excellente question ! Les **interfaces** sont un concept fondamental et puissant en Go, souvent utilisé pour écrire du code **modulaire**, **testable**, et **orienté comportement** plutôt que structure.

---

## 🧠 Qu’est-ce qu’une interface en Go ?

Une **interface** est un **type abstrait** qui décrit un **comportement** (un ensemble de méthodes).
**Toute valeur qui implémente les méthodes de l’interface satisfait cette interface.**

> 🟢 En Go, **il n’y a pas besoin d’implémentation explicite** : c’est implicite.

---

## 🧾 Exemple simple

```go
type Greeter interface {
    Greet() string
}

type French struct{}
func (f French) Greet() string {
    return "Bonjour"
}

type English struct{}
func (e English) Greet() string {
    return "Hello"
}

func saySomething(g Greeter) {
    fmt.Println(g.Greet())
}

func main() {
    saySomething(French{}) // Bonjour
    saySomething(English{}) // Hello
}
```

🔍 Ici :

* `Greeter` est une interface avec une méthode `Greet() string`.
* `French` et `English` satisfont automatiquement `Greeter`.
* `saySomething` accepte **n’importe quel type** qui implémente `Greet()`.

---

## ✅ Avantages des interfaces

| Avantage                 | Description                                                         |
| ------------------------ | ------------------------------------------------------------------- |
| 🔌 Découplage            | Le code utilise des comportements, pas des structures               |
| 🔄 Flexibilité           | On peut substituer différents comportements sans changer la logique |
| 🧪 Testabilité           | Permet de simuler des implémentations (mocks)                       |
| 🤝 Composition naturelle | Les interfaces peuvent être combinées facilement                    |

---

## 🧩 Exemple courant : `io.Writer`

```go
func writeMessage(w io.Writer, msg string) {
    w.Write([]byte(msg))
}
```

Tu peux passer :

* `os.Stdout`
* un `fichier`
* un `buffer mémoire` (`bytes.Buffer`)
* ou un mock

Tous implémentent la méthode `Write([]byte) (int, error)` → donc satisfont `io.Writer`.

---

## 🧱 Interface vide (`interface{}` ou `any`)

```go
func printAnything(x interface{}) {
    fmt.Println(x)
}
```

* `interface{}` peut accepter **n’importe quelle valeur**
* Utilisé dans les cas génériques ou quand le type exact n’est pas encore connu

---

## 🔄 Assertion et type switch

Pour utiliser une valeur `interface{}` comme son type original :

```go
var x interface{} = "hello"

s, ok := x.(string) // type assertion
if ok {
    fmt.Println("It's a string:", s)
}

switch v := x.(type) {
case int:
    fmt.Println("int", v)
case string:
    fmt.Println("string", v)
}
```

---

## 🧠 Résumé rapide

```go
type InterfaceName interface {
    Method1()
    Method2(param Type) ReturnType
}
```

* Implémentée automatiquement (pas besoin d’écrire "implements")
* Super utile pour les **injections**, **tests**, **stratégies**, **API flexibles**

---

Souhaites-tu un exemple avec **mock de test utilisant une interface** ou une interface **composée** (interface dans une interface) ?

---------------------------------------------------------------------------------------------------------------------------

Parfait ! Voyons **deux cas concrets** où les interfaces sont très utiles en Go :

---

## 📌 1. **Mock de test via interface**

Imaginons une application qui envoie des emails. On va :

* définir une interface `Mailer`
* créer une implémentation réelle
* puis une **fausse implémentation (mock)** pour les tests

---

### 📦 `mailer.go`

```go
package mailer

type Mailer interface {
    Send(to string, body string) error
}

type SmtpMailer struct{}

func (s SmtpMailer) Send(to string, body string) error {
    // envoi réel (exemple fictif)
    fmt.Printf("Sending to %s: %s\n", to, body)
    return nil
}
```

---

### 🧪 `mailer_test.go` (avec un mock)

```go
package mailer

import "testing"

// mockMailer enregistre les appels pour les tests
type mockMailer struct {
    sentTo  string
    sentMsg string
}

func (m *mockMailer) Send(to, body string) error {
    m.sentTo = to
    m.sentMsg = body
    return nil
}

func TestSend(t *testing.T) {
    mock := &mockMailer{}
    useMailer(mock) // utilise l'interface, pas l'implémentation réelle

    if mock.sentTo != "bob@example.com" || mock.sentMsg != "Hello" {
        t.Errorf("unexpected call: to=%s msg=%s", mock.sentTo, mock.sentMsg)
    }
}

// useMailer utilise l'interface Mailer
func useMailer(m Mailer) {
    m.Send("bob@example.com", "Hello")
}
```

✅ Ici, on teste la logique sans faire de vrai "Send" : **on isole et contrôle le comportement**.

---

## 🧱 2. **Interface composée (embedding)**

Tu peux **composer plusieurs interfaces** :

```go
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

// ReadWriter combine les deux
type ReadWriter interface {
    Reader
    Writer
}
```

🔁 Toute structure qui implémente `Read` **et** `Write` satisfait automatiquement `ReadWriter` (`os.File`, `bytes.Buffer`, etc.).

---

## 🧠 En résumé :

| Technique            | Utilité principale                             |
| -------------------- | ---------------------------------------------- |
| ✅ Mock par interface | Tester sans dépendre d’implémentations réelles |
| 🔁 Composition       | Regrouper des comportements réutilisables      |

---

Souhaites-tu que je t’aide à créer une **architecture modulaire** avec interfaces (ex. service-repository), ou voir une interface **avec des génériques (Go 1.18+)** ?

---------------------------------------------------------------------------------------------------------------------------

Excellent ! Tu vas maintenant voir comment créer une **architecture modulaire propre** en Go, en utilisant des **interfaces** pour séparer clairement :

* les **services métiers** (`UserService`)
* les **dépendances techniques** (`UserRepository`)
* avec des **mocks** pour les tests ✅

---

## 🧱 Objectif : Service + Repository pattern

On modélise un service `UserService` qui utilise un `UserRepository`. L’interface permet :

* de séparer la logique métier et la base de données
* de **remplacer le repository par un mock en test**

---

## 📦 1. `user.go` – Modèle

```go
package user

type User struct {
    ID   int
    Name string
}
```

---

## 📁 2. `repository.go` – Interface Repository

```go
package user

type UserRepository interface {
    GetByID(id int) (*User, error)
}
```

---

## 📁 3. `service.go` – Service qui dépend d’une interface

```go
package user

type UserService struct {
    Repo UserRepository
}

func (s *UserService) GetUsername(id int) (string, error) {
    user, err := s.Repo.GetByID(id)
    if err != nil {
        return "", err
    }
    return user.Name, nil
}
```

---

## 📁 4. `repo_in_memory.go` – Implémentation simple

```go
package user

import "errors"

type InMemoryRepo struct {
    users map[int]*User
}

func (r *InMemoryRepo) GetByID(id int) (*User, error) {
    u, ok := r.users[id]
    if !ok {
        return nil, errors.New("not found")
    }
    return u, nil
}
```

---

## 🧪 5. `service_test.go` – Test avec un mock

```go
package user

import (
    "testing"
)

type mockRepo struct{}

func (m *mockRepo) GetByID(id int) (*User, error) {
    return &User{ID: id, Name: "Mocked User"}, nil
}

func TestGetUsername(t *testing.T) {
    svc := UserService{Repo: &mockRepo{}}
    name, err := svc.GetUsername(123)

    if err != nil || name != "Mocked User" {
        t.Errorf("expected 'Mocked User', got '%s', err: %v", name, err)
    }
}
```

---

## ✅ Résultat

* `UserService` dépend **seulement d’un contrat** (`UserRepository`)
* en **prod** : tu peux brancher une vraie base (PostgreSQL, etc.)
* en **test** : tu utilises un mock en mémoire
* chaque partie est testable, réutilisable, maintenable

---

Souhaites-tu voir cette architecture connectée à une vraie base de données (`sqlx`, `gorm`) ou migrée avec **génériques** ?

---------------------------------------------------------------------------------------------------------------------------

Parfait ! Allons plus loin : je vais te montrer **comment brancher cette architecture sur une vraie base de données** en Go, avec `sqlx` (une surcouche pratique de `database/sql`).

---

## 🎯 Objectif

Brancher notre `UserService` à une vraie base PostgreSQL (ou SQLite) avec l’implémentation concrète du `UserRepository`.

---

## 🧰 Pré-requis : installer `sqlx`

```bash
go get github.com/jmoiron/sqlx
```

Et une base (ex : PostgreSQL ou SQLite via Docker/local)

---

## 🧩 Structure simplifiée

```
/user
  ├── model.go
  ├── repository.go         (interface)
  ├── service.go            (métier)
  ├── repo_sql.go           (implémentation réelle)
  └── service_test.go       (test avec mock)
```

---

## 📁 `model.go`

```go
package user

type User struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
}
```

---

## 📁 `repository.go` (interface)

```go
package user

type UserRepository interface {
    GetByID(id int) (*User, error)
}
```

---

## 📁 `repo_sql.go` (implémentation avec sqlx)

```go
package user

import (
    "github.com/jmoiron/sqlx"
)

type SQLRepo struct {
    DB *sqlx.DB
}

func (r *SQLRepo) GetByID(id int) (*User, error) {
    var user User
    err := r.DB.Get(&user, "SELECT id, name FROM users WHERE id = $1", id)
    if err != nil {
        return nil, err
    }
    return &user, nil
}
```

---

## 📁 `service.go`

```go
package user

type UserService struct {
    Repo UserRepository
}

func (s *UserService) GetUsername(id int) (string, error) {
    user, err := s.Repo.GetByID(id)
    if err != nil {
        return "", err
    }
    return user.Name, nil
}
```

---

## 🧪 Exemple d’utilisation dans `main.go`

```go
package main

import (
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq" // ou sqlite3
    "yourmodule/user"
)

func main() {
    db, err := sqlx.Connect("postgres", "user=postgres dbname=testdb sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

    repo := &user.SQLRepo{DB: db}
    svc := user.UserService{Repo: repo}

    name, err := svc.GetUsername(1)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Hello", name)
}
```

---

## 🧪 Pour tester localement avec SQLite :

Change :

```go
sqlx.Connect("sqlite3", "file:test.db")
```

Et adapte la requête SQL : `?` au lieu de `$1`.

---

## ✅ Avantages de cette architecture

| Aspect          | Avantage clé                              |
| --------------- | ----------------------------------------- |
| Interfaces      | Découplage métier / base de données       |
| sqlx            | Mapping structuré + ergonomie SQL         |
| Tests unitaires | Faciles à écrire avec des mocks           |
| Facile à migrer | Tu peux remplacer `sqlx` par `gorm`, etc. |

---

Souhaites-tu que je t’ajoute :

* une **méthode d’insertion avec `Create()`**,
* ou une version **générique** avec `Go 1.18+` pour éviter de dupliquer les repos ?


Super 👌 ! Passons maintenant à **une version générique de repository** avec Go 1.18+, qui évite de réécrire les mêmes fonctions pour chaque entité (`User`, `Product`, etc.).

---

## 🎯 Objectif

Créer un **repository générique** réutilisable pour n’importe quel modèle struct, avec `sqlx`.

---

## 🧩 Structure

On va créer :

* une interface `Repository[T any]`
* une implémentation `SQLRepository[T any]` qui fonctionne pour **tous les types**

---

## 📁 `repository.go` — interface générique

```go
package repository

type Repository[T any] interface {
    GetByID(id int) (*T, error)
    Create(entity *T) error
}
```

---

## 📁 `sql_repository.go` — implémentation générique avec `sqlx`

```go
package repository

import (
    "fmt"

    "github.com/jmoiron/sqlx"
)

type SQLRepository[T any] struct {
    DB    *sqlx.DB
    Table string
}

// Récupère une entité par son ID
func (r *SQLRepository[T]) GetByID(id int) (*T, error) {
    var t T
    query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.Table)
    err := r.DB.Get(&t, query, id)
    if err != nil {
        return nil, err
    }
    return &t, nil
}

// Insère une nouvelle entité
func (r *SQLRepository[T]) Create(entity *T) error {
    // ⚠️ Nécessite un insert préparé selon ta structure ou un outil comme squirrel/sqlbuilder
    return fmt.Errorf("Create not implemented: requires SQL builder or reflection")
}
```

> 💡 `Create()` ici est un **exemple d’extension**, mais pour de vraies insertions génériques, tu auras besoin d’un générateur SQL ou d'une lib de réflexion.

---

## 📁 `user/model.go`

```go
package user

type User struct {
    ID   int    `db:"id"`
    Name string `db:"name"`
}
```

---

## 📁 `main.go` — usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"

    "yourmodule/repository"
    "yourmodule/user"
)

func main() {
    db, err := sqlx.Connect("postgres", "user=postgres dbname=testdb sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

    userRepo := &repository.SQLRepository[user.User]{
        DB:    db,
        Table: "users",
    }

    user, err := userRepo.GetByID(1)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("User:", user.Name)
}
```

---

## ✅ Avantages

| Avantage                                 | Détail                           |
| ---------------------------------------- | -------------------------------- |
| 🔁 Réutilisable pour n’importe quel type | `Product`, `Order`, etc.         |
| 🧼 Évite la duplication                  | Un seul repo pour tous les types |
| 🔒 Fort typage                           | Chaque repo a un `T` bien défini |

---

Souhaites-tu que je t’ajoute une version complète de `Create()` générique avec la lib [`sqlbuilder`](https://github.com/huandu/go-sqlbuilder) ou un générateur SQL comme [`squirrel`](https://github.com/Masterminds/squirrel`) ?

Parfait 🙌 ! On va maintenant rendre **le `Create()` générique fonctionnel**, avec la bibliothèque [`squirrel`](https://github.com/Masterminds/squirrel), qui est idéale pour construire des requêtes SQL **dynamiques et sûres** en Go.

---

## 🧰 Prérequis : Installer squirrel

```bash
go get github.com/Masterminds/squirrel
```

---

## 🧱 Objectif

Dans le `SQLRepository[T]`, rendre `Create()` capable de générer une requête `INSERT INTO ...` pour **n'importe quelle struct**.

---

## ⚠️ Attention

Pour cela, on fait un **mapping struct -> champs SQL** à la volée. On utilisera `squirrel` **et `sqlx.NamedStruct()`** pour lier les paramètres proprement.

---

## 📁 `sql_repository.go` – version avec `Create()` générique

```go
package repository

import (
    "fmt"
    "reflect"

    sq "github.com/Masterminds/squirrel"
    "github.com/jmoiron/sqlx"
)

type SQLRepository[T any] struct {
    DB    *sqlx.DB
    Table string
}

// GetByID reste identique
func (r *SQLRepository[T]) GetByID(id int) (*T, error) {
    var t T
    query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.Table)
    err := r.DB.Get(&t, query, id)
    if err != nil {
        return nil, err
    }
    return &t, nil
}

// Create utilise squirrel pour générer l’INSERT dynamiquement
func (r *SQLRepository[T]) Create(entity *T) error {
    // Reflection pour extraire les champs de la struct
    val := reflect.Indirect(reflect.ValueOf(entity))
    typ := val.Type()

    fields := make([]string, 0)
    values := make(map[string]interface{})

    for i := 0; i < typ.NumField(); i++ {
        field := typ.Field(i)
        dbTag := field.Tag.Get("db")
        if dbTag == "" || dbTag == "-" || dbTag == "id" {
            continue // ignore les champs sans tag ou l'ID (si auto-incrémenté)
        }

        fields = append(fields, dbTag)
        values[dbTag] = val.Field(i).Interface()
    }

    // Génère l'insert via squirrel
    query, args, err := sq.
        Insert(r.Table).
        SetMap(values).
        PlaceholderFormat(sq.Dollar).
        ToSql()
    if err != nil {
        return err
    }

    _, err = r.DB.Exec(query, args...)
    return err
}
```

---

## 📁 Exemple de struct : `user/model.go`

```go
package user

type User struct {
    ID   int    `db:"id"`   // ignoré pour insert
    Name string `db:"name"` // utilisé
}
```

---

## 📁 `main.go` – utilisation

```go
package main

import (
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"

    "yourmodule/repository"
    "yourmodule/user"
)

func main() {
    db, err := sqlx.Connect("postgres", "user=postgres dbname=testdb sslmode=disable")
    if err != nil {
        log.Fatalln(err)
    }

    userRepo := &repository.SQLRepository[user.User]{
        DB:    db,
        Table: "users",
    }

    // Création d’un utilisateur
    newUser := &user.User{Name: "Jean Dupont"}
    err = userRepo.Create(newUser)
    if err != nil {
        log.Fatal("Erreur Create:", err)
    }

    fmt.Println("✅ Utilisateur ajouté avec succès !")
}
```

---

## ✅ Ce que tu gagnes

| Atout           | Détail                                                   |
| --------------- | -------------------------------------------------------- |
| 🔄 Réutilisable | Tu peux utiliser ce repo pour `Product`, `Article`, etc. |
| 🧼 Lisible      | Pas de SQL en dur : squirrel génère proprement           |
| 🔒 Sûr          | Pas d’injection SQL, noms + valeurs séparés              |
| ⚡️ Performant   | Compatible PostgreSQL, MySQL, SQLite                     |

---

Souhaites-tu que je t’aide à :

* gérer **`RETURNING id`** pour récupérer l’ID après insert (Postgres) ?
* ajouter une **méthode `Update()` générique** avec `squirrel` ?
* ou convertir tout ça en **module complet prêt à déployer** ?
