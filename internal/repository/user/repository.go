package user

import (
	"context"
	"errors"

	"authorization-service/internal/domain"
)

// ErrNotFound is returned when a user does not exist in storage.
var ErrNotFound = errors.New("user not found")

// Repository describes storage operations for users.
type Repository interface {
	// Create creates a new user record in storage.
	// It returns full User with ID and timestamps.
	Create(ctx context.Context, u domain.User) (domain.User, error)

	// GetByEmail looks up a user by email.
	GetByEmail(ctx context.Context, email string) (domain.User, error)
}
