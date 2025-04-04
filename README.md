# resumake

A CLI tool for generating professional resumes using AI.

## Overview

resumake is a command-line tool that uses the Gemini API to transform your stream-of-consciousness thoughts or existing resume into a polished, professional resume in Markdown format. Simply provide your professional experience, skills, and background, and resumake will generate a well-structured resume.

## Features

- Generate a professional resume from your input
- Enhance or refine an existing resume
- Combine existing resume with new input
- Output in clean Markdown format
- Simple command-line interface

## Installation

### Prerequisites

- Go 1.21 or higher
- A valid Gemini API key

### Installing from source

```bash
# Clone the repository
git clone https://github.com/phrazzld/resumake.git
cd resumake

# Build the binary
go build -o resumake

# Optional: Move the binary to your PATH
mv resumake /usr/local/bin/
```

## Configuration

resumake requires a Gemini API key to function. You can obtain one from the [Google AI Studio](https://makersuite.google.com/app/apikey).

Set your API key as an environment variable:

```bash
export GEMINI_API_KEY=your_api_key_here
```

For persistent configuration, add this to your shell profile (.bashrc, .zshrc, etc.):

```bash
echo 'export GEMINI_API_KEY=your_api_key_here' >> ~/.bashrc
source ~/.bashrc
```

## Usage

### Basic Usage

Run resumake and enter your professional experience:

```bash
resumake
```

This will prompt you to enter your professional information. Type or paste your content, then press Ctrl+D (Unix) or Ctrl+Z followed by Enter (Windows) to finish. resumake will generate a resume and save it as `resume_out.md`.

### Using an Existing Resume

Provide an existing resume file to refine or enhance it:

```bash
resumake -source existing_resume.md
```

resumake will use your existing resume as a foundation and still prompt you for additional input.

### Specifying Output File

To change the output filename:

```bash
resumake -output my_new_resume.md
```

## Example

Input:
```
I've been a software engineer for 5 years. Started at Amazon Web Services
as a junior developer in 2018 working on EC2. Promoted to SDE II in 2020.
Moved to Google in 2021 as a Senior Software Engineer working on
Cloud Storage. I know Python, Go, JavaScript, and some Rust.
Bachelor's in Computer Science from MIT in 2017.
```

Output (in `resume_out.md`):
```markdown
# Jane Doe

## Summary
Experienced Software Engineer with 5 years of expertise in cloud infrastructure at top-tier tech companies. Skilled in Python, Go, JavaScript, and Rust, with a proven track record of advancement and technical leadership.

## Experience
**Senior Software Engineer** | Google | 2021 - Present
- Working on Google Cloud Storage solutions
- [More details would be added based on input]

**Software Development Engineer II** | Amazon Web Services | 2020 - 2021
- Promoted to SDE II in recognition of technical contributions
- [More details would be added based on input]

**Software Development Engineer** | Amazon Web Services | 2018 - 2020
- Worked on EC2 services
- [More details would be added based on input]

## Skills
- **Programming Languages**: Python, Go, JavaScript, Rust
- **Cloud Platforms**: AWS, Google Cloud

## Education
**Bachelor of Science in Computer Science** | Massachusetts Institute of Technology | 2017
```

## Limitations

- resumake currently only outputs Markdown format
- The quality of the resume depends on the detail and clarity of your input
- API rate limits may apply depending on your Gemini API key tier

## Troubleshooting

### API Key Issues

If you see an error about the API key:
- Ensure the `GEMINI_API_KEY` environment variable is set correctly
- Verify your API key is valid and has not expired
- Check that you have quota available for the Gemini API

### Generation Issues

If the resume generation fails:
- Try providing more detailed input
- Ensure your input doesn't contain any content that might trigger safety filters
- If the response is truncated, try breaking your input into smaller, more focused parts

## License

[MIT License](LICENSE)
