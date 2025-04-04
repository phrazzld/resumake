// Package main provides the entry point for the resumake application.
//
// Resumake is a command-line tool that generates professional resumes using the Gemini API.
// It takes a stream-of-consciousness text input from the user (optionally combined with
// an existing resume) and transforms it into a polished, well-structured resume in Markdown format.
//
// The main package orchestrates the entire process flow:
// 1. Parsing command-line flags
// 2. Retrieving the API key from environment variables
// 3. Initializing the Gemini API client
// 4. Reading the source file (if provided)
// 5. Reading user input from stdin
// 6. Building a prompt from the inputs
// 7. Sending the prompt to the Gemini API
// 8. Processing the response and extracting Markdown content
// 9. Writing the generated resume to a file
// 10. Displaying a completion message with next steps
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/google/generative-ai-go/genai"
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
		// Check if this is a truncation error and we might be able to recover
		if response != nil && len(response.Candidates) > 0 && 
		   response.Candidates[0].FinishReason == genai.FinishReasonMaxTokens {
			fmt.Println("Warning: Response was truncated due to token limit")
			fmt.Println("Attempting to recover partial content...")
			
			// Try to recover partial content
			partialContent, recoverErr := api.TryRecoverPartialContent(response)
			if recoverErr == nil && partialContent != "" {
				fmt.Println("Successfully recovered partial content with truncation warning")
				markdownContent = partialContent
			} else {
				log.Fatalf("Error processing API response: %v", err)
			}
		} else {
			log.Fatalf("Error processing API response: %v", err)
		}
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
	
	// Display completion message
	fmt.Println("\n----- RESUME GENERATION COMPLETE -----")
	fmt.Printf("Output file:      %s\n", outputPath)
	fmt.Printf("Content length:   %d characters\n", len(markdownContent))
	fmt.Println("\nNext steps:")
	fmt.Printf("  * Review your resume at %s\n", outputPath)
	fmt.Println("  * Make any necessary edits to improve formatting")
	fmt.Println("  * Convert to other formats as needed (PDF, DOCX)")
	fmt.Println("---------------------------------------")
}