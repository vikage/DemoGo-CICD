package database

import "go.mongodb.org/mongo-driver/mongo"

// MongoDeleteResult delete result bridge
type MongoDeleteResult struct {
	internal *mongo.DeleteResult
}

// DeletedCount get deleted count
func (result *MongoDeleteResult) DeletedCount() int64 {
	return result.internal.DeletedCount
}
