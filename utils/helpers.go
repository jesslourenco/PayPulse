package utils

import (
	"encoding/json"
	"net/http"
)

func ErrorWithMessage(w http.ResponseWriter, status int, message string) {
	resp := ErrorResponse{
		Status:  status,
		Message: message,
	}

	payload, err := json.Marshal(resp)
	if err != nil {
		WithPayload(w, http.StatusUnprocessableEntity, []byte(`{"error": "Unable to Return Payload"}`))
	}
	WithPayload(w, status, payload)
}

func WithPayload(w http.ResponseWriter, status int, payload []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(payload)
}
