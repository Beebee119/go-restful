package main

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// Book is a placeholder for book
type Book struct {
	id     int
	name   string
	author string
}

func main() {
	db, err := sql.Open("sqlite3", "./books.db")
	log.Println(db)
	if err != nil {
		log.Println(err)
	}
	// Create table
	q := `CREATE TABLE IF NOT EXISTS books (
        id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
        "isbn" INTEGER,
        "author" VARCHAR(64),
        "name" VARCHAR(64));`
	statement, err := db.Prepare(q)
	if err != nil {
		log.Println("Error in creating table books")
	} else {
		log.Println("Successfully created table books")
	}
	statement.Exec()

	// Create
	statement, _ = db.Prepare("INSERT INTO books (name, author, isbn) VALUES (?, ?, ?)")
	statement.Exec("A tale of two cities", "Charles Young", 1404567)
	log.Println("Book innserted into the database")

	// Read
	books, _ := db.Query("SELECT id, name, author FROM books")
	var tempBook Book
	for books.Next() {
		books.Scan(&tempBook.id, &tempBook.name, &tempBook.author)
		log.Printf("ID:%d, Books:%s, Author:%s\n", tempBook.id, tempBook.name, tempBook.author)
	}
}
