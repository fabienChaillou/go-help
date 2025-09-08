En Go, `os.Pipe()` est une fonction de la bibliothèque standard `os` qui permet de créer **un canal de communication unidirectionnel** entre deux objets de type `*os.File`. Elle renvoie **deux fichiers** : un pour l’écriture (le writer) et un pour la lecture (le reader).

### Syntaxe

```go
r, w, err := os.Pipe()
```

* `r` : un *reader* (`*os.File`) — il lit les données.
* `w` : un *writer* (`*os.File`) — il écrit les données.
* `err` : une erreur en cas d’échec.

### Utilité

C’est souvent utilisé pour :

* simuler une communication entre processus (comme avec des pipes Unix),
* intercepter ou rediriger des sorties (`stdout`, `stderr`),
* créer des flux entre différentes parties d’un programme.

### Exemple simple

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	go func() {
		fmt.Fprintln(w, "Hello from pipe!")
		w.Close() // important de fermer le writer
	}()

	buf := make([]byte, 100)
	n, err := r.Read(buf)
	if err != nil && err != io.EOF {
		panic(err)
	}

	fmt.Printf("Received: %s", string(buf[:n]))
}
```

### Explication

* Le `goroutine` écrit une chaîne dans le `writer`.
* Le programme principal lit à partir du `reader`.
* On ferme le `writer` pour signaler la fin de l’écriture, ce qui est crucial pour que le lecteur n’attende pas indéfiniment.

Souhaites-tu un exemple d'utilisation dans un cas concret, comme rediriger `os.Stdout` ?

Voici un exemple concret d'utilisation de `os.Pipe()` pour **rediriger `os.Stdout`**, afin de capturer ce qu’un programme écrit dans la console :

---

### 🎯 Objectif :

Capturer ce qui est imprimé avec `fmt.Println()` en le redirigeant vers une variable plutôt que vers le terminal.

---

### ✅ Exemple de code :

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// Crée un pipe
	r, w, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	// Sauvegarde stdout original
	stdout := os.Stdout
	os.Stdout = w // Redirige stdout vers le writer du pipe

	// Ce texte sera capturé, pas affiché dans le terminal
	fmt.Println("Hello, redirected stdout!")

	// Ferme le writer pour signaler la fin
	w.Close()

	// Restaure stdout
	os.Stdout = stdout

	// Lis ce qui a été écrit dans le pipe
	out, _ := io.ReadAll(r)
	fmt.Printf("Captured output: %q\n", string(out))
}
```

---

### 🧠 Ce qui se passe :

1. `os.Stdout` est redirigé vers le `writer` du pipe.
2. Tout ce qui est imprimé avec `fmt.Print*` va dans ce pipe.
3. On ferme le `writer` pour libérer le flux.
4. On lit les données du `reader`.
5. `os.Stdout` est restauré à sa valeur d’origine.
6. Le texte capturé est affiché.

---

Souhaites-tu voir comment faire la même chose pour `os.Stderr` ou pour logger des tests unitaires ?

Parfait ! Voici comment faire la **même redirection avec `os.Stderr`** (utile pour capturer les erreurs ou logs), et un exemple de **capture dans un test unitaire**.

---

## 🔴 Rediriger `os.Stderr`

```go
package main

import (
	"fmt"
	"io"
	"os"
)

func main() {
	// Crée un pipe
	r, w, _ := os.Pipe()

	// Sauvegarde stderr original
	stderr := os.Stderr
	os.Stderr = w

	// Écrit une erreur simulée
	fmt.Fprintln(os.Stderr, "An error occurred!")

	// Ferme le writer et restaure stderr
	w.Close()
	os.Stderr = stderr

	// Lis la sortie d'erreur capturée
	out, _ := io.ReadAll(r)
	fmt.Printf("Captured stderr: %q\n", string(out))
}
```

---

## ✅ Exemple en test unitaire

### test\_capture\_stdout.go

```go
package main

import (
	"fmt"
	"io"
	"os"
	"testing"
)

func TestCaptureStdout(t *testing.T) {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	// Ce qui est censé être testé
	fmt.Print("Hello test!")

	// Clôture et restauration
	w.Close()
	os.Stdout = stdout

	out, _ := io.ReadAll(r)
	if string(out) != "Hello test!" {
		t.Errorf("Expected 'Hello test!', got %q", string(out))
	}
}
```

---

### 💡 Astuce :

Tu peux encapsuler la logique dans une fonction utilitaire pour réutiliser facilement la redirection de sortie :

```go
func captureOutput(f func()) string {
	r, w, _ := os.Pipe()
	stdout := os.Stdout
	os.Stdout = w

	f() // exécute la fonction fournie

	w.Close()
	os.Stdout = stdout

	out, _ := io.ReadAll(r)
	return string(out)
}
```

Et l'utiliser comme :

```go
output := captureOutput(func() {
	fmt.Print("inside function")
})
fmt.Println("Captured:", output)
```

---

Souhaites-tu une version similaire pour capturer la sortie de logs (`log.Println`) ou un exemple avec `stderr` dans un test ?
