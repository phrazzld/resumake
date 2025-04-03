# PLAN

**Objective:** Build a v0.1 CLI prototype of `resumake` in Go. This prototype will accept an optional existing resume file (`-source`) and a stream-of-consciousness text dump from the user via `stdin`. It will combine these inputs into a prompt, send them to the Gemini API using the `gemini-2.5-pro-exp-03-25` model via the user's API key, and write the resulting synthesized Markdown resume to `resume_out.md`. The focus is on the core end-to-end flow with minimal configuration and maximum simplicity for the initial implementation.

**Core Workflow (Engineer's View):**

1.  **Initialization:**
    *   Import necessary packages: `context`, `fmt`, `log`, `os`, `io`, `flag`, `github.com/google/generative-ai-go/genai`, `google.golang.org/api/option`.
    *   Retrieve the API key from the `GEMINI_API_KEY` environment variable using `os.Getenv`. If empty, log a fatal error (`log.Fatal`) indicating the variable is required.
    *   Initialize the Gemini client: `client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))`. Handle errors appropriately (log.Fatal). Ensure client cleanup: `defer client.Close()`.
    *   Get a handle to the specific generative model: `model := client.GenerativeModel("gemini-2.5-pro-exp-03-25")`. *Note: Confirm this model identifier is correct and available via the SDK.*
    *   Define the core system instruction text (see Prompt Construction below). Set it on the model instance: `model.SystemInstruction = &genai.Content{ Parts: []genai.Part{genai.Text(systemInstruction)} }`.

2.  **Argument Parsing:**
    *   Use the `flag` package.
    *   Define `sourcePath := flag.String("source", "", "Optional path to existing resume file (txt or md)")`.
    *   `flag.Parse()`.

3.  **Input Acquisition:**
    *   **Source Resume (Conditional):**
        *   Initialize `sourceContent := ""` (string).
        *   If `*sourcePath` is not empty:
            *   Read the file contents: `sourceBytes, err := os.ReadFile(*sourcePath)`. Handle file read errors (log.Fatal).
            *   Convert bytes to string: `sourceContent = string(sourceBytes)`.
    *   **Stream-of-Consciousness Input:**
        *   Print prompt to `stdout`: `fmt.Println("Enter your raw professional history below. Press Ctrl+D (Unix) or Ctrl+Z then Enter (Windows) when finished:")`.
        *   Read all from `stdin`: `socBytes, err := io.ReadAll(os.Stdin)`. Handle stdin read errors (log.Fatal).
        *   Convert bytes to string: `socContent := string(socBytes)`.

4.  **API Interaction:**
    *   **Prompt Construction:**
        *   Define the fixed system instruction (to be set on the model instance):
            ```text
            You are an expert resume writing assistant. Your goal is to synthesize the provided existing resume information (if any) and the raw stream-of-consciousness input into a single, coherent, professional resume formatted strictly in Markdown.

            Prioritize clarity, conciseness, and professional language. Structure the output logically with clear headings (e.g., Summary, Experience, Projects, Skills, Education). Infer structure and dates where possible, but do not fabricate information not present in the inputs. Focus on elevating the user's actual experience. Ensure the final output is only Markdown content.
            ```
        *   Construct the user input prompt string dynamically. Include clear delimiters:
            ```go
            userInputPrompt := fmt.Sprintf("### Existing Resume Content:\n%s\n\n### Raw Input Dump:\n%s", sourceContent, socContent)
            // Handle the case where sourceContent is empty gracefully (e.g., omit the 'Existing Resume Content' section).
            if sourceContent == "" {
                 userInputPrompt = fmt.Sprintf("### Raw Input Dump:\n%s", socContent)
            }
            ```
    *   **API Call:**
        *   Prepare the request content part: `reqPart := genai.Text(userInputPrompt)`.
        *   Make the generation call: `resp, err := model.GenerateContent(ctx, reqPart)`.
        *   Handle API call errors (`err != nil`): log.Fatal, including potentially parsing specific API error types if useful for debugging. Check for blocked responses or safety issues (`resp.Candidates[0].FinishReason != genai.FinishReasonStop` or check `SafetyRatings`).
    *   **Response Handling:**
        *   Check if `resp.Candidates` is non-empty and `resp.Candidates[0].Content` is not nil.
        *   Iterate through `resp.Candidates[0].Content.Parts`. Expect a single `genai.Text` part.
        *   Extract the text: `if textPart, ok := part.(genai.Text); ok { generatedMarkdown = string(textPart) }`. Handle cases where the response might not be text as an error.

5.  **Output Generation:**
    *   Define output filename: `outputFilename := "resume_out.md"`.
    *   Write the `generatedMarkdown` string to the file: `err := os.WriteFile(outputFilename, []byte(generatedMarkdown), 0644)`. Handle file write errors (log.Fatal).

6.  **Completion:**
    *   Print success message to `stdout`: `fmt.Printf("Successfully generated resume to %s\n", outputFilename)`.
    *   The program will exit implicitly with status 0 if no `log.Fatal` occurred.

**Technology Stack:**

*   Language: Go (latest stable version)
*   Key Go Packages: `context`, `flag`, `fmt`, `os`, `io`, `log`
*   External Go Packages: `github.com/google/generative-ai-go/genai`, `google.golang.org/api/option`
*   AI Model: `gemini-2.5-pro-exp-03-25` (via user API Key)

**Engineering Considerations (Prototype Focus):**

*   **API Key:** Strictly `GEMINI_API_KEY` environment variable.
*   **Error Handling:** `log.Fatal` for simplicity on errors. Informative messages.
*   **Input Handling:** Treat both source file and stdin dump as plain text for concatenation into the main user prompt part. Do not use `genai.Blob` or File API for v0.1.
*   **Prompting:** Use `model.SystemInstruction` for the fixed role/goal definition. Pass concatenated user inputs as a single `genai.Text` part to `GenerateContent`.
*   **Output Format:** Rely solely on prompt instructions to get Markdown. No schema enforcement (`ResponseMIMEType`, `ResponseSchema`).
*   **Dependencies:** Keep minimal.
*   **Concurrency:** None required.
*   **Testing:** Manual execution and verification of `resume_out.md`.
*   **Performance:** API latency is the bottleneck; client-side performance is not a concern.
*   **Security:** User responsible for API key security. User input is sent to Google APIs.

**Non-Goals (for this Prototype):**

*   Using other models (Flash, etc.).
*   Using advanced API features (File API for large files, Blobs, Context Caching, Function Calling, Structured Output schemas).
*   Interactive refinement loop.
*   Configuration files or complex flag parsing.
*   Job description tailoring, tone flags.
*   Integrations, other output formats (PDF/DOCX).
*   Sophisticated error recovery or input validation.
