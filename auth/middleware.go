package auth

import (
	"net/http"

	"log"
	"github.com/go-chi/jwtauth"
)

func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, claims, err := jwtauth.FromContext(r.Context())
        if err != nil || claims == nil {
            log.Println("Ошибка токена:", err)
            http.Error(w, "доступ запрещен", http.StatusForbidden)
            return
        }
        log.Println("Успешная авторизация для пользователя:", claims["username"])
        next.ServeHTTP(w, r)
    })
}