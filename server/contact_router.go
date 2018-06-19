package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dacode45/addressbook/models"

	"github.com/Dacode45/addressbook/storage"
	"github.com/gorilla/mux"
)

type contactRouter struct {
	userStorage storage.UserStorage
	jwtCoder    *JWTCoder
}

func NewContactRouter(u storage.UserStorage, config ServerConfig, router *mux.Router) *mux.Router {
	jwtCoder := NewJWTCoder(config.JWTSecret)
	cr := contactRouter{u, jwtCoder}

	router.HandleFunc("/", LoggedInMiddleware(jwtCoder, u, cr.allContactsEndPoint)).Methods("GET")
	router.HandleFunc("/{id}", LoggedInMiddleware(jwtCoder, u, cr.findContactEndPoint)).Methods("GET")
	router.HandleFunc("/", LoggedInMiddleware(jwtCoder, u, cr.createContactEndPoint)).Methods("POST")
	router.HandleFunc("/{id}", LoggedInMiddleware(jwtCoder, u, cr.updateContactEndPoint)).Methods("POST")
	router.HandleFunc("/{id}", LoggedInMiddleware(jwtCoder, u, cr.deleteContactEndPoint)).Methods("DELETE")
	return router
}

func (cr *contactRouter) allContactsEndPoint(w http.ResponseWriter, r *http.Request) {
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

func (cr *contactRouter) findContactEndPoint(w http.ResponseWriter, r *http.Request) {
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

func (cr *contactRouter) createContactEndPoint(w http.ResponseWriter, r *http.Request) {
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

func (cr *contactRouter) updateContactEndPoint(w http.ResponseWriter, r *http.Request) {
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

func (cr *contactRouter) deleteContactEndPoint(w http.ResponseWriter, r *http.Request) {
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

func decodeContact(r *http.Request) (models.Contact, error) {
	var c models.Contact
	if r.Body == nil {
		return c, fmt.Errorf("no request body")
	}
	err := json.NewDecoder(r.Body).Decode(&c)
	return c, err
}
