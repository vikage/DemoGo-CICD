package database

import (
	"context"
	"go-cicd/app/logger"
	"os"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoClient implement bridge for mongo client
type MongoClient struct {
	internal  *mongo.Client
	connected bool
}

var lock sync.Mutex
var client *MongoClient
var transactionClient *MongoClient

// NewMongoClientWithDefaultConfiguration create mongo client with default configuration
func NewMongoClientWithDefaultConfiguration() Client {
	return &MongoClient{internal: createConnection(readpref.SecondaryPreferredMode)}
}

// NewMongoTransactionClientWithDefaultConfiguration create mongo client with default configuration
func NewMongoTransactionClientWithDefaultConfiguration() Client {
	return &MongoClient{internal: createConnection(readpref.PrimaryMode)}
}

// GetMongoSingletonClient get mongodb singleton client
func GetMongoSingletonClient() Client {
	if client == nil {
		lock.Lock()
		client = NewMongoClientWithDefaultConfiguration().(*MongoClient)
		lock.Unlock()
	}

	return client
}

// GetMongoSingletonTransactionClient get transaction singleton client
func GetMongoSingletonTransactionClient() Client {
	if transactionClient == nil {
		lock.Lock()
		transactionClient = NewMongoTransactionClientWithDefaultConfiguration().(*MongoClient)
		lock.Unlock()
	}

	return transactionClient
}

// Ping ensure connected
func (client *MongoClient) Ping(ctx context.Context) error {
	return client.internal.Ping(ctx, nil)
}

// Database get database object
func (client *MongoClient) Database(name string, opts ...*options.DatabaseOptions) Database {
	// client.Ping(context.TODO())
	internal := client.internal.Database(name, opts...)
	return &MongoDatabase{db: internal}
}

// UseSession create mongo session
func (client *MongoClient) UseSession(ctx context.Context, fn func(SessionContext) error) error {
	return client.internal.UseSession(ctx, func(sessionContext mongo.SessionContext) error {
		context := MongoSessionContext{
			client: MongoClient{
				internal: sessionContext.Client(),
			},
			ctx: sessionContext,
		}

		err := sessionContext.StartTransaction()
		if err != nil {
			return err
		}

		err = fn(&context)
		sessionContext.EndSession(sessionContext)
		return err
	})
}

func createConnection(readprefMode readpref.Mode) *mongo.Client {
	connectStr := os.Getenv("MONGO_CONNECT_STR")
	clientOptions := options.Client().ApplyURI(connectStr)
	readPreference, err := readpref.New(readprefMode)

	if err != nil {
		panic(err)
	}

	ctx, cancel := MongoTimeoutContext()
	defer cancel()
	clientOptions.ReadPreference = readPreference
	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		logger.Debug("Can connect to server %s, error: %s", connectStr, err)
		panic(err)
	}

	return client
}
