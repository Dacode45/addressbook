package storage

import (
	"context"

	"github.com/Dacode45/addressbook/common"
	"github.com/Dacode45/addressbook/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type mongoUser struct {
	UserId   bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Username string        `bson:"username" json:"username"`
	Password string        `bson:"password" json:"password"`
}

func (u *mongoUser) toModel() *models.User {
	return &models.User{
		UserId:   u.UserId.Hex(),
		Username: u.Username,
		Password: u.Password,
	}
}

func usernameIndex() mgo.Index {
	return mgo.Index{
		Key:        []string{"username"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
}

func newMongoUser(u *models.User) *mongoUser {
	return &mongoUser{
		Username: u.Username,
		Password: u.Password,
	}
}

type MongoUserStorage struct {
	collection *mgo.Collection
	hash       common.Hash
}

func NewMongoUserStorage(session *MongoSession, dbName string, collectionName string, hash common.Hash) UserStorage {
	collection := session.GetCollection(dbName, collectionName)
	collection.EnsureIndex(usernameIndex())
	return &MongoUserStorage{
		collection,
		hash,
	}
}

func (s *MongoUserStorage) Login(ctx context.Context, c models.Credentials) (models.User, error) {
	model := mongoUser{}
	err := s.collection.Find(bson.M{"username": c.Username}).One(&model)

	err = s.hash.Compare(model.Password, c.Password)
	if err != nil {
		return models.User{}, err
	}

	return *model.toModel(), nil
}

func (s *MongoUserStorage) FindAll(ctx context.Context) ([]models.User, error) {
	var mUsers []mongoUser
	var users []models.User
	err := s.collection.Find(bson.M{}).All(&mUsers)
	for _, m := range mUsers {
		users = append(users, *m.toModel())
	}
	return users, err
}

func (s *MongoUserStorage) FindByUsername(ctx context.Context, username string) (models.User, error) {
	var model mongoUser
	err := s.collection.Find(bson.M{"username": username}).One(&model)
	return *model.toModel(), err
}

func (s *MongoUserStorage) Insert(ctx context.Context, user models.User) error {
	u := newMongoUser(&user)
	hashedPassword, err := s.hash.Generate(user.Password)
	if err != nil {
		return err
	}
	u.Password = hashedPassword
	return s.collection.Insert(newMongoUser(&user))
}

func (s *MongoUserStorage) Delete(ctx context.Context, user models.User) error {
	return s.collection.Remove(newMongoUser(&user))
}
