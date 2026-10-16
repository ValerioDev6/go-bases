package role

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
)

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func (s *SQLStore) GetAll(ctx context.Context) ([]Role, error) {
	const query = `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM roles
		WHERE deleted_at IS NULL
		ORDER BY id
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query roles: %w", err)
	}
	defer rows.Close()

	roles := make([]Role, 0)
	for rows.Next() {
		var role Role
		if err := scanRole(rows, &role); err != nil {
			return nil, fmt.Errorf("scan role: %w", err)
		}
		roles = append(roles, role)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate roles: %w", err)
	}
	return roles, nil
}

func (s *SQLStore) GetByID(ctx context.Context, id int64) (Role, error) {
	const query = `
		SELECT id, name, created_at, updated_at, deleted_at
		FROM roles
		WHERE id = $1 AND deleted_at IS NULL
	`
	var role Role
	if err := scanRole(s.db.QueryRowContext(ctx, query, id), &role); err != nil {
		return Role{}, storeError("query role", err)
	}
	return role, nil
}

func (s *SQLStore) Create(ctx context.Context, role Role) (Role, error) {
	const query = `
		INSERT INTO roles (name)
		VALUES ($1)
		RETURNING id, name, created_at, updated_at, deleted_at
	`
	if err := scanRole(s.db.QueryRowContext(ctx, query, role.Name), &role); err != nil {
		return Role{}, storeError("create role", err)
	}
	return role, nil
}

func (s *SQLStore) Update(ctx context.Context, id int64, role Role) (Role, error) {
	const query = `
		UPDATE roles
		SET name = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id, name, created_at, updated_at, deleted_at
	`
	if err := scanRole(s.db.QueryRowContext(ctx, query, id, role.Name), &role); err != nil {
		return Role{}, storeError("update role", err)
	}
	return role, nil
}

func (s *SQLStore) Delete(ctx context.Context, id int64) error {
	const query = `
		UPDATE roles
		SET deleted_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING id
	`
	var deletedID int64
	if err := s.db.QueryRowContext(ctx, query, id).Scan(&deletedID); err != nil {
		return storeError("delete role", err)
	}
	return nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanRole(row rowScanner, role *Role) error {
	return row.Scan(&role.ID, &role.Name, &role.CreatedAt, &role.UpdatedAt, &role.DeletedAt)
}

func storeError(operation string, err error) error {
	if errors.Is(err, sql.ErrNoRows) {
		return ErrNotFound
	}
	var postgresError *pgconn.PgError
	if errors.As(err, &postgresError) && postgresError.Code == "23505" {
		return ErrConflict
	}
	return fmt.Errorf("%s: %w", operation, err)
}
