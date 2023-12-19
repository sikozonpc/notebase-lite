package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	t "github.com/sikozonpc/notebase/types"
)

func NewMySQLStorage(cfg mysql.Config) (*MySQLStorage, error) {
	var err error
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	return &MySQLStorage{db: db}, nil
}

func (s *MySQLStorage) Init() error {
	// Initialize tables
	err := s.createHighlightsTable()
	if err != nil {
		return err
	}
	err = s.createUsersTable()
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStorage) createHighlightsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS highlights (
			id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
			text TEXT,
			location VARCHAR(255),
			note TEXT,
			userId INT UNSIGNED NOT NULL,
			bookId INT UNSIGNED NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStorage) createUsersTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
			firstName VARCHAR(255) NOT NULL,
			lastName VARCHAR(255) NOT NULL,
			email VARCHAR(500) NOT NULL,
			password VARCHAR(255) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (id),
			UNIQUE KEY (email),
			INDEX (email)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)
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
		&highlight.UserId,
		&highlight.BookId,
		&highlight.CreatedAt,
		&highlight.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return highlight, nil
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
