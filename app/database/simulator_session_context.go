package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

// SimulatorSessionContext session context bridge
type SimulatorSessionContext struct {
	client SimulatorClient
}

// Client get mongo client
func (simulator *SimulatorSessionContext) Client() Client {
	return &simulator.client
}

// CommitTransaction commit mongo transaction
func (simulator *SimulatorSessionContext) CommitTransaction(ctx context.Context) error {
	return nil
}

// AbortTransaction abort mongo transaction
func (simulator *SimulatorSessionContext) AbortTransaction(ctx context.Context) error {
	return nil
}

// StartTransaction start a transaction
func (simulator *SimulatorSessionContext) StartTransaction(opts ...*options.TransactionOptions) error {
	return nil
}

// Context Get context
func (simulator *SimulatorSessionContext) Context() context.Context {
	return context.TODO()
}
