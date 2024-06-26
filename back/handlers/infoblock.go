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

func DeleteInfoblockHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID из пути запроса
	vars := mux.Vars(r)
	InfoBlockID, err := strconv.ParseUint(vars["InfoBlockID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid InfoBlock ID", http.StatusBadRequest)
		return
	}

	err = repo.DeleteInformationalBlockByID(uint(InfoBlockID), userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}

func UpdateInfoblockHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	vars := mux.Vars(r)
	InfoBlockID, err := strconv.ParseUint(vars["InfoBlockID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid InfoBlock ID", http.StatusBadRequest)
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

	var InfoBlock *models.InformationalBlock
	InfoBlock, err = repo.GetInformationalBlockByID(uint(InfoBlockID))
	if err != nil {
		http.Error(w, "Invalid InfoBlock ID", http.StatusBadRequest)
		return
	}
	InfoBlock.Header = reqBody.Header
	InfoBlock.Body = reqBody.Body

	err = repo.UpdateInformationalBlocks(InfoBlock, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}
