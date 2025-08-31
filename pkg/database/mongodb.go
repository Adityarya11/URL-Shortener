package database

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectMongo() *mongo.Database {
	mongoURI := os.Getenv("MONGO_URI")
	dbName := os.Getenv("MONGO_DB")

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Failed to create Mongo client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("❌ Failed to connect to Mongo: %v", err)
	}

	// Test ping
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatalf("❌ MongoDB ping failed: %v", err)
	}

	log.Println("✅ Connected to MongoDB Atlas")

	return client.Database(dbName)
}
