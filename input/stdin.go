package input

import (
	"fmt"
	"io"
	"os"
)

// ReadFromReader reads input from the provided reader until EOF.
// It displays prompts through the provided writer.
// Returns the input as a string and any error that occurred.
func ReadFromReader(reader io.Reader, writer io.Writer) (string, error) {
	// Display prompt with instructions
	fmt.Fprintln(writer, "Enter your raw professional history below.")
	fmt.Fprintln(writer, "Press Ctrl+D (Unix) or Ctrl+Z then Enter (Windows) when finished:")
	
	// Read all content from the reader
	inputBytes, err := io.ReadAll(reader)
	if err != nil {
		return "", fmt.Errorf("error reading input: %w", err)
	}
	
	// Convert to string and return
	return string(inputBytes), nil
}

// ReadFromStdin is a convenience wrapper that reads from standard input
// and displays prompts on standard output.
// Returns the input as a string and any error that occurred.
func ReadFromStdin() (string, error) {
	return ReadFromReader(os.Stdin, os.Stdout)
}