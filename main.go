package main

import (
	"log"
	"net/http"

	"github.com/Dacode45/addressbook/server"
	"github.com/Dacode45/addressbook/storage"

	"github.com/gorilla/mux"
)

func init() {
	storage.SetContactDatabase(&storage.ContactMongo{
		Server:   "localhost",
		Database: "addressbook",
	})
	storage.DB.Connect()
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/contacts", server.AllContactsEndPoint).Methods("GET")
	router.HandleFunc("/api/v1/contacts/{id}", server.FindContactEndPoint).Methods("GET")
	router.HandleFunc("/api/v1/contacts", server.CreateContactEndPoint).Methods("POST")
	router.HandleFunc("/api/v1/contacts/{id}", server.UpdateContactEndPoint).Methods("POST")
	router.HandleFunc("/api/v2/contacts/{id}", server.DeleteContactEndPoint).Methods("DELETE")
	// http.Handle("/", router)
	// appengine.Main()
	if err := http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}
