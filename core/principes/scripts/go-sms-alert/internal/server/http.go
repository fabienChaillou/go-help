package server

import (
	"fmt"
	"go-sms-alert/internal/handler"
	"net/http"
	"strings"
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
