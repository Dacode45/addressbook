package storage

import (
	"context"
	"fmt"

	"github.com/Dacode45/addressbook/common"
	"github.com/Dacode45/addressbook/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// mongoContact is a mongodb specific implementaiton of the Contact struct
type mongoContact struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	FirstName string        `bson:"first_name" json:"first_name"`
	LastName  string        `bson:"last_name" json:"last_name"`
	Email     string        `bson:"email" json:"email"`
	Phone     string        `bson:"phone" json:"phone"`
}

// newMOngoContact creates a new MongodbContact from a Contact
func newMongoContact(c models.Contact, newID bool) *mongoContact {
	id := bson.ObjectId(c.ID)
	if newID {
		id = bson.NewObjectId()
	}
	return &mongoContact{
		ID:        id,
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
	}

}

// mongoContacts is a utility type for slices of mongoContacts
type mongoContacts []mongoContact

// finds a contact in list by id
func (contacts mongoContacts) findByID(id string) *mongoContact {
	if !bson.IsObjectIdHex(id) {
		return nil
	}
	for _, c := range contacts {
		if c.ID == bson.ObjectIdHex(id) {
			return &c
		}
	}
	return nil
}

// replaces (used for updating) an element of the contact slice with another
func (contacts mongoContacts) replaceWith(contact mongoContact) mongoContacts {
	update := contacts
	for i, c := range contacts {
		if c.ID == contact.ID {
			left, right := update[:i], update[i+1:]
			update = append(left, contact)
			update = append(update, right...)
			return update
		}
	}
	return update
}

// removes an elment of the contact slice all together
func (contacts mongoContacts) removeID(id string) mongoContacts {
	update := contacts
	for i, c := range contacts {
		if c.ID == bson.ObjectId(id) {
			update = append(contacts[:i], contacts[i+1:]...)
			return update
		}
	}
	return update
}

// converts to the Contact struct
func (c *mongoContact) toModel() *models.Contact {
	return &models.Contact{
		ID:        c.ID.Hex(),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
	}
}

// mongoUser creates a mongodb specific User
type mongoUser struct {
	UserID   bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username string        `bson:"username" json:"username"`
	Password string        `bson:"password" json:"password"`
	Contacts mongoContacts
}

// toModel transforms the mongo user to a User struct
func (u *mongoUser) toModel() *models.User {
	contacts := make([]models.Contact, len(u.Contacts))
	for i, c := range u.Contacts {
		contacts[i] = *c.toModel()
	}
	return &models.User{
		UserID:   u.UserID.Hex(),
		Username: u.Username,
		Password: u.Password,
		Contacts: contacts,
	}
}

// usernameIndex creates an index on the username field
func usernameIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

// newMongoUser creates a new mongoUser
func newMongoUser(u *models.User) *mongoUser {
	return &mongoUser{
		Username: u.Username,
		Password: u.Password,
	}
}

// MongoUserStorage implements the UserStorage interface
type MongoUserStorage struct {
	collection *mgo.Collection
	hash       common.Hash
}

// NewMongoUserStorage creates a new storage based of a session, database name, and collection name, as well as a password encoding hash
func NewMongoUserStorage(session *MongoSession, dbName string, collectionName string, hash common.Hash) UserStorage {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(usernameIndex())
	return &MongoUserStorage{
		collection,
		hash,
	}
}

// Login logs in a user
func (s *MongoUserStorage) Login(ctx context.Context, c models.Credentials) (*models.User, error) {
	model := mongoUser{}
	err := s.collection.Find(bson.M{"username": c.Username}).One(&model)
	err = s.hash.Compare(model.Password, c.Password)
	if err != nil {
		return nil, err
	}

	return model.toModel(), nil
}

// FindAll finds all users
func (s *MongoUserStorage) FindAll(ctx context.Context) ([]models.User, error) {
	var mUsers []mongoUser
	var users []models.User
	err := s.collection.Find(bson.M{}).All(&mUsers)
	for _, m := range mUsers {
		users = append(users, *m.toModel())
	}
	return users, err
}

// FindByUsername finds a user by username
func (s *MongoUserStorage) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var model mongoUser
	err := s.collection.Find(bson.M{"username": username}).One(&model)
	return model.toModel(), err
}

// Insert inserts a user into the db
func (s *MongoUserStorage) Insert(ctx context.Context, user models.User) error {
	u := newMongoUser(&user)
	hashedPassword, err := s.hash.Generate(user.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return s.collection.Insert(u)
}

// Delete removes a user from the db
func (s *MongoUserStorage) Delete(ctx context.Context, username string) error {
	return s.collection.Remove(bson.M{"username": username})
}

// Contact methods

// getUser gets a user from the databse. Utility function
func (s *MongoUserStorage) getUser(ctx context.Context, username string) (*mongoUser, error) {
	var user mongoUser
	err := s.collection.Find(bson.M{"username": username}).One(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateContact creates a new contact
func (s *MongoUserStorage) CreateContact(ctx context.Context, username string, contact models.Contact) (*models.Contact, error) {
	user, err := s.getUser(ctx, username)
	if err != nil {
		return nil, err
	}
	newContact := newMongoContact(contact, true)
	user.Contacts = append(user.Contacts, *newContact)
	err = s.collection.Update(bson.M{"_id": user.UserID}, bson.M{"$set": bson.M{"contacts": user.Contacts}})
	return newContact.toModel(), err
}

// FindContactById retrieves the specified contact
func (s *MongoUserStorage) FindContactById(ctx context.Context, username string, contactID string) (*models.Contact, error) {
	user, err := s.getUser(ctx, username)
	if err != nil {
		return nil, err
	}
	contact := user.Contacts.findByID(contactID)
	if contact == nil {
		return nil, fmt.Errorf("No contact with that id")
	}
	return contact.toModel(), nil
}

// FindAllContacts finds all the contacts of a user
func (s *MongoUserStorage) FindAllContacts(ctx context.Context, username string) ([]models.Contact, error) {
	user, err := s.getUser(ctx, username)
	if err != nil {
		return nil, err
	}
	return user.toModel().Contacts, nil
}

// UpdateContact updates a specific contact of a user
func (s *MongoUserStorage) UpdateContact(ctx context.Context, username string, update models.Contact) error {
	if !bson.IsObjectIdHex(update.ID) {
		return fmt.Errorf("Invalid id")
	}
	user, err := s.getUser(ctx, username)
	if err != nil {
		return nil
	}
	oldID := update.ID
	contact := newMongoContact(update, false)
	contact.ID = bson.ObjectIdHex(oldID)
	contacts := user.Contacts.replaceWith(*contact)
	err = s.collection.Update(bson.M{"_id": user.UserID}, bson.M{"$set": bson.M{"contacts": contacts}})
	return err
}

// DeleteContact deletes the contact
func (s *MongoUserStorage) DeleteContact(ctx context.Context, username string, contactID string) error {
	user, err := s.getUser(ctx, username)
	if err != nil {
		return nil
	}
	contacts := user.Contacts.removeID(contactID)
	err = s.collection.Update(bson.M{"_id": user.UserID}, bson.M{"$set": bson.M{"contacts": contacts}})
	return err
}
