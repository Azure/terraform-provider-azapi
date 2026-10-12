package services_test

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	datasourcetimeouts "github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestDataPlaneResourceDataSourceReadFoundryEvaluationLookups(t *testing.T) {
	const parentID = "account.services.ai.azure.com/api/projects/project"

	tests := []struct {
		name             string
		resourceType     string
		resourceName     string
		evaluationID     string
		responseBody     string
		wantURL          string
		wantResourceID   string
		wantEvaluationID string
	}{
		{
			name:           "evaluation",
			resourceType:   "Microsoft.Foundry/evaluation/versions@2025-05-01",
			resourceName:   "eval_123",
			responseBody:   `{"id":"eval_123","name":"evaluation"}`,
			wantURL:        "https://" + parentID + "/openai/v1/evals/eval_123",
			wantResourceID: parentID + "/openai/v1/evals/eval_123",
		},
		{
			name:             "evaluation run",
			resourceType:     "Microsoft.Foundry/evaluation/runs@2025-05-01",
			resourceName:     "run_456",
			evaluationID:     "eval_123",
			responseBody:     `{"id":"run_456","status":"completed"}`,
			wantURL:          "https://" + parentID + "/openai/v1/evals/eval_123/runs/run_456",
			wantResourceID:   parentID + "/openai/v1/evals/eval_123/runs/run_456",
			wantEvaluationID: "eval_123",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := context.Background()
			transport := &dataPlaneDataSourceReadTransport{responseBody: test.responseBody}
			dataPlaneClient, err := clients.NewDataPlaneClient(staticTokenCredential{}, &arm.ClientOptions{
				ClientOptions: policy.ClientOptions{
					Cloud:     cloud.AzurePublic,
					Transport: transport,
				},
			})
			if err != nil {
				t.Fatalf("creating data-plane client: %v", err)
			}

			dataSource := &services.DataPlaneResourceDataSource{
				ProviderData: &clients.Client{DataPlaneClient: dataPlaneClient},
			}
			var schemaResponse datasource.SchemaResponse
			dataSource.Schema(ctx, datasource.SchemaRequest{}, &schemaResponse)
			if schemaResponse.Diagnostics.HasError() {
				t.Fatalf("building data source schema: %+v", schemaResponse.Diagnostics)
			}

			model := services.DataPlaneResourceDataSourceModel{
				ID:                   types.StringNull(),
				Name:                 types.StringValue(test.resourceName),
				EvaluationID:         types.StringNull(),
				ParentID:             types.StringValue(parentID),
				Type:                 types.StringValue(test.resourceType),
				Body:                 types.DynamicNull(),
				ResponseExportValues: types.DynamicNull(),
				Output:               types.DynamicNull(),
				Timeouts: datasourcetimeouts.Value{
					Object: types.ObjectNull(map[string]attr.Type{"read": types.StringType}),
				},
				Retry: retry.NewRetryValueNull(),
			}
			if test.evaluationID != "" {
				model.EvaluationID = types.StringValue(test.evaluationID)
			}

			plan := tfsdk.Plan{Schema: schemaResponse.Schema}
			if diagnostics := plan.Set(ctx, &model); diagnostics.HasError() {
				t.Fatalf("setting data source config: %+v", diagnostics)
			}
			config := tfsdk.Config{Schema: schemaResponse.Schema, Raw: plan.Raw}

			var validateResponse datasource.ValidateConfigResponse
			dataSource.ValidateConfig(ctx, datasource.ValidateConfigRequest{Config: config}, &validateResponse)
			if validateResponse.Diagnostics.HasError() {
				t.Fatalf("validating data source config: %+v", validateResponse.Diagnostics)
			}

			readResponse := datasource.ReadResponse{
				State: tfsdk.State{
					Schema: schemaResponse.Schema,
					Raw:    tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(ctx), nil),
				},
			}
			dataSource.Read(ctx, datasource.ReadRequest{Config: config}, &readResponse)
			if readResponse.Diagnostics.HasError() {
				t.Fatalf("reading data source: %+v", readResponse.Diagnostics)
			}

			if len(transport.requests) != 1 {
				t.Fatalf("expected one GET request, got %d", len(transport.requests))
			}
			gotRequest := transport.requests[0]
			if gotRequest.method != http.MethodGet {
				t.Fatalf("expected GET request, got %q", gotRequest.method)
			}
			if gotRequest.url != test.wantURL {
				t.Fatalf("unexpected request URL: got %q, want %q", gotRequest.url, test.wantURL)
			}

			var result services.DataPlaneResourceDataSourceModel
			if diagnostics := readResponse.State.Get(ctx, &result); diagnostics.HasError() {
				t.Fatalf("reading data source state: %+v", diagnostics)
			}
			if result.ID.ValueString() != test.wantResourceID {
				t.Errorf("unexpected resource ID: got %q, want %q", result.ID.ValueString(), test.wantResourceID)
			}
			if result.Name.ValueString() != test.resourceName {
				t.Errorf("unexpected name: got %q, want %q", result.Name.ValueString(), test.resourceName)
			}
			if test.wantEvaluationID == "" {
				if !result.EvaluationID.IsNull() {
					t.Errorf("expected null evaluation_id, got %q", result.EvaluationID.ValueString())
				}
			} else if result.EvaluationID.ValueString() != test.wantEvaluationID {
				t.Errorf("unexpected evaluation_id: got %q, want %q", result.EvaluationID.ValueString(), test.wantEvaluationID)
			}
		})
	}
}

func TestDataPlaneResourceDataSourceValidateConfigRequiresEvaluationIDForRuns(t *testing.T) {
	ctx := context.Background()
	dataSource := &services.DataPlaneResourceDataSource{}
	var schemaResponse datasource.SchemaResponse
	dataSource.Schema(ctx, datasource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("building data source schema: %+v", schemaResponse.Diagnostics)
	}

	model := services.DataPlaneResourceDataSourceModel{
		ID:                   types.StringNull(),
		Name:                 types.StringValue("run_456"),
		EvaluationID:         types.StringNull(),
		ParentID:             types.StringValue("account.services.ai.azure.com/api/projects/project"),
		Type:                 types.StringValue("Microsoft.Foundry/evaluation/runs@2025-05-01"),
		Body:                 types.DynamicNull(),
		ResponseExportValues: types.DynamicNull(),
		Output:               types.DynamicNull(),
		Timeouts: datasourcetimeouts.Value{
			Object: types.ObjectNull(map[string]attr.Type{"read": types.StringType}),
		},
		Retry: retry.NewRetryValueNull(),
	}
	plan := tfsdk.Plan{Schema: schemaResponse.Schema}
	if diagnostics := plan.Set(ctx, &model); diagnostics.HasError() {
		t.Fatalf("setting data source config: %+v", diagnostics)
	}

	var response datasource.ValidateConfigResponse
	dataSource.ValidateConfig(ctx, datasource.ValidateConfigRequest{
		Config: tfsdk.Config{Schema: schemaResponse.Schema, Raw: plan.Raw},
	}, &response)
	if !response.Diagnostics.HasError() {
		t.Fatal("expected a validation error when evaluation_id is omitted for a run lookup")
	}
}

type dataPlaneDataSourceReadRequest struct {
	method string
	url    string
}

type dataPlaneDataSourceReadTransport struct {
	responseBody string
	requests     []dataPlaneDataSourceReadRequest
}

func (transport *dataPlaneDataSourceReadTransport) Do(request *http.Request) (*http.Response, error) {
	transport.requests = append(transport.requests, dataPlaneDataSourceReadRequest{
		method: request.Method,
		url:    request.URL.String(),
	})

	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(transport.responseBody)),
		Request:    request,
	}, nil
}
