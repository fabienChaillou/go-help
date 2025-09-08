Voici un exemple complet d'utilisation de [`GoMock`](https://github.com/golang/mock) dans un projet Go. Ce projet met en place une petite application qui récupère des utilisateurs depuis une base de données simulée à l’aide de `gomock`.

---

## 📦 Structure du projet

```
gomock-example/
├── go.mod
├── main.go
├── models/
│   └── user.go
├── repository/
│   ├── user_repository.go
│   └── mock_user_repository.go  # généré par gomock
├── service/
│   └── user_service.go
└── service_test/
    └── user_service_test.go
```

---

## 🔧 Étapes de préparation

1. **Initialiser le module :**

```bash
go mod init gomock-example
```

2. **Installer `gomock` et `mockgen` :**

```bash
go install go.uber.org/mock/mockgen@latest
go get go.uber.org/mock/gomock
```

---

## 📄 Fichiers source

### `models/user.go`

```go
package models

type User struct {
	ID   int
	Name string
}
```

---

### `repository/user_repository.go`

```go
package repository

import "gomock-example/models"

type UserRepository interface {
	GetUserByID(id int) (*models.User, error)
}
```

---

### `service/user_service.go`

```go
package service

import (
	"gomock-example/models"
	"gomock-example/repository"
)

type UserService struct {
	Repo repository.UserRepository
}

func (s *UserService) GetUserName(id int) (string, error) {
	user, err := s.Repo.GetUserByID(id)
	if err != nil {
		return "", err
	}
	return user.Name, nil
}
```

---

## 🔧 Générer le mock

Exécute cette commande depuis la racine du projet pour générer le mock :

```bash
mockgen -source=repository/user_repository.go -destination=repository/mock_user_repository.go -package=repository
```

---

## 🧪 Test avec GoMock

### `service_test/user_service_test.go`

```go
package service_test

import (
	"errors"
	"gomock-example/models"
	"gomock-example/repository"
	"gomock-example/service"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	mockRepo.
		EXPECT().
		GetUserByID(1).
		Return(&models.User{ID: 1, Name: "Alice"}, nil)

	userService := service.UserService{Repo: mockRepo}

	name, err := userService.GetUserName(1)

	assert.NoError(t, err)
	assert.Equal(t, "Alice", name)
}

func TestGetUserName_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)

	mockRepo.
		EXPECT().
		GetUserByID(2).
		Return(nil, errors.New("user not found"))

	userService := service.UserService{Repo: mockRepo}

	name, err := userService.GetUserName(2)

	assert.Error(t, err)
	assert.Equal(t, "", name)
}
```

---

Tu peux lancer les tests avec :

```bash
go test ./...
```
