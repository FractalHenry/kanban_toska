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

func GetUserInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	login := vars["login"]

	// Получаем пользователя из базы данных по логину
	user, err := repo.FindUserByLogin(login)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	// Создаем JSON ответ с информацией о пользователе
	userInfo := struct {
		Login           string `json:"login"`
		Email           string `json:"email"`
		EmailVisibility bool   `json:"emailVisibility"`
		UserDescription string `json:"userDescription"`
	}{
		Login:           user.Login,
		Email:           user.Email,
		EmailVisibility: user.EmailVisibility,
		UserDescription: user.UserDescription,
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userInfo)
}

func GetUserBoards(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	login := r.Header.Get("login")

	// Получаем доски из базы данных для данного пользователя
	boards, err := repo.GetUserBoards(login)
	if err != nil {
		http.Error(w, "Ошибка получения досок", http.StatusInternalServerError)
		return
	}

	// Создаем слайс для хранения неархивированных досок
	var activeBoards []struct {
		ID         uint   `json:"id"`
		Name       string `json:"name"`
		SpaceOwner string `json:"spaceOwner"`
	}

	for _, board := range boards {
		if !board.BoardInArchive {
			space, err := repo.GetSpaceByBoard(board.BoardID)
			if err != nil {
				http.Error(w, "Ошибка получения информации о пространстве", http.StatusInternalServerError)
				return
			}

			spaceOwner, err := repo.GetSpaceOwner(space.SpaceID)
			if err != nil {
				http.Error(w, "Ошибка получения информации о создателе пространства", http.StatusInternalServerError)
				return
			}

			activeBoards = append(activeBoards, struct {
				ID         uint   `json:"id"`
				Name       string `json:"name"`
				SpaceOwner string `json:"spaceOwner"`
			}{
				ID:         board.BoardID,
				Name:       board.BoardName,
				SpaceOwner: spaceOwner,
			})
		}
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(activeBoards)
}

func GetUserSpaces(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	login := r.Header.Get("login")

	// Получаем пространства из базы данных для данного пользователя
	spaces, err := repo.GetUserSpaces(login)
	if err != nil {
		http.Error(w, "Ошибка получения пространств", http.StatusInternalServerError)
		return
	}

	// Создаем слайс для хранения пространств с дополнительной информацией
	var spaceList []struct {
		SpaceOwner string `json:"SpaceOwner"`
		Boards     []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		} `json:"boards"`
		Users []string `json:"users"`
	}

	for _, space := range spaces {
		SpaceOwner, err := repo.GetSpaceOwner(space.SpaceID)
		if err != nil {
			http.Error(w, "Ошибка получения информации о создателе пространства", http.StatusInternalServerError)
			return
		}

		boards, err := repo.GetSpaceBoards(space.SpaceID)
		if err != nil {
			http.Error(w, "Ошибка получения досок для пространства", http.StatusInternalServerError)
			return
		}

		var boardList []struct {
			ID   uint   `json:"id"`
			Name string `json:"name"`
		}
		for _, board := range boards {
			boardList = append(boardList, struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			}{
				ID:   board.BoardID,
				Name: board.BoardName,
			})
		}

		users, err := repo.GetSpaceUsers(space.SpaceID)
		if err != nil {
			http.Error(w, "Ошибка получения пользователей для пространства", http.StatusInternalServerError)
			return
		}

		spaceList = append(spaceList, struct {
			SpaceOwner string `json:"SpaceOwner"`
			Boards     []struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			} `json:"boards"`
			Users []string `json:"users"`
		}{
			SpaceOwner: SpaceOwner,
			Boards:     boardList,
			Users:      users,
		})
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spaceList)
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
