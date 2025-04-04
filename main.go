package main

import (
	"context"
	"fmt"
	"log"

	"github.com/phrazzld/resumake/api"
	"github.com/phrazzld/resumake/input"
	"github.com/phrazzld/resumake/output"
	"github.com/phrazzld/resumake/prompt"
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
	
	// Build the prompt from source content and stdin input
	promptContent := prompt.GeneratePromptContent(sourceContent, stdinContent)
	
	// Display confirmation that the prompt has been built
	fmt.Println("Successfully built prompt from inputs")
	
	// Execute API request with the prompt content
	fmt.Println("Executing API request to generate resume...")
	response, err := api.ExecuteRequest(ctx, model, promptContent)
	if err != nil {
		log.Fatalf("Error executing API request: %v", err)
	}
	
	// Process the API response
	fmt.Println("Processing API response...")
	markdownContent, err := output.ProcessResponseContent(response)
	if err != nil {
		log.Fatalf("Error processing API response: %v", err)
	}
	
	// Display API response information
	fmt.Println("Successfully received and processed API response")
	fmt.Println("Validated and cleaned Markdown content")
	
	// Write the generated markdown to a file
	fmt.Println("Writing generated resume to file...")
	outputPath, err := output.WriteOutput(markdownContent, flags.OutputPath)
	if err != nil {
		log.Fatalf("Error writing output file: %v", err)
	}
	
	// Display success message
	fmt.Printf("Successfully generated resume with %d characters\n", len(markdownContent))
	fmt.Printf("Resume saved to: %s\n", outputPath)
}