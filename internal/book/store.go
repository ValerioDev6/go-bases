package book

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

type SQLStore struct {
	db *sql.DB
}

func NewSQLStore(db *sql.DB) *SQLStore {
	return &SQLStore{db: db}
}

func Migrate(ctx context.Context, db *sql.DB) error {
	const query = `
		CREATE TABLE IF NOT EXISTS books (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			author TEXT NOT NULL,
			price REAL NOT NULL
		)
	`
	if _, err := db.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("create books table: %w", err)
	}
	return nil
}

func (s *SQLStore) GetAll(ctx context.Context) ([]Book, error) {
	const query = "SELECT id, title, author, price FROM books ORDER BY id"
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query books: %w", err)
	}
	defer rows.Close()

	books := make([]Book, 0)
	for rows.Next() {
		var book Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
			return nil, fmt.Errorf("scan book: %w", err)
		}
		books = append(books, book)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate books: %w", err)
	}
	return books, nil
}

func (s *SQLStore) GetByID(ctx context.Context, id int) (Book, error) {
	const query = "SELECT id, title, author, price FROM books WHERE id = ?"
	var book Book
	if err := s.db.QueryRowContext(ctx, query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Price); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Book{}, ErrNotFound
		}
		return Book{}, fmt.Errorf("query book: %w", err)
	}
	return book, nil
}

func (s *SQLStore) Create(ctx context.Context, book Book) (Book, error) {
	const query = "INSERT INTO books (title, author, price) VALUES (?, ?, ?)"
	result, err := s.db.ExecContext(ctx, query, book.Title, book.Author, book.Price)
	if err != nil {
		return Book{}, fmt.Errorf("insert book: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return Book{}, fmt.Errorf("read inserted book id: %w", err)
	}
	book.ID = int(id)
	return book, nil
}

func (s *SQLStore) Update(ctx context.Context, id int, book Book) (Book, error) {
	const query = "UPDATE books SET title = ?, author = ?, price = ? WHERE id = ?"
	result, err := s.db.ExecContext(ctx, query, book.Title, book.Author, book.Price, id)
	if err != nil {
		return Book{}, fmt.Errorf("update book: %w", err)
	}
	if err := ensureAffected(result); err != nil {
		return Book{}, err
	}
	book.ID = id
	return book, nil
}

func (s *SQLStore) Delete(ctx context.Context, id int) error {
	const query = "DELETE FROM books WHERE id = ?"
	result, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete book: %w", err)
	}
	return ensureAffected(result)
}

func ensureAffected(result sql.Result) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("read affected rows: %w", err)
	}
	if rows == 0 {
		return ErrNotFound
	}
	return nil
}
