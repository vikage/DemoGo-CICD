package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// SimulatorClient simulator client bridge
type SimulatorClient struct {
}

// NewSimulatorClient create new simulator client
func NewSimulatorClient() Client {
	return &SimulatorClient{}
}

// ClearSimulatorData clear data
func ClearSimulatorData() {
	GlobalDatabases = make(map[string]*SimulatorDatabaseStorage)
}

// Ping ensure connected
func (client *SimulatorClient) Ping(ctx context.Context) error {
	return nil
}

// Database get database object
func (client *SimulatorClient) Database(name string, opts ...*options.DatabaseOptions) Database {
	storage := GlobalDatabases[name]
	if storage == nil {
		storage = &SimulatorDatabaseStorage{
			data: make(map[string]*SimulatorCollectionStorage),
		}

		GlobalDatabases[name] = storage
	}

	return &SimulatorDatabase{name: name, storage: storage}
}

// UseSession create mongo session
func (client *SimulatorClient) UseSession(ctx context.Context, fn func(SessionContext) error) error {
	context := SimulatorSessionContext{
		client: *client,
	}

	return fn(&context)
}
