package models

import "gopkg.in/mgo.v2/bson"

type Contact struct {
	ID        bson.ObjectId `bson:"_id" json:"id"`
	FirstName string        `bson:"first_name" json:"first_name"`
	LastName  string        `bson:"last_name" json:"last_name"`
	Email     string        `bson:"email" json:"email"`
	Phone     string        `bson:"phone" json:"phone"`
}
