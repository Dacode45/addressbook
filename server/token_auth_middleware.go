package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

var CREDENTIALS_KEY = "credentials"

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
		ctx := context.WithValue(r.Context(), CREDENTIALS_KEY, creds)
		r = r.WithContext(ctx)
		next(w, r)
	}
}
