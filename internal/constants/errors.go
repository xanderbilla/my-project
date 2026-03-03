package constants

// This file defines all error types and codes used throughout the application.
// Having all error codes in one place makes them easy to find and maintain.

// Error Types represent high-level categories of errors.
// Think of these as broad classifications that help organize different kinds of problems.
//
// Frontend applications use these to decide how to handle errors:
// - VALIDATION_ERROR -> Show field-specific error messages
// - AUTHENTICATION_ERROR -> Redirect to login page
// - SYSTEM_ERROR -> Show generic "try again later" message
const (
	// ErrorTypeValidation is used when user input doesn't meet validation rules
	// Example: empty required field, invalid email format, age too low
	ErrorTypeValidation = "VALIDATION_ERROR"

	// ErrorTypeAuthentication is used when a user fails to prove their identity
	// Example: wrong password, expired token, missing authorization header
	ErrorTypeAuthentication = "AUTHENTICATION_ERROR"

	// ErrorTypeAuthorization is used when a user doesn't have permission
	// Example: regular user trying to access admin-only features
	ErrorTypeAuthorization = "AUTHORIZATION_ERROR"

	// ErrorTypeNotFound is used when a requested resource doesn't exist
	// Example: trying to get user ID 999 when it doesn't exist in database
	ErrorTypeNotFound = "RESOURCE_NOT_FOUND"

	// ErrorTypeConflict is used when an operation would create a duplicate or violate constraints
	// Example: creating a user with an email that already exists
	ErrorTypeConflict = "RESOURCE_CONFLICT"

	// ErrorTypeRateLimit is used when a client makes too many requests
	// Example: more than 100 requests per minute from the same IP address
	ErrorTypeRateLimit = "RATE_LIMIT_EXCEEDED"

	// ErrorTypeSystem is used for unexpected server-side errors
	// Example: database connection failure, third-party API timeout
	ErrorTypeSystem = "SYSTEM_ERROR"

	// ErrorTypeBadRequest is used when the request is malformed
	// Example: invalid JSON, wrong Content-Type header, missing required header
	ErrorTypeBadRequest = "BAD_REQUEST"
)

// Error Codes are specific, machine-readable identifiers for individual errors.
// These are more specific than Error Types and allow programmatic error handling.
//
// For example, the frontend can check if code === "USER_NOT_FOUND" and show
// a specific message or take a specific action.
//
// Naming convention: Use UPPERCASE_WITH_UNDERSCORES for consistency.
const (
	// User-related error codes
	// These codes are used in user management operations

	// ErrCodeUserNotFound indicates a user lookup failed because the user doesn't exist
	ErrCodeUserNotFound = "USER_NOT_FOUND"

	// ErrCodeUserAlreadyExists indicates attempting to create a user that already exists
	// Usually happens with duplicate email or username
	ErrCodeUserAlreadyExists = "USER_ALREADY_EXISTS"

	// ErrCodeInvalidUserData indicates user data failed validation
	ErrCodeInvalidUserData = "INVALID_USER_DATA"

	// Validation-related error codes
	// These codes are used when request data fails validation

	// ErrCodeInvalidRequestBody indicates the overall request body is invalid
	ErrCodeInvalidRequestBody = "INVALID_REQUEST_BODY"

	// ErrCodeInvalidJSON indicates the request body is not valid JSON
	// Example: missing closing brace, trailing comma, unquoted keys
	ErrCodeInvalidJSON = "INVALID_JSON"

	// ErrCodeMissingField indicates a required field was not provided
	ErrCodeMissingField = "MISSING_REQUIRED_FIELD"

	// ErrCodeInvalidEmail indicates an email address doesn't meet format requirements
	ErrCodeInvalidEmail = "INVALID_EMAIL"

	// System-related error codes
	// These codes are used for server-side failures

	// ErrCodeDatabaseError indicates a database operation failed
	// Example: connection timeout, query error, transaction failure
	ErrCodeDatabaseError = "DATABASE_ERROR"

	// ErrCodeInternalError indicates an unexpected error occurred
	// This is a catch-all for errors that don't fit other categories
	ErrCodeInternalError = "INTERNAL_SERVER_ERROR"

	// ErrCodeServiceUnavailable indicates a required service is not available
	// Example: database is down, cache server unreachable
	ErrCodeServiceUnavailable = "SERVICE_UNAVAILABLE"

	// Request-related error codes
	// These codes are used when the request itself is problematic

	// ErrCodeInvalidContentType indicates the Content-Type header is wrong
	// Example: sending XML when API expects JSON
	ErrCodeInvalidContentType = "INVALID_CONTENT_TYPE"

	// ErrCodeMethodNotAllowed indicates using wrong HTTP method
	// Example: sending POST to an endpoint that only accepts GET
	ErrCodeMethodNotAllowed = "METHOD_NOT_ALLOWED"
)

// APIVersion identifies the version of your API.
// This helps with:
// - API versioning and backwards compatibility
// - Debugging (knowing which version a request used)
// - Gradual rollout of new features
//
// Common versioning strategies:
// - "v1", "v2", "v3" for major versions
// - "1.0", "1.1", "2.0" for more granular versioning
const APIVersion = "v1"
