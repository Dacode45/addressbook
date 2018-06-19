package mock

import (
	"github.com/Dacode45/addressbook/models"
	"github.com/icrowley/fake"
)

// FakeContacts generates count contacts
func FakeContacts(count int) []models.Contact {
	contacts := make([]models.Contact, count)
	for i := 0; i < count; i++ {
		contacts[i] = models.Contact{
			FirstName: fake.FirstName(),
			LastName:  fake.LastName(),
			Email:     fake.EmailAddress(),
			Phone:     fake.Phone(),
		}
	}
	return contacts
}
