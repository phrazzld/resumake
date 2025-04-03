# TODO

## Project Setup
- [x] Initialize Go module
  - Description: Create go.mod file with required module path
  - Dependencies: None
  - Priority: High

- [x] Add required dependencies
  - Description: Add google/generative-ai-go/genai and google.golang.org/api/option
  - Dependencies: Go module initialization
  - Priority: High

## Core Architecture
- [ ] Implement main package structure
  - Description: Create basic file structure with main.go and supporting packages
  - Dependencies: None
  - Priority: High

- [ ] Set up error handling utilities
  - Description: Create consistent error handling approach using log.Fatal for v0.1
  - Dependencies: None
  - Priority: Medium

## API Integration
- [ ] Add API key retrieval
  - Description: Get GEMINI_API_KEY from environment variables with appropriate error handling
  - Dependencies: None
  - Priority: High

- [ ] Implement Gemini client initialization
  - Description: Set up client with proper context and API key
  - Dependencies: API key retrieval
  - Priority: High

- [ ] Configure model with system instructions
  - Description: Initialize model with required system prompts for resume generation
  - Dependencies: Gemini client initialization
  - Priority: High

## Input Handling
- [ ] Implement CLI flag parsing
  - Description: Set up the -source flag for optional existing resume file
  - Dependencies: None
  - Priority: High

- [ ] Create file reading functionality
  - Description: Read source resume file if provided
  - Dependencies: CLI flag parsing
  - Priority: High

- [ ] Implement stdin reading
  - Description: Read stream-of-consciousness input from stdin with appropriate user prompts
  - Dependencies: None
  - Priority: High

## Prompt Construction
- [ ] Create dynamic prompt builder
  - Description: Combine existing resume and user input into properly formatted prompt
  - Dependencies: File reading and stdin reading
  - Priority: High

## API Call Handling
- [ ] Implement API request execution
  - Description: Send constructed prompt to Gemini API
  - Dependencies: Prompt construction, Gemini client
  - Priority: High

- [ ] Add response handling
  - Description: Process API response and extract generated markdown
  - Dependencies: API request execution
  - Priority: High

- [ ] Implement error and edge case handling
  - Description: Handle API errors, safety ratings, and unexpected response formats
  - Dependencies: API request and response handling
  - Priority: Medium

## Output Generation
- [ ] Create file writing functionality
  - Description: Write generated markdown to resume_out.md
  - Dependencies: Response handling
  - Priority: High

- [ ] Add completion messaging
  - Description: Display success message with output filename
  - Dependencies: File writing
  - Priority: Medium

## Testing
- [ ] Manual testing
  - Description: Test end-to-end workflow with sample inputs
  - Dependencies: All implementation tasks
  - Priority: High

## Documentation
- [ ] Add code documentation
  - Description: Add comments for public functions and types
  - Dependencies: Implementation
  - Priority: Medium

- [ ] Create usage documentation
  - Description: Add README with installation and usage instructions
  - Dependencies: Implementation and testing
  - Priority: Medium

## Assumptions
- User will provide GEMINI_API_KEY as environment variable
- Single output format (Markdown) is sufficient for v0.1
- The specified Gemini model "gemini-2.5-pro-exp-03-25" is available via the SDK
- Error handling prioritizes simplicity over recovery for v0.1
- No advanced API features (Blobs, File API, Context Caching) are needed