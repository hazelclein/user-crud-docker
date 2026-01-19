package command

import (
	"context"
	"user-crud/internal/domain"
	"user-crud/internal/infrastructure/cache"
	"user-crud/internal/infrastructure/tracing"
)

type UpdateUserCommand struct {
	ID    int64  `json:"-"`
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age" binding:"required,min=0,max=150"`
}

type UpdateUserHandler struct {
	repo  domain.UserRepository
	cache *cache.RedisCache
}

func NewUpdateUserHandler(repo domain.UserRepository, cache *cache.RedisCache) *UpdateUserHandler {
	return &UpdateUserHandler{repo: repo, cache: cache}
}

func (h *UpdateUserHandler) Handle(ctx context.Context, cmd UpdateUserCommand) (*domain.User, error) {
	ctx, span := tracing.StartSpan(ctx, "UpdateUserHandler.Handle")
	defer span.End()

	user, err := h.repo.GetByID(ctx, cmd.ID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	if user.Email != cmd.Email {
		existingUser, _ := h.repo.GetByEmail(ctx, cmd.Email)
		if existingUser != nil && existingUser.ID != cmd.ID {
			return nil, domain.ErrUserAlreadyExists
		}
	}

	if err := user.Update(cmd.Name, cmd.Email, cmd.Age); err != nil {
		return nil, err
	}

	if err := h.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	go h.cache.DeleteUser(context.Background(), cmd.ID)

	return user, nil
}