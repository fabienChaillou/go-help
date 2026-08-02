Voici un exemple complet de test en Go utilisant :

* pgx (driver PostgreSQL)
* pgxmock (mock compatible pgx)

On teste la requête :

```sql
SELECT id, type, bu, service, environment, action, realms, accounts, message, created_by, created_at, update_by, update_by
FROM rules
WHERE bu = $1
ORDER BY created_at DESC
```

---

# 1️⃣ Exemple Repository

```go
package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
)

type Rule struct {
	ID          int64
	Type        string
	BU          string
	Service     string
	Environment string
	Action      string
	Realms      []string
	Accounts    []string
	Message     string
	CreatedBy   string
	CreatedAt   time.Time
	UpdateBy    string
}

type RuleRepository struct {
	db pgx.Tx // ou interface compatible (voir test)
}

func NewRuleRepository(db pgx.Tx) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) ListByBU(ctx context.Context, bu string) ([]Rule, error) {
	query := `
SELECT id, type, bu, service, environment, action, realms, accounts, message, created_by, created_at, update_by, update_by
FROM rules
WHERE bu = $1
ORDER BY created_at DESC`

	rows, err := r.db.Query(ctx, query, bu)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []Rule

	for rows.Next() {
		var rule Rule
		err := rows.Scan(
			&rule.ID,
			&rule.Type,
			&rule.BU,
			&rule.Service,
			&rule.Environment,
			&rule.Action,
			&rule.Realms,
			&rule.Accounts,
			&rule.Message,
			&rule.CreatedBy,
			&rule.CreatedAt,
			&rule.UpdateBy,
			&rule.UpdateBy, // ⚠ doublon dans la requête
		)
		if err != nil {
			return nil, err
		}
		rules = append(rules, rule)
	}

	return rules, nil
}
```

---

# 2️⃣ Test avec pgxmock

```go
package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"yourmodule/repository"
)

func TestListByBU(t *testing.T) {
	mock, err := pgxmock.NewConn()
	assert.NoError(t, err)
	defer mock.Close(context.Background())

	repo := repository.NewRuleRepository(mock)

	query := `
SELECT id, type, bu, service, environment, action, realms, accounts, message, created_by, created_at, update_by, update_by
FROM rules
WHERE bu = $1
ORDER BY created_at DESC`

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "type", "bu", "service", "environment", "action",
		"realms", "accounts", "message",
		"created_by", "created_at", "update_by", "update_by",
	}).AddRow(
		int64(1),
		"security",
		"BU1",
		"svc1",
		"prod",
		"allow",
		[]string{"realm1"},
		[]string{"acc1"},
		"test message",
		"admin",
		now,
		"admin",
		"admin",
	)

	mock.ExpectQuery(query).
		WithArgs("BU1").
		WillReturnRows(rows)

	result, err := repo.ListByBU(context.Background(), "BU1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "BU1", result[0].BU)
	assert.Equal(t, "security", result[0].Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

# ⚠ Points importants

### 1️⃣ Attention au doublon

Dans ta requête tu as :

```sql
update_by, update_by
```

Ce qui veut dire :

* tu dois scanner deux colonnes
* même si c’est la même valeur

Probablement tu voulais :

```sql
update_by, updated_at
```

---

### 2️⃣ Si tu utilises pgxpool

Il faut injecter une interface :

```go
type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
```

Et utiliser cette interface dans ton repository.
`pgxmock` implémente cette interface.

---

# ✅ Résultat

Ce test :

* Mock la requête exacte
* Vérifie les arguments (`$1`)
* Simule un retour SQL
* Vérifie que le mapping fonctionne
* Vérifie que toutes les expectations sont respectées

---

Si tu veux, je peux te faire :

* 🔥 une version avec `pgxpool.Pool`
* 🔥 une version avec `sqlc`
* 🔥 une version avec `Squirrel`
* 🔥 un exemple avec gestion d’erreur testée
* 🔥 une version clean architecture avec interface mockable propre

Tu travailles toujours avec `pgx v5` ?

## version `pgxpool.Pool` && `Squirrel`

Parfait 👌
Voici une version propre avec :

* pgx v5
* pgxpool
* Squirrel
* pgxmock

Architecture propre + interface mockable.

---

# ✅ 1️⃣ Repository avec pgxpool + Squirrel

## Interface Queryer (mockable)

```go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
```

`pgxpool.Pool` implémente déjà cette méthode ✅
`pgxmock` aussi ✅

---

## Repository

```go
package repository

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
)

type Rule struct {
	ID          int64
	Type        string
	BU          string
	Service     string
	Environment string
	Action      string
	Realms      []string
	Accounts    []string
	Message     string
	CreatedBy   string
	CreatedAt   time.Time
	UpdateBy    string
}

type RuleRepository struct {
	db Queryer
}

func NewRuleRepository(db Queryer) *RuleRepository {
	return &RuleRepository{db: db}
}

func (r *RuleRepository) ListByBU(ctx context.Context, bu string) ([]Rule, error) {

	builder := sq.
		Select(
			"id",
			"type",
			"bu",
			"service",
			"environment",
			"action",
			"realms",
			"accounts",
			"message",
			"created_by",
			"created_at",
			"update_by",
		).
		From("rules").
		Where(sq.Eq{"bu": bu}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []Rule

	for rows.Next() {
		var rule Rule
		err := rows.Scan(
			&rule.ID,
			&rule.Type,
			&rule.BU,
			&rule.Service,
			&rule.Environment,
			&rule.Action,
			&rule.Realms,
			&rule.Accounts,
			&rule.Message,
			&rule.CreatedBy,
			&rule.CreatedAt,
			&rule.UpdateBy,
		)
		if err != nil {
			return nil, err
		}

		rules = append(rules, rule)
	}

	return rules, rows.Err()
}
```

---

# ✅ 2️⃣ Test avec pgxmock

```go
package repository_test

import (
	"context"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"yourmodule/repository"
)

func TestListByBU(t *testing.T) {

	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewRuleRepository(mock)

	// ⚠ On reconstruit EXACTEMENT la même requête Squirrel
	builder := sq.
		Select(
			"id",
			"type",
			"bu",
			"service",
			"environment",
			"action",
			"realms",
			"accounts",
			"message",
			"created_by",
			"created_at",
			"update_by",
		).
		From("rules").
		Where(sq.Eq{"bu": "BU1"}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	assert.NoError(t, err)

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "type", "bu", "service", "environment", "action",
		"realms", "accounts", "message",
		"created_by", "created_at", "update_by",
	}).AddRow(
		int64(1),
		"security",
		"BU1",
		"svc1",
		"prod",
		"allow",
		[]string{"realm1"},
		[]string{"acc1"},
		"test message",
		"admin",
		now,
		"admin",
	)

	mock.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(rows)

	result, err := repo.ListByBU(context.Background(), "BU1")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "BU1", result[0].BU)
	assert.Equal(t, "security", result[0].Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

# 🎯 Pourquoi cette version est clean ?

### ✔ Interface Queryer

Permet d’injecter :

* `*pgxpool.Pool`
* `pgxmock.PgxPoolIface`
* un fake maison si besoin

### ✔ Squirrel utilisé des deux côtés

On évite les erreurs :

* d’espaces
* d’ORDER BY
* de placeholders `$1`

### ✔ Test robuste

* vérifie la requête exacte générée
* vérifie les args
* vérifie le mapping
* vérifie les expectations

---

# 🔥 Bonus — Version production

En prod :

```go
pool, _ := pgxpool.New(ctx, databaseURL)
repo := repository.NewRuleRepository(pool)
```

---

# 🚀 Si tu veux aller plus loin

Je peux te faire :

* version générique type `BaseRepository[T]`
* test avec erreur simulée (`WillReturnError`)
* test avec plusieurs rows
* test transactionnel
* exemple clean architecture complet
* version compatible avec ton repo générique actuel

Vu que tu bosses déjà avec pgx + repository générique,
tu veux que je l’intègre dans ton pattern existant ?

## Version avec Generic to golang

Parfait 👌
On va faire une version **pgx v5 + pgxpool + Squirrel + Repository générique typé (Go 1.20+)**, propre et testable avec pgxmock.

Stack utilisée :

* pgx
* pgxpool
* Squirrel
* pgxmock

---

# 🎯 Objectif

Créer :

```go
repo := repository.NewRepository[Rule](pool, scanRule)
rules, err := repo.List(ctx, squirrelBuilder)
```

Donc :

* Repository générique `Repository[T]`
* Fonction de scan injectée (clean & flexible)
* Compatible pgxpool ET pgxmock
* Testable facilement

---

# ✅ 1️⃣ Interface Queryer (mockable)

```go
package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Queryer interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
```

✔ `*pgxpool.Pool` implémente déjà ça
✔ `pgxmock` aussi

---

# ✅ 2️⃣ Repository générique

```go
package repository

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type ScannerFunc[T any] func(rows interface {
	Scan(dest ...any) error
}) (T, error)

type Repository[T any] struct {
	db       Queryer
	scanFunc ScannerFunc[T]
}

func NewRepository[T any](db Queryer, scan ScannerFunc[T]) *Repository[T] {
	return &Repository[T]{
		db:       db,
		scanFunc: scan,
	}
}

func (r *Repository[T]) List(ctx context.Context, builder sq.SelectBuilder) ([]T, error) {

	query, args, err := builder.PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []T

	for rows.Next() {
		entity, err := r.scanFunc(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, entity)
	}

	return result, rows.Err()
}
```

---

# ✅ 3️⃣ Entité Rule + scanRule

```go
package repository

import (
	"time"

	"github.com/jackc/pgx/v5"
)

type Rule struct {
	ID          int64
	Type        string
	BU          string
	Service     string
	Environment string
	Action      string
	Realms      []string
	Accounts    []string
	Message     string
	CreatedBy   string
	CreatedAt   time.Time
	UpdateBy    string
}

func ScanRule(rows interface {
	Scan(dest ...any) error
}) (Rule, error) {

	var r Rule

	err := rows.Scan(
		&r.ID,
		&r.Type,
		&r.BU,
		&r.Service,
		&r.Environment,
		&r.Action,
		&r.Realms,
		&r.Accounts,
		&r.Message,
		&r.CreatedBy,
		&r.CreatedAt,
		&r.UpdateBy,
	)

	return r, err
}
```

---

# ✅ 4️⃣ Utilisation avec Squirrel

```go
builder := sq.
	Select(
		"id",
		"type",
		"bu",
		"service",
		"environment",
		"action",
		"realms",
		"accounts",
		"message",
		"created_by",
		"created_at",
		"update_by",
	).
	From("rules").
	Where(sq.Eq{"bu": "BU1"}).
	OrderBy("created_at DESC")

repo := repository.NewRepository[repository.Rule](pool, repository.ScanRule)

rules, err := repo.List(ctx, builder)
```

---

# ✅ 5️⃣ Test avec pgxmock

```go
package repository_test

import (
	"context"
	"testing"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"

	"yourmodule/repository"
)

func TestGenericRepository_ListRules(t *testing.T) {

	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := repository.NewRepository[repository.Rule](mock, repository.ScanRule)

	builder := sq.
		Select(
			"id",
			"type",
			"bu",
			"service",
			"environment",
			"action",
			"realms",
			"accounts",
			"message",
			"created_by",
			"created_at",
			"update_by",
		).
		From("rules").
		Where(sq.Eq{"bu": "BU1"}).
		OrderBy("created_at DESC").
		PlaceholderFormat(sq.Dollar)

	query, args, _ := builder.ToSql()

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"id", "type", "bu", "service", "environment", "action",
		"realms", "accounts", "message",
		"created_by", "created_at", "update_by",
	}).AddRow(
		int64(1),
		"security",
		"BU1",
		"svc1",
		"prod",
		"allow",
		[]string{"realm1"},
		[]string{"acc1"},
		"test message",
		"admin",
		now,
		"admin",
	)

	mock.ExpectQuery(query).
		WithArgs(args...).
		WillReturnRows(rows)

	result, err := repo.List(context.Background(), builder)

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, int64(1), result[0].ID)
	assert.Equal(t, "BU1", result[0].BU)

	assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

# 🧠 Pourquoi c’est puissant ?

✅ 100% générique
✅ Réutilisable pour toutes tes entités
✅ Pas de reflection
✅ Compatible pgxpool / transaction / mock
✅ Clean Architecture ready

---

# 🔥 Version encore plus avancée possible

Je peux te faire :

* `GetOne`
* `Insert`
* `Update`
* `Delete`
* support transaction `WithTx`
* repository base avec pagination
* repository générique + Specification pattern
* intégration avec ton repo générique existant sqlx/Squirrel

Vu que tu bosses déjà sur une archi MVC + validation + auth,
tu veux que je te fasse une base repository complète prête pour production ?


QueryMatcherRegexp into pashagolub module pgmock