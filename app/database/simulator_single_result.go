package database

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// SimulatorSingleResult single result bridge
type SimulatorSingleResult struct {
	err  error
	data *map[string]interface{}
}

// Decode decode struct from result
func (result *SimulatorSingleResult) Decode(val interface{}) error {
	if result.err != nil {
		return result.err
	}

	if result.data == nil {
		return mongo.ErrNoDocuments
	}

	data, err := json.Marshal(*(result.data))
	if err != nil {
		return err
	}

	return bson.UnmarshalExtJSON(data, false, val)
}
