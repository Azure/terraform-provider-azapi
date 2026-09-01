package services

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Test_readOverrideFromObject(t *testing.T) {
	value := types.ObjectValueMust(readOverrideAttributeTypes(), map[string]attr.Value{
		"method": types.StringValue("POST"),
		"action": types.StringValue("list"),
	})

	model, diags := readOverrideFromObject(context.Background(), value)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if model.Method.ValueString() != "POST" || model.Action.ValueString() != "list" {
		t.Fatalf("unexpected read override: %#v", model)
	}

	model, diags = readOverrideFromObject(context.Background(), types.ObjectNull(readOverrideAttributeTypes()))
	if diags.HasError() || model != nil {
		t.Fatalf("null read override should decode to nil without diagnostics")
	}
}
