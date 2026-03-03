package response

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	json "github.com/goccy/go-json"

	"github.com/xanderbilla/my-project/internal/constants"
	"github.com/xanderbilla/my-project/internal/types"
)

// Error sends a standardized error response to the client.
// This is the base function that all error responses use.
//
// The function:
//  1. Builds a standard error response structure
//  2. Logs the error for monitoring and debugging
//  3. Sends the error response as JSON
//
// Parameters:
//   - w: ResponseWriter to send the response to
//   - r: The original HTTP request
//   - status: HTTP status code (400, 404, 500, etc.)
//   - errorInfo: Detailed error information to include in the response
func Error(w http.ResponseWriter, r *http.Request, status int, errorInfo *types.ErrorInfo) {
	// Get the unique request ID for tracking
	requestID := getRequestID(r.Context())

	// Build the error response structure
	response := types.APIResponse{
		Success:   false,            // Indicates the request failed
		Status:    status,           // HTTP status code
		Timestamp: time.Now().UTC(), // When the error occurred
		RequestID: requestID,        // For tracking this specific request
		Path:      r.URL.Path,       // Which endpoint was called
		Method:    r.Method,         // HTTP method used
		Error:     errorInfo,        // Detailed error information
	}

	// Log the error for monitoring and debugging
	// This helps developers find and fix issues
	logError(r, status, errorInfo)

	// Set response headers and send JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		slog.Error("failed to encode error response", "error", err)
	}
}

// HandleAppError converts an AppError to an HTTP response.
// Use this when you have an AppError object that you want to send to the client.
//
// Example:
//
//	err := types.NewNotFoundError("User", "123")
//	response.HandleAppError(w, r, err)
func HandleAppError(w http.ResponseWriter, r *http.Request, err *types.AppError) {
	// Extract the HTTP status code and error info from AppError
	// and send them using the standard Error function
	Error(w, r, err.StatusCode, err.ToErrorInfo())
}

// ValidationError sends a validation error response with field-specific errors.
// Use this when user input fails validation checks.
//
// Parameters:
//   - validationErrors: Slice of field-level validation errors
//
// Example:
//
//	var errors []types.ValidationError
//	if user.Email == "" {
//	  errors = append(errors, types.ValidationError{
//	    Field: "email",
//	    RejectedValue: "",
//	    Message: "Email is required",
//	  })
//	}
//	response.ValidationError(w, r, errors)
func ValidationError(w http.ResponseWriter, r *http.Request, validationErrors []types.ValidationError) {
	errorInfo := &types.ErrorInfo{
		Type:             constants.ErrorTypeValidation,
		Code:             constants.ErrCodeInvalidRequestBody,
		UserMessage:      "Some fields in your request are invalid.",
		DeveloperMessage: "Request validation failed. Check validationErrors for details.",
		ValidationErrors: validationErrors, // Field-specific errors for the client to display
		Retryable:        false,            // Client needs to fix their input before retrying
	}

	// Send HTTP 422 Unprocessable Entity
	Error(w, r, http.StatusUnprocessableEntity, errorInfo)
}

// NotFound sends a 404 Not Found error response.
// Use this when a requested resource doesn't exist.
//
// Parameters:
//   - resourceType: Type of resource (e.g., "User", "Product", "Order")
//   - resourceID: The ID that was searched for
//
// Example:
//
//	response.NotFound(w, r, "User", "123")
//	// Returns: "The requested User could not be found."
func NotFound(w http.ResponseWriter, r *http.Request, resourceType, resourceID string) {
	// Create a NotFound error and send it
	appErr := types.NewNotFoundError(resourceType, resourceID)
	HandleAppError(w, r, appErr)
}

// BadRequest sends a 400 Bad Request error response.
// Use this when the client sends a malformed request.
//
// Parameters:
//   - code: Specific error code (e.g., "INVALID_JSON")
//   - userMsg: User-friendly message
//   - devMsg: Technical details for developers
//
// Example:
//
//	response.BadRequest(w, r,
//	  constants.ErrCodeInvalidJSON,
//	  "Invalid JSON format",
//	  "Failed to decode request body: unexpected end of JSON input")
func BadRequest(w http.ResponseWriter, r *http.Request, code, userMsg, devMsg string) {
	appErr := types.NewBadRequestError(code, userMsg, devMsg)
	HandleAppError(w, r, appErr)
}

// InternalError sends a 500 Internal Server Error response.
// Use this when an unexpected error occurs that the client didn't cause.
//
// This function automatically:
//   - Hides sensitive technical details from the user
//   - Logs the full error for developer investigation
//   - Marks the error as retryable
//
// Parameters:
//   - err: The underlying error that occurred
//
// Example:
//
//	if err := processData(); err != nil {
//	  response.InternalError(w, r, err)
//	  return
//	}
func InternalError(w http.ResponseWriter, r *http.Request, err error) {
	appErr := types.NewInternalError(err)
	HandleAppError(w, r, appErr)
}

// Conflict sends a 409 Conflict error response.
// Use this when an operation would create a duplicate or violate a constraint.
//
// Parameters:
//   - code: Specific error code (e.g., "USER_ALREADY_EXISTS")
//   - userMsg: User-friendly message
//   - devMsg: Technical details for developers
//
// Example:
//
//	response.Conflict(w, r,
//	  constants.ErrCodeUserAlreadyExists,
//	  "A user with this email already exists",
//	  "User with email john@example.com exists in database")
func Conflict(w http.ResponseWriter, r *http.Request, code, userMsg, devMsg string) {
	appErr := types.NewConflictError(code, userMsg, devMsg)
	HandleAppError(w, r, appErr)
}

// DatabaseError sends a database error response.
// Use this when a database operation fails.
//
// Parameters:
//   - operation: What database operation was attempted (e.g., "insert", "query", "update")
//   - err: The underlying database error
//
// Example:
//
//	if err := db.Insert(user); err != nil {
//	  response.DatabaseError(w, r, "insert", err)
//	  return
//	}
func DatabaseError(w http.ResponseWriter, r *http.Request, operation string, err error) {
	appErr := types.NewDatabaseError(operation, err)
	HandleAppError(w, r, appErr)
}

// Unauthorized sends a 401 Unauthorized error response.
// Use this when authentication fails (wrong password, invalid token, etc.).
//
// Parameters:
//   - reason: Technical reason for authentication failure
//
// Example:
//
//	response.Unauthorized(w, r, "JWT token expired")
func Unauthorized(w http.ResponseWriter, r *http.Request, reason string) {
	appErr := types.NewAuthenticationError(reason)
	HandleAppError(w, r, appErr)
}

// Forbidden sends a 403 Forbidden error response.
// Use this when an authenticated user doesn't have permission to perform an action.
//
// Parameters:
//   - action: What action was attempted (e.g., "delete", "update")
//   - resource: What resource was targeted (e.g., "User", "Order")
//
// Example:
//
//	response.Forbidden(w, r, "delete", "User")
//	// User sees: "You don't have permission to perform this action."
func Forbidden(w http.ResponseWriter, r *http.Request, action, resource string) {
	appErr := types.NewAuthorizationError(action, resource)
	HandleAppError(w, r, appErr)
}

// logError writes error information to the application logs.
// This is called automatically by the Error function.
//
// The logs include:
//   - Request ID for tracking
//   - Path and method that caused the error
//   - HTTP status code
//   - Error type and code
//   - Developer message with technical details
//
// These logs are essential for:
//   - Debugging issues
//   - Monitoring error rates
//   - Identifying patterns in failures
func logError(r *http.Request, status int, errorInfo *types.ErrorInfo) {
	slog.Error("API Error",
		"requestId", getRequestID(r.Context()),
		"path", r.URL.Path,
		"method", r.Method,
		"status", status,
		"errorType", errorInfo.Type,
		"errorCode", errorInfo.Code,
		"developerMessage", errorInfo.DeveloperMessage,
	)
}

// getRequestID retrieves the request ID from the request context.
// The request ID is added by the RequestID middleware earlier in the request chain.
//
// Returns "unknown" if no request ID is found (which should never happen in production
// if the RequestID middleware is properly configured).
func getRequestID(ctx context.Context) string {
	// Use the exported GetRequestID function from helpers.go
	return GetRequestID(ctx)
}
