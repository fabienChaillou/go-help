Le package `bufio` en Go (dans `bufio` pour *buffered I/O*) fournit des interfaces de lecture et d’écriture avec *buffering*, ce qui améliore les performances lors de la lecture/écriture sur des sources lentes comme les fichiers, les connexions réseau ou l'entrée standard.

### Pourquoi utiliser `bufio` ?

Les opérations de lecture/écriture sans buffer (comme avec `os.File.Read`) sont coûteuses car elles font appel directement au système. `bufio` réduit le nombre d'appels système en lisant ou écrivant de grandes quantités de données à la fois, puis en les traitant dans un tampon (*buffer*) en mémoire.

---

### Principales structures de `bufio`

#### `bufio.Reader`

Permet de lire efficacement depuis une source (comme un fichier, un réseau, etc.).

**Fonctions utiles :**

* `Read(p []byte) (n int, err error)`
* `ReadString(delim byte) (string, error)`
* `ReadBytes(delim byte) ([]byte, error)`
* `ReadLine()`
* `Peek(n int) ([]byte, error)` – lit sans avancer le pointeur de lecture.

**Exemple : Lire un fichier ligne par ligne**

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("exemple.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Erreur de lecture :", err)
    }
}
```

---

#### `bufio.Writer`

Permet d’écrire efficacement vers une destination (fichier, socket, etc.).

**Fonctions utiles :**

* `Write(p []byte) (n int, err error)`
* `WriteString(s string) (int, error)`
* `Flush()` – très important : vide le tampon vers la destination.

**Exemple : Écriture avec buffering**

```go
package main

import (
    "bufio"
    "os"
)

func main() {
    file, err := os.Create("exemple.txt")
    if err != nil {
        panic(err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    writer.WriteString("Hello, monde !\n")
    writer.Flush() // obligatoire pour écrire dans le fichier
}
```

---

#### `bufio.Scanner`

Utile pour lire facilement des fichiers ligne par ligne (plus simple que `bufio.Reader` pour ce cas).

**Limite** : ne convient pas aux très grandes lignes (> 64 Ko par défaut, sauf si on change la taille avec `scanner.Buffer()`).

---

### En résumé

| Structure       | Utilisation                                 |
| --------------- | ------------------------------------------- |
| `bufio.Reader`  | Lire efficacement avec des méthodes variées |
| `bufio.Writer`  | Écrire avec un tampon, nécessite `Flush()`  |
| `bufio.Scanner` | Lire facilement ligne par ligne             |

Souhaites-tu un exemple complet d’un programme avec lecture et écriture en utilisant `bufio` ?

Voici un exemple complet d’un programme Go qui :

1. Crée un fichier et y écrit plusieurs lignes avec `bufio.Writer`.
2. Rouvre le fichier et lit son contenu ligne par ligne avec `bufio.Scanner`.

---

### 🔧 Exemple complet avec `bufio.Writer` et `bufio.Scanner`

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    // Étape 1 : Créer et écrire dans un fichier avec bufio.Writer
    fileName := "demo.txt"
    file, err := os.Create(fileName)
    if err != nil {
        fmt.Println("Erreur lors de la création du fichier :", err)
        return
    }

    writer := bufio.NewWriter(file)
    lignes := []string{
        "Bonjour monde",
        "Voici une démonstration de bufio",
        "L'écriture est en buffer",
    }

    for _, ligne := range lignes {
        _, err := writer.WriteString(ligne + "\n")
        if err != nil {
            fmt.Println("Erreur lors de l'écriture :", err)
            return
        }
    }

    // Très important : vider le buffer dans le fichier
    writer.Flush()
    file.Close()

    // Étape 2 : Lire le fichier avec bufio.Scanner
    file, err = os.Open(fileName)
    if err != nil {
        fmt.Println("Erreur lors de l'ouverture du fichier :", err)
        return
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    fmt.Println("Contenu du fichier :")
    for scanner.Scan() {
        fmt.Println(scanner.Text())
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Erreur lors de la lecture :", err)
    }
}
```

---

### ✅ Résultat attendu (dans le terminal) :

```
Contenu du fichier :
Bonjour monde
Voici une démonstration de bufio
L'écriture est en buffer
```

---

Souhaites-tu que je te montre aussi comment lire caractère par caractère ou par mot avec `bufio.Reader` ?

Parfait ! Voici deux exemples pour lire un fichier avec `bufio.Reader`, cette fois :

---

## 🧩 Exemple 1 : Lecture caractère par caractère (`ReadByte`)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("demo.txt")
    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    fmt.Println("Lecture caractère par caractère :")
    for {
        char, err := reader.ReadByte()
        if err != nil {
            break // fin du fichier ou erreur
        }
        fmt.Printf("%c", char)
    }
}
```

---

## 🧩 Exemple 2 : Lecture mot par mot (en utilisant `ReadString(' ')`)

```go
package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("demo.txt")
    if err != nil {
        fmt.Println("Erreur :", err)
        return
    }
    defer file.Close()

    reader := bufio.NewReader(file)
    fmt.Println("\nLecture mot par mot :")
    for {
        word, err := reader.ReadString(' ')
        if err != nil {
            fmt.Print(word) // affiche le dernier mot sans espace
            break
        }
        fmt.Print(word)
    }
}
```
