package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xanderbilla/my-project/internal/constants"
	"github.com/xanderbilla/my-project/internal/types"
)

// TestGetAllUsers tests the GET /api/users endpoint
func TestGetAllUsers(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	getAllUsers(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getAllUsers should return 200, got %v", status)
	}

	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("getAllUsers should set Content-Type to application/json, got %v", contentType)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Response should have success=true")
	}

	if data, ok := response["data"].([]interface{}); !ok || len(data) == 0 {
		t.Error("Response should have data array with users")
	}
}

// TestGetUserByID_ValidID tests GET /api/users/{id} with valid ID
func TestGetUserByID_ValidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users/1", nil)
	req.SetPathValue("user_id", "1")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	getUserByID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("getUserByID with valid ID should return 200, got %v", status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Response should have success=true")
	}
}

// TestGetUserByID_InvalidID tests GET /api/users/{id} with invalid ID
func TestGetUserByID_InvalidID(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/users/999", nil)
	req.SetPathValue("user_id", "999")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	getUserByID(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("getUserByID with invalid ID should return 404, got %v", status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Error("Response should have success=false")
	}
}

// TestCreateUser_Valid tests POST /api/users with valid data
func TestCreateUser_Valid(t *testing.T) {
	user := types.User{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	createUser(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("createUser with valid data should return 201, got %v", status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || !success {
		t.Error("Response should have success=true")
	}
}

// TestCreateUser_InvalidJSON tests POST /api/users with invalid JSON
func TestCreateUser_InvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader([]byte("{invalid json}")))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	createUser(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("createUser with invalid JSON should return 400, got %v", status)
	}
}

// TestCreateUser_ValidationErrors tests POST /api/users with validation errors
func TestCreateUser_ValidationErrors(t *testing.T) {
	// Missing required fields
	user := types.User{
		Name: "", // Empty name should fail validation
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	createUser(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("createUser with validation errors should return 422, got %v", status)
	}

	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if success, ok := response["success"].(bool); !ok || success {
		t.Error("Response should have success=false")
	}

	// Check for error object with validationErrors
	if errorObj, ok := response["error"].(map[string]interface{}); !ok {
		t.Error("Response should have error object")
	} else if validationErrors, ok := errorObj["validationErrors"].([]interface{}); !ok || len(validationErrors) == 0 {
		t.Error("Error object should have validationErrors array")
	}
}

// TestUpdateUser_Valid tests PATCH /api/users/{id} with valid data
func TestUpdateUser_Valid(t *testing.T) {
	user := types.User{
		Name: "Updated User",
		Age:  30,
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPatch, "/api/users/1", bytes.NewReader(body))
	req.SetPathValue("user_id", "1")
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	updateUser(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("updateUser with valid data should return 200, got %v", status)
	}
}

// TestUpdateUser_NotFound tests PATCH /api/users/{id} with non-existent ID
func TestUpdateUser_NotFound(t *testing.T) {
	user := types.User{
		Name: "Updated User",
	}

	body, _ := json.Marshal(user)
	req := httptest.NewRequest(http.MethodPatch, "/api/users/999", bytes.NewReader(body))
	req.SetPathValue("user_id", "999")
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	updateUser(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("updateUser with invalid ID should return 404, got %v", status)
	}
}

// TestDeleteUser_Valid tests DELETE /api/users/{id} with valid ID
func TestDeleteUser_Valid(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/users/1", nil)
	req.SetPathValue("user_id", "1")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	deleteUser(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("deleteUser with valid ID should return 204, got %v", status)
	}
}

// TestDeleteUser_NotFound tests DELETE /api/users/{id} with non-existent ID
func TestDeleteUser_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/api/users/999", nil)
	req.SetPathValue("user_id", "999")
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	deleteUser(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("deleteUser with invalid ID should return 404, got %v", status)
	}
}

// TestHandleHealth tests GET /health
func TestHandleHealth(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rr := httptest.NewRecorder()

	handleHealth(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handleHealth should return 200, got %v", status)
	}

	var response map[string]interface{}
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatal("expected 'data' field in response")
	}

	status, ok := data["status"].(string)
	if !ok || status != "healthy" {
		t.Errorf("expected status 'healthy', got %v", status)
	}
}

// TestHandleFavicon tests GET /favicon.ico
func TestHandleFavicon(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/favicon.ico", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	handleFavicon(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handleFavicon should return 204, got %v", status)
	}

	if rr.Body.Len() != 0 {
		t.Error("handleFavicon should return empty body")
	}
}
