Bonne question — le **migrateur de base de données en Go** le plus utilisé s’appelle généralement **`golang-migrate/migrate`**. Je vais te montrer **comment l’installer et l’utiliser**, avec des exemples concrets (PostgreSQL et SQLite), en mode **CLI** et **Go**.

---

## 1️⃣ Qu’est-ce que `migrate` ?

`migrate` sert à :

* versionner ton schéma de base de données
* appliquer / annuler des migrations (`up` / `down`)
* garder ta DB synchronisée avec ton code

👉 Très courant dans les projets Go (Gin, Fiber, Clean Architecture, etc.).

---

## 2️⃣ Installation

### 🔹 Option A — via le binaire CLI (recommandé)

#### macOS

```bash
brew install golang-migrate
```

#### Linux

```bash
curl -L https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

#### Vérification

```bash
migrate -version
```

---

### 🔹 Option B — via Go (moins courant pour le CLI)

```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

---

## 3️⃣ Structure classique du projet

```text
project/
├── cmd/
├── internal/
├── migrations/
│   ├── 000001_init.up.sql
│   ├── 000001_init.down.sql
│   ├── 000002_add_users.up.sql
│   └── 000002_add_users.down.sql
└── main.go
```

📌 **Important**

* Les fichiers sont **numérotés**
* Chaque migration a un `.up.sql` et `.down.sql`

---

## 4️⃣ Créer une migration

```bash
migrate create -ext sql -dir migrations -seq init
```

Résultat :

```text
migrations/
├── 000001_init.up.sql
└── 000001_init.down.sql
```

---

## 5️⃣ Exemple de migration

### `000001_init.up.sql`

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now()
);
```

### `000001_init.down.sql`

```sql
DROP TABLE users;
```

---

## 6️⃣ Appliquer les migrations (CLI)

### 🔹 PostgreSQL

```bash
migrate \
  -path migrations \
  -database "postgres://user:password@localhost:5432/mydb?sslmode=disable" \
  up
```

### 🔹 SQLite

```bash
migrate \
  -path migrations \
  -database "sqlite3://db.sqlite" \
  up
```

---

## 7️⃣ Annuler une migration

```bash
migrate -path migrations -database DB_URL down 1
```

Ou tout annuler :

```bash
migrate down
```

---

## 8️⃣ Vérifier l’état

```bash
migrate -path migrations -database DB_URL version
```

---

## 9️⃣ Utiliser `migrate` **dans ton code Go**

⚠️ Moins fréquent, mais utile pour automatiser au démarrage.

### Installation

```bash
go get github.com/golang-migrate/migrate/v4
go get github.com/golang-migrate/migrate/v4/database/postgres
go get github.com/golang-migrate/migrate/v4/source/file
```

### Exemple PostgreSQL

```go
import (
    "log"

    "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
)

func runMigrations() {
    m, err := migrate.New(
        "file://migrations",
        "postgres://user:password@localhost:5432/mydb?sslmode=disable",
    )
    if err != nil {
        log.Fatal(err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        log.Fatal(err)
    }
}
```

📌 Appelé au démarrage de ton app (`main()`).

---

## 🔟 Bonnes pratiques

✅ Toujours versionner les migrations
✅ Une migration = une responsabilité
✅ Ne jamais modifier une migration déjà appliquée
✅ Utiliser le CLI en CI/CD
✅ Tester sur une DB vierge

---

## 🧠 Astuce pour ton contexte Go (Gin / MVC)

Vu que tu travailles souvent avec :

* Gin
* SQLite / PostgreSQL
* architecture MVC

👉 Je te recommande :

* **migrations via CLI**
* **script Makefile** ou `taskfile.yml`

Exemple `Makefile` :

```makefile
migrate-up:
	migrate -path migrations -database $(DB_URL) up
```

---

## 🚀 Tu veux aller plus loin ?

Je peux te fournir :

* un **exemple complet Gin + migrate + SQLite**
* un **setup Docker + migrate**
* un **Makefile prêt à l’emploi**
* ou comparer avec **goose** / **atlas**

👉 Dis-moi ton stack exact (DB + framework) et je t’adapte l’exemple.

Parfait — voici un **setup Docker propre et standard pour `golang-migrate` + PostgreSQL**, prêt à l’emploi 💪
(ce modèle est très utilisé en Go / Gin / FastAPI / CI).

---

## 🎯 Objectif

* PostgreSQL dans Docker
* `migrate` exécuté **dans un container**
* migrations SQL montées en volume
* reproductible local + CI

---

## 1️⃣ Arborescence du projet

```text
project/
├── docker-compose.yml
├── migrations/
│   ├── 000001_init.up.sql
│   ├── 000001_init.down.sql
│   └── 000002_add_users.up.sql
└── .env
```

---

## 2️⃣ Fichier `.env`

```env
POSTGRES_DB=mydb
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=db
POSTGRES_PORT=5432

DATABASE_URL=postgres://postgres:postgres@db:5432/mydb?sslmode=disable
```

---

## 3️⃣ `docker-compose.yml`

```yaml
version: "3.9"

services:
  db:
    image: postgres:16
    container_name: postgres
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
      retries: 5

  migrate:
    image: migrate/migrate:v4.17.1
    container_name: migrate
    depends_on:
      db:
        condition: service_healthy
    volumes:
      - ./migrations:/migrations
    command: [
      "-path", "/migrations",
      "-database", "${DATABASE_URL}",
      "up"
    ]
    restart: "no"

volumes:
  pgdata:
```

📌 **Important**

* `depends_on + healthcheck` évite que migrate démarre trop tôt
* le container `migrate` s’arrête après exécution (normal)

---

## 4️⃣ Exemple de migration

### `000001_init.up.sql`

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT now()
);
```

### `000001_init.down.sql`

```sql
DROP TABLE users;
```

---

## 5️⃣ Lancer les migrations 🚀

```bash
docker compose up migrate
```

Ou tout lancer :

```bash
docker compose up
```

Résultat attendu :

```text
migrate  | Applying migrations...
migrate  | Successfully applied 1 migration
```

---

## 6️⃣ Vérifier la DB

```bash
docker exec -it postgres psql -U postgres -d mydb
```

```sql
\dt
SELECT * FROM schema_migrations;
```

---

## 7️⃣ Commandes utiles

### 🔽 Rollback d’une migration

```bash
docker compose run --rm migrate down 1
```

### 🔍 Voir la version courante

```bash
docker compose run --rm migrate version
```

### 🧹 Reset complet

```bash
docker compose down -v
```

---

## 8️⃣ Variante : migrations **à la demande** (recommandé)

Ne pas lancer migrate automatiquement :

```yaml
command: ["sleep", "infinity"]
```

Puis :

```bash
docker compose run --rm migrate up
```

👉 Très pratique en **dev** et **CI/CD**

---

## 9️⃣ Bonnes pratiques en prod

✅ migrations exécutées **avant** le démarrage de l’app
✅ migrations **idempotentes**
✅ jamais modifier une migration déjà appliquée
✅ une migration = un changement logique
✅ versionner `migrations/`

---

## 🔥 Tu veux aller plus loin ?

Je peux te donner :

* un **Dockerfile Go + migrate + Gin**
* un **Makefile (migrate-up / migrate-down)**
* un **setup CI GitHub Actions**
* une **alternative avec `goose` ou `atlas`**

👉 Dis-moi si c’est pour **dev, prod ou CI** et si ton app est en **Go** (Gin / Fiber).

Parfait 👍
Voici un **setup simple et confortable pour du dev local avec Huma + PostgreSQL + migrate dans Docker**.
Objectif : *tu lances une commande, la DB est prête, les migrations sont appliquées, Huma démarre*.

---

## 🎯 Stack cible

* **Huma** (API Go)
* **PostgreSQL** (Docker)
* **golang-migrate** (Docker)
* Dev **local**, hot-reload possible

---

## 1️⃣ Arborescence recommandée

```text
project/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── api/
│   ├── db/
│   │   └── db.go
│   └── handlers/
├── migrations/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── docker-compose.yml
├── .env
├── go.mod
└── Makefile
```

---

## 2️⃣ `.env` (local)

```env
POSTGRES_DB=huma_dev
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
POSTGRES_HOST=db
POSTGRES_PORT=5432

DATABASE_URL=postgres://postgres:postgres@db:5432/huma_dev?sslmode=disable
```

---

## 3️⃣ `docker-compose.yml` (dev local)

```yaml
version: "3.9"

services:
  db:
    image: postgres:16
    container_name: huma-postgres
    environment:
      POSTGRES_DB: ${POSTGRES_DB}
      POSTGRES_USER: ${POSTGRES_USER}
      POSTGRES_PASSWORD: ${POSTGRES_PASSWORD}
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      retries: 5

  migrate:
    image: migrate/migrate:v4.17.1
    profiles: ["tools"]
    depends_on:
      db:
        condition: service_healthy
    volumes:
      - ./migrations:/migrations
    command: [
      "-path", "/migrations",
      "-database", "${DATABASE_URL}",
      "up"
    ]

volumes:
  pgdata:
```

📌

* `profiles: ["tools"]` → pas lancé automatiquement
* parfait pour du **dev local contrôlé**

---

## 4️⃣ Makefile (qualité de vie ✨)

```makefile
DB_URL=postgres://postgres:postgres@localhost:5432/huma_dev?sslmode=disable

db-up:
	docker compose up -d db

db-down:
	docker compose down -v

migrate-up:
	docker compose --profile tools run --rm migrate up

migrate-down:
	docker compose --profile tools run --rm migrate down 1

migrate-version:
	docker compose --profile tools run --rm migrate version
```

---

## 5️⃣ Connexion DB dans ton app Huma

### `internal/db/db.go`

```go
package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func Connect(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	return db
}
```

---

## 6️⃣ `main.go` (Huma minimal)

```go
package main

import (
	"os"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humagin"

	"github.com/gin-gonic/gin"
	"project/internal/db"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	dbConn := db.Connect(dsn)

	router := gin.Default()
	api := humagin.New(router, huma.DefaultConfig("Huma API", "1.0.0"))

	huma.Get(api, "/health", func() (*struct {
		Status string `json:"status"`
	}, error) {
		return &struct {
			Status string `json:"status"`
		}{Status: "ok"}, nil
	})

	_ = dbConn // injecte-le dans tes handlers ensuite

	router.Run(":8080")
}
```

---

## 7️⃣ Workflow dev local recommandé 🔄

```bash
make db-up
make migrate-up
go run ./cmd/api
```

👉 Une fois la DB prête, tu ne relances que l’API.

---

## 8️⃣ Astuce dev 💡

### 🔹 Reset DB complet

```bash
make db-down
make db-up
make migrate-up
```

### 🔹 Vérifier migrations

```bash
make migrate-version
```

---

## 9️⃣ Pourquoi ce setup marche bien avec Huma

✅ Huma reste **framework-agnostic**
✅ Migrations découplées de l’app
✅ Pas de migrations automatiques en prod
✅ Dev rapide, propre, reproductible

---

## 🚀 Prochaines améliorations possibles

* hot-reload avec `air`
* tests d’intégration avec DB éphémère
* injection DB propre dans handlers Huma
* Dockerfile pour l’API

Si tu veux, je peux :

* t’ajouter **air + Huma**
* structurer un **repo clean Huma**
* brancher **sqlx / squirrel**
* ou faire une version **SQLite dev / PostgreSQL prod**

👉 Dis-moi 👍
