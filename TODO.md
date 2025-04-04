# TODO

## API Client Management
- [x] Verify API client initialization approach
  - Description: Confirm the current code in `tui/model.go` properly initializes client/model in `initializeAPIClient`
  - Dependencies: None
  - Priority: High
- [x] Review and remove redundant `InitializeAPICmd`
  - Description: Analyze `tui/commands.go` and remove `InitializeAPICmd` if unused
  - Dependencies: Verification of initialization approach
  - Priority: High
- [x] Verify API client cleanup on all exit paths
  - Description: Ensure `cleanupAPIClient` is called reliably on all program exit paths
  - Dependencies: None
  - Priority: High

## Text Wrapping Functionality
- [x] Confirm shared text wrapping implementation
  - Description: Verify all instances of manual text wrapping are replaced with calls to `tui.wrapText`
  - Dependencies: None
  - Priority: Medium

## Context Management
- [x] Add context support to Model
  - Description: Add `context.Context` field to `tui.Model` and `WithContext` method
  - Dependencies: None
  - Priority: Medium
- [x] Update main.go with cancellable context
  - Description: Create root context and pass cancel function to signal handler
  - Dependencies: Model context support
  - Priority: Medium
- [x] Modify commands to use context
  - Description: Update `GenerateResumeCmd` to accept and use context for API requests
  - Dependencies: Model context support
  - Priority: Medium
