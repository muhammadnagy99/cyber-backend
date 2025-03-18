package main

import (
	"log"
	"net/http"

	api_handlers "cyber-rbac/internal/http"
	"cyber-rbac/internal/middleware"
	"cyber-rbac/internal/repository"
	"cyber-rbac/internal/usecase"
	"cyber-rbac/pkg/config"
	"cyber-rbac/pkg/logger"

	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

func main() {
	// Load Configuration
	cfg := config.LoadConfig()

	// Set up logging
	logger.Logger.Info("Starting RBAC server...")

	// Initialize database
	db, err := repository.NewPebbleDB(cfg.DatabasePath)
	if err != nil {
		logger.Logger.Fatal("Failed to initialize database: ", err)
	}

	// Initialize repositories
	roleRepo := repository.NewRoleRepository(db)
	permissionRepo := repository.NewPermissionRepository(db)
	rolePermissionRepo := repository.NewRolePermissionRepository(db)

	// Initialize use cases
	roleUseCase := usecase.NewRoleUseCase(roleRepo)
	permissionUseCase := usecase.NewPermissionUseCase(permissionRepo)
	rolePermissionUseCase := usecase.NewRolePermissionUseCase(rolePermissionRepo)

	// Initialize handlers
	roleHandler := api_handlers.NewRoleHandler(roleUseCase)
	permissionHandler := api_handlers.NewPermissionHandler(permissionUseCase)
	rolePermissionHandler := api_handlers.NewRolePermissionHandler(rolePermissionUseCase)

	// Setup Router
	router := mux.NewRouter()

	// Protected routes (Require Authentication)
	authRouter := router.PathPrefix("/").Subrouter()
	authRouter.Use(func(next http.Handler) http.Handler {
		return middleware.AuthMiddleware(cfg, next)
	})

	// Role routes (protected)
	authRouter.HandleFunc("/roles", roleHandler.CreateRole).Methods("POST")
	authRouter.HandleFunc("/roles", roleHandler.GetAllRoles).Methods("GET")
	authRouter.HandleFunc("/roles/{id}", roleHandler.GetRoleByID).Methods("GET")
	authRouter.HandleFunc("/roles/{id}", roleHandler.UpdateRole).Methods("PUT")
	authRouter.HandleFunc("/roles/{id}", roleHandler.DeleteRole).Methods("DELETE")

	// Permission routes (protected)
	authRouter.HandleFunc("/permissions", permissionHandler.CreatePermission).Methods("POST")
	authRouter.HandleFunc("/permissions", permissionHandler.GetAllPermissions).Methods("GET")
	authRouter.HandleFunc("/permissions/{id}", permissionHandler.GetPermissionByID).Methods("GET")
	authRouter.HandleFunc("/permissions/{id}", permissionHandler.UpdatePermission).Methods("PUT")
	authRouter.HandleFunc("/permissions/{id}", permissionHandler.DeletePermission).Methods("DELETE")

	// Role-Permission routes (protected)
	authRouter.HandleFunc("/roles/{role_id}/permissions/{permission_id}", rolePermissionHandler.AssignPermission).Methods("POST")
	authRouter.HandleFunc("/roles/{role_id}/permissions/{permission_id}", rolePermissionHandler.RemovePermission).Methods("DELETE")
	authRouter.HandleFunc("/roles/{role_id}/permissions", rolePermissionHandler.GetPermissions).Methods("GET")
	authRouter.HandleFunc("/roles/{role_id}/permissions/{permission_id}/check", rolePermissionHandler.HasPermission).Methods("GET")

	// Apply CORS Middleware
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}), // Include OPTIONS for preflight
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}), // Allow Authorization header
		handlers.AllowCredentials(), // Allow credentials (if needed)
	)(router)

	// Start server
	logger.Logger.Infof("Server is running on port %s...", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, corsHandler))
}
