package auth

import (
	"encoding/json"
	"net/http"
	"sync"
	"golang.org/x/crypto/bcrypt"
	"github.com/go-chi/jwtauth"
	
		"log"
	)

	var TokenAuth *jwtauth.JwtAuth

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
func UsersRegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	log.Printf("Регистрация пользователя: %s", user.Username)

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Ошибка хэширования пароля: %v", err)
		http.Error(w, "ошибка при хэшировании пароля", http.StatusInternalServerError)
		return
	}

	UserStorage.Lock()
	defer UserStorage.Unlock()

	if _, exists := UserStorage.users[user.Username]; exists {
		log.Printf("Пользователь %s уже существует", user.Username)
		http.Error(w, "пользователь уже существует", http.StatusBadRequest)
		return
	}

	UserStorage.users[user.Username] = string(hashedPassword)
	log.Printf("Пользователь %s успешно зарегистрирован", user.Username)
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
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Printf("Ошибка декодирования JSON: %v", err)
		http.Error(w, "неверный формат запроса", http.StatusBadRequest)
		return
	}

	log.Printf("Попытка входа пользователя: %s", user.Username)

	UserStorage.Lock()
	defer UserStorage.Unlock()

	storedPassword, exists := UserStorage.users[user.Username]
	if !exists {
		log.Printf("Пользователь %s не найден", user.Username)
		http.Error(w, "пользователь не существует или пароль не совпадает", http.StatusUnauthorized)
		return
	}

	// Сравниваем хэш пароля
	if err := bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(user.Password)); err != nil {
		log.Printf("Пароль не совпадает для пользователя %s: %v", user.Username, err)
		http.Error(w, "пользователь не существует или пароль не совпадает", http.StatusUnauthorized)
		return
	}

	// Генерируем JWT-токен
	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{"username": user.Username})
	if err != nil {
		log.Printf("Ошибка генерации JWT токена для пользователя %s: %v", user.Username, err)
		http.Error(w, "внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	log.Printf("Успешный вход пользователя: %s", user.Username)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
