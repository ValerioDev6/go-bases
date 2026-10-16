package role

import (
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("role not found")
	ErrNameRequired = errors.New("role name is required")
	ErrConflict     = errors.New("role name already exists")
)

type Role struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
