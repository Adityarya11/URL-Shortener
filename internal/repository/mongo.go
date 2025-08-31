package repository

import (
	"context"
	"log"
	"url-shortener/internal/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoRepo struct {
	collection *mongo.Collection
}

func NewMongoRepo(db *mongo.Database, collectionName string) *MongoRepo {
	return &MongoRepo{
		collection: db.Collection(collectionName),
	}
}

func (r *MongoRepo) Save(url *models.URL) error {
	_, err := r.collection.InsertOne(context.Background(), url)
	if err != nil {
		return err
	}
	log.Printf("✅ Saved to Mongo: %s -> %s\n", url.ShortCode, url.OriginalURL)
	return nil
}

func (r *MongoRepo) Find(shortCode string) (*models.URL, error) {
	var url models.URL
	err := r.collection.FindOne(context.Background(), bson.M{"shortCode": shortCode}).Decode(&url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *MongoRepo) IncrementClicks(shortCode string) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"shortCode": shortCode},
		bson.M{"$inc": bson.M{"clickCount": 1}},
	)
	return err
}
