package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
	"github.com/phrazzld/resumake/api"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("Resumake: A CLI tool for generating resumes")
	
	// Get API key from environment variables
	apiKey, err := api.GetAPIKey()
	if err != nil {
		log.Fatal(err)
	}
	
	// Initialize client for future use
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatalf("Failed to create Gemini client: %v", err)
	}
	defer client.Close()
	
	// This is just a placeholder and will be replaced with actual implementation
	fmt.Println("Successfully initialized Gemini client")
}