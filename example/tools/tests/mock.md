Pour tester un composant en Go en utilisant un *mock*, tu peux définir une **interface** correspondant aux méthodes à tester, puis utiliser un outil comme `mockgen` (de `gomock`) ou le faire manuellement.

### Étapes pour ajouter une interface et faire un mock (manuellement ou avec gomock)

#### 1. **Définir une interface**

Supposons que tu as un service qui dépend d’un repository, tu peux extraire une interface :

```go
// user.go
type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	GetByID(id int) (*User, error)
	Add(user *User) error
}
```

#### 2. **Utiliser l’interface dans le service**

```go
// user_service.go
type UserService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(user *User) error {
	return s.repo.Add(user)
}
```

#### 3. **Créer un mock (manuellement ou avec gomock)**

##### Option A: Manuellement

```go
// mock_user_repository.go
type MockUserRepository struct {
	AddFunc    func(user *User) error
	GetByIDFunc func(id int) (*User, error)
}

func (m *MockUserRepository) Add(user *User) error {
	return m.AddFunc(user)
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	return m.GetByIDFunc(id)
}
```

##### Option B: Avec `gomock`

Installe `mockgen` :

```sh
go install go.uber.org/mock/mockgen@latest
```

Génère le mock :

```sh
mockgen -source=user.go -destination=mock_user_repository.go -package=yourpackage
```

#### 4. **Écrire un test**

```go
// user_service_test.go
func TestCreateUser(t *testing.T) {
	mockRepo := &MockUserRepository{
		AddFunc: func(user *User) error {
			if user.Name == "" {
				return errors.New("name required")
			}
			return nil
		},
	}

	service := NewUserService(mockRepo)

	err := service.CreateUser(&User{Name: "Alice"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}
```

Souhaites-tu un exemple complet utilisant `gomock`, ou préfères-tu rester avec les mocks faits à la main ?

Parfait ! Voici un exemple **complet et exécutable** en Go utilisant un mock fait à la main. On va :

* Définir une interface `UserRepository`
* Créer un service `UserService` qui utilise cette interface
* Écrire un mock `MockUserRepository`
* Écrire un test unitaire avec ce mock

---

### 📁 Structure

```
/mockexample/
├── user.go
├── user_service.go
├── mock_user_repository.go
└── user_service_test.go
```

---

### 🧩 `user.go`

```go
package mockexample

type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	Add(user *User) error
	GetByID(id int) (*User, error)
}
```

---

### ⚙️ `user_service.go`

```go
package mockexample

type UserService struct {
	repo UserRepository
}

func NewUserService(r UserRepository) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) CreateUser(user *User) error {
	if user.Name == "" {
		return ErrInvalidName
	}
	return s.repo.Add(user)
}

func (s *UserService) FindUser(id int) (*User, error) {
	return s.repo.GetByID(id)
}
```

---

### ⚠️ `errors.go` (ajoute si nécessaire)

```go
package mockexample

import "errors"

var ErrInvalidName = errors.New("invalid name")
```

---

### 🧪 `mock_user_repository.go`

```go
package mockexample

type MockUserRepository struct {
	AddFunc     func(user *User) error
	GetByIDFunc func(id int) (*User, error)
}

func (m *MockUserRepository) Add(user *User) error {
	return m.AddFunc(user)
}

func (m *MockUserRepository) GetByID(id int) (*User, error) {
	return m.GetByIDFunc(id)
}
```

---

### ✅ `user_service_test.go`

```go
package mockexample

import (
	"errors"
	"testing"
)

func TestCreateUser_Success(t *testing.T) {
	mockRepo := &MockUserRepository{
		AddFunc: func(user *User) error {
			if user.Name == "Alice" {
				return nil
			}
			return errors.New("unexpected user")
		},
	}

	service := NewUserService(mockRepo)

	err := service.CreateUser(&User{Name: "Alice"})
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestCreateUser_InvalidName(t *testing.T) {
	mockRepo := &MockUserRepository{
		AddFunc: func(user *User) error {
			return nil // won't be called
		},
	}

	service := NewUserService(mockRepo)

	err := service.CreateUser(&User{Name: ""})
	if err != ErrInvalidName {
		t.Errorf("expected ErrInvalidName, got %v", err)
	}
}

func TestFindUser(t *testing.T) {
	mockRepo := &MockUserRepository{
		GetByIDFunc: func(id int) (*User, error) {
			return &User{ID: id, Name: "Bob"}, nil
		},
	}

	service := NewUserService(mockRepo)
	user, err := service.FindUser(42)
	if err != nil || user.ID != 42 || user.Name != "Bob" {
		t.Errorf("unexpected result: user=%v, err=%v", user, err)
	}
}
```

---

Tu peux tester tout cela avec :

```sh
go test ./...
```
