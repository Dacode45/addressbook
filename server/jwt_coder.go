package server

import (
	"fmt"

	"github.com/Dacode45/addressbook/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

// JWTToken is a wrapper around an actual jwt. Useful object for serialization
type JWTToken struct {
	Token string `json:"token"`
}

// JWTCoder encodes an object into a jwy. Requires a secret
type JWTCoder struct {
	secret string
}

// NewJWTCoder creates a JWTCoder with the given secret for encryption
func NewJWTCoder(secret string) *JWTCoder {
	return &JWTCoder{
		secret,
	}
}

// Create encodes a Credential object into a JWTToken
func (j *JWTCoder) Create(c models.Credentials) (JWTToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": c.Username,
		"password": c.Password,
	})
	tokenString, err := token.SignedString([]byte(j.secret))
	return JWTToken{tokenString}, err
}

// Decode decodes a jwt token into a Credentials object
func (j *JWTCoder) Decode(str string) (*models.Credentials, error) {
	token, err := jwt.Parse(str, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token")
		}
		return []byte(j.secret), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		var creds models.Credentials
		mapstructure.Decode(claims, &creds)
		return &creds, nil
	} else {
		return nil, fmt.Errorf("Invalid authorization token")
	}
}
