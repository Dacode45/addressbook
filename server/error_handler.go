package server

import (
	"net/http"
)

type ErrorHandler int

func (e ErrorHandler) Serve(err error) http.HandlerFunc {
	handler := JsonHandler(e)
	return handler.Serve(map[string]string{"error": err.Error()})
}

var (
	StatusBadRequest          = ErrorHandler(http.StatusBadRequest)
	StatusInternalServerError = ErrorHandler(http.StatusInternalServerError)
	StatusNotFound            = ErrorHandler(http.StatusNotFound)
	StatusUnauthorized        = ErrorHandler(http.StatusUnauthorized)
)
