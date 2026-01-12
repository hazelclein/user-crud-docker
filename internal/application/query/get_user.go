package query

import (
	"context"
	"user-crud/internal/domain"
)

type GetUserQuery struct {
	ID int64
}

type GetUserHandler struct {
	repo domain.UserRepository
}

func NewGetUserHandler(repo domain.UserRepository) *GetUserHandler {
	return &GetUserHandler{repo: repo}
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*domain.User, error) {
	user, err := h.repo.GetByID(ctx, query.ID)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}