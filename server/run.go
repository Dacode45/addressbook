package server

import (
	"github.com/Dacode45/addressbook/storage"
	"github.com/gorilla/mux"
)

func SetupDB() {
	storage.SetContactDatabase(&storage.ContactMongo{
		Server:   "localhost",
		Database: "addressbook",
	})
	storage.DB.Connect()
}

func SetupRouter() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/contacts", AllContactsEndPoint).Methods("GET")
	router.HandleFunc("/api/v1/contacts/{id}", FindContactEndPoint).Methods("GET")
	router.HandleFunc("/api/v1/contacts", CreateContactEndPoint).Methods("POST")
	router.HandleFunc("/api/v1/contacts/{id}", UpdateContactEndPoint).Methods("POST")
	router.HandleFunc("/api/v2/contacts/{id}", DeleteContactEndPoint).Methods("DELETE")
	return router
}
