package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	urlParser "net/url"
	"os"
	db "url_shortener/db"

	"github.com/gorilla/mux"
)

// UrlShortenerServer is a struct encapsulating a mux Router and a db UrlShortenerDb.
type UrlShortenerServer struct {
	Router      *mux.Router
	ShortenerDb *db.UrlShortenerDb
}

// NewUrlShortenerServer initializes a UrlShortenerServer
func NewUrlShortenerServer() *UrlShortenerServer {
	return &UrlShortenerServer{
		Router:      mux.NewRouter(),
		ShortenerDb: db.InitDB(),
	}
}

// StartServer sets up handlers and fires up a server listening on
// the port specified in the .env file.
func (s *UrlShortenerServer) StartServer() error {
	// Set up handlers
	s.Router.HandleFunc("/{url}", s.GetShortURL)
	s.Router.HandleFunc("/", s.PostURL)
	http.Handle("/", s.Router)

	// Start listening on the specified port
	port := ":" + os.Getenv("SERVER_PORT")
	err := http.ListenAndServe(port, nil)
	return err
}

// GetShortURL performs a Query operation to the db for the provided short URL.
// If a result is provided, it redirects to the corresponding original URL and with a
// 'See Other' status code. Otherwise, it prints a message.
func (s *UrlShortenerServer) GetShortURL(w http.ResponseWriter, r *http.Request) {
	shortUrl := mux.Vars(r)["url"]

	// Shortened URLs are always of length 5
	if len(shortUrl) != 5 {
		log.Println("Supports only short URLs of length 5")
	}
	// TODO: check if provided short URL is base58

	// Query the db for the corresponding original long URL
	queryResult := s.ShortenerDb.QueryShortURL(shortUrl)
	if queryResult == nil {
		// If not found, print message
		log.Println("No corresponding URL found for the short: " + shortUrl)
		return
	}
	// If found, redirect
	http.Redirect(w, r, queryResult.OriginalUrl, http.StatusSeeOther)
}

// PostURL performs an Insert operation to the db with the provided long URL.
// If it was successful, the corresponding short URL is printed. Otherwise,
// an information message is printed.
func (s *UrlShortenerServer) PostURL(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	url := string(bodyBytes)

	// Verify 'url' contains a valid url
	_, err := urlParser.ParseRequestURI(url)
	if err != nil {
		log.Println(err)
		return
	}

	// Perform Insert operation with the provided long URL
	insertionResult := s.ShortenerDb.InsertURL(url)
	if insertionResult == nil {
		log.Println("Unsuccessful insert operation for long url: " + url)
		return
	}
	fmt.Fprintf(w, "Shortened URL: %v\n", insertionResult.ShortUrl)
}
