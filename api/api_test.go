package api

import (
	"os"
	"testing"
)

func TestGetAPIKey(t *testing.T) {
	// Save original environment to restore after tests
	originalAPIKey := os.Getenv("GEMINI_API_KEY")
	defer os.Setenv("GEMINI_API_KEY", originalAPIKey)

	// Test case 1: API key is set correctly
	t.Run("API key is set", func(t *testing.T) {
		expected := "test-api-key-123"
		os.Setenv("GEMINI_API_KEY", expected)
		
		actual, err := GetAPIKey()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		if actual != expected {
			t.Errorf("Expected API key %q, got %q", expected, actual)
		}
	})

	// Test case 2: API key is missing
	t.Run("API key is missing", func(t *testing.T) {
		os.Unsetenv("GEMINI_API_KEY")
		
		_, err := GetAPIKey()
		if err == nil {
			t.Error("Expected error when API key is missing, got nil")
		}
	})
	
	// Test case 3: API key is empty
	t.Run("API key is empty", func(t *testing.T) {
		os.Setenv("GEMINI_API_KEY", "")
		
		_, err := GetAPIKey()
		if err == nil {
			t.Error("Expected error when API key is empty, got nil")
		}
	})
}