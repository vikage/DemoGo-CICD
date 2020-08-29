package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database database interface
type Database interface {
	Collection(name string, opts ...*options.CollectionOptions) Collection
}

// Client db client interface
type Client interface {
	Ping(ctx context.Context) error
	Database(name string, opts ...*options.DatabaseOptions) Database
	UseSession(ctx context.Context, fn func(SessionContext) error) error
}

// TransactionClient db transation client interface
type TransactionClient interface {
	Client
}

// SessionContext context bridge
type SessionContext interface {
	Context() context.Context
	Client() Client
	StartTransaction(opts ...*options.TransactionOptions) error
	CommitTransaction(ctx context.Context) error
	AbortTransaction(ctx context.Context) error
}

// InsertOneResult insert one result interface
type InsertOneResult interface {
	InsertedID() interface{}
}

// UpdateResult update result bridge interface
type UpdateResult interface {
	MatchedCount() int64  // The number of documents matched by the filter.
	ModifiedCount() int64 // The number of documents modified by the operation.
	UpsertedCount() int64 // The number of documents upserted by the operation.
	UpsertedID() interface{}
}

// SingleResult single result bridge interface
type SingleResult interface {
	Decode(v interface{}) error
}

// DeleteResult delete result interface
type DeleteResult interface {
	DeletedCount() int64
}

// Cursor cursor interface
type Cursor interface {
	Next(ctx context.Context) bool
	Decode(val interface{}) error
}

// Collection collection interface
type Collection interface {
	InsertOne(ctx context.Context, document interface{}, opts ...*options.InsertOneOptions) (InsertOneResult, error)
	ReplaceOne(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.ReplaceOptions) (UpdateResult, error)
	FindOneAndUpdate(ctx context.Context, filter interface{}, update interface{}, opts ...*options.FindOneAndUpdateOptions) SingleResult
	FindOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) SingleResult
	Find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Cursor, error)
	UpdateOne(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error)
	DeleteOne(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error)
	DeleteMany(ctx context.Context, filter interface{}, opts ...*options.DeleteOptions) (DeleteResult, error)
	FindOneAndDelete(ctx context.Context, filter interface{}, opts ...*options.FindOneAndDeleteOptions) SingleResult
	FindOneAndReplace(ctx context.Context, filter interface{}, replacement interface{}, opts ...*options.FindOneAndReplaceOptions) SingleResult
	UpdateMany(ctx context.Context, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (UpdateResult, error)
	Indexes() mongo.IndexView
	CountDocuments(ctx context.Context, filter interface{}, opts ...*options.CountOptions) (int64, error)
}
