package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetSpaceRolesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	spaceID, err := strconv.ParseUint(vars["spaceId"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор пространства", http.StatusBadRequest)
		return
	}

	roles, err := repo.GetSpaceRoles(uint(spaceID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(roles)
}

func CreateSpaceRoleHandler(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Header.Get("login")

	vars := mux.Vars(r)
	spaceId, err := strconv.ParseUint(vars["spaceId"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор пространства", http.StatusBadRequest)
		return
	}

	var reqBody struct {
		Name    string `json:"name"`
		IsAdmin bool   `json:"is_admin"`
		CanEdit bool   `json:"can_edit"`
	}

	err = json.NewDecoder(r.Body).Decode(&reqBody)
	if err != nil {
		http.Error(w, "Некорректные данные", http.StatusBadRequest)
		return
	}

	roleOnSpace := models.RoleOnSpace{
		SpaceID:         uint(spaceId),
		RoleOnSpaceName: reqBody.Name,
		IsAdmin:         reqBody.IsAdmin,
		CanEdit:         reqBody.CanEdit,
		IsOwner:         false,
	}

	err = repo.CreateRoleOnSpace(&roleOnSpace, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteSpaceRoleHandler(w http.ResponseWriter, r *http.Request) {
	userLogin := r.Header.Get("login")

	vars := mux.Vars(r)
	roleOnSpaceId, err := strconv.ParseUint(vars["roleOnSpaceId"], 10, 64)
	if err != nil {
		http.Error(w, "Неверный идентификатор роли пространства", http.StatusBadRequest)
		return
	}

	err = repo.DeleteRoleOnSpace(uint(roleOnSpaceId), userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
