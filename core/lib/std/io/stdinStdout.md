En Go, le package `os` fournit une interface pour interagir avec le système d'exploitation. Les objets `os.Stdin`, `os.Stdout` et `os.Stderr` représentent respectivement :

* **`os.Stdin`** : l'entrée standard (standard input) — généralement le clavier ou l'entrée d’un pipe.
* **`os.Stdout`** : la sortie standard (standard output) — généralement le terminal.
* **`os.Stderr`** : la sortie d'erreur standard — aussi généralement le terminal, mais séparée de `Stdout`.

Ces objets sont des *fichiers ouverts* (type `*os.File`) que tu peux lire ou écrire avec les fonctions classiques du package `io`, comme `io.ReadAll`, `fmt.Fprint`, etc.

---

### 🔹 Exemple d’utilisation de `os.Stdin`

Lire depuis l’entrée standard (par exemple, ce que l’utilisateur tape) :

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	fmt.Print("Tape quelque chose : ")
	input, _ := io.ReadAll(os.Stdin)
	fmt.Println("Tu as tapé :", string(input))
}
```

---

### 🔹 Exemple d’utilisation de `os.Stdout`

Tu peux écrire directement dans `os.Stdout` :

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stdout, "Ceci va dans la sortie standard")
}
```

---

### 🔹 Exemple combiné avec `bufio`

Pour lire ligne par ligne :

```go
package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	fmt.Print("Entrée : ")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		fmt.Fprintln(os.Stdout, "Tu as écrit :", scanner.Text())
	}
}
```

---

### 🔸 Remarques

* Tu peux rediriger ces flux dans un programme (utile pour les tests ou les pipes Unix).
* `os.Stdin` est en lecture seule, `os.Stdout` et `os.Stderr` sont en écriture.

Souhaites-tu un exemple dans un contexte spécifique (comme un script CLI ou une API) ?

## dans un context CLI

Parfait ! Voici un **exemple complet d'application CLI en Go** qui lit une entrée utilisateur depuis le terminal (via `os.Stdin`), et écrit une réponse sur `os.Stdout` (et sur `os.Stderr` en cas d'erreur).

---

### ✅ Exemple simple : CLI qui lit un nom et salue l'utilisateur

```go
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Affiche un message à l'utilisateur
	fmt.Fprint(os.Stdout, "Quel est ton nom ? ")

	// Crée un scanner pour lire depuis os.Stdin
	scanner := bufio.NewScanner(os.Stdin)

	// Lit une ligne
	if scanner.Scan() {
		nom := strings.TrimSpace(scanner.Text())

		// Si le nom est vide, affiche une erreur
		if nom == "" {
			fmt.Fprintln(os.Stderr, "Erreur : nom vide")
			os.Exit(1)
		}

		// Salutation
		fmt.Fprintf(os.Stdout, "Bonjour, %s !\n", nom)
	} else {
		// Gestion d'erreur de lecture
		fmt.Fprintln(os.Stderr, "Erreur lors de la lecture de l'entrée.")
		os.Exit(1)
	}
}
```

---

### 🔧 Comment l’utiliser

1. Sauvegarde ce fichier en `main.go`.
2. Exécute le programme :

```bash
go run main.go
```

3. Tape ton nom quand il te le demande.

---

### 💡 Bonus : redirections en ligne de commande

Tu peux aussi tester les flux standards comme ceci :

```bash
echo "Alice" | go run main.go           # utilise un pipe comme os.Stdin
go run main.go > sortie.txt             # redirige os.Stdout dans un fichier
go run main.go 2> erreur.txt            # redirige os.Stderr dans un fichier
```
