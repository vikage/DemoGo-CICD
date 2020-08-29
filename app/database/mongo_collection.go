package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoCollection bridge mongo collection
type MongoCollection struct {
	internal *mongo.Collection
}

// InsertOne insert one document to mongo
func (collection *MongoCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (InsertOneResult, error) {
	result, err := collection.internal.InsertOne(ctx, document, opts...)
	r1 := &MongoInsertOneResult{
		internal: result,
	}

	return r1, err
}

// ReplaceOne replace one doc in mongo
func (collection *MongoCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (UpdateResult, error) {
	result, err := collection.internal.ReplaceOne(ctx, filter, replacement, opts...)
	return &MongoUpdateResult{internal: result}, err
}

// FindOneAndUpdate find and update document in mongo
func (collection *MongoCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) SingleResult {
	result := collection.internal.FindOneAndUpdate(ctx, filter, update, opts...)
	return &MongoSingleResult{internal: result}
}

// FindOne find one document in mongo
func (collection *MongoCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult {
	result := collection.internal.FindOne(ctx, filter, opts...)
	return &MongoSingleResult{internal: result}
}

// Find find many document in mongo
func (collection *MongoCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	cursor, err := collection.internal.Find(ctx, filter, opts...)
	return &MongoCursor{internal: cursor}, err
}

// UpdateOne update one document in mongo
func (collection *MongoCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error) {
	result, err := collection.internal.UpdateOne(ctx, filter, update, opts...)
	return &MongoUpdateResult{internal: result}, err
}

// DeleteOne delete one document in mongo
func (collection *MongoCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error) {
	result, err := collection.internal.DeleteOne(ctx, filter, opts...)
	return &MongoDeleteResult{internal: result}, err
}

// DeleteMany delete many document in mongo
func (collection *MongoCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error) {
	result, err := collection.internal.DeleteMany(ctx, filter, opts...)
	return &MongoDeleteResult{internal: result}, err
}

// FindOneAndDelete find one and delete document in mongo. Return document that found
func (collection *MongoCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) SingleResult {
	result := collection.internal.FindOneAndDelete(ctx, filter, opts...)
	return &MongoSingleResult{internal: result}
}

// FindOneAndReplace find one and replace document in mongo. Return document that found
func (collection *MongoCollection) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) SingleResult {
	result := collection.internal.FindOneAndReplace(ctx, filter, replacement, opts...)
	return &MongoSingleResult{internal: result}
}

// UpdateMany update many document in mongo
func (collection *MongoCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error) {
	result, err := collection.internal.UpdateMany(ctx, filter, update, opts...)
	return &MongoUpdateResult{internal: result}, err
}

// Indexes get indexes in collection
func (collection *MongoCollection) Indexes() mongo.IndexView {
	return collection.internal.Indexes()
}

// CountDocuments count documents in collection by filter
func (collection *MongoCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return collection.internal.CountDocuments(ctx, filter, opts...)
}
