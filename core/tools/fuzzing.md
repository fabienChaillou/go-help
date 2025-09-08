Le **fuzzing** (ou **fuzz testing**) est une technique de test automatisé qui consiste à injecter des données aléatoires ou semi-aléatoires dans une fonction pour détecter des bugs, des crashs ou des comportements inattendus. En **Go**, le fuzzing est intégré nativement depuis **Go 1.18**, ce qui facilite grandement son utilisation.

---

### ✅ Pourquoi utiliser le fuzzing ?

* Pour détecter des erreurs que des tests unitaires classiques ne couvrent pas.
* Pour trouver des panics, des dépassements de tampon ou des comportements indésirables.
* Pour renforcer la robustesse et la sécurité du code.

---

### 📦 Exemple simple de fuzzing en Go

Supposons qu'on veuille tester une fonction qui inverse une chaîne de caractères :

```go
package reverse

func Reverse(s string) string {
    r := []rune(s)
    for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
        r[i], r[j] = r[j], r[i]
    }
    return string(r)
}
```

On peut écrire un **test de fuzzing** dans un fichier `reverse_test.go` :

```go
package reverse

import "testing"

func FuzzReverse(f *testing.F) {
    // Cas de test de base
    f.Add("hello")
    f.Add("Go")
    f.Add("")

    f.Fuzz(func(t *testing.T, orig string) {
        rev := Reverse(orig)
        doubleRev := Reverse(rev)

        // Vérifie que Reverse(Reverse(s)) == s
        if orig != doubleRev {
            t.Errorf("Expected %q, got %q", orig, doubleRev)
        }
    })
}
```

---

### 🚀 Exécution du fuzzing

Tu peux lancer le fuzzing avec la commande :

```bash
go test -fuzz=Fuzz
```

Cela va :

* Utiliser les cas de base définis avec `f.Add()`,
* Puis générer automatiquement d'autres entrées aléatoires,
* Et tenter de "casser" la fonction testée.

---

### 🧠 Ce qu’il faut savoir

* Le fuzzing en Go se fait uniquement sur des fonctions dont le nom commence par `Fuzz`.
* Il est particulièrement utile pour les fonctions de **parsing**, **encodage**, **cryptographie**, ou **gestion de chaînes ou fichiers**.
* Go stocke les cas qui échouent dans un répertoire de cache (`testdata/fuzz/...`), que tu peux rejouer.

---

Souhaites-tu un exemple plus avancé (ex : sur un parser JSON, ou une API) ?

Parfait ! Voici un **exemple plus avancé** de fuzzing en Go, appliqué à une fonction qui parse du JSON. C'est un cas réaliste et pertinent car le parsing est une source fréquente d'erreurs.

---

## 🔍 Cas : parser un JSON en Go

Supposons que tu aies cette fonction qui parse un JSON contenant un champ `"name"` :

```go
package parser

import "encoding/json"

type Person struct {
	Name string `json:"name"`
}

func ParsePerson(data []byte) (*Person, error) {
	var p Person
	err := json.Unmarshal(data, &p)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
```

---

## 🧪 Fichier de test avec fuzzing (`parser_test.go`)

```go
package parser

import "testing"

func FuzzParsePerson(f *testing.F) {
	// Quelques exemples valides
	f.Add([]byte(`{"name":"Alice"}`))
	f.Add([]byte(`{"name":"Bob"}`))
	f.Add([]byte(`{}`)) // cas sans nom
	f.Add([]byte(`""`)) // JSON invalide

	f.Fuzz(func(t *testing.T, data []byte) {
		person, err := ParsePerson(data)

		// Le fuzzing doit seulement faire échouer si la fonction panique
		// On ne vérifie pas le contenu ici, juste la robustesse
		if err == nil && person == nil {
			t.Errorf("Expected non-nil person when err is nil")
		}
	})
}
```

---

## ▶️ Lancer le fuzzing

Dans le dossier du module :

```bash
go test -fuzz=Fuzz
```

Go générera automatiquement des variations du JSON pour tenter de provoquer :

* des erreurs inattendues,
* des comportements étranges (panic, nil inattendu, etc).

---

## 📂 Cas où c’est utile

* Sécurité (tester des entrées non contrôlées venant de clients).
* Robustesse d’une API REST.
* Libs de parsing, encodage, protocole réseau, etc.

---

Souhaites-tu un exemple encore plus concret, par exemple avec un handler HTTP ou une base de données ?
