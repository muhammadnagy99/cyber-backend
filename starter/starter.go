package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	baseURL       = "http://localhost:8080"
	bearerToken   = "eyJ0eXAiOiJKV1QiLCJhbGciOiJFUzI1NiIsImtpZCI6IjY0YjZlYmEwM2RlZWE2ZTVjMjZjMTg1NDQ3ZmE4MDNjIn0.eyJzdWIiOiIyOTcyMTUxOTE5IiwibmFtZSI6IkJJR0JPU1MiLCJpYXQiOjEzMjEyMzEzMjEzMjF9.zN7mG-0pI2EBE2wsXu9jsdfud4uiqBiZDPgxrE0e2mJ4sD_CdesyQPANeEYp6c7log4haM8XbeMVr7P54oO-bQ"
)

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Permission struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createRole(role Role) {
	url := fmt.Sprintf("%s/roles", baseURL)
	jsonData, _ := json.Marshal(role)
	sendRequest("POST", url, jsonData)
}

func createPermission(permission Permission) {
	url := fmt.Sprintf("%s/permissions", baseURL)
	jsonData, _ := json.Marshal(permission)
	sendRequest("POST", url, jsonData)
}

func assignPermissionToRole(roleID, permissionID string) {
	url := fmt.Sprintf("%s/roles/%s/permissions/%s", baseURL, roleID, permissionID)
	sendRequest("POST", url, nil)
}

func sendRequest(method, url string, jsonData []byte) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}

defer resp.Body.Close()
	fmt.Printf("Response Status: %s for %s\n", resp.Status, url)
}

func main() {
	roles := []Role{
		{"super_admin", "Super Admin"},
		{"admin", "Admin"},
		{"security_officer", "Security Officer"},
		{"network_admin", "Network Administrator"},
		{"developer", "Developer"},
		{"devops_engineer", "DevOps Engineer"},
		{"compliance_officer", "Compliance Officer"},
		{"auditor", "Auditor"},
		{"support_agent", "Support Agent"},
		{"regular_user", "Regular User"},
	}

	permissions := []Permission{
		{"manage_users", "Manage Users"},
		{"manage_roles", "Manage Roles"},
		{"manage_permissions", "Manage Permissions"},
		{"view_logs", "View Logs"},
		{"manage_network", "Manage Network"},
		{"deploy_applications", "Deploy Applications"},
		{"access_sensitive_data", "Access Sensitive Data"},
		{"modify_configurations", "Modify Configurations"},
		{"monitor_security", "Monitor Security"},
		{"file_access", "File Access"},
		{"incident_response", "Incident Response"},
		{"view_reports", "View Reports"},
		{"provide_support", "Provide Support"},
		{"use_system", "Use System"},
	}

	// Create roles
	for _, role := range roles {
		createRole(role)
	}

	// Create permissions
	for _, permission := range permissions {
		createPermission(permission)
	}

	// Assign permissions to roles
	assignments := map[string][]string{
		"super_admin": {"manage_users", "manage_roles", "manage_permissions", "view_logs", "manage_network", "deploy_applications", "access_sensitive_data", "modify_configurations", "monitor_security", "file_access", "incident_response", "view_reports", "provide_support", "use_system"},
		"admin": {"manage_users", "manage_roles", "view_logs", "modify_configurations", "file_access", "view_reports", "use_system"},
		"security_officer": {"view_logs", "access_sensitive_data", "monitor_security", "incident_response", "view_reports", "use_system"},
		"network_admin": {"manage_network", "modify_configurations", "file_access", "use_system"},
		"developer": {"deploy_applications", "file_access", "use_system"},
		"devops_engineer": {"manage_network", "deploy_applications", "modify_configurations", "file_access", "view_logs", "use_system"},
		"compliance_officer": {"access_sensitive_data", "view_reports", "use_system"},
		"auditor": {"view_logs", "access_sensitive_data", "monitor_security", "view_reports", "use_system"},
		"support_agent": {"provide_support", "incident_response", "use_system"},
		"regular_user": {"use_system"},
	}

	for role, perms := range assignments {
		for _, perm := range perms {
			assignPermissionToRole(role, perm)
		}
	}
}