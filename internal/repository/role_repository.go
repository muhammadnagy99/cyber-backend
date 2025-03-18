package repository

import (
	"encoding/json"
	"errors"
	"strings"

	"cyber-rbac/internal/domain"
	"cyber-rbac/pkg/logger"

	"github.com/cockroachdb/pebble"
)

type RoleRepository interface {
	Create(role domain.Role) error
	Update(role domain.Role) error
	Delete(roleID string) error
	GetByID(roleID string) (domain.Role, error)
	GetAll() ([]domain.Role, error)
}

type RoleRepositoryImpl struct {
	db *PebbleDB
}

func NewRoleRepository(db *PebbleDB) RoleRepository {
	return &RoleRepositoryImpl{db: db}
}

func (r *RoleRepositoryImpl) Create(role domain.Role) error {
	if role.ID == "" || role.Name == "" {
		logger.Logger.Errorf("Invalid role data: %+v", role)
		return errors.New("invalid role data: ID and Name are required")
	}

	data, err := json.Marshal(role)
	if err != nil {
		logger.Logger.Errorf("Failed to encode role to JSON: %v", err)
		return err
	}

	err = r.db.db.Set([]byte("role:"+role.ID), data, &pebble.WriteOptions{Sync: true})
	if err != nil {
		logger.Logger.Errorf("Failed to store role in PebbleDB: %v", err)
		return err
	}

	logger.Logger.Infof("Role created successfully: %+v", role)
	return nil
}

func (r *RoleRepositoryImpl) Update(role domain.Role) error {
	_, err := r.GetByID(role.ID)
	if err != nil {
		if err == domain.ErrRoleNotFound {
			logger.Logger.Warnf("Role not found: %s", role.ID)
			return domain.ErrRoleNotFound
		}
		logger.Logger.Errorf("Error retrieving role %s: %v", role.ID, err)
		return err
	}

	data, err := json.Marshal(role)
	if err != nil {
		logger.Logger.Errorf("Failed to encode role to JSON: %v", err)
		return err
	}

	err = r.db.db.Set([]byte("role:"+role.ID), data, &pebble.WriteOptions{Sync: true})
	if err != nil {
		logger.Logger.Errorf("Failed to update role in PebbleDB: %v", err)
		return err
	}

	logger.Logger.Infof("Role updated successfully: %+v", role)
	return nil
}


func (r *RoleRepositoryImpl) Delete(roleID string) error {
	err := r.db.db.Delete([]byte("role:"+roleID), &pebble.WriteOptions{Sync: true})
	if err != nil {
		logger.Logger.Errorf("Failed to delete role %s: %v", roleID, err)
		return err
	}

	logger.Logger.Infof("Role deleted successfully: %s", roleID)
	return nil
}

func (r *RoleRepositoryImpl) GetByID(roleID string) (domain.Role, error) {
	var role domain.Role
	data, closer, err := r.db.db.Get([]byte("role:" + roleID))
	if err == pebble.ErrNotFound {
		logger.Logger.Warnf("Role not found: %s", roleID)
		return domain.Role{}, domain.ErrRoleNotFound
	}
	if err != nil {
		logger.Logger.Errorf("Failed to retrieve role %s: %v", roleID, err)
		return domain.Role{}, err
	}
	defer closer.Close()

	err = json.Unmarshal(data, &role)
	if err != nil {
		logger.Logger.Errorf("Failed to decode role JSON. Data: %s, Error: %v", string(data), err)
		return domain.Role{}, err
	}

	logger.Logger.Infof("Role retrieved successfully: %+v", role)
	return role, nil
}

func (r *RoleRepositoryImpl) GetAll() ([]domain.Role, error) {
	var roles []domain.Role

	iter, err := r.db.db.NewIter(nil)
	if err != nil {
		logger.Logger.Errorf("Failed to create PebbleDB iterator: %v", err)
		return nil, err
	}
	defer iter.Close()

	for iter.First(); iter.Valid(); iter.Next() {
		key := string(iter.Key())

		if strings.HasPrefix(key, "role:") {
			var role domain.Role
			data := iter.Value()

			err := json.Unmarshal(data, &role)
			if err != nil {
				logger.Logger.Errorf("Failed to decode role data. Key: %s, Error: %s", key, err)
				continue
			}

			roles = append(roles, role)
		}
	}

	if err := iter.Error(); err != nil {
		logger.Logger.Errorf("Error iterating over PebbleDB: %v", err)
		return nil, err
	}

	if len(roles) == 0 {
		logger.Logger.Info("No roles found in database")
		return []domain.Role{}, nil
	}

	logger.Logger.Infof("Retrieved %d roles from database", len(roles))
	return roles, nil
}
