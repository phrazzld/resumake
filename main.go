package main

import (
	"context"
	"fmt"
	"log"

	"github.com/phrazzld/resumake/api"
	"github.com/phrazzld/resumake/input"
)

func main() {
	fmt.Println("Resumake: A CLI tool for generating resumes")
	
	// Parse command-line flags
	flags, err := input.ParseFlags()
	if err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}
	
	// Get API key from environment variables
	apiKey, err := api.GetAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	
	// Initialize client and model for future use
	ctx := context.Background()
	client, model, err := api.InitializeClient(ctx, apiKey)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini client: %v", err)
	}
	defer client.Close()
	
	// Display initialization success message
	fmt.Printf("Successfully initialized Gemini client with model: %s\n", api.DefaultModelName)
	fmt.Println("Model configured with system instructions for resume generation")
	
	// Read source file if provided
	var sourceContent string
	var sourceFileRead bool
	
	if flags.SourcePath != "" {
		var err error
		sourceContent, sourceFileRead, err = input.ReadSourceFileFromFlags(flags)
		if err != nil {
			log.Fatalf("Error reading source file: %v", err)
		}
		
		if sourceFileRead {
			fmt.Printf("Successfully read source resume from: %s\n", flags.SourcePath)
		}
	}
	
	// Read stream-of-consciousness input from stdin
	stdinContent, err := input.ReadFromStdin()
	if err != nil {
		log.Fatalf("Error reading from stdin: %v", err)
	}
	
	if len(stdinContent) > 0 {
		fmt.Println("Successfully read input from stdin")
	} else {
		fmt.Println("Warning: No input received from stdin")
	}
	
	// Use model in the future for API calls
	_ = model // Prevent unused variable warning
	_ = sourceContent // Prevent unused variable warning
	_ = stdinContent // Prevent unused variable warning
}