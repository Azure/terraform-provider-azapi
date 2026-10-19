package customization

import (
	"encoding/base64"
	"reflect"
	"testing"
)

func TestWithKeyCreateParameters(t *testing.T) {
	modulus2048 := base64.RawURLEncoding.EncodeToString(append([]byte{0xc2}, make([]byte, 255)...))

	testCases := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name: "rsa key",
			input: map[string]interface{}{
				"key": map[string]interface{}{
					"kid":     "https://example.vault.azure.net/keys/rsa/0123",
					"kty":     "RSA",
					"key_ops": []interface{}{"wrapKey", "unwrapKey"},
					"n":       modulus2048,
					"e":       "AQAB",
				},
				"attributes": map[string]interface{}{"enabled": true},
				"tags":       map[string]interface{}{"env": "test"},
			},
			expected: map[string]interface{}{
				"key": map[string]interface{}{
					"kid":     "https://example.vault.azure.net/keys/rsa/0123",
					"kty":     "RSA",
					"key_ops": []interface{}{"wrapKey", "unwrapKey"},
					"n":       modulus2048,
					"e":       "AQAB",
				},
				"attributes":      map[string]interface{}{"enabled": true},
				"tags":            map[string]interface{}{"env": "test"},
				"kty":             "RSA",
				"key_ops":         []interface{}{"wrapKey", "unwrapKey"},
				"key_size":        float64(2048),
				"public_exponent": float64(65537),
			},
		},
		{
			name: "ec key",
			input: map[string]interface{}{
				"key": map[string]interface{}{
					"kid":     "https://example.vault.azure.net/keys/ec/0123",
					"kty":     "EC",
					"crv":     "P-256",
					"key_ops": []interface{}{"sign", "verify"},
					"x":       "AQ",
					"y":       "Ag",
				},
			},
			expected: map[string]interface{}{
				"key": map[string]interface{}{
					"kid":     "https://example.vault.azure.net/keys/ec/0123",
					"kty":     "EC",
					"crv":     "P-256",
					"key_ops": []interface{}{"sign", "verify"},
					"x":       "AQ",
					"y":       "Ag",
				},
				"kty":     "EC",
				"crv":     "P-256",
				"key_ops": []interface{}{"sign", "verify"},
			},
		},
		{
			name: "top level property is not overwritten",
			input: map[string]interface{}{
				"key": map[string]interface{}{"kty": "RSA"},
				"kty": "RSA-HSM",
			},
			expected: map[string]interface{}{
				"key": map[string]interface{}{"kty": "RSA"},
				"kty": "RSA-HSM",
			},
		},
		{
			name:     "no key",
			input:    map[string]interface{}{"attributes": map[string]interface{}{"enabled": true}},
			expected: map[string]interface{}{"attributes": map[string]interface{}{"enabled": true}},
		},
		{
			name:     "not an object",
			input:    "text",
			expected: "text",
		},
		{
			name:     "nil",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := withKeyCreateParameters(tc.input)
			if !reflect.DeepEqual(actual, tc.expected) {
				t.Fatalf("expected %#v, got %#v", tc.expected, actual)
			}
		})
	}
}

func TestWithKeyCreateParametersDoesNotModifyInput(t *testing.T) {
	input := map[string]interface{}{
		"key": map[string]interface{}{"kty": "EC", "crv": "P-384"},
	}

	_ = withKeyCreateParameters(input)

	if len(input) != 1 {
		t.Fatalf("expected the input to be unchanged, got %#v", input)
	}
}

func TestKeyVaultKeyReplaceTriggersRefs(t *testing.T) {
	customizedResource := GetCustomization("Microsoft.KeyVault/vaults/keys@7.4")
	if customizedResource == nil {
		t.Fatalf("expected a customization for Key Vault keys")
	}
	v, ok := (*customizedResource).(DataPlaneResourceWithReplaceTriggers)
	if !ok {
		t.Fatalf("expected the Key Vault key customization to implement DataPlaneResourceWithReplaceTriggers")
	}

	expected := []string{"kty", "key_size", "crv", "public_exponent"}
	if actual := v.ReplaceTriggersRefs(); !reflect.DeepEqual(actual, expected) {
		t.Fatalf("expected %v, got %v", expected, actual)
	}
}
