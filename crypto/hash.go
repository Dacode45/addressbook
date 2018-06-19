package crypto

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Hash uses bcrypt for encryptions
type Hash struct{}

// store hashes with their salt by splitting with the following deliminator
var delim = "||"

// Generate generates a hash of a string, and appends the salt used
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

// Compare checks a hash and a string. hash must have the salt appended with the deliminator
func (c *Hash) Compare(hash string, s string) error {
	parts := strings.Split(hash, delim)
	if len(parts) != 2 {
		return fmt.Errorf("Invalid hash, must have 2 parts")
	}
	incoming := []byte(s + parts[1])
	existing := []byte(parts[0])
	return bcrypt.CompareHashAndPassword(existing, incoming)
}
