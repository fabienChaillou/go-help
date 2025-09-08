L’interopérabilité avec C en Go se fait via un mécanisme appelé **`cgo`**. Il te permet d’appeler du code écrit en **C** directement depuis du code Go.

C’est très utile lorsque :

* Tu veux **réutiliser une bibliothèque C existante**.
* Tu as besoin de **fonctionnalités système bas-niveau** non exposées en Go.
* Tu veux **accélérer une portion critique** de ton programme avec du code optimisé C.

---

## ⚙️ Comment ça marche ? (`cgo`)

Go permet d’écrire du code comme ceci :

```go
/*
#include <stdlib.h>
#include <string.h>

int add(int a, int b) {
    return a + b;
}
*/
import "C"

import "fmt"

func main() {
    result := C.add(2, 3)
    fmt.Println("Résultat:", result)
}
```

---

## 📌 Points clés

### 1. **Bloc `import "C"`**

* Tu places ton code C juste **au-dessus** de cette ligne, dans un commentaire spécial.
* `go build` reconnaît cela et génère automatiquement les bindings.

### 2. **Appels C en Go**

* Les fonctions et types C deviennent accessibles via le préfixe `C.` :

  * `C.int`, `C.double`, `C.add()`, etc.

### 3. **Conversion de types**

Go et C ont des systèmes de types différents. Il faut convertir :

```go
var a C.int = C.int(5)
var b int = int(a)
```

---

## 🧠 Exemples de gestion de mémoire

### Côté Go → C

```go
str := C.CString("hello") // Go string → *C char
defer C.free(unsafe.Pointer(str)) // Libération mémoire
```

### C côté Go

```go
#include <stdlib.h>
char* greet() {
    return "hello from C";
}
```

```go
msg := C.greet()
fmt.Println(C.GoString(msg)) // Convertit *C.char → Go string
```

---

## 🛑 Limitations

* **`cgo` ralentit un peu la compilation.**
* Le code n’est plus 100 % portable (dépend du compilateur C).
* L’interaction Go ↔ C ajoute une **barrière de performances** (passage du GC, etc.).
* Plus complexe à déboguer.

---

## ✅ Astuce : éviter `cgo` si possible

Go fournit parfois des bindings "purs Go" aux bibliothèques C (comme SQLite, libgit2, etc.). Cela évite `cgo` et garde ton binaire **100 % statique** et **portable**.

---

Exemple concret avec `libcurl` en Go via `cgo`) ?

Parfait ! Voici un **exemple concret** d’utilisation de **`libcurl` via `cgo` en Go**.

---

## 📦 Objectif

Faire une **requête HTTP GET** en utilisant **`libcurl`**, une bibliothèque C populaire pour les requêtes réseau.

---

## 🧱 Prérequis

Tu dois avoir `libcurl` installé sur ton système :

### Sous Debian/Ubuntu :

```bash
sudo apt install libcurl4-openssl-dev
```

### Sous macOS (Homebrew) :

```bash
brew install curl
```

---

## 📄 Exemple de code Go + C (libcurl)

```go
package main

/*
#cgo LDFLAGS: -lcurl
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>

// Fonction de rappel pour stocker les données reçues
size_t write_callback(void *contents, size_t size, size_t nmemb, void *userp) {
    size_t realsize = size * nmemb;
    strncat(userp, contents, realsize);
    return realsize;
}

char* perform_request(const char* url) {
    CURL *curl;
    CURLcode res;

    static char buffer[10000]; // petite taille pour l'exemple
    buffer[0] = '\0';

    curl = curl_easy_init();
    if(curl) {
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, buffer);

        res = curl_easy_perform(curl);
        curl_easy_cleanup(curl);

        if (res != CURLE_OK) {
            return "Error";
        }
        return buffer;
    }
    return "Failed to initialize curl";
}
*/
import "C"

import (
    "fmt"
    "unsafe"
)

func main() {
    url := C.CString("https://httpbin.org/get")
    defer C.free(unsafe.Pointer(url))

    response := C.perform_request(url)
    fmt.Println(C.GoString(response))
}
```

---

## ⚙️ Explication rapide

* `#cgo LDFLAGS: -lcurl` : indique à Go de **lier libcurl**.
* `perform_request` est une fonction C appelée depuis Go.
* `write_callback` écrit les données reçues dans un buffer.
* `C.GoString(response)` convertit le `char*` en string Go.

---

## 🧠 Limites / À améliorer

* Le buffer C est **statique** → non sûr si tu fais plusieurs requêtes ou gros contenu.
* Pas de gestion fine des erreurs HTTP ou des en-têtes.
* `strncat` est dangereux sans vérification de taille restante (à remplacer par gestion dynamique pour du vrai usage).

---

Souhaites-tu que je te montre une version plus robuste avec allocation dynamique ou un wrapper Go plus sécurisé autour de `libcurl` ?

Parfait ! Voici une **version plus robuste** de l’exemple précédent, avec :

* 🧠 **allocation dynamique de mémoire** pour stocker la réponse,
* ✅ une **libération propre** de la mémoire côté Go,
* 📦 et un **wrapper Go** pour encapsuler l’appel C proprement.

---

## 🧱 Version améliorée avec allocation dynamique (`malloc` / `realloc`)

### 🧬 Code Go avec bloc C

```go
package main

/*
#cgo LDFLAGS: -lcurl
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>

// Structure pour accumuler les données
struct Buffer {
    char *data;
    size_t size;
};

// Callback pour écrire les données dans le buffer dynamique
size_t write_callback(void *contents, size_t size, size_t nmemb, void *userp) {
    size_t realsize = size * nmemb;
    struct Buffer *mem = (struct Buffer *)userp;

    char *ptr = realloc(mem->data, mem->size + realsize + 1);
    if(ptr == NULL) return 0; // erreur d'allocation

    mem->data = ptr;
    memcpy(&(mem->data[mem->size]), contents, realsize);
    mem->size += realsize;
    mem->data[mem->size] = '\0';

    return realsize;
}

// Fonction principale pour faire la requête GET
char* curl_get(const char* url) {
    CURL *curl;
    CURLcode res;

    struct Buffer chunk;
    chunk.data = malloc(1); // initialisation
    chunk.size = 0;

    curl = curl_easy_init();
    if(curl) {
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&chunk);

        res = curl_easy_perform(curl);
        curl_easy_cleanup(curl);

        if(res != CURLE_OK) {
            free(chunk.data);
            return NULL;
        }

        return chunk.data;
    }

    return NULL;
}
*/
import "C"

import (
    "fmt"
    "unsafe"
)

func main() {
    url := C.CString("https://httpbin.org/get")
    defer C.free(unsafe.Pointer(url))

    response := C.curl_get(url)
    if response == nil {
        fmt.Println("Erreur lors de la requête")
        return
    }
    defer C.free(unsafe.Pointer(response))

    fmt.Println("Réponse de libcurl :")
    fmt.Println(C.GoString(response))
}
```

---

## ✅ Avantages de cette version

* **Pas de buffer statique**, on alloue dynamiquement en fonction des données reçues.
* **Sécurité mémoire** : on libère proprement avec `C.free`.
* Peut gérer des **réponses plus grandes** ou être appelé plusieurs fois sans problème.

---

## ⚠️ Astuces & conseils

* Si tu veux en faire une **bibliothèque plus propre**, isole `curl_get` dans un fichier `.c` séparé.
* Pour la production, ajoute aussi :

  * Support des **timeouts** (`CURLOPT_TIMEOUT`)
  * Gestion des **en-têtes**, **POST**, **auth**, etc.

---

Souhaites-tu que je t’aide à créer un **wrapper Go plus "Go-like"** autour de ce code, par exemple une fonction comme `func CurlGet(url string) (string, error)` ?

Parfait ! Voici un **wrapper Go propre et idiomatique** autour de `libcurl`, qui expose une fonction :

```go
func CurlGet(url string) (string, error)
```

---

## ✅ Résultat attendu

Tu pourras utiliser `libcurl` comme ceci, dans du code Go natif :

```go
body, err := CurlGet("https://httpbin.org/get")
if err != nil {
    log.Fatal(err)
}
fmt.Println(body)
```

---

## 🧱 Implémentation complète

### 📁 `main.go`

```go
package main

/*
#cgo LDFLAGS: -lcurl
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>

struct Buffer {
    char *data;
    size_t size;
};

size_t write_callback(void *contents, size_t size, size_t nmemb, void *userp) {
    size_t realsize = size * nmemb;
    struct Buffer *mem = (struct Buffer *)userp;

    char *ptr = realloc(mem->data, mem->size + realsize + 1);
    if(ptr == NULL) return 0;

    mem->data = ptr;
    memcpy(&(mem->data[mem->size]), contents, realsize);
    mem->size += realsize;
    mem->data[mem->size] = '\0';

    return realsize;
}

char* curl_get(const char* url, char** errbuf) {
    CURL *curl;
    CURLcode res;

    struct Buffer chunk;
    chunk.data = malloc(1);
    chunk.size = 0;

    curl = curl_easy_init();
    if(curl) {
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&chunk);
        curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);

        // Erreur détaillée
        *errbuf = malloc(CURL_ERROR_SIZE);
        curl_easy_setopt(curl, CURLOPT_ERRORBUFFER, *errbuf);

        res = curl_easy_perform(curl);
        curl_easy_cleanup(curl);

        if(res != CURLE_OK) {
            free(chunk.data);
            return NULL;
        }

        return chunk.data;
    }

    return NULL;
}
*/
import "C"

import (
    "errors"
    "unsafe"
)

// CurlGet fait une requête HTTP GET via libcurl et retourne le corps de la réponse ou une erreur
func CurlGet(url string) (string, error) {
    cURL := C.CString(url)
    defer C.free(unsafe.Pointer(cURL))

    var errBuf *C.char

    cResponse := C.curl_get(cURL, &errBuf)
    if cResponse == nil {
        defer C.free(unsafe.Pointer(errBuf))
        return "", errors.New(C.GoString(errBuf))
    }
    defer C.free(unsafe.Pointer(cResponse))

    return C.GoString(cResponse), nil
}

func main() {
    body, err := CurlGet("https://httpbin.org/get")
    if err != nil {
        panic(err)
    }
    println("Réponse reçue :")
    println(body)
}
```

---

## 💡 Bonus : Pourquoi faire un wrapper ?

* Tu **encapsules `cgo`** dans une fonction propre à utiliser.
* Tu gardes ton application **Go idiomatique** (erreurs, chaînes, sécurité mémoire).
* Tu peux réutiliser ce code dans n’importe quel projet.

---

Souhaites-tu que je transforme cela en **module Go** ou en **paquet réutilisable (`curlwrapper`)** ?

Parfait, on va transformer ce code en **module Go réutilisable** appelé `curlwrapper`, avec une API simple : `curlwrapper.CurlGet(url string) (string, error)`.

---

## 🗂️ Structure du projet

Voici la structure de base :

```
curlwrapper/
├── go.mod
├── curl.go
├── curl.c
```

---

## 1. 📄 `go.mod`

Crée un nouveau module :

```bash
go mod init github.com/tonutilisateur/curlwrapper
```

---

## 2. 📄 `curl.c` (code C isolé)

```c
// curl.c
#include <stdlib.h>
#include <string.h>
#include <curl/curl.h>

struct Buffer {
    char *data;
    size_t size;
};

size_t write_callback(void *contents, size_t size, size_t nmemb, void *userp) {
    size_t realsize = size * nmemb;
    struct Buffer *mem = (struct Buffer *)userp;

    char *ptr = realloc(mem->data, mem->size + realsize + 1);
    if (ptr == NULL) return 0;

    mem->data = ptr;
    memcpy(&(mem->data[mem->size]), contents, realsize);
    mem->size += realsize;
    mem->data[mem->size] = '\0';

    return realsize;
}

char* curl_get(const char* url, char** errbuf) {
    CURL *curl;
    CURLcode res;

    struct Buffer chunk;
    chunk.data = malloc(1);
    chunk.size = 0;

    curl = curl_easy_init();
    if (curl) {
        curl_easy_setopt(curl, CURLOPT_URL, url);
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, write_callback);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&chunk);
        curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);

        *errbuf = malloc(CURL_ERROR_SIZE);
        curl_easy_setopt(curl, CURLOPT_ERRORBUFFER, *errbuf);

        res = curl_easy_perform(curl);
        curl_easy_cleanup(curl);

        if (res != CURLE_OK) {
            free(chunk.data);
            return NULL;
        }

        return chunk.data;
    }

    return NULL;
}
```

---

## 3. 📄 `curl.go` (wrapper Go)

```go
package curlwrapper

/*
#cgo LDFLAGS: -lcurl
#include "curl.c"
*/
import "C"
import (
    "errors"
    "unsafe"
)

// CurlGet effectue une requête HTTP GET via libcurl et retourne la réponse ou une erreur
func CurlGet(url string) (string, error) {
    cURL := C.CString(url)
    defer C.free(unsafe.Pointer(cURL))

    var errBuf *C.char

    cResponse := C.curl_get(cURL, &errBuf)
    if cResponse == nil {
        defer C.free(unsafe.Pointer(errBuf))
        return "", errors.New(C.GoString(errBuf))
    }
    defer C.free(unsafe.Pointer(cResponse))

    return C.GoString(cResponse), nil
}
```

---

## ✅ Utilisation dans un autre projet

Tu peux maintenant utiliser ton module dans un autre projet :

```go
package main

import (
    "fmt"
    "log"

    "github.com/tonutilisateur/curlwrapper"
)

func main() {
    body, err := curlwrapper.CurlGet("https://httpbin.org/get")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(body)
}
```

---

## 🚀 Astuces

* Tu peux publier ce module sur GitHub pour l'importer directement.
* Ajoute un `README.md` et un `go.mod tidy` pour nettoyer les dépendances.
* Pour les tests, tu peux simuler des URLs avec des serveurs locaux ou `httpbin`.

---

Souhaites-tu que je te génère un `README.md` pour ce module avec exemple et instructions ?

Voici un fichier **`README.md`** prêt à l'emploi pour ton module Go `curlwrapper` utilisant `libcurl` via `cgo` :

---

### 📄 `README.md`

````markdown
# curlwrapper

🌀 `curlwrapper` est un petit wrapper Go autour de [libcurl](https://curl.se/libcurl/), utilisant `cgo`, pour effectuer des requêtes HTTP GET avec toutes les capacités de libcurl.

---

## ✨ Fonctionnalités

- Utilise directement **libcurl** via `cgo`
- Prend en charge les redirections HTTP
- Allocation mémoire dynamique pour stocker le corps de la réponse
- API simple : `CurlGet(url string) (string, error)`

---

## 🛠️ Installation

Assurez-vous que **libcurl est installé** :

### Linux (Debian/Ubuntu)
```bash
sudo apt install libcurl4-openssl-dev
````

### macOS (Homebrew)

```bash
brew install curl
```

### Initialiser le module Go

```bash
go get github.com/tonutilisateur/curlwrapper
```

---

## 📦 Utilisation

```go
package main

import (
    "fmt"
    "log"

    "github.com/tonutilisateur/curlwrapper"
)

func main() {
    body, err := curlwrapper.CurlGet("https://httpbin.org/get")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(body)
}
```

---

## 🧬 Exemple d’appel natif

Le module encapsule cet appel `cgo` :

```go
func CurlGet(url string) (string, error)
```

Tu peux donc facilement intégrer `libcurl` à une application Go sans manipuler directement du code C.

---

## ⚠️ Notes

* Le code utilise `cgo`, donc :

  * Non compatible avec `go run` en mode `-trimpath`
  * Lien dynamique requis avec `libcurl` à l’exécution
* Pour un binaire statique pur Go, privilégier `net/http` (mais moins flexible que libcurl).

---

## 📜 Licence

MIT – libre à utiliser, modifier, et distribuer.

```

---

Souhaites-tu aussi que je t’aide à publier ce module sur GitHub ou à configurer un test unitaire minimal ?
```

