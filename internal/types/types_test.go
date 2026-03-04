package types

import (
	"errors"
	"net/http"
	"testing"
	"time"
)

// TestAppError_Error validates the Error() method
func TestAppError_Error(t *testing.T) {
	tests := []struct {
		name     string
		appError *AppError
		contains string
	}{
		{
			name: "ErrorWithoutOriginal",
			appError: &AppError{
				Code:             "USER_NOT_FOUND",
				DeveloperMessage: "User with ID 123 does not exist",
			},
			contains: "USER_NOT_FOUND: User with ID 123 does not exist",
		},
		{
			name: "ErrorWithOriginal",
			appError: &AppError{
				Code:             "DATABASE_ERROR",
				DeveloperMessage: "Failed to query database",
				OriginalError:    errors.New("connection timeout"),
			},
			contains: "DATABASE_ERROR: Failed to query database (original: connection timeout)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.appError.Error()
			if result != tt.contains {
				t.Errorf("expected %q, got %q", tt.contains, result)
			}
		})
	}
}

// TestAppError_Unwrap validates error unwrapping
func TestAppError_Unwrap(t *testing.T) {
	originalErr := errors.New("original error")
	appError := &AppError{
		Code:          "TEST_ERROR",
		OriginalError: originalErr,
	}

	unwrapped := appError.Unwrap()
	if unwrapped != originalErr {
		t.Errorf("expected unwrapped error to be original error")
	}

	// Test with no original error
	appErrorNoOriginal := &AppError{
		Code: "TEST_ERROR",
	}
	if appErrorNoOriginal.Unwrap() != nil {
		t.Error("expected nil when no original error")
	}
}

// TestAppError_ToErrorInfo validates conversion to ErrorInfo
func TestAppError_ToErrorInfo(t *testing.T) {
	appError := &AppError{
		Type:             "VALIDATION_ERROR",
		Code:             "INVALID_EMAIL",
		UserMessage:      "Email format is invalid",
		DeveloperMessage: "Email must be in format user@domain.com",
		Context: map[string]interface{}{
			"field": "email",
		},
		Retryable: false,
	}

	errorInfo := appError.ToErrorInfo()

	if errorInfo.Type != appError.Type {
		t.Errorf("expected Type %s, got %s", appError.Type, errorInfo.Type)
	}

	if errorInfo.Code != appError.Code {
		t.Errorf("expected Code %s, got %s", appError.Code, errorInfo.Code)
	}

	if errorInfo.UserMessage != appError.UserMessage {
		t.Errorf("expected UserMessage %s, got %s", appError.UserMessage, errorInfo.UserMessage)
	}

	if errorInfo.DeveloperMessage != appError.DeveloperMessage {
		t.Errorf("expected DeveloperMessage %s, got %s", appError.DeveloperMessage, errorInfo.DeveloperMessage)
	}

	if errorInfo.Retryable != appError.Retryable {
		t.Errorf("expected Retryable %v, got %v", appError.Retryable, errorInfo.Retryable)
	}
}

// TestNewValidationError validates validation error creation
func TestNewValidationError(t *testing.T) {
	err := NewValidationError("INVALID_EMAIL", "Email is invalid", "Email format must be user@domain.com")

	if err.Type != "VALIDATION_ERROR" {
		t.Errorf("expected Type VALIDATION_ERROR, got %s", err.Type)
	}

	if err.Code != "INVALID_EMAIL" {
		t.Errorf("expected Code INVALID_EMAIL, got %s", err.Code)
	}

	if err.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("expected StatusCode 422, got %d", err.StatusCode)
	}

	if err.Retryable {
		t.Error("validation errors should not be retryable")
	}
}

// TestNewNotFoundError validates not found error creation
func TestNewNotFoundError(t *testing.T) {
	err := NewNotFoundError("User", "123")

	if err.Type != "RESOURCE_NOT_FOUND" {
		t.Errorf("expected Type RESOURCE_NOT_FOUND, got %s", err.Type)
	}

	if err.Code != "User_NOT_FOUND" {
		t.Errorf("expected Code User_NOT_FOUND, got %s", err.Code)
	}

	if err.StatusCode != http.StatusNotFound {
		t.Errorf("expected StatusCode 404, got %d", err.StatusCode)
	}

	if err.Retryable {
		t.Error("not found errors should not be retryable")
	}

	// Check context
	if err.Context["resource"] != "User" {
		t.Errorf("expected resource User in context")
	}

	if err.Context["id"] != "123" {
		t.Errorf("expected id 123 in context")
	}
}

// TestNewConflictError validates conflict error creation
func TestNewConflictError(t *testing.T) {
	err := NewConflictError("USER_EXISTS", "User already exists", "User with email test@example.com already exists")

	if err.Type != "RESOURCE_CONFLICT" {
		t.Errorf("expected Type RESOURCE_CONFLICT, got %s", err.Type)
	}

	if err.StatusCode != http.StatusConflict {
		t.Errorf("expected StatusCode 409, got %d", err.StatusCode)
	}

	if err.Retryable {
		t.Error("conflict errors should not be retryable")
	}
}

// TestNewAuthenticationError validates authentication error creation
func TestNewAuthenticationError(t *testing.T) {
	err := NewAuthenticationError("Invalid credentials")

	if err.Type != "AUTHENTICATION_ERROR" {
		t.Errorf("expected Type AUTHENTICATION_ERROR, got %s", err.Type)
	}

	if err.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected StatusCode 401, got %d", err.StatusCode)
	}

	if err.Retryable {
		t.Error("authentication errors should not be retryable")
	}
}

// TestNewAuthorizationError validates authorization error creation
func TestNewAuthorizationError(t *testing.T) {
	err := NewAuthorizationError("delete", "User")

	if err.Type != "AUTHORIZATION_ERROR" {
		t.Errorf("expected Type AUTHORIZATION_ERROR, got %s", err.Type)
	}

	if err.StatusCode != http.StatusForbidden {
		t.Errorf("expected StatusCode 403, got %d", err.StatusCode)
	}

	if err.Retryable {
		t.Error("authorization errors should not be retryable")
	}
}

// TestNewBadRequestError validates bad request error creation
func TestNewBadRequestError(t *testing.T) {
	err := NewBadRequestError("INVALID_JSON", "Invalid JSON", "Failed to parse JSON body")

	if err.Type != "BAD_REQUEST" {
		t.Errorf("expected Type BAD_REQUEST, got %s", err.Type)
	}

	if err.StatusCode != http.StatusBadRequest {
		t.Errorf("expected StatusCode 400, got %d", err.StatusCode)
	}

	if err.Retryable {
		t.Error("bad request errors should not be retryable")
	}
}

// TestNewInternalError validates internal error creation
func TestNewInternalError(t *testing.T) {
	originalErr := errors.New("database connection failed")
	err := NewInternalError(originalErr)

	if err.Type != "SYSTEM_ERROR" {
		t.Errorf("expected Type SYSTEM_ERROR, got %s", err.Type)
	}

	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected StatusCode 500, got %d", err.StatusCode)
	}

	if !err.Retryable {
		t.Error("internal errors should be retryable")
	}

	if err.OriginalError != originalErr {
		t.Error("expected original error to be stored")
	}
}

// TestNewDatabaseError validates database error creation
func TestNewDatabaseError(t *testing.T) {
	originalErr := errors.New("connection timeout")
	err := NewDatabaseError("SELECT", originalErr)

	if err.Type != "SYSTEM_ERROR" {
		t.Errorf("expected Type SYSTEM_ERROR, got %s", err.Type)
	}

	if err.StatusCode != http.StatusInternalServerError {
		t.Errorf("expected StatusCode 500, got %d", err.StatusCode)
	}

	if !err.Retryable {
		t.Error("database errors should be retryable")
	}

	if err.Context["operation"] != "SELECT" {
		t.Error("expected operation in context")
	}
}

// TestNewRateLimitError validates rate limit error creation
func TestNewRateLimitError(t *testing.T) {
	err := NewRateLimitError(60)

	if err.Type != "RATE_LIMIT_EXCEEDED" {
		t.Errorf("expected Type RATE_LIMIT_EXCEEDED, got %s", err.Type)
	}

	if err.StatusCode != http.StatusTooManyRequests {
		t.Errorf("expected StatusCode 429, got %d", err.StatusCode)
	}

	if !err.Retryable {
		t.Error("rate limit errors should be retryable")
	}

	if err.Context["retryAfter"] != 60 {
		t.Error("expected retryAfter in context")
	}
}

// TestNewServiceUnavailableError validates service unavailable error creation
func TestNewServiceUnavailableError(t *testing.T) {
	err := NewServiceUnavailableError("PaymentService")

	if err.Type != "SYSTEM_ERROR" {
		t.Errorf("expected Type SYSTEM_ERROR, got %s", err.Type)
	}

	if err.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("expected StatusCode 503, got %d", err.StatusCode)
	}

	if !err.Retryable {
		t.Error("service unavailable errors should be retryable")
	}

	if err.Context["service"] != "PaymentService" {
		t.Error("expected service in context")
	}
}

// TestAPIResponse_Structure validates APIResponse fields
func TestAPIResponse_Structure(t *testing.T) {
	now := time.Now().UTC()
	response := APIResponse{
		Success:   true,
		Status:    200,
		Timestamp: now,
		RequestID: "test-123",
		Path:      "/api/users",
		Method:    "GET",
		Message:   "Success",
		Data:      map[string]string{"key": "value"},
	}

	if !response.Success {
		t.Error("expected Success to be true")
	}

	if response.Status != 200 {
		t.Errorf("expected Status 200, got %d", response.Status)
	}

	if !response.Timestamp.Equal(now) {
		t.Errorf("expected Timestamp %v, got %v", now, response.Timestamp)
	}

	if response.RequestID != "test-123" {
		t.Errorf("expected RequestID test-123, got %s", response.RequestID)
	}

	if response.Path != "/api/users" {
		t.Errorf("expected Path /api/users, got %s", response.Path)
	}

	if response.Method != "GET" {
		t.Errorf("expected Method GET, got %s", response.Method)
	}

	if response.Message != "Success" {
		t.Errorf("expected Message Success, got %s", response.Message)
	}

	if data, ok := response.Data.(map[string]string); !ok || data["key"] != "value" {
		t.Errorf("expected Data map with key=value, got %v", response.Data)
	}
}

// TestUser_Structure validates User struct
func TestUser_Structure(t *testing.T) {
	user := User{
		UserId: 1,
		Name:   "John Doe",
		Email:  "john@example.com",
		Age:    30,
	}

	if user.UserId != 1 {
		t.Errorf("expected UserId 1, got %d", user.UserId)
	}

	if user.Name != "John Doe" {
		t.Errorf("expected Name John Doe, got %s", user.Name)
	}

	if user.Email != "john@example.com" {
		t.Errorf("expected Email john@example.com, got %s", user.Email)
	}

	if user.Age != 30 {
		t.Errorf("expected Age 30, got %d", user.Age)
	}
}

// TestErrorInfo_Structure validates ErrorInfo struct
func TestErrorInfo_Structure(t *testing.T) {
	errorInfo := ErrorInfo{
		Type:             "VALIDATION_ERROR",
		Code:             "INVALID_EMAIL",
		UserMessage:      "Invalid email",
		DeveloperMessage: "Email format is incorrect",
		Context: map[string]interface{}{
			"field": "email",
		},
		Retryable:        false,
		ValidationErrors: []ValidationError{},
	}

	if errorInfo.Type != "VALIDATION_ERROR" {
		t.Errorf("expected Type VALIDATION_ERROR, got %s", errorInfo.Type)
	}

	if errorInfo.Code != "INVALID_EMAIL" {
		t.Errorf("expected Code INVALID_EMAIL, got %s", errorInfo.Code)
	}

	if errorInfo.UserMessage != "Invalid email" {
		t.Errorf("expected UserMessage 'Invalid email', got %s", errorInfo.UserMessage)
	}

	if errorInfo.DeveloperMessage != "Email format is incorrect" {
		t.Errorf("expected DeveloperMessage 'Email format is incorrect', got %s", errorInfo.DeveloperMessage)
	}

	if errorInfo.Context["field"] != "email" {
		t.Errorf("expected Context field=email, got %v", errorInfo.Context)
	}

	if errorInfo.Retryable {
		t.Error("expected Retryable to be false")
	}

	if len(errorInfo.ValidationErrors) != 0 {
		t.Errorf("expected empty ValidationErrors, got %d items", len(errorInfo.ValidationErrors))
	}
}

// TestValidationError_Structure validates ValidationError struct
func TestValidationError_Structure(t *testing.T) {
	valError := ValidationError{
		Field:         "email",
		Message:       "Email is required",
		RejectedValue: "",
	}

	if valError.Field != "email" {
		t.Errorf("expected Field email, got %s", valError.Field)
	}

	if valError.Message != "Email is required" {
		t.Errorf("expected Message 'Email is required', got %s", valError.Message)
	}

	if valError.RejectedValue != "" {
		t.Errorf("expected RejectedValue to be empty string, got %v", valError.RejectedValue)
	}
}
