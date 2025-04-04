package tui

import (
	"context"
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/generative-ai-go/genai"
	"github.com/phrazzld/resumake/api"
	"github.com/phrazzld/resumake/input"
	"github.com/phrazzld/resumake/output"
	"github.com/phrazzld/resumake/prompt"
)

// ReadSourceFileCmd returns a command that reads a source file
// and returns a FileReadResultMsg with the result.
func ReadSourceFileCmd(filePath string) tea.Cmd {
	return func() tea.Msg {
		// Skip file reading if path is empty
		if filePath == "" {
			return FileReadResultMsg{
				Success: true,
				Content: "",
				Error:   nil,
			}
		}

		content, err := input.ReadSourceFile(filePath)
		if err != nil {
			return FileReadResultMsg{
				Success: false,
				Content: "",
				Error:   fmt.Errorf("failed to read source file: %w", err),
			}
		}

		return FileReadResultMsg{
			Success: true,
			Content: content,
			Error:   nil,
		}
	}
}

// InitializeAPICmd returns a command that initializes the API client
// and returns an APIInitResultMsg with the result.
func InitializeAPICmd() tea.Cmd {
	return func() tea.Msg {
		// Get API key
		apiKey, err := api.GetAPIKey()
		if err != nil {
			return APIInitResultMsg{
				Success: false,
				Error:   fmt.Errorf("API key error: %w", err),
			}
		}

		// Initialize client and model
		ctx := context.Background()
		_, _, err = api.InitializeClient(ctx, apiKey)
		if err != nil {
			return APIInitResultMsg{
				Success: false,
				Error:   fmt.Errorf("failed to initialize API client: %w", err),
			}
		}

		return APIInitResultMsg{
			Success: true,
			Error:   nil,
		}
	}
}

// GenerateResumeCmd returns a command that generates a resume using the API
// and returns an APIResultMsg with the result.
func GenerateResumeCmd(client *genai.Client, model *genai.GenerativeModel, sourceContent, stdinContent, outputFlagPath string, dryRun bool) tea.Cmd {
	return func() tea.Msg {
		// Skip actual API call if this is a dry run (for testing)
		if dryRun {
			return APIResultMsg{
				Success:    true,
				Content:    "Test content (dry run)",
				OutputPath: outputFlagPath,
				Error:      nil,
			}
		}

		// Verify client and model are provided
		if client == nil || model == nil {
			return APIResultMsg{
				Success: false,
				Error:   fmt.Errorf("API client or model is nil"),
			}
		}
		
		// We don't need to close the client here since it's managed by the caller
		// The client lifecycle is now handled by the Model struct

		// Create context for the API request
		ctx := context.Background()
		
		// Build the prompt from source content and stdin input
		promptContent := prompt.GeneratePromptContent(sourceContent, stdinContent)

		// Execute API request with the prompt content
		response, err := api.ExecuteRequest(ctx, model, promptContent)
		if err != nil {
			return APIResultMsg{
				Success: false,
				Error:   fmt.Errorf("error executing API request: %w", err),
			}
		}

		// Process the API response
		markdownContent, err := output.ProcessResponseContent(response)
		truncatedMsg := ""

		// Handle truncation error
		if err != nil {
			// Check if this is a truncation error and we might be able to recover
			if response != nil && len(response.Candidates) > 0 &&
				response.Candidates[0].FinishReason == genai.FinishReasonMaxTokens {
				
				truncatedMsg = "Warning: Response was truncated due to token limit"
				
				// Try to recover partial content
				partialContent, recoverErr := api.TryRecoverPartialContent(response)
				if recoverErr == nil && partialContent != "" {
					markdownContent = partialContent
				} else {
					return APIResultMsg{
						Success: false,
						Error:   fmt.Errorf("error processing API response: %w", err),
					}
				}
			} else {
				return APIResultMsg{
					Success: false,
					Error:   fmt.Errorf("error processing API response: %w", err),
				}
			}
		}

		// Write the generated markdown to a file
		outputPath, err := output.WriteOutput(markdownContent, outputFlagPath)
		if err != nil {
			return APIResultMsg{
				Success: false,
				Error:   fmt.Errorf("error writing output file: %w", err),
			}
		}

		return APIResultMsg{
			Success:      true,
			Content:      markdownContent,
			OutputPath:   outputPath,
			TruncatedMsg: truncatedMsg,
			Error:        nil,
		}
	}
}

// SubmitStdinInputCmd returns a command that submits stdin input
// and returns a StdinSubmitMsg with the input.
func SubmitStdinInputCmd(content string) tea.Cmd {
	return func() tea.Msg {
		return StdinSubmitMsg{
			Content: content,
		}
	}
}

// SendProgressUpdateCmd returns a command that sends a progress update
// and returns a ProgressUpdateMsg with the update.
func SendProgressUpdateCmd(step, message string) tea.Cmd {
	return func() tea.Msg {
		return ProgressUpdateMsg{
			Step:    step,
			Message: message,
		}
	}
}