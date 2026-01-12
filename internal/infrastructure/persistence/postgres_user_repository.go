package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"user-crud/internal/domain"
	"user-crud/internal/application/query"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresUserRepository struct {
	db *pgxpool.Pool
}

func NewPostgresUserRepository(db *pgxpool.Pool) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, age, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Age,
		user.CreatedAt,
		user.UpdatedAt,
	).Scan(&user.ID)

	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, age, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, age, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Age,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, age, created_at, updated_at
		FROM users
		ORDER BY id
	`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Age,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *PostgresUserRepository) Update(ctx context.Context, user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1, email = $2, password_hash = $3, age = $4, updated_at = $5
		WHERE id = $6
	`

	result, err := r.db.Exec(
		ctx,
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Age,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

func (r *PostgresUserRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return domain.ErrUserNotFound
	}

	return nil
}

// Search searches users by name or email (ILIKE for case-insensitive)
func (r *PostgresUserRepository) Search(ctx context.Context, keyword string, page, limit int) ([]*domain.User, int64, error) {
	// Calculate offset
	offset := (page - 1) * limit

	// Search query
	searchQuery := `
		SELECT id, name, email, password_hash, age, created_at, updated_at
		FROM users
		WHERE name ILIKE $1 OR email ILIKE $1
		ORDER BY id
		LIMIT $2 OFFSET $3
	`

	// Count query
	countQuery := `
		SELECT COUNT(*)
		FROM users
		WHERE name ILIKE $1 OR email ILIKE $1
	`

	searchPattern := "%" + keyword + "%"

	// Get total count
	var total int64
	err := r.db.QueryRow(ctx, countQuery, searchPattern).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	rows, err := r.db.Query(ctx, searchQuery, searchPattern, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Age,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// FindWithFilters finds users with multiple filters
func (r *PostgresUserRepository) FindWithFilters(ctx context.Context, filters interface{}) ([]*domain.User, int64, error) {
	// Cast filters to ListUsersQuery
	q, ok := filters.(query.ListUsersQuery)
	if !ok {
		return nil, 0, fmt.Errorf("invalid filter type")
	}

	// Build WHERE clause
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Search filter
	if q.Search != "" {
		conditions = append(conditions, fmt.Sprintf("(name ILIKE $%d OR email ILIKE $%d)", argIndex, argIndex))
		args = append(args, "%"+q.Search+"%")
		argIndex++
	}

	// Age min filter
	if q.AgeMin > 0 {
		conditions = append(conditions, fmt.Sprintf("age >= $%d", argIndex))
		args = append(args, q.AgeMin)
		argIndex++
	}

	// Age max filter
	if q.AgeMax > 0 {
		conditions = append(conditions, fmt.Sprintf("age <= $%d", argIndex))
		args = append(args, q.AgeMax)
		argIndex++
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Validate sort field
	validSortFields := map[string]bool{
		"id":         true,
		"name":       true,
		"email":      true,
		"age":        true,
		"created_at": true,
	}
	sortBy := q.SortBy
	if !validSortFields[sortBy] {
		sortBy = "id"
	}

	// Validate order
	order := q.Order
	if order != "asc" && order != "desc" {
		order = "asc"
	}

	// Build ORDER BY clause
	orderClause := fmt.Sprintf("ORDER BY %s %s", sortBy, strings.ToUpper(order))

	// Calculate offset
	offset := (q.Page - 1) * q.Limit

	// Count query
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)

	// Get total count
	var total int64
	err := r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Main query with pagination
	mainQuery := fmt.Sprintf(`
		SELECT id, name, email, password_hash, age, created_at, updated_at
		FROM users
		%s
		%s
		LIMIT $%d OFFSET $%d
	`, whereClause, orderClause, argIndex, argIndex+1)

	args = append(args, q.Limit, offset)

	// Get users
	rows, err := r.db.Query(ctx, mainQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var user domain.User
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Age,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, total, nil
}