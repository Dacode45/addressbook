package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Dacode45/addressbook/models"
	"github.com/Dacode45/addressbook/storage"
	"github.com/gorilla/mux"
)

// userRouter handles the user routes
type userRouter struct {
	userStorage storage.UserStorage
	jwtCoder    *JWTCoder
}

// NewUserRouter creates a new userRouter
func NewUserRouter(u storage.UserStorage, config ServerConfig, router *mux.Router) *mux.Router {
	jwtCoder := NewJWTCoder(config.JWTSecret)
	userRouter := userRouter{u, jwtCoder}

	router.HandleFunc("", userRouter.CreateUserHandler).Methods("POST")
	router.HandleFunc("/login", userRouter.LoginHandler).Methods("POST")
	router.HandleFunc("/me", LoggedInMiddleware(jwtCoder, u, userRouter.GetLoggedInUser)).Methods("GET")
	router.HandleFunc("/{username}", userRouter.GetUserHandler).Methods("GET")
	return router
}

// CreateUserHandler creates a user from a POST request
func (ur *userRouter) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
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
	user.Password = ""

	StatusOK.Serve(user)(w, r)
}

// GetUserHandler gets a
func (ur *userRouter) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := ur.userStorage.FindByUsername(r.Context(), username)
	if err != nil {
		StatusNotFound.Serve(err)(w, r)
		return
	}
	user.Password = ""

	StatusOK.Serve(user)(w, r)
}

// GetLoggedInUser retireves the currently logged in user
func (ur *userRouter) GetLoggedInUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, ok := ctx.Value(ContextUserKey).(*models.User)
	if !ok || user == nil {
		StatusUnauthorized.Serve(fmt.Errorf("authentication failed"))(w, r)
		return
	}
	user.Password = ""
	StatusOK.Serve(user)(w, r)
}

// LoginHandler logs the user in and returns the JWT token for subsequent request
func (ur *userRouter) LoginHandler(w http.ResponseWriter, r *http.Request) {
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

// decodeUser decodes a user from the json body
func decodeUser(r *http.Request) (models.User, error) {
	var u models.User
	if r.Body == nil {
		return u, fmt.Errorf("no request body")
	}
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&u)
	return u, err
}

// decodeCredentials decodes a Credentials object from the json body
func decodeCredentials(r *http.Request) (models.Credentials, error) {
	var c models.Credentials
	if r.Body == nil {
		return c, fmt.Errorf("no request body")
	}
	err := json.NewDecoder(r.Body).Decode(&c)
	return c, err
}
