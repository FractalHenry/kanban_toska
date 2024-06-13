package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateMarkHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID задачи из пути запроса
	vars := mux.Vars(r)
	taskID, err := strconv.ParseUint(vars["taskid"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный task ID", http.StatusBadRequest)
		return
	}

	// Декодируем полученные данные
	var reqBody struct {
		Name  string `json:"name"`
		Color string `json:"color"`
	}

	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Некорректный request body", http.StatusBadRequest)
		return
	}

	// Создаем новую метку
	mark := &models.Mark{
		TaskID:    uint(taskID),
		MarkColor: reqBody.Color,
	}

	err = repo.CreateMark(mark, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем имя метки
	markName := &models.MarkName{
		Mark:     *mark,
		MarkID:   mark.MarkID,
		MarkName: reqBody.Name,
	}

	err = repo.CreateMarkName(markName, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
}
