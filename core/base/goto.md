En Go, le mot-clé `goto` permet de sauter l'exécution du programme à une **étiquette** définie plus loin ou plus haut dans le même bloc de fonction. C'est un mécanisme de saut inconditionnel.

### Syntaxe

```go
goto label

...

label:
    // instructions
```

### Exemple simple

```go
package main

import "fmt"

func main() {
    i := 0

start:
    if i >= 3 {
        return
    }
    fmt.Println(i)
    i++
    goto start
}
```

### Explication

* `goto start` saute à l’étiquette `start:`.
* C’est un moyen de faire des boucles manuelles (même si en Go on préfère `for`).
* Ce code affiche :

  ```
  0
  1
  2
  ```

---

### À éviter sauf cas particuliers

L’utilisation de `goto` est **fortement déconseillée** dans la majorité des cas car :

* Elle nuit à la lisibilité du code.
* Elle rend plus difficile la maintenance et la compréhension.

### Cas où `goto` peut être utile

* Sortie d’une fonction en cas d’erreur au milieu de beaucoup de code (par exemple, pour faire un "cleanup").

#### Exemple avec nettoyage :

```go
func doSomething() error {
    f, err := os.Open("file.txt")
    if err != nil {
        return err
    }

    // faire quelque chose avec f
    if someCondition {
        goto cleanup
    }

    // autre traitement
    // ...

cleanup:
    f.Close()
    return nil
}
```

Mais même dans ces cas, l’utilisation de `defer` est **souvent préférable** :

```go
f, err := os.Open("file.txt")
if err != nil {
    return err
}
defer f.Close()
```

Souhaites-tu un exemple dans un vrai cas pratique avec `goto`, ou tu veux savoir dans quel cas précis tu pourrais t’en servir ?

Voici un exemple **réaliste** d’utilisation de `goto` dans une fonction complexe, pour **éviter la duplication de code de nettoyage** dans une fonction où plusieurs erreurs peuvent survenir :

### 🧪 Cas pratique : Initialisation avec ressources multiples

```go
package main

import (
    "errors"
    "fmt"
)

func process() error {
    // Imaginons plusieurs étapes avec des ressources
    var connOpened bool
    var fileOpened bool

    // Étape 1 : ouvrir une connexion
    if err := openConnection(); err != nil {
        return err
    }
    connOpened = true

    // Étape 2 : ouvrir un fichier
    if err := openFile(); err != nil {
        goto cleanup
    }
    fileOpened = true

    // Étape 3 : traiter des données
    if err := processData(); err != nil {
        goto cleanup
    }

    fmt.Println("Traitement réussi !")
    return nil

cleanup:
    if fileOpened {
        fmt.Println("Fermeture fichier")
    }
    if connOpened {
        fmt.Println("Fermeture connexion")
    }
    return errors.New("échec pendant le traitement")
}

func openConnection() error {
    fmt.Println("Connexion ouverte")
    return nil
}

func openFile() error {
    fmt.Println("Fichier ouvert")
    return nil
}

func processData() error {
    fmt.Println("Erreur pendant le traitement des données")
    return errors.New("erreur de traitement")
}

func main() {
    err := process()
    if err != nil {
        fmt.Println("Erreur:", err)
    }
}
```

### 🧾 Sortie :

```
Connexion ouverte
Fichier ouvert
Erreur pendant le traitement des données
Fermeture fichier
Fermeture connexion
Erreur: échec pendant le traitement
```

---

### Pourquoi utiliser `goto` ici ?

* Pour **éviter d’appeler deux fois** le même bloc de nettoyage (`fermer fichier`, `fermer connexion`) à plusieurs endroits.
* Mais en Go, **`defer` est préféré** si chaque ressource peut être fermée directement après l'ouverture.

Souhaites-tu que je te montre **la version équivalente avec `defer`** pour comparaison ?

Voici **la version équivalente avec `defer`**, qui est plus idiomatique en Go (et plus lisible) :

### ✅ Version propre avec `defer`

```go
package main

import (
    "errors"
    "fmt"
)

func process() error {
    // Étape 1 : ouvrir une connexion
    if err := openConnection(); err != nil {
        return err
    }
    defer fmt.Println("Fermeture connexion") // sera exécuté à la fin

    // Étape 2 : ouvrir un fichier
    if err := openFile(); err != nil {
        return err
    }
    defer fmt.Println("Fermeture fichier")

    // Étape 3 : traiter des données
    if err := processData(); err != nil {
        return err
    }

    fmt.Println("Traitement réussi !")
    return nil
}

func openConnection() error {
    fmt.Println("Connexion ouverte")
    return nil
}

func openFile() error {
    fmt.Println("Fichier ouvert")
    return nil
}

func processData() error {
    fmt.Println("Erreur pendant le traitement des données")
    return errors.New("erreur de traitement")
}

func main() {
    err := process()
    if err != nil {
        fmt.Println("Erreur:", err)
    }
}
```

### 🧾 Sortie :

```
Connexion ouverte
Fichier ouvert
Erreur pendant le traitement des données
Fermeture fichier
Fermeture connexion
Erreur: erreur de traitement
```

---

### 🧠 Avantages de `defer`

* Le code est plus clair et **évite les `goto`**.
* Moins de duplication.
* Le nettoyage est **lié à l’allocation**, ce qui réduit les erreurs.

Utiliser `goto` peut encore avoir un intérêt dans certains cas très particuliers (performances critiques, code généré, etc.), mais dans 99% des cas, **`defer` est préférable** en Go.

Souhaites-tu un petit projet ou exemple en Go où tu pourrais tester cela toi-même ?
