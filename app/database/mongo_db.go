package database

import (
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDatabase bridge mongo database
type MongoDatabase struct {
	db *mongo.Database
}

// Collection get mongo collection
func (db *MongoDatabase) Collection(name string, opts ...*options.CollectionOptions) Collection {
	collection := db.db.Collection(name, opts...)
	return &MongoCollection{internal: collection}
}
