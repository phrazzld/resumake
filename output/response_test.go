package output

import (
	"testing"

	"github.com/google/generative-ai-go/genai"
)

func TestProcessResponseContent(t *testing.T) {
	tests := []struct {
		name     string
		response *genai.GenerateContentResponse
		wantErr  bool
	}{
		{
			name: "valid markdown response",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{
					{
						Content: &genai.Content{
							Parts: []genai.Part{
								genai.Text("# Resume\n\n## Skills\n\n- Go\n- Python"),
							},
						},
						FinishReason: genai.FinishReasonStop,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty response",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{},
			},
			wantErr: true,
		},
		{
			name:     "nil response",
			response: nil,
			wantErr:  true,
		},
		{
			name: "response with non-markdown content",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{
					{
						Content: &genai.Content{
							Parts: []genai.Part{
								genai.Text("Just plain text without markdown"),
							},
						},
						FinishReason: genai.FinishReasonStop,
					},
				},
			},
			wantErr: true,
		},
		{
			name: "response with error finish reason",
			response: &genai.GenerateContentResponse{
				Candidates: []*genai.Candidate{
					{
						Content: &genai.Content{
							Parts: []genai.Part{
								genai.Text("# Resume"),
							},
						},
						FinishReason: genai.FinishReasonSafety,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ProcessResponseContent(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("ProcessResponseContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExtractAndValidateMarkdown(t *testing.T) {
	tests := []struct {
		name         string
		responseText string
		wantErr      bool
	}{
		{
			name:         "valid markdown content",
			responseText: "# Resume\n\n## Skills\n\n- Go\n- Python",
			wantErr:      false,
		},
		{
			name:         "empty content",
			responseText: "",
			wantErr:      true,
		},
		{
			name:         "non-markdown content",
			responseText: "Just some plain text without markdown formatting",
			wantErr:      true,
		},
		{
			name:         "malformed markdown",
			responseText: "#Missing space\n##Also missing space",
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ExtractAndValidateMarkdown(tt.responseText)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractAndValidateMarkdown() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}