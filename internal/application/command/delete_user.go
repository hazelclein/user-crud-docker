package command

import (
	"context"
	"user-crud/internal/domain"
	"user-crud/internal/infrastructure/cache"
	"user-crud/internal/infrastructure/tracing"
)

type DeleteUserCommand struct {
	ID int64
}

type DeleteUserHandler struct {
	repo  domain.UserRepository
	cache *cache.RedisCache
}

func NewDeleteUserHandler(repo domain.UserRepository, cache *cache.RedisCache) *DeleteUserHandler {
	return &DeleteUserHandler{repo: repo, cache: cache}
}

func (h *DeleteUserHandler) Handle(ctx context.Context, cmd DeleteUserCommand) error {
	ctx, span := tracing.StartSpan(ctx, "DeleteUserHandler.Handle")
	defer span.End()

	_, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	if err := h.repo.Delete(ctx, cmd.ID); err != nil {
		return err
	}

	go h.cache.DeleteUser(context.Background(), cmd.ID)

	return nil
}