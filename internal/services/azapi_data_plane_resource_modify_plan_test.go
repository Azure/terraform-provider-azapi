package services

import (
	"context"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestDataPlaneResourceModifyPlanFoundryEvaluationReplacement(t *testing.T) {
	const initialBody = `{
		"name": "evaluation",
		"metadata": {"team": "platform"},
		"data_source_config": {
			"type": "custom",
			"include_sample_schema": false,
			"item_schema": {"type": "object", "properties": {"input": {"type": "string"}}}
		},
		"testing_criteria": [{"name": "TaskCompletion", "data_mapping": {"query": "{{item.input}}"}}]
	}`

	tests := []struct {
		name        string
		plannedBody string
		replace     bool
	}{
		{
			name: "changing include_sample_schema requires replacement",
			plannedBody: `{
				"name": "evaluation",
				"metadata": {"team": "platform"},
				"data_source_config": {
					"type": "custom",
					"include_sample_schema": true,
					"item_schema": {"type": "object", "properties": {"input": {"type": "string"}}}
				},
				"testing_criteria": [{"name": "TaskCompletion", "data_mapping": {"query": "{{item.input}}"}}]
			}`,
			replace: true,
		},
		{
			name: "changing nested data source configuration requires replacement",
			plannedBody: `{
				"name": "evaluation",
				"metadata": {"team": "platform"},
				"data_source_config": {
					"type": "custom",
					"include_sample_schema": false,
					"item_schema": {"type": "object", "properties": {"input": {"type": "number"}}}
				},
				"testing_criteria": [{"name": "TaskCompletion", "data_mapping": {"query": "{{item.input}}"}}]
			}`,
			replace: true,
		},
		{
			name: "changing testing criteria requires replacement",
			plannedBody: `{
				"name": "evaluation",
				"metadata": {"team": "platform"},
				"data_source_config": {
					"type": "custom",
					"include_sample_schema": false,
					"item_schema": {"type": "object", "properties": {"input": {"type": "string"}}}
				},
				"testing_criteria": [{"name": "Safety", "data_mapping": {"query": "{{item.input}}"}}]
			}`,
			replace: true,
		},
		{
			name: "changing name does not require replacement",
			plannedBody: `{
				"name": "renamed-evaluation",
				"metadata": {"team": "platform"},
				"data_source_config": {
					"type": "custom",
					"include_sample_schema": false,
					"item_schema": {"type": "object", "properties": {"input": {"type": "string"}}}
				},
				"testing_criteria": [{"name": "TaskCompletion", "data_mapping": {"query": "{{item.input}}"}}]
			}`,
			replace: false,
		},
		{
			name: "changing metadata does not require replacement",
			plannedBody: `{
				"name": "evaluation",
				"metadata": {"team": "application"},
				"data_source_config": {
					"type": "custom",
					"include_sample_schema": false,
					"item_schema": {"type": "object", "properties": {"input": {"type": "string"}}}
				},
				"testing_criteria": [{"name": "TaskCompletion", "data_mapping": {"query": "{{item.input}}"}}]
			}`,
			replace: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			response := modifyFoundryEvaluationPlan(t, initialBody, test.plannedBody)
			if response.Diagnostics.HasError() {
				t.Fatalf("modifying plan: %v", response.Diagnostics)
			}

			want := path.Paths{}
			if test.replace {
				want = path.Paths{path.Root("body")}
			}
			if got := response.RequiresReplace; len(got) != len(want) || (len(want) == 1 && got[0].String() != want[0].String()) {
				t.Fatalf("unexpected replacement paths: got %v, want %v", got, want)
			}
		})
	}
}

func TestDataPlaneResourceModifyPlanFoundryEvaluationUnknownValues(t *testing.T) {
	defaultStateBody := `{
		"name": "evaluation",
		"metadata": {"team": "platform"}
	}`
	tests := []struct {
		name      string
		stateBody string
		body      string
		replace   bool
	}{
		{
			name: "unknown data source configuration absent from state",
			body: `{
				"name": "evaluation",
				"metadata": {"team": "platform"},
				"data_source_config": "<unknown>"
			}`,
			replace: true,
		},
		{
			name: "unknown testing criteria absent from state",
			body: `{
				"name": "evaluation",
				"metadata": {"team": "platform"},
				"testing_criteria": "<unknown>"
			}`,
			replace: true,
		},
		{
			name: "unknown nested data source value replacing a null value",
			stateBody: `{
				"name": "evaluation",
				"metadata": {"team": "platform"},
				"data_source_config": {
					"type": "custom",
					"include_sample_schema": null
				}
			}`,
			body: `{
				"name": "evaluation",
				"metadata": {"team": "platform"},
				"data_source_config": {
					"type": "custom",
					"include_sample_schema": "<unknown>"
				}
			}`,
			replace: true,
		},
		{
			name: "unknown mutable metadata does not require replacement",
			body: `{
				"name": "evaluation",
				"metadata": "<unknown>"
			}`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stateBody := test.stateBody
			if stateBody == "" {
				stateBody = defaultStateBody
			}
			response := modifyFoundryEvaluationPlan(t, stateBody, test.body)
			if response.Diagnostics.HasError() {
				t.Fatalf("modifying plan: %v", response.Diagnostics)
			}

			want := path.Paths{}
			if test.replace {
				want = path.Paths{path.Root("body")}
			}
			if got := response.RequiresReplace; len(got) != len(want) || (len(want) == 1 && got[0].String() != want[0].String()) {
				t.Fatalf("unexpected replacement paths: got %v, want %v", got, want)
			}
		})
	}
}

func modifyFoundryEvaluationPlan(t *testing.T, stateBody, planBody string) resource.ModifyPlanResponse {
	t.Helper()

	ctx := context.Background()
	dataPlaneResource := &DataPlaneResource{}
	var schemaResponse resource.SchemaResponse
	dataPlaneResource.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("building data plane resource schema: %v", schemaResponse.Diagnostics)
	}

	statePlan := newFoundryEvaluationPlan(t, schemaResponse.Schema, stateBody)
	plan := newFoundryEvaluationPlan(t, schemaResponse.Schema, planBody)

	request := resource.ModifyPlanRequest{
		Config: tfsdk.Config{
			Schema: schemaResponse.Schema,
			Raw:    plan.Raw,
		},
		State: tfsdk.State{
			Schema: schemaResponse.Schema,
			Raw:    statePlan.Raw,
		},
		Plan: plan,
	}
	response := resource.ModifyPlanResponse{
		Plan: tfsdk.Plan{
			Schema: schemaResponse.Schema,
			Raw:    plan.Raw,
		},
	}

	dataPlaneResource.ModifyPlan(ctx, request, &response)
	return response
}

func newFoundryEvaluationPlan(t *testing.T, schema resourceschema.Schema, body string) tfsdk.Plan {
	t.Helper()

	ctx := context.Background()
	bodyValue, err := dynamic.FromJSONImplied([]byte(body))
	if err != nil {
		t.Fatalf("building body: %v", err)
	}

	emptyTimeouts := types.ObjectNull(map[string]attr.Type{
		"create": types.StringType,
		"update": types.StringType,
		"read":   types.StringType,
		"delete": types.StringType,
	})
	model := DataPlaneResourceModel{
		ID:                            types.StringValue("account.services.ai.azure.com/api/projects/project/openai/v1/evals/eval_123"),
		Name:                          types.StringNull(),
		EvaluationID:                  types.StringNull(),
		ParentID:                      types.StringValue("account.services.ai.azure.com/api/projects/project"),
		Type:                          types.StringValue("Microsoft.Foundry/evaluation/versions@2025-05-01"),
		Body:                          bodyValue,
		SensitiveBody:                 types.DynamicNull(),
		SensitiveBodyVersion:          types.MapValueMust(types.StringType, map[string]attr.Value{}),
		IgnoreCasing:                  types.BoolValue(false),
		IgnoreMissingProperty:         types.BoolValue(true),
		ReplaceTriggersExternalValues: types.DynamicNull(),
		ReplaceTriggersRefs:           types.ListNull(types.StringType),
		ResponseExportValues:          types.DynamicNull(),
		Retry: retry.RetryValue{
			ErrorMessageRegex:   types.ListNull(types.StringType),
			IntervalSeconds:     types.Int64Null(),
			MaxIntervalSeconds:  types.Int64Null(),
			Multiplier:          types.Float64Null(),
			RandomizationFactor: types.Float64Null(),
		},
		Locks:  types.ListNull(types.StringType),
		Output: types.DynamicNull(),
		Timeouts: timeouts.Value{
			Object: emptyTimeouts,
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

	plan := tfsdk.Plan{Schema: schema}
	if diagnostics := plan.Set(ctx, &model); diagnostics.HasError() {
		t.Fatalf("setting data plane resource plan: %v", diagnostics)
	}
	return plan
}
