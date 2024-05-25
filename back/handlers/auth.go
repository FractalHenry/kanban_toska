package handlers

import (
	"backend/models"
	"backend/utils"
	"encoding/json"
	"net/http"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

// db - глобальная переменная для хранения соединения с базой данных
var db *gorm.DB

// InitDatabase инициализирует соединение с базой данных и выполняет миграцию для модели User
func InitDatabase(database *gorm.DB) {
	db = database
	db.AutoMigrate(&models.User{})
}

// Register обрабатывает регистрацию нового пользователя
func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Декодируем JSON из тела запроса в объект user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}
	// Проверяем длину логина и пароля
	if len(user.Password) < 5 {
		http.Error(w, "Длинна логина должна быть минимум 6 символов", http.StatusBadRequest)
		return
	}
	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
		return
	}
	user.Password = string(hashedPassword)
	// Пытаемся создать нового пользователя в базе данных
	if err := db.Create(&user).Error; err != nil {
		http.Error(w, "Такой пользователь уже сущесвует", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// Login обрабатывает вход пользователя
func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	// Декодируем JSON из тела запроса в объект user
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}
	var dbUser models.User
	// Ищем пользователя в базе данных по логину
	if err := db.Where("login = ?", user.Login).First(&dbUser).Error; err != nil {
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}
	// Сравниваем хешированный пароль с введенным
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}
	// Генерируем JWT для пользователя
	token, err := utils.GenerateJWT(user.Login)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}
	// Возвращаем JWT в ответе
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ProtectedEndpoint - защищенный маршрут, доступный только аутентифицированным пользователям
func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("login")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello " + login})
}

// ProtectedEndpointWithLogin - защищенный маршрут с динамическим логином
func ProtectedEndpointWithLogin(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("login")
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello " + login})
}
