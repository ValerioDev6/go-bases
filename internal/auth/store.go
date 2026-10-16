package auth

import (
	"context"
	"database/sql"

	"web/internal/user"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, name, email, passwordHash string) (*user.User, error) {
	const query = `
		INSERT INTO users (name, email, password)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, name, email, passwordHash)
	if err != nil {
		return nil, err
	}
	return r.GetUserByEmail(ctx, email)
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*user.User, error) {
	const query = `
		SELECT id, role_id, name, email, password, created_at, updated_at, deleted_at
		FROM users
		WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL
	`

	var foundUser user.User
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&foundUser.ID,
		&foundUser.RoleID,
		&foundUser.Name,
		&foundUser.Email,
		&foundUser.PasswordHash,
		&foundUser.CreatedAt,
		&foundUser.UpdatedAt,
		&foundUser.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*user.User, error) {
	const query = `
		SELECT id, role_id, name, email, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`

	var foundUser user.User
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&foundUser.ID,
		&foundUser.RoleID,
		&foundUser.Name,
		&foundUser.Email,
		&foundUser.CreatedAt,
		&foundUser.UpdatedAt,
		&foundUser.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return &foundUser, nil
}
