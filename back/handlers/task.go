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

func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID задания из пути запроса
	vars := mux.Vars(r)
	taskID, err := strconv.ParseUint(vars["TaskID"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный task ID", http.StatusBadRequest)
		return
	}

	// Декодируем полученные данные
	var reqBody struct {
		Name        string `json:"name"`
		Color       string `json:"color"`
		Description string `json:"Dscription"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Некорректный request body", http.StatusBadRequest)
		return
	}

	// Получаем задание из базы данных
	task, err := repo.GetTaskByID(uint(taskID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Обновляем поля задания
	task.TaskName = reqBody.Name

	// Обновляем цвет задания, если он был передан
	if reqBody.Color != "" {
		taskColor, _ := repo.GetTaskColor(task.TaskID)
		TaskColor := &models.TaskColor{
			TaskColor: reqBody.Color,
			TaskID:    task.TaskID,
		}
		if taskColor == "" {
			err = repo.CreateTaskColor(TaskColor, userLogin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = repo.UpdateTaskColor(TaskColor, userLogin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	}

	// Обновляем описание задания, если оно было передано
	if reqBody.Description != "" {
		taskDescription, _ := repo.GetTaskDescription(task.TaskID)
		TaskDescription := &models.TaskDescription{
			TaskDescription: reqBody.Description,
			TaskID:          task.TaskID,
		}
		if taskDescription == "" {
			err = repo.CreateTaskDescription(TaskDescription, userLogin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			err = repo.UpdateTaskDescription(TaskDescription, userLogin)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	} else {
		// Удаляем описание задания, если оно не было передано
		err = repo.DeleteTaskDescription(task.TaskID, userLogin)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Обновляем задание в базе данных
	err = repo.UpdateTask(&task, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}

func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID задания из пути запроса
	vars := mux.Vars(r)
	taskID, err := strconv.ParseUint(vars["TaskID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid task ID", http.StatusBadRequest)
		return
	}

	// Удаляем задание из базы данных
	err = repo.DeleteTask(uint(taskID), userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}
