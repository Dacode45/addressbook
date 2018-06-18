package server

import (
	"fmt"

	"github.com/Dacode45/addressbook/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/mitchellh/mapstructure"
)

type JWTToken struct {
	Token string `json:"token"`
}

type JWTCoder struct {
	secret string
}

func NewJWTCoder(secret string) *JWTCoder {
	return &JWTCoder{
		secret,
	}
}

func (j *JWTCoder) Create(c models.Credentials) (JWTToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": c.Username,
		"password": c.Password,
	})
	tokenString, err := token.SignedString([]byte(j.secret))
	return JWTToken{tokenString}, err
}

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
