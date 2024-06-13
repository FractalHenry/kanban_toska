package handlers

import (
	"backend/models"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetSpaceRolesHandler(w http.ResponseWriter, r *http.Request) {
	spaceID, err := strconv.Atoi(mux.Vars(r)["spaceId"])
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
	var requestBody struct {
		Name    string `json:"name"`
		IsAdmin bool   `json:"is_admin"`
		CanEdit bool   `json:"can_edit"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	roleOnSpace := models.RoleOnSpace{
		RoleOnSpaceName: requestBody.Name,
		IsAdmin:         requestBody.IsAdmin,
		CanEdit:         requestBody.CanEdit,
		IsOwner:         false,
	}

	userLogin := r.Context().Value("userLogin").(string)
	err = repo.CreateRoleOnSpace(&roleOnSpace, userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func DeleteSpaceRoleHandler(w http.ResponseWriter, r *http.Request) {
	roleOnSpaceID, err := strconv.Atoi(mux.Vars(r)["roleOnSpaceId"])
	if err != nil {
		http.Error(w, "Неверный идентификатор роли пространства", http.StatusBadRequest)
		return
	}

	userLogin := r.Context().Value("userLogin").(string)

	err = repo.DeleteRoleOnSpace(uint(roleOnSpaceID), userLogin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
