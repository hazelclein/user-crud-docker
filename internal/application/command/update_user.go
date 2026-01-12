package command

import (
	"context"
	"user-crud/internal/domain"
)

// UpdateUserCommand represents the command to update a user
type UpdateUserCommand struct {
	ID    int64  `json:"-"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,min=0,max=150"`
}

// UpdateUserHandler handles user updates
type UpdateUserHandler struct {
	repo domain.UserRepository
}

// NewUpdateUserHandler creates a new UpdateUserHandler
func NewUpdateUserHandler(repo domain.UserRepository) *UpdateUserHandler {
	return &UpdateUserHandler{repo: repo}
}

// Handle executes the update user command
func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*domain.User, error) {
	// Get existing user
	user, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Check if email is being changed to an existing email
	if user.Email != cmd.Email {
		existingUser, _ := h.repo.GetByEmail(ctx, cmd.Email)
		if existingUser != nil && existingUser.ID != cmd.ID {
			return nil, domain.ErrUserAlreadyExists
		}
	}

	// Update user entity
	if err := user.Update(cmd.Name, cmd.Email, cmd.Age); err != nil {
		return nil, err
	}

	// Persist changes
	if err := h.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePasswordCommand represents the command to change user password
type ChangePasswordCommand struct {
	UserID      int64  `json:"-"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordHandler handles password changes
type ChangePasswordHandler struct {
	repo domain.UserRepository
}

// NewChangePasswordHandler creates a new ChangePasswordHandler
func NewChangePasswordHandler(repo domain.UserRepository) *ChangePasswordHandler {
	return &ChangePasswordHandler{repo: repo}
}

// Handle executes the change password command
func (h *ChangePasswordHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) error {
	// Get user
	user, err := h.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	// Update password (validates old password internally)
	if err := user.UpdatePassword(cmd.OldPassword, cmd.NewPassword); err != nil {
		return err
	}

	// Persist changes
	if err := h.repo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}