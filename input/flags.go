// Package input provides functionality for handling user input from various sources.
//
// It manages command-line flags, file reading, and stdin input gathering.
// The package ensures input is properly validated, processed, and prepared
// for use in resume generation, handling edge cases and providing clear
// error messages when needed.
package input

import (
	"flag"
	"os"
)

// Flags represents the command-line flags accepted by the application.
// It provides a structured way to access the flag values throughout the program.
type Flags struct {
	// SourcePath holds the path to an optional existing resume file.
	// If provided, this resume will be used as a basis for generation.
	SourcePath string

	// OutputPath holds the path where the generated resume will be written.
	// If not provided, a default path will be used.
	OutputPath string
}

// ParseFlags parses the command-line flags from os.Args and returns the results.
// This is a convenience wrapper around ParseFlagsWithArgs using the program's
// actual command-line arguments.
//
// Returns:
//   - Flags: The parsed flag values
//   - error: Any error that occurred during parsing
//
// Example:
//
//	flags, err := input.ParseFlags()
//	if err != nil {
//	    log.Fatalf("Error parsing flags: %v", err)
//	}
func ParseFlags() (Flags, error) {
	return ParseFlagsWithArgs(os.Args[1:])
}

// ParseFlagsWithArgs parses the given arguments to extract flag values.
// This function allows parsing arbitrary string slices instead of using os.Args,
// which makes it particularly useful for testing.
//
// Parameters:
//   - args: The command-line arguments to parse (excluding the program name)
//
// Returns:
//   - Flags: The parsed flag values
//   - error: Any error that occurred during parsing
//
// Example:
//
//	testArgs := []string{"-source", "my_resume.md", "-output", "new_resume.md"}
//	flags, err := input.ParseFlagsWithArgs(testArgs)
func ParseFlagsWithArgs(args []string) (Flags, error) {
	var flags Flags
	
	// Create a new flag set
	fs := flag.NewFlagSet("resumake", flag.ContinueOnError)
	
	// Define the source flag
	sourcePath := fs.String("source", "", "Optional path to existing resume file (txt or md)")
	
	// Define the output flag
	outputPath := fs.String("output", "", "Path for the output resume file (default: resume_out.md)")
	
	// Parse the flags
	err := fs.Parse(args)
	if err != nil {
		return flags, err
	}
	
	// Set the flags struct values
	flags.SourcePath = *sourcePath
	flags.OutputPath = *outputPath
	
	return flags, nil
}