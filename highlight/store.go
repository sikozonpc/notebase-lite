package highlight

import (
	"database/sql"

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
