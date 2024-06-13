package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateCheckListHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID доски из пути запроса
	vars := mux.Vars(r)
	taskID, err := strconv.ParseUint(vars["TaskID"], 10, 64)
	if err != nil {
		http.Error(w, "Неправильный board ID", http.StatusBadRequest)
		return
	}

	// Декодируем полученные данные
	var reqBody struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Неправильный request body", http.StatusBadRequest)
		return
	}

	// Создаем новую карточку
	checklist := &models.Checklist{
		ChecklistName: reqBody.Name,
		TaskID:        uint(taskID),
	}
	err = repo.CreateChecklist(checklist, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
}

func DeleteCheckListHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID задания из пути запроса
	vars := mux.Vars(r)
	checklistId, err := strconv.ParseUint(vars["checklistID"], 10, 64)
	if err != nil {
		http.Error(w, "Неправильный CheckList ID", http.StatusBadRequest)
		return
	}

	// Удаляем задание из базы данных
	err = repo.DeleteChecklist(uint(checklistId), userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}

func CreateCheckListElementHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID доски из пути запроса
	vars := mux.Vars(r)
	checklistID, err := strconv.ParseUint(vars["CheckListID"], 10, 64)
	if err != nil {
		http.Error(w, "Неправильный board ID", http.StatusBadRequest)
		return
	}

	// Декодируем полученные данные
	var reqBody struct {
		Name string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Неправильный request body", http.StatusBadRequest)
		return
	}

	// Создаем новый элемент чек-листа
	checklistElement := &models.ChecklistElement{
		ChecklistElementName: reqBody.Name,
		ChecklistID:          uint(checklistID),
		ElementOrder:         0,
	}

	err = repo.CreateChecklistElement(checklistElement, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
}
