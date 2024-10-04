package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/httplog/v2"
)

type responseMsg struct {
	Message string `json:"message"`
}

type responseErr struct {
	Error            string    `json:"error,omitempty"`
}

// encodeResponse encodes data as a JSON response.
func encodeResponse(w http.ResponseWriter, logger *httplog.Logger, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Error("Error while marshaling data", "err", err, "data", data)
		http.Error(w, `{"Error": "Internal server error"}`, http.StatusInternalServerError)
	}
}