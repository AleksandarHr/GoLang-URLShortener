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
	return client
}

func GetDbCollection(client *mongo.Client) *mongo.Collection {
	dbName := os.Getenv("MONGODB_DB_NAME")
	collectionName := os.Getenv("MONGODB_COLLECTION_NAME")
	return client.Database(dbName).Collection(collectionName)
}

func CloseDB(client *mongo.Client) {
	err := client.Disconnect(context.TODO())

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Successfully disconnected from MongoDB!")
}

func InsertURL(collection *mongo.Collection, url string) *DbEntry {
	shortUrl := algorithms.ShortenURL(url)
	dbEntry := DbEntry{url, shortUrl}
	fmt.Println(dbEntry)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if exists, existingShortUrl := previouslyShortened(collection, url); exists {
		fmt.Println("Record for provided long URL already exists.")
		return existingShortUrl
	}

	insertResult, err := collection.InsertOne(ctx, dbEntry)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)
	return &dbEntry
}

func QueryShortURL(collection *mongo.Collection, shortUrlQuery string) string {
	var queryResult DbEntry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := collection.FindOne(ctx, bson.D{{"shorturl", shortUrlQuery}}).Decode(&queryResult)
	if err != nil {
		log.Fatal(err)
	}

	return queryResult.OriginalUrl
}

func previouslyShortened(collection *mongo.Collection, originalUrl string) (bool, *DbEntry) {
	var queryResult DbEntry

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	fmt.Println("Checking for existance of " + originalUrl)
	err := collection.FindOne(ctx, bson.D{{"originalurl", originalUrl}}).Decode(&queryResult)
	if err != nil {
		// The filter did not match any documents in the collection
		if err == mongo.ErrNoDocuments {
			fmt.Println("No record exists for " + originalUrl)
			return false, nil
		}
		log.Fatal(err)
	}
	return true, &queryResult
}
