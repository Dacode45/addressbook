package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"

	"github.com/Dacode45/addressbook/models"
	"github.com/Dacode45/addressbook/server"

	"github.com/Dacode45/addressbook/mock"
	"github.com/Dacode45/addressbook/storage"
)

func Test_ContactRouter(t *testing.T) {
	t.Run("test contact api", should_retrieve_contacts)
}

func should_retrieve_contacts(t *testing.T) {
	session, uStorage := newStorage()
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	fakeUser, fakeContacts := populateDatabase(uStorage)

	cRouter := server.NewContactRouter(uStorage, config, mux.NewRouter())
	uRouter := server.NewUserRouter(uStorage, config, mux.NewRouter())

	// login
	creds, _ := json.Marshal(models.Credentials{
		Username: fakeUser.Username,
		Password: fakeUser.Password,
	})
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(creds))
	res := httptest.NewRecorder()
	uRouter.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	// get jwt token
	var token server.JWTToken
	err := json.NewDecoder(res.Body).Decode(&token)
	assert.NoError(t, err, "Failed to parse response")

	// Test the find all contacts route
	res = testEndpoint("GET", "/", nil, cRouter, token)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	var fetchedContacts []models.Contact
	err = json.NewDecoder(res.Body).Decode(&fetchedContacts)
	assert.NoError(t, err, "Failed to parse response")
	assert.Equal(t, 10, len(fetchedContacts), "Failed to fetched contacts")

	// Test finding one contact
	res = testEndpoint("GET", fmt.Sprintf("/%s", fetchedContacts[0].ID), nil, cRouter, token)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	var fetchedContact models.Contact
	err = json.NewDecoder(res.Body).Decode(&fetchedContact)
	assert.NoError(t, err, "Failed to parse response")
	assert.Equal(t, fakeContacts[0], fetchedContact, "contacts aren't equal")

	// Test creating a new contact
	newContact := mock.FakeContacts(1)[0]
	newContactStr, _ := json.Marshal(newContact)
	res = testEndpoint("POST", "/", bytes.NewBuffer(newContactStr), cRouter, token)
	assert.Equal(t, http.StatusOK, res.Code, "Ok response is expected")

	var parsedContact models.Contact
	err = json.NewDecoder(res.Body).Decode(&parsedContact)
	newContact.ID = parsedContact.ID
	assert.NoError(t, err, "Failed to parse the response")
	assert.Equal(t, newContact, parsedContact, "Contacts aren't equal")

	// Test updating contact
	newContact.FirstName = "test"
	newContact.LastName = "user"
	newContactStr, _ = json.Marshal(newContact)
	res = testEndpoint("POST", fmt.Sprintf("/%s", newContact.ID), bytes.NewBuffer(newContactStr), cRouter, token)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")

	err = json.NewDecoder(res.Body).Decode(&parsedContact)
	assert.NoError(t, err, "Failed to parse response")
	assert.Equal(t, newContact, parsedContact, "Contacts aren't equal")

	// Test deleting contact
	res = testEndpoint("DELETE", fmt.Sprintf("/%s", newContact.ID), nil, cRouter, token)
	t.Logf("%+v", res)
	assert.Equal(t, http.StatusOK, res.Code, "OK response is expected")
}

func testEndpoint(method string, url string, body io.Reader, router *mux.Router, token server.JWTToken) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", token.Token))
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)
	return res
}

func populateDatabase(uStorage storage.UserStorage) (models.User, []models.Contact) {
	ctx := context.Background()
	users := mock.FakeUsers(1)
	contacts := mock.FakeContacts(10)
	uStorage.Insert(ctx, users[0])
	for i, contact := range contacts {
		c, _ := uStorage.CreateContact(ctx, users[0].Username, contact)
		contacts[i] = *c
	}
	return users[0], contacts
}
