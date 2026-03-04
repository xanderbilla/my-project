package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xanderbilla/my-project/internal/constants"
)

// Test RequestID middleware with existing X-Request-ID header
func TestRequestID_WithExistingHeader(t *testing.T) {
	// Create a handler that checks if request ID is in context
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(constants.RequestIDKey)
		if requestID == nil {
			t.Error("Request ID should be in context")
			return
		}

		if id, ok := requestID.(string); !ok {
			t.Error("Request ID should be a string")
		} else if id != "test-request-id-123" {
			t.Errorf("Request ID should be 'test-request-id-123', got '%s'", id)
		}

		w.WriteHeader(http.StatusOK)
	})

	// Wrap handler with RequestID middleware
	wrappedHandler := RequestID(handler)

	// Create request with X-Request-ID header
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("X-Request-ID", "test-request-id-123")
	rr := httptest.NewRecorder()

	// Call the wrapped handler
	wrappedHandler.ServeHTTP(rr, req)

	// Check that X-Request-ID header is set in response
	if responseID := rr.Header().Get("X-Request-ID"); responseID != "test-request-id-123" {
		t.Errorf("Response X-Request-ID should be 'test-request-id-123', got '%s'", responseID)
	}
}

// Test RequestID middleware without X-Request-ID header (should generate UUID)
func TestRequestID_WithoutHeader(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Context().Value(constants.RequestIDKey)
		if requestID == nil {
			t.Error("Request ID should be in context")
			return
		}

		if id, ok := requestID.(string); !ok {
			t.Error("Request ID should be a string")
		} else if id == "" {
			t.Error("Request ID should not be empty")
		} else if len(id) < 10 {
			t.Error("Request ID should be a UUID (at least 10 characters)")
		}

		w.WriteHeader(http.StatusOK)
	})

	wrappedHandler := RequestID(handler)
	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	// Check that X-Request-ID header is set in response
	responseID := rr.Header().Get("X-Request-ID")
	if responseID == "" {
		t.Error("Response should have X-Request-ID header")
	}

	if len(responseID) < 10 {
		t.Error("Generated request ID should be a UUID")
	}
}

// Test Logger middleware
func TestLogger(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("test response"))
	})

	wrappedHandler := Logger(handler)

	req := httptest.NewRequest(http.MethodGet, "/test?foo=bar", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Handler should return 200, got %v", status)
	}

	if body := rr.Body.String(); body != "test response" {
		t.Errorf("Handler should return 'test response', got '%s'", body)
	}
}

// Test Logger middleware with different status codes
func TestLogger_StatusCodes(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
	}{
		{"Success", http.StatusOK},
		{"Created", http.StatusCreated},
		{"NoContent", http.StatusNoContent},
		{"BadRequest", http.StatusBadRequest},
		{"NotFound", http.StatusNotFound},
		{"InternalServerError", http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.statusCode)
			})

			wrappedHandler := Logger(handler)
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
			req = req.WithContext(ctx)
			rr := httptest.NewRecorder()

			wrappedHandler.ServeHTTP(rr, req)

			if status := rr.Code; status != tc.statusCode {
				t.Errorf("Handler should return %v, got %v", tc.statusCode, status)
			}
		})
	}
}

// Test Recovery middleware with panic
func TestRecovery_WithPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})

	wrappedHandler := Recovery(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// This should not panic - the middleware should catch it
	wrappedHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("Recovery should return 500 on panic, got %v", status)
	}

	// Check that response is JSON
	if contentType := rr.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Recovery should return JSON, got Content-Type: %v", contentType)
	}
}

// Test Recovery middleware without panic
func TestRecovery_NoPanic(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	wrappedHandler := Recovery(handler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	ctx := context.WithValue(req.Context(), constants.RequestIDKey, "test-request-id")
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Recovery should pass through normal requests, got %v", status)
	}

	if body := rr.Body.String(); body != "success" {
		t.Errorf("Recovery should pass through response body, got '%s'", body)
	}
}

// Test middleware chain (Recovery -> RequestID -> Logger)
func TestMiddlewareChain(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that request ID is in context
		requestID := r.Context().Value(constants.RequestIDKey)
		if requestID == nil {
			t.Error("Request ID should be in context")
		}

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})

	// Chain middlewares: Recovery -> RequestID -> Logger -> Handler
	wrappedHandler := Recovery(RequestID(Logger(handler)))

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	wrappedHandler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Middleware chain should work correctly, got status %v", status)
	}

	// Check that X-Request-ID header is set
	if responseID := rr.Header().Get("X-Request-ID"); responseID == "" {
		t.Error("Middleware chain should set X-Request-ID header")
	}
}
