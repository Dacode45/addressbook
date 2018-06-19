package server

import (
	"net/http"
)

var (
	StatusCreated = JsonHandler(http.StatusCreated)
	StatusOK      = JsonHandler(http.StatusOK)
	StatusOKCSV   = CSVHandler(http.StatusOK)
)
