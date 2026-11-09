package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

type APIResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, resp APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Println("failed to encode response:", err)
	}
}

func RespondError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, APIResponse{Success: false, Message: message})
}

func RespondSuccess(w http.ResponseWriter, status int, message string, data interface{}) {
	writeJSON(w, status, APIResponse{Success: true, Message: message, Data: data})
}
