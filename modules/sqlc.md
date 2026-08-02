Parfait, on va faire ça **pas à pas**, avec un **vrai workflow sqlc** comme en prod.
À la fin, tu sauras exactement **où écrire quoi**, **ce qui est généré**, et **comment l’utiliser avec pgx + transactions**.

---

# 🧬 sqlc — le guide clair et concret

## 🧠 L’idée de sqlc (en une phrase)

👉 **Tu écris du SQL**, sqlc génère **du Go typé**, sans ORM, sans magie runtime.

---

# 1️⃣ Installer sqlc

### CLI

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Vérifie :

```bash
sqlc version
```

---

# 2️⃣ Structure recommandée

```text
internal/
├── db/
│   ├── migrations/
│   │   └── 001_init.sql
│   ├── queries/
│   │   └── users.sql
│   └── generated/
│       └── sqlc/
│           ├── db.go
│           ├── models.go
│           └── users.sql.go
```

➡️ **Tu écris seulement** :

* `migrations/*.sql`
* `queries/*.sql`

➡️ **sqlc génère le reste**

---

# 3️⃣ Écrire le schéma SQL

### migrations/001_init.sql

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    action TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
```

---

# 4️⃣ Écrire les queries sqlc

### queries/users.sql

```sql
-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id, name, email;

-- name: GetUserByID :one
SELECT id, name, email
FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET name = $2, email = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;
```

### 🔑 Règles importantes

* `-- name:` → nom de la fonction Go
* `:one` → retourne une ligne
* `:many` → plusieurs lignes
* `:exec` → pas de retour
* `:execrows` → retourne le nombre de lignes affectées

---

# 5️⃣ Configurer sqlc

### sqlc.yaml

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "internal/db/migrations"
    queries: "internal/db/queries"
    gen:
      go:
        package: "sqlc"
        out: "internal/db/generated/sqlc"
        sql_package: "pgx/v5"
```

👉 `sql_package: pgx/v5` = **pgx natif + transactions compatibles**

---

# 6️⃣ Générer le code

```bash
sqlc generate
```

Ça génère :

```text
db.go        // DBTX, Queries, WithTx
models.go   // struct User
users.sql.go
```

---

# 7️⃣ Utiliser sqlc (sans transaction)

```go
pool, _ := pgxpool.New(ctx, dsn)
queries := sqlc.New(pool)

user, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
	Name:  "Alice",
	Email: "alice@mail.com",
})
```

➡️ **Typé, safe, zéro mapping manuel**

---

# 8️⃣ Utiliser sqlc avec transaction (important 🔥)

### Pattern officiel sqlc

```go
tx, err := pool.Begin(ctx)
if err != nil {
	return err
}
defer tx.Rollback(ctx)

qtx := queries.WithTx(tx)

user, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
	Name:  "Bob",
	Email: "bob@mail.com",
})
if err != nil {
	return err
}

err = qtx.CreateAuditLog(ctx, sqlc.CreateAuditLogParams{
	Action: "user_created",
	UserID: user.ID,
})
if err != nil {
	return err
}

return tx.Commit(ctx)
```

✔️ même API
✔️ même typage
✔️ rollback auto

---

# 9️⃣ sqlc + Repository Pattern (optionnel)

Tu peux **encapsuler sqlc** si tu veux garder une archi Clean.

```go
type UserRepository struct {
	q *sqlc.Queries
}

func (r *UserRepository) Get(ctx context.Context, id int64) (sqlc.User, error) {
	return r.q.GetUserByID(ctx, id)
}
```

Mais honnêtement 👇
👉 **sqlc = repository déjà très propre**

---

# 🔥 sqlc vs pgx vs gorm (rappel rapide)

| Critère             | sqlc     | pgx      | gorm        |
| ------------------- | -------- | -------- | ----------- |
| SQL explicite       | ✅        | ✅        | ❌           |
| Typage compile-time | ✅        | ❌        | ❌           |
| Transactions        | ✅ native | manuelle | opaque      |
| Boilerplate         | faible   | élevé    | très faible |
| Long terme          | 🏆       | 👍       | ⚠️          |

---

# 🚨 Erreurs fréquentes

❌ écrire du SQL dynamique (sqlc n’aime pas)
❌ mettre la logique métier dans les queries
❌ modifier le code généré
❌ oublier `WithTx`

---

# 🧠 Mon conseil réel

> **Projet backend sérieux en Go + PostgreSQL ?**
> 👉 **sqlc + pgx + transactions au niveau service**

---

Si tu veux la suite, je peux te montrer :

* 🧪 **sqlc + testcontainers**
* 🔁 **savepoints / nested tx**
* 🏗 **Clean Architecture complète**
* 🚀 **Gin + handlers + sqlc**
* 📦 **sqlc dans un monorepo**

Tu choisis 😄


#############################################

Parfait, là on passe en **mode vraie prod** 😄
Je te montre **comment faire proprement** :

1. **sqlc dans un monorepo**
2. **sqlc + testcontainers (tests d’intégration fiables)**

Je reste pragmatique, avec des patterns que tu peux copier-coller.

---

# 1️⃣ sqlc dans un **monorepo**

## 🎯 Problème à résoudre

* plusieurs services
* chacun sa DB (ou schéma)
* pas de conflits de génération
* code SQL clair et isolé

---

## 🧱 Structure de monorepo recommandée

```text
monorepo/
├── go.work
├── services/
│   ├── user-service/
│   │   ├── go.mod
│   │   ├── internal/
│   │   │   └── db/
│   │   │       ├── migrations/
│   │   │       ├── queries/
│   │   │       └── generated/sqlc/
│   │   └── sqlc.yaml
│   ├── order-service/
│   │   ├── go.mod
│   │   ├── internal/db/...
│   │   └── sqlc.yaml
│   └── payment-service/
│       └── ...
```

👉 **1 service = 1 sqlc.yaml = 1 DB boundary**

---

## 🧠 go.work (racine)

```bash
go work init ./services/user-service ./services/order-service
```

✔️ partage des versions Go
✔️ build rapide
✔️ isolation logique

---

## 📜 sqlc.yaml (par service)

### services/user-service/sqlc.yaml

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "internal/db/migrations"
    queries: "internal/db/queries"
    gen:
      go:
        package: "sqlc"
        out: "internal/db/generated/sqlc"
        sql_package: "pgx/v5"
```

---

## 🔁 Génération

Depuis la racine :

```bash
sqlc generate -f services/user-service/sqlc.yaml
sqlc generate -f services/order-service/sqlc.yaml
```

Ou via Makefile :

```makefile
sqlc:
	sqlc generate -f services/user-service/sqlc.yaml
	sqlc generate -f services/order-service/sqlc.yaml
```

---

## 🧠 Bonnes pratiques monorepo + sqlc

✅ **jamais partager les queries entre services**
✅ **une DB = un service**
✅ **sqlc reste interne (`internal/`)**
✅ **migration versionnée par service**

> sqlc = contrat DB → service
> pas un contrat inter-services

---

# 2️⃣ sqlc + **Testcontainers**

Là on teste **le vrai SQL**, pas des mocks bidons.

---

## 🎯 Objectif

* lancer PostgreSQL dans Docker
* appliquer migrations
* tester les queries sqlc
* teardown automatique

---

## 📦 Dépendances

```bash
go get github.com/testcontainers/testcontainers-go
go get github.com/testcontainers/testcontainers-go/modules/postgres
```

---

## 🧪 Setup commun de test

### internal/db/testutil/postgres.go

```go
package testutil

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func NewTestPostgres(ctx context.Context) (*pgxpool.Pool, func(), error) {
	container, err := postgres.Run(
		ctx,
		"postgres:15",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		return nil, nil, err
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, nil, err
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		pool.Close()
		container.Terminate(ctx)
	}

	return pool, cleanup, nil
}
```

---

## 🔄 Appliquer les migrations

Tu peux utiliser `golang-migrate` ou exec simple :

```go
func ApplyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
		CREATE TABLE users (
			id BIGSERIAL PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT NOT NULL UNIQUE
		);
	`)
	return err
}
```

(en vrai → `golang-migrate`, mais là on reste focus)

---

## 🧪 Test sqlc simple

### internal/db/generated/sqlc/users_test.go

```go
func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	db, cleanup, err := testutil.NewTestPostgres(ctx)
	require.NoError(t, err)
	defer cleanup()

	require.NoError(t, testutil.ApplyMigrations(ctx, db))

	queries := sqlc.New(db)

	user, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  "Alice",
		Email: "alice@test.com",
	})
	require.NoError(t, err)
	require.NotZero(t, user.ID)
}
```

---

## 🧪 Test transaction sqlc

```go
func TestCreateUserWithTx(t *testing.T) {
	ctx := context.Background()

	db, cleanup, _ := testutil.NewTestPostgres(ctx)
	defer cleanup()

	testutil.ApplyMigrations(ctx, db)

	queries := sqlc.New(db)

	tx, err := db.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  "Bob",
		Email: "bob@test.com",
	})
	require.NoError(t, err)

	require.NoError(t, tx.Commit(ctx))
}
```

---

## 🔥 Pourquoi c’est le combo parfait

| Aspect         | Mock | Testcontainers |
| -------------- | ---- | -------------- |
| SQL réel       | ❌    | ✅              |
| Contraintes DB | ❌    | ✅              |
| Index / unique | ❌    | ✅              |
| Confiance prod | ⚠️   | 🏆             |

👉 sqlc **mérite** des tests d’intégration

---

# 3️⃣ Astuces avancées 💡

### ⚡ Accélérer les tests

* un container par package
* reuse container (`TESTCONTAINERS_RYUK_DISABLED`)
* migrations en `TestMain`

### 🧼 Isolation

* un schéma par test
* rollback global
* truncate tables

### 🔁 CI

* aucun PostgreSQL à installer
* Docker suffit
* reproductible à 100%

---

# 🧠 Résumé clair

✔️ sqlc s’intègre **parfaitement** dans un monorepo
✔️ 1 service = 1 sqlc boundary
✔️ sqlc + testcontainers = **confiance maximale**
✔️ zéro ORM, zéro magie

---

Si tu veux aller encore plus loin, je peux te montrer :

* 🧪 **TestMain optimisé (1 container / suite)**
* 🔁 **savepoints sqlc**
* 🏗 **Clean Architecture complète monorepo**
* 🚀 **Gin + handlers testés**
* 📦 **Makefile / CI GitHub Actions**

Tu veux quoi pour la suite ? 😄

#####################################################

Parfait 😎 ! On va faire un **niveau supérieur de tests avec sqlc** :

* **Savepoints SQLC** → transactions imbriquées sûres.
* **TestMain optimisé** → un seul container PostgreSQL pour toute la suite de tests.

Je te montre **un pattern testé en production**.

---

# 1️⃣ Savepoints avec sqlc

## 🎯 Objectif

* Pouvoir faire des sous-transactions
* Rollback partiel si une étape échoue
* Test de cas complexes sans rollback global

---

### 💡 Exemple simple : repository + savepoint

```go
func TestSavepoint(t *testing.T) {
	ctx := context.Background()

	db, cleanup, _ := testutil.NewTestPostgres(ctx)
	defer cleanup()
	testutil.ApplyMigrations(ctx, db)

	queries := sqlc.New(db)

	// transaction principale
	tx, err := db.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	qtx := queries.WithTx(tx)

	// création user
	user, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  "Charlie",
		Email: "charlie@test.com",
	})
	require.NoError(t, err)

	// --- savepoint ---
	_, err = tx.Exec(ctx, "SAVEPOINT sp1")
	require.NoError(t, err)

	// tentative d'opération qui pourrait échouer
	_, err = qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  "Duplicate",
		Email: "charlie@test.com", // violation unique
	})
	if err != nil {
		// rollback partiel → retourne au savepoint
		_, err2 := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT sp1")
		require.NoError(t, err2)
	}

	// opération valide après rollback partiel
	user2, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  "Delta",
		Email: "delta@test.com",
	})
	require.NoError(t, err)

	// commit final
	require.NoError(t, tx.Commit(ctx))
	require.NotZero(t, user2.ID)
}
```

### ✅ Points clés

* `SAVEPOINT sp1` → point de retour
* `ROLLBACK TO SAVEPOINT sp1` → rollback partiel
* `COMMIT` final → persiste tout ce qui n’a pas été rollbacké
* Compatible **sqlc** + **pgx** (tx natif)

---

# 2️⃣ TestMain optimisé : 1 container pour toute la suite

Au lieu de créer un container **par test**, on crée **1 container global**, on applique migrations, et chaque test utilise ce pool.

---

### 🧪 internal/db/testutil/testmain.go

```go
package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

var TestDB *pgxpool.Pool
var cleanup func()
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()

	var err error
	TestDB, cleanup, err = NewTestPostgres(ctx)
	if err != nil {
		panic(err)
	}

	// appliquer migrations une seule fois
	if err := ApplyMigrations(ctx, TestDB); err != nil {
		panic(err)
	}

	// lancer la suite de tests
	code := m.Run()

	// teardown
	cleanup()

	os.Exit(code)
}
```

---

### ⚡ Exemple d’utilisation dans un test

```go
func TestCreateUser(t *testing.T) {
	queries := sqlc.New(TestDB)

	user, err := queries.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:  "Alice",
		Email: "alice@test.com",
	})
	require.NoError(t, err)
	require.NotZero(t, user.ID)
}
```

### ✅ Avantages

* 1 container → tests beaucoup plus rapides
* migrations appliquées **une seule fois**
* tous les tests partagent la même base isolée
* compatible avec **savepoints** pour rollback partiel

---

# 3️⃣ Combiner savepoints + TestMain

* Crée container global → migrations appliquées
* Chaque test démarre une transaction
* Si besoin de rollback partiel → **savepoint**
* Commit final ou rollback global selon scénario

```go
func TestUserComplexFlow(t *testing.T) {
	tx, err := TestDB.Begin(context.Background())
	require.NoError(t, err)
	defer tx.Rollback(context.Background())

	qtx := sqlc.New(TestDB).WithTx(tx)

	// création user
	_, _ = qtx.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:  "X",
		Email: "x@test.com",
	})

	// savepoint pour essais risqués
	_, _ = tx.Exec(context.Background(), "SAVEPOINT sp1")
	_, err = qtx.CreateUser(context.Background(), sqlc.CreateUserParams{
		Name:  "Y",
		Email: "x@test.com", // duplicate email
	})
	if err != nil {
		_, _ = tx.Exec(context.Background(), "ROLLBACK TO SAVEPOINT sp1")
	}

	// commit final
	require.NoError(t, tx.Commit(context.Background()))
}
```

---

# 4️⃣ Résumé 🔥

| Concept         | Objectif                                          |
| --------------- | ------------------------------------------------- |
| Savepoint       | rollback partiel dans une transaction             |
| TestMain global | 1 container PostgreSQL pour toute la suite        |
| sqlc + tx       | typé, safe, transactions natives                  |
| pgx             | driver performant pour transactions et savepoints |

**Avantages combinés**

* Tests rapides et réalistes
* Confiance totale sur SQL
* Maintenance simple dans un monorepo
* Isolation test parfaite via transactions + savepoints

---

Si tu veux, je peux te faire **un exemple complet de monorepo Go avec sqlc**, **Gin**, **transactions**, **savepoints**, et **tests intégrés** prêt à cloner.

Veux‑tu que je fasse ça ?


##################### EXAMPLE COMPLET ###########################

Parfait 😎 ! Je vais te construire un **exemple complet et prêt à l’emploi** :

* **Monorepo Go**
* **Service `user-service`**
* **Gin pour l’API**
* **sqlc + PostgreSQL (pgx)**
* **Transactions + savepoints**
* **Testcontainers + TestMain optimisé**

Je vais te faire une **structure claire**, tu pourras cloner et adapter directement.

---

# 1️⃣ Structure du monorepo

```text
monorepo/
├── go.work
├── services/
│   └── user-service/
│       ├── go.mod
│       ├── main.go
│       ├── sqlc.yaml
│       ├── internal/
│       │   ├── db/
│       │   │   ├── migrations/
│       │   │   │   └── 001_init.sql
│       │   │   ├── queries/
│       │   │   │   └── users.sql
│       │   │   └── generated/sqlc/       # généré par sqlc
│       │   ├── domain/
│       │   │   └── user.go
│       │   ├── repository/
│       │   │   └── user_repository.go
│       │   ├── service/
│       │   │   └── user_service.go
│       │   └── testutil/
│       │       ├── postgres.go
│       │       └── testmain.go
│       └── handler/
│           └── user_handler.go
```

---

# 2️⃣ SQLC Config

### services/user-service/sqlc.yaml

```yaml
version: "2"
sql:
  - engine: "postgresql"
    schema: "internal/db/migrations"
    queries: "internal/db/queries"
    gen:
      go:
        package: "sqlc"
        out: "internal/db/generated/sqlc"
        sql_package: "pgx/v5"
```

---

# 3️⃣ Migration initiale

### internal/db/migrations/001_init.sql

```sql
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE audit_logs (
    id BIGSERIAL PRIMARY KEY,
    action TEXT NOT NULL,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
```

---

# 4️⃣ Queries SQL

### internal/db/queries/users.sql

```sql
-- name: CreateUser :one
INSERT INTO users (name, email)
VALUES ($1, $2)
RETURNING id, name, email;

-- name: GetUserByID :one
SELECT id, name, email
FROM users
WHERE id = $1;

-- name: UpdateUser :exec
UPDATE users
SET name = $2, email = $3
WHERE id = $1;

-- name: DeleteUser :exec
DELETE FROM users
WHERE id = $1;

-- name: CreateAuditLog :exec
INSERT INTO audit_logs (action, user_id)
VALUES ($1, $2);
```

---

# 5️⃣ Domaine

### internal/domain/user.go

```go
package domain

type User struct {
	ID    int64
	Name  string
	Email string
}
```

---

# 6️⃣ Repository

### internal/repository/user_repository.go

```go
package repository

import (
	"context"
	"monorepo/services/user-service/internal/domain"
	"monorepo/services/user-service/internal/db/generated/sqlc"
)

type UserRepository struct {
	queries *sqlc.Queries
}

func NewUserRepository(queries *sqlc.Queries) *UserRepository {
	return &UserRepository{queries: queries}
}

func (r *UserRepository) Create(ctx context.Context, name, email string) (domain.User, error) {
	u, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name:  name,
		Email: email,
	})
	return domain.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}, err
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.queries.GetUserByID(ctx, id)
	return domain.User{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}, err
}
```

---

# 7️⃣ Service (transactions + savepoints)

### internal/service/user_service.go

```go
package service

import (
	"context"
	"monorepo/services/user-service/internal/domain"
	"monorepo/services/user-service/internal/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserService struct {
	db   *pgxpool.Pool
	repo *repository.UserRepository
}

func NewUserService(db *pgxpool.Pool, repo *repository.UserRepository) *UserService {
	return &UserService{db: db, repo: repo}
}

// Create user + audit log dans une transaction
func (s *UserService) CreateUserWithAudit(ctx context.Context, name, email string) (domain.User, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return domain.User{}, err
	}
	defer tx.Rollback(ctx)

	qtx := s.repo.queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, sqlc.CreateUserParams{Name: name, Email: email})
	if err != nil {
		return domain.User{}, err
	}

	// savepoint pour test rollback partiel
	_, _ = tx.Exec(ctx, "SAVEPOINT sp1")

	err = qtx.CreateAuditLog(ctx, sqlc.CreateAuditLogParams{
		Action: "user_created",
		UserID: user.ID,
	})
	if err != nil {
		_, _ = tx.Exec(ctx, "ROLLBACK TO SAVEPOINT sp1")
	}

	return domain.User{ID: user.ID, Name: user.Name, Email: user.Email}, tx.Commit(ctx)
}
```

---

# 8️⃣ Gin handler

### handler/user_handler.go

```go
package handler

import (
	"net/http"
	"monorepo/services/user-service/internal/service"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.svc.CreateUserWithAudit(c.Request.Context(), req.Name, req.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
```

---

# 9️⃣ Main

### main.go

```go
package main

import (
	"context"
	"log"
	"monorepo/services/user-service/internal/db/generated/sqlc"
	"monorepo/services/user-service/internal/repository"
	"monorepo/services/user-service/internal/service"
	"monorepo/services/user-service/handler"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	dsn := "postgres://user:password@localhost:5432/mydb?sslmode=disable"

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatal(err)
	}

	queries := sqlc.New(pool)
	repo := repository.NewUserRepository(queries)
	svc := service.NewUserService(pool, repo)
	h := handler.NewUserHandler(svc)

	r := gin.Default()
	r.POST("/users", h.CreateUser)

	r.Run(":8080")
}
```

---

# 🔟 Tests + Testcontainers + TestMain

### internal/testutil/postgres.go

```go
package testutil

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func NewTestPostgres(ctx context.Context) (*pgxpool.Pool, func(), error) {
	container, err := postgres.Run(ctx,
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
	)
	if err != nil {
		return nil, nil, err
	}

	dsn, _ := container.ConnectionString(ctx, "sslmode=disable")
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, nil, err
	}

	cleanup := func() {
		pool.Close()
		container.Terminate(ctx)
	}

	return pool, cleanup, nil
}

func ApplyMigrations(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, `
	CREATE TABLE users (id BIGSERIAL PRIMARY KEY, name TEXT, email TEXT UNIQUE);
	CREATE TABLE audit_logs (id BIGSERIAL PRIMARY KEY, action TEXT, user_id BIGINT, created_at TIMESTAMPTZ DEFAULT now());
	`)
	return err
}
```

### internal/testutil/testmain.go

```go
package testutil

import (
	"context"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

var TestDB *pgxpool.Pool
var cleanup func()
var ctx context.Context

func TestMain(m *testing.M) {
	ctx = context.Background()
	var err error
	TestDB, cleanup, err = NewTestPostgres(ctx)
	if err != nil {
		panic(err)
	}

	if err := ApplyMigrations(ctx, TestDB); err != nil {
		panic(err)
	}

	code := m.Run()
	cleanup()
	os.Exit(code)
}
```

### Exemple de test avec savepoints

```go
package repository_test

import (
	"context"
	"testing"

	"monorepo/services/user-service/internal/db/generated/sqlc"
	"monorepo/services/user-service/internal/repository"
	"monorepo/services/user-service/internal/testutil"

	"github.com/stretchr/testify/require"
)

func TestCreateUserWithSavepoint(t *testing.T) {
	ctx := context.Background()

	tx, err := testutil.TestDB.Begin(ctx)
	require.NoError(t, err)
	defer tx.Rollback(ctx)

	queries := sqlc.New(testutil.TestDB).WithTx(tx)
	repo := repository.NewUserRepository(queries)

	user, err := repo.Create(ctx, "Savepoint", "sp@test.com")
	require.NoError(t, err)
	require.NotZero(t, user.ID)

	_, err = tx.Exec(ctx, "SAVEPOINT sp1")
	require.NoError(t, err)

	// tentative d'opération fail
	_, err = repo.Create(ctx, "Duplicate", "sp@test.com")
	if err != nil {
		_, err2 := tx.Exec(ctx, "ROLLBACK TO SAVEPOINT sp1")
		require.NoError(t, err2)
	}

	require.NoError(t, tx.Commit(ctx))
}
```

---

✅ Avec ça tu as :

* Monorepo Go + sqlc par service
* Gin REST API CRUD
* Repository pattern
* Transactions + savepoints
* Testcontainers + TestMain global

---

