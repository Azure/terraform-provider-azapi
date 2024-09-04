package docstrings

import "testing"

func TestAddBackquotes(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "Hello, %s!",
			expected: "Hello, `!",
		},
		{
			input:    "This is a %s test %s.",
			expected: "This is a ` test `.",
		},
		{
			input:    "This is a consecutive test %s%s.",
			expected: "This is a consecutive test ``.",
		},
	}

	for _, tc := range testCases {
		result := addBackquotes(tc.input)
		if result != tc.expected {
			t.Errorf("Expected: %s, but got: %s", tc.expected, result)
		}
	}
}
