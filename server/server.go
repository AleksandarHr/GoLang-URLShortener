package server

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	urlParser "net/url"
	"os"
	"url_shortener/db"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var Client = db.InitDB()
var Collection = db.GetDbCollection(Client)

func InitServer() {
	router.HandleFunc("/{url}", GetShortURL)
	router.HandleFunc("/", PostURL)
	http.Handle("/", router)

	port := ":" + os.Getenv("SERVER_PORT")
	log.Fatal(http.ListenAndServe(port, nil))
}

func GetShortURL(w http.ResponseWriter, r *http.Request) {
	// TODO: Read the requested shortened URL from the vars
	// TODO: Query the db
	// TODO: Redirect to the original URL
	shortUrl := mux.Vars(r)["url"]
	// vars contains the shortened url to be lookedup and redirected
	queryResult := db.QueryShortURL(Collection, shortUrl)
	if queryResult != nil {
		http.Redirect(w, r, queryResult.OriginalUrl, http.StatusSeeOther)
	}
}

func PostURL(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	// responseString contains the requested URL to be shortened
	url := string(bodyBytes)
	// Verify 'url' contains a valid url
	_, err := urlParser.ParseRequestURI(url)
	if err != nil {
		log.Println(err)
		return
	}
	// TODO: shorten the URL
	// 1. Check if the long URL already exists in db. If so, return its short counterpart
	// 2. If not, run the shortening algorithm
	// 3. Store in db
	insertionResult := db.InsertURL(Collection, url)
	if insertionResult != nil {
		// 4. Print the shortened URL
		fmt.Fprintf(w, "Shortened URL: %v\n", insertionResult.ShortUrl)
	}
}
