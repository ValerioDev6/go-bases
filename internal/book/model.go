package book

import "errors"

var (
	ErrNotFound      = errors.New("book not found")
	ErrTitleRequired = errors.New("book title is required")
)

type Book struct {
	ID     int     `json:"id"`
	Title  string  `json:"title"`
	Author string  `json:"author"`
	Price  float64 `json:"price"`
}
