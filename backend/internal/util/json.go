package util

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes any data as JSON
func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// ReadJSON reads and decodes JSON from request body
func ReadJSON(r *http.Request, dest any) error {
	return json.NewDecoder(r.Body).Decode(dest)
}
