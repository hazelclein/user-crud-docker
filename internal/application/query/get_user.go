package query

import (
	"context"
	"log"

	"user-crud/internal/domain"
	"user-crud/internal/infrastructure/cache"
	"user-crud/internal/infrastructure/tracing"
)

type GetUserQuery struct {
	ID int64
}

type GetUserHandler struct {
	repo  domain.UserRepository
	cache *cache.RedisCache
}

func NewGetUserHandler(repo domain.UserRepository, cache *cache.RedisCache) *GetUserHandler {
	return &GetUserHandler{
		repo:  repo,
		cache: cache,
	}
}

func (h *GetUserHandler) Handle(ctx context.Context, query GetUserQuery) (*domain.User, error) {
	ctx, span := tracing.StartSpan(ctx, "GetUserHandler.Handle")
	defer span.End()

	// Try cache first
	ctx, cacheSpan := tracing.StartSpan(ctx, "cache.GetUser")
	user, err := h.cache.GetUser(ctx, query.ID)
	cacheSpan.End()

	if err != nil {
		log.Printf("Cache error: %v", err)
	}

	if user != nil {
		log.Printf("Cache HIT for user ID: %d", query.ID)
		return user, nil
	}

	log.Printf("Cache MISS for user ID: %d", query.ID)

	// Get from database
	ctx, dbSpan := tracing.StartSpan(ctx, "repository.GetByID")
	user, err = h.repo.GetByID(ctx, query.ID)
	dbSpan.End()

	if err != nil {
		return nil, domain.ErrUserNotFound
	}

	// Store in cache (async)
	go func() {
		if err := h.cache.SetUser(context.Background(), user); err != nil {
			log.Printf("Failed to cache user: %v", err)
		}
	}()

	return user, nil
}