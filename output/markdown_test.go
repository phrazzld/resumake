package output

import (
	"testing"
)

func TestValidateMarkdown(t *testing.T) {
	tests := []struct {
		name    string
		content string
		wantErr bool
	}{
		{
			name:    "valid markdown with headers",
			content: "# Resume\n\n## Skills\n\n- Go\n- Python",
			wantErr: false,
		},
		{
			name:    "valid markdown with lists",
			content: "# John Doe\n\n- Software Engineer\n- 5+ years experience",
			wantErr: false,
		},
		{
			name:    "empty content",
			content: "",
			wantErr: true,
		},
		{
			name:    "non-markdown content",
			content: "Hello world without any markdown syntax",
			wantErr: true,
		},
		{
			name:    "malformed markdown",
			content: "# Missing newline\n## Another header without proper spacing",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMarkdown(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCleanMarkdown(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected string
	}{
		{
			name:     "already clean content",
			content:  "# Resume\n\n## Skills\n\n- Go\n- Python",
			expected: "# Resume\n\n## Skills\n\n- Go\n- Python",
		},
		{
			name:     "extra whitespace",
			content:  "# Resume\n\n\n\n## Skills\n\n- Go\n- Python  ",
			expected: "# Resume\n\n## Skills\n\n- Go\n- Python",
		},
		{
			name:     "inconsistent newlines",
			content:  "# Resume\r\n## Skills\r\n- Go\r\n- Python",
			expected: "# Resume\n\n## Skills\n\n- Go\n- Python",
		},
		{
			name:     "leading and trailing whitespace",
			content:  "  \n  # Resume\n\n## Skills\n\n- Go\n- Python\n  ",
			expected: "# Resume\n\n## Skills\n\n- Go\n- Python",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CleanMarkdown(tt.content); got != tt.expected {
				t.Errorf("CleanMarkdown() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestPrepareForOutput(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		wantErr  bool
		expected string
	}{
		{
			name:     "valid markdown",
			content:  "# Resume\n\n## Skills\n\n- Go\n- Python",
			wantErr:  false,
			expected: "# Resume\n\n## Skills\n\n- Go\n- Python",
		},
		{
			name:     "valid markdown that needs cleaning",
			content:  "  # Resume\n\n\n## Skills\n\n- Go\n- Python  ",
			wantErr:  false,
			expected: "# Resume\n\n## Skills\n\n- Go\n- Python",
		},
		{
			name:     "empty content",
			content:  "",
			wantErr:  true,
			expected: "",
		},
		{
			name:     "non-markdown content",
			content:  "Just plain text",
			wantErr:  true,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := PrepareForOutput(tt.content)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrepareForOutput() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.expected {
				t.Errorf("PrepareForOutput() = %v, want %v", got, tt.expected)
			}
		})
	}
}