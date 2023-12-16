package main

import (
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:   "root",
		Passwd: "mypassword",
		Net:    "tcp",
		Addr:   "127.0.0.1:3306",
		AllowNativePasswords: true,
		DBName: "highlights",
		ParseTime: true,
	}

	store, err := NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if err := store.Init(); err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(":3000", store)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
