L’opérateur `...` (trois points) en Go est utilisé dans plusieurs contextes avec la **librairie standard**, notamment pour :

### 1. **Fonction avec un nombre variable d'arguments (`variadic functions`)**

Permet de passer un nombre variable d'arguments à une fonction.

```go
package main

import (
	"fmt"
)

func somme(nombres ...int) int {
	total := 0
	for _, n := range nombres {
		total += n
	}
	return total
}

func main() {
	fmt.Println(somme(1, 2, 3))       // 6
	fmt.Println(somme(10, 20, 30, 40)) // 100
}
```

### 2. **Passage d’un slice comme arguments variadiques**

On peut utiliser `...` pour **décomposer un slice** et le passer à une fonction variadique :

```go
n := []int{4, 5, 6}
fmt.Println(somme(n...)) // 15
```

### 3. **Déclaration d’un tableau avec taille implicite**

Le `...` peut aussi être utilisé pour laisser le compilateur déterminer la taille d’un tableau :

```go
jours := [...]string{"Lundi", "Mardi", "Mercredi"}
fmt.Println(len(jours)) // 3
```

---

### example avec des chaines de caractères

Voici un exemple utilisant l'opérateur `...` pour manipuler des chaînes en Go avec la librairie standard :

### 1. **Concaténation de chaînes avec une fonction variadique**

L'opérateur `...` permet de passer un nombre variable d'arguments à une fonction qui combine des chaînes.

```go
package main

import (
	"fmt"
	"strings"
)

// Concaténer un nombre variable de chaînes
func concatener(separateur string, mots ...string) string {
	return strings.Join(mots, separateur)
}

func main() {
	// Exemple avec des arguments variadiques
	resultat := concatener(" ", "Bonjour", "à", "toi", "!")
	fmt.Println(resultat) // "Bonjour à toi !"

	// Exemple avec un slice
	mots := []string{"Go", "est", "super"}
	fmt.Println(concatener("-", mots...)) // "Go-est-super"
}
```

### 2. **Conversion et manipulation avec des slices de chaînes**

L’opérateur `...` est utile lorsqu’on veut manipuler des collections de chaînes et les passer à des fonctions variadiques comme `strings.Join`.

```go
package main

import (
	"fmt"
	"strings"
)

func main() {
	mots := []string{"Je", "code", "en", "Go"}
	
	// Utiliser Join avec un slice
	phrase := strings.Join(mots, " ")
	fmt.Println(phrase) // "Je code en Go"

	// Ajouter des mots supplémentaires
	motsSup := []string{"et", "c'est", "cool"}
	phraseComplete := strings.Join(append(mots, motsSup...), " ")
	fmt.Println(phraseComplete) // "Je code en Go et c'est cool"
}
```

Ces exemples montrent comment le `...` peut simplifier la manipulation de chaînes dans des scénarios pratiques, tout en maintenant la lisibilité du code.


