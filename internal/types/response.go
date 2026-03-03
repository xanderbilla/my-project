package types

import "time"

// APIResponse is the standard wrapper structure for all API responses.
// This ensures consistency across all endpoints in your API.
// Every response from the API will follow this same structure, whether it's a success or error.
//
// Example success response:
//
//	{
//	  "success": true,
//	  "status": 200,
//	  "timestamp": "2026-03-04T10:30:00Z",
//	  "requestId": "abc-123",
//	  "path": "/api/users",
//	  "method": "GET",
//	  "message": "Users retrieved successfully",
//	  "data": [...],
//	  "meta": {...}
//	}
type APIResponse struct {
	// Success indicates whether the request was successful (true) or failed (false)
	Success bool `json:"success"`

	// Status is the HTTP status code (200, 404, 500, etc.)
	Status int `json:"status"`

	// Timestamp shows when the response was generated (in UTC timezone)
	Timestamp time.Time `json:"timestamp"`

	// RequestID is a unique identifier for this specific request
	// Useful for tracking requests in logs and debugging
	RequestID string `json:"requestId"`

	// Path is the API endpoint that was called (e.g., "/api/users/123")
	Path string `json:"path"`

	// Method is the HTTP method used (GET, POST, PUT, DELETE, etc.)
	Method string `json:"method"`

	// Message is a human-readable description of the result
	// Example: "User created successfully"
	// Note: omitempty means this field won't appear in JSON if it's empty
	Message string `json:"message,omitempty"`

	// Data contains the actual response data (user object, list of items, etc.)
	// This field is only present in successful responses
	// interface{} means it can hold any type of data
	Data interface{} `json:"data,omitempty"`

	// Meta contains additional metadata about the response
	// Like API version, pagination info, execution time, etc.
	Meta *Meta `json:"meta,omitempty"`

	// Error contains detailed error information
	// This field is only present when Success is false
	Error *ErrorInfo `json:"error,omitempty"`
}

// Meta contains additional metadata about the API response.
// This is useful for providing extra information that doesn't belong in the main data.
type Meta struct {
	// APIVersion helps clients know which version of the API they're using
	APIVersion string `json:"apiVersion"`

	// ExecutionTimeMs shows how long the request took to process (in milliseconds)
	// Formatted to 4 decimal places for precision (e.g., 0.1234 ms)
	// Useful for performance monitoring
	ExecutionTimeMs float64 `json:"executionTimeMs"`

	// The following fields are used for pagination (splitting large datasets into pages)

	// Page is the current page number (starts from 1)
	Page int `json:"page,omitempty"`

	// PageSize is how many items are shown per page
	PageSize int `json:"pageSize,omitempty"`

	// TotalCount is the total number of items available across all pages
	TotalCount int64 `json:"totalCount,omitempty"`
}

// ErrorInfo contains detailed information about an error.
// This structure follows enterprise best practices by separating:
// - User-facing messages (safe to show in UI)
// - Developer messages (detailed info for debugging)
// - Machine-readable codes (for programmatic error handling)
type ErrorInfo struct {
	// Type is a high-level category of the error
	// Examples: "VALIDATION_ERROR", "AUTHENTICATION_ERROR", "SYSTEM_ERROR"
	Type string `json:"type"`

	// Code is a specific machine-readable error code
	// Examples: "USER_NOT_FOUND", "INVALID_EMAIL", "DATABASE_ERROR"
	// Frontend can use this to decide how to handle the error
	Code string `json:"code"`

	// UserMessage is a safe, user-friendly message to show in the UI
	// Example: "The requested user could not be found."
	// Never include sensitive system information here
	UserMessage string `json:"userMessage"`

	// DeveloperMessage contains detailed technical information
	// Example: "User with ID 123 does not exist in database table 'users'"
	// This helps developers debug issues but should not be shown to end users
	DeveloperMessage string `json:"developerMessage"`

	// ValidationErrors contains field-specific validation errors
	// Only present when Type is "VALIDATION_ERROR"
	ValidationErrors []ValidationError `json:"validationErrors,omitempty"`

	// Context provides additional contextual information about the error
	// Example: {"userId": 123, "attemptedAction": "delete"}
	Context map[string]interface{} `json:"context,omitempty"`

	// Retryable indicates whether the client should retry this request
	// true: temporary issue, retry might work (e.g., network timeout)
	// false: permanent issue, retrying won't help (e.g., invalid input)
	Retryable bool `json:"retryable"`
}

// ValidationError represents a single field validation error.
// Used when user input fails validation checks.
//
// Example: If a user submits an empty email field, you'd get:
//
//	{
//	  "field": "email",
//	  "rejectedValue": "",
//	  "message": "Email is required"
//	}
type ValidationError struct {
	// Field is the name of the field that failed validation
	Field string `json:"field"`

	// RejectedValue is the actual value that was rejected
	// Helps users understand what they submitted
	RejectedValue interface{} `json:"rejectedValue"`

	// Message explains why the value was rejected
	Message string `json:"message"`
}
