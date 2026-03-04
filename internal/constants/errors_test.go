package constants

import "testing"

// TestErrorTypeConstants checks that error type constants are properly defined
func TestErrorTypeConstants(t *testing.T) {
	testCases := []struct {
		name     string
		constant string
	}{
		{"ErrorTypeValidation", ErrorTypeValidation},
		{"ErrorTypeAuthentication", ErrorTypeAuthentication},
		{"ErrorTypeAuthorization", ErrorTypeAuthorization},
		{"ErrorTypeNotFound", ErrorTypeNotFound},
		{"ErrorTypeConflict", ErrorTypeConflict},
		{"ErrorTypeRateLimit", ErrorTypeRateLimit},
		{"ErrorTypeSystem", ErrorTypeSystem},
		{"ErrorTypeBadRequest", ErrorTypeBadRequest},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.constant == "" {
				t.Errorf("%s should not be empty", tc.name)
			}
		})
	}
}

// TestErrorCodeConstants checks that error code constants are properly defined
func TestErrorCodeConstants(t *testing.T) {
	testCases := []struct {
		name     string
		constant string
	}{
		{"ErrCodeUserNotFound", ErrCodeUserNotFound},
		{"ErrCodeUserAlreadyExists", ErrCodeUserAlreadyExists},
		{"ErrCodeInvalidUserData", ErrCodeInvalidUserData},
		{"ErrCodeInvalidRequestBody", ErrCodeInvalidRequestBody},
		{"ErrCodeInvalidJSON", ErrCodeInvalidJSON},
		{"ErrCodeMissingField", ErrCodeMissingField},
		{"ErrCodeInvalidEmail", ErrCodeInvalidEmail},
		{"ErrCodeDatabaseError", ErrCodeDatabaseError},
		{"ErrCodeInternalError", ErrCodeInternalError},
		{"ErrCodeServiceUnavailable", ErrCodeServiceUnavailable},
		{"ErrCodeInvalidContentType", ErrCodeInvalidContentType},
		{"ErrCodeMethodNotAllowed", ErrCodeMethodNotAllowed},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.constant == "" {
				t.Errorf("%s should not be empty", tc.name)
			}
		})
	}
}

// TestAPIVersion checks that API version is defined
func TestAPIVersion(t *testing.T) {
	if APIVersion == "" {
		t.Error("APIVersion should not be empty")
	}
}
