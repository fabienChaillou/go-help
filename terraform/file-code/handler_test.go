package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/example/terraform-module/api"
	"github.com/example/terraform-module/terraform"
)

func newTestServer() *httptest.Server {
	registry := terraform.NewRegistry()
	dispatcher := terraform.NewDispatcher(registry)
	handler := api.NewHandler(dispatcher)
	mux := http.NewServeMux()
	handler.Register(mux)
	return httptest.NewServer(mux)
}

func postJSON(t *testing.T, server *httptest.Server, payload any) *http.Response {
	t.Helper()
	body, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	resp, err := http.Post(server.URL+"/terraform/command", "application/json", bytes.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	return resp
}

func decode(t *testing.T, resp *http.Response) map[string]any {
	t.Helper()
	var out map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		t.Fatal(err)
	}
	return out
}

func TestAPI_ValidWorkspaceList(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{
		"action":      "workspace",
		"sub_command": "list",
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body := decode(t, resp)
	if body["success"] != true {
		t.Error("expected success=true")
	}
}

func TestAPI_InvalidAction(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{"action": "nuke"})
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
	body := decode(t, resp)
	if body["success"] != false {
		t.Error("expected success=false")
	}
}

func TestAPI_InvalidJSON(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp, _ := http.Post(srv.URL+"/terraform/command", "application/json", bytes.NewBufferString("not-json"))
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", resp.StatusCode)
	}
}

func TestAPI_DestroyWithoutAutoApprove(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{"action": "destroy"})
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
	body := decode(t, resp)
	result := body["result"].(map[string]any)
	if result["help"] == "" {
		t.Error("expected help text in response")
	}
}

func TestAPI_WorkspaceNew_Valid(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{
		"action":      "workspace",
		"sub_command": "new",
		"name":        "staging",
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

func TestAPI_ImportMissingFields(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{"action": "import"})
	if resp.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", resp.StatusCode)
	}
}

func TestAPI_PlanCleanPayload(t *testing.T) {
	srv := newTestServer()
	defer srv.Close()

	resp := postJSON(t, srv, map[string]any{
		"action": "plan",
		"out":    "tfplan",
		"target": []string{"aws_instance.web"},
	})
	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
	body := decode(t, resp)
	result := body["result"].(map[string]any)
	payload := result["payload"].(map[string]any)
	if payload["out"] != "tfplan" {
		t.Error("expected out=tfplan in cleaned payload")
	}
}
