package services

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestNormalizeDataPlaneResourcePlanBody(t *testing.T) {
	planBody := types.DynamicValue(types.ObjectValueMust(
		map[string]attr.Type{
			"name":          types.StringType,
			"source_url":    types.StringType,
			"source_sha256": types.StringType,
		},
		map[string]attr.Value{
			"name":          types.StringValue("example-dataset"),
			"source_url":    types.StringUnknown(),
			"source_sha256": types.StringUnknown(),
		},
	))
	stateBody := types.DynamicValue(types.ObjectValueMust(
		map[string]attr.Type{
			"version":         types.StringType,
			"type":            types.StringType,
			"format":          types.StringType,
			"computed_sha256": types.StringType,
		},
		map[string]attr.Value{
			"version":         types.StringValue("1"),
			"type":            types.StringValue("uri_file"),
			"format":          types.StringValue("jsonl"),
			"computed_sha256": types.StringValue("provider-output"),
		},
	))

	normalizedBody, err := normalizeDataPlaneResourcePlanBody(
		context.Background(),
		"Microsoft.Foundry/datasets/versions@2025-05-01",
		planBody,
		stateBody,
	)
	if err != nil {
		t.Fatalf("normalizing Foundry plan body: %v", err)
	}

	attributes, ok := normalizedBody.UnderlyingValue().(types.Object)
	if !ok {
		t.Fatalf("unexpected normalized body type: %T", normalizedBody.UnderlyingValue())
	}
	values := attributes.Attributes()
	for field, expected := range map[string]string{
		"version": "1",
		"type":    "uri_file",
		"format":  "jsonl",
	} {
		value, ok := values[field].(types.String)
		if !ok || value.ValueString() != expected {
			t.Fatalf("expected %q to be copied from state, got %#v", field, values[field])
		}
	}
	for _, field := range []string{"source_url", "source_sha256"} {
		value, ok := values[field].(types.String)
		if !ok || !value.IsUnknown() {
			t.Fatalf("unknown %q was not preserved: %#v", field, values[field])
		}
	}
	if _, exists := values["computed_sha256"]; exists {
		t.Fatalf("provider output was copied into the request body: %#v", values)
	}

	unchangedBody, err := normalizeDataPlaneResourcePlanBody(
		context.Background(),
		"Microsoft.KeyVault/vaults/secrets@7.4",
		planBody,
		stateBody,
	)
	if err != nil {
		t.Fatalf("normalizing non-customized plan body: %v", err)
	}
	if !unchangedBody.Equal(planBody) {
		t.Fatalf("non-customized plan body changed: got %s, want %s", unchangedBody, planBody)
	}
}

func TestGetCreateResponseFuncOnlyForOptInCustomization(t *testing.T) {
	foundryFunc, foundryOK := getCreateResponseFunc(&DataPlaneResourceModel{
		Type: types.StringValue("Microsoft.Foundry/datasets/versions@2025-05-01"),
	})
	if !foundryOK || foundryFunc == nil {
		t.Fatal("expected Foundry dataset customization to expose a create response")
	}

	_, otherOK := getCreateResponseFunc(&DataPlaneResourceModel{
		Type: types.StringValue("Microsoft.KeyVault/vaults/secrets@7.4"),
	})
	if otherOK {
		t.Fatal("non-opted-in customizations must not expose a create response")
	}
}
