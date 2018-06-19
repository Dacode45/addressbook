package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Dacode45/addressbook/models"
	"github.com/Dacode45/addressbook/storage"
)

// TokenAuthMiddleware is simple middleware to parse JWT from the authorization header
// Expects the authorization: Bearer <token> format, but Bearer isn't required
func (coder *JWTCoder) TokenAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")
		if authorizationHeader == "" {
			StatusUnauthorized.Serve(fmt.Errorf("route requires jwt authorization"))(w, r)
			return
		}
		bearerToken := strings.Split(authorizationHeader, " ")
		if len(bearerToken) != 2 {
			StatusBadRequest.Serve(fmt.Errorf("invalid authorization header"))(w, r)
			return
		}
		creds, err := coder.Decode(bearerToken[1])
		if err != nil {
			StatusInternalServerError.Serve(err)(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), ContextCredentialsKey, creds)
		r = r.WithContext(ctx)
		next(w, r)
	}
}

// LoggedInMiddleware logs the user in based of their jwt token
func LoggedInMiddleware(jwtCoder *JWTCoder, userStorage storage.UserStorage, next http.HandlerFunc) http.HandlerFunc {
	return jwtCoder.TokenAuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		creds, ok := ctx.Value(ContextCredentialsKey).(*models.Credentials)
		if !ok || creds == nil {
			StatusUnauthorized.Serve(fmt.Errorf("no jwt passed"))(w, r)
			return
		}
		user, err := userStorage.Login(r.Context(), *creds)
		if err != nil {
			StatusUnauthorized.Serve(err)(w, r)
			return
		}
		ctx = context.WithValue(ctx, ContextUserKey, user)
		r = r.WithContext(ctx)
		next(w, r)
	})
}
