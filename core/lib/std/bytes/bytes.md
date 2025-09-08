La librairie standard `bytes` en Go fournit des fonctions pour manipuler des tranches de **bytes** (`[]byte`). Elle est très utile pour le traitement de texte ou de données binaires, notamment lorsque tu veux éviter les conversions constantes entre `string` et `[]byte`.

---

### 📦 Importation

```go
import "bytes"
```

---

### 🔧 Fonctions principales

Voici les fonctions et types les plus utilisés dans `bytes` :

---

#### 1. **Comparaison**

```go
bytes.Compare(a, b []byte) int
```

* Renvoie `0` si `a == b`
* `< 0` si `a < b`
* `> 0` si `a > b`

---

#### 2. **Égalité**

```go
bytes.Equal(a, b []byte) bool
```

* Compare deux slices de bytes (plus rapide que `string(a) == string(b)`).

---

#### 3. **Recherche**

```go
bytes.Contains(s, sub []byte) bool
```

* Vérifie si `sub` est contenu dans `s`.

```go
bytes.Index(s, sub []byte) int
```

* Renvoie la première position de `sub` dans `s`, ou `-1` si non trouvé.

---

#### 4. **Découpage et jointure**

```go
bytes.Split(s, sep []byte) [][]byte
```

* Découpe `s` par le séparateur `sep`.

```go
bytes.Join(slices [][]byte, sep []byte) []byte
```

* Joint plusieurs slices en un seul avec un séparateur.

---

#### 5. **Trimming**

```go
bytes.Trim(s []byte, cutset string) []byte
```

* Supprime les caractères en début et fin qui appartiennent à `cutset`.

```go
bytes.TrimSpace(s []byte) []byte
```

* Supprime les espaces en début et fin.

---

#### 6. **Remplacement**

```go
bytes.Replace(s, old, new []byte, n int) []byte
```

* Remplace `n` occurrences de `old` par `new`. `n = -1` pour toutes.

---

### 📘 Type utile : `bytes.Buffer`

```go
var buf bytes.Buffer
buf.WriteString("Hello, ")
buf.Write([]byte("world!"))
fmt.Println(buf.String()) // "Hello, world!"
```

* Un tampon efficace pour construire dynamiquement des chaînes ou données binaires.

---

### ✅ Exemple complet

```go
package main

import (
    "bytes"
    "fmt"
)

func main() {
    data := []byte("hello world")
    
    if bytes.Contains(data, []byte("world")) {
        fmt.Println("Contient 'world'")
    }

    upper := bytes.Replace(data, []byte("world"), []byte("Golang"), 1)
    fmt.Println(string(upper)) // "hello Golang"

    parts := bytes.Split(upper, []byte(" "))
    for _, part := range parts {
        fmt.Println(string(part))
    }
}
```

---

Souhaites-tu un exemple d'utilisation dans un contexte API ou manipulation de fichiers ?

Parfait ! Voici un **exemple simple d’utilisation du package `bytes` dans un contexte API REST en Go avec Gin**. On va simuler un endpoint qui **traite des données brutes (`[]byte`)** envoyées dans une requête POST (par exemple, du texte brut ou un JSON qu'on veut filtrer ou manipuler au niveau bas).

---

### 🎯 Objectif

Créer un endpoint `/clean` qui :

1. Reçoit du texte brut (`Content-Type: text/plain`)
2. Nettoie les espaces superflus
3. Convertit tout en minuscules
4. Remplace certains mots (ex. `"golang"` → `"Go"`)

---

### 📦 Exemple complet

```go
package main

import (
    "bytes"
    "fmt"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.POST("/clean", func(c *gin.Context) {
        // Lire le corps brut de la requête
        body, err := c.GetRawData()
        if err != nil {
            c.String(http.StatusBadRequest, "Erreur lecture corps: %v", err)
            return
        }

        // Nettoyage : trim des espaces
        cleaned := bytes.TrimSpace(body)

        // Mise en minuscules (en passant temporairement par string)
        cleanedLower := bytes.ToLower(cleaned)

        // Remplacer "golang" par "Go"
        final := bytes.ReplaceAll(cleanedLower, []byte("golang"), []byte("Go"))

        c.Data(http.StatusOK, "text/plain", final)
    })

    r.Run(":8080")
}
```

---

### 🧪 Test avec `curl`

```bash
curl -X POST http://localhost:8080/clean \
  -H "Content-Type: text/plain" \
  --data "    GOLANG est génial !     "
```

**Réponse :**

```
Go est génial !
```

---

### ✅ Pourquoi utiliser `bytes` ici ?

* **Performance** : on évite les conversions inutiles entre `[]byte` et `string`.
* **Simplicité** : les fonctions de `bytes` sont pratiques pour du texte brut.
* **Cas d’usage réel** : API qui reçoivent des logs, du texte, ou du contenu non JSON.

Souhaites-tu voir le même type de logique appliqué à un traitement de fichiers ?

Parfait, voyons maintenant un **exemple de traitement de fichiers avec `bytes` dans un contexte d’API en Go**. L’objectif est de recevoir un **fichier texte via un upload**, puis de :

1. Lire son contenu (`multipart/form-data`)
2. Nettoyer le texte : supprimer les espaces en début/fin, convertir en minuscules, remplacer certains mots
3. Retourner le texte transformé dans la réponse

---

### 📁 Contexte : API `POST /upload`

Le client envoie un fichier `.txt`, et l’API renvoie le contenu modifié.

---

### ✅ Exemple complet

```go
package main

import (
    "bytes"
    "fmt"
    "io"
    "net/http"

    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    r.POST("/upload", func(c *gin.Context) {
        // Récupère le fichier depuis le champ "file"
        file, _, err := c.Request.FormFile("file")
        if err != nil {
            c.String(http.StatusBadRequest, "Erreur fichier : %v", err)
            return
        }
        defer file.Close()

        // Lit tout le contenu du fichier
        content, err := io.ReadAll(file)
        if err != nil {
            c.String(http.StatusInternalServerError, "Erreur lecture : %v", err)
            return
        }

        // Traitement du contenu avec le package bytes
        cleaned := bytes.TrimSpace(content)
        cleaned = bytes.ToLower(cleaned)
        cleaned = bytes.ReplaceAll(cleaned, []byte("golang"), []byte("Go"))

        // Réponse
        c.Data(http.StatusOK, "text/plain", cleaned)
    })

    r.Run(":8080")
}
```

---

### 🧪 Test avec `curl` :

Supposons que tu as un fichier `sample.txt` contenant :

```
    GOLANG est puissant.
    Merci GOLANG !
```

Commandes :

```bash
curl -X POST http://localhost:8080/upload \
  -F "file=@sample.txt"
```

**Réponse :**

```
go est puissant.
merci go !
```

---

### 🔍 Remarques

* `bytes.TrimSpace` nettoie les débuts/fins de ligne.
* `bytes.ReplaceAll` est très utile pour des traitements ciblés (nettoyage, censure, aliasing).
* Ce pattern est courant dans les microservices de **traitement de fichiers**, **import/export**, **filtrage de contenu**, etc.
