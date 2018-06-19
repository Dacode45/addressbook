package server

import (
	"io"
	"net/http"
)

// StatusHandler serves a status as text in case we can't send json
type StatusHandler int

// Serve serves a status and it's official status text
func (s StatusHandler) Serve(w http.ResponseWriter, r *http.Request) {
	code := int(s)
	w.WriteHeader(code)
	io.WriteString(w, http.StatusText(code))
}

var (
	NotFoundHandler       = StatusHandler(http.StatusNotFound)
	ServerErrorHandler    = StatusHandler(http.StatusInternalServerError)
	NotImplementedHandler = StatusHandler(http.StatusNotImplemented)
	NotLegalHandler       = StatusHandler(http.StatusNotAcceptable)
)
