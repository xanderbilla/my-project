package response

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	json "github.com/goccy/go-json"

	"github.com/xanderbilla/my-project/internal/constants"
	"github.com/xanderbilla/my-project/internal/types"
)

// This file contains helper functions that make working with HTTP requests easier.
// These functions handle common tasks like extracting parameters, decoding JSON,
// and working with pagination.

// Context Helpers
// These functions help you work with Go's context, which stores request-scoped data.

// GetRequestID extracts the request ID from the request context.
// The request ID is added by the RequestID middleware and helps track individual requests.
//
// Use this when you need to include the request ID in logs or error messages.
//
// Example:
//
//	requestID := response.GetRequestID(r.Context())
//	log.Printf("Processing request %s", requestID)
func GetRequestID(ctx context.Context) string {
	if id, ok := ctx.Value(constants.RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}

// Pagination Helpers
// These functions help implement pagination, which splits large datasets into pages.

// PaginationParams holds pagination parameters extracted from a request.
// Page numbers start from 1, not 0.
type PaginationParams struct {
	// Page is the current page number (starts from 1)
	Page int

	// PageSize is how many items to show per page
	PageSize int

	// Offset is calculated automatically and used in database queries
	// Example: If Page=2 and PageSize=20, then Offset=20
	Offset int
}

// GetPaginationParams extracts pagination parameters from URL query string.
// It looks for "page" and "pageSize" query parameters.
//
// Default values:
//   - page: 1 (first page)
//   - pageSize: 20 (20 items per page)
//
// Constraints:
//   - page must be > 0
//   - pageSize must be between 1 and 100 (to prevent excessive data transfer)
//
// Example URL: /api/users?page=2&pageSize=50
//
// Usage:
//
//	params := response.GetPaginationParams(r)
//	users := database.GetUsers(params.Offset, params.PageSize)
func GetPaginationParams(r *http.Request) PaginationParams {
	params := PaginationParams{
		Page:     1,  // Default to first page
		PageSize: 20, // Default page size
	}

	// Try to parse the "page" query parameter
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		// Convert string to integer
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			params.Page = page
		}
		// If conversion fails or page <= 0, keep the default
	}

	// Try to parse the "pageSize" query parameter
	if sizeStr := r.URL.Query().Get("pageSize"); sizeStr != "" {
		// Convert string to integer, enforce max of 100
		if size, err := strconv.Atoi(sizeStr); err == nil && size > 0 && size <= 100 {
			params.PageSize = size
		}
		// If conversion fails or out of range, keep the default
	}

	// Calculate offset for database queries
	// Example: Page 1 → Offset 0, Page 2 → Offset 20, Page 3 → Offset 40
	params.Offset = (params.Page - 1) * params.PageSize

	return params
}

// PaginatedResponse sends a success response with pagination metadata.
// Use this when returning a list of items that's split across multiple pages.
//
// Parameters:
//   - data: The items for the current page
//   - page: Current page number
//   - pageSize: Number of items per page
//   - totalCount: Total number of items across all pages
//
// The response includes:
//   - The items array
//   - Current page number
//   - Page size
//   - Total item count
//   - Total number of pages (calculated automatically)
//
// Example:
//
//	users := getUsersForPage(2, 20)  // Get page 2, 20 items per page
//	totalUsers := countAllUsers()     // Get total count
//	response.PaginatedResponse(w, r, "Users retrieved", users, 2, 20, totalUsers)
func PaginatedResponse(w http.ResponseWriter, r *http.Request, message string, data interface{}, page, pageSize int, totalCount int64) {
	// Calculate total number of pages
	// Example: 45 items with pageSize 20 = 3 pages (20 + 20 + 5)
	totalPages := (totalCount + int64(pageSize) - 1) / int64(pageSize)

	// Send success response with pagination info
	Success(w, r, http.StatusOK, message, map[string]interface{}{
		"items": data,
		"pagination": map[string]interface{}{
			"page":       page,
			"pageSize":   pageSize,
			"totalCount": totalCount,
			"totalPages": totalPages,
		},
	})
}

// JSON Helpers
// These functions help with encoding and decoding JSON safely.

// DecodeJSON safely decodes a JSON request body into a Go struct.
// It performs several validations to prevent common issues.
//
// Checks performed:
//  1. Content-Type must be "application/json"
//  2. Request body size is limited (10MB maximum)
//  3. JSON must be valid
//  4. Unknown fields are rejected (prevents typos)
//  5. Only a single JSON object is allowed
//
// Parameters:
//   - r: The HTTP request containing the JSON body
//   - dst: Pointer to the struct to decode into
//
// Returns an error if any validation fails.
//
// Example:
//
//	var user types.User
//	if err := response.DecodeJSON(r, &user); err != nil {
//	  response.BadRequest(w, r, "INVALID_JSON", "Invalid JSON", err.Error())
//	  return
//	}
func DecodeJSON(r *http.Request, dst interface{}) error {
	// Check Content-Type header
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return fmt.Errorf("content-type must be application/json")
	}

	// Limit request body size to 10MB to prevent abuse
	// 10<<20 means 10 * 2^20 = 10,485,760 bytes = 10MB
	r.Body = http.MaxBytesReader(nil, r.Body, 10<<20)

	// Create a JSON decoder
	dec := json.NewDecoder(r.Body)

	// Reject unknown fields to catch typos in client code
	// Example: If struct has "name" field but JSON has "namee", it will error
	dec.DisallowUnknownFields()

	// Decode the JSON into the destination struct
	if err := dec.Decode(dst); err != nil {
		// Check if it's a type mismatch error (e.g., string instead of int)
		// Only treat "cannot unmarshal" errors as type mismatches
		// Other errors (syntax errors, etc.) should be treated as bad requests
		errMsg := err.Error()
		if strings.Contains(errMsg, "cannot unmarshal") {
			return &TypeMismatchError{Original: err}
		}
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Check if there's extra data after the JSON object
	// dec.More() returns true if there's more data to read
	if dec.More() {
		return fmt.Errorf("request body must contain only a single JSON object")
	}

	return nil
}

// TypeMismatchError represents a JSON type mismatch error
type TypeMismatchError struct {
	Original error
}

func (e *TypeMismatchError) Error() string {
	return e.Original.Error()
}

// DecodeJSONWithValidation decodes JSON and runs validation in one step.
// This combines JSON decoding and validation for convenience.
//
// Parameters:
//   - w, r: HTTP response and request
//   - dst: Pointer to struct to decode into
//   - validateFunc: Optional validation function to run after decoding
//
// Returns true if successful, false if there was an error.
// If false, an error response was already sent to the client.
//
// Example:
//
//	var user types.User
//	validateFunc := func(data interface{}) []types.ValidationError {
//	  u := data.(*types.User)
//	  return validateUser(*u)
//	}
//	if !response.DecodeJSONWithValidation(w, r, &user, validateFunc) {
//	  return  // Error response already sent
//	}
//	// user is valid, continue processing
func DecodeJSONWithValidation(w http.ResponseWriter, r *http.Request, dst interface{}, validateFunc func(interface{}) []types.ValidationError) bool {
	// Get request ID for logging
	requestID := GetRequestID(r.Context())

	// Try to decode the JSON
	if err := DecodeJSON(r, dst); err != nil {
		// Check if it's a type mismatch error - treat as validation error
		if typeMismatchErr, ok := err.(*TypeMismatchError); ok {
			validationErrors := parseTypeMismatchError(typeMismatchErr.Original)
			slog.Warn("JSON type mismatch detected", "requestId", requestID, "error", typeMismatchErr.Original.Error(), "path", r.URL.Path)
			ValidationError(w, r, validationErrors)
			return false
		}

		slog.Warn("JSON decode failed", "requestId", requestID, "error", err.Error(), "path", r.URL.Path)
		BadRequest(w, r, "INVALID_JSON", "Invalid JSON format", err.Error())
		return false
	}

	// Run validation if a validation function was provided
	if validateFunc != nil {
		if errors := validateFunc(dst); len(errors) > 0 {
			slog.Warn("request validation failed", "requestId", requestID, "error_count", len(errors), "path", r.URL.Path)
			ValidationError(w, r, errors)
			return false
		}
	}

	// Success - JSON decoded and validated
	return true
}

// parseTypeMismatchError converts a Go JSON type mismatch error into user-friendly validation errors
func parseTypeMismatchError(err error) []types.ValidationError {
	errMsg := err.Error()
	slog.Debug("parsing type mismatch error", "error", errMsg)

	// Extract field name and type information from error message
	// Error formats:
	// - "json: cannot unmarshal string into Go struct field User.age of type int"
	// - "json: cannot unmarshal number \" into Go struct field User.age of type int"
	var field string
	var expectedType string
	var receivedType string

	// Parse field name (e.g., "User.age" -> "age")
	if strings.Contains(errMsg, "field ") {
		parts := strings.Split(errMsg, "field ")
		if len(parts) > 1 {
			fieldParts := strings.Split(parts[1], " ")
			if len(fieldParts) > 0 {
				// Extract just the field name (after the last dot)
				fullField := fieldParts[0]
				fieldNameParts := strings.Split(fullField, ".")
				if len(fieldNameParts) > 0 {
					field = fieldNameParts[len(fieldNameParts)-1]
				}
			}
		}
	}

	// Parse expected type (e.g., "of type int")
	if strings.Contains(errMsg, "of type ") {
		parts := strings.Split(errMsg, "of type ")
		if len(parts) > 1 {
			expectedType = strings.TrimSpace(parts[1])
		}
	}

	// Determine received type from unmarshal error
	if strings.Contains(errMsg, "cannot unmarshal") {
		if strings.Contains(errMsg, "unmarshal string") {
			receivedType = "string"
		} else if strings.Contains(errMsg, "unmarshal number") {
			// "number \"" or similar means a quoted number (string)
			if strings.Contains(errMsg, "number \"") || strings.Contains(errMsg, "number \\\"") {
				receivedType = "string"
			} else {
				receivedType = "number"
			}
		} else if strings.Contains(errMsg, "unmarshal bool") {
			receivedType = "boolean"
		} else if strings.Contains(errMsg, "unmarshal array") || strings.Contains(errMsg, "unmarshal object") {
			receivedType = "object or array"
		}
	}

	// If we couldn't detect it, default to string (most common case)
	if receivedType == "" {
		receivedType = "string"
	}

	// Build user-friendly message
	var message string
	if field != "" && expectedType != "" {
		switch expectedType {
		case "int", "int64", "int32", "float64", "float32":
			message = fmt.Sprintf("%s must be a number, but received a %s", field, receivedType)
		case "string":
			message = fmt.Sprintf("%s must be a string, but received a %s", field, receivedType)
		case "bool":
			message = fmt.Sprintf("%s must be a boolean (true/false), but received a %s", field, receivedType)
		default:
			message = fmt.Sprintf("%s must be a %s, but received a %s", field, expectedType, receivedType)
		}
	} else if field != "" {
		message = fmt.Sprintf("%s has an invalid data type", field)
	} else {
		message = "Invalid data type in request"
	}

	return []types.ValidationError{
		{
			Field:         field,
			Message:       message,
			RejectedValue: nil, // We don't have access to the actual value here
		},
	}
}

// Query Parameter Helpers
// These functions extract and parse URL query parameters.

// GetQueryParam gets a string query parameter with a default value.
// If the parameter is missing or empty, the default value is returned.
//
// Example URL: /api/users?sort=name
//
//	sort := response.GetQueryParam(r, "sort", "id")  // Returns "name"
//	sort := response.GetQueryParam(r, "filter", "")   // Returns "" (default)
func GetQueryParam(r *http.Request, key, defaultValue string) string {
	value := r.URL.Query().Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// GetQueryParamInt gets an integer query parameter with a default value.
// If the parameter is missing, empty, or not a valid integer, the default is returned.
//
// Example URL: /api/products?minPrice=100
//
//	minPrice := response.GetQueryParamInt(r, "minPrice", 0)  // Returns 100
//	maxPrice := response.GetQueryParamInt(r, "maxPrice", 0)  // Returns 0 (default)
func GetQueryParamInt(r *http.Request, key string, defaultValue int) int {
	valueStr := r.URL.Query().Get(key)
	if valueStr == "" {
		return defaultValue
	}

	// Try to convert string to integer
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}

// GetQueryParamBool gets a boolean query parameter.
// Recognizes "true", "1" as true, everything else as false.
//
// Example URL: /api/users?active=true
//
//	active := response.GetQueryParamBool(r, "active")  // Returns true
//	deleted := response.GetQueryParamBool(r, "deleted") // Returns false
func GetQueryParamBool(r *http.Request, key string) bool {
	value := r.URL.Query().Get(key)
	return value == "true" || value == "1"
}

// Path Parameter Helpers
// These functions extract values from URL path parameters (route variables).

// GetPathParam gets a path parameter and validates it's not empty.
// Returns the value and true if found, or empty string and false if not found.
//
// Example route: /api/users/{user_id}
// Example URL: /api/users/123
//
// Usage:
//
//	userID, ok := response.GetPathParam(r, "user_id")
//	if !ok {
//	  response.BadRequest(w, r, "MISSING_PARAM", "User ID required", "user_id path parameter missing")
//	  return
//	}
func GetPathParam(r *http.Request, key string) (string, bool) {
	value := r.PathValue(key)
	if value == "" {
		return "", false
	}
	return value, true
}

// GetPathParamInt gets an integer path parameter.
// Returns the value and true if found and valid, or 0 and false otherwise.
//
// Example route: /api/posts/{post_id}
// Example URL: /api/posts/456
//
// Usage:
//
//	postID, ok := response.GetPathParamInt(r, "post_id")
//	if !ok {
//	  response.BadRequest(w, r, "INVALID_PARAM", "Invalid post ID", "post_id must be an integer")
//	  return
//	}
func GetPathParamInt(r *http.Request, key string) (int, bool) {
	valueStr := r.PathValue(key)
	if valueStr == "" {
		return 0, false
	}

	// Try to convert to integer
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, false
	}

	return value, true
}

// Request Body Helpers

// ReadBody reads the entire request body into a byte slice.
// Use this when you need raw body data (not JSON).
//
// Parameters:
//   - r: HTTP request
//   - maxSize: Maximum allowed body size in bytes
//
// Example:
//
//	body, err := response.ReadBody(r, 1<<20)  // Max 1MB
//	if err != nil {
//	  response.BadRequest(w, r, "BODY_ERROR", "Invalid body", err.Error())
//	  return
//	}
func ReadBody(r *http.Request, maxSize int64) ([]byte, error) {
	// Limit body size to prevent memory exhaustion
	r.Body = http.MaxBytesReader(nil, r.Body, maxSize)

	// Read entire body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Header Helpers

// GetAuthToken extracts a Bearer token from the Authorization header.
// Returns the token and true if found, or empty string and false if not found.
//
// Expected header format: "Authorization: Bearer <token>"
//
// Example:
//
//	token, ok := response.GetAuthToken(r)
//	if !ok {
//	  response.Unauthorized(w, r, "Missing or invalid Authorization header")
//	  return
//	}
//	// Verify token...
func GetAuthToken(r *http.Request) (string, bool) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", false
	}

	// Split "Bearer <token>" into ["Bearer", "<token>"]
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", false
	}

	return parts[1], true
}

// SetCacheHeaders sets HTTP cache headers to enable caching.
// Use this for responses that don't change frequently.
//
// Parameters:
//   - maxAge: How long (in seconds) the response can be cached
//
// Example:
//
//	response.SetCacheHeaders(w, 3600)  // Cache for 1 hour
func SetCacheHeaders(w http.ResponseWriter, maxAge int) {
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", maxAge))
}

// SetNoCacheHeaders disables all caching for this response.
// Use this for dynamic data that should never be cached.
//
// Example:
//
//	response.SetNoCacheHeaders(w)
//	// Send response with latest data...
func SetNoCacheHeaders(w http.ResponseWriter) {
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
}

// CORS Helpers
// CORS (Cross-Origin Resource Sharing) allows browsers to make requests to your API
// from different domains.

// SetCORSHeaders sets basic CORS headers that allow requests from any origin.
// Use this in development or for public APIs.
//
// WARNING: For production, you should restrict allowed origins.
func SetCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
}

// HandleCORS handles OPTIONS preflight requests.
// Browsers send OPTIONS requests before actual requests to check if CORS is allowed.
//
// Returns true if this was an OPTIONS request (and response was sent).
// Returns false if this is a normal request (continue processing).
//
// Example:
//
//	if response.HandleCORS(w, r) {
//	  return  // OPTIONS request handled
//	}
//	// Continue with normal request processing
func HandleCORS(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodOptions {
		SetCORSHeaders(w)
		w.WriteHeader(http.StatusNoContent)
		return true
	}
	return false
}

// Validation Helpers
// These functions help build validation errors.

// BuildValidationError creates a validation error for a single field.
//
// Example:
//
//	err := response.BuildValidationError("email", "", "Email is required")
func BuildValidationError(field string, rejectedValue interface{}, message string) types.ValidationError {
	return types.ValidationError{
		Field:         field,
		RejectedValue: rejectedValue,
		Message:       message,
	}
}

// ValidateRequired checks if a required string field is not empty.
// If validation fails, an error is appended to the errors slice.
//
// Example:
//
//	var errors []types.ValidationError
//	response.ValidateRequired("name", user.Name, &errors)
//	response.ValidateRequired("email", user.Email, &errors)
func ValidateRequired(field, value string, errors *[]types.ValidationError) {
	if strings.TrimSpace(value) == "" {
		*errors = append(*errors, BuildValidationError(field, value, field+" is required"))
	}
}

// ValidateEmail validates email format with a basic check.
// This is a simple validation - for production, use a proper email validation library.
//
// Example:
//
//	var errors []types.ValidationError
//	response.ValidateEmail("email", user.Email, &errors)
func ValidateEmail(field, email string, errors *[]types.ValidationError) {
	// Basic check: must contain @ and a dot
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		*errors = append(*errors, BuildValidationError(field, email, "Invalid email format"))
	}
}

// ValidateMinLength validates that a string meets minimum length requirement.
//
// Example:
//
//	var errors []types.ValidationError
//	response.ValidateMinLength("password", user.Password, 8, &errors)
func ValidateMinLength(field, value string, minLength int, errors *[]types.ValidationError) {
	if len(value) < minLength {
		*errors = append(*errors, BuildValidationError(
			field,
			value,
			fmt.Sprintf("%s must be at least %d characters", field, minLength),
		))
	}
}

// ValidateMaxLength validates that a string doesn't exceed maximum length.
//
// Example:
//
//	var errors []types.ValidationError
//	response.ValidateMaxLength("bio", user.Bio, 500, &errors)
func ValidateMaxLength(field, value string, maxLength int, errors *[]types.ValidationError) {
	if len(value) > maxLength {
		*errors = append(*errors, BuildValidationError(
			field,
			value,
			fmt.Sprintf("%s must not exceed %d characters", field, maxLength),
		))
	}
}

// ValidateRange validates that a numeric value is within a specified range.
//
// Example:
//
//	var errors []types.ValidationError
//	response.ValidateRange("age", user.Age, 18, 120, &errors)
func ValidateRange(field string, value, min, max int, errors *[]types.ValidationError) {
	if value < min || value > max {
		*errors = append(*errors, BuildValidationError(
			field,
			value,
			fmt.Sprintf("%s must be between %d and %d", field, min, max),
		))
	}
}
