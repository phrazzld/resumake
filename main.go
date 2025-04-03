package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func main() {
	fmt.Println("Resumake: A CLI tool for generating resumes")
	
	// Validate that we can import and use the dependencies
	// This is just a placeholder and will be replaced with actual implementation
	if os.Getenv("RESUMAKE_TEST_API") == "true" {
		ctx := context.Background()
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			log.Fatal("GEMINI_API_KEY environment variable is required")
		}
		
		client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
		if err != nil {
			log.Fatalf("Failed to create client: %v", err)
		}
		defer client.Close()
		
		fmt.Println("Successfully initialized Gemini client")
	}
}