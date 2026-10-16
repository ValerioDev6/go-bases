package book

import (
	"context"
	"strings"
)

type Store interface {
	GetAll(ctx context.Context) ([]Book, error)
	GetByID(ctx context.Context, id int) (Book, error)
	Create(ctx context.Context, book Book) (Book, error)
	Update(ctx context.Context, id int, book Book) (Book, error)
	Delete(ctx context.Context, id int) error
}

type Service struct {
	store Store
}

func NewService(store Store) *Service {
	return &Service{store: store}
}

func (s *Service) GetAll(ctx context.Context) ([]Book, error) {
	return s.store.GetAll(ctx)
}

func (s *Service) GetByID(ctx context.Context, id int) (Book, error) {
	return s.store.GetByID(ctx, id)
}

func (s *Service) Create(ctx context.Context, book Book) (Book, error) {
	if strings.TrimSpace(book.Title) == "" {
		return Book{}, ErrTitleRequired
	}
	return s.store.Create(ctx, book)
}

func (s *Service) Update(ctx context.Context, id int, book Book) (Book, error) {
	if strings.TrimSpace(book.Title) == "" {
		return Book{}, ErrTitleRequired
	}
	return s.store.Update(ctx, id, book)
}

func (s *Service) Delete(ctx context.Context, id int) error {
	return s.store.Delete(ctx, id)
}
