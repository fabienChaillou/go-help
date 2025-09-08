En Go, **l’atomicité** fait référence à l'exécution d'une opération ou d'un ensemble d'opérations de manière **indivisible** : soit tout se produit, soit rien ne se produit. Cela est crucial lorsqu’on travaille avec des **concurrences** (goroutines) ou avec des **transactions en base de données**, afin d’éviter des états inconsistants.

---

## 🔹 1. Atomicité au niveau des **variables partagées** (concurrence)

Quand plusieurs goroutines accèdent à une même variable, il faut assurer que les lectures/écritures ne se "marchent pas dessus". Le package [`sync/atomic`](https://pkg.go.dev/sync/atomic) de la bibliothèque standard fournit des primitives sûres.

### ➤ Exemple : Compteur atomique

```go
package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter int64 = 0
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		}()
	}

	wg.Wait()
	fmt.Println("Final counter:", counter) // Attendu : 1000
}
```

✅ Ici, `atomic.AddInt64` garantit que chaque incrément est effectué de manière **atomique**, sans besoin de mutex.

---

## 🔹 2. Atomicité au niveau des **transactions SQL**

Quand tu interagis avec une base de données (ex. PostgreSQL ou SQLite), tu veux souvent garantir que **plusieurs requêtes soient traitées comme une seule unité de travail**. En Go (avec `database/sql`, `sqlx`, etc.), cela se fait via des **transactions**.

### ➤ Exemple : Transaction avec `sqlx`

```go
tx, err := db.Beginx()
if err != nil {
	log.Fatal(err)
}

_, err = tx.Exec("INSERT INTO users (name) VALUES (?)", "Alice")
if err != nil {
	tx.Rollback() // annule tout
	log.Fatal(err)
}

_, err = tx.Exec("INSERT INTO accounts (user_id, balance) VALUES (?, ?)", 1, 100)
if err != nil {
	tx.Rollback()
	log.Fatal(err)
}

err = tx.Commit() // tout est validé ensemble
if err != nil {
	log.Fatal(err)
}
```

✅ Si l’une des requêtes échoue, tout est annulé. C’est ça, **l’atomicité transactionnelle**.

---

## 🟦 Résumé

| Contexte              | Moyen de garantir l’atomicité              |
| --------------------- | ------------------------------------------ |
| Variables partagées   | `sync/atomic` ou `sync.Mutex`              |
| Accès base de données | Transactions (`db.Begin()`, `tx.Commit()`) |

Souhaites-tu un exemple concret d'atomicité dans ton code Go actuel (API, base SQLite, etc.) ?

Parfait ! Voici un **exemple concret d'atomicité** dans une API Go avec SQLite (comme tu utilises Gin, sqlx, et SQLite).

---

## 🎯 Contexte

Imaginons une API qui crée un **utilisateur** et lui associe automatiquement un **compte** bancaire avec un solde initial.

Les deux opérations doivent être **atomiques** :

* Création de l'utilisateur
* Création du compte

Si une seule des deux échoue, il faut **annuler l’ensemble** (rollback).

---

## 🧱 Modèle simplifié

```go
type User struct {
	ID   int    `db:"id"`
	Name string `db:"name"`
}

type Account struct {
	ID      int `db:"id"`
	UserID  int `db:"user_id"`
	Balance int `db:"balance"`
}
```

---

## 🚀 Handler Gin avec atomicité (transaction `sqlx`)

```go
func CreateUserWithAccount(c *gin.Context) {
	type Request struct {
		Name         string `json:"name"`
		InitialFunds int    `json:"initial_funds"`
	}
	var req Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}

	tx, err := db.Beginx()
	if err != nil {
		c.JSON(500, gin.H{"error": "failed to begin transaction"})
		return
	}

	// Étape 1 : créer l'utilisateur
	res, err := tx.Exec("INSERT INTO users (name) VALUES (?)", req.Name)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "failed to create user"})
		return
	}
	userID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "failed to retrieve user ID"})
		return
	}

	// Étape 2 : créer le compte lié
	_, err = tx.Exec("INSERT INTO accounts (user_id, balance) VALUES (?, ?)", userID, req.InitialFunds)
	if err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "failed to create account"})
		return
	}

	// Commit si tout s’est bien passé
	if err := tx.Commit(); err != nil {
		c.JSON(500, gin.H{"error": "failed to commit transaction"})
		return
	}

	c.JSON(200, gin.H{"message": "user and account created", "user_id": userID})
}
```

---

## ✅ Ce code garantit l’atomicité :

* Si **l’une des deux requêtes échoue**, tout est annulé (rollback).
* L’opération est **indivisible du point de vue de la base de données**.
