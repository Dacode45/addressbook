package storage

import (
	"context"

	"github.com/Dacode45/addressbook/models"
	"gopkg.in/mgo.v2/bson"
)

type ContactDB interface {
	NewObjectId() bson.ObjectId
	FindAll(context.Context) ([]models.Contact, error)
	FindById(context.Context, string) (models.Contact, error)
	Insert(context.Context, models.Contact) error
	Delete(context.Context, models.Contact) error
	Update(context.Context, models.Contact) error
	Connect()
}

var DB ContactDB

func SetContactDatabase(db ContactDB) {
	DB = db
}
