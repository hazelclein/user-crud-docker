package domain

import (
	"errors"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User represents the user domain entity
type User struct {
	ID           int64     `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Never expose password in JSON
	Age          int       `json:"age"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewUser creates a new user with validation and password hashing
func NewUser(name, email, password string, age int) (*User, error) {
	// Trim whitespace
	name = strings.TrimSpace(name)
	email = strings.TrimSpace(email)
	password = strings.TrimSpace(password)
	
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}
	if email == "" {
		return nil, errors.New("email cannot be empty")
	}
	if password == "" {
		return nil, errors.New("password cannot be empty")
	}
	if len(password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	if age < 0 || age > 150 {
		return nil, errors.New("age must be between 0 and 150")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	now := time.Now()
	return &User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Age:          age,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// Update updates user fields with validation
func (u *User) Update(name, email string, age int) error {
	if name == "" {
		return errors.New("name cannot be empty")
	}
	if email == "" {
		return errors.New("email cannot be empty")
	}
	if age < 0 || age > 150 {
		return errors.New("age must be between 0 and 150")
	}

	u.Name = name
	u.Email = email
	u.Age = age
	u.UpdatedAt = time.Now()

	return nil
}

// UpdatePassword updates user password with validation
func (u *User) UpdatePassword(oldPassword, newPassword string) error {
	// Verify old password
	if err := u.ComparePassword(oldPassword); err != nil {
		return errors.New("old password is incorrect")
	}

	// Validate new password
	if newPassword == "" {
		return errors.New("new password cannot be empty")
	}
	if len(newPassword) < 8 {
		return errors.New("new password must be at least 8 characters")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	u.PasswordHash = string(hashedPassword)
	u.UpdatedAt = time.Now()

	return nil
}

// SetPassword sets a new password without verifying old password (for reset password)
func (u *User) SetPassword(newPassword string) error {
	if newPassword == "" {
		return errors.New("password cannot be empty")
	}
	if len(newPassword) < 8 {
		return errors.New("password must be at least 8 characters")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	u.PasswordHash = string(hashedPassword)
	u.UpdatedAt = time.Now()

	return nil
}

// ComparePassword compares given password with stored hash
func (u *User) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}

// ToPublicUser returns user without sensitive information
func (u *User) ToPublicUser() *PublicUser {
	return &PublicUser{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Age:       u.Age,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// PublicUser represents user data for public API responses
type PublicUser struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Common domain errors
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserData   = errors.New("invalid user data")
	ErrInvalidPassword   = errors.New("invalid password")
)