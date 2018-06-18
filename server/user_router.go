package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dacode45/addressbook/models"
	"github.com/Dacode45/addressbook/storage"
	"github.com/gorilla/mux"
)

type userRouter struct {
	userStorage storage.UserStorage
	jwtCoder    *JWTCoder
}

func NewUserRouter(u storage.UserStorage, config ServerConfig, router *mux.Router) *mux.Router {
	jwtCoder := NewJWTCoder(config.JWTSecret)
	userRouter := userRouter{u, jwtCoder}

	router.HandleFunc("/", userRouter.createUserHandler).Methods("POST")
	router.HandleFunc("/login", userRouter.loginHandler).Methods("POST")
	router.HandleFunc("/me", jwtCoder.TokenAuthMiddleware(userRouter.getLoggedInUser)).Methods("GET")
	router.HandleFunc("/{username}", userRouter.getUserHandler).Methods("GET")
	return router
}

func (ur *userRouter) createUserHandler(w http.ResponseWriter, r *http.Request) {
	user, err := decodeUser(r)
	if err != nil {
		StatusBadRequest.Serve(err)(w, r)
		return
	}

	err = ur.userStorage.Insert(r.Context(), user)
	if err != nil {
		StatusInternalServerError.Serve(err)(w, r)
		return
	}

	StatusOK.Serve(user)(w, r)
}

func (ur *userRouter) getUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := ur.userStorage.FindByUsername(r.Context(), username)
	if err != nil {
		StatusNotFound.Serve(err)(w, r)
		return
	}

	StatusOK.Serve(user)(w, r)
}

func (ur *userRouter) getLoggedInUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creds, ok := ctx.Value(CREDENTIALS_KEY).(*models.Credentials)
	if !ok || creds == nil {
		StatusUnauthorized.Serve(fmt.Errorf("no jwt passed"))(w, r)
		return
	}
	user, err := ur.userStorage.Login(r.Context(), *creds)
	if err != nil {
		StatusUnauthorized.Serve(fmt.Errorf("not logged in"))(w, r)
		return
	}
	StatusOK.Serve(user)(w, r)
}

func (ur *userRouter) loginHandler(w http.ResponseWriter, r *http.Request) {
	credentials, err := decodeCredentials(r)
	if err != nil {
		StatusBadRequest.Serve(err)(w, r)
		return
	}

	_, err = ur.userStorage.Login(r.Context(), credentials)
	if err != nil {
		StatusUnauthorized.Serve(fmt.Errorf("Incorrect password"))(w, r)
		return
	}
	// User is logged in send jwt token
	var token JWTToken
	token, err = ur.jwtCoder.Create(credentials)
	if err != nil {
		StatusInternalServerError.Serve(err)(w, r)
	}
	StatusOK.Serve(token)(w, r)
}

func decodeUser(r *http.Request) (models.User, error) {
	var u models.User
	if r.Body == nil {
		return u, fmt.Errorf("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	return u, err
}

func decodeCredentials(r *http.Request) (models.Credentials, error) {
	var c models.Credentials
	if r.Body == nil {
		return c, fmt.Errorf("no request body")
	}
	err := json.NewDecoder(r.Body).Decode(&c)
	return c, err
}
