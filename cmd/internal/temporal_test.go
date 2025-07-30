package internal_test

import (
	"cooking-with-datastar/cmd/internal"
	"testing"
)

func TestDisplayMinutesSeconds(t *testing.T) {
	tt := []struct {
		name     string
		input    int
		expected string
	}{
		{"negative ninety nite seconds", -99, "00:00"},
		{"negative one second", -1, "00:00"},
		{"min", 0, "00:00"},
		{"one second", 1, "00:01"},
		{"ten seconds", 10, "00:10"},
		{"one minute twenty two seconds", 142, "02:22"},
		{"max", 3599, "59:59"},
		{"one hour", 3600, "59:59"},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			result := internal.DisplayMinutesSeconds(tc.input)

			if result != tc.expected {
				t.Logf("want '%s', got '%s'", tc.expected, result)
				t.Fail()
			}
		})
	}
}
