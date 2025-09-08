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
