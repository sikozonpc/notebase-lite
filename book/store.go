package book

import (
	"database/sql"
	"fmt"

	t "github.com/sikozonpc/notebase/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetBookByISBN(ISBN string) (*t.Book, error) {
	rows, err := s.db.Query("SELECT * FROM books WHERE isbn = ?", ISBN)
	if err != nil {
		return nil, err
	}

	book := new(t.Book)
	for rows.Next() {
		book, err = scanRowsIntoBook(rows)
		if err != nil {
			return nil, err
		}
	}

	if book.ISBN == "" {
		return nil, fmt.Errorf("book not found")
	}

	return book, nil
}

func (s *Store) CreateBook(book t.Book) error {
	_, err := s.db.Exec("INSERT INTO books (isbn, title, authors) VALUES (?, ?, ?)", book.ISBN, book.Title, book.Authors)
	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoBook(rows *sql.Rows) (*t.Book, error) {
	b := new(t.Book)
	err := rows.Scan(&b.ISBN, &b.Title, &b.Authors)
	if err != nil {
		return nil, err
	}

	return b, nil
}
