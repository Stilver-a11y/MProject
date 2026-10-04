package handlers

import (
	"encoding/json"
	"mproject/Models"
	"mproject/database"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var JwtSecret = []byte("change-me-to-a-long-random-secret")

type RegisterRequest struct {
	Username	string	`json:"username"`
	Email	string	`json:"email"`
	Password	string	`json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		//"Метод не поддерживается"
		return
	}

	var req RegisterRequest
	if err := json.NewDecoder(r.Body). Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		// "Неверное тело запроса"
		return
	}
	
// VALIDATION
// ВАЛИДАЦИЯ
if req.Username == "" || req.Email == "" || len(req.Password) < 6 {
	http.Error(w, "All fields are requered, password must be at least 6 characters", http.StatusBadRequest)
	// "Все поля обязательны для заполнения, пароль должен быть не менее чем из 6 символов"
	return
}

// CHECK IF user ALREADY EXIST
// ПРОВЕРКА, СУЩЕСТВУЕТ ЛИ УЖЕ ПОЛЬЗОВАТЕЛЬ

var existing Models.User
err := database.Db.QueryRow(
	r.Context(),
	"SELECT id FROM users WHERE email = $1 OR username = $2",
	// "Выбрать id из таблицы users (пользователи), где email равен $1 ИЛИ username равно $2"
	req.Email, req.Username,
).Scan(&existing.ID)
if err == nil {
	http.Error(w, "User with this email or username already exist", http.StatusConflict)
	// "Пользователь с таким эмеилом уже существует"
	return
}

// HASH PASSWORD
// ХЭШИРОВАНИЕ ПАРОЛЯ

// Hash password
hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
if err != nil {
    http.Error(w, "Failed to process password", http.StatusInternalServerError)
    return
}

// Inset user
_, err = database.Db.Exec(
	r.Context(),
	"INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)",
	req.Username, req.Email, string(hash), 
)

if err != nil {
	http.Error(w, "Failed to create user:"+err.Error(), http.StatusInternalServerError)
// ОШИБКА СОЗДАНИЯ ПОЛЬЗОВАТЕЛЯ 
	return
}

w.WriteHeader(http.StatusCreated)
json.NewEncoder(w).Encode(map[string]string{
	"massage": "User registered successfuly",
})
}

type LoginRequest struct {
	Email	string	`json:"email"`
	Password	string	`json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		// МЕТОД НЕ ПОДДЕРЖИВАЕТСЯ
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body). Decode (&req); err != nil {
		http.Error (w, "Invalid request body", http.StatusBadRequest)
		// НЕВЕРНОЕ ТЕЛО ЗАПРОСА
		return
	}
	
	if req.Email == "" || req.Password == "" {
		http.Error (w, "Email and Password required", http.StatusBadRequest)
		//
		return
	}

	var user Models.User
	err := database.Db.QueryRow(
		r.Context(),
		"SELECT id, email, password_hash FROM users WHERE email = $1",
		// "Выбрать id из таблицы users (пользователи), где email равен $1 ИЛИ username равно $2"
		req.Email,	
	).Scan(&user.ID, &user.Email, &user.PasswordHash)
	if err != nil {
	http.Error (w, "Invalid email and password", http.StatusUnauthorized)
	return
}

if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
	http.Error(w, "invalid email and password", http.StatusUnauthorized)
	return
}

token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
	"user.id": user.ID,
	"exp": time.Now().Add(72 * time.Hour).Unix(),
})

signed, err := token.SignedString(JwtSecret)
if err != nil {
	http.Error (w, "Failed to create token", http.StatusInternalServerError)
	return
}

w.Header().Set("Content-Type", "application/json")
json.NewEncoder(w).Encode(map[string]string{
	"token": signed,
})
}