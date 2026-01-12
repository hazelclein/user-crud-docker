package command

import (
	"context"
	"user-crud/internal/domain"
)

// CreateUserCommand represents the command to create a user
type CreateUserCommand struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Age      int    `json:"age" binding:"required,min=0,max=150"`
}

// CreateUserHandler handles user creation
type CreateUserHandler struct {
	repo domain.UserRepository
}

// NewCreateUserHandler creates a new CreateUserHandler
func NewCreateUserHandler(repo domain.UserRepository) *CreateUserHandler {
	return &CreateUserHandler{repo: repo}
}

// Handle executes the create user command
func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
	// Check if user already exists
	existingUser, _ := h.repo.GetByEmail(ctx, cmd.Email)
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Create new user domain entity (with password hashing)
	user, err := domain.NewUser(cmd.Name, cmd.Email, cmd.Password, cmd.Age)
	if err != nil {
		return nil, err
	}

	// Persist to repository
	if err := h.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}