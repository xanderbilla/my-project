package types

import (
	"fmt"
	"net/http"
)

// AppError represents a custom application error with rich context information.
// This is the core error type used throughout the application.
//
// Unlike Go's standard error type which only has a string message, AppError provides:
// - Separate messages for users and developers
// - Machine-readable error codes
// - HTTP status codes
// - Additional context data
// - Ability to wrap underlying errors
//
// Example usage:
//
//	err := types.NewNotFoundError("User", "123")
//	// Creates an error indicating user 123 was not found
type AppError struct {
	// Type categorizes the error at a high level
	// Examples: "VALIDATION_ERROR", "AUTHENTICATION_ERROR", "SYSTEM_ERROR"
	Type string

	// Code is a specific machine-readable identifier for this error
	// Examples: "USER_NOT_FOUND", "INVALID_EMAIL", "DATABASE_ERROR"
	// Frontend applications can use this to handle specific errors programmatically
	Code string

	// UserMessage is a safe, user-friendly message that can be displayed in the UI
	// Never include sensitive information like database details or stack traces here
	UserMessage string

	// DeveloperMessage contains technical details for debugging
	// This helps developers understand what went wrong
	DeveloperMessage string

	// Context stores additional data about the error
	// Example: {"userId": 123, "attemptedAction": "delete", "tableName": "users"}
	Context map[string]interface{}

	// Retryable indicates if retrying the request might succeed
	// true: temporary issues (network timeout, service busy)
	// false: permanent issues (invalid input, resource not found)
	Retryable bool

	// StatusCode is the HTTP status code that should be returned
	// Examples: 404 (Not Found), 400 (Bad Request), 500 (Internal Server Error)
	StatusCode int

	// OriginalError stores the underlying error if this AppError wraps another error
	// Useful for maintaining the error chain for debugging
	OriginalError error
}

// Error implements the error interface, allowing AppError to be used as a standard Go error.
// This method is called when you convert the error to a string (e.g., in logs).
//
// Returns a formatted string with the error code and messages.
// If there's an original error, it's included in the output.
func (e *AppError) Error() string {
	if e.OriginalError != nil {
		return fmt.Sprintf("%s: %s (original: %v)", e.Code, e.DeveloperMessage, e.OriginalError)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.DeveloperMessage)
}

// Unwrap returns the original error for error unwrapping.
// This is part of Go's error wrapping convention introduced in Go 1.13.
// It allows functions like errors.Is() and errors.As() to work with wrapped errors.
func (e *AppError) Unwrap() error {
	return e.OriginalError
}

// ToErrorInfo converts AppError to ErrorInfo for inclusion in API responses.
// ErrorInfo is the structure that gets sent to clients as JSON.
// This separation allows us to have a rich internal error type (AppError)
// while controlling exactly what information is exposed in the API response.
func (e *AppError) ToErrorInfo() *ErrorInfo {
	return &ErrorInfo{
		Type:             e.Type,
		Code:             e.Code,
		UserMessage:      e.UserMessage,
		DeveloperMessage: e.DeveloperMessage,
		Context:          e.Context,
		Retryable:        e.Retryable,
	}
}

// The following functions are constructor functions that create specific types of errors.
// Using these constructors ensures consistency across the application.
// Instead of manually creating AppError structs, you call these functions.

// NewValidationError creates an error for invalid user input.
// Use this when the client sends data that doesn't meet validation requirements.
//
// Parameters:
//   - code: specific error code (e.g., "INVALID_EMAIL")
//   - userMsg: message to show to the user
//   - devMsg: technical details for developers
//
// Returns HTTP 422 (Unprocessable Entity) status code.
func NewValidationError(code, userMsg, devMsg string) *AppError {
	return &AppError{
		Type:             "VALIDATION_ERROR",
		Code:             code,
		UserMessage:      userMsg,
		DeveloperMessage: devMsg,
		StatusCode:       http.StatusUnprocessableEntity, // HTTP 422
		Retryable:        false,                          // Retrying won't help - the client needs to fix their input
	}
}

// NewNotFoundError creates an error for when a requested resource doesn't exist.
// Use this when a user tries to access something that can't be found (user, product, order, etc.).
//
// Parameters:
//   - resource: the type of resource (e.g., "User", "Product", "Order")
//   - id: the identifier that was searched for
//
// Returns HTTP 404 (Not Found) status code.
//
// Example:
//
//	err := NewNotFoundError("User", "123")
//	// User message: "The requested User could not be found."
//	// Developer message: "User with ID 123 does not exist."
func NewNotFoundError(resource, id string) *AppError {
	return &AppError{
		Type:             "RESOURCE_NOT_FOUND",
		Code:             fmt.Sprintf("%s_NOT_FOUND", resource), // e.g., "USER_NOT_FOUND"
		UserMessage:      fmt.Sprintf("The requested %s could not be found.", resource),
		DeveloperMessage: fmt.Sprintf("%s with ID %s does not exist.", resource, id),
		StatusCode:       http.StatusNotFound, // HTTP 404
		Retryable:        false,               // The resource doesn't exist, so retrying won't help
		// Store the resource type and ID in context for debugging
		Context: map[string]interface{}{
			"resource": resource,
			"id":       id,
		},
	}
}

// NewConflictError creates an error for resource conflicts.
// Use this when an operation fails because it would create a duplicate or violate a constraint.
//
// Parameters:
//   - code: specific error code (e.g., "USER_ALREADY_EXISTS")
//   - userMsg: message to show to the user
//   - devMsg: technical details for developers
//
// Returns HTTP 409 (Conflict) status code.
//
// Example:
//
//	err := NewConflictError("USER_ALREADY_EXISTS",
//	  "A user with this email already exists",
//	  "User with email john@example.com already exists in database")
func NewConflictError(code, userMsg, devMsg string) *AppError {
	return &AppError{
		Type:             "RESOURCE_CONFLICT",
		Code:             code,
		UserMessage:      userMsg,
		DeveloperMessage: devMsg,
		StatusCode:       http.StatusConflict, // HTTP 409
		Retryable:        false,               // The conflict needs to be resolved, retrying won't help
	}
}

// NewAuthenticationError creates an error for authentication failures.
// Use this when a user fails to prove their identity (wrong password, invalid token, etc.).
//
// Parameters:
//   - reason: technical reason for the authentication failure
//
// Returns HTTP 401 (Unauthorized) status code.
//
// Example:
//
//	err := NewAuthenticationError("Invalid JWT token: token expired")
//	// User sees: "Authentication failed. Please check your credentials."
//	// Developers see: "Invalid JWT token: token expired"
func NewAuthenticationError(reason string) *AppError {
	return &AppError{
		Type:             "AUTHENTICATION_ERROR",
		Code:             "AUTHENTICATION_FAILED",
		UserMessage:      "Authentication failed. Please check your credentials.",
		DeveloperMessage: reason,
		StatusCode:       http.StatusUnauthorized, // HTTP 401
		Retryable:        false,                   // User needs to provide correct credentials
	}
}

// NewAuthorizationError creates an error for authorization failures.
// Use this when an authenticated user tries to do something they don't have permission for.
//
// Note: Authentication vs Authorization:
//   - Authentication (401): "Who are you?" - proving identity
//   - Authorization (403): "What are you allowed to do?" - permission check
//
// Parameters:
//   - action: what the user tried to do (e.g., "delete", "update")
//   - resource: what they tried to do it on (e.g., "User", "Order")
//
// Returns HTTP 403 (Forbidden) status code.
//
// Example:
//
//	err := NewAuthorizationError("delete", "User")
//	// User sees: "You don't have permission to perform this action."
//	// Developers see: "User not authorized to delete on User"
func NewAuthorizationError(action, resource string) *AppError {
	return &AppError{
		Type:             "AUTHORIZATION_ERROR",
		Code:             "INSUFFICIENT_PERMISSIONS",
		UserMessage:      "You don't have permission to perform this action.",
		DeveloperMessage: fmt.Sprintf("User not authorized to %s on %s", action, resource),
		StatusCode:       http.StatusForbidden, // HTTP 403
		Retryable:        false,                // User needs elevated permissions, retrying won't help
		// Store what was attempted for debugging
		Context: map[string]interface{}{
			"action":   action,
			"resource": resource,
		},
	}
}

// NewBadRequestError creates an error for malformed requests.
// Use this when the client sends a request that's structurally wrong
// (invalid JSON, wrong content type, missing required headers, etc.).
//
// Parameters:
//   - code: specific error code (e.g., "INVALID_JSON", "INVALID_CONTENT_TYPE")
//   - userMsg: message to show to the user
//   - devMsg: technical details for developers
//
// Returns HTTP 400 (Bad Request) status code.
func NewBadRequestError(code, userMsg, devMsg string) *AppError {
	return &AppError{
		Type:             "BAD_REQUEST",
		Code:             code,
		UserMessage:      userMsg,
		DeveloperMessage: devMsg,
		StatusCode:       http.StatusBadRequest, // HTTP 400
		Retryable:        false,                 // Client needs to fix the request format
	}
}

// NewInternalError creates an error for unexpected server-side failures.
// Use this when something goes wrong that the client didn't cause.
//
// This error type hides technical details from users (security best practice)
// but logs the full error for developers to investigate.
//
// Parameters:
//   - err: the underlying error that occurred
//
// Returns HTTP 500 (Internal Server Error) status code.
//
// Example:
//
//	if err := processData(); err != nil {
//	  return NewInternalError(err)
//	}
func NewInternalError(err error) *AppError {
	return &AppError{
		Type:             "SYSTEM_ERROR",
		Code:             "INTERNAL_SERVER_ERROR",
		UserMessage:      "An unexpected error occurred. Please try again later.",
		DeveloperMessage: "Internal server error",
		StatusCode:       http.StatusInternalServerError, // HTTP 500
		Retryable:        true,                           // Might be a temporary issue, worth retrying
		OriginalError:    err,                            // Store the original error for logging
	}
}

// NewDatabaseError creates an error for database operation failures.
// Use this when insert, update, delete, or query operations fail.
//
// Parameters:
//   - operation: what database operation was attempted (e.g., "insert", "query", "update")
//   - err: the underlying database error
//
// Returns HTTP 500 (Internal Server Error) status code.
//
// Example:
//
//	if err := db.Insert(user); err != nil {
//	  return NewDatabaseError("insert", err)
//	}
func NewDatabaseError(operation string, err error) *AppError {
	return &AppError{
		Type:             "SYSTEM_ERROR",
		Code:             "DATABASE_ERROR",
		UserMessage:      "A database error occurred. Please try again later.",
		DeveloperMessage: fmt.Sprintf("Database %s failed", operation),
		StatusCode:       http.StatusInternalServerError, // HTTP 500
		Retryable:        true,                           // Database might recover, worth retrying
		OriginalError:    err,                            // Store the actual database error
		// Store what operation failed for debugging
		Context: map[string]interface{}{
			"operation": operation,
		},
	}
}

// NewRateLimitError creates an error for when a client exceeds rate limits.
// Use this to protect your API from being overwhelmed by too many requests.
//
// Rate limiting prevents abuse and ensures fair resource usage.
//
// Parameters:
//   - retryAfter: how many seconds the client should wait before trying again
//
// Returns HTTP 429 (Too Many Requests) status code.
//
// Example:
//
//	err := NewRateLimitError(60)  // Tell client to wait 60 seconds
func NewRateLimitError(retryAfter int) *AppError {
	return &AppError{
		Type:             "RATE_LIMIT_EXCEEDED",
		Code:             "TOO_MANY_REQUESTS",
		UserMessage:      "You have exceeded the rate limit. Please try again later.",
		DeveloperMessage: fmt.Sprintf("Rate limit exceeded. Retry after %d seconds", retryAfter),
		StatusCode:       http.StatusTooManyRequests, // HTTP 429
		Retryable:        true,                       // Client can retry after waiting
		// Tell the client when they can try again
		Context: map[string]interface{}{
			"retryAfter": retryAfter,
		},
	}
}

// NewServiceUnavailableError creates an error for when a service is down or overloaded.
// Use this when an external dependency (database, cache, third-party API) is unavailable.
//
// Parameters:
//   - service: name of the unavailable service (e.g., "database", "cache", "payment-api")
//
// Returns HTTP 503 (Service Unavailable) status code.
//
// Example:
//
//	if !database.IsHealthy() {
//	  return NewServiceUnavailableError("database")
//	}
func NewServiceUnavailableError(service string) *AppError {
	return &AppError{
		Type:             "SYSTEM_ERROR",
		Code:             "SERVICE_UNAVAILABLE",
		UserMessage:      "The service is temporarily unavailable. Please try again later.",
		DeveloperMessage: fmt.Sprintf("%s service is unavailable", service),
		StatusCode:       http.StatusServiceUnavailable, // HTTP 503
		Retryable:        true,                          // Service might come back up
		// Store which service is down for monitoring
		Context: map[string]interface{}{
			"service": service,
		},
	}
}

// WithContext adds additional context information to an existing error.
// This is useful for adding error-specific details after creating the error.
//
// Returns the same error (allowing method chaining).
//
// Example:
//
//	err := types.NewNotFoundError("User", userID)
//	err.WithContext("attemptedBy", currentUserEmail)
//	err.WithContext("requestOrigin", "mobile-app")
func (e *AppError) WithContext(key string, value interface{}) *AppError {
	// Initialize the Context map if it doesn't exist yet
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	// Add the key-value pair to the context
	e.Context[key] = value
	// Return the error to allow chaining (e.g., err.WithContext("a", 1).WithContext("b", 2))
	return e
}

// WithOriginalError adds an underlying error to an existing AppError.
// Use this when you want to wrap an existing error with additional context.
//
// Returns the same error (allowing method chaining).
//
// Example:
//
//	dbErr := database.Query()
//	appErr := types.NewDatabaseError("query", nil)
//	return appErr.WithOriginalError(dbErr)
func (e *AppError) WithOriginalError(err error) *AppError {
	// Store the original error for debugging and logging
	e.OriginalError = err
	// Return the error to allow method chaining
	return e
}
