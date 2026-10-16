package role

import (
	"context"
	"strings"
)

type Store interface {
	GetAll(ctx context.Context) ([]Role, error)
	GetByID(ctx context.Context, id int64) (Role, error)
	Create(ctx context.Context, role Role) (Role, error)
	Update(ctx context.Context, id int64, role Role) (Role, error)
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetAll(ctx context.Context) ([]Role, error) {
	return s.store.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int64) (Role, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, role Role) (Role, error) {
	role.Name = strings.TrimSpace(role.Name)
	if role.Name == "" {
		return Role{}, ErrNameRequired
	}
	return s.store.Create(ctx, role)
}

func (s *Service) Update(ctx context.Context, id int64, role Role) (Role, error) {
	role.Name = strings.TrimSpace(role.Name)
	if role.Name == "" {
		return Role{}, ErrNameRequired
	}
	return s.store.Update(ctx, id, role)
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.store.Delete(ctx, id)
}
