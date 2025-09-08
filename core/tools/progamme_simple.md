Voici une fonction Go simple qui montre comment utiliser `signal` pour intercepter une interruption du système (comme `Ctrl+C`) et `recover` pour gérer une panique :

```go
package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// Capture des signaux système (ex: Ctrl+C)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine pour intercepter les signaux
	go func() {
		sig := <-sigChan
		fmt.Println("\nSignal reçu :", sig)
		os.Exit(0)
	}()

	// Appel d'une fonction avec recover pour gérer une panique
	safeFunction()
}

func safeFunction() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Panic récupérée :", r)
		}
	}()

	fmt.Println("Fonction en cours d'exécution...")
	// Simule une panique
	panic("Quelque chose s'est mal passé !")
}
```

### Explication :

* `signal.Notify` permet d'intercepter les signaux système.
* `recover()` permet de capturer une panique et d’éviter que le programme ne se termine brutalement.
* Cette fonction affiche un message avant de récupérer la panique.

Souhaites-tu intégrer cette gestion dans une application web ou serveur Go ?

Parfait ! Voici un **exemple d'application web en Go avec Gin** qui utilise `signal` pour gérer proprement l'arrêt du serveur, et `recover` pour éviter les crashs à cause des paniques dans les handlers.

---

### 🔧 Exemple complet avec Gin, signal et recover :

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Création du routeur avec recover intégré
	router := gin.New()
	router.Use(gin.Recovery()) // Intercepte les paniques automatiquement

	// Route simple
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Serveur opérationnel !")
	})

	// Route qui déclenche une panique volontaire
	router.GET("/panic", func(c *gin.Context) {
		panic("Erreur fatale simulée")
	})

	// Configuration du serveur
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Lancement du serveur dans une goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur : %s\n", err)
		}
	}()
	fmt.Println("Serveur démarré sur http://localhost:8080")

	// Gestion des signaux système
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Attend Ctrl+C

	fmt.Println("\nArrêt du serveur en cours...")

	// Shutdown propre du serveur
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Erreur à l'arrêt : %s\n", err)
	}

	fmt.Println("Serveur arrêté proprement")
}
```

---

### ✅ Ce que ça fait :

* **`gin.Recovery()`** intercepte toute panique dans les routes et évite un crash.
* **`signal.Notify`** permet d’arrêter le serveur proprement avec `Ctrl+C`.
* Le serveur écoute sur `:8080`.

Souhaites-tu ajouter un middleware de logging ou de validation en plus ?
