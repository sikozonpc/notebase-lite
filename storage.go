package main

import (
	"database/sql"

	t "github.com/sikozonpc/notebase/types"
)

type Storage interface {
	GetHighlights() ([]*t.Highlight, error)
	GetHighlightByID(id int) (*t.Highlight, error)
	CreateHighlight(t.Highlight) error
	DeleteHighlight(id int) error

	CreateUser(t.User) error
	GetUserByEmail(email string) (*t.User, error)
}

type MySQLStorage struct {
	db *sql.DB
}

func (s *MySQLStorage) GetHighlights() ([]*t.Highlight, error) {
	rows, err := s.db.Query("SELECT * FROM highlights")
	if err != nil {
		return nil, err
	}

	var highlights []*t.Highlight
	for rows.Next() {
		h, err := scanRowsIntoHighlight(rows)
		if err != nil {
			return nil, err
		}

		highlights = append(highlights, h)
	}

	return highlights, nil
}

func (s *MySQLStorage) CreateHighlight(highlight t.Highlight) error {
	_, err := s.db.Exec("INSERT INTO highlights (text, location, note, userId, bookId) VALUES (?, ?, ?, ?, ?)", highlight.Text, highlight.Location, highlight.Note, highlight.UserId, highlight.BookId)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStorage) GetHighlightByID(id int) (*t.Highlight, error) {
	rows, err := s.db.Query("SELECT * FROM highlights WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	h := new(t.Highlight)
	for rows.Next() {
		h, err = scanRowsIntoHighlight(rows)
		if err != nil {
			return nil, err
		}
	}

	return h, nil
}

func (s *MySQLStorage) DeleteHighlight(id int) error {
	_, err := s.db.Exec("DELETE FROM highlights WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStorage) CreateUser(user t.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStorage) GetUserByEmail(email string) (*t.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE email = ?", email)
	if err != nil {
		return nil, err
	}

	u := new(t.User)
	for rows.Next() {
		u, err = scanRowsIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	return u, nil
}
