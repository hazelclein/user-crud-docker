package domain

import (
	"context"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id int64) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	
	// Search & Filter methods
	Search(ctx context.Context, keyword string, page, limit int) ([]*User, int64, error)
	FindWithFilters(ctx context.Context, filters interface{}) ([]*User, int64, error)
}