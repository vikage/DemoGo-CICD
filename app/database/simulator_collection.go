package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SimulatorCollection bridge simulator collection
type SimulatorCollection struct {
	name    string
	storage *SimulatorCollectionStorage
}

// InsertOne insert one document to mongo
func (collection *SimulatorCollection) InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (InsertOneResult, error) {
	return collection.storage.Insert(document)
}

// ReplaceOne replace one doc in mongo
func (collection *SimulatorCollection) ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (UpdateResult, error) {
	return collection.storage.ReplaceOne(filter, replacement, opts...)
}

// FindOneAndUpdate find and update document in mongo
func (collection *SimulatorCollection) FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) SingleResult {
	return collection.storage.FindOneAndUpdate(filter, update, opts...)
}

// FindOne find one document in mongo
func (collection *SimulatorCollection) FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult {
	return collection.storage.FindOne(filter, opts...)
}

// Find find many document in mongo
func (collection *SimulatorCollection) Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error) {
	return collection.storage.Find(filter)
}

// UpdateOne update one document in mongo
func (collection *SimulatorCollection) UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error) {
	return collection.storage.UpdateOne(filter, update, opts...)
}

// DeleteOne delete one document in mongo
func (collection *SimulatorCollection) DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error) {
	return collection.storage.DeleteOne(filter, opts...)
}

// DeleteMany delete many document in mongo
func (collection *SimulatorCollection) DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error) {
	return collection.storage.DeleteMany(filter, opts...)
}

// FindOneAndDelete find one and delete document in mongo. Return document that found
func (collection *SimulatorCollection) FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) SingleResult {
	return collection.storage.FindOneAndDelete(filter, opts...)
}

// FindOneAndReplace find one and replace document in mongo. Return document that found
func (collection *SimulatorCollection) FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) SingleResult {
	return collection.storage.FindOneAndReplace(filter, replacement, opts...)
}

// UpdateMany update many document in mongo
func (collection *SimulatorCollection) UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error) {
	return collection.storage.UpdateMany(filter, update, opts...)
}

// Indexes get indexes in collection
func (collection *SimulatorCollection) Indexes() mongo.IndexView {
	// return collection.internal.Indexes()
	return mongo.IndexView{}
}

// CountDocuments count documents in simulator database
func (collection *SimulatorCollection) CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	return collection.storage.CountDocuments(filter, opts...)
}
