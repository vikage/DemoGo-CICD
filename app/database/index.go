package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

// EnsureIndexes create indexes if need
func EnsureIndexes() {
	dbClient := ResolveDatabaseClient()
	db := dbClient.Database("go-cicd")

	ctx, cancel := MongoTimeoutContext()
	defer cancel()

	EnsureIndex(ctx, db.Collection("User"), []string{"email"}, options.Index().SetBackground(true).SetName("email"))
}

// EnsureIndex create index if need
func EnsureIndex(ctx context.Context, c Collection, keys []string, opts *options.IndexOptions) error {
	// get indices
	batchSize := int32(10)
	duration := 10 * time.Second
	cur, err := c.Indexes().List(ctx, &options.ListIndexesOptions{
		BatchSize: &batchSize,
		MaxTime:   &duration,
	})

	if err != nil {
		return err
	}

	// check for requested
	found := false
	for cur.Next(ctx) {
		d := bson.D{}

		if err := cur.Decode(&d); err != nil {
			return err
		}

		v := d.Map()["name"]
		if v != nil && v == opts.Name {
			found = true
			break
		}
	}

	if found {
		return nil
	}

	// create if required
	ks := bsonx.Doc{}
	for _, k := range keys {
		ks = ks.Append(k, bsonx.Int32(int32(1)))
	}

	index := mongo.IndexModel{
		Keys:    ks,
		Options: opts,
	}
	_, err = c.Indexes().CreateOne(ctx, index, &options.CreateIndexesOptions{MaxTime: &duration})
	return err
}
