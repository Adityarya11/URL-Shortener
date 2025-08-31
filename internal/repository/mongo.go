package repository

import (
	"context"
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
	return err
}

func (r *MongoRepo) Find(shortCode string) (*models.URL, error) {
	var url models.URL
	err := r.collection.FindOne(context.Background(), bson.M{"shortcode": shortCode}).Decode(&url)
	if err != nil {
		return nil, err
	}
	return &url, nil
}

func (r *MongoRepo) IncrementClicks(shortCode string) error {
	_, err := r.collection.UpdateOne(
		context.Background(),
		bson.M{"shortcode": shortCode},
		bson.M{"$inc": bson.M{"clickcount": 1}},
	)
	return err
}
