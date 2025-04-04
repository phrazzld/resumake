# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands
- Build: `go build -o resumake`
- Run: `GEMINI_API_KEY=your_key ./resumake [-source existing_resume.md]`
- Help: `./resumake --help` or `./resumake -h`
- Test: `go test ./...` 
- Single test: `go test -run TestName ./path/to/package`
- Lint: `golangci-lint run`
- Architect: `architect --task "description" *.go */*.go` (generates a PLAN.md file for implementing features)

## Code Style Guidelines
- Formatting: Use `gofmt` for code formatting
- Imports: Group standard library imports first, then third-party, then project packages
- Error handling: Check and handle all errors; prefer explicit error messages
- Naming: Follow Go conventions (CamelCase for exported, camelCase for unexported)
- Types: Use strong typing; document interfaces thoroughly
- Documentation: Add comments for exported functions, types, and variables
- Structure: Organize code into packages based on functionality (api, input, output)
- Testing: Write unit tests for core functionality

## Project-Specific Guidelines
- API Key: Always retrieve from GEMINI_API_KEY environment variable
- Input handling: Support both file input (--source) and stdin
- Output: Generate Markdown format (resume_out.md by default)
- Keep code simple and focused on the core workflow for the prototype

## Git Workflow
- Always use conventional commits format: type(scope): description
  - Types: feat, fix, docs, style, refactor, test, chore
  - Example: feat(input): add support for stdin reading
- Keep commits focused on single logical changes
- Write descriptive commit messages in imperative mood