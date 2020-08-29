package base

// ContextKey use for http context
type ContextKey string

const (
	// UserAuthenticatedKey key for user authenticate
	UserAuthenticatedKey ContextKey = "user"
)
