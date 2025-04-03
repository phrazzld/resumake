# BACKLOG

*This is a loose, unordered list of potential future improvements and features beyond the initial prototype.*

**Core AI & Processing:**

*   [AI] Allow selection between Gemini models (e.g., 2.5 Pro, 2.0 Flash) via flag (`--model`).
*   [AI] Implement logic to automatically choose model based on input size, complexity, or user preference (cost/speed).
*   [AI] Develop more sophisticated prompt engineering, potentially breaking down the task (e.g., first structure, then elaborate).
*   [AI] Introduce prompt templates or allow user-provided system instructions/prompts.
*   [AI] Add functionality to tailor the resume based on a provided job description input (`--job-desc file.txt`).
*   [AI] Implement explicit tone customization flags (`--tone understated`, `--tone punchy`, `--tone corporate`).
*   [AI] Explore techniques for better handling of ambiguity or conflicting info in the input dump.
*   [AI] Add an interactive refinement loop (e.g., TUI asking "Regenerate summary?", "Rephrase section X?", "Focus more on Y?").
*   [AI] Improve deduplication logic within the AI prompt or via post-processing.
*   [AI] Add capability to explicitly extract and structure specific sections (e.g., `resumake --extract-skills`).
*   [AI] Explore using `genai.Blob` with appropriate MIME types (text/plain, text/markdown, application/pdf) for the `--source` input instead of pure text concatenation. (**From API Docs**)
*   [AI] Investigate forcing JSON output using `model.ResponseMIMEType` and `model.ResponseSchema` for specific structured extraction tasks. (**From API Docs**)
*   [AI] Utilize Context Caching (`client.CreateCachedContent`) if processing the *same* source resume repeatedly. (**From API Docs**)
*   [AI] Expose model parameters (`Temperature`, `TopP`, `TopK`, `MaxOutputTokens`) via flags. (**From API Docs**)
*   [AI] Evaluate specific Gemini models (Flash-Lite, 1.5 variants) for cost/speed trade-offs. (**From API Docs**)

**Input & Output:**

*   [Input] Support reading from `stdin` via pipe (`cat file.txt | resumake`).
*   [Input] Support other input formats besides plain text/Markdown (e.g., JSON Resume Schema).
*   [Input] Implement a more robust stream-of-consciousness capture (e.g., opening `$EDITOR`).
*   [Input] Accept PDF as `--source` input, utilizing the SDK's PDF handling. (**From API Docs**)
*   [Input] Handle large source files (>20MB) using the File API (`client.UploadFileFromPath`). (**From API Docs**)
*   [Output] Allow specifying the output filename via a flag (`--output my_resume.md`).
*   [Output] Add export functionality to PDF format (perhaps via `pandoc` or similar).
*   [Output] Add export functionality to DOCX format (perhaps via `pandoc` or similar).
*   [Output] Support different Markdown flavors or structuring options (e.g., CommonMark, GFM).
*   [Output] Add optional metadata/frontmatter to the Markdown output.

**Features & Integrations:**

*   [Feature] Generate a basic cover letter based on the resume and an optional job description.
*   [Feature] Implement resume diffing (`resumake diff resume_v1.md resume_v2.md`).
*   [Feature] Add basic version control or snapshotting capabilities (e.g., timestamped outputs).
*   [Integration] Integrate with GitHub API to pull repository contribution data.
*   [Integration] Integrate with LinkedIn API to pull profile data (high complexity/approval needed).

**User Experience & Developer Experience:**

*   [UX] Implement a configuration file (`~/.config/resumake/config.yaml`?) for API key and default settings.
*   [UX] Add progress indicators/spinners, especially during the API call.
*   [UX] Provide more verbose logging/debugging options (`--verbose`).
*   [UX] Implement better, more specific error messages and suggestions (e.g., check API key validity).
*   [UX] Add basic input validation (e.g., check if source file exists before reading).
*   [UX] Explore adding a lightweight TUI (Terminal User Interface) for interaction (e.g., using `bubbletea`).
*   [DX] Set up automated testing (unit, integration with mock API).
*   [DX] Set up CI/CD pipeline for builds and releases (GitHub Actions).
*   [DX] Improve code structure (e.g., separate packages for `api`, `input`, `output`).
*   [DX] Add code comments and documentation (`godoc`).
*   [UX] Add cost estimation/warning based on token usage before sending API request.
*   [UX] Improve handling of API rate limits (e.g., exponential backoff).

**Distribution & Platform:**

*   [Platform] Create build targets for different OS/architectures (Linux, macOS, Windows) using `goreleaser`.
*   [Distribution] Package for common package managers (Homebrew, Scoop, Apt, etc.).
*   [Platform] Explore a lightweight web UI wrapper around the core Go logic.
