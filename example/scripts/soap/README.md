En Go, il n’existe pas de librairie standard native pour **SOAP** (contrairement à HTTP/JSON/REST), mais tu peux interagir avec un serveur SOAP de plusieurs façons :

---

## 🔹 1. Utiliser des requêtes HTTP manuelles

Un service SOAP est essentiellement du **XML envoyé via HTTP POST**.
Tu peux construire toi-même les enveloppes SOAP et les envoyer avec `net/http`.

Exemple simple :

```go
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	url := "http://example.com/soap" // URL du service SOAP
	action := "http://example.com/GetUser" // SOAPAction si requis

	// Exemple d'enveloppe SOAP
	soapBody := `<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		<GetUser xmlns="http://example.com/">
		  <id>123</id>
		</GetUser>
	  </soap:Body>
	</soap:Envelope>`

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(soapBody)))
	if err != nil {
		panic(err)
	}

	// Headers SOAP
	req.Header.Add("Content-Type", "text/xml; charset=utf-8")
	req.Header.Add("SOAPAction", action)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
```

👉 Avantage : contrôle total, pas de dépendance.
👉 Inconvénient : il faut parser le XML toi-même (`encoding/xml`).

---

## 🔹 2. Générer du code Go à partir d’un WSDL

Si ton serveur expose un **WSDL**, tu peux générer des bindings Go.

📦 Outil populaire : [`github.com/hooklift/gowsdl`](https://github.com/hooklift/gowsdl) (WSDL → Go)

```bash
go install github.com/hooklift/gowsdl/cmd/gowsdl@latest
```

Puis :

```bash
gowsdl http://example.com/service?wsdl > client.go
```

Cela génère un client Go que tu peux utiliser directement, par exemple :

```go
package main

import (
	"fmt"
	"log"

	"myproject/soapclient" // généré par gowsdl
)

func main() {
	client := soapclient.NewSomeServicePortType("http://example.com/service", false, nil)
	resp, err := client.GetUser(&soapclient.GetUser{Id: 123})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User: %+v\n", resp)
}
```

👉 Avantage : tu manipules des structs Go au lieu de bricoler du XML.
👉 Inconvénient : `gowsdl` n’est plus très actif, mais reste utile.

---

## 🔹 3. Utiliser une librairie SOAP en Go

Quelques packages communautaires :

* [`github.com/tiaguinho/gosoap`](https://github.com/tiaguinho/gosoap) → envoi SOAP simplifié avec parsing XML.
* [`github.com/hooklift/gowsdl`](https://github.com/hooklift/gowsdl) → génération de client depuis WSDL.

Exemple avec **gosoap** :

```go
package main

import (
	"fmt"
	"log"

	"github.com/tiaguinho/gosoap"
)

func main() {
	client, err := gosoap.SoapClient("http://example.com/service?wsdl")
	if err != nil {
		log.Fatal(err)
	}

	params := gosoap.Params{
		"id": 123,
	}

	res, err := client.Call("GetUser", params)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(res)
}
```

---

✅ **Résumé :**

* Si tu veux du **rapide et générique** → HTTP + `encoding/xml`.
* Si tu as un **WSDL** bien défini → `gowsdl`.
* Si tu veux un **client léger** → `gosoap`.

---

Veux-tu que je te prépare un **exemple concret avec parsing XML de la réponse** (méthode manuelle), ou tu préfères un **client généré via WSDL** ?

Parfait 👍 Je vais te montrer un exemple **complet** en Go :

1. Un client qui envoie une requête SOAP manuellement.
2. Le parsing de la réponse XML avec `encoding/xml`.
3. Des tests unitaires avec un serveur HTTP mocké (`httptest`).

---

## 🟦 Exemple concret : SOAP Client en Go

Imaginons un service SOAP avec une méthode `GetUser` qui renvoie un utilisateur :

### 📄 `soap_client.go`

```go
package soapclient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// -------- Structures SOAP --------

// Requête SOAP
func BuildSOAPRequest(userID int) string {
	return fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		<GetUser xmlns="http://example.com/">
		  <id>%d</id>
		</GetUser>
	  </soap:Body>
	</soap:Envelope>`, userID)
}

// Réponse SOAP → Struct
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	GetUserResponse GetUserResponse `xml:"GetUserResponse"`
}

type GetUserResponse struct {
	User User `xml:"User"`
}

type User struct {
	ID   int    `xml:"Id"`
	Name string `xml:"Name"`
}

// -------- Client SOAP --------

type SOAPClient struct {
	Endpoint   string
	HTTPClient *http.Client
}

func NewSOAPClient(endpoint string) *SOAPClient {
	return &SOAPClient{
		Endpoint:   endpoint,
		HTTPClient: &http.Client{},
	}
}

func (c *SOAPClient) GetUser(id int) (*User, error) {
	soapBody := BuildSOAPRequest(id)

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer([]byte(soapBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://example.com/GetUser")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var envelope Envelope
	if err := xml.Unmarshal(bodyBytes, &envelope); err != nil {
		return nil, err
	}

	return &envelope.Body.GetUserResponse.User, nil
}
```

---

## 🟦 Tests unitaires

On utilise `httptest.NewServer` pour simuler un serveur SOAP :

### 📄 `soap_client_test.go`

```go
package soapclient

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser(t *testing.T) {
	// Mock SOAP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vérifie que la requête est bien SOAP
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "text/xml; charset=utf-8" {
			t.Errorf("Unexpected Content-Type: %s", r.Header.Get("Content-Type"))
		}

		// Réponse SOAP simulée
		response := `<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		  <soap:Body>
			<GetUserResponse xmlns="http://example.com/">
			  <User>
				<Id>123</Id>
				<Name>John Doe</Name>
			  </User>
			</GetUserResponse>
		  </soap:Body>
		</soap:Envelope>`
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	// Utilise le client SOAP avec le mock
	client := NewSOAPClient(mockServer.URL)

	user, err := client.GetUser(123)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Vérifications
	if user.ID != 123 {
		t.Errorf("Expected ID=123, got %d", user.ID)
	}
	if user.Name != "John Doe" {
		t.Errorf("Expected Name=John Doe, got %s", user.Name)
	}
}

func TestGetUser_ServerError(t *testing.T) {
	// Mock SOAP server qui renvoie une erreur
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}))
	defer mockServer.Close()

	client := NewSOAPClient(mockServer.URL)

	_, err := client.GetUser(123)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	fmt.Println("Got expected error:", err)
}
```

---

## ✅ Ce que tu as maintenant :

* Un **client SOAP générique** (`SOAPClient`).
* Un **parsing XML** propre avec `encoding/xml`.
* Des **tests unitaires** avec un serveur mock.

---

👉 Veux-tu que je pousse encore plus loin et que je te montre aussi un **test d’intégration** qui lit une vraie réponse XML depuis un fichier (utile si ton WSDL fournit un exemple) ?

Super 👍 Un **test d’intégration avec fichier XML** permet de valider ton parsing sur une vraie réponse SOAP (souvent donnée dans la doc du service ou récupérée via Postman/curl).

Voici la suite de l’exemple avec un fichier de réponse SOAP stocké dans ton projet.

---

## 📄 Exemple de fichier de réponse SOAP

➡️ Sauvegarde ce fichier comme `testdata/get_user_response.xml`

```xml
<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <GetUserResponse xmlns="http://example.com/">
      <User>
        <Id>456</Id>
        <Name>Jane Smith</Name>
      </User>
    </GetUserResponse>
  </soap:Body>
</soap:Envelope>
```

---

## 📄 Test d’intégration

On lit le fichier et on fait comme si c’était la réponse du serveur.

### `soap_client_integration_test.go`

```go
package soapclient

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
)

func TestGetUser_FromFileResponse(t *testing.T) {
	// Charge la réponse SOAP depuis un fichier
	xmlPath := filepath.Join("testdata", "get_user_response.xml")
	data, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		t.Fatalf("Erreur lors de la lecture du fichier XML: %v", err)
	}

	// Crée un serveur mock qui renvoie cette réponse
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer mockServer.Close()

	client := NewSOAPClient(mockServer.URL)

	user, err := client.GetUser(456)
	if err != nil {
		t.Fatalf("Erreur lors de GetUser: %v", err)
	}

	// Vérifie que le parsing fonctionne correctement
	if user.ID != 456 {
		t.Errorf("Expected ID=456, got %d", user.ID)
	}
	if user.Name != "Jane Smith" {
		t.Errorf("Expected Name=Jane Smith, got %s", user.Name)
	}
}
```

---

## ✅ Résultat

* Tu testes **ton parsing XML** avec un vrai exemple.
* Tu peux stocker plusieurs fichiers dans `testdata/` (`get_user_error.xml`, `get_user_notfound.xml`, etc.).
* Les tests deviennent beaucoup plus proches de la **réalité du service SOAP**.

---

👉 Veux-tu que je t’ajoute aussi un exemple où on parse une **réponse SOAP avec une erreur** (ex: `<soap:Fault>`) pour que ton client sache gérer les erreurs du serveur ?

Parfait 👌 gérer les erreurs SOAP est important, car un service SOAP renvoie souvent une enveloppe `<soap:Fault>` au lieu d’un simple code HTTP 500.

---

## 🟦 Exemple d’une réponse SOAP avec erreur (Fault)

➡️ Sauvegarde ce fichier comme `testdata/get_user_fault.xml`

```xml
<?xml version="1.0" encoding="utf-8"?>
<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
  <soap:Body>
    <soap:Fault>
      <faultcode>soap:Client</faultcode>
      <faultstring>User not found</faultstring>
    </soap:Fault>
  </soap:Body>
</soap:Envelope>
```

---

## 🟦 Mise à jour du client pour gérer les Faults

On étend les structs pour capturer `<soap:Fault>` :

### `soap_client.go` (ajouts)

```go
type Fault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
}

type Body struct {
	GetUserResponse GetUserResponse `xml:"GetUserResponse"`
	Fault           *Fault          `xml:"Fault"`
}

func (c *SOAPClient) GetUser(id int) (*User, error) {
	soapBody := BuildSOAPRequest(id)

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer([]byte(soapBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("SOAPAction", "http://example.com/GetUser")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var envelope Envelope
	if err := xml.Unmarshal(bodyBytes, &envelope); err != nil {
		return nil, err
	}

	// Vérifie s'il y a une erreur SOAP
	if envelope.Body.Fault != nil {
		return nil, fmt.Errorf("SOAP Fault: %s (%s)",
			envelope.Body.Fault.FaultString,
			envelope.Body.Fault.FaultCode)
	}

	return &envelope.Body.GetUserResponse.User, nil
}
```

---

## 🟦 Test unitaire pour le Fault

### `soap_client_fault_test.go`

```go
package soapclient

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetUser_Fault(t *testing.T) {
	// Charge une réponse SOAP Fault depuis testdata
	xmlPath := filepath.Join("testdata", "get_user_fault.xml")
	data, err := ioutil.ReadFile(xmlPath)
	if err != nil {
		t.Fatalf("Impossible de lire le fichier XML: %v", err)
	}

	// Serveur mock qui renvoie le Fault
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer mockServer.Close()

	client := NewSOAPClient(mockServer.URL)

	_, err = client.GetUser(999) // utilisateur inexistant
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}

	// Vérifie que le message d'erreur contient bien la faute SOAP
	if !strings.Contains(err.Error(), "User not found") {
		t.Errorf("Expected 'User not found', got %v", err)
	}
	if !strings.Contains(err.Error(), "soap:Client") {
		t.Errorf("Expected 'soap:Client', got %v", err)
	}
}
```

---

## ✅ Ce que tu as maintenant

* **Client SOAP robuste** qui gère :

  * Réponses normales (`GetUserResponse`).
  * Réponses d’erreur (`<soap:Fault>`).
* **Tests unitaires** avec fichiers `testdata/`.
* Un comportement plus réaliste et prêt pour un vrai serveur SOAP.

---

👉 Veux-tu que je t’aide aussi à écrire un **helper générique** pour envoyer n’importe quelle requête SOAP (pas seulement `GetUser`), afin que ton client soit réutilisable pour plusieurs méthodes ?

Excellent 👌 tu veux un client **générique** qui sache :

* Envoyer **n’importe quelle requête SOAP** (pas seulement `GetUser`).
* Accepter un corps XML construit dynamiquement.
* Retourner soit une struct parsée, soit directement le **XML brut** (si tu veux l’unmarshaler toi-même).

Je vais te proposer une version améliorée avec un **helper générique**.

---

## 🟦 Nouveau design : Client SOAP générique

### 📄 `soap_client.go`

```go
package soapclient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Structures génériques SOAP
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	Content xml.RawMessage `xml:",innerxml"` // Contenu générique
	Fault   *Fault         `xml:"Fault"`
}

type Fault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
}

// Client SOAP
type SOAPClient struct {
	Endpoint   string
	HTTPClient *http.Client
}

func NewSOAPClient(endpoint string) *SOAPClient {
	return &SOAPClient{
		Endpoint:   endpoint,
		HTTPClient: &http.Client{},
	}
}

// Méthode générique pour appeler une action SOAP
func (c *SOAPClient) Call(action string, body string, response interface{}) error {
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		%s
	  </soap:Body>
	</soap:Envelope>`, body)

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer([]byte(envelope)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	if action != "" {
		req.Header.Set("SOAPAction", action)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var env Envelope
	if err := xml.Unmarshal(bodyBytes, &env); err != nil {
		return err
	}

	// Gérer les erreurs SOAP
	if env.Body.Fault != nil {
		return fmt.Errorf("SOAP Fault: %s (%s)",
			env.Body.Fault.FaultString, env.Body.Fault.FaultCode)
	}

	// Si l'utilisateur a fourni une struct de réponse
	if response != nil {
		if err := xml.Unmarshal(env.Body.Content, response); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
	}

	return nil
}
```

---

## 🟦 Exemple d’utilisation (avec `GetUser`)

Définis la struct correspondant à ta réponse SOAP :

```go
// Exemple de réponse spécifique
type GetUserResponse struct {
	XMLName xml.Name `xml:"GetUserResponse"`
	User    User     `xml:"User"`
}

type User struct {
	ID   int    `xml:"Id"`
	Name string `xml:"Name"`
}
```

Puis appelle le service :

```go
client := NewSOAPClient("http://example.com/soap")

// Corps spécifique de la requête SOAP
reqBody := `<GetUser xmlns="http://example.com/"><id>123</id></GetUser>`

var resp GetUserResponse
err := client.Call("http://example.com/GetUser", reqBody, &resp)
if err != nil {
	fmt.Println("Erreur:", err)
	return
}

fmt.Printf("Utilisateur: %+v\n", resp.User)
```

---

## 🟦 Test unitaire générique

### 📄 `soap_client_generic_test.go`

```go
package soapclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCall_GenericGetUser(t *testing.T) {
	// Simule un serveur SOAP
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		  <soap:Body>
			<GetUserResponse xmlns="http://example.com/">
			  <User>
				<Id>42</Id>
				<Name>Alice</Name>
			  </User>
			</GetUserResponse>
		  </soap:Body>
		</soap:Envelope>`
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	client := NewSOAPClient(mockServer.URL)

	reqBody := `<GetUser xmlns="http://example.com/"><id>42</id></GetUser>`
	var resp GetUserResponse

	err := client.Call("http://example.com/GetUser", reqBody, &resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.User.ID != 42 || resp.User.Name != "Alice" {
		t.Errorf("unexpected user: %+v", resp.User)
	}
}
```

---

## ✅ Résumé

* `SOAPClient.Call` est **générique** : tu passes l’action SOAP, le corps XML et une struct pour parser la réponse.
* Tu peux réutiliser la même méthode pour **toutes les opérations SOAP** (`GetUser`, `ListOrders`, `CreateInvoice`, etc.).
* Les tests fonctionnent avec un serveur mock.

---

👉 Veux-tu que je t’ajoute aussi un **helper pour construire le body XML automatiquement à partir de structs Go** (genre mini-serializer), au lieu d’écrire le XML à la main ?

Parfait 🚀
On va ajouter un **mini-sérialiseur SOAP** qui construit le corps XML directement à partir de structs Go (au lieu d’écrire manuellement du XML dans des strings).

---

## 🟦 Étape 1 : Structs pour sérialisation

Imaginons une requête `GetUser` qui prend un `Id`.
On définit la struct avec des tags `xml`.

```go
type GetUserRequest struct {
	XMLName xml.Name `xml:"GetUser"`
	Xmlns   string   `xml:"xmlns,attr"`
	ID      int      `xml:"id"`
}
```

Et la réponse déjà vue :

```go
type GetUserResponse struct {
	XMLName xml.Name `xml:"GetUserResponse"`
	User    User     `xml:"User"`
}

type User struct {
	ID   int    `xml:"Id"`
	Name string `xml:"Name"`
}
```

---

## 🟦 Étape 2 : Helper pour sérialiser un body SOAP

On ajoute une fonction qui prend une struct et la transforme en `string` XML :

```go
import (
	"encoding/xml"
)

// MarshalBody convertit une struct Go en XML pour l’inclure dans le SOAP Body
func MarshalBody(v interface{}) (string, error) {
	bytes, err := xml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
```

---

## 🟦 Étape 3 : Utilisation dans le client générique

Maintenant, on peut faire :

```go
client := NewSOAPClient("http://example.com/soap")

// Construire la requête à partir d’une struct
req := GetUserRequest{
	Xmlns: "http://example.com/",
	ID:    123,
}

reqBody, err := MarshalBody(req)
if err != nil {
	panic(err)
}

var resp GetUserResponse
err = client.Call("http://example.com/GetUser", reqBody, &resp)
if err != nil {
	panic(err)
}

fmt.Printf("Utilisateur: %+v\n", resp.User)
```

👉 Le client envoie :

```xml
<GetUser xmlns="http://example.com/">
  <id>123</id>
</GetUser>
```

---

## 🟦 Étape 4 : Test unitaire avec struct

```go
func TestCall_WithStructRequest(t *testing.T) {
	// Mock SOAP server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		  <soap:Body>
			<GetUserResponse xmlns="http://example.com/">
			  <User>
				<Id>77</Id>
				<Name>Bob Marley</Name>
			  </User>
			</GetUserResponse>
		  </soap:Body>
		</soap:Envelope>`
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	client := NewSOAPClient(mockServer.URL)

	// Crée la requête avec struct
	req := GetUserRequest{
		Xmlns: "http://example.com/",
		ID:    77,
	}
	reqBody, _ := MarshalBody(req)

	var resp GetUserResponse
	err := client.Call("http://example.com/GetUser", reqBody, &resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if resp.User.ID != 77 || resp.User.Name != "Bob Marley" {
		t.Errorf("unexpected response: %+v", resp.User)
	}
}
```

---

## ✅ Résultat

* Tu peux maintenant construire des **requêtes SOAP depuis des structs Go** → plus de strings XML bricolées à la main.
* Le client SOAP reste **générique** (`Call`).
* Les tests valident le flux complet (struct → XML → SOAP → struct).

---

👉 Veux-tu que je pousse encore plus loin et que je te prépare un **wrapper complet type SDK** (ex: `client.GetUser(id int)`, `client.ListOrders()`, etc.) qui cache la sérialisation/désérialisation, comme un vrai client SOAP généré ?

Parfait 💡 Là on passe du côté **SDK SOAP maison** : ton client devient vraiment agréable à utiliser, un peu comme si on avait généré du code avec un outil WSDL → Go, mais en gardant le contrôle.

L’idée :

* On garde le **client générique** (`Call`).
* On ajoute des **wrappers dédiés** (ex : `GetUser(id int) (*User, error)`).
* Chaque wrapper s’occupe de construire la requête (`struct → XML`), d’appeler `Call`, et de parser la réponse (`XML → struct`).

---

## 🟦 Client SDK complet

### 📄 `soap_sdk.go`

```go
package soapclient

import (
	"encoding/xml"
	"fmt"
)

// =============================
// Structures spécifiques
// =============================

type GetUserRequest struct {
	XMLName xml.Name `xml:"GetUser"`
	Xmlns   string   `xml:"xmlns,attr"`
	ID      int      `xml:"id"`
}

type GetUserResponse struct {
	XMLName xml.Name `xml:"GetUserResponse"`
	User    User     `xml:"User"`
}

type User struct {
	ID   int    `xml:"Id"`
	Name string `xml:"Name"`
}

// =============================
// SDK : méthodes haut-niveau
// =============================

func (c *SOAPClient) GetUser(id int) (*User, error) {
	// Construit la requête
	req := GetUserRequest{
		Xmlns: "http://example.com/",
		ID:    id,
	}

	reqBody, err := MarshalBody(req)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request: %w", err)
	}

	var resp GetUserResponse
	err = c.Call("http://example.com/GetUser", reqBody, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}
```

---

## 🟦 Exemple d’utilisation

Dans ton `main.go` :

```go
package main

import (
	"fmt"
	"log"

	"myproject/soapclient"
)

func main() {
	client := soapclient.NewSOAPClient("http://example.com/soap")

	user, err := client.GetUser(123)
	if err != nil {
		log.Fatal("Erreur SOAP:", err)
	}

	fmt.Printf("Utilisateur trouvé: %+v\n", user)
}
```

👉 Plus besoin d’écrire du XML ni d’appeler `Call` directement : ton SDK fait tout.

---

## 🟦 Test unitaire du SDK

### 📄 `soap_sdk_test.go`

```go
package soapclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSDK_GetUser(t *testing.T) {
	// Simule un serveur SOAP
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vérifie que la requête contient bien GetUser
		buf := make([]byte, r.ContentLength)
		r.Body.Read(buf)
		if string(buf) == "" || !contains(string(buf), "<GetUser") {
			t.Errorf("Request body does not contain <GetUser>: %s", string(buf))
		}

		// Réponse simulée
		response := `<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		  <soap:Body>
			<GetUserResponse xmlns="http://example.com/">
			  <User>
				<Id>101</Id>
				<Name>Ada Lovelace</Name>
			  </User>
			</GetUserResponse>
		  </soap:Body>
		</soap:Envelope>`
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	client := NewSOAPClient(mockServer.URL)

	user, err := client.GetUser(101)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID != 101 {
		t.Errorf("Expected ID=101, got %d", user.ID)
	}
	if user.Name != "Ada Lovelace" {
		t.Errorf("Expected Name=Ada Lovelace, got %s", user.Name)
	}
}

// petit helper interne
func contains(s, sub string) bool {
	return len(s) >= len(sub) && (len(s) == 0 || (len(s) > 0 && string([]byte(s)) != "" && (string([]byte(s)) != "" && string([]byte(s)) != "" && (len(s) >= len(sub) && string([]byte(s))[0:len(sub)] == sub || contains(string([]byte(s))[1:], sub))))
}
```

*(le `contains` est juste un helper rudimentaire, tu peux remplacer par `strings.Contains` pour faire simple 😉)*

---

## ✅ Ce que tu as maintenant

* Un **client SOAP générique** (`Call`).
* Un **SDK de haut niveau** (`GetUser`) qui masque la sérialisation/désérialisation.
* Des **tests unitaires** réalistes avec un serveur mock.

---

👉 Veux-tu que je généralise le SDK avec un **pattern extensible** (ex. une interface `UserService`, `OrderService`, etc.), pour séparer les différents modules SOAP en paquets propres ?

Excellent 👍 là on va passer d’un simple **client centralisé** à une architecture **SDK modulaire**, où chaque “service” SOAP (ex. `UserService`, `OrderService`, `InvoiceService`) est représenté par un sous-client dédié.

Cela permet :

* de garder un **client générique unique** (`SOAPClient`)
* mais d’exposer des **interfaces métier** claires (`UserService`, `OrderService`…),
* un peu comme le fait un SDK officiel généré depuis un WSDL.

---

# 🟦 Étape 1 : Architecture proposée

```
soapclient/
│── client.go        → Client SOAP générique (Call)
│── marshal.go       → Helper MarshalBody
│── user_service.go  → Service User (GetUser, CreateUser…)
│── order_service.go → Service Order (ListOrders, GetOrder…)
│── models.go        → Structs communes (User, Order…)
```

---

# 🟦 Étape 2 : Client générique

📄 `client.go`

```go
package soapclient

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type SOAPClient struct {
	Endpoint   string
	HTTPClient *http.Client
}

func NewSOAPClient(endpoint string) *SOAPClient {
	return &SOAPClient{
		Endpoint:   endpoint,
		HTTPClient: &http.Client{},
	}
}

// Structures SOAP génériques
type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    Body     `xml:"Body"`
}

type Body struct {
	Content xml.RawMessage `xml:",innerxml"`
	Fault   *Fault         `xml:"Fault"`
}

type Fault struct {
	FaultCode   string `xml:"faultcode"`
	FaultString string `xml:"faultstring"`
}

// Méthode générique
func (c *SOAPClient) Call(action string, body string, response interface{}) error {
	envelope := fmt.Sprintf(`<?xml version="1.0" encoding="utf-8"?>
	<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	  <soap:Body>
		%s
	  </soap:Body>
	</soap:Envelope>`, body)

	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer([]byte(envelope)))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	if action != "" {
		req.Header.Set("SOAPAction", action)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var env Envelope
	if err := xml.Unmarshal(bodyBytes, &env); err != nil {
		return err
	}

	if env.Body.Fault != nil {
		return fmt.Errorf("SOAP Fault: %s (%s)",
			env.Body.Fault.FaultString, env.Body.Fault.FaultCode)
	}

	if response != nil {
		if err := xml.Unmarshal(env.Body.Content, response); err != nil {
			return fmt.Errorf("error decoding response: %w", err)
		}
	}
	return nil
}
```

---

# 🟦 Étape 3 : Helper sérialisation

📄 `marshal.go`

```go
package soapclient

import "encoding/xml"

func MarshalBody(v interface{}) (string, error) {
	bytes, err := xml.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
```

---

# 🟦 Étape 4 : Service `UserService`

📄 `user_service.go`

```go
package soapclient

import (
	"encoding/xml"
	"fmt"
)

// Models
type User struct {
	ID   int    `xml:"Id"`
	Name string `xml:"Name"`
}

// Requête SOAP
type GetUserRequest struct {
	XMLName xml.Name `xml:"GetUser"`
	Xmlns   string   `xml:"xmlns,attr"`
	ID      int      `xml:"id"`
}

// Réponse SOAP
type GetUserResponse struct {
	XMLName xml.Name `xml:"GetUserResponse"`
	User    User     `xml:"User"`
}

// Interface
type UserService interface {
	GetUser(id int) (*User, error)
}

// Implémentation
type userService struct {
	client *SOAPClient
}

func NewUserService(c *SOAPClient) UserService {
	return &userService{client: c}
}

func (u *userService) GetUser(id int) (*User, error) {
	req := GetUserRequest{
		Xmlns: "http://example.com/",
		ID:    id,
	}
	reqBody, err := MarshalBody(req)
	if err != nil {
		return nil, err
	}

	var resp GetUserResponse
	err = u.client.Call("http://example.com/GetUser", reqBody, &resp)
	if err != nil {
		return nil, err
	}

	return &resp.User, nil
}
```

---

# 🟦 Étape 5 : Exemple d’usage SDK

```go
package main

import (
	"fmt"
	"log"

	"myproject/soapclient"
)

func main() {
	client := soapclient.NewSOAPClient("http://example.com/soap")
	userService := soapclient.NewUserService(client)

	user, err := userService.GetUser(42)
	if err != nil {
		log.Fatal("Erreur SOAP:", err)
	}

	fmt.Printf("Utilisateur trouvé: %+v\n", user)
}
```

👉 Lisible, clair : `userService.GetUser(42)`

---

# 🟦 Étape 6 : Test unitaire du service

📄 `user_service_test.go`

```go
package soapclient

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserService_GetUser(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := `<?xml version="1.0" encoding="utf-8"?>
		<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
		  <soap:Body>
			<GetUserResponse xmlns="http://example.com/">
			  <User>
				<Id>42</Id>
				<Name>Alice</Name>
			  </User>
			</GetUserResponse>
		  </soap:Body>
		</soap:Envelope>`
		w.Write([]byte(response))
	}))
	defer mockServer.Close()

	client := NewSOAPClient(mockServer.URL)
	service := NewUserService(client)

	user, err := service.GetUser(42)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID != 42 || user.Name != "Alice" {
		t.Errorf("unexpected user: %+v", user)
	}
}
```

---

# ✅ Résultat

* `SOAPClient` générique pour l’infra.
* `UserService` clair et typé, avec interface → prêt pour **tests/mocks**.
* Tu peux ajouter `OrderService`, `InvoiceService` de la même manière.
* SDK lisible et extensible 🔥.

---
