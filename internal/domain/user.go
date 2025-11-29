package domain

import "time"

// User is a domain model representing an application user.
// It is decoupled from both protobuf and database details.
type User struct {
	ID            int64
	Email         string
	Login         string
	PasswordHash  string
	EmailVerified bool

	GithubID *string
	GoogleID *string

	CreatedAt time.Time
	UpdatedAt time.Time
}
