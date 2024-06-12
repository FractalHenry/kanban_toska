package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
