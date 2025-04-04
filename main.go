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
	
	// Display source path if provided
	if flags.SourcePath != "" {
		fmt.Printf("Using source resume file: %s\n", flags.SourcePath)
	}
	
	// Use model in the future for API calls
	_ = model // Prevent unused variable warning
}