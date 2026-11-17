package services_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

const keyVaultKeyRotationPolicyType = "Microsoft.KeyVault/vaults/keys/rotationPolicy@7.6"

func TestDataPlaneResourceCreateKeyVaultKeyRotationPolicy(t *testing.T) {
	ctx := context.Background()
	policyBody := `{"lifetimeActions":[{"trigger":{"timeAfterCreate":"P90D"},"action":{"type":"Rotate"}},{"trigger":{"timeBeforeExpiry":"P30D"},"action":{"type":"Notify"}}],"attributes":{"expiryTime":"P1Y"}}`
	// Every key has a rotation policy, so create must not fail because the policy already exists.
	transport := &sequenceTransport{responses: []sequenceResponse{
		{
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       policyBody,
			path:       "/keys/test/rotationpolicy",
		},
		{
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			body:       policyBody,
			path:       "/keys/test/rotationpolicy",
		},
	}}
	r := newDataPlaneResourceWithTransport(t, transport)
	resourceSchema := dataPlaneResourceSchema(t, r)

	model := newDataPlaneResourceModel(keyVaultKeyRotationPolicyType, "example.vault.azure.net/keys/test", "", dynamicFromJSON(t, policyBody))
	model.ID = types.StringUnknown()
	plan := tfsdk.Plan{Schema: resourceSchema}
	require.False(t, plan.Set(ctx, model).HasError())
	state := tfsdk.State{
		Schema: resourceSchema,
		Raw:    tftypes.NewValue(resourceSchema.Type().TerraformType(ctx), nil),
	}

	var diagnostics diag.Diagnostics
	r.CreateUpdate(ctx, tfsdk.Config{Schema: resourceSchema, Raw: plan.Raw}, plan, &state, &diagnostics, newTestPrivateData())
	require.False(t, diagnostics.HasError(), "%+v", diagnostics)
	assertRequestMethods(t, transport, []string{http.MethodPut, http.MethodGet})

	var actual services.DataPlaneResourceModel
	require.False(t, state.Get(ctx, &actual).HasError())
	require.Equal(t, "example.vault.azure.net/keys/test/rotationpolicy", actual.ID.ValueString())
}

func TestDataPlaneResourceDeleteKeyVaultKeyRotationPolicy(t *testing.T) {
	testCases := []struct {
		name       string
		statusCode int
		response   string
	}{
		{
			name:       "sends a policy with no lifetime actions",
			statusCode: http.StatusOK,
			response:   `{"lifetimeActions":[{"trigger":{"timeBeforeExpiry":"P30D"},"action":{"type":"Notify"}}]}`,
		},
		{
			name:       "key already deleted",
			statusCode: http.StatusNotFound,
			response:   `{"error":{"code":"KeyNotFound","message":"A key with (name/id) test was not found in this key vault."}}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			// A DELETE request to the rotation policy URL deletes the key, so delete must send a PUT.
			transport := &sequenceTransport{responses: []sequenceResponse{
				{
					method:     http.MethodPut,
					statusCode: tc.statusCode,
					body:       tc.response,
					path:       "/keys/test/rotationpolicy",
				},
			}}
			recorder := &requestBodyRecorder{transport: transport}
			r := newDataPlaneResourceWithTransport(t, recorder)
			resourceSchema := dataPlaneResourceSchema(t, r)

			state := tfsdk.State{Schema: resourceSchema}
			require.False(t, state.Set(ctx, newDataPlaneResourceModel(keyVaultKeyRotationPolicyType, "example.vault.azure.net/keys/test", "", dynamicFromJSON(t, `{}`))).HasError())

			response := resource.DeleteResponse{State: state}
			r.Delete(ctx, resource.DeleteRequest{State: state}, &response)
			require.False(t, response.Diagnostics.HasError(), "%+v", response.Diagnostics)
			assertRequestMethods(t, transport, []string{http.MethodPut})
			require.Len(t, recorder.bodies, 1)
			require.JSONEq(t, `{"lifetimeActions":[]}`, recorder.bodies[0])
		})
	}
}

type requestBodyRecorder struct {
	transport policy.Transporter
	bodies    []string
}

func (r *requestBodyRecorder) Do(request *http.Request) (*http.Response, error) {
	if request.Body != nil {
		body, err := io.ReadAll(request.Body)
		if err != nil {
			return nil, err
		}
		request.Body = io.NopCloser(bytes.NewReader(body))
		r.bodies = append(r.bodies, string(body))
	}
	return r.transport.Do(request)
}
