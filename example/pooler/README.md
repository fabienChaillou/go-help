J'ai créé un gestionnaire de pools multi-connexions en Go (fichier `pgx_multi_pooler.go`) dans la **zone de code** à côté de la conversation — tu peux l'ouvrir pour voir le code complet, les commentaires et un exemple d'utilisation.

Résumé rapide (ce que fait le code) :

* gère plusieurs pools nommés (`AddPool`, `GetPool`, `ClosePool`, `CloseAll`) basés sur `pgxpool` de pgx v5,
* logique de retry simple à la création d'un pool,
* wrappers pratiques `Query` et `Exec` pour exécuter des requêtes sur un pool nommé,
* health-checks récurrents dans une goroutine (avec callback `onFailure`),
* exemple d'utilisation inclus en commentaire.

Si tu veux que je :

* transforme ça en module/paquet réutilisable avec tests unitaires,
* ajoute des métriques (Prometheus) / logging structuré / prepared statements partagés,
* ajoute gestion du failover automatique entre pools (round-robin, priorité, poids) ou un pool de lecture/écriture séparé,

— dis-moi ce que tu veux ajouter et je l'intègre directement dans le fichier.

J'ai transformé le code en un module réutilisable et j'ai ajouté un guide pour les tests unitaires dans un fichier `pooler_test.go`.
Le module est maintenant prêt à être importé dans d'autres projets, et les tests couvrent :

* **Création et récupération d'un pool** (`AddPool` / `GetPool`),
* **Fermeture d'un pool spécifique** (`ClosePool`),
* **Fermeture de tous les pools** (`CloseAll`).

```go
// Package dbpool provides a reusable multi-connection pool manager for PostgreSQL using pgx v5.
// It allows managing multiple named connection pools, provides query/exec helpers,
// and includes built-in health check functionality.

package dbpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pooler manages multiple named pgx connection pools.
// It provides helpers to add/get/close pools, run health checks and convenience query/exec wrappers.
type Pooler struct {
	mu    sync.RWMutex
	pools map[string]*pgxpool.Pool
}

// NewPooler creates an empty Pooler instance.
func NewPooler() *Pooler {
	return &Pooler{
		pools: make(map[string]*pgxpool.Pool),
	}
}

// AddPool creates a pool for the given name and connection string.
// If a pool already exists under the name, it's returned.
// ctx is used for the initial connection and should have a reasonable timeout.
// maxAttempts controls simple retry logic (with backoff) for establishing the pool.
func (m *Pooler) AddPool(ctx context.Context, name, connString string, maxAttempts int) (*pgxpool.Pool, error) {
	if name == "" {
		return nil, errors.New("pool name required")
	}
	m.mu.RLock()
	if p, ok := m.pools[name]; ok {
		m.mu.RUnlock()
		return p, nil
	}
	m.mu.RUnlock()

	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		poolCfg, err := pgxpool.ParseConfig(connString)
		if err != nil {
			return nil, fmt.Errorf("invalid connection string: %w", err)
		}
		if poolCfg.MaxConns == 0 {
			poolCfg.MaxConns = 10
		}

		p, err := pgxpool.NewWithConfig(ctx, poolCfg)
		if err == nil {
			m.mu.Lock()
			defer m.mu.Unlock()
			if existing, ok := m.pools[name]; ok {
				p.Close()
				return existing, nil
			}
			m.pools[name] = p
			return p, nil
		}

		lastErr = err
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(100*(1<<uint(i))) * time.Millisecond):
		}
	}

	return nil, fmt.Errorf("could not create pool %q after %d attempts: %w", name, maxAttempts, lastErr)
}

// GetPool returns the pool by name or nil if not present.
func (m *Pooler) GetPool(name string) *pgxpool.Pool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pools[name]
}

// ClosePool closes and removes the pool with the given name.
func (m *Pooler) ClosePool(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p, ok := m.pools[name]; ok {
		p.Close()
		delete(m.pools, name)
	}
}

// CloseAll closes all pools managed by the Pooler.
func (m *Pooler) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, p := range m.pools {
		p.Close()
		delete(m.pools, k)
	}
}

// Exec runs an SQL command on a named pool.
func (m *Pooler) Exec(ctx context.Context, poolName string, sql string, args ...interface{}) (int64, error) {
	p := m.GetPool(poolName)
	if p == nil {
		return 0, fmt.Errorf("pool %q not found", poolName)
	}
	cmd, err := p.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}

// HealthCheck starts background health checks on all pools at the given interval.
// The onFailure callback is invoked when a pool fails a health check.
func (m *Pooler) StartHealthChecks(cancelCtx context.Context, interval time.Duration, onFailure func(name string, err error)) {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-cancelCtx.Done():
				return
			case <-ticker.C:
				m.mu.RLock()
				for name, p := range m.pools {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					err := p.Ping(ctx)
					cancel()
					if err != nil && onFailure != nil {
						go onFailure(name, err)
					}
				}
				m.mu.RUnlock()
			}
		}
	}()
}

// ===================
// TESTS
// ===================
// Create a file named pooler_test.go in the same package with the following content:
//
// package dbpool
//
// import (
//     "context"
//     "testing"
//     "time"
// )
//
// func TestAddAndGetPool(t *testing.T) {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()
//
//     p := NewPooler()
//     pool, err := p.AddPool(ctx, "test", "postgres://user:pass@localhost:5432/dbname?sslmode=disable", 1)
//     if err != nil {
//         t.Fatalf("failed to add pool: %v", err)
//     }
//     if pool == nil {
//         t.Fatal("expected pool to be non-nil")
//     }
//     if got := p.GetPool("test"); got == nil {
//         t.Fatal("GetPool returned nil, expected pool")
//     }
// }
//
// func TestClosePool(t *testing.T) {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()
//
//     p := NewPooler()
//     _, _ = p.AddPool(ctx, "test", "postgres://user:pass@localhost:5432/dbname?sslmode=disable", 1)
//     p.ClosePool("test")
//     if got := p.GetPool("test"); got != nil {
//         t.Fatal("expected pool to be nil after ClosePool")
//     }
// }
//
// func TestCloseAll(t *testing.T) {
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()
//
//     p := NewPooler()
//     _, _ = p.AddPool(ctx, "db1", "postgres://user:pass@localhost:5432/dbname?sslmode=disable", 1)
//     _, _ = p.AddPool(ctx, "db2", "postgres://user:pass@localhost:5432/dbname?sslmode=disable", 1)
//
//     p.CloseAll()
//     if len(p.pools) != 0 {
//         t.Fatal("expected all pools to be closed")
//     }
// }
````

Souhaites-tu que je t’ajoute un `go.mod` et un exemple de structure de projet complète avec les commandes `go test` et `go run` ?

Voici une structure complète pour ton projet Go, avec le module `dbpool`, un `go.mod`, des tests unitaires, et un exemple de démarrage :

---

## **Structure de projet**

```
multi-pooler/
├── go.mod
├── go.sum
├── cmd/
│   └── example/
│       └── main.go
├── internal/
│   └── dbpool/
│       ├── pooler.go
│       └── pooler_test.go
└── README.md
```

---

## **1. go.mod**

Ce fichier déclare ton module et ses dépendances.

```go
module github.com/tonusername/multi-pooler

go 1.22

require (
	github.com/jackc/pgx/v5 v5.6.0
)
```

> ⚡️ Remplace `github.com/tonusername/multi-pooler` par le nom de ton dépôt GitHub.

Pour l'initialiser :

```bash
go mod init github.com/tonusername/multi-pooler
go get github.com/jackc/pgx/v5
```

---

## **2. internal/dbpool/pooler.go**

Le code que nous avons transformé précédemment.
Ce fichier contient la logique du **multi-pooler**.

```go
// internal/dbpool/pooler.go
package dbpool

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Pooler manages multiple named pgx connection pools.
type Pooler struct {
	mu    sync.RWMutex
	pools map[string]*pgxpool.Pool
}

// NewPooler creates an empty Pooler instance.
func NewPooler() *Pooler {
	return &Pooler{
		pools: make(map[string]*pgxpool.Pool),
	}
}

// AddPool creates a new pool or returns an existing one if it already exists.
func (m *Pooler) AddPool(ctx context.Context, name, connString string, maxAttempts int) (*pgxpool.Pool, error) {
	if name == "" {
		return nil, errors.New("pool name required")
	}
	m.mu.RLock()
	if p, ok := m.pools[name]; ok {
		m.mu.RUnlock()
		return p, nil
	}
	m.mu.RUnlock()

	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		poolCfg, err := pgxpool.ParseConfig(connString)
		if err != nil {
			return nil, fmt.Errorf("invalid connection string: %w", err)
		}
		if poolCfg.MaxConns == 0 {
			poolCfg.MaxConns = 10
		}

		p, err := pgxpool.NewWithConfig(ctx, poolCfg)
		if err == nil {
			m.mu.Lock()
			defer m.mu.Unlock()
			if existing, ok := m.pools[name]; ok {
				p.Close()
				return existing, nil
			}
			m.pools[name] = p
			return p, nil
		}

		lastErr = err
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(time.Duration(100*(1<<uint(i))) * time.Millisecond):
		}
	}

	return nil, fmt.Errorf("could not create pool %q after %d attempts: %w", name, maxAttempts, lastErr)
}

// GetPool retrieves a pool by its name.
func (m *Pooler) GetPool(name string) *pgxpool.Pool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.pools[name]
}

// ClosePool closes and removes a single pool.
func (m *Pooler) ClosePool(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p, ok := m.pools[name]; ok {
		p.Close()
		delete(m.pools, name)
	}
}

// CloseAll closes and removes all pools.
func (m *Pooler) CloseAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, p := range m.pools {
		p.Close()
		delete(m.pools, k)
	}
}

// Exec executes an SQL command on a named pool.
func (m *Pooler) Exec(ctx context.Context, poolName string, sql string, args ...interface{}) (int64, error) {
	p := m.GetPool(poolName)
	if p == nil {
		return 0, fmt.Errorf("pool %q not found", poolName)
	}
	cmd, err := p.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return cmd.RowsAffected(), nil
}

// StartHealthChecks periodically checks all pools.
func (m *Pooler) StartHealthChecks(cancelCtx context.Context, interval time.Duration, onFailure func(name string, err error)) {
	if interval <= 0 {
		interval = 30 * time.Second
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-cancelCtx.Done():
				return
			case <-ticker.C:
				m.mu.RLock()
				for name, p := range m.pools {
					ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
					err := p.Ping(ctx)
					cancel()
					if err != nil && onFailure != nil {
						go onFailure(name, err)
					}
				}
				m.mu.RUnlock()
			}
		}
	}()
}
```

---

## **3. internal/dbpool/pooler\_test.go**

Tests unitaires pour le module.

```go
package dbpool

import (
	"context"
	"testing"
	"time"
)

// ATTENTION : nécessite une instance PostgreSQL locale ou en Docker
// export DATABASE_URL=postgres://user:pass@localhost:5432/testdb?sslmode=disable
const testConn = "postgres://user:pass@localhost:5432/testdb?sslmode=disable"

func TestAddAndGetPool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewPooler()
	pool, err := p.AddPool(ctx, "test", testConn, 1)
	if err != nil {
		t.Fatalf("failed to add pool: %v", err)
	}
	if pool == nil {
		t.Fatal("expected pool to be non-nil")
	}
	if got := p.GetPool("test"); got == nil {
		t.Fatal("GetPool returned nil, expected pool")
	}
}

func TestClosePool(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewPooler()
	_, _ = p.AddPool(ctx, "test", testConn, 1)
	p.ClosePool("test")
	if got := p.GetPool("test"); got != nil {
		t.Fatal("expected pool to be nil after ClosePool")
	}
}

func TestCloseAll(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := NewPooler()
	_, _ = p.AddPool(ctx, "db1", testConn, 1)
	_, _ = p.AddPool(ctx, "db2", testConn, 1)

	p.CloseAll()
	if len(p.pools) != 0 {
		t.Fatal("expected all pools to be closed")
	}
}
```

---

## **4. Exemple d’utilisation (cmd/example/main.go)**

```go
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/tonusername/multi-pooler/internal/dbpool"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	p := dbpool.NewPooler()

	connStr := "postgres://user:pass@localhost:5432/testdb?sslmode=disable"

	_, err := p.AddPool(ctx, "primary", connStr, 3)
	if err != nil {
		panic(err)
	}

	healthCtx, healthCancel := context.WithCancel(context.Background())
	defer healthCancel()

	p.StartHealthChecks(healthCtx, 10*time.Second, func(name string, err error) {
		fmt.Printf("Health check failed for pool %s: %v\n", name, err)
	})

	affected, err := p.Exec(ctx, "primary", "CREATE TABLE IF NOT EXISTS demo (id SERIAL PRIMARY KEY, name TEXT);")
	if err != nil {
		fmt.Printf("Exec error: %v\n", err)
	} else {
		fmt.Printf("Rows affected: %d\n", affected)
	}

	fmt.Println("Multi-pooler initialized successfully!")
}
```

---

## **5. Commandes utiles**

### Lancer les tests :

```bash
go test ./...
```

### Exécuter l’exemple :

```bash
go run ./cmd/example
```

### Mettre à jour les dépendances :

```bash
go mod tidy
```

---

## **6. README.md**

````markdown
# Multi-Pooler PostgreSQL (Go + pgx v5)

Ce module gère plusieurs pools PostgreSQL avec `pgx/v5`.

## Fonctionnalités
- Multi-connection pools
- Retry avec backoff exponentiel
- Health checks automatiques
- Tests unitaires inclus

## Installation

```bash
go get github.com/tonusername/multi-pooler
````

## Utilisation

Voir `cmd/example/main.go`.

## Tests

```bash
go test ./...
```
