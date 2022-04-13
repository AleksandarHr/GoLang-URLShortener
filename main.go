package main

import (
	// "os"
	"fmt"
	"log"
	"net/http"
	server "url_shortener/server"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("Starting Chaos URL Shortener.")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Start server
	s := server.NewUrlShortenerServer()
	if err := s.StartServer(); err != http.ErrServerClosed {
		// unexpected error, close db and log error as fatal
		s.ShortenerDb.CloseDB()
		log.Fatalf("ListenAndServe(): %v", err)
	}
}
