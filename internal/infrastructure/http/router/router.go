package router

import (
	"user-crud/internal/infrastructure/http/handler"
	"user-crud/internal/infrastructure/http/middleware"

	_ "user-crud/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/time/rate"
)

func SetupRouter(h *handler.Handler) *gin.Engine {
	// Release mode
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// Global middleware
	r.Use(
		gin.Recovery(),
		gin.Logger(),
		middleware.TracingMiddleware("user-crud-api"),
		middleware.CircuitBreakerMiddleware(),
	)

	// Rate limiter global
	rateLimiter := middleware.NewRateLimiter(rate.Limit(10), 20)
	r.Use(rateLimiter.Middleware())

	// ===== Infra endpoints (ROOT) =====
	r.GET("/health", h.HealthCheck)
	r.GET("/metrics", h.Metrics)

	// Swagger (infra, bukan API bisnis)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ===== API v1 =====
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			users := v1.Group("/users")
			{
				users.POST("", h.CreateUser)
				users.GET("", h.ListUsers)
				users.GET("/search", h.SearchUsers)
				users.GET("/:id", h.GetUser)
				users.PUT("/:id", h.UpdateUser)
				users.DELETE("/:id", h.DeleteUser)
				users.PUT("/:id/change-password", h.ChangePassword)
			}
		}
	}

	return r
}
