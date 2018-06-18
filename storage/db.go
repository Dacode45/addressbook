package storage

import (
	"context"

	"github.com/Dacode45/addressbook/models"
	"gopkg.in/mgo.v2/bson"
)

type UserStorage interface {
	Login(context.Context, models.Credentials) (models.User, error)
	FindAll(context.Context) ([]models.User, error)
	FindByUsername(context.Context, string) (models.User, error)
	Insert(context.Context, models.User) error
	Delete(context.Context, models.User) error
}

type ContactStorage interface {
	NewObjectId() bson.ObjectId
	FindAll(context.Context) ([]models.Contact, error)
	FindById(context.Context, string) (models.Contact, error)
	Insert(context.Context, models.Contact) error
	Delete(context.Context, models.Contact) error
	Update(context.Context, models.Contact) error
	Connect()
}

var userStorage UserStorage
var contactStorage ContactStorage

func SetGlobalUserStorage(u UserStorage) {
	userStorage = u
}

func SetGlobalContactStorage(c ContactStorage) {
	contactStorage = c
}
