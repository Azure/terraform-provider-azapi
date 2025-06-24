package fix

import (
	"reflect"
	"testing"
)

func Test_fixUserAssignedIdentities(t *testing.T) {
	testcases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},

		{
			name:     "empty input",
			input:    map[string]interface{}{},
			expected: map[string]interface{}{},
		},

		{
			name: "input with no identity",
			input: map[string]interface{}{
				"name": "testResource",
				"type": "Microsoft.Test/testResources",
			},
			expected: map[string]interface{}{
				"name": "testResource",
				"type": "Microsoft.Test/testResources",
			},
		},

		{
			name: "input with identity but no userAssignedIdentities",
			input: map[string]interface{}{
				"name": "testResource",
				"type": "Microsoft.Test/testResources",
				"identity": map[string]interface{}{
					"type": "SystemAssigned",
				},
			},
			expected: map[string]interface{}{
				"name": "testResource",
				"type": "Microsoft.Test/testResources",
				"identity": map[string]interface{}{
					"type": "SystemAssigned",
				},
			},
		},

		{
			name: "input with userAssignedIdentities",
			input: map[string]interface{}{
				"name": "testResource",
				"type": "Microsoft.Test/testResources",
				"identity": map[string]interface{}{
					"type": "UserAssigned",
					"userAssignedIdentities": map[string]interface{}{
						"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/testGroup/providers/Microsoft.ManagedIdentity/userAssignedIdentities/testIdentity1": map[string]interface{}{
							"clientId":    "12345678-1234-1234-1234-123456789012",
							"principalId": "12345678-1234-1234-1234-123456789012",
						},
					},
				},
			},
			expected: map[string]interface{}{
				"name": "testResource",
				"type": "Microsoft.Test/testResources",
				"identity": map[string]interface{}{
					"type": "UserAssigned",
					"userAssignedIdentities": map[string]interface{}{
						"/subscriptions/12345678-1234-1234-1234-123456789012/resourceGroups/testGroup/providers/Microsoft.ManagedIdentity/userAssignedIdentities/testIdentity1": map[string]interface{}{},
					},
				},
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			result := fixUserAssignedIdentities(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("expected %v, got %v", tc.expected, result)
			}
		})
	}

}
