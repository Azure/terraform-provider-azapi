package migration_test

import (
	"context"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/Azure/terraform-provider-azapi/internal/services/migration"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAzapiResourceActionMigrationsIncludeSensitiveBody(t *testing.T) {
	ctx := context.Background()
	currentSchema := actionResourceSchema(t, ctx)

	tests := []struct {
		name     string
		upgrader resource.StateUpgrader
		state    any
	}{
		{
			name:     "v0",
			upgrader: migration.AzapiResourceActionMigrationV0ToV2(ctx),
			state: actionResourceV0State{
				ID:                   types.StringValue("test"),
				Type:                 types.StringValue("Microsoft.Resources/resourceGroups@2021-04-01"),
				ResourceId:           types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000"),
				Action:               types.StringValue("read"),
				Method:               types.StringValue("GET"),
				Body:                 types.StringNull(),
				When:                 types.StringValue("apply"),
				Locks:                types.ListNull(types.StringType),
				ResponseExportValues: types.ListNull(types.StringType),
				Output:               types.StringNull(),
				Timeouts:             nullTimeouts(),
			},
		},
		{
			name:     "v1",
			upgrader: migration.AzapiResourceActionMigrationV1ToV2(ctx),
			state: actionResourceV1State{
				ID:                   types.StringValue("test"),
				Type:                 types.StringValue("Microsoft.Resources/resourceGroups@2021-04-01"),
				ResourceId:           types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000"),
				Action:               types.StringValue("read"),
				Method:               types.StringValue("GET"),
				Body:                 types.DynamicNull(),
				When:                 types.StringValue("apply"),
				Locks:                types.ListNull(types.StringType),
				ResponseExportValues: types.ListNull(types.StringType),
				Output:               types.DynamicNull(),
				Timeouts:             nullTimeouts(),
			},
		},
		{
			name:     "v2",
			upgrader: migration.AzapiResourceActionMigrationV2ToV3(ctx),
			state: actionResourceV2State{
				ID:                            types.StringValue("test"),
				Type:                          types.StringValue("Microsoft.Resources/resourceGroups@2021-04-01"),
				ResourceId:                    types.StringValue("/subscriptions/00000000-0000-0000-0000-000000000000"),
				Action:                        types.StringValue("read"),
				Method:                        types.StringValue("GET"),
				Body:                          types.DynamicNull(),
				When:                          types.StringValue("apply"),
				Locks:                         types.ListNull(types.StringType),
				IgnoreNotFound:                types.BoolValue(false),
				Exist:                         types.BoolValue(true),
				ResponseExportValues:          types.DynamicNull(),
				SensitiveResponseExportValues: types.DynamicNull(),
				Output:                        types.DynamicNull(),
				SensitiveOutput:               types.DynamicNull(),
				Timeouts:                      nullTimeouts(),
				Retry:                         retry.NewRetryValueNull(),
				Headers:                       types.MapNull(types.StringType),
				QueryParameters:               types.MapNull(types.ListType{ElemType: types.StringType}),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priorState := tfsdk.State{Schema: *tt.upgrader.PriorSchema}
			if diags := priorState.Set(ctx, tt.state); diags.HasError() {
				t.Fatalf("setting prior state returned diagnostics: %v", diags)
			}

			response := resource.UpgradeStateResponse{
				State: tfsdk.State{Schema: currentSchema},
			}
			tt.upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{State: &priorState}, &response)
			if response.Diagnostics.HasError() {
				t.Fatalf("upgrading state returned diagnostics: %v", response.Diagnostics)
			}

			var upgraded services.ActionResourceModel
			if diags := response.State.Get(ctx, &upgraded); diags.HasError() {
				t.Fatalf("reading upgraded state returned diagnostics: %v", diags)
			}
			if !upgraded.SensitiveBody.IsNull() {
				t.Fatalf("expected sensitive_body to be null, got %s", upgraded.SensitiveBody.String())
			}
			if !upgraded.SensitiveBodyVersion.IsNull() {
				t.Fatalf("expected sensitive_body_version to be null, got %s", upgraded.SensitiveBodyVersion.String())
			}
		})
	}
}

func actionResourceSchema(t *testing.T, ctx context.Context) schema.Schema {
	t.Helper()

	var actionResource services.ActionResource
	var response resource.SchemaResponse
	actionResource.Schema(ctx, resource.SchemaRequest{}, &response)
	if response.Diagnostics.HasError() {
		t.Fatalf("reading action resource schema returned diagnostics: %v", response.Diagnostics)
	}
	return response.Schema
}

func nullTimeouts() timeouts.Value {
	return timeouts.Value{
		Object: types.ObjectNull(map[string]attr.Type{
			"create": types.StringType,
			"update": types.StringType,
			"read":   types.StringType,
			"delete": types.StringType,
		}),
	}
}

type actionResourceV0State struct {
	ID                   types.String   `tfsdk:"id"`
	Type                 types.String   `tfsdk:"type"`
	ResourceId           types.String   `tfsdk:"resource_id"`
	Action               types.String   `tfsdk:"action"`
	Method               types.String   `tfsdk:"method"`
	Body                 types.String   `tfsdk:"body"`
	When                 types.String   `tfsdk:"when"`
	Locks                types.List     `tfsdk:"locks"`
	ResponseExportValues types.List     `tfsdk:"response_export_values"`
	Output               types.String   `tfsdk:"output"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type actionResourceV1State struct {
	ID                   types.String   `tfsdk:"id"`
	Type                 types.String   `tfsdk:"type"`
	ResourceId           types.String   `tfsdk:"resource_id"`
	Action               types.String   `tfsdk:"action"`
	Method               types.String   `tfsdk:"method"`
	Body                 types.Dynamic  `tfsdk:"body"`
	When                 types.String   `tfsdk:"when"`
	Locks                types.List     `tfsdk:"locks"`
	ResponseExportValues types.List     `tfsdk:"response_export_values"`
	Output               types.Dynamic  `tfsdk:"output"`
	Timeouts             timeouts.Value `tfsdk:"timeouts"`
}

type actionResourceV2State struct {
	ID                            types.String     `tfsdk:"id"`
	Type                          types.String     `tfsdk:"type"`
	ResourceId                    types.String     `tfsdk:"resource_id"`
	Action                        types.String     `tfsdk:"action"`
	Method                        types.String     `tfsdk:"method"`
	Body                          types.Dynamic    `tfsdk:"body"`
	When                          types.String     `tfsdk:"when"`
	Locks                         types.List       `tfsdk:"locks"`
	IgnoreNotFound                types.Bool       `tfsdk:"ignore_not_found"`
	Exist                         types.Bool       `tfsdk:"exist"`
	ResponseExportValues          types.Dynamic    `tfsdk:"response_export_values"`
	SensitiveResponseExportValues types.Dynamic    `tfsdk:"sensitive_response_export_values"`
	Output                        types.Dynamic    `tfsdk:"output"`
	SensitiveOutput               types.Dynamic    `tfsdk:"sensitive_output"`
	Timeouts                      timeouts.Value   `tfsdk:"timeouts"`
	Retry                         retry.RetryValue `tfsdk:"retry"`
	Headers                       types.Map        `tfsdk:"headers"`
	QueryParameters               types.Map        `tfsdk:"query_parameters"`
}
