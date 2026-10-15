package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"web/internal/service"
	"web/internal/store"
	"web/internal/transport"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// TODO:  conectar a sqlite

	db, err := sql.Open("sqlite3", "./books.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// TODO: crear table si no existe
	q := `
		CREATE TABLE IF NOT EXISTS books (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				title TEXT NOT NULL,
				author TEXT NOT NULL,
				price INTEGER NOT NULL
		)
	`
	if _, err := db.Exec(q); err != nil {
		log.Fatal(err.Error())
	}

	// TODO: Injectar las dependencias

	bookStore := store.New(db)
	bookService := service.NewService(bookStore)
	bookHandler := transport.NewBookHandler(bookService)

	// TODO: Configurar las rutas

	http.HandleFunc("/books", bookHandler.HandleBooks)
	http.HandleFunc("/books/", bookHandler.HandleBookById)

	fmt.Println("Serviro ejecutandose en http://localhost:8080")

	// TODO: Empezar y escuchar al servidor
	log.Fatal(http.ListenAndServe(":8080", nil))
}
