package command

import (
	"context"
	"user-crud/internal/domain"
)

type DeleteUserCommand struct {
	ID int64
}

type DeleteUserHandler struct {
	repo domain.UserRepository
}

func NewDeleteUserHandler(repo domain.UserRepository) *DeleteUserHandler {
	return &DeleteUserHandler{repo: repo}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	_, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return err
	}

	return nil
}