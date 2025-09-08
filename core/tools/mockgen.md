Voici un exemple d’utilisation de `mockgen` de [GoMock](https://github.com/golang/mock) avec la directive `//go:generate` pour générer des mocks automatiquement dans un projet Go.

---

### 📦 Exemple de structure de projet

```bash
myapp/
├── go.mod
├── main.go
├── service/
│   ├── user.go          # L’interface à mocker
│   ├── user_test.go     # Les tests utilisant le mock
│   └── mock_user.go     # Le fichier généré
```

---

### 📄 `service/user.go` — interface à mocker

```go
package service

type UserService interface {
    GetUser(id int) (string, error)
}
```

---

### 📄 `service/mock_user.go` — fichier mock généré (ne pas écrire à la main)

Ajoute une directive `go:generate` dans `user.go` (ou un fichier dédié) :

```go
//go:generate mockgen -source=user.go -destination=mock_user.go -package=service
```

---

### 📄 `service/user_test.go` — test avec le mock

```go
package service

import (
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
)

func TestGetUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockUserService := NewMockUserService(ctrl)
    mockUserService.EXPECT().GetUser(1).Return("Alice", nil)

    name, err := mockUserService.GetUser(1)
    assert.NoError(t, err)
    assert.Equal(t, "Alice", name)
}
```

---

### 🔧 Commandes à exécuter

1. Installer `mockgen` :

```bash
go install go.uber.org/mock/mockgen@latest
```

2. Générer les mocks :

```bash
go generate ./...
```

---

## example avec mockgen et go generate

Parfait ! Voici un exemple de projet MVC en Go qui utilise `mockgen` avec `go:generate` pour générer des mocks dans un test unitaire. Ce projet simule une API simple de gestion d'utilisateurs avec les couches `controller`, `service`, et `repository`.

---

## 📦 Structure du projet

```
myapp/
├── go.mod
├── main.go
├── controller/
│   └── user_controller.go
├── service/
│   ├── user_service.go          # Contient l’interface UserService
│   ├── mock_user_service.go     # Fichier mock (généré)
│   └── user_service_test.go     # Tests unitaires
├── repository/
│   └── user_repository.go
```

---

### ✅ Étape 1 : Interface à mocker (`UserService`)

**`service/user_service.go`**

```go
package service

type UserService interface {
    GetUserName(id int) (string, error)
}

//go:generate mockgen -source=user_service.go -destination=mock_user_service.go -package=service
```

---

### ✅ Étape 2 : Implémentation (exemple simple)

**`repository/user_repository.go`**

```go
package repository

var users = map[int]string{
    1: "Alice",
    2: "Bob",
}

func GetUserByID(id int) (string, error) {
    name, ok := users[id]
    if !ok {
        return "", fmt.Errorf("user not found")
    }
    return name, nil
}
```

**`service/user_service_impl.go`**

```go
package service

import "myapp/repository"

type userServiceImpl struct{}

func NewUserService() UserService {
    return &userServiceImpl{}
}

func (s *userServiceImpl) GetUserName(id int) (string, error) {
    return repository.GetUserByID(id)
}
```

---

### ✅ Étape 3 : Contrôleur

**`controller/user_controller.go`**

```go
package controller

import (
    "fmt"
    "myapp/service"
)

type UserController struct {
    Service service.UserService
}

func (uc *UserController) ShowUserName(id int) {
    name, err := uc.Service.GetUserName(id)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Println("User:", name)
}
```

---

### ✅ Étape 4 : Test avec mock

**`service/user_service_test.go`**

```go
package service

import (
    "testing"
    "github.com/golang/mock/gomock"
    "github.com/stretchr/testify/assert"
)

func TestUserService(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mock := NewMockUserService(ctrl)
    mock.EXPECT().GetUserName(1).Return("TestUser", nil)

    name, err := mock.GetUserName(1)
    assert.NoError(t, err)
    assert.Equal(t, "TestUser", name)
}
```

---

### ✅ Étape 5 : Utilisation dans `main.go`

```go
package main

import (
    "myapp/controller"
    "myapp/service"
)

func main() {
    userController := controller.UserController{
        Service: service.NewUserService(),
    }

    userController.ShowUserName(1)
}
```

---

### 🛠 Commandes à exécuter

1. Installer `mockgen` :

```bash
go install go.uber.org/mock/mockgen@latest
```

2. Générer les mocks :

```bash
go generate ./...
```

3. Lancer les tests :

```bash
go test ./...
```

---
