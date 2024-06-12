package handlers

import (
	"backend/models"
	"backend/repository"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var repo *repository.Repository

func InitHandlers(db *gorm.DB) {
	repo = repository.NewRepository(db)

	models := []interface{}{
		&models.User{},
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
		&models.RoleOnBoard{},
		&models.RoleOnSpace{},
		&models.BoardRoleOnBoard{},
		&models.UserRoleOnSpace{},
		&models.UserBoardRoleOnBoard{},
	}

	for _, model := range models {
		err := db.AutoMigrate(model)
		if err != nil {
			log.Fatalf("AutoMigrate failed for model %T: %v", model, err)
		}
	}
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
	if _, _, err := repo.CreateSpaceWithOwnerRole("test", user.Login, ""); err != nil {
		http.Error(w, "Ошибка создания пространства", http.StatusInternalServerError)
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
