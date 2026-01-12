package query

import (
	"context"
	"user-crud/internal/domain"
)

// ListUsersQuery represents the query to list users with filters
type ListUsersQuery struct {
	Search   string // Search by name or email
	AgeMin   int    // Minimum age filter
	AgeMax   int    // Maximum age filter
	SortBy   string // Sort field: "name", "email", "age", "created_at"
	Order    string // Sort order: "asc" or "desc"
	Page     int    // Page number (starts from 1)
	Limit    int    // Items per page
}

// ListUsersResult represents paginated user list result
type ListUsersResult struct {
	Users      []*domain.User `json:"users"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// ListUsersHandler handles listing users with filters
type ListUsersHandler struct {
	repo domain.UserRepository
}

// NewListUsersHandler creates a new ListUsersHandler
func NewListUsersHandler(repo domain.UserRepository) *ListUsersHandler {
	return &ListUsersHandler{repo: repo}
}

// Handle executes the list users query with filters
func (h *ListUsersHandler) Handle(ctx context.Context, query ListUsersQuery) (*ListUsersResult, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100 // Maximum 100 items per page
	}
	if query.SortBy == "" {
		query.SortBy = "id"
	}
	if query.Order == "" {
		query.Order = "asc"
	}

	// Get filtered users from repository
	users, total, err := h.repo.FindWithFilters(ctx, query)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / query.Limit
	if int(total)%query.Limit > 0 {
		totalPages++
	}

	return &ListUsersResult{
		Users:      users,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
	}, nil
}

// SearchUsersQuery represents the query to search users
type SearchUsersQuery struct {
	Keyword string
	Page    int
	Limit   int
}

// SearchUsersHandler handles user search
type SearchUsersHandler struct {
	repo domain.UserRepository
}

// NewSearchUsersHandler creates a new SearchUsersHandler
func NewSearchUsersHandler(repo domain.UserRepository) *SearchUsersHandler {
	return &SearchUsersHandler{repo: repo}
}

// Handle executes the search users query
func (h *SearchUsersHandler) Handle(ctx context.Context, query SearchUsersQuery) (*ListUsersResult, error) {
	// Set defaults
	if query.Page < 1 {
		query.Page = 1
	}
	if query.Limit < 1 {
		query.Limit = 10
	}
	if query.Limit > 100 {
		query.Limit = 100
	}

	// Search users
	users, total, err := h.repo.Search(ctx, query.Keyword, query.Page, query.Limit)
	if err != nil {
		return nil, err
	}

	// Calculate total pages
	totalPages := int(total) / query.Limit
	if int(total)%query.Limit > 0 {
		totalPages++
	}

	return &ListUsersResult{
		Users:      users,
		Total:      total,
		Page:       query.Page,
		Limit:      query.Limit,
		TotalPages: totalPages,
	}, nil
}