package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string
	Password string
}

var UserStorage = struct{
	sync.Mutex 
	users map[string]string
}{users : make(map[string]string)}

// UsersRegisterHandler обрабатывает регистрацию пользователей
// @Summary Регистрация нового пользователя
// @Description Регистрирует нового пользователя, сохраняя в памяти его имя пользователя и хэшированный пароль
// @Tags authorization
// @Accept json
// @Produce json
// @Param user body User true "Учетные данные пользователя"
// @Success 201 {string} string "Пользователь успешно зарегистрирован"
// @Failure 400 {string} string "Неверный формат запроса"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /api/register [post]
func UsersRegisterHandler(w http.ResponseWriter, r *http.Request){
	var user User 
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	//хэшируем пароль перед сохранением 
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "ошибка при хэшировании пароля", http.StatusInternalServerError)
		return
	}
	UserStorage.Lock()
	defer UserStorage.Unlock()
	
	if _, exists := UserStorage.users[user.Username]; exists {
		http.Error(w, "Пользователь уже существует", http.StatusBadRequest)
		return
	}

	UserStorage.users[user.Username] = string(hashedPassword)
	w.WriteHeader(http.StatusCreated)
}

// LoginHandler обрабатывает вход пользователя в систему
// @Summary Вход в систему пользователя
// @Description Выполняет вход пользователя и возвращает JWT-токен, если учетные данные верны
// @Tags authorization
// @Accept json
// @Produce json
// @Param user body User true "Учетные данные пользователя"
// @Success 200 {object} map[string]string "JWT-токен"
// @Failure 400 {string} string "Неверный формат запроса"
// @Success 200 {string} string "Пользователь не существует или пароль не совпадает"
// @Router /api/login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request){
	var user User 
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil{
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	UserStorage.Lock()
	defer UserStorage.Unlock()

	storedPassword, exists := UserStorage.users[user.Username]
	if !exists {
		http.Error(w, "пользователь не существует или пароль не совпадает", http.StatusOK)
		return
	}

	// cравниваем хэш пароля
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		http.Error(w, "пользователь не существует или пароль не совпадает", http.StatusOK)
		return
	}

	// генерируем JWT-токен
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"username": user.Username})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})


}