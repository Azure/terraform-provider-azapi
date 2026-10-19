package services_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/require"
)

const (
	keyVaultKeyType               = "Microsoft.KeyVault/vaults/keys@7.4"
	keyVaultKeyRotationPolicyType = "Microsoft.KeyVault/vaults/keys/rotationpolicy@7.6"
)

func TestDataPlaneResourceModifyPlanKeyVaultKeyReplaceTriggers(t *testing.T) {
	rsaKey := `{"kty":"RSA","key_size":2048,"key_ops":["encrypt","decrypt"],"attributes":{"exp":1830297600},"tags":{"env":"test"}}`

	testCases := []struct {
		name            string
		resourceType    string
		stateBody       types.Dynamic
		planBody        types.Dynamic
		requiresReplace bool
	}{
		{
			name:            "key type changed",
			resourceType:    keyVaultKeyType,
			stateBody:       dynamicFromJSON(t, rsaKey),
			planBody:        dynamicFromJSON(t, `{"kty":"EC","crv":"P-256","key_ops":["encrypt","decrypt"],"attributes":{"exp":1830297600},"tags":{"env":"test"}}`),
			requiresReplace: true,
		},
		{
			name:            "key size changed",
			resourceType:    keyVaultKeyType,
			stateBody:       dynamicFromJSON(t, rsaKey),
			planBody:        dynamicFromJSON(t, `{"kty":"RSA","key_size":3072,"key_ops":["encrypt","decrypt"],"attributes":{"exp":1830297600},"tags":{"env":"test"}}`),
			requiresReplace: true,
		},
		{
			name:            "curve changed",
			resourceType:    keyVaultKeyType,
			stateBody:       dynamicFromJSON(t, `{"kty":"EC","crv":"P-256"}`),
			planBody:        dynamicFromJSON(t, `{"kty":"EC","crv":"P-384"}`),
			requiresReplace: true,
		},
		{
			name:            "updatable properties changed",
			resourceType:    keyVaultKeyType,
			stateBody:       dynamicFromJSON(t, rsaKey),
			planBody:        dynamicFromJSON(t, `{"kty":"RSA","key_size":2048,"key_ops":["sign","verify"],"attributes":{"exp":1861920000},"tags":{"env":"prod"}}`),
			requiresReplace: false,
		},
		{
			name:            "imported key has no body in state",
			resourceType:    keyVaultKeyType,
			stateBody:       types.DynamicNull(),
			planBody:        dynamicFromJSON(t, rsaKey),
			requiresReplace: false,
		},
		{
			name:            "create-only property removed from config",
			resourceType:    keyVaultKeyType,
			stateBody:       dynamicFromJSON(t, rsaKey),
			planBody:        dynamicFromJSON(t, `{"kty":"RSA","key_ops":["encrypt","decrypt"],"attributes":{"exp":1830297600},"tags":{"env":"test"}}`),
			requiresReplace: false,
		},
		{
			name:         "create-only property unknown in plan",
			resourceType: keyVaultKeyType,
			stateBody:    dynamicFromJSON(t, `{"kty":"RSA"}`),
			planBody: types.DynamicValue(types.ObjectValueMust(
				map[string]attr.Type{"kty": types.StringType},
				map[string]attr.Value{"kty": types.StringUnknown()},
			)),
			requiresReplace: false,
		},
		{
			name:            "resource type without create-only properties",
			resourceType:    "Microsoft.KeyVault/vaults/secrets@7.4",
			stateBody:       dynamicFromJSON(t, `{"kty":"RSA"}`),
			planBody:        dynamicFromJSON(t, `{"kty":"EC"}`),
			requiresReplace: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			r := &services.DataPlaneResource{}
			resourceSchema := dataPlaneResourceSchema(t, r)

			state := tfsdk.State{Schema: resourceSchema}
			require.False(t, state.Set(ctx, newDataPlaneResourceModel(tc.resourceType, "example.vault.azure.net", "test", tc.stateBody)).HasError())
			plan := tfsdk.Plan{Schema: resourceSchema}
			require.False(t, plan.Set(ctx, newDataPlaneResourceModel(tc.resourceType, "example.vault.azure.net", "test", tc.planBody)).HasError())

			request := resource.ModifyPlanRequest{
				Config: tfsdk.Config{Schema: resourceSchema, Raw: plan.Raw},
				Plan:   plan,
				State:  state,
			}
			response := resource.ModifyPlanResponse{Plan: plan}
			r.ModifyPlan(ctx, request, &response)
			require.False(t, response.Diagnostics.HasError(), "%+v", response.Diagnostics)

			require.Equal(t, tc.requiresReplace, response.RequiresReplace.Contains(path.Root("body")))
		})
	}
}

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

func dynamicFromJSON(t *testing.T, input string) types.Dynamic {
	t.Helper()
	out, err := dynamic.FromJSONImplied([]byte(input))
	require.NoError(t, err)
	return out
}

func newDataPlaneResourceWithTransport(t *testing.T, transport policy.Transporter) *services.DataPlaneResource {
	t.Helper()
	dataPlaneClient, err := clients.NewDataPlaneClient(staticTokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud:     cloud.AzurePublic,
			Transport: transport,
		},
	})
	require.NoError(t, err)

	return &services.DataPlaneResource{
		ProviderData: &clients.Client{
			Features:        features.Default(),
			DataPlaneClient: dataPlaneClient,
		},
	}
}

func dataPlaneResourceSchema(t *testing.T, r *services.DataPlaneResource) schema.Schema {
	t.Helper()
	var schemaResponse resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &schemaResponse)
	require.False(t, schemaResponse.Diagnostics.HasError(), "%+v", schemaResponse.Diagnostics)
	return schemaResponse.Schema
}

func newDataPlaneResourceModel(resourceType, parentID, name string, body types.Dynamic) services.DataPlaneResourceModel {
	nameValue := types.StringNull()
	if name != "" {
		nameValue = types.StringValue(name)
	}
	id, _ := parse.NewDataPlaneResourceId(name, parentID, resourceType)
	return services.DataPlaneResourceModel{
		ID:                            types.StringValue(id.ID()),
		Name:                          nameValue,
		ParentID:                      types.StringValue(parentID),
		Type:                          types.StringValue(resourceType),
		Body:                          body,
		SensitiveBody:                 types.DynamicNull(),
		SensitiveBodyVersion:          types.MapNull(types.StringType),
		IgnoreCasing:                  types.BoolValue(false),
		IgnoreMissingProperty:         types.BoolValue(true),
		ReplaceTriggersExternalValues: types.DynamicNull(),
		ReplaceTriggersRefs:           types.ListNull(types.StringType),
		ResponseExportValues:          types.DynamicNull(),
		Retry:                         retry.NewRetryValueNull(),
		Locks:                         types.ListNull(types.StringType),
		Output:                        types.DynamicNull(),
		Timeouts: timeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{
				"create": types.StringType,
				"update": types.StringType,
				"read":   types.StringType,
				"delete": types.StringType,
			}),
		},
		CreateHeaders:         types.MapNull(types.StringType),
		CreateQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		UpdateHeaders:         types.MapNull(types.StringType),
		UpdateQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		DeleteHeaders:         types.MapNull(types.StringType),
		DeleteQueryParameters: types.MapNull(types.ListType{ElemType: types.StringType}),
		ReadHeaders:           types.MapNull(types.StringType),
		ReadQueryParameters:   types.MapNull(types.ListType{ElemType: types.StringType}),
	}
}
