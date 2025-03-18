package usecase

import (
	"cyber-rbac/internal/repository"
	"cyber-rbac/pkg/logger"
)

type RolePermissionUseCase struct {
	repo repository.RolePermissionRepository
}

func NewRolePermissionUseCase(repo repository.RolePermissionRepository) *RolePermissionUseCase {
	return &RolePermissionUseCase{repo: repo}
}

func (uc *RolePermissionUseCase) AssignPermission(roleID, permissionID string) error {
	return uc.repo.AssignPermissionToRole(roleID, permissionID)
}

func (uc *RolePermissionUseCase) RemovePermission(roleID, permissionID string) error {
	return uc.repo.RemovePermissionFromRole(roleID, permissionID)
}

func (uc *RolePermissionUseCase) GetPermissions(roleID string) ([]string, error) {
	return uc.repo.GetPermissionsForRole(roleID)
}

func (uc *RolePermissionUseCase) HasPermission(roleID, permissionID string) (bool, error) {
	hasPerm, err := uc.repo.HasPermission(roleID, permissionID)
	if err != nil {
		logger.Logger.Errorf("Error checking if role '%s' has permission '%s': %v", roleID, permissionID, err)
		return false, err
	}
	return hasPerm, nil
}