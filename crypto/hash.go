package crypto

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Hash struct{}

var delim = "||"

func (c *Hash) Generate(s string) (string, error) {
	salt := uuid.New().String()
	saltedBytes := []byte(s + salt)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash + delim + salt, nil
}

func (c *Hash) Compare(hash string, s string) error {
	parts := strings.Split(hash, delim)
	if len(parts) != 2 {
		return fmt.Errorf("Invalid hash, must have 2 parts")
	}
	incoming := []byte(s + parts[1])
	existing := []byte(parts[0])
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
