package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
