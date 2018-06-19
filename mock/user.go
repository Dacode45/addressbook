package mock

import (
	"fmt"

	"github.com/Dacode45/addressbook/models"
	"github.com/icrowley/fake"
)

// FakeUsers generates count users
func FakeUsers(count int) []models.User {
	users := make([]models.User, count)
	for i := 0; i < count; i++ {
		users[i] = models.User{
			Username: fmt.Sprintf("%s%d", fake.UserName(), i),
			Password: fake.SimplePassword(),
		}
	}
	return users
}
