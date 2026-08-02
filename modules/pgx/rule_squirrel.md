Voici la conversion de ta requête **GORM → pgx + Squirrel**.

Je reprends exactement la logique :

* filtres `bu`, `service`
* environnement avec fallback `EnvAll`
* action avec fallback `ActAll`
* filtres `accounts` et `realms` avec `ANY`
* ORDER BY avec le `CASE WHEN type='ALLOW'`
* tri secondaire `created_at DESC`

---

# 🔹 Version pgx + Squirrel

```go
import (
	"context"
	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

func (r *Repo) Evaluate(ctx context.Context, evalCtx EvalContext) (bool, *Rule, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	qb := psql.
		Select("*").
		From("rules").
		Where("bu = ?", evalCtx.BU).
		Where("service = ?", evalCtx.Service).
		Where("(environment = ? OR environment = ?)", evalCtx.Environment, EnvAll).
		Where("(action = ? OR action = ?)", evalCtx.Action, ActAll).
		OrderBy("CASE WHEN type='ALLOW' THEN 0 ELSE 1 END", "created_at DESC")

	// filter accounts
	if evalCtx.Account != "" {
		qb = qb.Where(
			"? = ANY(accounts) OR accounts IS NULL OR cardinality(accounts)=0",
			evalCtx.Account,
		)
	}

	// filter realms
	if evalCtx.RealmName != "" {
		qb = qb.Where(
			"? = ANY(realms) OR realms IS NULL OR cardinality(realms)=0",
			evalCtx.RealmName,
		)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return false, nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
	if err != nil {
		return false, nil, err
	}
	defer rows.Close()

	var rules []Rule
	for rows.Next() {
		var rule Rule
		if err := rows.Scan(
			&rule.ID,
			&rule.BU,
			&rule.Service,
			&rule.Environment,
			&rule.Action,
			&rule.Type,
			&rule.Accounts,
			&rule.Realms,
			&rule.CreatedAt,
		); err != nil {
			return false, nil, err
		}
		rules = append(rules, rule)
	}

	if err := rows.Err(); err != nil {
		return false, nil, err
	}

	// même logique métier que ton code GORM

	for _, ru := range rules {
		if ru.Type == RuleAllow {
			return true, &ru, nil
		}
	}

	if len(rules) > 0 {
		return false, &rules[0], nil
	}

	return false, nil, nil
}
```

---

# 🔥 Points importants

### 1️⃣ Placeholder `$1, $2`

On utilise :

```go
PlaceholderFormat(squirrel.Dollar)
```

Sinon pgx ne fonctionnera pas correctement.

---

### 2️⃣ OrderBy avec CASE

En GORM tu avais :

```go
Clauses(clause.OrderBy{
  Expression: clause.Expr{SQL: "CASE WHEN type='ALLOW' THEN 0 ELSE 1 END, created_at DESC"},
})
```

En Squirrel :

```go
OrderBy("CASE WHEN type='ALLOW' THEN 0 ELSE 1 END", "created_at DESC")
```

---

### 3️⃣ ANY avec pgx

`? = ANY(accounts)` fonctionne parfaitement avec pgx si :

* `accounts` est un `text[]`
* ton struct `Rule` mappe bien vers `[]string`

---

# 🎯 Version plus propre (si tu veux optimiser)

Vu que ton ORDER garantit que :

* ALLOW arrive avant DENY
* le plus récent arrive en premier

Tu pourrais juste faire :

```go
qb = qb.Limit(1)
```

Et éviter de charger toutes les rules en mémoire.

---

Si tu veux, je peux te faire :

* ✅ une version optimisée avec `LIMIT 1`
* ✅ une version avec `pgxscan` (beaucoup plus propre)
* ✅ une version repository clean architecture compatible avec ton projet REST
* ✅ une version transactionnelle avec `FOR UPDATE`

Tu veux laquelle ?

Parfait 👌
On va faire une version :

* ✅ optimisée avec `ORDER BY + LIMIT 1`
* ✅ utilisant `pgxpool`
* ✅ testable proprement
* ✅ avec **tests unitaires via pgxmock**
* ✅ architecture compatible repository propre

---

# 🎯 Idée d’optimisation

Ton tri garantit déjà :

```sql
ORDER BY 
  CASE WHEN type='ALLOW' THEN 0 ELSE 1 END,
  created_at DESC
```

Donc :

* ALLOW prioritaire
* le plus récent en premier

👉 Il suffit donc de `LIMIT 1`
👉 Plus besoin de charger toutes les rules en mémoire
👉 Complexité réduite
👉 Moins de RAM
👉 Plus performant

---

# 📦 Repository

## Interface DB (pour mock)

On définit une petite interface pour pouvoir mocker pgx :

```go
type DB interface {
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}
```

`*pgxpool.Pool` implémente déjà ça 👍

---

## Implémentation Evaluate optimisée

```go
import (
	"context"
	"errors"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

type Repo struct {
	db DB
}

func (r *Repo) Evaluate(ctx context.Context, evalCtx EvalContext) (bool, *Rule, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	qb := psql.
		Select("*").
		From("rules").
		Where("bu = ?", evalCtx.BU).
		Where("service = ?", evalCtx.Service).
		Where("(environment = ? OR environment = ?)", evalCtx.Environment, EnvAll).
		Where("(action = ? OR action = ?)", evalCtx.Action, ActAll).
		OrderBy("CASE WHEN type='ALLOW' THEN 0 ELSE 1 END").
		OrderBy("created_at DESC").
		Limit(1)

	if evalCtx.Account != "" {
		qb = qb.Where(
			"? = ANY(accounts) OR accounts IS NULL OR cardinality(accounts)=0",
			evalCtx.Account,
		)
	}

	if evalCtx.RealmName != "" {
		qb = qb.Where(
			"? = ANY(realms) OR realms IS NULL OR cardinality(realms)=0",
			evalCtx.RealmName,
		)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return false, nil, err
	}

	var rule Rule
	err = r.db.QueryRow(ctx, sql, args...).Scan(
		&rule.ID,
		&rule.BU,
		&rule.Service,
		&rule.Environment,
		&rule.Action,
		&rule.Type,
		&rule.Accounts,
		&rule.Realms,
		&rule.CreatedAt,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil, nil
	}
	if err != nil {
		return false, nil, err
	}

	return rule.Type == RuleAllow, &rule, nil
}
```

---

# 🧪 Tests unitaires avec pgxmock

Dépendance :

```bash
go get github.com/pashagolub/pgxmock/v3
```

---

## Test complet

```go
import (
	"context"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func TestEvaluate_Allow(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := &Repo{db: mock}

	evalCtx := EvalContext{
		BU:          "finance",
		Service:     "billing",
		Environment: "prod",
		Action:      "read",
	}

	rows := pgxmock.NewRows([]string{
		"id", "bu", "service", "environment", "action",
		"type", "accounts", "realms", "created_at",
	}).AddRow(
		1,
		"finance",
		"billing",
		"prod",
		"read",
		"ALLOW",
		[]string{},
		[]string{},
		time.Now(),
	)

	mock.ExpectQuery("SELECT (.+) FROM rules").
		WillReturnRows(rows)

	allowed, rule, err := repo.Evaluate(context.Background(), evalCtx)

	assert.NoError(t, err)
	assert.NotNil(t, rule)
	assert.True(t, allowed)
	assert.Equal(t, "ALLOW", rule.Type)

	assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

## Test : aucune règle

```go
func TestEvaluate_NoRule(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := &Repo{db: mock}

	mock.ExpectQuery("SELECT (.+) FROM rules").
		WillReturnError(pgx.ErrNoRows)

	allowed, rule, err := repo.Evaluate(context.Background(), EvalContext{})

	assert.NoError(t, err)
	assert.Nil(t, rule)
	assert.False(t, allowed)
}
```

---

# 🚀 Ce que tu gagnes

| Avant                          | Maintenant          |
| ------------------------------ | ------------------- |
| Chargement de toutes les rules | 1 seule rule        |
| Boucle en mémoire              | Plus besoin         |
| Plus lent si table large       | Optimisé            |
| Plus difficile à tester        | Mockable proprement |

---

# 💡 Encore mieux ?

Vu ton niveau (repo générique, archi clean…), je peux te proposer :

* 🔥 Version avec `pgxscan` pour éviter le Scan manuel
* 🔥 Version avec transaction + `FOR UPDATE`
* 🔥 Version benchmarkée
* 🔥 Version avec index SQL optimisé pour ce tri spécifique

Tu veux qu’on pousse jusqu’au niveau “production grade enterprise” ? 😄

######## tests with func generic #############

Parfait 😎
On va faire une vraie version **ultra-générique, type-safe avec Go generics**, compatible avec ton `Rule`, et testable proprement.

Objectifs :

* ✅ Générique `Repository[T]`
* ✅ Champs autorisés (sécurité anti-injection)
* ✅ QueryOptions typée
* ✅ Compatible pgx + squirrel
* ✅ Tests unitaires sur `Rule`

---

# 🧠 1️⃣ Principe type-safe

Le problème classique d’un query builder générique :

```go
Where("created_at = ?", ...)
```

👉 string libre = risque d’injection + erreurs runtime.

Solution :

* Définir une **whitelist de colonnes autorisées**
* Mapper proprement struct ↔ table
* Repository générique `BaseRepository[T]`

---

# 🏗️ 2️⃣ Définition générique

## Query Options

```go
type Filter struct {
	Field string
	Op    string
	Value any
}

type Sort struct {
	Field string
	Desc  bool
}

type QueryOptions struct {
	Filters []Filter
	OrderBy []Sort
	Limit   uint64
	Offset  uint64
}
```

---

# 🏛️ 3️⃣ Repository Générique

```go
type DBTX interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type BaseRepository[T any] struct {
	db            DBTX
	table         string
	columns       []string
	allowedFields map[string]struct{}
	scanFunc      func(pgx.Rows) (T, error)
}
```

---

# 🏗️ 4️⃣ Constructeur générique

```go
func NewBaseRepository[T any](
	db DBTX,
	table string,
	columns []string,
	scanFunc func(pgx.Rows) (T, error),
) *BaseRepository[T] {

	allowed := make(map[string]struct{})
	for _, c := range columns {
		allowed[c] = struct{}{}
	}

	return &BaseRepository[T]{
		db:            db,
		table:         table,
		columns:       columns,
		allowedFields: allowed,
		scanFunc:      scanFunc,
	}
}
```

---

# 🚀 5️⃣ Find générique type-safe

```go
func (r *BaseRepository[T]) Find(
	ctx context.Context,
	opts QueryOptions,
) ([]T, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	qb := psql.
		Select(r.columns...).
		From(r.table)

	// WHERE
	for _, f := range opts.Filters {

		if _, ok := r.allowedFields[f.Field]; !ok {
			return nil, fmt.Errorf("field not allowed: %s", f.Field)
		}

		qb = qb.Where(
			squirrel.Expr(f.Field+" "+f.Op+" ?", f.Value),
		)
	}

	// ORDER BY
	for _, s := range opts.OrderBy {

		if _, ok := r.allowedFields[s.Field]; !ok {
			return nil, fmt.Errorf("field not allowed: %s", s.Field)
		}

		order := s.Field
		if s.Desc {
			order += " DESC"
		}

		qb = qb.OrderBy(order)
	}

	if opts.Limit > 0 {
		qb = qb.Limit(opts.Limit)
	}

	if opts.Offset > 0 {
		qb = qb.Offset(opts.Offset)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sql, args...)
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

# 🧩 6️⃣ Implémentation pour Rule

## Scan Rule

```go
func ScanRule(rows pgx.Rows) (Rule, error) {
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
		&r.UpdatedBy,
		&r.UpdatedAt,
	)
	return r, err
}
```

---

## Instanciation Repo Rule

```go
ruleRepo := NewBaseRepository[Rule](
	db,
	"rules",
	[]string{
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
		"updated_by",
		"updated_at",
	},
	ScanRule,
)
```

---

# 🧪 7️⃣ Tests Unitaires sur Rule

On utilise `pgxmock`.

```go
func TestRuleRepository_Find_WithLimit(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewBaseRepository[Rule](
		mock,
		"rules",
		[]string{
			"id","type","bu","service","environment","action",
			"realms","accounts","message","created_by",
			"created_at","updated_by","updated_at",
		},
		ScanRule,
	)

	id := uuid.New()

	rows := pgxmock.NewRows([]string{
		"id","type","bu","service","environment","action",
		"realms","accounts","message","created_by",
		"created_at","updated_by","updated_at",
	}).AddRow(
		id,"ALLOW","finance","billing","prod","read",
		[]string{},[]string{},"msg","admin",
		time.Now(),nil,nil,
	)

	mock.ExpectQuery("SELECT (.+) FROM rules").
		WillReturnRows(rows)

	result, err := repo.Find(context.Background(), QueryOptions{
		Filters: []Filter{
			{Field: "bu", Op: "=", Value: "finance"},
		},
		Limit: 1,
	})

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, id, result[0].ID)
}
```

---

# 🧪 Test sécurité (champ interdit)

```go
func TestRuleRepository_InvalidField(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewBaseRepository[Rule](
		mock,
		"rules",
		[]string{"id","type"},
		ScanRule,
	)

	_, err := repo.Find(context.Background(), QueryOptions{
		Filters: []Filter{
			{Field: "DROP TABLE", Op: "=", Value: "x"},
		},
	})

	assert.Error(t, err)
}
```

---

# 🎯 Résultat

Tu as maintenant :

* 🔥 Repository générique type-safe
* 🔥 Protection contre injection SQL
* 🔥 Compatible pgx
* 🔥 Compatible Clean Architecture
* 🔥 Testable
* 🔥 Compatible pagination REST

---

# 🚀 Si tu veux aller encore plus loin

Je peux te faire :

* ✅ version avec `FindOne`
* ✅ version avec `Count`
* ✅ version avec `Insert/Update/Delete` générique
* ✅ version avec auto-mapping via reflection (plus besoin de scanFunc)
* ✅ version compatible RAG/log filtering (dans ton autre projet)

Tu veux qu’on passe en mode **repository générique production-grade complet** ? 😄

Parfait 👌
On va ajouter un **`Count` générique type-safe** à ton `BaseRepository[T]`, réutilisant exactement la même logique de filtres que `Find`.

Objectifs :

* ✅ Même `QueryOptions`
* ✅ Même validation des champs autorisés
* ✅ Requête `SELECT COUNT(*)`
* ✅ Compatible pgx
* ✅ Test unitaire sur `Rule`

---

# 🧠 1️⃣ Méthode `Count` générique

On réutilise le même builder que `Find`, mais :

* on fait `Select("COUNT(*)")`
* on ignore `OrderBy`
* on ignore `Limit/Offset`

---

## 🔹 Implémentation

```go
func (r *BaseRepository[T]) Count(
	ctx context.Context,
	opts QueryOptions,
) (int64, error) {

	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

	qb := psql.
		Select("COUNT(*)").
		From(r.table)

	// WHERE
	for _, f := range opts.Filters {

		if _, ok := r.allowedFields[f.Field]; !ok {
			return 0, fmt.Errorf("field not allowed: %s", f.Field)
		}

		qb = qb.Where(
			squirrel.Expr(f.Field+" "+f.Op+" ?", f.Value),
		)
	}

	sql, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	err = r.db.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}
```

---

# 🔥 Pourquoi on ignore OrderBy / Limit ?

Parce que :

```sql
SELECT COUNT(*) FROM ...
```

* ORDER BY inutile
* LIMIT fausserait le résultat

Donc on garde uniquement les filtres.

---

# 🧪 2️⃣ Test unitaire sur Rule

```go
func TestRuleRepository_Count(t *testing.T) {
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	repo := NewBaseRepository[Rule](
		mock,
		"rules",
		[]string{
			"id","type","bu","service","environment","action",
			"realms","accounts","message","created_by",
			"created_at","updated_by","updated_at",
		},
		ScanRule,
	)

	mock.ExpectQuery("SELECT COUNT\\(\\*\\) FROM rules").
		WithArgs("finance").
		WillReturnRows(
			pgxmock.NewRows([]string{"count"}).
				AddRow(int64(3)),
		)

	count, err := repo.Count(context.Background(), QueryOptions{
		Filters: []Filter{
			{Field: "bu", Op: "=", Value: "finance"},
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, int64(3), count)
	assert.NoError(t, mock.ExpectationsWereMet())
}
```

---

# 🧪 Test champ interdit

```go
func TestRuleRepository_Count_InvalidField(t *testing.T) {
	mock, _ := pgxmock.NewPool()
	defer mock.Close()

	repo := NewBaseRepository[Rule](
		mock,
		"rules",
		[]string{"id","type"},
		ScanRule,
	)

	_, err := repo.Count(context.Background(), QueryOptions{
		Filters: []Filter{
			{Field: "DROP TABLE", Op: "=", Value: "x"},
		},
	})

	assert.Error(t, err)
}
```

---

# 🚀 Maintenant ton repository supporte

| Méthode           | Fonction              |
| ----------------- | --------------------- |
| `Find`            | Liste paginée         |
| `Count`           | Total pour pagination |
| Validation champs | Anti injection        |
| Generics          | Type-safe             |
| pgx               | Production ready      |

---

# 💎 Bonus : Pagination REST complète

Maintenant tu peux faire :

```go
items, _ := repo.Find(ctx, opts)
total, _ := repo.Count(ctx, opts)
```

Et retourner :

```json
{
  "data": [...],
  "total": 42,
  "limit": 10,
  "offset": 0
}
```

---

Si tu veux, prochaine étape possible :

* 🔥 méthode `FindWithPagination()` qui retourne `{Data, Total}`
* 🔥 support automatique des opérateurs `IN`, `LIKE`, `ILIKE`
* 🔥 support JSONB filters
* 🔥 auto-mapping struct tags → colonnes SQL (full generic reflection)

Tu veux qu’on le rende encore plus propre ? 😄

## go conditon query
```go
func buildRuleQuery(filter RuleFilter) (string, []any, error) {

	CleanStringArray((*pq.StringArray)(&filter.Accounts))
	CleanStringArray((*pq.StringArray)(&filter.Realms))

	// fmt.Printf("RuleFilter: type: %v, Account: %v, Action: %v, BU: %v, Environment: %v, Realms: %v, Service: %v", filter.Type, filter.Accounts, filter.Action, filter.BU, filter.Environment, filter.Realms, filter.Service)
	builder := sq.
	Select(
	"id", "type", "bu", "service", "environment", "action",
	"realms", "accounts", "message",
	"created_by", "created_at", "updated_by", "updated_at",
	).
	From("rules").
	PlaceholderFormat(sq.Dollar)

	// BU
	if filter.BU != "" {
		builder = builder.Where(sq.Eq{"bu": filter.BU})
	}

	// Service
	if filter.Service != "" {
		builder = builder.Where(sq.Eq{"service": filter.Service})
	}

	// type
	if filter.Type != "" {
		builder = builder.Where(sq.Eq{"type": filter.Type})
	}

	// Environment (vide => ALL)
	if filter.Environment != "" {
		builder = builder.Where(
		sq.Or{
			sq.Eq{"environment": filter.Environment},
			sq.Eq{"environment": EnvALL},
		})
	}

	// Action (vide => ALL)
	if filter.Action != "" {
		builder = builder.Where(
		sq.Or{
			sq.Eq{"action": filter.Action},
			sq.Eq{"action": ActAll},
		})
	}

	// Realms (@>)
	if len(filter.Realms) > 0 {
		builder = builder.Where(
			sq.Expr(
			"? = ANY(realms) OR cardinality(realms) = 0",
			filter.Realms,
		),
	)
}

	// Accounts (&&)
	if len(filter.Accounts) > 0 {
		builder = builder.Where(
			sq.Expr(
				"? = ANY(accounts) OR cardinality(accounts) = 0",
				filter.Accounts,
			),
		)
	}

	builder = builder.OrderBy("created_at DESC")

	return builder.ToSql()
}

func CleanStringArray(a *pq.StringArray) {
	cleaned := (*a)[:0]
	for _, v := range *a {
		if strings.TrimSpace(v) != "" {
			cleaned = append(cleaned, v)
		}
	}
	*a = cleaned
}
```