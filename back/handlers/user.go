package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

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
		SpaceName  string `json:"SpaceName"`
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

		var userLogins []string
		for _, user := range *users {
			userLogins = append(userLogins, user.Login)
		}

		spaceList = append(spaceList, struct {
			SpaceName  string `json:"SpaceName"`
			SpaceID    uint   `json:"spaceId"`
			SpaceOwner string `json:"SpaceOwner"`
			Boards     []struct {
				ID   uint   `json:"id"`
				Name string `json:"name"`
			} `json:"boards"`
			Users []string `json:"users"`
		}{
			SpaceName:  space.SpaceName,
			SpaceID:    space.SpaceID,
			SpaceOwner: SpaceOwner,
			Boards:     boardList,
			Users:      userLogins,
		})
	}

	// Отправляем JSON ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spaceList)
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
