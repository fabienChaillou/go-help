Parfait 👌 on va faire ça propre et dynamique avec **Squirrel**.

On veut générer quelque chose comme :

```sql
SELECT 
  (type = 'ALLOW') AS is_allow,
  message
FROM rules
WHERE bu = $1
  AND service = $2
  AND environment = $3
  AND action = $4
ORDER BY (type = 'ALLOW') DESC,
         created_at DESC
LIMIT 1;
```

Avec :

Masterminds/squirrel
et compatible PostgreSQL

---

# ✅ Fonction complète avec filtres

```go
package repository

import (
	sq "github.com/Masterminds/squirrel"
)

type RuleFilter struct {
	BU          string
	Service     string
	Environment string
	Action      string
}

func BuildSelectPrioritizedRule(filter RuleFilter) (string, []interface{}, error) {
	builder := sq.
		Select(
			"(type = 'ALLOW') AS is_allow",
			"message",
		).
		From("rules").
		Where(sq.Eq{
			"bu":          filter.BU,
			"service":     filter.Service,
			"environment": filter.Environment,
			"action":      filter.Action,
		}).
		OrderBy("(type = 'ALLOW') DESC").
		OrderBy("created_at DESC").
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	return builder.ToSql()
}
```

---

# 🧪 Exemple d'utilisation

```go
filter := RuleFilter{
	BU:          "finance",
	Service:     "billing",
	Environment: "prod",
	Action:      "read",
}

sql, args, err := BuildSelectPrioritizedRule(filter)
if err != nil {
	log.Fatal(err)
}

row := db.QueryRowContext(ctx, sql, args...)

var isAllow bool
var message string

err = row.Scan(&isAllow, &message)
if err != nil {
	// gérer sql.ErrNoRows
}

if isAllow {
	return message
}
return false
```

---

# 🧠 Important

Si une des colonnes peut être optionnelle, il vaut mieux construire le `WHERE` dynamiquement :

```go
builder := sq.
	Select("(type = 'ALLOW') AS is_allow", "message").
	From("rules")

if filter.BU != "" {
	builder = builder.Where(sq.Eq{"bu": filter.BU})
}
if filter.Service != "" {
	builder = builder.Where(sq.Eq{"service": filter.Service})
}
if filter.Environment != "" {
	builder = builder.Where(sq.Eq{"environment": filter.Environment})
}
if filter.Action != "" {
	builder = builder.Where(sq.Eq{"action": filter.Action})
}

builder = builder.
	OrderBy("(type = 'ALLOW') DESC").
	OrderBy("created_at DESC").
	Limit(1).
	PlaceholderFormat(sq.Dollar)
```

---

# 🚀 Bonus performance

Pour que cette requête soit rapide sur une grosse table, index recommandé :

```sql
CREATE INDEX idx_rules_lookup
ON rules (bu, service, environment, action, created_at DESC);
```

---

Si tu veux, je peux maintenant :

* 🔥 te faire la version repository complète (interface + impl + test mockgen)
* 🔥 te proposer un index partiel optimisé pour ALLOW
* 🔥 analyser le plan d’exécution PostgreSQL avec EXPLAIN

Oui 👌 je te fais un test propre avec `sqlmock`, adapté à ton architecture repository en Go.

On va tester :

* ✅ que la requête est bien générée
* ✅ qu’elle retourne `is_allow` et `message`
* ✅ le comportement métier (retour message si ALLOW)

On utilise :

DATA-DOG/go-sqlmock
Masterminds/squirrel

---

# 🎯 Exemple de méthode repository à tester

Imaginons :

```go
func (r *Repository) GetPrioritizedRule(ctx context.Context, filter RuleFilter) (any, error) {
	sql, args, err := BuildSelectPrioritizedRule(filter)
	if err != nil {
		return nil, err
	}

	row := r.db.QueryRowContext(ctx, sql, args...)

	var isAllow bool
	var message string

	if err := row.Scan(&isAllow, &message); err != nil {
		return nil, err
	}

	if isAllow {
		return message, nil
	}

	return false, nil
}
```

---

# ✅ Test avec sqlmock

```go
package repository_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
)

func TestGetPrioritizedRule_Allow(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	filter := RuleFilter{
		BU:          "finance",
		Service:     "billing",
		Environment: "prod",
		Action:      "read",
	}

	// Regex car Squirrel génère du SQL formaté
	query := regexp.QuoteMeta(`
SELECT (type = 'ALLOW') AS is_allow, message
FROM rules
WHERE bu = $1 AND service = $2 AND environment = $3 AND action = $4
ORDER BY (type = 'ALLOW') DESC, created_at DESC
LIMIT 1`)

	rows := sqlmock.NewRows([]string{"is_allow", "message"}).
		AddRow(true, "access granted")

	mock.ExpectQuery(query).
		WithArgs("finance", "billing", "prod", "read").
		WillReturnRows(rows)

	result, err := repo.GetPrioritizedRule(context.Background(), filter)

	require.NoError(t, err)
	require.Equal(t, "access granted", result)
	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

# ✅ Test quand ce n’est PAS ALLOW

```go
func TestGetPrioritizedRule_Deny(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	repo := NewRepository(db)

	filter := RuleFilter{
		BU:          "finance",
		Service:     "billing",
		Environment: "prod",
		Action:      "read",
	}

	query := regexp.QuoteMeta(`
SELECT (type = 'ALLOW') AS is_allow, message
FROM rules
WHERE bu = $1 AND service = $2 AND environment = $3 AND action = $4
ORDER BY (type = 'ALLOW') DESC, created_at DESC
LIMIT 1`)

	rows := sqlmock.NewRows([]string{"is_allow", "message"}).
		AddRow(false, "denied")

	mock.ExpectQuery(query).
		WithArgs("finance", "billing", "prod", "read").
		WillReturnRows(rows)

	result, err := repo.GetPrioritizedRule(context.Background(), filter)

	require.NoError(t, err)
	require.Equal(t, false, result)
	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

# 🧠 Pourquoi c’est un bon test ?

✔ Teste la vraie query générée
✔ Vérifie les args
✔ Vérifie le mapping Scan
✔ Vérifie la logique métier
✔ 100% isolé (pas besoin de PostgreSQL)

---

# 🚀 Bonus (plus robuste)

Si tu veux éviter que le test casse à cause d’un changement mineur de format SQL, tu peux matcher seulement :

```go
mock.ExpectQuery("SELECT .* FROM rules")
```

---

Si tu veux, je peux te faire :

* 🔥 version avec `go:generate` + mockgen
* 🔥 test table-driven
* 🔥 test d’erreur `sql.ErrNoRows`
* 🔥 test d’intégration avec testcontainer PostgreSQL

Dis-moi ce que tu veux approfondir 👌

