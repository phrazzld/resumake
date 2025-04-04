package prompt

import (
	"testing"
)

func TestBuildPrompt(t *testing.T) {
	tests := []struct {
		name          string
		sourceContent string
		stdinContent  string
		want          string
	}{
		{
			name:          "both source and stdin provided",
			sourceContent: "# Existing Resume\n\nSkills: Go, Python",
			stdinContent:  "I also know JavaScript and worked at Google",
			want:          "EXISTING RESUME:\n# Existing Resume\n\nSkills: Go, Python\n\nUSER INPUT:\nI also know JavaScript and worked at Google",
		},
		{
			name:          "only source content",
			sourceContent: "# Existing Resume\n\nSkills: Go, Python",
			stdinContent:  "",
			want:          "EXISTING RESUME:\n# Existing Resume\n\nSkills: Go, Python\n\nUSER INPUT:\n(No additional input provided)",
		},
		{
			name:          "only stdin content",
			sourceContent: "",
			stdinContent:  "I know JavaScript and worked at Google",
			want:          "EXISTING RESUME:\n(No existing resume provided)\n\nUSER INPUT:\nI know JavaScript and worked at Google",
		},
		{
			name:          "neither source nor stdin content",
			sourceContent: "",
			stdinContent:  "",
			want:          "EXISTING RESUME:\n(No existing resume provided)\n\nUSER INPUT:\n(No additional input provided)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BuildPrompt(tt.sourceContent, tt.stdinContent); got != tt.want {
				t.Errorf("BuildPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePromptContent(t *testing.T) {
	tests := []struct {
		name          string
		sourceContent string
		stdinContent  string
		wantParts     int
	}{
		{
			name:          "both source and stdin provided",
			sourceContent: "# Resume",
			stdinContent:  "Input",
			wantParts:     1, // Should generate a single text part
		},
		{
			name:          "empty inputs still produce valid prompt",
			sourceContent: "",
			stdinContent:  "",
			wantParts:     1, // Should generate a single text part even with empty inputs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			content := GeneratePromptContent(tt.sourceContent, tt.stdinContent)
			if len(content.Parts) != tt.wantParts {
				t.Errorf("GeneratePromptContent() returned %d parts, want %d", len(content.Parts), tt.wantParts)
			}
		})
	}
}