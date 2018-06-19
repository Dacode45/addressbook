package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Dacode45/addressbook/models"

	"github.com/gorilla/mux"

	"github.com/Dacode45/addressbook/mock"
	"github.com/Dacode45/addressbook/server"
	"github.com/Dacode45/addressbook/storage"
)

const (
	mongoUrl           = "localhost"
	dbName             = "addressbook_router_test_db"
	userCollectionName = "user"
)

var config = server.ServerConfig{
	JWTSecret: "secret",
}

func Test_UserRouter(t *testing.T) {
	t.Run("test user creation", should_create_user)
	t.Run("test user retrieval", should_retrieve_user)
	t.Run("test should login", should_login_user)
}

func should_create_user(t *testing.T) {
	session, uStorage := newStorage()

	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	router := server.NewUserRouter(uStorage, config, mux.NewRouter())

	user := models.User{
		Username: "testUser",
		Password: "testPassword",
	}
	u, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(u))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	var dbUser *models.User
	var err error
	dbUser, err = uStorage.FindByUsername(context.Background(), user.Username)
	t.Log(dbUser)
	assert.NoError(t, err, "Failed to retrieve user")
}

func should_retrieve_user(t *testing.T) {
	session, uStorage := newStorage()

	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	router := server.NewUserRouter(uStorage, config, mux.NewRouter())

	user := models.User{
		Username: "testUser",
		Password: "testPassword",
	}
	u, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(u))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	var dbUser *models.User
	var err error
	dbUser, err = uStorage.FindByUsername(context.Background(), user.Username)
	t.Log(dbUser)
	assert.NoError(t, err, "Failed to retrieve user")
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	req, _ = http.NewRequest("GET", fmt.Sprintf("/%s", dbUser.Username), nil)
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	var fetched models.User
	err = json.NewDecoder(res.Body).Decode(&fetched)
	assert.NoError(t, err, "Failed to parse response")
	assert.Equal(t, fetched.Username, dbUser.Username, "Unexpected result when fetching")
}

func should_login_user(t *testing.T) {
	session, uStorage := newStorage()

	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	router := server.NewUserRouter(uStorage, config, mux.NewRouter())

	user := models.User{
		Username: "testUser",
		Password: "testPassword",
	}

	// Ensure that you can't log in first
	creds, _ := json.Marshal(models.Credentials{
		Username: user.Username,
		Password: user.Password,
	})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(creds))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	u, _ := json.Marshal(user)
	req, _ = http.NewRequest("POST", "/", bytes.NewBuffer(u))
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)

	var dbUser *models.User
	var err error
	dbUser, err = uStorage.FindByUsername(context.Background(), user.Username)
	t.Log(dbUser)
	assert.NoError(t, err, "Failed to retrieve user")
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	req, _ = http.NewRequest("POST", "/login", bytes.NewBuffer(creds))
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	coder := server.NewJWTCoder(config.JWTSecret)

	var token server.JWTToken
	err = json.NewDecoder(res.Body).Decode(&token)
	assert.NoError(t, err, "Failed to parse response")

	var decodedCreds *models.Credentials
	decodedCreds, err = coder.Decode(token.Token)
	assert.NoError(t, err, "Failed to parse jwt token")
	assert.Equal(t, decodedCreds.Username, dbUser.Username, "Unexpected result when logging in")

	// test the me route
	req, _ = http.NewRequest("GET", "/me", nil)
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token.Token))
	res = httptest.NewRecorder()
	router.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")
	t.Logf("%+v", res)

	var fetched models.User
	err = json.NewDecoder(res.Body).Decode(&fetched)
	assert.NoError(t, err, "Failed to parse response")
	assert.Equal(t, decodedCreds.Username, fetched.Username, "Unexpected result when fetching")
}

func newStorage() (*storage.MongoSession, storage.UserStorage) {
	session, err := storage.NewMongoSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}

	mockHash := mock.Hash{}
	uStorage := storage.NewMongoUserStorage(session.Copy(), dbName, userCollectionName, &mockHash)
	return session, uStorage
}
