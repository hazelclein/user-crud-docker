package router

import (
	"user-crud/internal/infrastructure/http/handler"

	"github.com/gin-gonic/gin"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	// Set Gin to release mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	// Health check endpoint
	r.GET("/health", h.HealthCheck)

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// User routes
		v1.POST("/users", h.CreateUser)
		v1.GET("/users", h.ListUsers)                           // With filters & pagination
		v1.GET("/users/search", h.SearchUsers)                  // Search endpoint
		v1.GET("/users/:id", h.GetUser)
		v1.PUT("/users/:id", h.UpdateUser)
		v1.DELETE("/users/:id", h.DeleteUser)
		v1.PUT("/users/:id/change-password", h.ChangePassword)
	}

	return r
}