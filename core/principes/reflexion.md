Pour installer **GoMock** et son outil de génération **mockgen** sur macOS, voici les étapes à suivre :

---

### ✅ Prérequis

Assure-toi que Go est déjà installé :

```bash
go version
```

---

### 1. 📦 Installer `GoMock` (librairie)

Aucune installation spéciale n’est requise pour la bibliothèque elle-même : elle s’installe automatiquement avec `go get` ou lors du `go mod tidy`.

Tu peux l’ajouter à ton projet avec :

```bash
go get go.uber.org/mock
```

---

### 2. 🔧 Installer `mockgen` (outil CLI)

#### Option 1 : via `go install` (recommandé)

```bash
go install go.uber.org/mock/mockgen@latest
```

Cela installe l’exécutable `mockgen` dans ton `$GOPATH/bin` ou `$HOME/go/bin`.

#### Vérifie l’installation :

```bash
mockgen --version
```

Si tu obtiens une erreur de commande introuvable, ajoute `$HOME/go/bin` à ton `PATH` :

```bash
echo 'export PATH=$PATH:$HOME/go/bin' >> ~/.zshrc
source ~/.zshrc
```

(ou `~/.bash_profile` si tu utilises bash)

---

### 3. 📄 Générer un mock avec `mockgen`

Voici un exemple de commande :

```bash
mockgen -source=your_interface.go -destination=your_interface_mock.go -package=yourpackage
```

---

## Example with mockgen
Voici un exemple complet d'utilisation de **GoMock** avec `mockgen`, incluant une interface, un mock généré, et un test unitaire.

---

## 🧩 Exemple : Interface `UserRepository` avec un test mocké

### 1. 📄 Interface à mocker (dans `user/repository.go`)

```go
package user

type User struct {
	ID   int
	Name string
}

type UserRepository interface {
	GetUserByID(id int) (*User, error)
}
```

---

### 2. 🧪 Fichier de test (dans `user/repository_test.go`)

On suppose que le mock est généré dans `user/mocks/mock_user_repository.go`.

Voici comment utiliser le mock dans un test :

```go
package user_test

import (
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/tonrepo/user"                       // change selon ton import réel
	mock_user "github.com/tonrepo/user/mocks"       // idem

	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock_user.NewMockUserRepository(ctrl)

	expectedUser := &user.User{ID: 1, Name: "Alice"}

	// Expectation
	mockRepo.EXPECT().
		GetUserByID(1).
		Return(expectedUser, nil)

	// Call
	u, err := mockRepo.GetUserByID(1)

	assert.NoError(t, err)
	assert.Equal(t, expectedUser, u)
}
```

---

### 3. 🛠️ Génération du mock avec `mockgen`

Tu peux générer le mock avec cette commande :

```bash
mockgen -source=user/repository.go -destination=user/mocks/mock_user_repository.go -package=mocks
```

> 🔁 Si tu veux l’automatiser, tu peux l’ajouter à un `Makefile` ou script.

---

### 4. 📦 Dépendances utiles

Ajoute ceci dans ton `go.mod` :

```bash
go get github.com/golang/mock/gomock
go get github.com/stretchr/testify
```

---
