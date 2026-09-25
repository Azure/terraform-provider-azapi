package services

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDataPlaneResourceEphemeralSchemaRequestOptions(t *testing.T) {
	var response ephemeral.SchemaResponse
	resource := DataPlaneResourceEphemeral{}
	resource.Schema(context.Background(), ephemeral.SchemaRequest{}, &response)

	headers, ok := response.Schema.Attributes["headers"].(schema.MapAttribute)
	if !ok || !headers.IsOptional() || headers.ElementType != types.StringType {
		t.Fatalf("expected optional map(string) headers attribute, got %#v", response.Schema.Attributes["headers"])
	}
	queryParameters, ok := response.Schema.Attributes["query_parameters"].(schema.MapAttribute)
	expectedElementType := types.ListType{ElemType: types.StringType}
	if !ok || !queryParameters.IsOptional() || !queryParameters.ElementType.Equal(expectedElementType) {
		t.Fatalf("expected optional map(list(string)) query_parameters attribute, got %#v", response.Schema.Attributes["query_parameters"])
	}
}

func TestDataPlaneResourceEphemeralRequestOptions(t *testing.T) {
	model := &DataPlaneResourceEphemeralModel{
		Headers: types.MapValueMust(types.StringType, map[string]attr.Value{
			"x-ms-version": types.StringValue("2026-04-06"),
			"x-custom":     types.StringValue("custom-value"),
		}),
		QueryParameters: types.MapValueMust(types.ListType{ElemType: types.StringType}, map[string]attr.Value{
			"$filter": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("PartitionKey eq 'example'"),
			}),
			"$select": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("PartitionKey"),
				types.StringValue("RowKey"),
			}),
		}),
		Retry: retry.NewRetryValueNull(),
	}

	options := dataPlaneResourceEphemeralRequestOptions(model)

	expectedHeaders := map[string]string{
		"x-ms-version": "2026-04-06",
		"x-custom":     "custom-value",
	}
	if !reflect.DeepEqual(options.Headers, expectedHeaders) {
		t.Fatalf("expected headers %#v, got %#v", expectedHeaders, options.Headers)
	}
	expectedQueryParameters := map[string]string{
		"$filter": "PartitionKey eq 'example'",
		"$select": "PartitionKey,RowKey",
	}
	if !reflect.DeepEqual(options.QueryParameters, expectedQueryParameters) {
		t.Fatalf("expected query parameters %#v, got %#v", expectedQueryParameters, options.QueryParameters)
	}
	if options.RetryOptions != nil || options.LastRetryError != nil {
		t.Fatal("expected null retry configuration to remain disabled")
	}

	transport := &ephemeralRequestCaptureTransport{}
	dataPlaneClient, err := clients.NewDataPlaneClient(ephemeralStaticTokenCredential{}, &arm.ClientOptions{
		ClientOptions: policy.ClientOptions{
			Cloud:     cloud.AzurePublic,
			Transport: transport,
		},
	})
	if err != nil {
		t.Fatalf("building data plane client: %v", err)
	}
	if _, err := dataPlaneClient.Get(context.Background(), parse.DataPlaneResourceId{
		AzureResourceId: "management.azure.com/entity",
		ApiVersion:      "2026-04-06",
	}, options); err != nil {
		t.Fatalf("sending request: %v", err)
	}
	if got := transport.request.Header.Get("x-ms-version"); got != "2026-04-06" {
		t.Fatalf("expected x-ms-version header, got %q", got)
	}
	if got := transport.request.Header.Get("x-custom"); got != "custom-value" {
		t.Fatalf("expected custom header, got %q", got)
	}
	if got := transport.request.URL.Query().Get("$filter"); got != "PartitionKey eq 'example'" {
		t.Fatalf("expected $filter query parameter, got %q", got)
	}
	if got := transport.request.URL.Query().Get("$select"); got != "PartitionKey,RowKey" {
		t.Fatalf("expected $select query parameter, got %q", got)
	}

	model.Retry = retry.NewRetryValueMust(
		map[string]attr.Type{
			"interval_seconds":     types.Int64Type,
			"max_interval_seconds": types.Int64Type,
			"multiplier":           types.Float64Type,
			"randomization_factor": types.Float64Type,
			"error_message_regex":  types.ListType{ElemType: types.StringType},
		},
		map[string]attr.Value{
			"interval_seconds":     types.Int64Value(1),
			"max_interval_seconds": types.Int64Value(2),
			"multiplier":           types.Float64Value(1),
			"randomization_factor": types.Float64Value(0),
			"error_message_regex": types.ListValueMust(types.StringType, []attr.Value{
				types.StringValue("AuthorizationFailure"),
			}),
		},
	)
	options = dataPlaneResourceEphemeralRequestOptions(model)
	if options.RetryOptions == nil || options.LastRetryError == nil {
		t.Fatal("expected configured retry behavior to be retained")
	}
}

type ephemeralStaticTokenCredential struct{}

func (ephemeralStaticTokenCredential) GetToken(context.Context, policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		Token:     "test-token",
		ExpiresOn: time.Now().Add(time.Hour),
	}, nil
}

type ephemeralRequestCaptureTransport struct {
	request *http.Request
}

func (t *ephemeralRequestCaptureTransport) Do(request *http.Request) (*http.Response, error) {
	t.request = request
	return &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body:    io.NopCloser(strings.NewReader(`{}`)),
		Request: request,
	}, nil
}

func TestStorageTableEntityResponseStateShaping(t *testing.T) {
	responseBody := map[string]interface{}{
		"PartitionKey": "pk",
		"RowKey":       "rk",
		"Timestamp":    "2026-09-13T00:00:00Z",
		"odata.etag":   `W/"datetime'2026-09-13T00%3A00%3A00.0000000Z'"`,
		"outputs":      "remote-value",
	}
	for _, test := range []struct {
		name           string
		configuredBody map[string]interface{}
		expectedBody   map[string]interface{}
	}{
		{
			name: "server-only identity and metadata stay out of body",
			configuredBody: map[string]interface{}{
				"outputs": "configured-value",
			},
			expectedBody: map[string]interface{}{
				"outputs": "remote-value",
			},
		},
		{
			name: "explicit identity keys remain in body",
			configuredBody: map[string]interface{}{
				"PartitionKey": "pk",
				"RowKey":       "rk",
				"outputs":      "configured-value",
			},
			expectedBody: map[string]interface{}{
				"PartitionKey": "pk",
				"RowKey":       "rk",
				"outputs":      "remote-value",
			},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			managedBody := utils.UpdateObject(test.configuredBody, responseBody, utils.UpdateJsonOption{
				IgnoreMissingProperty: true,
			})
			if !reflect.DeepEqual(managedBody, test.expectedBody) {
				t.Fatalf("expected managed body %#v, got %#v", test.expectedBody, managedBody)
			}
		})
	}

	exportAll := types.DynamicValue(types.ListValueMust(types.StringType, []attr.Value{
		types.StringValue("*"),
	}))
	output, err := buildOutputFromBody(responseBody, exportAll, nil)
	if err != nil {
		t.Fatalf("building output: %v", err)
	}
	outputValue, ok := output.UnderlyingValue().(types.Dynamic)
	if !ok {
		t.Fatalf("expected dynamic output value, got %T", output.UnderlyingValue())
	}
	outputJSON, err := dynamic.ToJSON(outputValue)
	if err != nil {
		t.Fatalf("converting output to JSON: %v", err)
	}
	var actualOutput map[string]interface{}
	if err := json.Unmarshal(outputJSON, &actualOutput); err != nil {
		t.Fatalf("decoding output: %v", err)
	}
	if !reflect.DeepEqual(actualOutput, responseBody) {
		t.Fatalf("expected full response output %#v, got %#v", responseBody, actualOutput)
	}
}
