package utils

import "testing"

func TestStringIdentifierAtPath(t *testing.T) {
	testcases := []struct {
		path     string
		input    map[string]interface{}
		expected string
	}{
		{
			path:     "name",
			input:    map[string]interface{}{"name": "example"},
			expected: "example",
		},
		{
			path: "subnet.id",
			input: map[string]interface{}{
				"subnet": map[string]interface{}{"id": "/virtualNetworks/example/subnets/default"},
			},
			expected: "/virtualNetworks/example/subnets/default",
		},
	}

	for _, testcase := range testcases {
		t.Run(testcase.path, func(t *testing.T) {
			if actual := stringIdentifierAtPath(testcase.input, testcase.path); actual != testcase.expected {
				t.Errorf("expected %q, got %q", testcase.expected, actual)
			}
		})
	}
}
