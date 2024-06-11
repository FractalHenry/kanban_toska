package handlers

import (
	"backend/models"
	"backend/repository"
	"backend/utils"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

var repo *repository.Repository

func InitHandlers(db *gorm.DB) {
	repo = repository.NewRepository(db)

	// Автоматическое создание таблиц
	db.AutoMigrate(
		&models.User{},
		&models.RoleOnBoard{},
		&models.BoardRoleOnBoard{},
		&models.RoleOnSpace{},
		&models.UserRoleOnSpace{},
		&models.UserBoardRoleOnBoard{},
		&models.Space{},
		&models.Board{},
		&models.Card{},
		&models.Task{},
		&models.Mark{},
		&models.MarkName{},
		&models.Checklist{},
		&models.ChecklistElement{},
		&models.TaskColor{},
		&models.TaskDescription{},
		&models.TaskDateStart{},
		&models.TaskDateEnd{},
		&models.TaskNotification{},
		&models.Notification{},
	)
}

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}
	if len(user.Password) < 5 {
		http.Error(w, "Длинна логина должна быть минимум 6 символов", http.StatusBadRequest)
		return
	}
	if err := repo.CreateUser(&user); err != nil {
		http.Error(w, "Такой пользователь уже сущесвует", http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Неверный ввод", http.StatusBadRequest)
		return
	}
	dbUser, err := repo.FindUserByLogin(user.Login)
	if err != nil {
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Неправильный логин или пароль", http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateJWT(user.Login)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// ProtectedEndpoint - защищенный маршрут, доступный только аутентифицированным пользователям
func ProtectedEndpoint(w http.ResponseWriter, r *http.Request) {
	login := r.URL.Path
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello " + login})
}

func ProtectedEndpointWithLogin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]
	json.NewEncoder(w).Encode(map[string]string{"message": "Hello " + name})
}
