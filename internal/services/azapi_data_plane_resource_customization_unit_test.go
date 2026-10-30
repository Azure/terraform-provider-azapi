package services_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/stretchr/testify/require"
)

const keyVaultKeyType = "Microsoft.KeyVault/vaults/keys@7.4"

func TestDataPlaneResourceReadKeyVaultKey(t *testing.T) {
	stateBody := `{"kty":"RSA","key_size":2048,"key_ops":["encrypt","decrypt","wrapKey","unwrapKey"],"attributes":{"exp":1830297600}}`

	testCases := []struct {
		name         string
		keyOps       []string
		keySize      int
		expectedBody string
	}{
		{
			name:         "key matches config",
			keyOps:       []string{"encrypt", "decrypt", "wrapKey", "unwrapKey"},
			keySize:      2048,
			expectedBody: stateBody,
		},
		{
			name:         "key operations changed outside of terraform",
			keyOps:       []string{"encrypt", "decrypt"},
			keySize:      2048,
			expectedBody: `{"kty":"RSA","key_size":2048,"key_ops":["encrypt","decrypt"],"attributes":{"exp":1830297600}}`,
		},
		{
			name:         "key size differs from config",
			keyOps:       []string{"encrypt", "decrypt", "wrapKey", "unwrapKey"},
			keySize:      4096,
			expectedBody: `{"kty":"RSA","key_size":4096,"key_ops":["encrypt","decrypt","wrapKey","unwrapKey"],"attributes":{"exp":1830297600}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			transport := &sequenceTransport{responses: []sequenceResponse{
				{
					method:     http.MethodGet,
					statusCode: http.StatusOK,
					body:       keyBundle(t, tc.keyOps, tc.keySize),
					path:       "/keys/test",
				},
			}}
			r := newDataPlaneResourceWithTransport(t, transport)
			resourceSchema := dataPlaneResourceSchema(t, r)

			model := newDataPlaneResourceModel(keyVaultKeyType, "example.vault.azure.net", "test", dynamicFromJSON(t, stateBody))
			state := tfsdk.State{Schema: resourceSchema}
			require.False(t, state.Set(ctx, model).HasError())

			response := resource.ReadResponse{State: state}
			r.Read(ctx, resource.ReadRequest{State: state}, &response)
			require.False(t, response.Diagnostics.HasError(), "%+v", response.Diagnostics)
			assertRequestMethods(t, transport, []string{http.MethodGet})

			var actual services.DataPlaneResourceModel
			require.False(t, response.State.Get(ctx, &actual).HasError())
			actualBody, err := dynamic.ToJSON(actual.Body)
			require.NoError(t, err)
			require.JSONEq(t, tc.expectedBody, string(actualBody))
		})
	}
}

func keyBundle(t *testing.T, keyOps []string, keySize int) string {
	t.Helper()
	bundle := map[string]interface{}{
		"key": map[string]interface{}{
			"kid":     "https://example.vault.azure.net/keys/test/0123456789abcdef0123456789abcdef",
			"kty":     "RSA",
			"key_ops": keyOps,
			"n":       base64.RawURLEncoding.EncodeToString(append([]byte{0xc2}, make([]byte, keySize/8-1)...)),
			"e":       "AQAB",
		},
		"attributes": map[string]interface{}{
			"enabled":       true,
			"exp":           1830297600,
			"created":       1792193468,
			"updated":       1792193468,
			"recoveryLevel": "CustomizedRecoverable+Purgeable",
		},
	}
	out, err := json.Marshal(bundle)
	require.NoError(t, err)
	return string(out)
}
