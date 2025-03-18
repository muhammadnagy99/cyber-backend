package repository

import (
	"encoding/json"
	"errors"
	"cyber-rbac/pkg/logger"
	"github.com/cockroachdb/pebble"
)

type RolePermissionRepository interface {
	AssignPermissionToRole(roleID, permissionID string) error
	RemovePermissionFromRole(roleID, permissionID string) error
	GetPermissionsForRole(roleID string) ([]string, error)
	HasPermission(roleID, permissionID string) (bool, error)
}

type RolePermissionRepositoryImpl struct {
	db *PebbleDB
}

func NewRolePermissionRepository(db *PebbleDB) RolePermissionRepository {
	return &RolePermissionRepositoryImpl{db: db}
}

func (r *RolePermissionRepositoryImpl) AssignPermissionToRole(roleID, permissionID string) error {
	logger.Logger.Infof("Assigning permission '%s' to role '%s'", permissionID, roleID)

	// Check if the permission exists
	data, closer, err := r.db.db.Get([]byte("permission:" + permissionID))
	logger.Logger.Infof("permission '%s''", data)

	if err == pebble.ErrNotFound {
		logger.Logger.Warnf("Permission '%s' does not exist", permissionID)
		return errors.New("permission not found")
	}
	if err != nil {
		logger.Logger.Errorf("Error retrieving permission '%s': %v", permissionID, err)
		return err
	}
	defer closer.Close()

	// Retrieve existing permissions for the role
	permissions, err := r.GetPermissionsForRole(roleID)
	if err != nil && !errors.Is(err, pebble.ErrNotFound) {
		logger.Logger.Errorf("Error retrieving permissions for role '%s': %v", roleID, err)
		return err
	}

	// Avoid duplicate assignments
	for _, perm := range permissions {
		if perm == permissionID {
			logger.Logger.Warnf("Permission '%s' is already assigned to role '%s'", permissionID, roleID)
			return nil
		}
	}

	permissions = append(permissions, permissionID)
	data, err = json.Marshal(permissions)
	if err != nil {
		logger.Logger.Errorf("Failed to encode permissions for role '%s': %v", roleID, err)
		return err
	}

	err = r.db.db.Set([]byte("role_permissions:"+roleID), data, &pebble.WriteOptions{Sync: true})
	if err != nil {
		logger.Logger.Errorf("Failed to store role-permission mapping for role '%s': %v", roleID, err)
		return err
	}

	logger.Logger.Infof("Successfully assigned permission '%s' to role '%s'", permissionID, roleID)
	return nil
}


func (r *RolePermissionRepositoryImpl) RemovePermissionFromRole(roleID, permissionID string) error {
	logger.Logger.Infof("Removing permission '%s' from role '%s'", permissionID, roleID)
	permissions, err := r.GetPermissionsForRole(roleID)
	if err != nil {
		logger.Logger.Errorf("Error retrieving permissions for role '%s': %v", roleID, err)
		return err
	}

	var updatedPermissions []string
	found := false
	for _, perm := range permissions {
		if perm != permissionID {
			updatedPermissions = append(updatedPermissions, perm)
		} else {
			found = true
		}
	}

	if !found {
		logger.Logger.Warnf("Permission '%s' not found in role '%s'", permissionID, roleID)
		return nil
	}

	data, err := json.Marshal(updatedPermissions)
	if err != nil {
		logger.Logger.Errorf("Failed to encode updated permissions for role '%s': %v", roleID, err)
		return err
	}

	err = r.db.db.Set([]byte("role_permissions:"+roleID), data, &pebble.WriteOptions{Sync: true})
	if err != nil {
		logger.Logger.Errorf("Failed to update role-permission mapping for role '%s': %v", roleID, err)
		return err
	}

	logger.Logger.Infof("Successfully removed permission '%s' from role '%s'", permissionID, roleID)
	return nil
}

func (r *RolePermissionRepositoryImpl) GetPermissionsForRole(roleID string) ([]string, error) {
	logger.Logger.Infof("Retrieving permissions for role '%s'", roleID)
	data, closer, err := r.db.db.Get([]byte("role_permissions:" + roleID))
	if err == pebble.ErrNotFound {
		logger.Logger.Warnf("No permissions found for role '%s'", roleID)
		return []string{}, nil
	}
	if err != nil {
		logger.Logger.Errorf("Failed to retrieve role permissions for role '%s': %v", roleID, err)
		return nil, err
	}
	defer closer.Close()

	var permissions []string
	err = json.Unmarshal(data, &permissions)
	if err != nil {
		logger.Logger.Errorf("Failed to decode permissions for role '%s': %v", roleID, err)
		return nil, err
	}

	logger.Logger.Infof("Successfully retrieved permissions for role '%s': %v", roleID, permissions)
	return permissions, nil
}

func (r *RolePermissionRepositoryImpl) HasPermission(roleID, permissionID string) (bool, error) {
	permissions, err := r.GetPermissionsForRole(roleID)
	if err != nil {
		if errors.Is(err, pebble.ErrNotFound) {
			logger.Logger.Warnf("Role '%s' not found when checking permission '%s'", roleID, permissionID)
			return false, nil
		}
		logger.Logger.Errorf("Error retrieving permissions for role '%s': %v", roleID, err)
		return false, err
	}

	// Check if the role has the given permission
	for _, perm := range permissions {
		if perm == permissionID {
			logger.Logger.Infof("Role '%s' has permission '%s'", roleID, permissionID)
			return true, nil
		}
	}

	logger.Logger.Infof("Role '%s' does NOT have permission '%s'", roleID, permissionID)
	return false, nil
}