package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dacode45/addressbook/models"
	"github.com/Dacode45/addressbook/storage"
	"github.com/gorilla/mux"
)

func AllContactsEndPoint(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	movies, err := storage.DB.FindAll(ctx)
	if err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusOK.Serve(movies)(w, r)
}

func FindContactEndPoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	ctx := r.Context()
	movie, err := storage.DB.FindById(ctx, params["id"])
	if err != nil {
		StatusBadRequest.Serve(fmt.Errorf("Invalid Movie ID"))(w, r)
		return
	}
	StatusOK.Serve(movie)(w, r)
}

func CreateContactEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		StatusBadRequest.Serve(err)(w, r)
		return
	}
	contact.ID = storage.DB.NewObjectId()
	if err := storage.DB.Insert(ctx, contact); err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusCreated.Serve(contact)(w, r)
}

func UpdateContactEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		StatusBadRequest.Serve(fmt.Errorf("Invalid request payload"))(w, r)
		return
	}
	if err := storage.DB.Update(ctx, contact); err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusOK.Serve(map[string]string{"result": "success"})
}

func DeleteContactEndPoint(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	ctx := r.Context()
	var contact models.Contact
	if err := json.NewDecoder(r.Body).Decode(&contact); err != nil {
		StatusBadRequest.Serve(fmt.Errorf("Invalid request payload"))(w, r)
		return
	}
	if err := storage.DB.Delete(ctx, contact); err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}
	StatusOK.Serve(map[string]string{"result": "success"})(w, r)
}
