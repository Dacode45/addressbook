package storage

import (
	"context"

	"github.com/Dacode45/addressbook/models"
)

type UserStorage interface {
	Login(context.Context, models.Credentials) (*models.User, error)
	FindAll(context.Context) ([]models.User, error)
	FindByUsername(context.Context, string) (*models.User, error)
	Insert(context.Context, models.User) error
	Delete(context.Context, string) error

	// CRUD On Contacts
	CreateContact(context.Context, string, models.Contact) (*models.Contact, error)
	FindAllContacts(context.Context, string) ([]models.Contact, error)
	FindContactById(context.Context, string, string) (*models.Contact, error)
	UpdateContact(context.Context, string, models.Contact) error
	DeleteContact(context.Context, string, string) error
}
