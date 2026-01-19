package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"user-crud/internal/application/command"
	"user-crud/internal/application/query"
	"user-crud/internal/domain"
	"user-crud/internal/infrastructure/cache"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Handler struct {
	createUserHandler     *command.CreateUserHandler
	updateUserHandler     *command.UpdateUserHandler
	deleteUserHandler     *command.DeleteUserHandler
	changePasswordHandler *command.ChangePasswordHandler
	getUserHandler        *query.GetUserHandler
	listUsersHandler      *query.ListUsersHandler
	searchUsersHandler    *query.SearchUsersHandler
	db                    *pgxpool.Pool
	cache                 *cache.RedisCache
}

func NewHandler(
	createUserHandler *command.CreateUserHandler,
	updateUserHandler *command.UpdateUserHandler,
	deleteUserHandler *command.DeleteUserHandler,
	changePasswordHandler *command.ChangePasswordHandler,
	getUserHandler *query.GetUserHandler,
	listUsersHandler *query.ListUsersHandler,
	searchUsersHandler *query.SearchUsersHandler,
	db *pgxpool.Pool,
	cache *cache.RedisCache,
) *Handler {
	return &Handler{
		createUserHandler:     createUserHandler,
		updateUserHandler:     updateUserHandler,
		deleteUserHandler:     deleteUserHandler,
		changePasswordHandler: changePasswordHandler,
		getUserHandler:        getUserHandler,
		listUsersHandler:      listUsersHandler,
		searchUsersHandler:    searchUsersHandler,
		db:                    db,
		cache:                 cache,
	}
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the service is healthy
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 503 {object} map[string]interface{}
// @Router /health [get]
func (h *Handler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// Check database
	dbStatus := "connected"
	if err := h.db.Ping(ctx); err != nil {
		dbStatus = "disconnected"
	}

	// Check Redis
	redisStatus := "connected"
	if err := h.cache.Ping(ctx); err != nil {
		redisStatus = "disconnected"
	}

	status := "healthy"
	statusCode := http.StatusOK
	if dbStatus != "connected" || redisStatus != "connected" {
		status = "unhealthy"
		statusCode = http.StatusServiceUnavailable
	}

	c.JSON(statusCode, gin.H{
		"status":    status,
		"database":  dbStatus,
		"cache":     redisStatus,
		"timestamp": time.Now(),
	})
}

// Metrics godoc
// @Summary Get metrics
// @Description Get application metrics
// @Tags metrics
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /metrics [get]
func (h *Handler) Metrics(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Metrics endpoint - integrate with Prometheus here",
	})
}

// CreateUser godoc
// @Summary Create a new user
// @Description Create a new user with name, email, password, and age
// @Tags users
// @Accept json
// @Produce json
// @Param user body command.CreateUserCommand true "User data"
// @Success 201 {object} map[string]interface{} "User created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 409 {object} map[string]interface{} "User already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	var cmd command.CreateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	user, err := h.createUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if err == domain.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "user with this email already exists",
			})
			return
		}
		if err.Error() == "password cannot be empty" ||
			err.Error() == "password must be at least 8 characters" ||
			err.Error() == "name cannot be empty" ||
			err.Error() == "email cannot be empty" {
			c.JSON(http.StatusBadRequest, gin.H{
				"status":  "error",
				"message": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   user.ToPublicUser(),
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Description Get a single user by their ID (with Redis caching)
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User found"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid user id",
		})
		return
	}

	user, err := h.getUserHandler.Handle(c.Request.Context(), query.GetUserQuery{ID: id})
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "user not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user.ToPublicUser(),
	})
}

// ListUsers godoc
// @Summary List users with filters
// @Description Get paginated list of users with optional filters
// @Tags users
// @Produce json
// @Param search query string false "Search by name or email"
// @Param age_min query int false "Minimum age"
// @Param age_max query int false "Maximum age"
// @Param sort query string false "Sort field (id, name, email, age, created_at)"
// @Param order query string false "Sort order (asc, desc)"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{} "Users list"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	search := c.Query("search")
	ageMin, _ := strconv.Atoi(c.Query("age_min"))
	ageMax, _ := strconv.Atoi(c.Query("age_max"))
	sortBy := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	q := query.ListUsersQuery{
		Search: search,
		AgeMin: ageMin,
		AgeMax: ageMax,
		SortBy: sortBy,
		Order:  order,
		Page:   page,
		Limit:  limit,
	}

	result, err := h.listUsersHandler.Handle(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	publicUsers := make([]*domain.PublicUser, len(result.Users))
	for i, user := range result.Users {
		publicUsers[i] = user.ToPublicUser()
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"data":        publicUsers,
		"total":       result.Total,
		"page":        result.Page,
		"limit":       result.Limit,
		"total_pages": result.TotalPages,
	})
}

// SearchUsers godoc
// @Summary Search users
// @Description Search users by keyword
// @Tags users
// @Produce json
// @Param q query string true "Search keyword"
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{} "Search results"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/search [get]
func (h *Handler) SearchUsers(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "search keyword is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	q := query.SearchUsersQuery{
		Keyword: keyword,
		Page:    page,
		Limit:   limit,
	}

	result, err := h.searchUsersHandler.Handle(c.Request.Context(), q)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	publicUsers := make([]*domain.PublicUser, len(result.Users))
	for i, user := range result.Users {
		publicUsers[i] = user.ToPublicUser()
	}

	c.JSON(http.StatusOK, gin.H{
		"status":      "success",
		"data":        publicUsers,
		"total":       result.Total,
		"page":        result.Page,
		"limit":       result.Limit,
		"total_pages": result.TotalPages,
	})
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user information
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body command.UpdateUserCommand true "User data"
// @Success 200 {object} map[string]interface{} "User updated"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 409 {object} map[string]interface{} "Email already exists"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid user id",
		})
		return
	}

	var cmd command.UpdateUserCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	cmd.ID = id
	user, err := h.updateUserHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "user not found",
			})
			return
		}
		if err == domain.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "error",
				"message": "user with this email already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user.ToPublicUser(),
	})
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete a user by ID
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{} "User deleted"
// @Failure 400 {object} map[string]interface{} "Invalid user ID"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid user id",
		})
		return
	}

	err = h.deleteUserHandler.Handle(c.Request.Context(), command.DeleteUserCommand{ID: id})
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "user not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "user deleted successfully",
	})
}

// ChangePassword godoc
// @Summary Change user password
// @Description Change password for a user
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param password body command.ChangePasswordCommand true "Password data"
// @Success 200 {object} map[string]interface{} "Password changed"
// @Failure 400 {object} map[string]interface{} "Invalid input"
// @Failure 401 {object} map[string]interface{} "Incorrect old password"
// @Failure 404 {object} map[string]interface{} "User not found"
// @Failure 500 {object} map[string]interface{} "Internal server error"
// @Router /users/{id}/change-password [put]
func (h *Handler) ChangePassword(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "invalid user id",
		})
		return
	}

	var cmd command.ChangePasswordCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	cmd.UserID = id
	err = h.changePasswordHandler.Handle(c.Request.Context(), cmd)
	if err != nil {
		if err == domain.ErrUserNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"status":  "error",
				"message": "user not found",
			})
			return
		}
		if err.Error() == "old password is incorrect" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"status":  "error",
				"message": "old password is incorrect",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "password changed successfully",
	})
}