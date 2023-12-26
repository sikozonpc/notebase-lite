package highlight

import (
	"database/sql"
	"log"

	t "github.com/sikozonpc/notebase/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserHighlights(userID int) ([]*t.Highlight, error) {
	rows, err := s.db.Query("SELECT * FROM highlights WHERE userId = ?", userID)
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

func (s *Store) CreateHighlight(highlight t.Highlight) error {
	_, err := s.db.Exec("INSERT INTO highlights (text, location, note, userId, bookId) VALUES (?, ?, ?, ?, ?)", highlight.Text, highlight.Location, highlight.Note, highlight.UserID, highlight.BookID)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) CreateHighlights(highlights []t.Highlight) error {
	tx, err := s.db.Begin()
	if err != nil {
		log.Println("Error starting transaction: ", err)
		return err
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// Create a slice to hold the values
	values := []interface{}{}

	query := "INSERT INTO highlights (text, location, note, userId, bookId) VALUES "
	for _, h := range highlights {
		query += "(?, ?, ?, ?, ?),"
		values = append(values, h.Text, h.Location, h.Note, h.UserID, h.BookID)
	}

	// Remove the last comma
	query = query[:len(query)-1]

	_, err = tx.Exec(query, values...)
	if err != nil {
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Store) GetHighlightByID(id, userID int) (*t.Highlight, error) {
	rows, err := s.db.Query("SELECT * FROM highlights WHERE id = ? AND userId = ?", id, userID)
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

func (s *Store) DeleteHighlight(id int) error {
	_, err := s.db.Exec("DELETE FROM highlights WHERE id = ?", id)
	if err != nil {
		return err
	}

	return nil
}

func scanRowsIntoHighlight(rows *sql.Rows) (*t.Highlight, error) {
	highlight := new(t.Highlight)

	err := rows.Scan(
		&highlight.ID,
		&highlight.Text,
		&highlight.Location,
		&highlight.Note,
		&highlight.UserID,
		&highlight.BookID,
		&highlight.CreatedAt,
		&highlight.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return highlight, nil
}
