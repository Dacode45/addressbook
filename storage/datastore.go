package storage

import (
	"context"

	"github.com/Dacode45/addressbook/models"
	"google.golang.org/appengine/datastore"
)

type ContactDatastore struct {
}

func (db *ContactDatastore) Insert(ctx context.Context, contact models.Contact) error {
	key := datastore.NewIncompleteKey(ctx, "Contact", nil)
	_, err := datastore.Put(ctx, key, &contact)
	return err
}
