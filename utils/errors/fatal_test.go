package errors

import (
	"bytes"
	stderrors "errors"
	"os"
	"testing"
)

func TestHandleErrorFatal(t *testing.T) {
	// Since HandleErrorFatal calls log.Fatal which exits the program,
	// we need to capture the log output without actually exiting.
	// We'll redirect log output to our buffer and mock the exit function.
	
	// Save original log output and restore after test
	originalOutput := logOutput
	t.Cleanup(func() {
		logOutput = originalOutput
		exitFunc = os.Exit
	})
	
	var logBuffer bytes.Buffer
	logOutput = &logBuffer
	
	// Mock exit function
	exitCalled := false
	exitCode := 0
	exitFunc = func(code int) {
		exitCalled = true
		exitCode = code
		// Don't actually exit in tests
	}
	
	// Test case
	testErr := stderrors.New("connection failed")
	HandleErrorFatal("connecting to server", testErr)
	
	// Verify exit was called with code 1
	if !exitCalled {
		t.Error("Exit function was not called")
	}
	if exitCode != 1 {
		t.Errorf("Expected exit code 1, got %d", exitCode)
	}
	
	// Verify log message
	logMsg := logBuffer.String()
	expectedMsg := "Error connecting to server: connection failed\n"
	if logMsg != expectedMsg {
		t.Errorf("Expected log message %q, got %q", expectedMsg, logMsg)
	}
}