package main

import (
	"fmt"
	"log"

	"github.com/sikozonpc/notebase/config"
	"github.com/sikozonpc/notebase/db"
)

func main() {
	mongoClient, err := db.ConnectToMongo(config.Envs.MongoURI)
	if err != nil {
		log.Fatal(err)
	}

	server := NewAPIServer(fmt.Sprintf(":%s", config.Envs.Port), mongoClient)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}
