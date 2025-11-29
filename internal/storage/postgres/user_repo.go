package postgres

import (
	"authorization-service/internal/domain"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"

	userrepo "authorization-service/internal/repository/user"
)

// UserRepository is a Postgres implementation of user.Repository.
type UserRepository struct {
	log  *slog.Logger
	pool *pgxpool.Pool
}

// NewUserRepository constructs a new Postgres-backed user repository.
func NewUserRepository(log *slog.Logger, pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		log:  log,
		pool: pool,
	}
}

// Ensure interface implementation at compile time.
var _ userrepo.Repository = (*UserRepository)(nil)

// Create inserts a new user into the database and returns full entity with ID and timestamps.
func (r *UserRepository) Create(ctx context.Context, u domain.User) (domain.User, error) {
	const op = "UserRepository.Create"
	query := `
		INSERT INTO users (
			email,
			login,
			password_hash,
			email_verified,
			github_id,
			google_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING 
			id,
			email,
			login,
			password_hash,
			email_verified,
			github_id,
			google_id,
			created_at,
			updated_at
	`
	var (
		dbGithub sql.NullString
		dbGoogle sql.NullString
		res      domain.User
	)

	err := r.pool.QueryRow(ctx, query,
		u.Email,
		u.Login,
		u.PasswordHash,
		u.EmailVerified,
		u.GithubID,
		u.GoogleID,
	).Scan(
		&res.ID,
		&res.Email,
		&res.Login,
		&res.PasswordHash,
		&res.EmailVerified,
		&dbGithub,
		&dbGoogle,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		r.log.Error(op+" failed",
			slog.String("email", u.Email),
			slog.Any("err", err),
		)
		return domain.User{}, err
	}

	if dbGithub.Valid {
		g := dbGithub.String
		res.GithubID = &g
	}
	if dbGoogle.Valid {
		g := dbGoogle.String
		res.GoogleID = &g
	}

	return res, nil
}

// GetByEmail looks up a user by email.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	const op = "UserRepository.GetByEmail"

	query := `
		SELECT
			id,
			email,
			login,
			password_hash,
			email_verified,
			github_id,
			google_id,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
		LIMIT 1
	`

	var (
		u        domain.User
		dbGithub sql.NullString
		dbGoogle sql.NullString
	)

	err := r.pool.QueryRow(ctx, query, email).Scan(
		&u.ID,
		&u.Email,
		&u.Login,
		&u.PasswordHash,
		&u.EmailVerified,
		&dbGithub,
		&dbGoogle,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.User{}, userrepo.ErrNotFound
		}

		r.log.Error(op+" failed",
			slog.String("email", email),
			slog.Any("err", err),
		)
		return domain.User{}, err
	}

	if dbGithub.Valid {
		g := dbGithub.String
		u.GithubID = &g
	}
	if dbGoogle.Valid {
		g := dbGoogle.String
		u.GoogleID = &g
	}

	return u, nil
}
