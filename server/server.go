package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"url_shortener/db"

	"github.com/gorilla/mux"
)

var router = mux.NewRouter()
var Client = db.InitDB()
var Collection = db.GetDbCollection(Client)

func InitServer() int {
	router.HandleFunc("/{url}", GetShortURL)
	router.HandleFunc("/", PostURL)
	http.Handle("/", router)

	port := ":" + os.Getenv("SERVER_PORT")
	if err := http.ListenAndServe(port, nil); err != nil {
		// TODO: Add error handling/reporting
		return -1
	}

	return 0
}

func GetShortURL(w http.ResponseWriter, r *http.Request) {
	// TODO: Read the requested shortened URL from the vars
	// TODO: Query the db
	// TODO: Redirect to the original URL
	shortUrl := mux.Vars(r)["url"]
	// vars contains the shortened url to be lookedup and redirected
	originalUrl := db.QueryShortURL(Collection, shortUrl)
	http.Redirect(w, r, originalUrl, http.StatusSeeOther)
}

func PostURL(w http.ResponseWriter, r *http.Request) {
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	// responseString contains the requested URL to be shortened
	url := string(bodyBytes)
	fmt.Println(url)
	// TODO: shorten the URL
	// 1. Check if the long URL already exists in db. If so, return its short counterpart
	// 2. If not, run the shortening algorithm
	// 3. Store in db
	result := db.InsertURL(Collection, url)
	// 4. Return the shortened URL
	fmt.Fprintf(w, "Shortened URL: %v\n", result.ShortUrl)
}
