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
