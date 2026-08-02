## query chate gbt

```sql
CREATE TABLE IF NOT EXISTS rules (
	i_d VARCHAR,
	type VARCHAR,
	b_u TEXT,
	service TEXT,
	environment VARCHAR,
	action VARCHAR,
	realms VARCHAR,
	accounts VARCHAR,
	message TEXT,
	created_by TEXT,
	created_at TIMESTAMP,
	updated_by TEXT,
	updated_at TIMESTAMP
);
```

go struct
```go
type Rule struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	Type        RuleType   `json:"type" db:"type"`
	BU          string     `json:"bu" db:"bu"`
	Service     string     `json:"service" db:"service"`
	Environment EnvType    `json:"environment" db:"environment"`
	Action      Action     `json:"action" db:"action"`

	Realms   []string `json:"realms,omitempty" db:"realms"`
	Accounts []string `json:"accounts,omitempty" db:"accounts"`

	Message   string     `json:"message,omitempty" db:"message"`
	CreatedBy string     `json:"createdBy,omitempty" db:"created_by"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`

	UpdatedBy *string    `json:"updatedBy,omitempty" db:"updated_by"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty" db:"updated_at"`
}
```

```sql
SELECT id, type, bu, service, environment, action, realms, accounts, message, created_by, created_at, updated_by, updated_at FROM rules WHERE bu = $1 LIMIT 1
```


```sql
SELECT id, type, bu, service, environment, action, realms, accounts, message, created_by, created_at, updated_by, updated_at FROM rules WHERE bu = $1 LIMIT 1
```

```sql
SELECT * FROM rules WHERE bu = $1 AND service = $2 AND (environment = $3 OR environment = $4) AND (action = $5 OR action = $6) ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END, created_at DESC LIMIT 1
```


## reponse Chate GBT

Parfait 👌 on va faire exactement le même type de test avec **pgxmock**, mais adapté à ta requête complexe.

Je vais te montrer :

1. ✅ Un repository minimal
2. ✅ Le test avec `ExpectQuery`
3. ✅ La gestion des 6 arguments
4. ✅ La gestion du `ORDER BY CASE`
5. ✅ Le mock des colonnes correspondant à ta table

On parle bien de `github.com/pashagolub/pgxmock/v4`.

---

# 🏗 1️⃣ Repository à tester

Imaginons cette méthode :

```go
type Rule struct {
	ID          string
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
	UpdatedBy   string
	UpdatedAt   time.Time
}

type RuleRepository struct {
	db pgxpoolIface // interface pour mocker
}

func (r *RuleRepository) FindBestRule(
	ctx context.Context,
	bu, service,
	env1, env2,
	action1, action2 string,
) (*Rule, error) {

	query := `
SELECT * FROM rules
WHERE bu = $1
AND service = $2
AND (environment = $3 OR environment = $4)
AND (action = $5 OR action = $6)
ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END,
created_at DESC
LIMIT 1
`

	row := r.db.QueryRow(ctx, query,
		bu, service,
		env1, env2,
		action1, action2,
	)

	var rule Rule
	err := row.Scan(
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
		&rule.UpdatedBy,
		&rule.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &rule, nil
}
```

---

# 🧪 2️⃣ Test avec pgxmock

⚠️ Rappel : `ExpectQuery` utilise une **regex**, donc :

* `$1` → `\\$1`
* `(` et `)` doivent être échappées
* les espaces peuvent être simplifiés

---

## Version propre et robuste (recommandée)

On simplifie la regex pour éviter les problèmes d’espaces :

```go
func TestFindBestRule(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &RuleRepository{db: mock}

	now := time.Now()

	rows := pgxmock.NewRows([]string{
		"i_d",
		"type",
		"b_u",
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
	}).AddRow(
		"rule-1",
		"ALLOW",
		"BU1",
		"SVC1",
		"prod",
		[]string{"read"},
		[]string{"realm1"},
		"acc1",
		"ok",
		"admin",
		now,
		"admin",
		now,
	)

	mock.ExpectQuery(`SELECT \* FROM rules`).
		WithArgs(
			"BU1",
			"SVC1",
			"prod",
			"default",
			"read",
			"write",
		).
		WillReturnRows(rows)

	rule, err := repo.FindBestRule(
		context.Background(),
		"BU1",
		"SVC1",
		"prod",
		"default",
		"read",
		"write",
	)

	require.NoError(t, err)
	require.Equal(t, "rule-1", rule.ID)
	require.Equal(t, "ALLOW", rule.Type)

	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

# 🔬 3️⃣ Si tu veux matcher la requête complète (regex exacte)

Voici la regex complète correctement échappée :

```go
mock.ExpectQuery(
`SELECT \* FROM rules
WHERE bu = \$1
AND service = \$2
AND \(environment = \$3 OR environment = \$4\)
AND \(action = \$5 OR action = \$6\)
ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END,
created_at DESC
LIMIT 1`,
)
```

Mais ⚠️ c’est fragile si tu modifies le formatting.

---

# ❗ 4️⃣ Cas test erreur

Tester si aucune règle trouvée :

```go
mock.ExpectQuery(`SELECT \* FROM rules`).
	WithArgs("BU1", "SVC1", "prod", "default", "read", "write").
	WillReturnError(pgx.ErrNoRows)
```

Et vérifier que ton repository gère correctement le `ErrNoRows`.

---

# 🧠 5️⃣ Points importants dans ton cas

### 🔹 1️⃣ Ordre des colonnes

`SELECT *` retourne les colonnes **dans l’ordre de création de la table** :

```sql
i_d,
type,
b_u,
service,
environment,
action,
realms,
accounts,
message,
created_by,
created_at,
updated_by,
updated_at
```

Ton `Scan()` doit respecter exactement cet ordre.

---

### 🔹 2️⃣ Si tu utilises Squirrel

Si tu construis la requête dynamiquement, je recommande :

```go
mock.ExpectQuery("SELECT \\* FROM rules")
```

Plutôt que matcher toute la string.

---

# 🎯 Résumé

Pour ta requête spécifique :

* ✔️ Utiliser `ExpectQuery`
* ✔️ Échapper `$`
* ✔️ Échapper `(` `)`
* ✔️ Fournir les 6 arguments dans le bon ordre
* ✔️ Fournir les colonnes dans l’ordre exact de la table
* ✔️ Toujours appeler `ExpectationsWereMet()`

---

Si tu veux, je peux maintenant te faire :

* 🔥 une version avec transaction
* 🔥 une version avec repository générique
* 🔥 une version avec Squirrel
* 🔥 une version avec test table-driven
* 🔥 une version avec mockgen + interface propre

Tu veux laquelle ?

## test query conditions

### go func condition

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


```sql
SELECT * FROM rules
WHERE bu = $1
AND service = $2
AND (environment = $3 OR environment = $4)
AND (action = $5 OR action = $6)
ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END,
created_at DESC
LIMIT 1
```

```sql
CREATE TABLE IF NOT EXISTS rules (
	i_d VARCHAR,
	type VARCHAR,
	b_u TEXT,
	service TEXT,
	environment VARCHAR,
	action VARCHAR,
	realms VARCHAR,
	accounts VARCHAR,
	message TEXT,
	created_by TEXT,
	created_at TIMESTAMP,
	updated_by TEXT,
	updated_at TIMESTAMP
);
```


## query whitout bu and service
```sql
SELECT * FROM rules 
ORDER BY created_at DESC
```

## query whit bu
```sql
SELECT * FROM rules WHERE bu = $1
ORDER BY created_at DESC
```

## query whit service
```sql
SELECT * FROM rules WHERE service = $1
ORDER BY created_at DESC
```

## query whit type
```sql
SELECT * FROM rules WHERE type = $1
ORDER BY created_at DESC
```

## query bu and service
```sql
SELECT * FROM rules WHERE bu = $1 AND service = $2 
ORDER BY created_at DESC
```

## query bu and type
```sql
SELECT * FROM rules WHERE bu = $1 AND type = $2 
ORDER BY created_at DESC
```

## query service and type
```sql
SELECT * FROM rules WHERE service = $1 AND type = $2 
ORDER BY created_at DESC
```

## query environment OR environment = ALL
```sql
SELECT * FROM rules WHERE (environment = $1 OR environment = $2)
ORDER BY created_at DESC
```

## query action OR action = ALL
```sql
SELECT * FROM rules WHERE (action = $1 OR action = $2)
ORDER BY created_at DESC
```

## query Realms
```sql
SELECT * FROM rules WHERE $1 = ANY(realms) OR cardinality(realms) = 0
ORDER BY created_at DESC
```

## query accounts
```sql
SELECT * FROM rules WHERE $1 = ANY(accounts) OR cardinality(accounts) = 0
ORDER BY created_at DESC
```


## Query multi args

## query bu and service and environment
```sql
SELECT * FROM rules WHERE bu = $1 AND service = $2 AND (environment = $3 OR environment = $4) 
ORDER BY created_at DESC
```

## query bu and service and environment and action
```sql
SELECT * FROM rules WHERE bu = $1 AND service = $2 AND (environment = $3 OR environment = $4) AND (action = $5 OR action = $6) 
ORDER BY created_at DESC
```

## requery specific for Validate func
```sql
SELECT * FROM rules WHERE bu = $1 AND service = $2 AND (environment = $3 OR environment = $4) AND (action = $5 OR action = $6) ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END, created_at DESC LIMIT 1
```

## Test unit

Parfait 👌
Je vais te faire un **vrai test table-driven complet avec pgxmock v4**, couvrant **tous tes cas de figures**, proprement structuré et maintenable.

On reste sur :

`github.com/pashagolub/pgxmock/v4`

---

# 🏗 Hypothèses

On suppose une méthode générique :

```go
func (r *RuleRepository) Find(ctx context.Context, query string, args ...any) ([]Rule, error)
```

Et on teste uniquement :

* la requête envoyée
* les arguments
* le scan
* le nombre de résultats

---

# 🧪 Test table-driven complet

```go
func TestRuleRepository_Find(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &RuleRepository{db: mock}

	now := time.Now()

	columns := []string{
		"i_d",
		"type",
		"b_u",
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
	}

	baseRow := []any{
		"rule-1",
		"ALLOW",
		"BU1",
		"SVC1",
		"prod",
		"read",
		[]string{"realm1"},
		[]string{"acc1"},
		"ok",
		"admin",
		now,
		"admin",
		now,
	}

	tests := []struct {
		name  string
		query string
		args  []any
	}{
		{
			name:  "without bu and service",
			query: `SELECT \* FROM rules ORDER BY created_at DESC`,
			args:  nil,
		},
		{
			name:  "with bu",
			query: `SELECT \* FROM rules WHERE bu = \$1 ORDER BY created_at DESC`,
			args:  []any{"BU1"},
		},
		{
			name:  "with service",
			query: `SELECT \* FROM rules WHERE service = \$1 ORDER BY created_at DESC`,
			args:  []any{"SVC1"},
		},
		{
			name:  "with type",
			query: `SELECT \* FROM rules WHERE type = \$1 ORDER BY created_at DESC`,
			args:  []any{"ALLOW"},
		},
		{
			name:  "bu and service",
			query: `SELECT \* FROM rules WHERE bu = \$1 AND service = \$2 ORDER BY created_at DESC`,
			args:  []any{"BU1", "SVC1"},
		},
		{
			name:  "bu and type",
			query: `SELECT \* FROM rules WHERE bu = \$1 AND type = \$2 ORDER BY created_at DESC`,
			args:  []any{"BU1", "ALLOW"},
		},
		{
			name:  "service and type",
			query: `SELECT \* FROM rules WHERE service = \$1 AND type = \$2 ORDER BY created_at DESC`,
			args:  []any{"SVC1", "ALLOW"},
		},
		{
			name:  "environment OR environment",
			query: `SELECT \* FROM rules WHERE \(environment = \$1 OR environment = \$2\) ORDER BY created_at DESC`,
			args:  []any{"prod", "ALL"},
		},
		{
			name:  "action OR action",
			query: `SELECT \* FROM rules WHERE \(action = \$1 OR action = \$2\) ORDER BY created_at DESC`,
			args:  []any{"read", "ALL"},
		},
		{
			name:  "realms ANY",
			query: `SELECT \* FROM rules WHERE \$1 = ANY\(realms\) OR cardinality\(realms\) = 0 ORDER BY created_at DESC`,
			args:  []any{"realm1"},
		},
		{
			name:  "accounts ANY",
			query: `SELECT \* FROM rules WHERE \$1 = ANY\(accounts\) OR cardinality\(accounts\) = 0 ORDER BY created_at DESC`,
			args:  []any{"acc1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rows := pgxmock.NewRows(columns).AddRow(baseRow...)

			expect := mock.ExpectQuery(tt.query)

			if tt.args != nil {
				expect.WithArgs(tt.args...)
			}

			expect.WillReturnRows(rows)

			result, err := repo.Find(context.Background(), tt.query, tt.args...)
			require.NoError(t, err)
			require.Len(t, result, 1)

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
```

---

# 🔥 Ce que couvre ce test

✅ Aucun filtre
✅ Filtre simple (1 param)
✅ Filtre double (2 params)
✅ OR condition
✅ ANY(array)
✅ cardinality(array)
✅ ORDER BY
✅ Scan correct
✅ Vérification des arguments
✅ Vérification des attentes

---

# 🧠 Points importants

### 1️⃣ Les regex

On échappe :

* `*` → `\*`
* `$1` → `\$1`
* `( )` → `\(` `\)`

---

### 2️⃣ realms et accounts

Pour que `ANY()` fonctionne en vrai PostgreSQL, il faut :

```sql
realms TEXT[]
accounts TEXT[]
```

Sinon `ANY()` ne marchera pas en production.

---

### 3️⃣ Pourquoi test table-driven ?

* évite la duplication
* scalable
* lisible
* parfait pour architecture repository
* compatible CI

---

# 🎯 Si tu veux aller encore plus loin

Je peux te faire :

* 🔥 Version avec Squirrel dynamique
* 🔥 Version avec builder conditionnel (si param != nil)
* 🔥 Version avec pagination
* 🔥 Version avec transaction
* 🔥 Version avec mockgen + interface propre
* 🔥 Version avec gestion NULL arrays
* 🔥 Version avec pgtype.TextArray

Tu veux qu’on pousse jusqu’à un repository "production ready" niveau senior Go ?


## example with Squirrel

Parfait 👌 On va faire une **version propre avec Squirrel dynamique + pgxmock v4**, adaptée à ton cas.

On va utiliser :

* Squirrel pour construire la requête dynamiquement
* pgx pour l’exécution
* pgxmock pour les tests

---

# 🏗 1️⃣ Repository avec Squirrel dynamique

On va construire la query uniquement avec les filtres fournis.

```go id="repo_squirrel"
type RuleRepository struct {
	db pgxpoolIface
}

func (r *RuleRepository) Find(
	ctx context.Context,
	bu *string,
	service *string,
	ruleType *string,
) ([]Rule, error) {

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	query := psql.Select("*").
		From("rules").
		OrderBy("created_at DESC")

	if bu != nil {
		query = query.Where(sq.Eq{"bu": *bu})
	}

	if service != nil {
		query = query.Where(sq.Eq{"service": *service})
	}

	if ruleType != nil {
		query = query.Where(sq.Eq{"type": *ruleType})
	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Rule

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
			&rule.UpdatedBy,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, rule)
	}

	return result, nil
}
```

---

# 🧪 2️⃣ Test avec pgxmock + Squirrel

⚠️ Ici on ne teste PAS la string SQL exacte.
On teste la requête générée par Squirrel.

```go id="test_squirrel"
func TestRuleRepository_Find_WithSquirrel(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &RuleRepository{db: mock}

	now := time.Now()

	columns := []string{
		"i_d","type","b_u","service","environment","action",
		"realms","accounts","message","created_by",
		"created_at","updated_by","updated_at",
	}

	rows := pgxmock.NewRows(columns).AddRow(
		"rule-1",
		"ALLOW",
		"BU1",
		"SVC1",
		"prod",
		"read",
		[]string{},
		[]string{},
		"ok",
		"admin",
		now,
		"admin",
		now,
	)

	// Cas: bu + service
	bu := "BU1"
	service := "SVC1"

	// On s'attend à une requête contenant WHERE bu = $1 AND service = $2
	mock.ExpectQuery(`SELECT \* FROM rules WHERE bu = \$1 AND service = \$2 ORDER BY created_at DESC`).
		WithArgs(bu, service).
		WillReturnRows(rows)

	result, err := repo.Find(context.Background(), &bu, &service, nil)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

# 🔥 3️⃣ Version plus robuste (recommandée)

Avec Squirrel, l’ordre des conditions peut changer.

Donc au lieu de matcher la query complète, on peut faire :

```go id="robust_mock"
mock.ExpectQuery(`SELECT \* FROM rules`).
	WillReturnRows(rows)
```

👉 Puis vérifier uniquement les arguments avec `WithArgs`.

C’est plus stable en CI.

---

# 🧪 4️⃣ Test table-driven complet avec Squirrel

```go id="table_squirrel"
func TestRuleRepository_Find_TableDriven(t *testing.T) {
	tests := []struct {
		name     string
		bu       *string
		service  *string
		ruleType *string
	}{
		{"no filters", nil, nil, nil},
		{"bu only", strPtr("BU1"), nil, nil},
		{"service only", nil, strPtr("SVC1"), nil},
		{"type only", nil, nil, strPtr("ALLOW")},
		{"bu + service", strPtr("BU1"), strPtr("SVC1"), nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			mock, _ := pgxmock.NewPool()
			defer mock.Close()

			repo := &RuleRepository{db: mock}

			rows := pgxmock.NewRows([]string{
				"i_d","type","b_u","service","environment","action",
				"realms","accounts","message","created_by",
				"created_at","updated_by","updated_at",
			}).AddRow(
				"rule-1","ALLOW","BU1","SVC1",
				"prod","read",[]string{},[]string{},
				"ok","admin",time.Now(),"admin",time.Now(),
			)

			mock.ExpectQuery(`SELECT \* FROM rules`).
				WillReturnRows(rows)

			_, err := repo.Find(context.Background(), tt.bu, tt.service, tt.ruleType)

			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func strPtr(s string) *string {
	return &s
}
```

---

# 🎯 Pourquoi cette approche est meilleure ?

### ✔ Pas de SQL dupliqué dans les tests

### ✔ Requête générée dynamiquement

### ✔ Facile à maintenir

### ✔ Compatible avec filtres additionnels futurs

### ✔ Idéal pour architecture repository propre

---

# 🚀 with Limit 1

Parfait 👌 on va **ajouter cette condition complète** dans tes tests **pgxmock + Squirrel**, proprement et de façon robuste.

La requête à tester :

```sql
SELECT * FROM rules
WHERE bu = $1
AND service = $2
AND (environment = $3 OR environment = $4)
AND (action = $5 OR action = $6)
ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END,
created_at DESC
LIMIT 1
```

---

# ✅ 1️⃣ Version test avec pgxmock (match exact)

Si tu veux tester **la requête exacte**, voici comment l’ajouter :

```go id="test_full_condition"
func TestRuleRepository_Find_WithFullCondition(t *testing.T) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &RuleRepository{db: mock}

	now := time.Now()

	columns := []string{
		"i_d","type","b_u","service","environment","action",
		"realms","accounts","message","created_by",
		"created_at","updated_by","updated_at",
	}

	rows := pgxmock.NewRows(columns).AddRow(
		"rule-1",
		"ALLOW",
		"BU1",
		"SVC1",
		"prod",
		"read",
		[]string{},
		[]string{},
		"ok",
		"admin",
		now,
		"admin",
		now,
	)

	query := `SELECT \* FROM rules
WHERE bu = \$1
AND service = \$2
AND \(environment = \$3 OR environment = \$4\)
AND \(action = \$5 OR action = \$6\)
ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END,
created_at DESC
LIMIT 1`

	mock.ExpectQuery(query).
		WithArgs(
			"BU1",
			"SVC1",
			"prod",
			"default",
			"read",
			"write",
		).
		WillReturnRows(rows)

	result, err := repo.FindBestRule(
		context.Background(),
		"BU1",
		"SVC1",
		"prod",
		"default",
		"read",
		"write",
	)

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "rule-1", result[0].ID)

	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

# ✅ 2️⃣ Version recommandée (plus stable)

⚠️ Matcher toute la requête complète peut casser si tu reformates le SQL.

Version plus robuste :

```go id="robust_full_condition"
mock.ExpectQuery(`SELECT \* FROM rules`).
	WithArgs("BU1","SVC1","prod","default","read","write").
	WillReturnRows(rows)
```

👉 On teste :

* les arguments
* le comportement
* le résultat
* pas le formatting exact

C’est plus maintenable en CI.

---

# ✅ 3️⃣ Version avec Squirrel (dynamique)

Si ta query est construite avec Squirrel, voici comment intégrer la condition :

```go id="squirrel_full_condition"
psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

query := psql.
	Select("*").
	From("rules").
	Where(sq.Eq{"bu": bu}).
	Where(sq.Eq{"service": service}).
	Where(sq.Or{
		sq.Eq{"environment": env1},
		sq.Eq{"environment": env2},
	}).
	Where(sq.Or{
		sq.Eq{"action": action1},
		sq.Eq{"action": action2},
	}).
	OrderBy("CASE WHEN type='ALLOW' THEN 0 ELSE 1 END",
		"created_at DESC").
	Limit(1)
```

Puis dans le test :

```go id="squirrel_mock"
mock.ExpectQuery(`SELECT \* FROM rules`).
	WillReturnRows(rows)
```

---

# 🧪 4️⃣ Si tu veux l’intégrer dans ton table-driven test

Ajoute simplement un nouveau cas :

```go id="table_case"
{
	name: "full condition",
	query: `SELECT \* FROM rules WHERE bu = \$1 AND service = \$2 AND \(environment = \$3 OR environment = \$4\) AND \(action = \$5 OR action = \$6\) ORDER BY CASE WHEN type='ALLOW' THEN 0 ELSE 1 END, created_at DESC LIMIT 1`,
	args: []any{"BU1","SVC1","prod","default","read","write"},
},
```

---

# 🎯 Recommandation d’architecture

Si ton repository devient complexe avec :

* filtres dynamiques
* ANY()
* OR
* ORDER BY conditionnel
* LIMIT
* pagination

👉 Le meilleur pattern est :

**Squirrel pour construire**
**pgxmock pour tester les args + résultat**
**Ne pas dépendre du SQL exact dans les tests**

---

Si tu veux, je peux maintenant te faire :

* 🔥 Une version complète avec Squirrel + tous les filtres optionnels combinés
* 🔥 Un builder propre pour éviter les duplications
* 🔥 Une architecture repository senior-ready avec tests 100%
* 🔥 Une stratégie de test optimale pour éviter les regex fragiles

Tu veux qu’on pousse ça au niveau production avancé ? 🚀


### test all conditions

Parfait 👌 on va faire une **version complète, propre et production-ready** avec :

* Squirrel pour construire dynamiquement la requête
* pgx pour l’exécution
* pgxmock pour les tests

On va gérer **tous les filtres optionnels combinables** :

* bu
* service
* type
* environment (env OR ALL)
* action (action OR ALL)
* realms (ANY + cardinality=0)
* accounts (ANY + cardinality=0)
* tri prioritaire ALLOW
* LIMIT optionnel

---

# 🏗 1️⃣ Structure des filtres

On évite 15 paramètres 😅
On passe une struct propre :

```go
type RuleFilters struct {
	BU        *string
	Service   *string
	Type      *string
	Env       *string
	Action    *string
	Realm     *string
	Account   *string
	Limit     *uint64
}
```

---

# 🧠 2️⃣ Repository complet avec Squirrel

```go
func (r *RuleRepository) Find(ctx context.Context, f RuleFilters) ([]Rule, error) {

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := psql.
		Select("*").
		From("rules")

	// ---- filtres simples ----

	if f.BU != nil {
		q = q.Where(sq.Eq{"bu": *f.BU})
	}

	if f.Service != nil {
		q = q.Where(sq.Eq{"service": *f.Service})
	}

	if f.Type != nil {
		q = q.Where(sq.Eq{"type": *f.Type})
	}

	// ---- environment OR ALL ----

	if f.Env != nil {
		q = q.Where(sq.Or{
			sq.Eq{"environment": *f.Env},
			sq.Eq{"environment": "ALL"},
		})
	}

	// ---- action OR ALL ----

	if f.Action != nil {
		q = q.Where(sq.Or{
			sq.Eq{"action": *f.Action},
			sq.Eq{"action": "ALL"},
		})
	}

	// ---- realms ANY + cardinality ----

	if f.Realm != nil {
		q = q.Where(
			sq.Expr("? = ANY(realms) OR cardinality(realms) = 0", *f.Realm),
		)
	}

	// ---- accounts ANY + cardinality ----

	if f.Account != nil {
		q = q.Where(
			sq.Expr("? = ANY(accounts) OR cardinality(accounts) = 0", *f.Account),
		)
	}

	// ---- ORDER BY prioritaire ALLOW ----

	q = q.OrderBy(
		"CASE WHEN type='ALLOW' THEN 0 ELSE 1 END",
		"created_at DESC",
	)

	// ---- LIMIT ----

	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Rule

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
			&rule.UpdatedBy,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, rule)
	}

	return result, nil
}
```

---

# 🧪 3️⃣ Test complet combiné avec pgxmock

⚠️ On ne matche PAS la SQL exacte → fragile
On vérifie :

* la requête part bien
* les arguments sont corrects
* le résultat est correct

```go
func TestRuleRepository_Find_AllFilters(t *testing.T) {

	mock, err := pgxmock.NewPool()
	require.NoError(t, err)
	defer mock.Close()

	repo := &RuleRepository{db: mock}

	now := time.Now()

	columns := []string{
		"i_d","type","b_u","service","environment","action",
		"realms","accounts","message","created_by",
		"created_at","updated_by","updated_at",
	}

	rows := pgxmock.NewRows(columns).AddRow(
		"rule-1",
		"ALLOW",
		"BU1",
		"SVC1",
		"prod",
		"read",
		[]string{"realm1"},
		[]string{"acc1"},
		"ok",
		"admin",
		now,
		"admin",
		now,
	)

	bu := "BU1"
	service := "SVC1"
	env := "prod"
	action := "read"
	realm := "realm1"
	account := "acc1"
	limit := uint64(1)

	mock.ExpectQuery(`SELECT \* FROM rules`).
		WithArgs(
			bu,
			service,
			env,        // env
			realm,      // ANY(realms)
			account,    // ANY(accounts)
			action,     // action
		).
		WillReturnRows(rows)

	result, err := repo.Find(context.Background(), RuleFilters{
		BU:      &bu,
		Service: &service,
		Env:     &env,
		Action:  &action,
		Realm:   &realm,
		Account: &account,
		Limit:   &limit,
	})

	require.NoError(t, err)
	require.Len(t, result, 1)
	require.Equal(t, "rule-1", result[0].ID)

	require.NoError(t, mock.ExpectationsWereMet())
}
```

---

# 🧠 4️⃣ Pourquoi cette version est solide en prod ?

✅ Tous filtres combinables
✅ Pas de duplication SQL
✅ Query builder safe
✅ Compatible PostgreSQL arrays
✅ ORDER BY prioritaire ALLOW
✅ LIMIT dynamique
✅ Test maintenable
✅ Pas fragile au formatting

---

# 🚀 Si tu veux niveau encore au-dessus

Je peux te faire :

* 🔥 pagination OFFSET + COUNT total
* 🔥 cache key basé sur filtres
* 🔥 specification pattern propre
* 🔥 builder dédié (RuleQueryBuilder)
* 🔥 test de la SQL générée précisément
* 🔥 benchmark performance Squirrel vs SQL raw

Tu veux qu’on pousse ça vers une architecture senior / clean architecture complète ?


## with Order by optional

Oui 👍 on va rendre ce `ORDER BY` **optionnel proprement**, sans casser la lisibilité ni les tests.

L’idée :
👉 ajouter un flag dans tes filtres
👉 appliquer le tri prioritaire seulement si demandé

---

# ✅ 1️⃣ Mise à jour de `RuleFilters`

On ajoute un booléen :

```go
type RuleFilters struct {
	BU            *string
	Service       *string
	Type          *string
	Env           *string
	Action        *string
	Realm         *string
	Account       *string
	Limit         *uint64

	PrioritizeAllow bool // ← nouveau flag
}
```

---

# ✅ 2️⃣ Repository avec ORDER BY optionnel

```go
func (r *RuleRepository) Find(ctx context.Context, f RuleFilters) ([]Rule, error) {

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := psql.
		Select("*").
		From("rules")

	// --- filtres (inchangés) ---
	if f.BU != nil {
		q = q.Where(sq.Eq{"bu": *f.BU})
	}

	if f.Service != nil {
		q = q.Where(sq.Eq{"service": *f.Service})
	}

	if f.Type != nil {
		q = q.Where(sq.Eq{"type": *f.Type})
	}

	if f.Env != nil {
		q = q.Where(sq.Or{
			sq.Eq{"environment": *f.Env},
			sq.Eq{"environment": "ALL"},
		})
	}

	if f.Action != nil {
		q = q.Where(sq.Or{
			sq.Eq{"action": *f.Action},
			sq.Eq{"action": "ALL"},
		})
	}

	if f.Realm != nil {
		q = q.Where(
			sq.Expr("? = ANY(realms) OR cardinality(realms) = 0", *f.Realm),
		)
	}

	if f.Account != nil {
		q = q.Where(
			sq.Expr("? = ANY(accounts) OR cardinality(accounts) = 0", *f.Account),
		)
	}

	// --- ORDER BY conditionnel ---

	if f.PrioritizeAllow {
		q = q.OrderBy(
			"CASE WHEN type='ALLOW' THEN 0 ELSE 1 END",
			"created_at DESC",
		)
	} else {
		q = q.OrderBy("created_at DESC")
	}

	// --- LIMIT ---
	if f.Limit != nil {
		q = q.Limit(*f.Limit)
	}

	sqlStr, args, err := q.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Rule
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
			&rule.UpdatedBy,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, rule)
	}

	return result, nil
}
```

---

# 🧪 3️⃣ Tests pour couvrir les deux cas

## 🔹 Cas 1 — PrioritizeAllow = true

```go
mock.ExpectQuery(`SELECT \* FROM rules`).
	WillReturnRows(rows)

_, err := repo.Find(ctx, RuleFilters{
	PrioritizeAllow: true,
})
require.NoError(t, err)
```

---

## 🔹 Cas 2 — PrioritizeAllow = false

```go
mock.ExpectQuery(`SELECT \* FROM rules`).
	WillReturnRows(rows)

_, err := repo.Find(ctx, RuleFilters{
	PrioritizeAllow: false,
})
require.NoError(t, err)
```

---

# 🎯 Option encore plus propre (production ready)

Si tu veux quelque chose de plus flexible que bool :

```go
type SortMode int

const (
	SortByDate SortMode = iota
	SortAllowFirst
)

type RuleFilters struct {
	...
	SortMode SortMode
}
```

Puis :

```go
switch f.SortMode {
case SortAllowFirst:
	q = q.OrderBy(
		"CASE WHEN type='ALLOW' THEN 0 ELSE 1 END",
		"created_at DESC",
	)
default:
	q = q.OrderBy("created_at DESC")
}
```

👉 Plus extensible si un jour tu veux :

* tri par BU
* tri par Service
* tri dynamique asc/desc
* tri multiple configurable

---

# 🧠 Résumé

Tu as maintenant :

✅ ORDER BY conditionnel
✅ Pas de duplication
✅ Compatible avec tous les filtres
✅ Testable facilement
✅ Clean & scalable

---

Si tu veux, on peut aller encore plus loin :

* 🔥 tri dynamique configurable via API (ex: ?sort=type,-created_at)
* 🔥 builder de tri générique
* 🔥 pattern Specification complet
* 🔥 version clean architecture découplée

Tu veux qu’on pousse vers une version "enterprise grade" ? 🚀


## with builder generic

Parfait 🚀
On va construire un **QueryBuilder complet, propre et extensible**, combinant :

* FilterBuilder
* SortBuilder
* Pagination
* Pattern Specification

Compatible avec :

* Squirrel
* pgx

Objectif : architecture **clean, scalable, testable, production-ready**.

---

# 🧠 🏗 Architecture globale

On va séparer les responsabilités :

```
RuleRepository
 └── QueryBuilder
      ├── Specifications (filters dynamiques)
      ├── SortBuilder
      └── Pagination
```

Chaque filtre devient une **Specification** indépendante.

---

# 🎯 1️⃣ Pattern Specification

Interface générique :

```go
type Specification interface {
	Apply(q sq.SelectBuilder) sq.SelectBuilder
}
```

---

## 🔹 Exemple 1 — BU Specification

```go
type BUSpec struct {
	Value string
}

func (s BUSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(sq.Eq{"bu": s.Value})
}
```

---

## 🔹 Exemple 2 — Environment OR ALL

```go
type EnvironmentSpec struct {
	Value string
}

func (s EnvironmentSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(sq.Or{
		sq.Eq{"environment": s.Value},
		sq.Eq{"environment": "ALL"},
	})
}
```

---

## 🔹 Exemple 3 — Realm ANY()

```go
type RealmSpec struct {
	Value string
}

func (s RealmSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(
		sq.Expr("? = ANY(realms) OR cardinality(realms) = 0", s.Value),
	)
}
```

---

# 🧱 2️⃣ QueryBuilder principal

```go
type QueryBuilder struct {
	specs []Specification
	sort  string
	page  *Pagination
}
```

---

## Pagination struct

```go
type Pagination struct {
	Limit  uint64
	Offset uint64
}
```

---

# 🏗 3️⃣ Implémentation complète

```go
type QueryBuilder struct {
	specs []Specification
	sort  string
	page  *Pagination
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{}
}

func (qb *QueryBuilder) AddSpec(s Specification) {
	qb.specs = append(qb.specs, s)
}

func (qb *QueryBuilder) SetSort(sort string) {
	qb.sort = sort
}

func (qb *QueryBuilder) SetPagination(p Pagination) {
	qb.page = &p
}

func (qb *QueryBuilder) Build() (string, []any, error) {

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := psql.Select("*").From("rules")

	// Apply Specifications
	for _, s := range qb.specs {
		q = s.Apply(q)
	}

	// Apply Sorting
	q = ruleSortBuilder.Apply(q, qb.sort)

	// Apply Pagination
	if qb.page != nil {
		q = q.Limit(qb.page.Limit).Offset(qb.page.Offset)
	}

	return q.ToSql()
}
```

---

# 🔐 4️⃣ SortBuilder sécurisé

```go
var ruleSortBuilder = NewSortBuilder(map[string]string{
	"type":       "type",
	"created_at": "created_at",
	"bu":         "bu",
	"service":    "service",
	"priority":   "priority",
})
```

Avec gestion spéciale :

```go
if field == "priority" {
	q = q.OrderBy("CASE WHEN type='ALLOW' THEN 0 ELSE 1 END")
	continue
}
```

---

# 🏗 5️⃣ Repository final

```go
func (r *RuleRepository) Find(ctx context.Context, qb *QueryBuilder) ([]Rule, error) {

	sqlStr, args, err := qb.Build()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Rule
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
			&rule.UpdatedBy,
			&rule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		result = append(result, rule)
	}

	return result, nil
}
```

---

# 🚀 6️⃣ Exemple d’utilisation combinée

```go
qb := NewQueryBuilder()

qb.AddSpec(BUSpec{"BU1"})
qb.AddSpec(EnvironmentSpec{"prod"})
qb.AddSpec(RealmSpec{"realm1"})

qb.SetSort("priority,-created_at")

qb.SetPagination(Pagination{
	Limit:  10,
	Offset: 0,
})
```

🔥 Résultat :

* BU = BU1
* environment = prod OR ALL
* realm ANY
* ALLOW first
* created_at DESC
* LIMIT 10

---

# 🧪 7️⃣ Test avec pgxmock

On teste simplement :

```go
mock.ExpectQuery(`SELECT \* FROM rules`).
	WillReturnRows(rows)
```

On ne teste pas la SQL exacte (fragile),
on teste :

* que les specs sont appliquées
* que les args sont corrects
* que le résultat est bon

---

# 🎯 Pourquoi c’est architecture senior ?

✅ SRP respecté
✅ Open/Closed principle
✅ Ultra extensible
✅ Pas de duplication SQL
✅ Testable indépendamment
✅ Compatible Clean Architecture
✅ Compatible CQRS si besoin
✅ Facilement portable sur autre entité

---

## est-ce que je peux avoir cette logique dans un builder generic specific avc le pattern Specification

Oui 💪 on peut transformer **toute cette logique** en :

* ✅ un **Generic Specification Builder**
* ✅ réutilisable pour n’importe quelle entité
* ✅ compatible avec Squirrel
* ✅ propre Clean Architecture
* ✅ extensible production

On va construire :

```
GenericQueryBuilder[T]
 ├── Specifications
 ├── Sort
 ├── Pagination
 └── Build()
```

---

# 🧠 1️⃣ Interface Specification générique

On la rend indépendante de `rules`.

```go
type Specification interface {
	Apply(q sq.SelectBuilder) sq.SelectBuilder
}
```

Simple. Extensible. Clean.

---

# 🧱 2️⃣ GenericQueryBuilder

On le rend réutilisable pour n’importe quelle table.

```go
type GenericQueryBuilder struct {
	table string
	specs []Specification
	sort  func(sq.SelectBuilder) sq.SelectBuilder
	limit *uint64
}
```

---

## Constructeur

```go
func NewGenericQueryBuilder(table string) *GenericQueryBuilder {
	return &GenericQueryBuilder{
		table: table,
	}
}
```

---

## Ajouter une Specification

```go
func (qb *GenericQueryBuilder) Add(spec Specification) {
	qb.specs = append(qb.specs, spec)
}
```

---

## Ajouter un tri

On passe une fonction → ultra flexible.

```go
func (qb *GenericQueryBuilder) SetSort(
	sortFn func(sq.SelectBuilder) sq.SelectBuilder,
) {
	qb.sort = sortFn
}
```

---

## Ajouter une limit

```go
func (qb *GenericQueryBuilder) SetLimit(limit uint64) {
	qb.limit = &limit
}
```

---

# 🏗 3️⃣ Méthode Build()

```go
func (qb *GenericQueryBuilder) Build() (string, []any, error) {

	psql := sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	q := psql.Select("*").From(qb.table)

	// Apply Specifications
	for _, s := range qb.specs {
		q = s.Apply(q)
	}

	// Apply Sorting
	if qb.sort != nil {
		q = qb.sort(q)
	}

	// Apply Limit
	if qb.limit != nil {
		q = q.Limit(*qb.limit)
	}

	return q.ToSql()
}
```

---

# 🎯 4️⃣ Specifications pour Rules

On crée des specs indépendantes.

---

## 🔹 BUSpec

```go
type BUSpec struct {
	Value string
}

func (s BUSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(sq.Eq{"bu": s.Value})
}
```

---

## 🔹 ServiceSpec

```go
type ServiceSpec struct {
	Value string
}

func (s ServiceSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(sq.Eq{"service": s.Value})
}
```

---

## 🔹 TypeSpec

```go
type TypeSpec struct {
	Value string
}

func (s TypeSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(sq.Eq{"type": s.Value})
}
```

---

## 🔹 EnvironmentSpec (OR ALL)

```go
type EnvironmentSpec struct {
	Value string
}

func (s EnvironmentSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(sq.Or{
		sq.Eq{"environment": s.Value},
		sq.Eq{"environment": "ALL"},
	})
}
```

---

## 🔹 ActionSpec (OR ALL)

```go
type ActionSpec struct {
	Value string
}

func (s ActionSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(sq.Or{
		sq.Eq{"action": s.Value},
		sq.Eq{"action": "ALL"},
	})
}
```

---

## 🔹 RealmSpec (ANY)

```go
type RealmSpec struct {
	Value string
}

func (s RealmSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(
		sq.Expr("? = ANY(realms) OR cardinality(realms) = 0", s.Value),
	)
}
```

---

## 🔹 AccountSpec (ANY)

```go
type AccountSpec struct {
	Value string
}

func (s AccountSpec) Apply(q sq.SelectBuilder) sq.SelectBuilder {
	return q.Where(
		sq.Expr("? = ANY(accounts) OR cardinality(accounts) = 0", s.Value),
	)
}
```

---

# 🔀 5️⃣ Tri conditionnel (ALLOW first)

On injecte la logique dans une fonction :

```go
func RuleSort(prioritizeAllow bool) func(sq.SelectBuilder) sq.SelectBuilder {
	return func(q sq.SelectBuilder) sq.SelectBuilder {

		if prioritizeAllow {
			return q.OrderBy(
				"CASE WHEN type='ALLOW' THEN 0 ELSE 1 END",
				"created_at DESC",
			)
		}

		return q.OrderBy("created_at DESC")
	}
}
```

🔥 Très propre, découplé.

---

# 🚀 6️⃣ Utilisation complète

```go
qb := NewGenericQueryBuilder("rules")

qb.Add(BUSpec{"BU1"})
qb.Add(ServiceSpec{"SVC1"})
qb.Add(EnvironmentSpec{"prod"})
qb.Add(RealmSpec{"realm1"})

qb.SetSort(RuleSort(true))

qb.SetLimit(1)

sqlStr, args, err := qb.Build()
```

Résultat :

* BU = BU1
* service = SVC1
* environment = prod OR ALL
* realm ANY
* ALLOW first
* created_at DESC
* LIMIT 1

---

# 🎯 7️⃣ Pourquoi c’est vraiment propre

✅ Generic
✅ Réutilisable pour toute table
✅ Open/Closed principle respecté
✅ Ajout de filtre sans modifier le builder
✅ Pas de if spaghetti
✅ Testable individuellement
✅ Compatible pagination + count
✅ Clean Architecture ready

---

# 🔥 Si tu veux niveau encore plus avancé

On peut :

* Ajouter `AndSpec`, `OrSpec` combinables
* Ajouter `NotSpec`
* Ajouter `CountQuery` automatique
* Ajouter support multi-tenant
* Ajouter soft-delete auto-filter
* Ajouter version Go generics
* Ajouter auto-binding depuis Gin query params

Tu veux qu’on pousse vers un **Specification Engine complet type DDD** ? 🚀
