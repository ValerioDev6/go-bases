package store

import (
	"database/sql"
	"errors"
	"web/internal/model"
)

var ErrNotFound = errors.New("libro no encontrado")

type Store interface {
	GetAll() ([]*model.Book, error)
	GetByID(id int) (*model.Book, error)
	Create(book *model.Book) (*model.Book, error)
	Update(id int, book *model.Book) (*model.Book, error)
	Delete(id int) error
}

type store struct {
	db *sql.DB
}

func New(db *sql.DB) Store {
	return &store{db: db}
}

func (s *store) GetAll() ([]*model.Book, error) {
	q := "SELECT id, title, author, price FROM books"
	rows, err := s.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []*model.Book
	for rows.Next() {
		var b model.Book
		if err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Price); err != nil {
			return nil, err
		}
		books = append(books, &b)
	}
	return books, nil
}

func (s *store) GetByID(id int) (*model.Book, error) {
	q := "SELECT id, title, author, price FROM books where id=?"

	var b model.Book
	err := s.db.QueryRow(q, id).Scan(&b.ID, &b.Title, &b.Author, &b.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &b, nil
}

func (s *store) Create(book *model.Book) (*model.Book, error) {
	q := `INSERT INTO books (title, author, price) VALUES (?, ?, ?)`
	res, err := s.db.Exec(q, book.Title, book.Author, book.Price)
	if err != nil {
		return nil, err
	}
	id, _ := res.LastInsertId()
	book.ID = int(id)
	return book, nil
}

func (s *store) Update(id int, book *model.Book) (*model.Book, error) {
	q := "UPDATE books SET title=?, author=?, price=? WHERE id=?"
	_, err := s.db.Exec(q, book.Title, book.Author, book.Price, id)
	if err != nil {
		return nil, err
	}
	book.ID = id
	return book, nil
}

func (s *store) Delete(id int) error {
	q := "DELETE FROM books WHERE id=?"
	_, err := s.db.Exec(q, id)
	if err != nil {
		return err
	}
	return nil
}
