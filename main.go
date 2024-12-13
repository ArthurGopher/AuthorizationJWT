/* Я долго пытался понять почему код не запускается через go run main.go но запускается через go run . 
но так и не понял, пожалуйста напишите в комментарии 

Проверил работу кода в постман, все отлично работает:
1) сначала нужно зарегистрироваться через http://localhost:8080/api/login и ввести в тело запроса: {
    "username": "user123",
    "password": "password321"
}
2)после этого делаем то же самое но по маршруту http://localhost:8080/api/login и получаем в ответ ТОКЕН
3)далее при вводе http://localhost:8080/api/address/search или /geocode и тела запроса нужно 
в bearer token вставить наш токен и тогда мы получим нужный ответ от сервера
4) если токен не ввести то получим статус 403 

*/
package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"

 
   
	"github.com/swaggo/http-swagger" 

	_ "github.com/ArthurGopher/AuthorizationJWT/docs"
)

var tokenAuth *jwtauth.JwtAuth

func main() {
	// cоздаем JWT-токен с секретным ключом
	tokenAuth = jwtauth.New("HS256", []byte("secret-key"), nil)
	r := chi.NewRouter()


	r.Use(middleware.Recoverer)
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), 
	))
	r.Post("/api/register", UsersRegisterHandler)
	r.Post("/api/login", LoginHandler)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth)) // Проверка токена
		r.Use(AuthMiddleware)              // Middleware для авторизации

		r.Post("/api/address/search", SearchHandler)
		r.Post("/api/address/geocode", GeocodeHandler)
	})

	http.ListenAndServe(":8080", r)
}
