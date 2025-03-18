package domain

import "errors"

var (
	ErrRoleNotFound       = errors.New("role not found")
	ErrPermissionNotFound = errors.New("permission not found")
	ErrInvalidRoleID      = errors.New("invalid role ID")
	ErrInvalidPermissionID = errors.New("invalid permission ID")
	ErrEmptyFields        = errors.New("fields cannot be empty")
)
