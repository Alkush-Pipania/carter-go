package middleware

import (
	"net/http"
	"strings"
)

type Jwtservice interface {
	Verify(token string) (string, error)
}

func AuthMiddleware(jwt Jwtservice) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if auth == "" {
				http.Error(w, "missing Authorized Header", http.StatusUnauthorized)
				return
			}
			parts := strings.Split(auth, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "invalid Authorization Header", http.StatusUnauthorized)
				return
			}
			user_id, err := jwt.Verify(parts[1])
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			r.Header.Set("user_id", user_id)
			next.ServeHTTP(w, r)
		})
	}
}
