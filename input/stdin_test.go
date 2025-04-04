package input

import (
	"bytes"
	"strings"
	"testing"
)

func TestReadFromReader(t *testing.T) {
	// Test case 1: Successfully read input from reader
	t.Run("Read input from reader", func(t *testing.T) {
		// Setup test input reader
		testInput := "This is a test input.\nIt has multiple lines.\n"
		reader := strings.NewReader(testInput)
		
		// Setup output writer to capture prompts
		var outBuf bytes.Buffer
		
		// Call ReadFromReader with reader and writer
		input, err := ReadFromReader(reader, &outBuf)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		
		// Verify the returned input matches the mock input
		if input != testInput {
			t.Errorf("Expected input %q, got %q", testInput, input)
		}
		
		// Verify prompt was displayed in writer
		output := outBuf.String()
		if !strings.Contains(output, "Enter your raw professional history") {
			t.Errorf("Expected prompt to contain instructions, got: %q", output)
		}
	})
	
	// Test case 2: Empty input
	t.Run("Empty input", func(t *testing.T) {
		// Setup empty input reader
		reader := strings.NewReader("")
		
		// Create a discard writer for output
		var outBuf bytes.Buffer
		
		// Call ReadFromReader
		input, err := ReadFromReader(reader, &outBuf)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error for empty input, got %v", err)
		}
		
		// Verify an empty string is returned
		if input != "" {
			t.Errorf("Expected empty input, got %q", input)
		}
	})
	
	// Test case 3: Large input
	t.Run("Large input", func(t *testing.T) {
		// Generate large input (1000 lines)
		var largeInput strings.Builder
		for i := 0; i < 1000; i++ {
			largeInput.WriteString(strings.Repeat("a", 50) + "\n")
		}
		testInput := largeInput.String()
		
		// Setup reader with large input
		reader := strings.NewReader(testInput)
		
		// Create a discard writer for output
		var outBuf bytes.Buffer
		
		// Call ReadFromReader
		input, err := ReadFromReader(reader, &outBuf)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error for large input, got %v", err)
		}
		
		// Verify the entire input is correctly read
		if len(input) != len(testInput) {
			t.Errorf("Expected input length %d, got %d", len(testInput), len(input))
		}
		
		if input != testInput {
			t.Errorf("Large input was not correctly read")
		}
	})
	
	// Test case 4: Special characters
	t.Run("Input with special characters", func(t *testing.T) {
		// Setup input with special characters
		testInput := "Unicode: 😊 ñ é\nSymbols: © ® ™\nTabs:\t\t\t"
		reader := strings.NewReader(testInput)
		
		// Create a buffer for output
		var outBuf bytes.Buffer
		
		// Call ReadFromReader
		input, err := ReadFromReader(reader, &outBuf)
		
		// Verify no error occurred
		if err != nil {
			t.Errorf("Expected no error for input with special characters, got %v", err)
		}
		
		// Verify special characters are preserved
		if input != testInput {
			t.Errorf("Expected input %q, got %q", testInput, input)
		}
	})
}