package command

import (
	"context"
	"user-crud/internal/domain"
	"user-crud/internal/infrastructure/cache"
	"user-crud/internal/infrastructure/tracing"
)

type ChangePasswordCommand struct {
	UserID      int64  `json:"-"`
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type ChangePasswordHandler struct {
	repo  domain.UserRepository
	cache *cache.RedisCache
}

func NewChangePasswordHandler(repo domain.UserRepository, cache *cache.RedisCache) *ChangePasswordHandler {
	return &ChangePasswordHandler{repo: repo, cache: cache}
}

func (h *ChangePasswordHandler) Handle(ctx context.Context, cmd ChangePasswordCommand) error {
	ctx, span := tracing.StartSpan(ctx, "ChangePasswordHandler.Handle")
	defer span.End()

	user, err := h.repo.GetByID(ctx, cmd.UserID)
	if err != nil {
		return domain.ErrUserNotFound
	}

	if err := user.UpdatePassword(cmd.OldPassword, cmd.NewPassword); err != nil {
		return err
	}

	if err := h.repo.Update(ctx, user); err != nil {
		return err
	}

	go h.cache.DeleteUser(context.Background(), cmd.UserID)

	return nil
}