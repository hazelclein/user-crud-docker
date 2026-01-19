package command

import (
	"context"
	"user-crud/internal/domain"
	"user-crud/internal/infrastructure/cache"
	"user-crud/internal/infrastructure/tracing"
)

type CreateUserCommand struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Age      int    `json:"age" binding:"required,min=0,max=150"`
}

type CreateUserHandler struct {
	repo  domain.UserRepository
	cache *cache.RedisCache
}

func NewCreateUserHandler(repo domain.UserRepository, cache *cache.RedisCache) *CreateUserHandler {
	return &CreateUserHandler{repo: repo, cache: cache}
}

func (h *CreateUserHandler) Handle(ctx context.Context, cmd CreateUserCommand) (*domain.User, error) {
	ctx, span := tracing.StartSpan(ctx, "CreateUserHandler.Handle")
	defer span.End()

	existingUser, _ := h.repo.GetByEmail(ctx, cmd.Email)
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	user, err := domain.NewUser(cmd.Name, cmd.Email, cmd.Password, cmd.Age)
	if err != nil {
		return nil, err
	}

	if err := h.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}