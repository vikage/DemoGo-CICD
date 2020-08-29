package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// MongoCursor cursor bridge
type MongoCursor struct {
	internal *mongo.Cursor
}

// Next next pointer
func (cursor *MongoCursor) Next(ctx context.Context) bool {
	return cursor.internal.Next(ctx)
}

// Decode decode to struct
func (cursor *MongoCursor) Decode(val interface{}) error {
	return cursor.internal.Decode(val)
}
