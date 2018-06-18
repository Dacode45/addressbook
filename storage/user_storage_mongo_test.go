package storage_test

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/icrowley/fake"

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
	t.Run("Inser User", should_insert_user)
	t.Run("Query checks", should_query_users)
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

	users := fakeUsers(numFake)
	ctx := context.Background()
	for _, u := range users {
		assert.NoError(t, uStorage.Insert(ctx, u), "Unable to create user")
	}

	// Test find all
	var findAll []models.User
	findAll, err = uStorage.FindAll(ctx)
	t.Logf("%+v", findAll)
	assert.NoError(t, err, "Failed to find all users")
	assert.Equal(t, len(findAll), numFake, "Incorrect number of users")

	// Find by username
	var byUsername models.User
	byUsername, err = uStorage.FindByUsername(ctx, findAll[0].Username)
	assert.NoError(t, err, "Failed to find a user")
	assert.Equal(t, byUsername.Username, findAll[0].Username, "Incorrect user fetched")

	// Deletion
	err = uStorage.Delete(ctx, findAll[0])
	assert.NoError(t, err, "Failed to delete user")
	byUsername, err = uStorage.FindByUsername(ctx, findAll[0].Username)
	assert.Error(t, err, "Failed to delete user")

}

func fakeUsers(count int) []models.User {
	users := make([]models.User, count)
	for i := 0; i < count; i++ {
		users[i] = models.User{
			Username: fmt.Sprintf("%s%d", fake.UserName(), i),
			Password: fake.SimplePassword(),
		}
	}
	return users
}
