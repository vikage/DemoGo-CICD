package database

import "go.mongodb.org/mongo-driver/mongo"

// MongoInsertOneResult insert one result bridge
type MongoInsertOneResult struct {
	internal *mongo.InsertOneResult
}

// InsertedID get inserted id
func (result *MongoInsertOneResult) InsertedID() interface{} {
	return result.internal.InsertedID
}
