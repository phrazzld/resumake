package errors

import (
	stderrors "errors"
	"testing"
)

func TestFormatErrorMessage(t *testing.T) {
	testCases := []struct {
		name     string
		context  string
		err      error
		expected string
	}{
		{
			name:     "basic error with context",
			context:  "loading file",
			err:      stderrors.New("file not found"),
			expected: "Error loading file: file not found",
		},
		{
			name:     "context with colon",
			context:  "parsing input at line 42:",
			err:      stderrors.New("unexpected token"),
			expected: "Error parsing input at line 42: unexpected token",
		},
		{
			name:     "empty context",
			context:  "",
			err:      stderrors.New("network timeout"),
			expected: "Error: network timeout",
		},
		{
			name:     "nil error",
			context:  "processing data",
			err:      nil,
			expected: "Error processing data: <nil>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := FormatErrorMessage(tc.context, tc.err)
			if result != tc.expected {
				t.Errorf("Expected %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestWrapError(t *testing.T) {
	testCases := []struct {
		name            string
		context         string
		originalErr     error
		expectedMessage string
	}{
		{
			name:            "basic wrapping",
			context:         "loading configuration",
			originalErr:     stderrors.New("file not found"),
			expectedMessage: "Error loading configuration: file not found",
		},
		{
			name:            "wrapping with empty context",
			context:         "",
			originalErr:     stderrors.New("connection timeout"),
			expectedMessage: "Error: connection timeout",
		},
		{
			name:            "wrapping nil error",
			context:         "processing data",
			originalErr:     nil,
			expectedMessage: "Error processing data: <nil>",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			wrappedErr := WrapError(tc.context, tc.originalErr)
			
			// Check error message
			if wrappedErr == nil {
				if tc.originalErr != nil {
					t.Fatal("WrapError returned nil for non-nil error")
				}
			} else if wrappedErr.Error() != tc.expectedMessage {
				t.Errorf("Expected message %q, got %q", tc.expectedMessage, wrappedErr.Error())
			}
			
			// Check unwrapping (when original error isn't nil)
			if tc.originalErr != nil && wrappedErr != nil {
				unwrapped := stderrors.Unwrap(wrappedErr)
				if unwrapped != tc.originalErr {
					t.Errorf("Unwrapped error doesn't match original. Got %v, expected %v", unwrapped, tc.originalErr)
				}
			}
		})
	}
}