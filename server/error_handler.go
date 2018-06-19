package server

import (
	"net/http"
)

// ErrorHandler provides a consistent way of sending errors as json
type ErrorHandler int

// Serve serves an error in the format {"error": "<error>"}
func (e ErrorHandler) Serve(err error) http.HandlerFunc {
	handler := JSONHandler(e)
	return handler.Serve(map[string]string{"error": err.Error()})
}

var (
	// StatusBadRequest sets the StatusBadRequest
	StatusBadRequest = ErrorHandler(http.StatusBadRequest)
	// StatusInternalServerError sets the StatusInternalServerError
	StatusInternalServerError = ErrorHandler(http.StatusInternalServerError)
	// StatusNotFound sets the StatusNotFound
	StatusNotFound = ErrorHandler(http.StatusNotFound)
	// StatusUnauthorized sets the StatusUnauthorized
	StatusUnauthorized = ErrorHandler(http.StatusUnauthorized)
)
