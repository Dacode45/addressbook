package server

import (
	"io"
	"net/http"
)

type StatusHandler int

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
