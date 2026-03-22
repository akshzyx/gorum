package util

import (
	"log"
	"net/http"
)

// WriteJSONError writes a standardized error response
func WriteJSONError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{"error": message})
}

// BadRequest sends a 400 response with logging
func BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[400] %s %s - %v", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

// InternalServerError sends a 500 response with logging
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[500] %s %s - %v", r.Method, r.URL.Path, err)
	WriteJSONError(w, http.StatusInternalServerError, "internal server error")
}

// NotFound sends a 404 response
func NotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("[404] %s %s", r.Method, r.URL.Path)
	WriteJSONError(w, http.StatusNotFound, "not found")
}
