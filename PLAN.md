# PLAN.md: Address Code Review Issues for TUI Implementation

## 1. Overview

This plan outlines the technical steps required to address the key issues identified in the `CODE_REVIEW.md` document for the `resumake` TUI implementation. The goal is to improve the application's robustness, maintainability, user experience, and resource management by resolving API client redundancy, code duplication, unclear error messages, terminal compatibility concerns, spinner animation inconsistency, incomplete signal handling, and missing cleanup calls on exit.

## 2. Task Breakdown

| Task                                                    | Description                                                                                                                                                              | Effort | Affected Files/Modules                                                               |
| :------------------------------------------------------ | :----------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :----- | :----------------------------------------------------------------------------------- |
| **1. Fix API Client Initialization Redundancy**         | Ensure API client is initialized only once in `tui.Model` and reused. Verify `InitializeAPICmd` is removed or its purpose clarified.                                       | M      | `tui/model.go`, `tui/commands.go`, `tui/api_client_test.go`                          |
| **2. Fix Duplicated Text Wrapping Logic**               | Consolidate the `wrapText` function into `tui/utils.go` and update view rendering functions to use the shared utility. Ensure tests pass.                               | S      | `tui/views.go`, `tui/utils.go`, `tui/utils_test.go`                                  |
| **3. Fix Unclear Truncation Recovery Error Messages**   | Modify error handling in `tui/commands.go` -> `GenerateResumeCmd` to include both original processing error and recovery error when `TryRecoverPartialContent` fails.      | S      | `tui/commands.go`, `tui/commands_test.go`                                            |
| **4. Address Terminal Compatibility Concerns**          | Execute the manual testing plan across target terminals. Document findings and implement necessary fallback styles or adjustments.                                       | M-L    | N/A (Testing Process), potentially `tui/views.go`, `tui/styles.go`                   |
| **5. Fix Spinner Animation Inconsistency**              | Ensure `spinner.Tick` command is consistently added in `tui.Model.Update` when `state == stateGenerating`, regardless of other messages. Verify tests pass.                | S      | `tui/model.go`, `tui/spinner_test.go`                                                |
| **6. Enhance Signal Handling & Context Cancellation**   | Introduce a root `context.Context` with cancellation. Pass it down to relevant commands (like `GenerateResumeCmd`). Cancel the context in the signal handler.             | M      | `main.go`, `tui/model.go`, `tui/commands.go`                                         |
| **7. Ensure Exit Handlers Call Cleanup**                | Verify that `cleanupAPIClient` is called reliably on all exit paths (`tea.QuitMsg`, `Ctrl+C`, `Esc`, final state transitions). Leverage `TestExitHandlersCallCleanup`.       | S      | `tui/model.go`, `tui/api_client_test.go`                                             |

**Effort Estimation:** S = Small (1-2 hours), M = Medium (2-4 hours), L = Large (4+ hours)

## 3. Implementation Details

### 3.1. Fix API Client Initialization Redundancy

*   **Goal:** Ensure the `genai.Client` and `genai.GenerativeModel` are initialized strictly once and managed correctly.
*   **Approach:**
    1.  **Verify Current State:** The current code in `tui/model.go` initializes the client/model in the `initializeAPIClient` helper, called during the transition from `stateWelcome` if `apiKeyOk` is true. The instances are stored in `m.apiClient` and `m.apiModel`. This seems correct.
    2.  **Review `InitializeAPICmd`:** Analyze `tui/commands.go`. The `InitializeAPICmd` appears unused in the main TUI flow, as initialization happens directly in the `Model.Update` function.
    3.  **Action:** Remove `InitializeAPICmd` from `tui/commands.go` and its corresponding message `APIInitResultMsg` from `tui/messages.go` if it's confirmed to be unused. Remove related tests if necessary.
    4.  **Refinement (Optional):** Consider moving the blocking calls (`api.GetAPIKey`, `api.InitializeClient`) within `initializeAPIClient` into a command to prevent potential UI freezes, although marked as Low-Medium risk in the review. This would involve:
        *   Creating a new command (e.g., `PerformAPIInitializationCmd`).
        *   Returning this command from `Update` when transitioning from `stateWelcome`.
        *   Creating a new message (e.g., `APIInitializationCompleteMsg`) returned by the command, carrying the client/model instances or an error.
        *   Handling this new message in `Update` to store the client/model or transition to the error state.

### 3.2. Fix Duplicated Text Wrapping Logic

*   **Goal:** Use the shared `tui/utils.go -> wrapText` function consistently.
*   **Approach:**
    1.  **Verify Current State:** The `diff.md` shows `tui/utils.go` and `tui/utils_test.go` were created, and `wrapText` implemented. `tui/views.go` was updated to use `wrapText`.
    2.  **Action:** Confirm that all instances of manual text wrapping in `tui/views.go` (specifically `renderWelcomeView` and `renderGeneratingView`) have been replaced with calls to `tui.wrapText`. Ensure `tui/utils_test.go` provides adequate coverage for `wrapText`. No code changes should be needed if the diff was applied correctly.

### 3.3. Fix Unclear Truncation Recovery Error Messages

*   **Goal:** Provide more context when API response truncation recovery fails.
*   **Approach:**
    1.  **Locate Code:** Find the error handling block in `tui/commands.go` -> `GenerateResumeCmd` where `api.TryRecoverPartialContent` is called.
    2.  **Verify Current State:** The `diff.md` shows the error message was updated:
        ```go
        // Inside GenerateResumeCmd, within the error handling for output.ProcessResponseContent
        // ...
        } else {
            // *** IMPROVED ERROR MESSAGE ***
            return APIResultMsg{
                Success: false,
                Error:   fmt.Errorf("error processing API response: %w (recovery failed: %w)", err, recoverErr),
            }
        }
        // ...
        ```
    3.  **Action:** Confirm this change is present and correctly wraps both the original processing error (`err`) and the recovery error (`recoverErr`). Ensure `tui/commands_test.go -> TestTruncationRecoveryErrorMessageImplementation` passes. No code changes should be needed if the diff was applied correctly.

### 3.4. Address Terminal Compatibility Concerns

*   **Goal:** Ensure the TUI renders acceptably across common terminal environments.
*   **Approach:**
    1.  **Execute Testing Plan:** Follow the procedure outlined in the previous `PLAN.md` (section 3.4):
        *   **Target Terminals:** GNOME Terminal, Konsole, xterm (Linux); Terminal.app, iTerm2 (macOS); Windows Terminal (Windows).
        *   **Procedure:** Build and run `resumake` on each terminal. Test all states, check rendering (colors, borders, alignment, wrapping, spinner), cursor behavior, and resizing.
    2.  **Document Findings:** Record any rendering issues, specifying terminal, OS, and the problem. Use screenshots.
    3.  **Implement Fixes (If Needed):**
        *   Prioritize fixes for major layout breaks or unreadable text.
        *   Leverage `lipgloss` features for compatibility where possible.
        *   Consider simplifying styles (e.g., `lipgloss.NormalBorder()` instead of `lipgloss.RoundedBorder()`) if specific features cause widespread issues.
        *   If necessary, use `termenv` to detect terminal capabilities and apply conditional styling, though this adds complexity.

### 3.5. Fix Spinner Animation Inconsistency

*   **Goal:** Ensure the loading spinner animates smoothly during the `stateGenerating`.
*   **Approach:**
    1.  **Locate Code:** Examine the `tui.Model.Update` function in `tui/model.go`.
    2.  **Verify Current State:** The `diff.md` shows the following logic was added:
        ```go
        // Handle spinner updates based on state
        if m.state == stateGenerating {
            var spinnerCmd tea.Cmd
            // Update the spinner regardless of msg type to ensure animation consistency
            m.spinner, spinnerCmd = m.spinner.Update(msg)

            // Always ensure the spinner keeps ticking by adding the tick command
            cmds = append(cmds, spinnerCmd, m.spinner.Tick) // Ensures Tick is always added
        }
        // ... Batch commands ...
        ```
        This ensures `m.spinner.Tick` is added to the command batch whenever the state is `stateGenerating`.
    3.  **Action:** Confirm this logic is present and correctly ensures the `spinner.Tick` command is always returned when in the generating state. Verify `tui/spinner_test.go` passes. No code changes should be needed if the diff was applied correctly.

### 3.6. Enhance Signal Handling & Context Cancellation

*   **Goal:** Allow long-running operations (like API calls) to be cancelled gracefully upon receiving a signal (Ctrl+C).
*   **Approach:**
    1.  **Create Root Context:** In `main.go`, create a cancellable context before starting the Bubble Tea program.
        ```go
        // main.go
        import (
            "context"
            // ... other imports
        )

        func main() {
            // ... flag parsing ...

            // Create a cancellable context
            ctx, cancel := context.WithCancel(context.Background())
            defer cancel() // Ensure cancel is called eventually

            model := tui.NewModel().WithContext(ctx) // Add method to pass context

            // ... setup signal handling ...
            p := setupProgramWithSignalHandling(model, cancel) // Pass cancel func

            // ... run program ...
        }

        func setupProgramWithSignalHandling(model tea.Model, cancel context.CancelFunc) *tea.Program {
            // ... existing setup ...
            go func() {
                sig := <-signalCh
                log.Printf("Received signal: %v", sig)
                cancel() // <<< Cancel the context here
                p.Send(tea.QuitMsg{})
            }()
            return p
        }
        ```
    2.  **Store Context in Model:** Add a `context.Context` field to `tui.Model`. Add a `WithContext` method to set it during initialization.
        ```go
        // tui/model.go
        type Model struct {
            // ... existing fields ...
            ctx context.Context
        }

        func (m Model) WithContext(ctx context.Context) Model {
            m.ctx = ctx
            return m
        }
        ```
    3.  **Pass Context to Commands:** Modify commands that perform potentially long-running operations (like `GenerateResumeCmd`) to accept and use the context from the model.
        ```go
        // tui/commands.go
        func GenerateResumeCmd(ctx context.Context, client *genai.Client, /*...*/) tea.Cmd {
            return func() tea.Msg {
                // ... existing setup ...
                // Use the passed context for the API request
                response, err := api.ExecuteRequest(ctx, model, promptContent)
                // ... handle response ...
            }
        }

        // tui/model.go -> Update (when calling GenerateResumeCmd)
        cmds = append(cmds, GenerateResumeCmd(m.ctx, m.apiClient, /*...*/))
        ```
    4.  **Update API Layer:** Ensure `api.ExecuteRequest` and `api.InitializeClient` accept and use the passed context. (They already do).

### 3.7. Ensure Exit Handlers Call Cleanup

*   **Goal:** Guarantee `cleanupAPIClient` is called on all program exit paths.
*   **Approach:**
    1.  **Verify Current State:** The `diff.md` shows `cleanupAPIClient` is called in `tui.Model.Update` for:
        *   `tea.QuitMsg`
        *   `tea.KeyCtrlC`
        *   `tea.KeyEsc`
        *   `tea.KeyEnter` when in `stateResultSuccess` or `stateResultError`.
    2.  **Signal Handler:** The signal handler sends `tea.QuitMsg`, which triggers the cleanup logic in `Update`.
    3.  **Action:** Confirm these calls are present. Run the existing test `tui/api_client_test.go -> TestExitHandlersCallCleanup` to verify. This test mocks `cleanupAPIClient` and checks if it's called for the relevant messages and states. No code changes should be needed if the diff was applied correctly and the test passes.

## 4. Potential Challenges & Considerations

*   **API Client Lifecycle:** Ensuring `client.Close()` is called *exactly once* on all exit paths requires careful state management. The current approach of calling it within `Update` handlers for exit messages/keys seems robust.
*   **Terminal Compatibility:** Fixing rendering issues across diverse terminals can be time-consuming and may require compromises on visual fidelity. A pragmatic approach might be needed, prioritizing functionality over perfect rendering everywhere.
*   **Context Propagation:** Adding context cancellation requires passing the context through relevant layers (`main` -> `tui.Model` -> `tui.Commands` -> `api`). This adds minor complexity but improves robustness.
*   **Testing Asynchronous Operations:** Testing context cancellation thoroughly might require more complex integration tests or mocks that simulate blocking operations.

## 5. Testing Strategy

*   **Unit Tests:**
    *   Verify `tui/utils_test.go` covers `wrapText` edge cases.
    *   Verify `tui/commands_test.go -> TestTruncationRecoveryErrorMessageImplementation` passes.
    *   Verify `tui/spinner_test.go` confirms consistent ticking in the generating state.
    *   Verify `tui/api_client_test.go -> TestAPIClientInitialization` passes (confirms single initialization).
    *   Verify `tui/api_client_test.go -> TestExitHandlersCallCleanup` passes (confirms cleanup is called).
    *   Add unit tests for context propagation if significant logic is added (e.g., checking if context is passed to commands).
*   **Integration Tests:**
    *   The existing `main_test.go` tests help flag parsing but are less effective for the TUI.
    *   Consider adding integration tests that simulate sending signals (e.g., `syscall.Kill(cmd.Process.Pid, syscall.SIGINT)`) and check for graceful shutdown, although this can be platform-dependent and complex.
*   **Manual Testing:**
    *   **Crucial** for **Terminal Compatibility:** Execute the testing plan from section 3.4 across all target terminals.
    *   Verify spinner animation smoothness during the generation phase.
    *   Test all exit paths (Ctrl+C, Esc, Enter in final states) to ensure clean shutdown without errors.
    *   Test context cancellation by sending Ctrl+C during the "Generating" phase and observing if the process exits promptly (may require adding logging or simulating a long API call).

## 6. Open Questions

1.  **Terminal Compatibility Strategy:** If significant compatibility issues are found, what is the preferred approach: aiming for broad compatibility with simpler styles, or defining a minimum set of required terminal features?
2.  **Context Cancellation Necessity:** While good practice, is immediate cancellation of the *current* API request strictly necessary for this prototype, or is ensuring cleanup via `tea.QuitMsg` sufficient for now? (Recommendation: Implement context cancellation for robustness).
3.  **Blocking Calls in `Update`:** Is the potential minor UI freeze during initial API client setup acceptable (Low-Medium risk), or should it be moved to a command as suggested in 3.1? (Recommendation: Defer unless observed to be a problem).
