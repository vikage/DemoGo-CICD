package database

import "go.mongodb.org/mongo-driver/mongo"

// MongoUpdateResult update result bridge
type MongoUpdateResult struct {
	internal *mongo.UpdateResult
}

// MatchedCount get matched count
func (result *MongoUpdateResult) MatchedCount() int64 {
	return result.internal.MatchedCount
}

// ModifiedCount get modified count
func (result *MongoUpdateResult) ModifiedCount() int64 {
	return result.internal.ModifiedCount
}

// UpsertedCount get upserted count
func (result *MongoUpdateResult) UpsertedCount() int64 {
	return result.internal.UpsertedCount
}

// UpsertedID get upserted id
func (result *MongoUpdateResult) UpsertedID() interface{} {
	return result.internal.UpsertedID
}
