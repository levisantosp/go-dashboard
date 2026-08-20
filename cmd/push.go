package main

import (
	"context"
	"log"

	"dash/db"
	"dash/utils"
)

func main() {
	utils.LoadEnv()
	db.Connect()

	err := db.Client.Schema.Create(context.Background())
	if err != nil {
		log.Fatal(err)
	}
}
