# TODO

## API Client Initialization Refactor
- [x] Add API client fields to Model struct
  - Description: Add `apiClient` and `apiModel` fields to the `tui.Model` struct to store instances
  - Dependencies: None
  - Priority: High

- [x] Modify Model.Update for centralized initialization
  - Description: Update the Model.Update method to initialize the API client once when transitioning to a state that requires it
  - Dependencies: Added API client fields to Model
  - Priority: High

- [x] Restructure GenerateResumeCmd to use passed clients
  - Description: Modify GenerateResumeCmd to accept client and model instances as parameters instead of initializing them
  - Dependencies: Added API client fields and modified Update method
  - Priority: High

- [ ] Implement proper client cleanup on exit
  - Description: Add logic to ensure client.Close() is called when the application exits
  - Dependencies: Added API client fields
  - Priority: Medium

- [ ] Update tests affected by API client changes
  - Description: Modify tests to account for the new API client initialization approach
  - Dependencies: All API client refactoring complete
  - Priority: Medium

## Text Wrapping Utility
- [ ] Create utils.go file
  - Description: Create new file in tui package for utility functions
  - Dependencies: None
  - Priority: High

- [ ] Implement wrapText function
  - Description: Extract wrap function from views.go into a shared utility function
  - Dependencies: Created utils.go file
  - Priority: High

- [ ] Update view rendering to use shared utility
  - Description: Replace local wrap functions in renderWelcomeView and renderGeneratingView with calls to the new utility
  - Dependencies: Implemented wrapText function
  - Priority: High

- [ ] Create utils_test.go for testing
  - Description: Add unit tests for the wrapText function, covering various edge cases
  - Dependencies: Implemented wrapText function
  - Priority: Medium

## Error Message Improvements
- [ ] Enhance truncation recovery error messages
  - Description: Update error handling in GenerateResumeCmd to include both original and recovery errors
  - Dependencies: None
  - Priority: Medium

- [ ] Test improved error message handling
  - Description: Verify error message improvements work correctly with test cases
  - Dependencies: Enhanced error messages
  - Priority: Low

## Terminal Compatibility
- [ ] Define list of target terminals for testing
  - Description: Create a list of terminals to test across Linux, macOS, and Windows
  - Dependencies: None
  - Priority: Medium

- [ ] Create testing procedure document
  - Description: Document steps to test TUI rendering across different terminals
  - Dependencies: Defined target terminals
  - Priority: Medium

- [ ] Execute terminal compatibility testing
  - Description: Test the TUI on the identified terminals and document issues
  - Dependencies: Created testing procedure document
  - Priority: Medium

- [ ] Implement compatibility fixes if needed
  - Description: Address any rendering issues found during terminal testing
  - Dependencies: Completed terminal testing
  - Priority: Low

## Spinner Animation Fix
- [ ] Review current spinner implementation
  - Description: Analyze the spinner animation issues mentioned in the CODE_REVIEW.md
  - Dependencies: None
  - Priority: High

- [ ] Implement spinner animation improvements
  - Description: Ensure spinner continues to animate correctly when state changes
  - Dependencies: Reviewed spinner implementation
  - Priority: High

- [ ] Test spinner animation
  - Description: Verify spinner animation works consistently across all states
  - Dependencies: Implemented spinner improvements
  - Priority: Medium