package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetBoardDetailsHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем ID доски из пути запроса
	vars := mux.Vars(r)
	boardID, err := strconv.ParseUint(vars["boardId"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный ID доски", http.StatusBadRequest)
		return
	}

	// Получаем детальную информацию о доске
	board, cards, tasks, infoBlocks, err := repo.GetBoardDetails(uint(boardID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Получаем пользователей и их роли на доске
	users, err := repo.GetBoardUsers(uint(boardID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем ответ в формате JSON
	response := struct {
		ID         uint                        `json:"id"`
		Cards      []models.Card               `json:"cards"`
		Tasks      []models.Task               `json:"tasks"`
		InfoBlocks []models.InformationalBlock `json:"infoBlocks,omitempty"`
		Users      []struct {
			Login       string `json:"login"`
			RoleOnBoard string `json:"role_on_board"`
		} `json:"users"`
	}{
		ID:         board.BoardID,
		Cards:      cards,
		Tasks:      tasks,
		InfoBlocks: infoBlocks,
		Users:      nil, // Инициализируем пустой слайс
	}

	// Заполняем информацию о пользователях
	for _, user := range users {
		response.Users = append(response.Users, struct {
			Login       string `json:"login"`
			RoleOnBoard string `json:"role_on_board"`
		}{
			Login:       user.Login,
			RoleOnBoard: user.RoleOnBoard.RoleOnBoardName,
		})
	}

	// Отправляем ответ в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func CreateBoardHandler(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Header.Get("login")

	var boardData struct {
		SpaceID   uint   `json:"spaceid"`
		BoardName string `json:"boardname"`
	}

	err := json.NewDecoder(r.Body).Decode(&boardData)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	board := &models.Board{
		SpaceID:   boardData.SpaceID,
		BoardName: boardData.BoardName,
	}

	err = repo.CreateBoard(board, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
