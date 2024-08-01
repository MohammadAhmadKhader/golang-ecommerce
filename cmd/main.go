package main

import (
	"database/sql"
	"log"

	"github.com/mohammadahmadkhader/golang-ecommerce/cmd/api"
	"github.com/mohammadahmadkhader/golang-ecommerce/utils"
)

func main() {
	db ,err := utils.StartMySqlDB()
	if err != nil {
		log.Fatal(err)
	}
	
	err = initStorage(db)
	if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", db)

	err = server.Run()
	if err != nil {
		log.Fatal(err)
	}
	
	defer db.Close()
}


func initStorage(db *sql.DB) error {
	err := db.Ping()
	if err != nil {
		return err
	}

	log.Println("DB is successfully connected")
	return nil
}