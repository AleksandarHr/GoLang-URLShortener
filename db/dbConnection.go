package db

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"
	"url_shortener/algorithms"

	"github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DbEntry struct {
	OriginalUrl string
	ShortUrl    string
}

func InitDB() *mongo.Client {
	err := godotenv.Load()
	mongoURI := os.Getenv("MONGODB_URI")

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetServerAPIOptions(serverAPIOptions)
	// Add connection timeout of 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
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

	// Initialize mongo db client
	// client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// // Calling Connect does not block for server discovery.
	// err = client.Connect(ctx)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Check if a MongoDB server has been found and connected to
	// err = client.Ping(ctx, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	return client
}

func GetDbCollection(client *mongo.Client) *mongo.Collection {
	dbName := os.Getenv("MONGODB_DB_NAME")
	collectionName := os.Getenv("MONGODB_COLLECTION_NAME")
	return client.Database(dbName).Collection(collectionName)
}

func CloseDB() {

}

func InsertURL(collection *mongo.Collection, url string) DbEntry {
	shortUrl := algorithms.ShortenURL(url)
	dbEntry := DbEntry{url, shortUrl}
	fmt.Println(dbEntry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertResult, err := collection.InsertOne(ctx, dbEntry)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return dbEntry
}

func QueryShortURL() {

}

func QueryLongURL() {

}
