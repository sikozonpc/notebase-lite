package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
)

type MySQLStorage struct {
	db *sql.DB
}

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

func (s *MySQLStorage) Init() (*sql.DB, error) {
	// Initialize tables
	err := s.createHighlightsTable()
	if err != nil {
		return nil, err
	}
	err = s.createUsersTable()
	if err != nil {
		return nil, err
	}

	err = s.createBookTable()
	if err != nil {
		return nil, err
	}

	return s.db, nil
}

func (s *MySQLStorage) createHighlightsTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS highlights (
			id INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
			text TEXT,
			location VARCHAR(255),
			note TEXT,
			userId INT UNSIGNED NOT NULL,
			bookId VARCHAR(50) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updatedAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			
			FOREIGN KEY (userId) REFERENCES users(id),
			FOREIGN KEY (bookId) REFERENCES books(isbn),
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
			isActive BOOLEAN NOT NULL DEFAULT TRUE,
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

func (s *MySQLStorage) createBookTable() error {
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS books (
			isbn VARCHAR(50) NOT NULL,
			title VARCHAR(500) NOT NULL,
			authors VARCHAR(500) NOT NULL,
			createdAt TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			
			PRIMARY KEY (isbn),
			UNIQUE KEY (isbn),
			INDEX (isbn)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8;
	`)
	if err != nil {
		return err
	}

	return nil
}
