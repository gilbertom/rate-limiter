package middlewares

import (
	"context"
	"net/http"
)

// GetTokenFromHeader gets the Token API_KEY from header.
func GetTokenFromHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct := context.WithValue(r.Context(), "Token", r.Header.Get("API_KEY"))
		r = r.WithContext(ct)
		next.ServeHTTP(w, r)
	})
}
