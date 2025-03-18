package http

import (
	"encoding/json"
	"net/http"

	"cyber-rbac/internal/usecase"
	"cyber-rbac/pkg/logger"

	"github.com/gorilla/mux"
)

type RolePermissionHandler struct {
	uc *usecase.RolePermissionUseCase
}

func NewRolePermissionHandler(uc *usecase.RolePermissionUseCase) *RolePermissionHandler {
	return &RolePermissionHandler{uc: uc}
}

func (h *RolePermissionHandler) AssignPermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID := vars["role_id"]
	permissionID := vars["permission_id"]

	if roleID == "" || permissionID == "" {
		http.Error(w, `{"status":"error", "message":"Role ID and Permission ID are required"}`, http.StatusBadRequest)
		return
	}

	if err := h.uc.AssignPermission(roleID, permissionID); err != nil {
		http.Error(w, `{"status":"error", "message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "success", "message": "Permission assigned successfully"}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *RolePermissionHandler) RemovePermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID := vars["role_id"]
	permissionID := vars["permission_id"]

	if roleID == "" || permissionID == "" {
		http.Error(w, `{"status":"error", "message":"Role ID and Permission ID are required"}`, http.StatusBadRequest)
		return
	}

	if err := h.uc.RemovePermission(roleID, permissionID); err != nil {
		http.Error(w, `{"status":"error", "message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]string{"status": "success", "message": "Permission removed successfully"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *RolePermissionHandler) GetPermissions(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID := vars["role_id"]

	if roleID == "" {
		http.Error(w, `{"status":"error", "message":"Role ID is required"}`, http.StatusBadRequest)
		return
	}

	permissions, err := h.uc.GetPermissions(roleID)
	if err != nil {
		http.Error(w, `{"status":"error", "message":"`+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"status": "success",
		"permissions": permissions,
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func (h *RolePermissionHandler) HasPermission(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roleID := vars["role_id"]
	permissionID := vars["permission_id"]

	hasPerm, err := h.uc.HasPermission(roleID, permissionID)
	if err != nil {
		logger.Logger.Errorf("Failed to check permission for role '%s': %v", roleID, err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Return JSON response
	response := map[string]interface{}{
		"role":        roleID,
		"permission":  permissionID,
		"has_access":  hasPerm,
		"status_code": http.StatusOK,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}