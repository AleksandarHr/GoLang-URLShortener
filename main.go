package main

import (
	// "os"
	"fmt"
	"log"
	server "url_shortener/server"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Chaos URL Shortener.")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// db.InitDB()
	// Start server
	server.InitServer()

	fmt.Println("Closing Chaos URL Shortener.")
}
