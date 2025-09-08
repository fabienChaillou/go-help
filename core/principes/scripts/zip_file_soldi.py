import os
from zipfile import ZipFile

project_files = {
    "go-sms-alert/go.mod": "module go-sms-alert\n\ngo 1.20\n",
    "go-sms-alert/main.go": """package main

import "go-sms-alert/internal/server"

func main() {
server.StartHTTPServer()
}
""",
    "go-sms-alert/internal/server/http.go": """package server

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
""",
    "go-sms-alert/internal/handler/command.go": """package handler

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
""",
    "go-sms-alert/test/handler_test.go": """package handler_test

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
""",
}

with ZipFile("go-sms-alert.zip", "w") as zipf:
    for filepath, content in project_files.items():
        os.makedirs(os.path.dirname(filepath), exist_ok=True)

with open(filepath, "w") as f:
    f.write(content)
    zipf.write(filepath)

print("Fichier zip généré : go-sms-alert.zip")
