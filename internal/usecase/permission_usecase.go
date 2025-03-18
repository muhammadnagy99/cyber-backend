package usecase

import (
	"errors"
	"cyber-rbac/internal/domain"
	"cyber-rbac/internal/repository"
)

type PermissionUseCase struct {
	repo repository.PermissionRepository
}

func NewPermissionUseCase(repo repository.PermissionRepository) *PermissionUseCase {
	return &PermissionUseCase{repo: repo}
}

func (uc *PermissionUseCase) CreatePermission(permission domain.Permission) error {
	if permission.ID == "" || permission.Name == "" {
		return errors.New("permission ID and name cannot be empty")
	}
	return uc.repo.Create(permission)
}

func (uc *PermissionUseCase) UpdatePermission(permission domain.Permission) error {
	if permission.ID == "" {
		return errors.New("permission ID is required")
	}
	return uc.repo.Update(permission)
}

func (uc *PermissionUseCase) DeletePermission(permissionID string) error {
	if permissionID == "" {
		return errors.New("permission ID is required")
	}
	return uc.repo.Delete(permissionID)
}

func (uc *PermissionUseCase) GetPermissionByID(permissionID string) (domain.Permission, error) {
	if permissionID == "" {
		return domain.Permission{}, errors.New("permission ID is required")
	}
	return uc.repo.GetByID(permissionID)
}

func (uc *PermissionUseCase) GetAllPermissions() ([]domain.Permission, error) {
	return uc.repo.GetAll()
}
