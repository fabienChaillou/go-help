## SOLID PRINCIPES

Les principes SOLID sont un ensemble de cinq principes de conception orientée objet qui facilitent la **maintenance**, **l’évolutivité** et la **lisibilité** du code. Bien que Go (Golang) ne soit pas un langage orienté objet traditionnel comme Java ou C#, il permet tout de même d’appliquer ces principes grâce à ses **interfaces**, **structures**, et **fonctions**.


Voici une explication des principes SOLID adaptés à Go :

1. S - Single Responsibility Principle (SRP)


Une fonction ou une structure doit avoir une seule raison de changer.


En Go :

Cela signifie qu’un struct ou un package ne doit faire qu’une seule chose.

Chaque fonction ou méthode doit être courte, claire, et spécialisée.


Exemple :

```go
type Invoice struct {
    Amount float64
}

func (i Invoice) CalculateTax() float64 {
    return i.Amount * 0.2
}

// Séparer la logique d’affichage :

type InvoicePrinter struct {}

func (p InvoicePrinter) Print(i Invoice) {
    fmt.Println("Total:", i.Amount)
}
```

2. O - Open/Closed Principle (OCP)


Le code doit être ouvert à l’extension mais fermé à la modification.


En Go :

On applique ce principe en utilisant des interfaces pour permettre l’extension sans modifier le code existant.

Exemple :

```go
type Shape interface {
    Area() float64
}

type Circle struct {
    Radius float64
}

func (c Circle) Area() float64 {
    return math.Pi * c.Radius * c.Radius
}

func TotalArea(shapes []Shape) float64 {
    var total float64
    for _, s := range shapes {
        total += s.Area()
    }
    return total
}
```

Ajouter une nouvelle forme (comme Rectangle) ne nécessite pas de modifier TotalArea.
3. L - Liskov Substitution Principle (LSP)


Les objets d’une classe dérivée doivent pouvoir remplacer ceux de la classe de base sans altérer le comportement.


En Go :

Cela revient à dire qu’un type concret doit satisfaire correctement l’interface qu’il implémente.



Exemple :

```go
type Bird interface {
    Fly()
}

type Sparrow struct {}

func (s Sparrow) Fly() {
    fmt.Println("Sparrow flying")
}
```

Si un type comme Penguin ne peut pas voler, il ne doit pas implémenter l’interface Bird car cela violerait le LSP.

4. I - Interface Segregation Principle (ISP)


Les clients ne doivent pas être forcés de dépendre d’interfaces qu’ils n’utilisent pas.


En Go :

Il est recommandé de créer de petites interfaces spécifiques.

Exemple :

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}
```

Plutôt que d’utiliser io.ReadWriter si le client n’a besoin que de Read, on ne l’oblige pas à implémenter Write.

5. D - Dependency Inversion Principle (DIP)

Le code doit dépendre d’abstractions (interfaces), pas de concrétions (types concrets).

En Go : Utilise des interfaces pour injecter les dépendances.


Exemple :

```go
type Notifier interface {
    Send(message string) error
}

type EmailNotifier struct {}

func (e EmailNotifier) Send(msg string) error {
    fmt.Println("Sending email:", msg)
    return nil
}

type Alert struct {
    notifier Notifier
}

func (a Alert) Trigger() {
    a.notifier.Send("Alert triggered!")
}
```

Alert ne dépend que de l’interface Notifier, donc on peut injecter n’importe quel type qui l’implémente.


## Example

Voici un exemple de petit projet Go respectant les principes SOLID : un système de notifications avec alertes critiques, extensible et modulaire.

Contexte du projet


Un service d’alertes doit envoyer des notifications via différents canaux : e-mail, SMS, Slack, etc. Chaque type de notification doit être facilement ajoutable sans modifier le cœur du système.

1. Structure du projet
/alertsystem
├── main.go
├── notifier
│   ├── notifier.go       // interfaces et implémentations (Email, Slack, etc.)
├── alert
│   ├── alert.go          // logique d'alerte
2. Code résumé


notifier/notifier.go
package notifier

```go
import "fmt"

// Interface (DIP + ISP)
type Notifier interface {
    Send(message string) error
}

// Email notifier (SRP)
type EmailNotifier struct{}

func (e EmailNotifier) Send(message string) error {
    fmt.Println("[Email] Sending:", message)
    return nil
}

// Slack notifier (extensible — OCP)
type SlackNotifier struct{}

func (s SlackNotifier) Send(message string) error {
    fmt.Println("[Slack] Sending:", message)
    return nil
}

// alert/alert.go
package alert

import "alertsystem/notifier"

// Alerte (SRP + DIP)
type Alert struct {
    Notifiers []notifier.Notifier
}

func (a Alert) Trigger(msg string) {
    for _, n := range a.Notifiers {
        _ = n.Send(msg) // Erreurs ignorées ici pour simplifier
    }
}

// main.go
package main

import (
    "alertsystem/alert"
    "alertsystem/notifier"
)

func main() {
    email := notifier.EmailNotifier{}
    slack := notifier.SlackNotifier{}

    a := alert.Alert{
        Notifiers: []notifier.Notifier{email, slack},
    }

    a.Trigger("CRITICAL: Disk space low!")
}
```

Comment ça respecte SOLID
Principe

Application

SRP

Chaque composant a une responsabilité claire.

OCP

On peut ajouter de nouveaux notifiers sans modifier Alert.

LSP

Tous les notifiers peuvent être utilisés de manière interchangeable.

ISP

L’interface Notifier est fine et spécifique.

DIP

Alert dépend de l’interface, pas des implémentations concrètes.

Souhaites-tu que je t’aide à étendre ce projet avec une base de données, une API REST ou un autre canal de notification ?

### example avc canal de notification

Ajoutons un nouveau canal de notification : une notification SMS. Cela montre comment appliquer le principe Open/Closed : on étend le système sans modifier les composants existants.

Étapes


1. Créer une nouvelle struct SMSNotifier


Dans notifier/notifier.go :

```go
type SMSNotifier struct{}

func (s SMSNotifier) Send(message string) error {
    fmt.Println("[SMS] Sending:", message)
    return nil
}
```

Cette nouvelle struct implémente l’interface Notifier, donc elle est interopérable avec le reste du système.

2. Utiliser SMS dans main.go

```go
package main

import (
    "alertsystem/alert"
    "alertsystem/notifier"
)

func main() {
    email := notifier.EmailNotifier{}
    slack := notifier.SlackNotifier{}
    sms := notifier.SMSNotifier{} // Nouveau canal

    a := alert.Alert{
        Notifiers: []notifier.Notifier{email, slack, sms},
    }

    a.Trigger("CRITICAL: CPU usage over 90%!")
}
```

3. Résultat à l’exécution
[Email] Sending: CRITICAL: CPU usage over 90%!
[Slack] Sending: CRITICAL: CPU usage over 90%!
[SMS] Sending: CRITICAL: CPU usage over 90%!
Conclusion


Tu n’as rien modifié dans Alert, ni dans les autres notifiers : tu as simplement étendu le comportement, ce qui est un exemple clair d’Open/Closed Principle.



Souhaites-tu que l’on ajoute une vraie intégration SMS via une API (comme Twilio) ou qu’on continue avec d’autres canaux (ex. : webhook, push mobile) ?



Suite example



Bien sûr ! Pour intégrer l’envoi de SMS via Twilio dans ton projet Go tout en respectant les principes SOLID, voici comment procéder.

🔧 Étape 1 : Installer la bibliothèque Twilio pour Go


Commence par installer la bibliothèque officielle de Twilio pour Go :

```shell
go get github.com/twilio/twilio-go
```

📦 Étape 2 : Créer un canal SMS avec Twilio


Dans le fichier notifier/sms_notifier.go, définis une nouvelle structure SMSNotifier qui implémente l’interface Notifier :

```go
package notifier

import (
    "fmt"
    "github.com/twilio/twilio-go"
    "github.com/twilio/twilio-go/rest/api/v2010"
    "os"
)

type SMSNotifier struct {
    From string
}

func NewSMSNotifier(from string) SMSNotifier {
    return SMSNotifier{From: from}
}

func (s SMSNotifier) Send(message string) error {
    client := twilio.NewRestClient()

    params := &api.CreateMessageParams{}
    params.SetBody(message)
    params.SetFrom(s.From)
    params.SetTo("+1234567890") // Remplace par le numéro du destinataire

    resp, err := client.Api.CreateMessage(params)
    if err != nil {
        return fmt.Errorf("failed to send SMS: %w", err)
    }

    fmt.Println("SMS envoyé :", *resp.Body)
    return nil
}
```

🧪 Étape 3 : Utiliser le SMSNotifier dans `main.go`


Dans le fichier main.go, crée une instance de SMSNotifier et utilise-la pour envoyer un message :

```go
package main

import (
    "alertsystem/alert"
    "alertsystem/notifier"
    "os"
)

func main() {
    sms := notifier.NewSMSNotifier(os.Getenv("TWILIO_PHONE_NUMBER"))

    a := alert.Alert{
        Notifiers: []notifier.Notifier{sms},
    }

    a.Trigger("ALERTE : Espace disque faible !")
}
```

✅ Étape 4 : Définir les variables d’environnement


Avant d’exécuter le programme, assure-toi de définir les variables d’environnement nécessaires :

```shell
export TWILIO_ACCOUNT_SID="votre_account_sid"
export TWILIO_AUTH_TOKEN="votre_auth_token"
export TWILIO_PHONE_NUMBER="votre_numero_twilio"
```

Remplace les valeurs par celles fournies par Twilio.

📌 Remarques importantes
Si tu utilises un compte Twilio en version d’essai, tu devras vérifier le numéro de téléphone du destinataire dans le console Twilio.

Pour plus de détails et d’exemples, consulte la documentation officielle de Twilio pour Go :

Programmable Messaging Quickstart - Go

Souhaites-tu que je t’aide à configurer un serveur HTTP pour recevoir et répondre à des messages entrants via Twilio ?

--------------------------------------------------------------------------------------------------------------------------------------

Parfait ! Voici comment configurer un serveur HTTP en Go pour recevoir des messages entrants via Twilio Webhooks, tout en respectant les principes SOLID.

Contexte


Twilio peut envoyer des webhooks HTTP POST à ton serveur chaque fois qu’un SMS est reçu. Ton serveur doit :

Écouter sur un port HTTP.

Extraire les données du message.

Répondre éventuellement (Twilio peut afficher ta réponse comme un SMS).

Étape 1 : Créer le serveur HTTP


Dans main.go ou dans un fichier séparé (server/server.go) :

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func handleIncomingSMS(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    from := r.FormValue("From")
    body := r.FormValue("Body")

    fmt.Printf("Message reçu de %s : %s\n", from, body)

    // Répondre à Twilio avec un message TwiML
    w.Header().Set("Content-Type", "application/xml")
    fmt.Fprintf(w, `<Response><Message>Merci pour votre message !</Message></Response>`)
}

func main() {
    http.HandleFunc("/sms", handleIncomingSMS)

    fmt.Println("Serveur en écoute sur http://localhost:8080/sms")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Étape 2 : Démarrer le serveur localement

`go run main.go`

Étape 3 : Exposer le serveur local à Internet (avec ngrok)


Twilio doit pouvoir accéder à ton serveur.

ngrok http 8080
Tu recevras une URL du style https://abcd1234.ngrok.io/sms.
Étape 4 : Configurer le webhook Twilio


Dans Twilio Console > Messaging > Sandbox ou numéro actif, mets :

https://abcd1234.ngrok.io/sms
comme A MESSAGE COMES IN URL (en POST).

Ce que tu obtiens
Un système complet capable d’envoyer des alertes par SMS.

Un serveur qui reçoit et traite les réponses ou SMS entrants.

Un design propre et modulaire compatible avec SOLID.

Souhaites-tu que je t’aide à intégrer une logique plus poussée, comme stocker les messages dans une base de données ou déclencher une action en fonction du contenu reçu ?



Parfait ! On va ajouter une logique métier pour déclencher des actions différentes selon le contenu du SMS reçu — tout en gardant une architecture propre et modulable.

Objectif


Quand un utilisateur envoie un SMS :

S’il contient “status”, le serveur répond avec un état du système.

S’il contient “reboot”, il simule un redémarrage de service.

Sinon, il répond : “Commande inconnue”.

Étape 1 : Créer une interface CommandHandler


Dans un nouveau fichier handler/command.go :

```go
package handler

import (
    "fmt"
    "strings"
)

// Interface unique pour appliquer DIP et OCP
type CommandHandler interface {
    CanHandle(cmd string) bool
    Handle(cmd string) string
}
```

Étape 2 : Implémenter quelques commandes


Toujours dans handler/command.go :

```go
type StatusCommand struct{}

func (s StatusCommand) CanHandle(cmd string) bool {
    return strings.ToLower(cmd) == "status"
}

func (s StatusCommand) Handle(cmd string) string {
    return "Système opérationnel."
}

type RebootCommand struct{}

func (r RebootCommand) CanHandle(cmd string) bool {
    return strings.ToLower(cmd) == "reboot"
}

func (r RebootCommand) Handle(cmd string) string {
    // On simule une action
    return "Redémarrage du système en cours..."
}
```

Étape 3 : Adapter le serveur HTTP


Dans main.go :

```go
package main

import (
    "fmt"
    "log"
    "net/http"
    "strings"

    "alertsystem/handler"
)

var commandHandlers = []handler.CommandHandler{
    handler.StatusCommand{},
    handler.RebootCommand{},
}

func handleIncomingSMS(w http.ResponseWriter, r *http.Request) {
    if err := r.ParseForm(); err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        return
    }

    body := strings.TrimSpace(r.FormValue("Body"))
    fmt.Println("Commande reçue :", body)

    response := "Commande inconnue."
    for _, h := range commandHandlers {
        if h.CanHandle(body) {
            response = h.Handle(body)
            break
        }
    }

    // Réponse XML à Twilio
    w.Header().Set("Content-Type", "application/xml")
    fmt.Fprintf(w, `<Response><Message>%s</Message></Response>`, response)
}

func main() {
    http.HandleFunc("/sms", handleIncomingSMS)
    fmt.Println("Serveur en écoute sur http://localhost:8080/sms")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

Résultat
Envoie “status” : tu obtiens “Système opérationnel.”

Envoie “reboot” : “Redémarrage du système en cours…”

Tout autre message : “Commande inconnue.”

Avantage SOLID
Tu peux ajouter une nouvelle commande (ex: ShutdownCommand) sans modifier le reste du code, juste en ajoutant une nouvelle struct qui implémente CommandHandler.

Souhaites-tu que je t’aide à structurer tout cela dans un vrai projet Go (avec modules, dossiers, tests, etc.) ?


Super ! Voici comment structurer un vrai projet Go basé sur notre système d’alertes + commandes SMS, en respectant les bonnes pratiques de projet Go et les principes SOLID.

Structure du projet

```shell
go-sms-alert/
├── go.mod
├── main.go
├── internal/
│   ├── server/           # Serveur HTTP et routing
│   │   └── http.go
│   ├── handler/          # Logique de commandes (status, reboot, etc.)
│   │   └── command.go
│   └── notifier/         # Interface Notifier et implémentations (Email, SMS, etc.)
│       └── notifier.go
└── test/
    └── handler_test.go   # Tests unitaires des commandes
```

Étape 1 : Initialiser le projet
mkdir go-sms-alert && cd go-sms-alert
go mod init go-sms-alert
Étape 2 : Exemple de contenu

```go
// internal/handler/command.go
package handler

import "strings"

type CommandHandler interface {
    CanHandle(cmd string) bool
    Handle(cmd string) string
}

type StatusCommand struct{}
func (s StatusCommand) CanHandle(cmd string) bool {
    return strings.ToLower(cmd) == "status"
}
func (s StatusCommand) Handle(cmd string) string {
    return "Système opérationnel."
}

type RebootCommand struct{}
func (r RebootCommand) CanHandle(cmd string) bool {
    return strings.ToLower(cmd) == "reboot"
}
func (r RebootCommand) Handle(cmd string) string {
    return "Redémarrage en cours..."
}

// internal/server/http.go

package server

import (
    "fmt"
    "net/http"
    "strings"
    "go-sms-alert/internal/handler"
)

var commandHandlers = []handler.CommandHandler{
    handler.StatusCommand{},
    handler.RebootCommand{},
}

func StartHTTPServer() {
    http.HandleFunc("/sms", func(w http.ResponseWriter, r *http.Request) {
        if err := r.ParseForm(); err != nil {
            http.Error(w, "Bad request", http.StatusBadRequest)
            return
        }

        body := strings.TrimSpace(r.FormValue("Body"))
        response := "Commande inconnue."
        for _, h := range commandHandlers {
            if h.CanHandle(body) {
                response = h.Handle(body)
                break
            }
        }

        w.Header().Set("Content-Type", "application/xml")
        fmt.Fprintf(w, `<Response><Message>%s</Message></Response>`, response)
    })

    fmt.Println("Serveur en écoute sur http://localhost:8080/sms")
    http.ListenAndServe(":8080", nil)
}

// main.go

package main

import "go-sms-alert/internal/server"

func main() {
    server.StartHTTPServer()
}
```

Étape 3 : Lancer le serveur
go run main.go
Étape 4 : Ajouter un test unitaire (optionnel mais recommandé)


test/handler_test.go

```go
package handler_test

import (
    "go-sms-alert/internal/handler"
    "testing"
)

func TestStatusCommand(t *testing.T) {
    cmd := handler.StatusCommand{}
    if !cmd.CanHandle("status") {
        t.Error("status should be handled")
    }
    if cmd.Handle("status") != "Système opérationnel." {
        t.Error("unexpected status response")
    }
}
```
