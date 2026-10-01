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

func Test_readOverrideConflictsWithSensitiveBody(t *testing.T) {
	readOverride := types.ObjectValueMust(readOverrideAttributeTypes(), map[string]attr.Value{
		"method": types.StringValue("POST"),
		"action": types.StringValue("list"),
	})

	tests := []struct {
		name          string
		readOverride  types.Object
		sensitiveBody types.Dynamic
		want          bool
	}{
		{
			name:          "both configured",
			readOverride:  readOverride,
			sensitiveBody: types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{"secret": types.StringType}, map[string]attr.Value{"secret": types.StringValue("value")})),
			want:          true,
		},
		{
			name:          "sensitive body null",
			readOverride:  readOverride,
			sensitiveBody: types.DynamicNull(),
		},
		{
			name:          "read override null",
			readOverride:  types.ObjectNull(readOverrideAttributeTypes()),
			sensitiveBody: types.DynamicValue(types.StringValue("value")),
		},
		{
			name:          "sensitive body unknown",
			readOverride:  readOverride,
			sensitiveBody: types.DynamicUnknown(),
			want:          true,
		},
		{
			name:          "read override unknown",
			readOverride:  types.ObjectUnknown(readOverrideAttributeTypes()),
			sensitiveBody: types.DynamicValue(types.StringValue("value")),
			want:          true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := readOverrideConflictsWithSensitiveBody(test.readOverride, test.sensitiveBody); got != test.want {
				t.Fatalf("readOverrideConflictsWithSensitiveBody() = %t, want %t", got, test.want)
			}
		})
	}
}
