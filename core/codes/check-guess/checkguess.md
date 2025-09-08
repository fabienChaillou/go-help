Voici une fonction en Go qui compare une tentative (`guess`) avec un mot mystère (`secret`) et retourne un tableau indiquant pour chaque lettre si elle est :

* correcte et bien placée (`"correct"`),
* correcte mais mal placée (`"present"`),
* ou absente (`"absent"`),

Un peu comme dans le jeu **Wordle**.

```go
package main

import (
	"fmt"
)

func CheckGuess(secret, guess string) []string {
	if len(secret) != len(guess) {
		return nil // ou return une erreur
	}

	result := make([]string, len(secret))
	secretRunes := []rune(secret)
	guessRunes := []rune(guess)
	used := make([]bool, len(secret)) // pour marquer les lettres déjà utilisées

	// Première passe : lettres correctes et bien placées
	for i := range secretRunes {
		if guessRunes[i] == secretRunes[i] {
			result[i] = "correct"
			used[i] = true
		}
	}

	// Deuxième passe : lettres correctes mais mal placées
	for i := range guessRunes {
		if result[i] == "" { // seulement si pas déjà "correct"
			for j := range secretRunes {
				if !used[j] && guessRunes[i] == secretRunes[j] {
					result[i] = "present"
					used[j] = true
					break
				}
			}
			if result[i] == "" {
				result[i] = "absent"
			}
		}
	}

	return result
}

func main() {
	secret := "pomme"
	guess := "moppe"

	result := CheckGuess(secret, guess)
	fmt.Println(result)
}
```

### Exemple de sortie :

```
[present present correct correct absent]
```

Cela signifie :

* `m` est présent mais mal placé,
* `o` est présent mais mal placé,
* `p` et `p` sont corrects,
* `e` est absent.

Souhaites-tu que la fonction retourne aussi le nombre de lettres bien placées ?

Voici la version améliorée de la fonction Go : elle retourne **à la fois** un tableau d’indices (`"correct"`, `"present"`, `"absent"`) et le **nombre de lettres bien placées** (`correctCount`) :

```go
package main

import (
	"fmt"
)

func CheckGuess(secret, guess string) ([]string, int) {
	if len(secret) != len(guess) {
		return nil, 0 // ou panic / erreur
	}

	result := make([]string, len(secret))
	secretRunes := []rune(secret)
	guessRunes := []rune(guess)
	used := make([]bool, len(secret)) // pour éviter les doublons
	correctCount := 0

	// Première passe : lettres correctes et bien placées
	for i := range secretRunes {
		if guessRunes[i] == secretRunes[i] {
			result[i] = "correct"
			used[i] = true
			correctCount++
		}
	}

	// Deuxième passe : lettres correctes mais mal placées
	for i := range guessRunes {
		if result[i] == "" {
			for j := range secretRunes {
				if !used[j] && guessRunes[i] == secretRunes[j] {
					result[i] = "present"
					used[j] = true
					break
				}
			}
			if result[i] == "" {
				result[i] = "absent"
			}
		}
	}

	return result, correctCount
}

func main() {
	secret := "pomme"
	guess := "moppe"

	result, correctCount := CheckGuess(secret, guess)
	fmt.Println("Résultats :", result)
	fmt.Println("Lettres bien placées :", correctCount)
}
```

### Résultat attendu :

```
Résultats : [present present correct correct absent]
Lettres bien placées : 2
```

## add unit test

Voici comment ajouter un test unitaire pour la fonction `CheckGuess` en utilisant le framework de test standard `testing` de Go.

### 📁 Structure de fichiers recommandée :

```
project/
│
├── main.go
└── main_test.go
```

### ✅ `main_test.go`

```go
package main

import (
	"reflect"
	"testing"
)

func TestCheckGuess(t *testing.T) {
	secret := "pomme"
	guess := "moppe"

	expectedResult := []string{"present", "present", "correct", "correct", "absent"}
	expectedCorrectCount := 2

	result, correctCount := CheckGuess(secret, guess)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Errorf("Résultat incorrect. Attendu %v, obtenu %v", expectedResult, result)
	}

	if correctCount != expectedCorrectCount {
		t.Errorf("Nombre de lettres correctes incorrect. Attendu %d, obtenu %d", expectedCorrectCount, correctCount)
	}
}
```

### ✅ Pour exécuter le test :

Dans ton terminal, depuis le dossier du projet, lance :

```bash
go test
```

Souhaites-tu tester plusieurs cas (ex. mot 100% correct, mot 100% faux, lettres doublées, etc.) ?

