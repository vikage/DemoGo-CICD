package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoSessionContext session context bridge
type MongoSessionContext struct {
	client MongoClient
	ctx    mongo.SessionContext
}

// Client get mongo client
func (context *MongoSessionContext) Client() Client {
	return &context.client
}

// CommitTransaction commit mongo transaction
func (context *MongoSessionContext) CommitTransaction(ctx context.Context) error {
	return context.ctx.CommitTransaction(ctx)
}

// AbortTransaction abort mongo transaction
func (context *MongoSessionContext) AbortTransaction(ctx context.Context) error {
	return context.ctx.AbortTransaction(ctx)
}

// StartTransaction start a transaction
func (context *MongoSessionContext) StartTransaction(opts ...*options.TransactionOptions) error {
	return context.ctx.StartTransaction(opts...)
}

// Context get context
func (context *MongoSessionContext) Context() context.Context {
	return context.ctx
}
