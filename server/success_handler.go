package server

import (
	"net/http"
)

var (
	// StatusCreated serves json with the StatusCreatedCode
	StatusCreated = JSONHandler(http.StatusCreated)
	// StatusOK serves json with the StatusOKCode
	StatusOK = JSONHandler(http.StatusOK)
	// StatusOkCSV serves csv with the StatusOkCSVCode
	StatusOKCSV = CSVHandler(http.StatusOK)
)
