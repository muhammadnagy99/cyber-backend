package http

import (
	"encoding/json"
	"net/http"

	"cyber-rbac/internal/domain"
	"cyber-rbac/internal/usecase"

	"github.com/gorilla/mux"
)

type RoleHandler struct {
	uc *usecase.RoleUseCase
}

func NewRoleHandler(uc *usecase.RoleUseCase) *RoleHandler {
	return &RoleHandler{uc: uc}
}

// CreateRole - Handles role creation
func (h *RoleHandler) CreateRole(w http.ResponseWriter, r *http.Request) {
	var role domain.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, `{"error": "Invalid request payload, expected JSON with 'id' and 'name'"}`, http.StatusBadRequest)
		return
	}

	if err := h.uc.CreateRole(role); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Role created successfully",
		"status":  "success",
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetRoleByID - Retrieves a specific role by ID
func (h *RoleHandler) GetRoleByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID := vars["id"]

	role, err := h.uc.GetRoleByID(roleID)
	if err != nil {
		http.Error(w, `{"error": "Role not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(role)
}

// GetAllRoles - Retrieves all roles
func (h *RoleHandler) GetAllRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := h.uc.GetAllRoles()
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(roles)
}

// UpdateRole - Updates a role
func (h *RoleHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	var role domain.Role
	if err := json.NewDecoder(r.Body).Decode(&role); err != nil {
		http.Error(w, `{"error": "Invalid request payload, expected JSON with 'id' and 'name'"}`, http.StatusBadRequest)
		return
	}

	if err := h.uc.UpdateRole(role); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Role updated successfully",
		"status":  "success",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteRole - Deletes a role
func (h *RoleHandler) DeleteRole(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID := vars["id"]

	if err := h.uc.DeleteRole(roleID); err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"message": "Role deleted successfully",
		"status":  "success",
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
