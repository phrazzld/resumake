// Package main provides the entry point for the resumake application.
//
// Resumake is a command-line tool that generates professional resumes using the Gemini API.
// It takes a stream-of-consciousness text input from the user (optionally combined with
// an existing resume) and transforms it into a polished, well-structured resume in Markdown format.
//
// The main package orchestrates the entire process flow by initializing and running a
// Terminal User Interface (TUI) built with the BubbleTea framework.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/phrazzld/resumake/input"
	"github.com/phrazzld/resumake/tui"
)

func main() {
	fmt.Println("Resumake: A CLI tool for generating resumes")
	
	// Parse command-line flags
	flags, err := input.ParseFlags()
	if err != nil {
		// Check if the error is due to the user requesting help
		if err == flag.ErrHelp {
			// The flag package already printed usage information
			// We just need to exit cleanly
			os.Exit(0)
		}
		// For any other parsing error, log fatally
		log.Fatalf("Error parsing flags: %v", err)
	}
	
	// Initialize the Bubble Tea model with flags for pre-filling inputs
	model := tui.NewModel()
	
	// If a source path was provided via flags, pre-fill it in the model
	if flags.SourcePath != "" {
		model = model.WithSourcePath(flags.SourcePath)
	}
	
	// If an output path was provided via flags, set it in the model
	if flags.OutputPath != "" {
		model = model.WithOutputPath(flags.OutputPath)
	}
	
	// Set up signal handling for graceful shutdown
	p := setupProgramWithSignalHandling(model)
	
	// Run the program
	if _, err := p.Run(); err != nil {
		log.Fatalf("Error running TUI: %v", err)
	}
	
	// Program finished successfully
	fmt.Println("\nResumake finished.")
}

// setupProgramWithSignalHandling creates a new Bubble Tea program with the given model
// and sets up signal handling for graceful shutdown.
func setupProgramWithSignalHandling(model tea.Model) *tea.Program {
	// Create a new Bubble Tea program with our model
	p := tea.NewProgram(model, tea.WithAltScreen())
	
	// Create a channel to listen for signals
	signalCh := make(chan os.Signal, 1)
	
	// Set up signal notification
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)
	
	// Start a goroutine to handle signals
	go func() {
		sig := <-signalCh
		
		// Log the signal that was received
		log.Printf("Received signal: %v", sig)
		
		// Clean up by sending a QuitMsg to the program
		// This ensures the cleanupAPIClient function is called
		// before exiting
		p.Send(tea.QuitMsg{})
	}()
	
	return p
}