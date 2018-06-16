package storage

import (
	"context"
	"fmt"
	"log"

	"github.com/Dacode45/addressbook/models"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type ContactMongo struct {
	Server   string
	Database string
	db       *mgo.Database
}

const (
	COLLECTION = "contacts"
)

func (c *ContactMongo) Connect() {
	session, err := mgo.Dial(c.Server)
	if err != nil {
		log.Fatal(err)
	}
	c.db = session.DB(c.Database)
}

func (c *ContactMongo) FindAll(ctx context.Context) ([]models.Contact, error) {
	var contacts []models.Contact
	err := c.db.C(COLLECTION).Find(bson.M{}).All(&contacts)
	return contacts, err
}

func (c *ContactMongo) FindById(ctx context.Context, id string) (models.Contact, error) {
	var contact models.Contact
	if !bson.IsObjectIdHex(id) {
		return contact, fmt.Errorf("Invalid Id")
	}
	err := c.db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&contact)
	return contact, err
}

func (c *ContactMongo) Insert(ctx context.Context, contact models.Contact) error {
	err := c.db.C(COLLECTION).Insert(&contact)
	return err
}

func (c *ContactMongo) Delete(ctx context.Context, contact models.Contact) error {
	err := c.db.C(COLLECTION).Remove(&contact)
	return err
}

func (c *ContactMongo) Update(ctx context.Context, contact models.Contact) error {
	err := c.db.C(COLLECTION).UpdateId(contact.ID, &contact)
	return err
}

func (c *ContactMongo) NewObjectId() bson.ObjectId {
	return bson.NewObjectId()
}
