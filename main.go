package main

import (
	// "os"
	"fmt"
	server "url_shortener/server"
)


func main() {
	fmt.Println("Starting Chaos URL Shortener.")
	// port := os.Getenv("SERVER_PORT")
	// fmt.Println("Server listening on port " + port)
	// Start server
	server.InitServer()

	fmt.Println("Closing Chaos URL Shortener.")
}