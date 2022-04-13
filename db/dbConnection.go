// Package db provides functionality for initializing, querying, and inserting entries to mongo db.
package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"url_shortener/algorithms"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DbEntry is a struct containing an original url and the corresponding shortened url.
type DbEntry struct {
	OriginalUrl string
	ShortUrl    string
}

// UrlShortenerDb is a struct encapsulating a mongo db Client and a mongo db Collection.
type UrlShortenerDb struct {
	Client     *mongo.Client
	Collection *mongo.Collection
}

// InitDB initializes a connection with mongo db using the URI provided in the .env file.
// It pings the db after initiating the connection to make sure it has successfully connected.
// It returns the created mongo Client.
func InitDB() *UrlShortenerDb {
	// Load mongodb URI from .env file
	err := godotenv.Load()
	mongoURI := os.Getenv("MONGODB_URI")

	// Set connection options
	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPIOptions)

	// Add connection timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Try to connect to mongo db
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check if a MongoDB server has been found and connected to
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Successfully connected to MongoDB!")

	urlShortenerDb := &UrlShortenerDb{Client: client}
	urlShortenerDb.GetDbCollection()
	return urlShortenerDb
}

// GetDbCollection gets the relevant collection from the relevant database in mongo.
// It returns a pointer to the retrieved Collection.
func (db *UrlShortenerDb) GetDbCollection() {
	dbName := os.Getenv("MONGODB_DB_NAME")
	collectionName := os.Getenv("MONGODB_COLLECTION_NAME")
	db.Collection = db.Client.Database(dbName).Collection(collectionName)
}

// CloseDB closes the existing connection to mongo db.
func (db *UrlShortenerDb) CloseDB() {
	err := db.Client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully disconnected from MongoDB!")
}

// InsertURL takes a URL string and checks the database if it was previously shortened.
// If yes, then returns the corresponding short URL. Otherwise, it performs the shortening
// algorithm, stores the new entry in the database, and returns a pointer to a DbEntry with
// the original URL and the corresponding shortened URL.
func (db *UrlShortenerDb) InsertURL(url string) *DbEntry {

	// Check the db if the provided URL has already been shortened before.
	if exists, existingShortUrl := db.previouslyShortened(url); exists {
		return existingShortUrl
	}

	// If not, perform shortening algorithm.
	shortUrl := algorithms.ShortenURL(url)
	dbEntry := DbEntry{url, shortUrl}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Insert the record with the URL pair in the db
	insertResult, err := db.Collection.InsertOne(ctx, dbEntry)
	if err != nil {
		log.Println(err)
		return nil
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return &dbEntry
}

// QueryShortURL takes a shortened URL string and queries the database for its corresponding
// original long querry. If a record corresponding to the provided short URL is found, it
// returns a pointer to a DbEntry with the original long and the shortened URL. Otherwise, it
// logs an error.
func (db *UrlShortenerDb) QueryShortURL(shortUrlQuery string) *DbEntry {
	var queryResult DbEntry

	// Create a 5-second timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Querry the database for a record corresponding to the provided shortened URL
	err := db.Collection.FindOne(ctx, bson.D{{"shorturl", shortUrlQuery}}).Decode(&queryResult)
	if err != nil {
		log.Println(err)
		return nil
	}

	return &queryResult
}

// previouslyShortened checks if a provided long URL has been previously shortened.
// It querries the database to see if a record exists for the provided long URL.
// If so, it returns a pointer to a DbEntry with the relevant long and shortened URLs.
func (db *UrlShortenerDb) previouslyShortened(originalUrl string) (bool, *DbEntry) {
	var queryResult DbEntry

	// Create a 5-second timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Querry the database for a record corresponding to the provided long URL
	err := db.Collection.FindOne(ctx, bson.D{{"originalurl", originalUrl}}).Decode(&queryResult)
	if err != nil {
		// The filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			return false, nil
		}
		log.Println(err)
		return false, nil
	}

	// A record already exists in the database for the long URL. Return corresponding pair of URLs
	return true, &queryResult
}
