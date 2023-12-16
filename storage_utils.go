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
	err := s.createHighlightsTable()
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