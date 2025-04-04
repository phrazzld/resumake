// Package errors provides error handling utilities for the resumake application.
//
// It offers a consistent approach to error handling, formatting, and reporting
// throughout the application. For v0.1, the focus is on simple error handling
// using log.Fatal, with a design that allows for more sophisticated error
// handling strategies in future versions.
package errors

import (
	"fmt"
	"io"
	"log"
	"os"
)

// FormatErrorMessage creates a consistently formatted error message combining
// a context string with an error. This ensures uniform error message formatting
// throughout the application.
//
// Parameters:
//   - context: A string describing the context in which the error occurred
//   - err: The error that occurred
//
// Returns:
//   - string: A formatted error message combining the context and error
//
// Example:
//
//	msg := errors.FormatErrorMessage("parsing config", err)
//	// Returns: "Error parsing config: <error message>"
func FormatErrorMessage(context string, err error) string {
	// If context is empty, use a generic prefix
	if context == "" {
		return fmt.Sprintf("Error: %v", err)
	}
	
	// Check if the context already ends with a colon to avoid double colons
	if len(context) > 0 && context[len(context)-1] == ':' {
		return fmt.Sprintf("Error %s %v", context, err)
	}
	
	// Return formatted error message with context
	return fmt.Sprintf("Error %s: %v", context, err)
}

// wrappedError implements the error interface and wraps an underlying error
// with additional context while supporting the Unwrap method.
type wrappedError struct {
	msg string
	err error
}

// Error returns the formatted error message.
func (w *wrappedError) Error() string {
	return w.msg
}

// Unwrap returns the original error.
func (w *wrappedError) Unwrap() error {
	return w.err
}

// WrapError wraps an error with additional context, creating a new error
// that both carries a descriptive message and maintains the original error
// for use with errors.Is/As/Unwrap.
//
// Parameters:
//   - context: A string describing the context in which the error occurred
//   - err: The original error to wrap (can be nil)
//
// Returns:
//   - error: A new error wrapping the original with additional context
//
// Example:
//
//	originalErr := io.EOF
//	wrappedErr := errors.WrapError("reading config file", originalErr)
//	// wrappedErr.Error() returns "Error reading config file: EOF"
//	// errors.Unwrap(wrappedErr) returns io.EOF
func WrapError(context string, err error) error {
	// If the original error is nil, still return a descriptive error
	// to avoid null pointer issues downstream
	if err == nil {
		return &wrappedError{
			msg: FormatErrorMessage(context, nil),
			err: nil,
		}
	}
	
	return &wrappedError{
		msg: FormatErrorMessage(context, err),
		err: err,
	}
}

// Variables that can be overridden for testing
var (
	// logOutput is the destination for log messages (defaults to os.Stderr)
	logOutput io.Writer = os.Stderr
	
	// exitFunc is the function called to exit the program (defaults to os.Exit)
	exitFunc = os.Exit
)

// HandleErrorFatal logs an error message with context and exits the program.
// This provides a consistent way to handle fatal errors throughout the application.
//
// Parameters:
//   - context: A string describing the context in which the error occurred
//   - err: The error that caused the program to exit
//
// Example:
//
//	if err != nil {
//	    errors.HandleErrorFatal("reading configuration file", err)
//	}
//
// Note: This function never returns as it calls os.Exit(1).
func HandleErrorFatal(context string, err error) {
	// Create a custom logger that writes to the configured output
	logger := log.New(logOutput, "", 0)
	
	// Format the error message with context
	message := FormatErrorMessage(context, err)
	
	// Log the error and exit
	logger.Println(message)
	exitFunc(1)
}

// CheckErrorNil executes a handler function only if the error is not nil.
// This simplifies error handling code by removing repetitive nil checks.
//
// Parameters:
//   - err: The error to check
//   - handler: Function to call if the error is not nil
//
// Example:
//
//	err := doSomething()
//	errors.CheckErrorNil(err, func(e error) {
//	    log.Printf("Failed to do something: %v", e)
//	})
//
// A common use case is to combine with HandleErrorFatal:
//
//	err := doSomething()
//	errors.CheckErrorNil(err, func(e error) {
//	    errors.HandleErrorFatal("doing something", e)
//	})
func CheckErrorNil(err error, handler func(error)) {
	if err != nil {
		handler(err)
	}
}