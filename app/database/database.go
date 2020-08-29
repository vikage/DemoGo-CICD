package database

import (
	"context"
	"fmt"
	"time"
)

// MongoTimeoutContext get timeout context for mongo
func MongoTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}

// MongoTimeoutContextWithSeconds get timeout for mongo with time
func MongoTimeoutContextWithSeconds(seconds int64) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Duration(seconds)*time.Second)
}

// ErrDatabaseCanNotConnect define error for case can not connect
var ErrDatabaseCanNotConnect = fmt.Errorf("Can not connect to database")

// UseSession create session
func UseSession(client Client, fn func(sessonContext SessionContext)) error {
	if client == nil {
		return ErrDatabaseCanNotConnect
	}

	ctx, cancel := MongoTimeoutContextWithSeconds(10)
	defer cancel()

	return client.UseSession(ctx, func(sessionContext SessionContext) error {
		fn(sessionContext)
		return nil
	})
}
