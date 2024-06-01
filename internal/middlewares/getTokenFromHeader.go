package middlewares

import (
	"context"
	"net/http"
)

// GetTokenFromHeader gets the Token API_KEY from header.
func GetTokenFromHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("API_KEY")
		if token != "" {
			ct := context.WithValue(r.Context(), "Token", token)
			r = r.WithContext(ct)
			next.ServeHTTP(w, r)
			return
		}
		ct := context.WithValue(r.Context(), "Token", "token not present")
		r = r.WithContext(ct)
		next.ServeHTTP(w, r)
	})
}
