package input

import (
	"flag"
	"os"
)

// Flags represents the command-line flags for the application
type Flags struct {
	SourcePath string // Path to the optional source resume file
}

// ParseFlags parses the command-line flags and returns the results
func ParseFlags() (Flags, error) {
	return ParseFlagsWithArgs(os.Args[1:])
}

// ParseFlagsWithArgs parses the given arguments instead of os.Args,
// which is useful for testing
func ParseFlagsWithArgs(args []string) (Flags, error) {
	var flags Flags
	
	// Create a new flag set
	fs := flag.NewFlagSet("resumake", flag.ContinueOnError)
	
	// Define the source flag
	sourcePath := fs.String("source", "", "Optional path to existing resume file (txt or md)")
	
	// Parse the flags
	err := fs.Parse(args)
	if err != nil {
		return flags, err
	}
	
	// Set the flags struct values
	flags.SourcePath = *sourcePath
	
	return flags, nil
}