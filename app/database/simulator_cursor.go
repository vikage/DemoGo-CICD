package database

import (
	"context"
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

// SimulatorCursor cursor bridge
type SimulatorCursor struct {
	data         []map[string]interface{}
	currentIndex int
}

// Next next pointer
func (cursor *SimulatorCursor) Next(ctx context.Context) bool {
	if cursor.currentIndex+1 < len(cursor.data) {
		cursor.currentIndex++
		return true
	}

	return false
}

// Decode decode to struct
func (cursor *SimulatorCursor) Decode(val interface{}) error {
	document := cursor.data[cursor.currentIndex]
	data, err := json.Marshal(document)
	if err != nil {
		return err
	}

	return bson.UnmarshalExtJSON(data, false, val)
}
