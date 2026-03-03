package constants

// ContextKey is a custom type for context keys to avoid collisions.
// Using a custom type instead of string prevents key collisions
// when different packages use context.WithValue().
//
// Why this matters:
// If two packages both use context.WithValue(ctx, "requestID", value),
// they could overwrite each other's values. By using a custom type,
// we ensure our keys are unique.
//
// Best practice from the Go documentation:
// "The provided key must be comparable and should not be of type
// string or any other built-in type to avoid collisions between
// packages using context."
type ContextKey string

// Context key constants
const (
	// RequestIDKey is the context key for request IDs
	RequestIDKey ContextKey = "requestID"
)
