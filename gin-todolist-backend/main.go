// main.go
package main

import (
	"log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
  if err != nil {
    log.Fatal("Error loading .env file")
  }

	dbGetAllTasks()
	router := Router()
	router.Run()
	
}