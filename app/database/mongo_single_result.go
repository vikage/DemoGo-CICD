package database

import "go.mongodb.org/mongo-driver/mongo"

// MongoSingleResult single result bridge
type MongoSingleResult struct {
	internal *mongo.SingleResult
}

// Decode decode struct from result
func (result *MongoSingleResult) Decode(v interface{}) error {
	return result.internal.Decode(v)
}
