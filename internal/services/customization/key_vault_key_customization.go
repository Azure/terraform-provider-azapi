package customization

import (
	"context"
	"encoding/base64"
	"math/big"
	"net/http"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
)

// KeyVaultKeyCustomization supports keys that Key Vault generates.
//
// Create: POST   {vault}/keys/{name}/create
// Get:    GET    {vault}/keys/{name}
// Update: PATCH  {vault}/keys/{name}, which updates the current key version in place
// Delete: DELETE {vault}/keys/{name}
type KeyVaultKeyCustomization struct {
}

func (k KeyVaultKeyCustomization) CreateFunc() CreateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "create", id.ApiVersion, http.MethodPost, body, options)
		return err
	}
}

func (k KeyVaultKeyCustomization) UpdateFunc() UpdateFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, body interface{}, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.Action(ctx, id.AzureResourceId, "", id.ApiVersion, http.MethodPatch, body, options)
		return err
	}
}

func (k KeyVaultKeyCustomization) ReadFunc() ReadFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) (interface{}, error) {
		responseBody, err := client.DataPlaneClient.Get(ctx, id, options)
		if err != nil {
			return nil, err
		}
		return withKeyCreateParameters(responseBody), nil
	}
}

// withKeyCreateParameters adds the create key parameters to a key bundle returned by Get Key, at the
// same paths as in the create request, so that the configured body can be compared with the key.
// Get Key returns kty, key_ops and crv under "key", and the RSA modulus and public exponent as the
// JSON web key values "n" and "e".
func withKeyCreateParameters(responseBody interface{}) interface{} {
	bundle, ok := responseBody.(map[string]interface{})
	if !ok {
		return responseBody
	}
	key, ok := bundle["key"].(map[string]interface{})
	if !ok {
		return responseBody
	}

	parameters := make(map[string]interface{})
	for _, name := range []string{"kty", "key_ops", "crv"} {
		if v, ok := key[name]; ok && v != nil {
			parameters[name] = v
		}
	}
	if n := decodeBase64URL(key["n"]); len(n) != 0 {
		parameters["key_size"] = float64(len(n) * 8)
	}
	if e := decodeBase64URL(key["e"]); len(e) != 0 {
		parameters["public_exponent"] = float64(new(big.Int).SetBytes(e).Int64())
	}

	out := make(map[string]interface{}, len(bundle)+len(parameters))
	for name, v := range bundle {
		out[name] = v
	}
	for name, v := range parameters {
		if _, ok := out[name]; !ok {
			out[name] = v
		}
	}
	return out
}

func decodeBase64URL(input interface{}) []byte {
	v, ok := input.(string)
	if !ok {
		return nil
	}
	b, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(v, "="))
	if err != nil {
		return nil
	}
	return b
}

func (k KeyVaultKeyCustomization) DeleteFunc() DeleteFunc {
	return func(ctx context.Context, client clients.Client, id parse.DataPlaneResourceId, options clients.RequestOptions) error {
		_, err := client.DataPlaneClient.DeleteThenPoll(ctx, id, options)
		return err
	}
}

func (k KeyVaultKeyCustomization) GetResourceType() string {
	return "Microsoft.KeyVault/vaults/keys"
}

var _ DataPlaneResource = &KeyVaultKeyCustomization{}
