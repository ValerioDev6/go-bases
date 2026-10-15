package service

import (
	"errors"
	"web/internal/model"
	"web/internal/store"
)

type Service struct {
	store store.Store
}

func NewService(s store.Store) *Service {
	return &Service{store: s}
}

func (s *Service) GetAllBooks() ([]*model.Book, error) {
	books, err := s.store.GetAll()
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *Service) GetBookByID(id int) (*model.Book, error) {
	book, err := s.store.GetByID(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return book, nil
}

func (s *Service) CreateBook(book *model.Book) (*model.Book, error) {
	if book.Title == "" {
		return nil, errors.New("book title is required")
	}

	newBook, err := s.store.Create(book)
	if err != nil {
		return nil, err
	}

	return newBook, nil
}

func (s *Service) UpdateBook(id int, book *model.Book) (*model.Book, error) {
	if book.Title == "" {
		return nil, errors.New("book title is required")
	}

	updatedBook, err := s.store.Update(id, book)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return nil, errors.New("book not found")
		}
		return nil, err
	}

	return updatedBook, nil
}

func (s *Service) DeleteBook(id int) error {
	err := s.store.Delete(id)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return errors.New("book not found")
		}
		return err
	}

	return nil
}
