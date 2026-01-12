package handler

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"user-crud/internal/application/command"
	"user-crud/internal/application/query"
	"user-crud/internal/domain"

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
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "error",
			"message":   "database connection failed",
			"timestamp": time.Now(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"database":  "connected",
		"timestamp": time.Now(),
	})
}

// CreateUser handles user creation
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

// GetUser handles getting a user by ID
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

// ListUsers handles listing users with filters and pagination
func (h *Handler) ListUsers(c *gin.Context) {
	// Parse query parameters
	search := c.Query("search")
	ageMin, _ := strconv.Atoi(c.Query("age_min"))
	ageMax, _ := strconv.Atoi(c.Query("age_max"))
	sortBy := c.DefaultQuery("sort", "id")
	order := c.DefaultQuery("order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	// Create query
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

	// Convert to public users
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

// SearchUsers handles user search
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

	// Convert to public users
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

// UpdateUser handles user updates
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

// DeleteUser handles user deletion
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

// ChangePassword handles password change requests
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