package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func DeleteSpaceHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID пространства из пути запроса
	vars := mux.Vars(r)
	spaceID, err := strconv.ParseUint(vars["spaceId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid space ID", http.StatusBadRequest)
		return
	}

	// Удаляем пространство из базы данных
	err = repo.DeleteSpace(uint(spaceID), userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}

func UpdateSpaceHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Получаем ID пространства из пути запроса
	vars := mux.Vars(r)
	spaceID, err := strconv.ParseUint(vars["spaceId"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid space ID", http.StatusBadRequest)
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

	// Обновляем название пространства
	err = repo.UpdateSpaceName(uint(spaceID), reqBody.Name, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusOK)
}

func CreateSpaceHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем логин пользователя из заголовка
	userLogin := r.Header.Get("login")

	// Декодируем полученные данные
	var reqBody struct {
		SpaceName       string `json:"spaceName"`
		RoleOnSpaceName string `json:"roleOnSpaceName"`
	}
	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Создаем новое пространство с ролью владельца
	space, ownerRole, err := repo.CreateSpaceWithOwnerRole(reqBody.SpaceName, userLogin, reqBody.RoleOnSpaceName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Отправляем ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		SpaceID       uint   `json:"spaceId"`
		SpaceName     string `json:"spaceName"`
		OwnerRoleID   uint   `json:"ownerRoleId"`
		OwnerRoleName string `json:"ownerRoleName"`
	}{
		SpaceID:       space.SpaceID,
		SpaceName:     space.SpaceName,
		OwnerRoleID:   ownerRole.RoleOnSpaceID,
		OwnerRoleName: ownerRole.RoleOnSpaceName,
	})
}

func GetUsersNotOnSpace(w http.ResponseWriter, r *http.Request) {
	// Декодируем полученные данные
	var reqBody struct {
		SpaceID string `json:"spaceID"`
	}

	err := json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
}

func AddUserToSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	spaceID, err := strconv.ParseUint(vars["spaceID"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid space ID", http.StatusBadRequest)
		return
	}

	userLogin := r.Header.Get("login")

	var reqBody struct {
		TargetUserLogin string `json:"targetUserLogin"`
		RoleOnSpaceID   uint   `json:"roleOnSpaceID"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = repo.AssociateUserWithRoleOnSpace(userLogin, reqBody.TargetUserLogin, reqBody.RoleOnSpaceID, uint(spaceID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func RemoveUserFromSpace(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.ParseUint(vars["spaceID"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный space ID", http.StatusBadRequest)
		return
	}

	userLogin := r.Header.Get("login")

	var reqBody struct {
		TargetUserLogin string `json:"targetUserLogin"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Некорректный request body", http.StatusBadRequest)
		return
	}

	err = repo.DeleteUserRoleOnSpace(userLogin, reqBody.TargetUserLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteSpaceRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_, err := strconv.ParseUint(vars["spaceID"], 10, 64)
	if err != nil {
		http.Error(w, "Некорректный space ID", http.StatusBadRequest)
		return
	}

	userLogin := r.Header.Get("login")

	var reqBody struct {
		RoleOnSpaceID uint `json:"roleOnSpaceID"`
	}
	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Некорректный request body", http.StatusBadRequest)
		return
	}

	err = repo.DeleteRoleOnSpace(reqBody.RoleOnSpaceID, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
