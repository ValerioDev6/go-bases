package user

import (
	"context"
	"errors"
	"strings"
)

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetAllUsers(ctx context.Context) ([]*User, error) {
	users, err := s.store.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *Service) GetUserByID(ctx context.Context, id int64) (*User, error) {
	user, err := s.store.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

func (s *Service) CreateUser(ctx context.Context, user *User) (*User, error) {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	if user.Name == "" {
		return nil, ErrNameRequired
	}
	if user.Email == "" {
		return nil, ErrEmailRequired
	}
	if user.RoleID <= 0 {
		return nil, ErrRoleRequired
	}

	newUser, err := s.store.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (s *Service) UpdateUser(ctx context.Context, id int64, user *User) (*User, error) {
	user.Name = strings.TrimSpace(user.Name)
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	if user.Name == "" {
		return nil, ErrNameRequired
	}
	if user.Email == "" {
		return nil, ErrEmailRequired
	}
	if user.RoleID <= 0 {
		return nil, ErrRoleRequired
	}

	updatedUser, err := s.store.Update(ctx, id, user)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return updatedUser, nil
}

func (s *Service) DeleteUser(ctx context.Context, id int64) error {
	err := s.store.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return err
	}

	return nil
}
