package handlers

import (
	"backend/models"
	"backend/repository"
	"backend/utils"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

func CreateInfoblockHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID доски из пути запроса
	vars := mux.Vars(r)
	boardID, err := strconv.ParseUint(vars["boardId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	// Декодируем полученные данные
	var reqBody struct {
		Header string `json:"header"`
		Body   string `json:"body"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем новый инфоблок
	InformationalBlock := &models.InformationalBlock{
		Header:  reqBody.Header,
		Body:    reqBody.Body,
		BoardID: uint(boardID),
	}
	err = repo.CreateInformationalBlock(InformationalBlock, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
}

func CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID доски и карточки из пути запроса
	vars := mux.Vars(r)

	cardID, err := strconv.ParseUint(vars["cardId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid card ID", http.StatusBadRequest)
		return
	}

	// Декодируем полученные данные
	var reqBody struct {
		Name  string `json:"name"`
		Color string `json:"color"` // Цвет не является обязательным
	}

	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем новый таск
	task := &models.Task{
		TaskName: reqBody.Name,
		CardID:   uint(cardID),
	}

	err = repo.CreateTask(task, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Если был указан цвет, создаем запись о цвете таска
	if reqBody.Color != "" {
		taskColor := &models.TaskColor{
			TaskColor: reqBody.Color,
			TaskID:    task.TaskID,
		}
		err = repo.CreateTaskColor(taskColor, userLogin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
}

func CreateCardHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID доски из пути запроса
	vars := mux.Vars(r)
	boardID, err := strconv.ParseUint(vars["boardId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid board ID", http.StatusBadRequest)
		return
	}

	// Декодируем полученные данные
	var reqBody struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем новую карточку
	card := &models.Card{
		CardName: reqBody.Name,
		BoardID:  uint(boardID),
	}
	err = repo.CreateCard(card, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
}

func UpdateUserDescription(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Декодируем полученные данные
	var reqBody struct {
		NewDescription string `json:"newDescription"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Обновляем описание пользователя
	err = repo.UpdateUserDescription(userLogin, reqBody.NewDescription)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}

func CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Header.Get("login")

	var boardData struct {
		BoardName string `json:"boardname"`
	}

	err := json.NewDecoder(r.Body).Decode(&boardData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	board := &models.Board{
		SpaceID:   1,
		BoardName: boardData.BoardName,
	}

	err = repo.CreateBoard(board, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GetUserLogin(w http.ResponseWriter, r *http.Request) {
	login := r.Header.Get("login")

	user, err := repo.FindUserByLogin(login)
	if err != nil {
		http.Error(w, "Пользователь не найден", http.StatusNotFound)
		return
	}

	response := struct {
		Login string `json:"login"`
	}{
		Login: user.Login,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		SpaceID    uint   `json:"spaceId"`
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
			SpaceID    uint   `json:"spaceId"`
			SpaceOwner string `json:"SpaceOwner"`
			Boards     []struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			} `json:"boards"`
			Users []string `json:"users"`
		}{
			SpaceID:    space.SpaceID,
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
