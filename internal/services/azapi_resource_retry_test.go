package services_test

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/features"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestAzapiResourceCreateRetriesAuthorizationFailureDuringExistenceCheck(t *testing.T) {
	transport := &sequenceTransport{responses: []sequenceResponse{
		{
			method:     http.MethodGet,
			statusCode: http.StatusForbidden,
			body:       `{"error":{"code":"AuthorizationFailed","message":"role assignment has not propagated"}}`,
		},
		{
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
			body:       `{"error":{"code":"ResourceNotFound","message":"resource does not exist"}}`,
		},
		{
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       resourceGroupResponse,
		},
		{
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			body:       resourceGroupResponse,
		},
	}}

	diagnostics := createResourceWithRetry(t, transport, []string{"AuthorizationFailed"})
	if diagnostics.HasError() {
		t.Fatalf("creating resource: %+v", diagnostics)
	}

	assertRequestMethods(t, transport, []string{http.MethodGet, http.MethodGet, http.MethodPut, http.MethodGet})
}

func TestAzapiResourceCreateDoesNotRetryNotFoundDuringExistenceCheck(t *testing.T) {
	transport := &sequenceTransport{responses: []sequenceResponse{
		{
			method:     http.MethodGet,
			statusCode: http.StatusNotFound,
			body:       `{"error":{"code":"FooResourceNotFound","message":"resource does not exist"}}`,
		},
		{
			method:     http.MethodPut,
			statusCode: http.StatusOK,
			body:       resourceGroupResponse,
		},
		{
			method:     http.MethodGet,
			statusCode: http.StatusOK,
			body:       resourceGroupResponse,
		},
	}}

	diagnostics := createResourceWithRetry(t, transport, []string{"FooResourceNotFound"})
	if diagnostics.HasError() {
		t.Fatalf("creating resource: %+v", diagnostics)
	}

	assertRequestMethods(t, transport, []string{http.MethodGet, http.MethodPut, http.MethodGet})
}

const resourceGroupResponse = `{
	"id":"/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/acctest-rg",
	"name":"acctest-rg",
	"type":"Microsoft.Resources/resourceGroups",
	"location":"eastus",
	"properties":{}
}`

func createResourceWithRetry(t *testing.T, transport policy.Transporter, regexes []string) diag.Diagnostics {
	t.Helper()

	ctx := context.Background()
	resourceClient, err := clients.NewResourceClient(staticTokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{Transport: transport},
	})
	if err != nil {
		t.Fatalf("building resource client: %+v", err)
	}

	azapiResource := &services.AzapiResource{
		ProviderData: &clients.Client{
			Features:       features.Default(),
			ResourceClient: resourceClient,
		},
	}

	var schemaResponse resource.SchemaResponse
	azapiResource.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("building resource schema: %+v", schemaResponse.Diagnostics)
	}

	model := services.NewDefaultAzapiResourceModel()
	model.Name = types.StringValue("acctest-rg")
	model.ParentID = types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000")
	model.Type = types.StringValue("Microsoft.Resources/resourceGroups@2021-04-01")
	model.Location = types.StringValue("eastus")
	model.Body = types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{}))
	model.Retry = testRetryValue(regexes)

	plan := tfsdk.Plan{Schema: schemaResponse.Schema}
	if diagnostics := plan.Set(ctx, &model); diagnostics.HasError() {
		t.Fatalf("setting plan: %+v", diagnostics)
	}
	config := tfsdk.Config{
		Schema: schemaResponse.Schema,
		Raw:    plan.Raw,
	}

	state := tfsdk.State{
		Schema: schemaResponse.Schema,
		Raw:    tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil),
	}

	var diagnostics diag.Diagnostics
	azapiResource.CreateUpdate(ctx, config, plan, &state, &diagnostics, newTestPrivateData())
	return diagnostics
}

func testRetryValue(regexes []string) retry.RetryValue {
	regexValues := make([]attr.Value, len(regexes))
	for i, expression := range regexes {
		regexValues[i] = types.StringValue(expression)
	}

	return retry.NewRetryValueMust(
		map[string]attr.Type{
			"interval_seconds":     types.Int64Type,
			"max_interval_seconds": types.Int64Type,
			"multiplier":           types.Float64Type,
			"randomization_factor": types.Float64Type,
			"error_message_regex":  types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"interval_seconds":     types.Int64Value(1),
			"max_interval_seconds": types.Int64Value(1),
			"multiplier":           types.Float64Value(1),
			"randomization_factor": types.Float64Value(0),
			"error_message_regex":  types.ListValueMust(types.StringType, regexValues),
		},
	)
}

type sequenceResponse struct {
	method     string
	statusCode int
	body       string
}

type sequenceTransport struct {
	mu        sync.Mutex
	responses []sequenceResponse
	methods   []string
}

func (t *sequenceTransport) Do(request *http.Request) (*http.Response, error) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.methods = append(t.methods, request.Method)
	if len(t.responses) == 0 {
		return nil, fmt.Errorf("unexpected %s request", request.Method)
	}

	response := t.responses[0]
	t.responses = t.responses[1:]
	if request.Method != response.method {
		return nil, fmt.Errorf("expected %s request, got %s", response.method, request.Method)
	}

	return &http.Response{
		StatusCode: response.statusCode,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body:    io.NopCloser(strings.NewReader(response.body)),
		Request: request,
	}, nil
}

func assertRequestMethods(t *testing.T, transport *sequenceTransport, expected []string) {
	t.Helper()

	transport.mu.Lock()
	defer transport.mu.Unlock()

	if strings.Join(transport.methods, ",") != strings.Join(expected, ",") {
		t.Fatalf("expected request methods %v, got %v", expected, transport.methods)
	}
	if len(transport.responses) != 0 {
		t.Fatalf("expected all responses to be consumed, got %d remaining", len(transport.responses))
	}
}

type testPrivateData struct {
	values map[string][]byte
}

func newTestPrivateData() *testPrivateData {
	return &testPrivateData{values: make(map[string][]byte)}
}

func (d *testPrivateData) GetKey(_ context.Context, key string) ([]byte, diag.Diagnostics) {
	return d.values[key], nil
}

func (d *testPrivateData) SetKey(_ context.Context, key string, value []byte) diag.Diagnostics {
	if value == nil {
		delete(d.values, key)
	} else {
		d.values[key] = value
	}
	return nil
}
