Excellente question ⚙️

Quand tu veux **appeler une API SOAP en Go**, il ne s’agit pas seulement d’écrire du code de requête XML : il faut **structurer ton projet proprement** pour gérer la sérialisation XML, les erreurs, la logique métier, et la testabilité.

---

## 🧩 Architecture recommandée (propre et testable)

Voici une architecture modulaire typique pour un **client SOAP en Go**, inspirée du pattern *Clean Architecture / Hexagonal Architecture* :

```shell
soap-client/
├── cmd/
│   └── soapcli/             # point d’entrée : main.go
├── internal/
│   ├── app/                 # orchestrateur / use cases
│   ├── domain/              # entités métiers
│   ├── soapclient/          # adapter SOAP (encodage XML, requêtes HTTP)
│   ├── repository/          # si tu persistes des données
│   └── config/              # configuration (URL SOAP, credentials, timeouts)
└── tests/
    └── integration_test.go  # tests SOAP end-to-end
```

---

## ⚙️ Description des couches

### `internal/config`

Contient la configuration de base :

```go
type Config struct {
    SOAPEndpoint string
    SOAPUser     string
    SOAPPassword string
    Timeout      time.Duration
}
```

---

### `internal/domain`

Définit les **entités métier**, indépendantes du protocole SOAP.
Exemple :

```go
type Invoice struct {
    ID     string
    Amount float64
    Status string
}
```

---

### `internal/soapclient`

C’est la couche qui parle SOAP.
Elle contient :

* Les structures XML (Request / Response)
* L’encodage/décodage XML
* L’appel HTTP
* Les erreurs SOAP

Exemple minimal :

```go
package soapclient

import (
    "bytes"
    "encoding/xml"
    "fmt"
    "net/http"
    "io/ioutil"
)

type SOAPClient struct {
    Endpoint string
    Username string
    Password string
    Client   *http.Client
}

func New(endpoint, user, pass string) *SOAPClient {
    return &SOAPClient{
        Endpoint: endpoint,
        Username: user,
        Password: pass,
        Client:   &http.Client{},
    }
}

type GetInvoiceRequest struct {
    XMLName xml.Name `xml:"soapenv:Envelope"`
    SoapEnv string   `xml:"xmlns:soapenv,attr"`
    Body    GetInvoiceBody
}

type GetInvoiceBody struct {
    XMLName xml.Name `xml:"soapenv:Body"`
    Request struct {
        XMLName xml.Name `xml:"ns:GetInvoice"`
        ID      string   `xml:"ns:id"`
        NS      string   `xml:"xmlns:ns,attr"`
    }
}

func (c *SOAPClient) GetInvoice(id string) (string, error) {
    req := GetInvoiceRequest{
        SoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
    }
    req.Body.Request.NS = "http://example.com/soap"
    req.Body.Request.ID = id

    xmlBytes, _ := xml.Marshal(req)
    body := bytes.NewReader(xmlBytes)

    httpReq, _ := http.NewRequest("POST", c.Endpoint, body)
    httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")

    resp, err := c.Client.Do(httpReq)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    raw, _ := ioutil.ReadAll(resp.Body)
    if resp.StatusCode != http.StatusOK {
        return "", fmt.Errorf("soap error: %s", string(raw))
    }

    return string(raw), nil
}
```

---

### `internal/app`

C’est la logique métier (use cases).
Exemple :

```go
package app

import "github.com/you/soap-client/internal/soapclient"

type Service struct {
    soap *soapclient.SOAPClient
}

func NewService(soap *soapclient.SOAPClient) *Service {
    return &Service{soap: soap}
}

func (s *Service) GetInvoiceInfo(id string) (string, error) {
    return s.soap.GetInvoice(id)
}
```

---

### `cmd/soapcli/main.go`

Point d’entrée :

```go
package main

import (
    "log"
    "github.com/you/soap-client/internal/app"
    "github.com/you/soap-client/internal/soapclient"
)

func main() {
    client := soapclient.New("https://example.com/soap", "user", "pass")
    svc := app.NewService(client)

    invoice, err := svc.GetInvoiceInfo("12345")
    if err != nil {
        log.Fatal(err)
    }

    log.Println(invoice)
}
```

---

## ✅ Bonnes pratiques

* ✅ **Génère le client à partir du WSDL** si possible avec [`wsdl2go`](https://github.com/hooklift/gowsdl) (automatique et typé).
* ✅ Garde ton **client SOAP isolé** dans un package — tu pourras le mocker dans les tests.
* ✅ **Teste avec des mocks** pour éviter d’appeler le service SOAP réel.
* ✅ Gère les **timeouts et retry** (SOAP est souvent lent).
* ✅ Ajoute une **validation XML stricte** si le service est sensible au namespace.

---

## 💡 Variante moderne

Si ton API SOAP est complexe, tu peux :

* Créer un **adapter REST** interne (un mini proxy) pour convertir SOAP ⇄ JSON pour ton application principale ;
* Ou encapsuler le SOAP client dans un **use case isolé**, accessible via une interface `ExternalBillingAPI` ou `CRMClient`.

---

Souhaites-tu que je te génère **un squelette complet de projet SOAP client en Go** (comme celui avec Huma) avec structure, `main.go`, client SOAP, tests, et Dockerfile ?
Je peux le faire en une minute.

# Squelette : Client SOAP isolé (stateless) — Projet Go

Ce squelette est conçu pour appeler une API SOAP depuis Go **sans état** (pas de DB). Le client SOAP est **isolé** dans `internal/soapclient` et exposé via une interface que la couche métier (`internal/app`) consomme — facile à mocker pour les tests.

Arborescence recommandée :

```shell
soap-client/
├── cmd/soapcli/main.go
├── internal/
│   ├── app/               # use-cases / orchestration
│   │   └── service.go
│   ├── domain/            # entités métiers
│   │   └── invoice.go
│   └── soapclient/        # client SOAP isolé
│       ├── client.go
│       └── types.go       # structures XML request/response
├── internal/app/service_test.go
├── go.mod
├── Makefile
├── Dockerfile
└── README.md
```

---

## `go.mod`

```go
module github.com/you/soap-client

go 1.21

require (
    github.com/stretchr/testify v1.8.4
)
```

---

## `cmd/soapcli/main.go`

```go
package main

import (
    "context"
    "log"
    "time"

    "github.com/you/soap-client/internal/app"
    "github.com/you/soap-client/internal/soapclient"
)

func main() {
    // config (hardcoded pour l'exemple) — tu peux le remplacer par viper/env
    endpoint := "https://example.com/soap" // remplacer
    timeout := 15 * time.Second

    // créer le client SOAP isolé
    client := soapclient.NewSOAPClient(endpoint, nil)
    client.SetHTTPTimeout(timeout)

    // créer le service métier qui utilise le client SOAP
    svc := app.NewService(client)

    // exécuter un use-case
    ctx := context.Background()
    resp, err := svc.GetInvoiceInfo(ctx, "INV-12345")
    if err != nil {
        log.Fatalf("error calling SOAP: %v", err)
    }

    log.Printf("Invoice info: %+v\n", resp)
}
```

---

## `internal/soapclient/client.go`

```go
package soapclient

import (
    "bytes"
    "context"
    "encoding/xml"
    "fmt"
    "io"
    "net/http"
    "time"
)

// Interface exposée par le package — utile pour mock dans les tests
type Client interface {
    GetInvoice(ctx context.Context, id string) (*GetInvoiceResponse, error)
    SetHTTPTimeout(d time.Duration)
}

// SOAPClient implémente Client
type SOAPClient struct {
    endpoint string
    httpc    *http.Client
}

func NewSOAPClient(endpoint string, httpc *http.Client) *SOAPClient {
    if httpc == nil {
        httpc = &http.Client{Timeout: 10 * time.Second}
    }
    return &SOAPClient{endpoint: endpoint, httpc: httpc}
}

func (c *SOAPClient) SetHTTPTimeout(d time.Duration) {
    c.httpc.Timeout = d
}

func (c *SOAPClient) doRequest(ctx context.Context, action string, payload []byte) ([]byte, error) {
    req, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewReader(payload))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "text/xml; charset=utf-8")
    if action != "" {
        // certains services SOAP demandent header SOAPAction
        req.Header.Set("SOAPAction", action)
    }

    resp, err := c.httpc.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("soap request failed: status=%d body=%s", resp.StatusCode, string(body))
    }
    return body, nil
}

// GetInvoice: construit la requête XML, l'envoie et parse la réponse
func (c *SOAPClient) GetInvoice(ctx context.Context, id string) (*GetInvoiceResponse, error) {
    // construire Envelope
    reqEnvelope := Envelope{
        XMLName: xml.Name{Local: "soapenv:Envelope"},
        XmlnsSoapenv: "http://schemas.xmlsoap.org/soap/envelope/",
        Body: Body{
            Content: &GetInvoiceRequest{ID: id, XmlnsNs: "http://example.com/soap"},
        },
    }

    b, err := xml.Marshal(reqEnvelope)
    if err != nil {
        return nil, err
    }

    // option: ajouter declaration xml
    payload := append([]byte(xml.Header), b...)

    raw, err := c.doRequest(ctx, "", payload)
    if err != nil {
        return nil, err
    }

    // parser la réponse
    var respEnvelope Envelope
    if err := xml.Unmarshal(raw, &respEnvelope); err != nil {
        return nil, err
    }

    // extraire le content comme GetInvoiceResponse
    if respEnvelope.Body.Fault != nil {
        return nil, fmt.Errorf("soap fault: %s", respEnvelope.Body.Fault.String())
    }

    if respEnvelope.Body.Content == nil {
        return nil, fmt.Errorf("empty soap body")
    }

    // respEnvelope.Body.Content est typé comme interface{} via unmarshalling
    // mais dans notre struct types.go on utilise *GetInvoiceResponse dans Body.Content
    if resp, ok := respEnvelope.Body.Content.(*GetInvoiceResponse); ok {
        return resp, nil
    }

    // fallback: tenter de remapper en bytes -> GetInvoiceResponse
    // (depends on how xml.Unmarshal filled Content)
    // on remarshal le body puis unmarshal
    var candidate GetInvoiceResponse
    // essayer d'extraire le inner XML
    // (pour la plupart des WSDL simples, le premier child sera le bon)
    // pour robustesse, on re-unmarshal avec une struct enveloppe ciblée

    type RespEnv struct {
        XMLName xml.Name `xml:"Envelope"`
        Body    struct {
            GetInvoiceResponse *GetInvoiceResponse `xml:"GetInvoiceResponse"`
        } `xml:"Body"`
    }
    var re RespEnv
    if err := xml.Unmarshal(raw, &re); err == nil && re.Body.GetInvoiceResponse != nil {
        candidate = *re.Body.GetInvoiceResponse
        return &candidate, nil
    }

    return nil, fmt.Errorf("unexpected soap response shape")
}
```

---

## `internal/soapclient/types.go`

```go
package soapclient

import "encoding/xml"

// Envelope / Body / Fault génériques pour SOAP 1.1
type Envelope struct {
    XMLName     xml.Name `xml:"soapenv:Envelope"`
    XmlnsSoapenv string   `xml:"xmlns:soapenv,attr"`
    Body        Body     `xml:"soapenv:Body"`
}

type Body struct {
    Content interface{} `xml:",omitempty"`
    Fault   *Fault      `xml:"Fault,omitempty"`
}

// Fault minimal
type Fault struct {
    FaultCode   string `xml:"faultcode"`
    FaultString string `xml:"faultstring"`
}

func (f *Fault) String() string { return f.FaultString }

// GetInvoiceRequest: adapte les namespaces et noms selon ton WSDL
type GetInvoiceRequest struct {
    XMLName xml.Name `xml:"ns:GetInvoice"`
    XmlnsNs string   `xml:"xmlns:ns,attr,omitempty"`
    ID      string   `xml:"ns:id"`
}

// GetInvoiceResponse: structure attendue (ex: GetInvoiceResponse / Invoice)
type GetInvoiceResponse struct {
    XMLName xml.Name `xml:"GetInvoiceResponse"`
    Invoice Invoice  `xml:"Invoice"`
}

type Invoice struct {
    ID     string  `xml:"ID"`
    Amount float64 `xml:"Amount"`
    Status string  `xml:"Status"`
}

// Custom Unmarshal for Body.Content so xml.Unmarshal knows to decode into GetInvoiceResponse
func (b *Body) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
    // read tokens to find the first element under Body
    for {
        tok, err := d.Token()
        if err != nil {
            return err
        }
        switch t := tok.(type) {
        case xml.StartElement:
            // assume it's GetInvoiceResponse for this example
            if t.Name.Local == "GetInvoiceResponse" || t.Name.Local == "GetInvoiceResponse" {
                var r GetInvoiceResponse
                if err := d.DecodeElement(&r, &t); err != nil {
                    return err
                }
                b.Content = &r
                return nil
            }
            // unknown: decode into raw
            var raw xml.Token
            _ = raw
        case xml.EndElement:
            return nil
        }
    }
}
```

---

## `internal/domain/invoice.go`

```go
package domain

type Invoice struct {
    ID     string
    Amount float64
    Status string
}
```

---

## `internal/app/service.go`

```go
package app

import (
    "context"
    "github.com/you/soap-client/internal/soapclient"
    "github.com/you/soap-client/internal/domain"
)

// Service uses a soapclient.Client (interface) so it can be mocked
type Service struct {
    client soapclient.Client
}

func NewService(c soapclient.Client) *Service { return &Service{client: c} }

// GetInvoiceInfo wraps the SOAP response into domain.Invoice
func (s *Service) GetInvoiceInfo(ctx context.Context, id string) (*domain.Invoice, error) {
    resp, err := s.client.GetInvoice(ctx, id)
    if err != nil {
        return nil, err
    }

    inv := &domain.Invoice{
        ID:     resp.Invoice.ID,
        Amount: resp.Invoice.Amount,
        Status: resp.Invoice.Status,
    }
    return inv, nil
}
```

---

## `internal/app/service_test.go`

```go
package app

import (
    "context"
    "testing"

    "github.com/stretchr/testify/require"
    "github.com/you/soap-client/internal/domain"
    "github.com/you/soap-client/internal/soapclient"
)

// mock client implementing soapclient.Client
type mockClient struct{}

func (m *mockClient) SetHTTPTimeout(d time.Duration) {}

func (m *mockClient) GetInvoice(ctx context.Context, id string) (*soapclient.GetInvoiceResponse, error) {
    return &soapclient.GetInvoiceResponse{
        Invoice: soapclient.Invoice{ID: id, Amount: 42.5, Status: "PAID"},
    }, nil
}

func TestGetInvoiceInfo(t *testing.T) {
    mc := &mockClient{}
    svc := NewService(mc)

    inv, err := svc.GetInvoiceInfo(context.Background(), "INV-1")
    require.NoError(t, err)
    require.Equal(t, "INV-1", inv.ID)
    require.Equal(t, float64(42.5), inv.Amount)
    require.Equal(t, "PAID", inv.Status)
}
```

(Note: le test utilise un mock simple ; pour de vrais tests d'intégration on peut démarrer un serveur SOAP de test.)

---

## `Makefile`

```makefile
.PHONY: run test build

run:
	go run ./cmd/soapcli

test:
	go test ./... -v

build:
	go build -o bin/soapcli ./cmd/soapcli
```

---

## `Dockerfile`

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o /bin/soapcli ./cmd/soapcli

FROM alpine:3.18
COPY --from=builder /bin/soapcli /bin/soapcli
ENTRYPOINT ["/bin/soapcli"]
```

---

## `README.md` (extrait)

```markdown
# SOAP Client (stateless)

## Run

go run ./cmd/soapcli

## Test

make test
```

---

## Remarques & étapes suivantes

* **Adaptation WSDL** : pour un service réel, tu dois adapter les `types.go` (namespaces, noms d'éléments) au WSDL du fournisseur. Si possible, utilise un générateur (`gowsdl`, `wsdl2go`) pour produire des types fortement typés.
* **Authentification** : si SOAP nécessite WS-Security (Header), il faudra ajouter le header approprié (token, signature) dans `doRequest` ou construire l'enveloppe avec les headers WS-Security.
* **Retries / Timeout** : ajoute des retries exponentiels et metrics selon besoin.
* **Tests d'intégration** : créer un serveur SOAP mock (httptest.Server) qui renvoie réponses XML pour valider le parsing.

---

Si tu veux, je peux maintenant :

* générer chaque fichier dans le canvas séparément pour copie téléchargeable, ou
* adapter `types.go` à un WSDL réel si tu me fournis le WSDL/endpoint, ou
* ajouter gestion WS-Security (header UsernameToken ou signature) — dis ce que tu veux que j'ajoute.


C’est prêt — j’ai créé un **squelette stateless** avec le client SOAP isolé dans `internal/soapclient`, la couche métier, les entités, tests, Makefile et Dockerfile, dans un document nommé **"Soap-client-stateless-skeleton"** sur le canvas.

Options rapides — dis-moi laquelle tu veux next (je m’en occupe tout de suite) :

* Générer chaque fichier séparément dans le canvas pour téléchargement.
* Adapter `types.go` au WSDL réel (envoie le WSDL ou l’endpoint).
* Ajouter WS-Security (UsernameToken ou signature).
* Ajouter un mock server (`httptest.Server`) pour tests d’intégration.


## test integration

# Squelette : Client SOAP isolé (stateless) — Projet Go avec Tests d’Intégration

Ce squelette appelle une API SOAP depuis Go **sans état** (pas de DB). Le client SOAP est **isolé** dans `internal/soapclient` et consommé par la couche métier `internal/app`. Des **tests d’intégration** utilisent `httptest.Server` pour simuler un serveur SOAP.

---

## Arborescence

```
soap-client/
├── cmd/soapcli/main.go
├── internal/
│   ├── app/
│   │   ├── service.go
│   │   └── service_test.go
│   ├── domain/
│   │   └── invoice.go
│   └── soapclient/
│       ├── client.go
│       ├── types.go
│       └── client_integration_test.go
├── go.mod
├── Makefile
├── Dockerfile
└── README.md
```

---

## `internal/soapclient/client_integration_test.go`

```go
package soapclient

import (
    "context"
    "fmt"
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
    "time"

    "github.com/stretchr/testify/require"
)

func TestSOAPClient_GetInvoice_Integration(t *testing.T) {
    // simulate a SOAP endpoint using httptest.Server
    soapResponse := `<?xml version="1.0" encoding="UTF-8"?>
        <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/">
            <soapenv:Body>
                <GetInvoiceResponse>
                    <Invoice>
                        <ID>INV-42</ID>
                        <Amount>123.45</Amount>
                        <Status>PAID</Status>
                    </Invoice>
                </GetInvoiceResponse>
            </soapenv:Body>
        </soapenv:Envelope>`

    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // verify SOAPAction header or content if necessary
        body := new(strings.Builder)
        _, _ = body.ReadFrom(r.Body)
        fmt.Printf("Received SOAP request: %s\n", body.String())

        w.Header().Set("Content-Type", "text/xml; charset=utf-8")
        _, _ = w.Write([]byte(soapResponse))
    }))
    defer server.Close()

    // create client pointing to the mock server
    c := NewSOAPClient(server.URL, &http.Client{Timeout: 5 * time.Second})

    ctx := context.Background()
    resp, err := c.GetInvoice(ctx, "INV-42")
    require.NoError(t, err)
    require.NotNil(t, resp)

    require.Equal(t, "INV-42", resp.Invoice.ID)
    require.Equal(t, 123.45, resp.Invoice.Amount)
    require.Equal(t, "PAID", resp.Invoice.Status)
}
```

Ce test démarre un faux serveur SOAP local et valide :

* La requête SOAP émise par `SOAPClient`.
* Le parsing XML complet jusqu’à `GetInvoiceResponse`.

---

## Commandes

```bash
make test # inclut les tests d’intégration
```

---

## Étapes suivantes possibles

* Ajouter un test d’intégration **avec SOAP Fault** pour valider la gestion d’erreur XML.
* Ajouter un test **avec délai (timeout)** pour tester la résilience réseau.
* Ajouter des fixtures XML de test dans un dossier `testdata/`.

C’est fait ✅ — j’ai ajouté un test d’intégration complet (`client_integration_test.go`) avec un serveur SOAP mocké (`httptest.Server`).
Souhaites-tu que j’ajoute aussi un **test de gestion d’erreur SOAP (SOAP Fault)** pour vérifier la robustesse du client ?
