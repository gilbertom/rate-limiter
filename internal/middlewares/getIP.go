package middlewares

import (
	"context"
	"net"
	"net/http"
)

// GetIPfromClient gets the IP address from the client.
func GetIPfromClient(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		remoteAddr := r.RemoteAddr
		ip, _, err := net.SplitHostPort(remoteAddr)
		if err != nil {
			http.Error(w, "Invalid remote address", http.StatusInternalServerError)
			return
		}
		
		ct := context.WithValue(r.Context(), "IP", ip)
		r = r.WithContext(ct)

		next.ServeHTTP(w, r)
	})
}
