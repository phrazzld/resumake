package tui

import (
	"strings"
	"testing"
)

func TestWrapText(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		width    int
		expected string
	}{
		{
			name:     "Empty text",
			text:     "",
			width:    80,
			expected: "",
		},
		{
			name:     "Single word",
			text:     "Hello",
			width:    10,
			expected: "Hello",
		},
		{
			name:     "No wrapping needed",
			text:     "Short text",
			width:    20,
			expected: "Short text",
		},
		{
			name:     "Simple wrap",
			text:     "This is a longer text that should wrap",
			width:    10,
			expected: "This is a\nlonger\ntext that\nshould\nwrap",
		},
		{
			name:     "Zero width defaults to 80",
			text:     "Text with zero width",
			width:    0,
			expected: "Text with zero width",
		},
		{
			name:     "Negative width defaults to 80",
			text:     "Text with negative width",
			width:    -5,
			expected: "Text with negative width",
		},
		{
			name:     "Long words",
			text:     "Supercalifragilisticexpialidocious is a very long word",
			width:    10,
			expected: "Supercalif\nragilistic\nexpialidoc\nious\nis a very\nlong word",
		},
		{
			name:     "Multiple spaces",
			text:     "Text   with   multiple   spaces",
			width:    10,
			expected: "Text with\nmultiple\nspaces",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := wrapText(tt.text, tt.width)
			
			// Check if result matches expected
			if result != tt.expected {
				t.Errorf("wrapText(%q, %d) = %q, want %q", tt.text, tt.width, result, tt.expected)
			}
			
			// Additional check: ensure no line exceeds the width
			if tt.width > 0 {
				lines := strings.Split(result, "\n")
				for i, line := range lines {
					if len(line) > tt.width {
						t.Errorf("Line %d exceeds width %d: %q (length: %d)", i+1, tt.width, line, len(line))
					}
				}
			}
		})
	}
}