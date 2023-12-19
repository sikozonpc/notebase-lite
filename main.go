package main

import (
	"log"

	"github.com/go-sql-driver/mysql"
)

func main() {
	cfg := mysql.Config{
		User:                 Configs.DBUser,
		Passwd:               Configs.DBPassword,
		Net:                  "tcp",
		Addr:                 Configs.DBAddress,
		DBName:               Configs.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
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
