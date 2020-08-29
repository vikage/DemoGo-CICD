package database

import (
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SimulatorDatabase bridge simulator database
type SimulatorDatabase struct {
	name    string
	storage *SimulatorDatabaseStorage
}

// Collection get mongo collection
func (db *SimulatorDatabase) Collection(name string, opts ...*options.CollectionOptions) Collection {
	storage := db.storage.data[name]
	if storage == nil {
		storage = &SimulatorCollectionStorage{
			data: make([]map[string]interface{}, 0),
		}

		db.storage.data[name] = storage
	}

	return &SimulatorCollection{
		name:    name,
		storage: storage,
	}
}
