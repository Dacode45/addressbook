package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dacode45/addressbook/models"
	"github.com/gocarina/gocsv"

	"github.com/Dacode45/addressbook/storage"
	"github.com/gorilla/mux"
)

type contactRouter struct {
	userStorage storage.UserStorage
	jwtCoder    *JWTCoder
}

// NewContactRouter generates a router for handling the contacts api. Requires access to our user storage
func NewContactRouter(u storage.UserStorage, config ServerConfig, router *mux.Router) *mux.Router {
	jwtCoder := NewJWTCoder(config.JWTSecret)
	cr := contactRouter{u, jwtCoder}

	router.HandleFunc("", LoggedInMiddleware(jwtCoder, u, cr.AllContactsEndPoint)).Methods("GET")
	// export import csv
	router.HandleFunc("/export", LoggedInMiddleware(jwtCoder, u, cr.ExportAllContactsEndpoint)).Methods("GET")
	router.HandleFunc("/{id}", LoggedInMiddleware(jwtCoder, u, cr.FindContactEndPoint)).Methods("GET")
	router.HandleFunc("", LoggedInMiddleware(jwtCoder, u, cr.CreateContactEndPoint)).Methods("POST")
	router.HandleFunc("/import", LoggedInMiddleware(jwtCoder, u, cr.ImportContactsEndPoint)).Methods("POST")
	router.HandleFunc("/{id}", LoggedInMiddleware(jwtCoder, u, cr.UpdateContactEndPoint)).Methods("POST")
	router.HandleFunc("/{id}", LoggedInMiddleware(jwtCoder, u, cr.DeleteContactEndPoint)).Methods("DELETE")
	return router
}

// ExportAllContactsEndpoints exports all user contacts as csv. Limits to 1000000 byte body
func (cr *contactRouter) ExportAllContactsEndpoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(ContextUserKey).(*models.User)
	if !ok || user == nil {
		StatusUnauthorized.Serve(fmt.Errorf("Unautorized"))(w, r)
		return
	}

	contacts, err := cr.userStorage.FindAllContacts(ctx, user.Username)
	if err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusOKCSV.Serve("contacts.csv", contacts)(w, r)
}

// AllContactsEndPoint retrieves all user contacts as json
func (cr *contactRouter) AllContactsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(ContextUserKey).(*models.User)
	if !ok || user == nil {
		StatusUnauthorized.Serve(fmt.Errorf("Unautorized"))(w, r)
		return
	}

	contacts, err := cr.userStorage.FindAllContacts(ctx, user.Username)
	if err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusOK.Serve(contacts)(w, r)
}

// FindContactEndPoint searches for a given contact
func (cr *contactRouter) FindContactEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx := r.Context()
	user, ok := ctx.Value(ContextUserKey).(*models.User)
	if !ok || user == nil {
		StatusUnauthorized.Serve(fmt.Errorf("Unautorized"))(w, r)
		return
	}

	contact, err := cr.userStorage.FindContactById(ctx, user.Username, params["id"])
	if err != nil {
		StatusNotFound.Serve(err)(w, r)
		return
	}
	StatusOK.Serve(contact)(w, r)
}

// CreateContactEndPoint creates a given contact from a json body
func (cr *contactRouter) CreateContactEndPoint(w http.ResponseWriter, r *http.Request) {
	contact, err := decodeContact(r)
	if err != nil {
		StatusBadRequest.Serve(err)(w, r)
		return
	}

	ctx := r.Context()
	user, ok := ctx.Value(ContextUserKey).(*models.User)
	if !ok || user == nil {
		StatusUnauthorized.Serve(fmt.Errorf("Unautorized"))(w, r)
		return
	}

	var newContact *models.Contact
	newContact, err = cr.userStorage.CreateContact(ctx, user.Username, contact)
	if err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusOK.Serve(newContact)(w, r)
}

// ImportContactsEndPoint imports a csv file for contacts
func (cr *contactRouter) ImportContactsEndPoint(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1000000)
	contacts, err := decodeContacts(r)
	if err != nil {
		StatusBadRequest.Serve(err)(w, r)
		return
	}

	ctx := r.Context()
	user, ok := ctx.Value(ContextUserKey).(*models.User)
	if !ok || user == nil {
		StatusUnauthorized.Serve(fmt.Errorf("Unautorized"))(w, r)
		return
	}
	var newContacts = make([]*models.Contact, len(contacts))
	for i, contact := range contacts {
		var newContact *models.Contact
		newContact, err = cr.userStorage.CreateContact(ctx, user.Username, contact)
		if err != nil {
			StatusInternalServerError.Serve(err)(w, r)
			return
		}
		newContacts[i] = newContact
	}
	StatusOK.Serve(newContacts)(w, r)
}

// UpdateContactEndPoint updates the fields of a given contact
func (cr *contactRouter) UpdateContactEndPoint(w http.ResponseWriter, r *http.Request) {
	contact, err := decodeContact(r)
	if err != nil {
		StatusBadRequest.Serve(err)(w, r)
		return
	}

	ctx := r.Context()
	user, ok := ctx.Value(ContextUserKey).(*models.User)
	if !ok || user == nil {
		StatusUnauthorized.Serve(fmt.Errorf("Unautorized"))(w, r)
		return
	}

	err = cr.userStorage.UpdateContact(ctx, user.Username, contact)
	if err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusOK.Serve(contact)(w, r)
}

// DelecteContactEndPoint removes a contact
func (cr *contactRouter) DeleteContactEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx := r.Context()
	user, ok := ctx.Value(ContextUserKey).(*models.User)
	if !ok || user == nil {
		StatusUnauthorized.Serve(fmt.Errorf("Unautorized"))(w, r)
		return
	}

	err := cr.userStorage.DeleteContact(ctx, user.Username, params["id"])
	if err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusOK.Serve(map[string]string{"msg": "Success"})(w, r)
}

// decodeContact returns a contact from a json body
func decodeContact(r *http.Request) (models.Contact, error) {
	defer r.Body.Close()
	var c models.Contact
	if r.Body == nil {
		return c, fmt.Errorf("no request body")
	}
	err := json.NewDecoder(r.Body).Decode(&c)
	return c, err
}

// decodeContacts returns several contacts from a csv body
func decodeContacts(r *http.Request) ([]models.Contact, error) {
	defer r.Body.Close()
	var c []models.Contact
	if r.Body == nil {
		return c, fmt.Errorf("no request body")
	}
	err := gocsv.Unmarshal(r.Body, &c)
	return c, err
}
