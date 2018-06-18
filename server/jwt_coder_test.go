package server_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Dacode45/addressbook/models"
	"github.com/Dacode45/addressbook/server"
)

func Test_JWTCoder(t *testing.T) {
	t.Run("jwt coder test", should_code_and_decode)
}

func should_code_and_decode(t *testing.T) {
	coder := server.NewJWTCoder("secret")
	creds := models.Credentials{
		Username: "testUser",
		Password: "testPassword",
	}
	token, err := coder.Create(creds)
	assert.NoError(t, err, "Failed to sign jwt")

	var decoded *models.Credentials
	decoded, err = coder.Decode(token.Token)
	assert.NoError(t, err, "Failed to decode jwt")
	assert.Equal(t, *decoded, creds, "Encoding missmatch")
}
