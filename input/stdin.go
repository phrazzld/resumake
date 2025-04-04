package input

import (
	"fmt"
	"io"
	"os"
)

// ReadFromReader reads input from the provided reader until EOF (end of file).
// It displays user prompts through the provided writer, instructing the user on
// how to terminate input. This function is useful for collecting multi-line
// textual input from various sources.
//
// Parameters:
//   - reader: The source to read input from (e.g., stdin)
//   - writer: The destination to write prompts to (e.g., stdout)
//
// Returns:
//   - string: The collected input as a string
//   - error: Any error that occurred during reading
//
// Example:
//
//	content, err := input.ReadFromReader(strings.NewReader("test input"), io.Discard)
//	if err != nil {
//	    log.Fatalf("Error reading input: %v", err)
//	}
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

// ReadFromStdin is a convenience wrapper around ReadFromReader that reads from
// the system's standard input (os.Stdin) and displays prompts on standard output
// (os.Stdout). This is the primary function used for collecting user input in
// interactive terminal sessions.
//
// Returns:
//   - string: The collected input as a string
//   - error: Any error that occurred during reading
//
// Example:
//
//	content, err := input.ReadFromStdin()
//	if err != nil {
//	    log.Fatalf("Error reading from stdin: %v", err)
//	}
//	fmt.Printf("Read %d characters from stdin\n", len(content))
func ReadFromStdin() (string, error) {
	return ReadFromReader(os.Stdin, os.Stdout)
}