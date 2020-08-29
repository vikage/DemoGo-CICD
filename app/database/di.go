package database

import (
	"go-cicd/app/di"
	"go-cicd/app/di/gdi"
	"reflect"
)

var (
	// DatabaseClientType use for database.Client
	DatabaseClientType = reflect.TypeOf((*Client)(nil)).Elem()
	// TransactionClientType use for database.Client transaction
	TransactionClientType = reflect.TypeOf((*TransactionClient)(nil)).Elem()
)

// ResolveDatabaseClient get database client from di
func ResolveDatabaseClient() Client {
	db, err := di.DefaultContainer.Resolve(DatabaseClientType)
	if err != nil {
		return nil
	}

	return db.(Client)
}

// ResolveTransactionClient get transaction client from di
func ResolveTransactionClient() Client {
	db, err := di.DefaultContainer.Resolve(TransactionClientType)
	if err != nil {
		return nil
	}

	return db.(Client)
}

// RegisterDependencyInContainer register database dependency
func RegisterDependencyInContainer(container *gdi.Container) {
	container.Register(DatabaseClientType, GetMongoSingletonClient)
	container.Register(TransactionClientType, GetMongoSingletonTransactionClient)
}
