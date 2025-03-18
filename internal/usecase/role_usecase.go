package usecase

import (
	"errors"
	"cyber-rbac/internal/domain"
	"cyber-rbac/internal/repository"
)

type RoleUseCase struct {
	repo repository.RoleRepository
}

func NewRoleUseCase(repo repository.RoleRepository) *RoleUseCase {
	return &RoleUseCase{repo: repo}
}

func (uc *RoleUseCase) CreateRole(role domain.Role) error {
	if role.ID == "" || role.Name == "" {
		return errors.New("role ID and name cannot be empty")
	}
	return uc.repo.Create(role)
}

func (uc *RoleUseCase) UpdateRole(role domain.Role) error {
	if role.ID == "" {
		return errors.New("role ID is required")
	}
	return uc.repo.Update(role)
}

func (uc *RoleUseCase) DeleteRole(roleID string) error {
	if roleID == "" {
		return errors.New("role ID is required")
	}
	return uc.repo.Delete(roleID)
}

func (uc *RoleUseCase) GetRoleByID(roleID string) (domain.Role, error) {
	if roleID == "" {
		return domain.Role{}, errors.New("role ID is required")
	}
	return uc.repo.GetByID(roleID)
}

func (uc *RoleUseCase) GetAllRoles() ([]domain.Role, error) {
	return uc.repo.GetAll()
}
