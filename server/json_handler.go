package server

import (
	"encoding/json"
	"net/http"
)

// JsonHandler serves responses as json
type JSONHandler int

// Serve payload as json
func (j JSONHandler) Serve(payload interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := json.Marshal(payload)
		if err != nil {
			ServerErrorHandler.Serve(w, r)
			return
		}
		code := int(j)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(code)
		w.Write(msg)
	}
}
