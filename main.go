/*
 1. сначала нужно зарегистрироваться через http://localhost:8080/api/register и ввести в тело запроса: {
    "username": "user123",
    "password": "password321"
    }

2)после этого делаем то же самое но по маршруту http://localhost:8080/api/login и получаем в ответ ТОКЕН
3)далее при вводе http://localhost:8080/api/address/search или /geocode и тела запрос (например {"query": "Москва"}) нужно
в bearer token вставить наш токен и тогда мы получим нужный ответ от сервера
4) если токен не ввести то получим статус 403
*/
package main

import (
	"net/http"
	"net/http/pprof"

	"github.com/ArthurGopher/AuthorizationJWT/auth"
	"github.com/ArthurGopher/AuthorizationJWT/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/ArthurGopher/AuthorizationJWT/docs"
)

func main() {
	// cоздаем JWT-токен с секретным ключом

	auth.TokenAuth = jwtauth.New("HS256", []byte("secret-key"), nil)
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
	))
	r.Post("/api/register", auth.UsersRegisterHandler)
	r.Post("/api/login", auth.LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(auth.TokenAuth)) // Проверка токена
		r.Use(auth.AuthMiddleware)              // Middleware для авторизации

		r.Post("/api/address/search", handlers.SearchHandler)
		r.Post("/api/address/geocode", handlers.GeocodeHandler)

		// добавляем маршруты для pprof
		r.Route("/mycustompath/pprof", func(r chi.Router) {
			r.HandleFunc("/", pprof.Index)
			r.HandleFunc("/cmdline", pprof.Cmdline)
			r.HandleFunc("/profile", pprof.Profile)
			r.HandleFunc("/symbol", pprof.Symbol)
			r.HandleFunc("/trace", pprof.Trace)
			r.HandleFunc("/allocs", pprof.Handler("allocs").ServeHTTP)
			r.HandleFunc("/block", pprof.Handler("block").ServeHTTP)
			r.HandleFunc("/goroutine", pprof.Handler("goroutine").ServeHTTP)
			r.HandleFunc("/heap", pprof.Handler("heap").ServeHTTP)
			r.HandleFunc("/mutex", pprof.Handler("mutex").ServeHTTP)
			r.HandleFunc("/threadcreate", pprof.Handler("threadcreate").ServeHTTP)
		})
	})

	http.ListenAndServe(":8080", r)
}
