package errors

import (
	stderrors "errors"
	"testing"
)

func TestCheckErrorNil(t *testing.T) {
	// Create a counter to track how many times handler is called
	handlerCallCount := 0
	
	// Create a simple handler that increments the counter
	handler := func(err error) {
		handlerCallCount++
	}
	
	// Test with nil error - handler should not be called
	CheckErrorNil(nil, handler)
	if handlerCallCount != 0 {
		t.Errorf("Handler was called for nil error, call count: %d", handlerCallCount)
	}
	
	// Test with non-nil error - handler should be called
	testError := stderrors.New("test error")
	CheckErrorNil(testError, handler)
	if handlerCallCount != 1 {
		t.Errorf("Handler was not called for non-nil error, call count: %d", handlerCallCount)
	}
	
	// Verify the error was passed correctly to the handler
	var capturedError error
	captureHandler := func(err error) {
		capturedError = err
	}
	
	CheckErrorNil(testError, captureHandler)
	if capturedError != testError {
		t.Errorf("Handler received wrong error. Got %v, expected %v", capturedError, testError)
	}
}