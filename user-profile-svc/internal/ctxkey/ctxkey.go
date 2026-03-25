package ctxkey

// UserIDKey is the context key type for storing authenticated user id.
// We use a dedicated type to avoid key collisions across packages.
type UserIDKey struct{}
