package internal_test

import (
	"cooking-with-datastar/cmd/internal"
	"testing"
)

func TestToStartCase(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"kebab case", "kebab-case", "Kebab case"},
		{"snake case", "snake_case", "Snake case"},
		{"camel case", "camelCase", "Camel case"},
		{"pascal case", "PascalCase", "Pascal case"},
		{"spaces", "one two", "One two"},
		{"mixed case", "Some-words_go here", "Some words go here"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := internal.ToStartCase(tc.input)

			if result != tc.expected {
				t.Logf("want '%s', got '%s'", tc.expected, result)
				t.Fail()
			}
		})
	}
}
