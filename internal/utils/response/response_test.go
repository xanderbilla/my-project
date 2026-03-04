package response

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xanderbilla/my-project/internal/constants"
	"github.com/xanderbilla/my-project/internal/types"
)

// Test OK response
func TestOK(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	testData := map[string]string{"message": "test data"}
	OK(rr, req, "Test success", testData)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("OK should return 200, got %v", status)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("OK should set Content-Type to application/json, got %v", contentType)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Response should have success=true")
	}

	if message, ok := response["message"].(string); !ok || message != "Test success" {
		t.Errorf("Response message should be 'Test success', got '%s'", message)
	}
}

// Test Created response
func TestCreated(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	testData := map[string]string{"id": "123"}
	Created(rr, req, "Resource created", testData)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("Created should return 201, got %v", status)
	}
}

// Test NoContent response
func TestNoContent(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/test", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	NoContent(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("NoContent should return 204, got %v", status)
	}
}

// Test BadRequest error
func TestBadRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	BadRequest(rr, req, "INVALID_FORMAT", "Invalid request", "Bad request details")

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("BadRequest should return 400, got %v", status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := response["success"].(bool); ok && success {
		t.Error("Error response should have success=false")
	}

	if errorInfo, ok := response["error"].(map[string]interface{}); !ok {
		t.Error("Response should have error object")
	} else {
		if errType, ok := errorInfo["type"].(string); !ok || errType != "BAD_REQUEST" {
			t.Errorf("Error type should be BAD_REQUEST, got %v", errType)
		}
	}
}

// Test NotFound error
func TestNotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test/123", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	NotFound(rr, req, "User", "123")

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("NotFound should return 404, got %v", status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if errorInfo, ok := response["error"].(map[string]interface{}); !ok {
		t.Error("Response should have error object")
	} else {
		if errType, ok := errorInfo["type"].(string); !ok || errType != "RESOURCE_NOT_FOUND" {
			t.Errorf("Error type should be RESOURCE_NOT_FOUND, got %v", errType)
		}

		if context, ok := errorInfo["context"].(map[string]interface{}); !ok {
			t.Error("Error should have context")
		} else {
			if id, ok := context["id"].(string); !ok || id != "123" {
				t.Errorf("Context should have id=123, got %v", id)
			}
		}
	}
}

// Test InternalError
func TestInternalError(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	testErr := errors.New("database connection failed")
	InternalError(rr, req, testErr)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("InternalError should return 500, got %v", status)
	}
}

// Test ValidationError
func TestValidationError(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	errors := []types.ValidationError{
		{Field: "email", RejectedValue: "not-email", Message: "Invalid email format"},
		{Field: "age", RejectedValue: 10, Message: "Age must be at least 18"},
	}

	ValidationError(rr, req, errors)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("ValidationError should return 422, got %v", status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if errorInfo, ok := response["error"].(map[string]interface{}); !ok {
		t.Error("Response should have error object")
	} else {
		if validationErrors, ok := errorInfo["validationErrors"].([]interface{}); !ok {
			t.Error("Error should have validationErrors array")
		} else if len(validationErrors) != 2 {
			t.Errorf("Should have 2 validation errors, got %d", len(validationErrors))
		}
	}
}

// Test DecodeJSON helper
func TestDecodeJSON_Valid(t *testing.T) {
	data := map[string]string{"name": "Test", "email": "test@example.com"}
	body, _ := json.Marshal(data)

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]string
	err := DecodeJSON(req, &result)

	if err != nil {
		t.Errorf("DecodeJSON should not return error for valid JSON, got: %v", err)
	}

	if result["name"] != "Test" {
		t.Errorf("Expected name='Test', got '%s'", result["name"])
	}
}

// Test DecodeJSON helper with invalid JSON
func TestDecodeJSON_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewReader([]byte("{invalid}")))
	req.Header.Set("Content-Type", "application/json")

	var result map[string]string
	err := DecodeJSON(req, &result)

	if err == nil {
		t.Error("DecodeJSON should return error for invalid JSON")
	}
}

// Test GetQueryParam helper
func TestGetQueryParam(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?name=John&age=30", nil)

	name := GetQueryParam(req, "name", "default")
	if name != "John" {
		t.Errorf("Expected name='John', got '%s'", name)
	}

	missing := GetQueryParam(req, "missing", "default")
	if missing != "default" {
		t.Errorf("Expected default value, got '%s'", missing)
	}
}

// Test GetQueryParamInt helper
func TestGetQueryParamInt(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?page=5&limit=20", nil)

	page := GetQueryParamInt(req, "page", 1)
	if page != 5 {
		t.Errorf("Expected page=5, got %d", page)
	}

	missing := GetQueryParamInt(req, "missing", 10)
	if missing != 10 {
		t.Errorf("Expected default value 10, got %d", missing)
	}
}

// Test GetQueryParamBool helper
func TestGetQueryParamBool(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?active=true&deleted=false", nil)

	active := GetQueryParamBool(req, "active")
	if !active {
		t.Error("Expected active=true")
	}

	deleted := GetQueryParamBool(req, "deleted")
	if deleted {
		t.Error("Expected deleted=false")
	}

	missing := GetQueryParamBool(req, "missing")
	if missing {
		t.Error("Expected missing parameter to return false")
	}
}

// Test ValidateRequired helper
func TestValidateRequired(t *testing.T) {
	errors := make([]types.ValidationError, 0)

	// Test with empty string
	ValidateRequired("name", "", &errors)
	if len(errors) != 1 {
		t.Error("ValidateRequired should add error for empty string")
	}

	// Test with non-empty string
	errors = make([]types.ValidationError, 0)
	ValidateRequired("name", "John", &errors)
	if len(errors) != 0 {
		t.Error("ValidateRequired should not add error for non-empty string")
	}
}

// Test ValidateEmail helper
func TestValidateEmail(t *testing.T) {
	errors := make([]types.ValidationError, 0)

	// Test valid email
	ValidateEmail("email", "test@example.com", &errors)
	if len(errors) != 0 {
		t.Error("ValidateEmail should not add error for valid email")
	}

	// Test invalid email
	errors = make([]types.ValidationError, 0)
	ValidateEmail("email", "not-an-email", &errors)
	if len(errors) == 0 {
		t.Error("ValidateEmail should add error for invalid email")
	}

	// Test empty email
	errors = make([]types.ValidationError, 0)
	ValidateEmail("email", "", &errors)
	if len(errors) == 0 {
		t.Error("ValidateEmail should add error for empty email")
	}
}

// Test ValidateMinLength helper
func TestValidateMinLength(t *testing.T) {
	errors := make([]types.ValidationError, 0)

	// Test string meeting minimum length
	ValidateMinLength("password", "12345", 5, &errors)
	if len(errors) != 0 {
		t.Error("ValidateMinLength should not add error when length is sufficient")
	}

	// Test string below minimum length
	errors = make([]types.ValidationError, 0)
	ValidateMinLength("password", "123", 5, &errors)
	if len(errors) == 0 {
		t.Error("ValidateMinLength should add error when length is insufficient")
	}
}

// Test ValidateMaxLength helper
func TestValidateMaxLength(t *testing.T) {
	errors := make([]types.ValidationError, 0)

	// Test string within maximum length
	ValidateMaxLength("name", "John", 10, &errors)
	if len(errors) != 0 {
		t.Error("ValidateMaxLength should not add error when length is within limit")
	}

	// Test string exceeding maximum length
	errors = make([]types.ValidationError, 0)
	ValidateMaxLength("name", "VeryLongNameThatExceedsLimit", 10, &errors)
	if len(errors) == 0 {
		t.Error("ValidateMaxLength should add error when length exceeds limit")
	}
}

// Test ValidateRange helper
func TestValidateRange(t *testing.T) {
	errors := make([]types.ValidationError, 0)

	// Test value within range
	ValidateRange("age", 25, 18, 120, &errors)
	if len(errors) != 0 {
		t.Error("ValidateRange should not add error when value is within range")
	}

	// Test value below minimum
	errors = make([]types.ValidationError, 0)
	ValidateRange("age", 10, 18, 120, &errors)
	if len(errors) == 0 {
		t.Error("ValidateRange should add error when value is below minimum")
	}

	// Test value above maximum
	errors = make([]types.ValidationError, 0)
	ValidateRange("age", 150, 18, 120, &errors)
	if len(errors) == 0 {
		t.Error("ValidateRange should add error when value exceeds maximum")
	}
}

// Test GetPaginationParams helper
func TestGetPaginationParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?page=2&pageSize=50", nil)

	params := GetPaginationParams(req)

	if params.Page != 2 {
		t.Errorf("Expected page=2, got %d", params.Page)
	}

	if params.PageSize != 50 {
		t.Errorf("Expected pageSize=50, got %d", params.PageSize)
	}
}

// Test GetPaginationParams with defaults
func TestGetPaginationParams_Defaults(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test", nil)

	params := GetPaginationParams(req)

	if params.Page != 1 {
		t.Errorf("Expected default page=1, got %d", params.Page)
	}

	if params.PageSize != 20 {
		t.Errorf("Expected default pageSize=20, got %d", params.PageSize)
	}
}

// Test GetPaginationParams with invalid values
func TestGetPaginationParams_Invalid(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/test?page=0&pageSize=999", nil)

	params := GetPaginationParams(req)

	if params.Page != 1 {
		t.Errorf("Expected page=1 for invalid page value, got %d", params.Page)
	}

	if params.PageSize != 20 {
		t.Errorf("Expected pageSize=20 (default) when exceeding max, got %d", params.PageSize)
	}
}
