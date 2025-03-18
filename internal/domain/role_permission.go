package domain

type RolePermission struct {
	RoleID       string   `json:"role_id"`
	Permissions  []string `json:"permissions"`
}