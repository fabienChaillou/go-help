Le package `embed` de la bibliothèque standard de Go permet **d'inclure des fichiers statiques dans un binaire** lors de la compilation. Cela signifie que tu peux embarquer des fichiers (comme des templates HTML, fichiers CSS/JS, images, fichiers de configuration, etc.) directement dans l'exécutable Go, sans avoir à les lire depuis le disque au runtime.

### 🔧 Utilisation de base

Voici un exemple simple :

```go
package main

import (
    "embed"
    "fmt"
)

//go:embed hello.txt
var content string

func main() {
    fmt.Println(content)
}
```

> Ce code va embarquer le fichier `hello.txt` dans la variable `content`.

---

### 📁 Types supportés

Tu peux utiliser `embed` avec **trois types de variables** :

1. `string` → pour des fichiers texte
2. `[]byte` → pour des fichiers binaires
3. `embed.FS` → pour accéder à un système de fichiers virtuel (pour plusieurs fichiers ou répertoires)

---

### 📂 Exemple avec `embed.FS` (multi-fichiers)

```go
package main

import (
    "embed"
    "fmt"
    "io/fs"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
    data, err := staticFiles.ReadFile("static/test.txt")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
}
```

> Tous les fichiers sous `static/` seront embarqués, et tu peux les lire avec les méthodes de `fs.FS`.

---

### ⚠️ Contraintes

* Les directives `//go:embed` **doivent être placées juste au-dessus de la variable** à laquelle elles s’appliquent.
* La variable doit être de type `string`, `[]byte`, ou `embed.FS` uniquement.
* Les fichiers doivent exister au moment de la compilation.

---

### ✅ Cas d'usage typiques

* Templates HTML pour un serveur web
* Fichiers statiques (CSS, JS, images)
* Fichiers de configuration par défaut
* Embedding de contenu pour des outils CLI portables

Souhaites-tu un exemple dans un projet complet (ex : serveur web Gin qui sert des fichiers statiques embarqués) ?

### Just a example with txt file and the tests

Voici un exemple minimaliste qui montre comment embarquer un fichier `.txt` avec le package `embed`, et comment écrire un test pour vérifier que le contenu est bien chargé.

---

### 📁 Structure du projet

```
example/
├── data.txt
├── main.go
└── main_test.go
```

---

### 📄 `data.txt`

```txt
Bonjour depuis Go embed!
```

---

### 📄 `main.go`

```go
package main

import (
	"embed"
	"fmt"
)

//go:embed data.txt
var data string

func GetData() string {
	return data
}

func main() {
	fmt.Println(GetData())
}
```

---

### 📄 `main_test.go`

```go
package main

import (
	"strings"
	"testing"
)

func TestGetData(t *testing.T) {
	content := GetData()

	if !strings.Contains(content, "Bonjour") {
		t.Errorf("Le contenu ne contient pas 'Bonjour' : %s", content)
	}

	if len(content) == 0 {
		t.Error("Le contenu est vide")
	}
}
```

---

### ✅ Pour tester

Lance simplement :

```bash
go test
```

Et tu devrais voir un test qui passe si `data.txt` contient bien `"Bonjour depuis Go embed!"`.

Souhaites-tu un exemple équivalent avec un `[]byte` ou `embed.FS` à la place ?

Parfait ! Voici le **même exemple**, mais cette fois en utilisant `[]byte` à la place de `string`.

---

### 📁 Structure

```
example/
├── data.txt
├── main.go
└── main_test.go
```

---

### 📄 `data.txt`

```txt
Bonjour depuis Go embed!
```

---

### 📄 `main.go`

```go
package main

import (
	"embed"
	"fmt"
)

//go:embed data.txt
var data []byte

func GetData() []byte {
	return data
}

func main() {
	fmt.Println(string(GetData()))
}
```

---

### 📄 `main_test.go`

```go
package main

import (
	"bytes"
	"testing"
)

func TestGetData(t *testing.T) {
	content := GetData()

	if len(content) == 0 {
		t.Error("Le contenu est vide")
	}

	if !bytes.Contains(content, []byte("Bonjour")) {
		t.Errorf("Le contenu ne contient pas 'Bonjour' : %s", content)
	}
}
```

---

### ✅ Exécution du test

```bash
go test
```

Ce test vérifie que :

* Le contenu n’est pas vide
* Il contient bien le mot `"Bonjour"`

Souhaites-tu que je t’ajoute aussi un exemple avec `embed.FS` et plusieurs fichiers ?

Voici un exemple avec `embed.FS`, qui permet d’embarquer **plusieurs fichiers** ou même des **dossiers entiers**.

---

### 📁 Structure

```
example/
├── assets/
│   ├── a.txt
│   └── b.txt
├── main.go
└── main_test.go
```

---

### 📄 `assets/a.txt`

```txt
Contenu du fichier A
```

### 📄 `assets/b.txt`

```txt
Contenu du fichier B
```

---

### 📄 `main.go`

```go
package main

import (
	"embed"
	"io/fs"
)

//go:embed assets/*
var content embed.FS

func GetFileContent(name string) (string, error) {
	data, err := content.ReadFile("assets/" + name)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func ListFiles() ([]fs.DirEntry, error) {
	return content.ReadDir("assets")
}
```

---

### 📄 `main_test.go`

```go
package main

import (
	"strings"
	"testing"
)

func TestGetFileContent(t *testing.T) {
	content, err := GetFileContent("a.txt")
	if err != nil {
		t.Fatalf("Erreur de lecture de a.txt: %v", err)
	}
	if !strings.Contains(content, "fichier A") {
		t.Errorf("Le contenu ne contient pas 'fichier A' : %s", content)
	}
}

func TestListFiles(t *testing.T) {
	files, err := ListFiles()
	if err != nil {
		t.Fatalf("Erreur lors de la lecture du répertoire assets: %v", err)
	}

	if len(files) < 2 {
		t.Errorf("Il manque des fichiers dans assets/: %d trouvés", len(files))
	}
}
```

---

### ✅ Pour exécuter

Assure-toi que tu es dans le dossier `example`, puis :

```bash
go test
```

Ce test vérifie :

* Que `a.txt` est bien embarqué et contient le texte attendu.
* Que les fichiers dans `assets/` sont bien listés via `ReadDir`.

Souhaites-tu un exemple avec lecture récursive ou intégration à un serveur web (Gin, par ex.) ?
