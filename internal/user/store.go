package user

import (
	"context"
	"database/sql"
	"errors"
)

type Store interface {
	GetAll(ctx context.Context) ([]*User, error)
	GetByID(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, user *User) (*User, error)
	Update(ctx context.Context, id int64, user *User) (*User, error)
	Delete(ctx context.Context, id int64) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) GetAll(ctx context.Context) ([]*User, error) {
	const query = `
		SELECT id, role_id, name, email, created_at, updated_at, deleted_at
		FROM users
		WHERE deleted_at IS NULL
		ORDER BY id
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*User, 0)
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID,
			&user.RoleID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return users, nil
}

func (s *store) GetByID(ctx context.Context, id int64) (*User, error) {
	const query = `
		SELECT id, role_id, name, email, created_at, updated_at, deleted_at
		FROM users
		WHERE id = $1 AND deleted_at IS NULL
	`
	var user User
	err := s.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *store) Create(ctx context.Context, user *User) (*User, error) {
	const roleQuery = `
		SELECT EXISTS(
			SELECT 1 FROM roles WHERE id = $1 AND deleted_at IS NULL
		)
	`
	var roleExists bool
	if err := s.db.QueryRowContext(ctx, roleQuery, user.RoleID).Scan(&roleExists); err != nil {
		return nil, err
	}
	if !roleExists {
		return nil, ErrRoleNotFound
	}

	const emailQuery = `
		SELECT EXISTS(
			SELECT 1 FROM users WHERE LOWER(email) = LOWER($1) AND deleted_at IS NULL
		)
	`
	var emailExists bool
	if err := s.db.QueryRowContext(ctx, emailQuery, user.Email).Scan(&emailExists); err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrConflict
	}

	const query = `
		INSERT INTO users (role_id, name, email)
		VALUES ($1, $2, $3)
		RETURNING id, role_id, name, email, created_at, updated_at, deleted_at
	`
	err := s.db.QueryRowContext(ctx, query, user.RoleID, user.Name, user.Email).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *store) Update(ctx context.Context, id int64, user *User) (*User, error) {
	const roleQuery = `
		SELECT EXISTS(
			SELECT 1 FROM roles WHERE id = $1 AND deleted_at IS NULL
		)
	`
	var roleExists bool
	if err := s.db.QueryRowContext(ctx, roleQuery, user.RoleID).Scan(&roleExists); err != nil {
		return nil, err
	}
	if !roleExists {
		return nil, ErrRoleNotFound
	}

	const emailQuery = `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE LOWER(email) = LOWER($1)
				AND id <> $2
				AND deleted_at IS NULL
		)
	`
	var emailExists bool
	if err := s.db.QueryRowContext(ctx, emailQuery, user.Email, id).Scan(&emailExists); err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ErrConflict
	}

	const query = `
		UPDATE users
		SET role_id = $2, name = $3, email = $4, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, role_id, name, email, created_at, updated_at, deleted_at
	`
	err := s.db.QueryRowContext(ctx, query, id, user.RoleID, user.Name, user.Email).Scan(
		&user.ID,
		&user.RoleID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *store) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE users
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id
	`
	var deletedID int64
	err := s.db.QueryRowContext(ctx, query, id).Scan(&deletedID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	return err
}
