package storage_test

import (
	"context"
	"log"
	"testing"

	"github.com/Dacode45/addressbook/mock"
	"github.com/Dacode45/addressbook/models"
	"github.com/Dacode45/addressbook/storage"
	"github.com/stretchr/testify/assert"
)

const (
	mongoUrl           = "localhost"
	dbName             = "addressbook_test_db"
	userCollectionName = "user"
)

func Test_MongoUserStorage(t *testing.T) {
	t.Run("Insert user", should_insert_user)
	t.Run("Query users", should_query_users)
	t.Run("Query contacts", should_query_contacts)
}

func should_insert_user(t *testing.T) {
	session, err := storage.NewMongoSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	mockHash := mock.Hash{}
	uStorage := storage.NewMongoUserStorage(session.Copy(), dbName, userCollectionName, &mockHash)

	user := models.User{
		Username: "test_user",
		Password: "test_password",
	}
	ctx := context.Background()
	err = uStorage.Insert(ctx, user)

	assert.NoError(t, err, "Unable to create user")
	// ensure 1
	var results []models.User
	session.GetCollection(dbName, userCollectionName).Find(nil).All(&results)

	count := len(results)
	assert.Equal(t, count, 1, "Incorrect number of results.")
	assert.Equal(t, results[0].Username, user.Username, "Incorrect Username")
}

func should_query_users(t *testing.T) {
	session, err := storage.NewMongoSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	mockHash := mock.Hash{}
	uStorage := storage.NewMongoUserStorage(session.Copy(), dbName, userCollectionName, &mockHash)

	numFake := 10

	users := mock.FakeUsers(numFake)
	ctx := context.Background()
	for _, u := range users {
		assert.NoError(t, uStorage.Insert(ctx, u), "Unable to create user")
	}

	// Test find all
	var findAll []models.User
	findAll, err = uStorage.FindAll(ctx)
	t.Logf("%+v", findAll)
	assert.NoError(t, err, "Failed to find a users")
	assert.Equal(t, len(findAll), numFake, "Incorrect number of users")

	// Find by username
	var byUsername *models.User
	byUsername, err = uStorage.FindByUsername(ctx, findAll[0].Username)
	assert.NoError(t, err, "Failed to find a user")
	assert.Equal(t, byUsername.Username, findAll[0].Username, "Incorrect user fetched")

	// Deletion
	err = uStorage.Delete(ctx, findAll[0].Username)
	assert.NoError(t, err, "Failed to delete user")
	byUsername, err = uStorage.FindByUsername(ctx, findAll[0].Username)
	assert.Error(t, err, "Failed to delete user")
}

func should_query_contacts(t *testing.T) {
	session, err := storage.NewMongoSession(mongoUrl)
	if err != nil {
		log.Fatalf("Unable to connect to mongo: %s", err)
	}
	defer func() {
		session.DropDatabase(dbName)
		session.Close()
	}()

	mockHash := mock.Hash{}
	uStorage := storage.NewMongoUserStorage(session.Copy(), dbName, userCollectionName, &mockHash)

	numFake := 10

	fakeUser := mock.FakeUsers(1)[0]
	fakeContacts := mock.FakeContacts(numFake)

	ctx := context.Background()

	assert.NoError(t, uStorage.Insert(ctx, fakeUser), "Unable to create user")
	for i, contact := range fakeContacts {
		c, e := uStorage.CreateContact(ctx, fakeUser.Username, contact)
		assert.NoError(t, e, "Failed to insert Contact")
		fakeContacts[i] = *c
	}
	t.Logf("fake contacts: %+v", fakeContacts)

	// Check that we can find all
	var allContacts []models.Contact
	allContacts, err = uStorage.FindAllContacts(ctx, fakeUser.Username)
	assert.NoError(t, err, "Failed to fetch all contacts")
	assert.Equal(t, numFake, len(allContacts))

	// Check that we can find by id
	var fakeContact *models.Contact
	fakeContact, err = uStorage.FindContactById(ctx, fakeUser.Username, fakeContacts[0].ID)
	assert.NoError(t, err, "Failed to fetch contact")
	assert.Equal(t, *fakeContact, fakeContacts[0], "Incorrect contact fetched")

	// check that we can update
	t.Logf("fakeContact: %+v", fakeContact)
	fakeContact.FirstName = "test"
	fakeContact.LastName = "user"
	err = uStorage.UpdateContact(ctx, fakeUser.Username, *fakeContact)
	assert.NoError(t, err, "Failed to update")
	var update *models.Contact
	update, err = uStorage.FindContactById(ctx, fakeUser.Username, fakeContact.ID)
	assert.NoError(t, err, "Failed to find contact")
	assert.Equal(t, fakeContact, update, "failed to update contact")

	// check that we can delete
	err = uStorage.DeleteContact(ctx, fakeUser.Username, fakeContact.ID)
	assert.NoError(t, err, "Failed to delete")
	fakeContact, err = uStorage.FindContactById(ctx, fakeUser.UserID, fakeContact.ID)
	assert.Error(t, err, "failed to delete contact")
}
