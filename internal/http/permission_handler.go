package http

import (
	"encoding/json"
	"net/http"

	"cyber-rbac/internal/domain"
	"cyber-rbac/internal/usecase"

	"github.com/gorilla/mux"
)

type PermissionHandler struct {
	uc *usecase.PermissionUseCase
}

func NewPermissionHandler(uc *usecase.PermissionUseCase) *PermissionHandler {
	return &PermissionHandler{uc: uc}
}

func (h *PermissionHandler) CreatePermission(w http.ResponseWriter, r *http.Request) {
	var permission domain.Permission
	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		http.Error(w, `{"status":"error", "message":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if err := h.uc.CreatePermission(permission); err != nil {
		http.Error(w, `{"status":"error", "message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "success", "message": "Permission created successfully"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *PermissionHandler) GetPermissionByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	permissionID := vars["id"]

	permission, err := h.uc.GetPermissionByID(permissionID)
	if err != nil {
		http.Error(w, `{"status":"error", "message":"Permission not found"}`, http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permission)
}

func (h *PermissionHandler) GetAllPermissions(w http.ResponseWriter, r *http.Request) {
	permissions, err := h.uc.GetAllPermissions()
	if err != nil {
		http.Error(w, `{"status":"error", "message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(permissions)
}

func (h *PermissionHandler) UpdatePermission(w http.ResponseWriter, r *http.Request) {
	var permission domain.Permission
	if err := json.NewDecoder(r.Body).Decode(&permission); err != nil {
		http.Error(w, `{"status":"error", "message":"Invalid request payload"}`, http.StatusBadRequest)
		return
	}

	if err := h.uc.UpdatePermission(permission); err != nil {
		http.Error(w, `{"status":"error", "message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "success", "message": "Permission updated successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *PermissionHandler) DeletePermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	permissionID := vars["id"]

	if err := h.uc.DeletePermission(permissionID); err != nil {
		http.Error(w, `{"status":"error", "message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "success", "message": "Permission deleted successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}