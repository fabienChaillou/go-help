package api

import (
	"encoding/json"
	"net/http"

	"github.com/example/terraform-module/terraform"
)

// Handler exposes the single POST /terraform/command route
type Handler struct {
	dispatcher *terraform.Dispatcher
}

func NewHandler(d *terraform.Dispatcher) *Handler {
	return &Handler{dispatcher: d}
}

// APIResponse is the unified response envelope
type APIResponse struct {
	Success bool                    `json:"success"`
	Result  *terraform.CommandResult `json:"result,omitempty"`
	Error   string                  `json:"error,omitempty"`
}

// Register wires the route onto an existing mux
func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("POST /terraform/command", h.HandleCommand)
}

// HandleCommand is the single POST endpoint
//
//	POST /terraform/command
//	Content-Type: application/json
//
//	{ "action": "workspace", "sub_command": "new", "name": "staging" }
func (h *Handler) HandleCommand(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ── Parse body ──────────────────────────────────────────────────────────
	var payload map[string]any
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		writeJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   "invalid JSON: " + err.Error(),
		})
		return
	}

	// ── Dispatch ────────────────────────────────────────────────────────────
	result := h.dispatcher.Dispatch(payload)

	status := http.StatusOK
	if !result.Valid {
		status = http.StatusUnprocessableEntity
	}

	writeJSON(w, status, APIResponse{
		Success: result.Valid,
		Result:  &result,
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
