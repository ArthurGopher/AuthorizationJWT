package main

import (
	"net/http"

	"github.com/go-chi/jwtauth"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, claims, err := jwtauth.FromContext(r.Context())
		if err != nil || claims == nil {
			http.Error(w, "доступ запрещен", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
