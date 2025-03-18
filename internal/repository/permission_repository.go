package repository

import (
	"encoding/json"
	"strings"

	"cyber-rbac/internal/domain"
	"cyber-rbac/pkg/logger"

	"github.com/cockroachdb/pebble"
)

type PermissionRepository interface {
	Create(permission domain.Permission) error
	Update(permission domain.Permission) error
	Delete(permissionID string) error
	GetByID(permissionID string) (domain.Permission, error)
	GetAll() ([]domain.Permission, error)
}

type PermissionRepositoryImpl struct {
	db *PebbleDB
}

func NewPermissionRepository(db *PebbleDB) PermissionRepository {
	return &PermissionRepositoryImpl{db: db}
}

func (p *PermissionRepositoryImpl) Create(permission domain.Permission) error {
	data, err := json.Marshal(permission)
	if err != nil {
		logger.Logger.Errorf("Failed to encode permission to JSON: %v", err)
		return err
	}

	err = p.db.db.Set([]byte("permission:"+permission.ID), data, &pebble.WriteOptions{Sync: true})
	if err != nil {
		logger.Logger.Errorf("Failed to store permission in PebbleDB: %v", err)
		return err
	}

	logger.Logger.Infof("Permission created successfully: %+v", permission)

	storedData, closer, err := p.db.db.Get([]byte("permission:" + permission.ID))
	if err != nil {
		logger.Logger.Errorf("Failed to verify stored permission data for ID: %s, Error: %v", permission.ID, err)
		return err
	}
	defer closer.Close()

	logger.Logger.Infof("Raw stored data for permission ID=%s: %s", permission.ID, string(storedData))

	return nil
}

func (p *PermissionRepositoryImpl) Update(permission domain.Permission) error {
	_, err := p.GetByID(permission.ID)
	if err != nil {
		if err == domain.ErrPermissionNotFound {
			logger.Logger.Warnf("Permission not found: %s", permission.ID)
			return domain.ErrPermissionNotFound
		}
		logger.Logger.Errorf("Error retrieving permission %s: %v", permission.ID, err)
		return err
	}

	data, err := json.Marshal(permission)
	if err != nil {
		logger.Logger.Errorf("Failed to encode permission to JSON: %v", err)
		return err
	}

	err = p.db.db.Set([]byte("permission:"+permission.ID), data, &pebble.WriteOptions{Sync: true})
	if err != nil {
		logger.Logger.Errorf("Failed to update permission in PebbleDB: %v", err)
		return err
	}

	logger.Logger.Infof("Permission updated successfully: %+v", permission)
	return nil
}

func (p *PermissionRepositoryImpl) Delete(permissionID string) error {
	err := p.db.db.Delete([]byte("permission:"+permissionID), &pebble.WriteOptions{Sync: true})
	if err != nil {
		logger.Logger.Errorf("Failed to delete permission %s: %v", permissionID, err)
		return err
	}

	logger.Logger.Infof("Permission deleted successfully: %s", permissionID)
	return nil
}

func (p *PermissionRepositoryImpl) GetByID(permissionID string) (domain.Permission, error) {
	var permission domain.Permission
	data, closer, err := p.db.db.Get([]byte("permission:" + permissionID))
	if err == pebble.ErrNotFound {
		logger.Logger.Warnf("Permission not found: %s", permissionID)
		return domain.Permission{}, domain.ErrPermissionNotFound
	}
	if err != nil {
		logger.Logger.Errorf("Failed to retrieve permission %s: %v", permissionID, err)
		return domain.Permission{}, err
	}
	defer closer.Close()

	err = json.Unmarshal(data, &permission)
	if err != nil {
		logger.Logger.Errorf("Failed to decode permission JSON. Data: %s, Error: %v", string(data), err)
		return domain.Permission{}, err
	}

	logger.Logger.Infof("Permission retrieved successfully: %+v", permission)
	return permission, nil
}

func (p *PermissionRepositoryImpl) GetAll() ([]domain.Permission, error) {
	var permissions []domain.Permission

	iter, err := p.db.db.NewIter(nil)
	if err != nil {
		logger.Logger.Errorf("Failed to create PebbleDB iterator: %v", err)
		return nil, err
	}
	defer iter.Close()

	logger.Logger.Info("Starting iteration over PebbleDB keys...")

	// Iterate through all keys in the database
	for iter.First(); iter.Valid(); iter.Next() {
		key := string(iter.Key())

		// Log every key found in the database
		logger.Logger.Infof("Found key in DB: %s", key)

		// ✅ Use strings.HasPrefix() to avoid errors with short keys
		if strings.HasPrefix(key, "permission:") {
			var permission domain.Permission
			data := iter.Value()

			// Log raw JSON data before decoding
			logger.Logger.Infof("Raw stored data for key %s: %s", key, string(data))

			// Ensure data is not empty before decoding
			if len(data) == 0 {
				logger.Logger.Warnf("Skipping empty value for key: %s", key)
				continue
			}

			// Decode JSON properly
			err := json.Unmarshal(data, &permission)
			if err != nil {
				logger.Logger.Errorf("Failed to decode permission data. Key: %s, Error: %s", key, err)
				continue // Skip corrupted data instead of stopping execution
			}

			permissions = append(permissions, permission)
		}
	}

	// Handle iteration errors
	if err := iter.Error(); err != nil {
		logger.Logger.Errorf("Error iterating over PebbleDB: %v", err)
		return nil, err
	}

	// Log retrieved count and return an empty array if no results
	if len(permissions) == 0 {
		logger.Logger.Info("No permissions found in database")
		return []domain.Permission{}, nil
	}

	logger.Logger.Infof("Successfully retrieved %d permissions from database", len(permissions))
	return permissions, nil
}
