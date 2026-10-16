package user

import (
	"errors"
	"time"
)

var (
	ErrNotFound      = errors.New("user not found")
	ErrNameRequired  = errors.New("user name is required")
	ErrEmailRequired = errors.New("user email is required")
	ErrRoleRequired  = errors.New("user role_id is required")
	ErrRoleNotFound  = errors.New("role not found")
	ErrConflict      = errors.New("user email already exists")
)

type User struct {
	ID           int64      `json:"id"`
	RoleID       int64      `json:"role_id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	DeletedAt    *time.Time `json:"deleted_at,omitempty"`
}
