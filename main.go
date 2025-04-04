package main

import (
	"context"
	"fmt"
	"log"

	"github.com/phrazzld/resumake/api"
)

func main() {
	fmt.Println("Resumake: A CLI tool for generating resumes")
	
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
	
	// This is just a placeholder and will be replaced with actual implementation
	fmt.Printf("Successfully initialized Gemini client with model: %s\n", api.DefaultModelName)
	
	// Use model in the future for API calls
	_ = model // Prevent unused variable warning
}