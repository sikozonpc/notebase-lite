package user

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

func (s *Store) CreateUser(user t.User) error {
	_, err := s.db.Exec("INSERT INTO users (firstName, lastName, email, password) VALUES (?, ?, ?, ?)", user.FirstName, user.LastName, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) GetUserByEmail(email string) (*t.User, error) {
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

func (s *Store) GetUserByID(id string) (*t.User, error) {
	rows, err := s.db.Query("SELECT * FROM users WHERE id = ?", id)
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

func scanRowsIntoUser(rows *sql.Rows) (*t.User, error) {
	user := new(t.User)

	err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
