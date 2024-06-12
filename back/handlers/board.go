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
	boardUsers, err := repo.GetBoardUsers(uint(boardID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	space, _ := repo.GetSpaceByBoard(uint(boardID))
	spaceUsers, err := repo.GetSpaceUsers(space.SpaceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Создаем ответ в формате JSON
	response := struct {
		ID         uint                        `json:"id"`
		Cards      []models.Card               `json:"cards"`
		Tasks      []models.Task               `json:"tasks"`
		InfoBlocks []models.InformationalBlock `json:"infoBlocks"`
		BoardUsers []struct {
			Login   string `json:"login"`
			CanEdit bool   `json:"can_edit"`
		} `json:"boardUsers"`
		SpaceUsers []struct {
			Login   string `json:"login"`
			IsAdmin bool   `json:"is_admin"`
			IsOwner bool   `json:"is_owner"`
			CanEdit bool   `json:"can_edit"`
		} `json:"spaceUsers"`
	}{
		ID:         board.BoardID,
		Cards:      cards,
		Tasks:      tasks,
		InfoBlocks: infoBlocks,
		BoardUsers: make([]struct {
			Login   string `json:"login"`
			CanEdit bool   `json:"can_edit"`
		}, 0),
		SpaceUsers: make([]struct {
			Login   string `json:"login"`
			IsAdmin bool   `json:"is_admin"`
			IsOwner bool   `json:"is_owner"`
			CanEdit bool   `json:"can_edit"`
		}, 0),
	}

	// Заполняем информацию о пользователях
	if boardUsers != nil {
		for _, user := range *boardUsers {
			response.BoardUsers = append(response.BoardUsers, struct {
				Login   string `json:"login"`
				CanEdit bool   `json:"can_edit"`
			}{
				Login:   user.Login,
				CanEdit: user.RoleOnBoard.CanEdit,
			})
		}
	}

	if spaceUsers != nil {
		for _, user := range *spaceUsers {
			response.SpaceUsers = append(response.SpaceUsers, struct {
				Login   string `json:"login"`
				IsAdmin bool   `json:"is_admin"`
				IsOwner bool   `json:"is_owner"`
				CanEdit bool   `json:"can_edit"`
			}{
				Login:   user.Login,
				IsAdmin: user.RoleOnSpace.IsAdmin,
				IsOwner: user.RoleOnSpace.IsOwner,
				CanEdit: user.RoleOnSpace.CanEdit,
			})
		}
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
