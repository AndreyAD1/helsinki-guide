package middlewares

import (
	"net/http"
)

func GetBasicAuthHandler(
	next http.Handler,
	expectedUsername,
	expectedPassword string,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if !ok || user != expectedUsername || password != expectedPassword {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
