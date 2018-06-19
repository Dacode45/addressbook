package server

import (
	"fmt"
	"net/http"

	"github.com/gocarina/gocsv"
)

type CSVHandler int

func (c CSVHandler) Serve(filename string, payload interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msg, err := gocsv.MarshalString(payload)
		if err != nil {
			ServerErrorHandler.Serve(w, r)
			return
		}
		code := int(c)
		w.Header().Set("Content-Type", "text/csv")
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
		w.WriteHeader(code)
		w.Write([]byte(msg))
	}
}
