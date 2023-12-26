package main

import (
	"fmt"
	"log"

	"github.com/go-sql-driver/mysql"
	"github.com/sikozonpc/notebase/config"
	"github.com/sikozonpc/notebase/db"
)

func main() {
	cfg := mysql.Config{
		User:                 config.Envs.DBUser,
		Passwd:               config.Envs.DBPassword,
		Net:                  "tcp",
		Addr:                 config.Envs.DBAddress,
		DBName:               config.Envs.DBName,
		AllowNativePasswords: true,
		ParseTime:            true,
	}

	store, err := db.NewMySQLStorage(cfg)
	if err != nil {
		log.Fatal(err)
	}

	db, err := store.Init()
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	server := NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), db)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
