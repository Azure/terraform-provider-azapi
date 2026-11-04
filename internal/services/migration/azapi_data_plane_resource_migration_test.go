package migration_test

import (
	"context"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/services"
	"github.com/Azure/terraform-provider-azapi/internal/services/migration"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestAzapiDataPlaneResourceMigrationsIncludeEvaluationID(t *testing.T) {
	ctx := context.Background()
	currentSchema := dataPlaneResourceSchema(t, ctx)
	tests := []struct {
		name     string
		upgrader resource.StateUpgrader
		state    any
	}{
		{
			name:     "v0",
			upgrader: migration.AzapiDataPlaneResourceMigrationV0ToV2(ctx),
			state: dataPlaneResourceV0State{
				ID:                    types.StringValue("legacy-id"),
				Name:                  types.StringValue("legacy-resource"),
				ParentID:              types.StringValue("/subscriptions/example"),
				Type:                  types.StringValue("Microsoft.Resources/resourceGroups@2021-04-01"),
				Body:                  types.StringValue("{}"),
				IgnoreCasing:          types.BoolValue(false),
				IgnoreMissingProperty: types.BoolValue(true),
				ResponseExportValues:  types.ListNull(types.StringType),
				Locks:                 types.ListNull(types.StringType),
				Output:                types.StringValue("{}"),
				Timeouts:              nullTimeouts(),
			},
		},
		{
			name:     "v1",
			upgrader: migration.AzapiDataPlaneResourceMigrationV1ToV2(ctx),
			state: dataPlaneResourceV1State{
				ID:                    types.StringValue("legacy-id"),
				Name:                  types.StringValue("legacy-resource"),
				ParentID:              types.StringValue("/subscriptions/example"),
				Type:                  types.StringValue("Microsoft.Resources/resourceGroups@2021-04-01"),
				Body:                  types.DynamicNull(),
				IgnoreCasing:          types.BoolValue(false),
				IgnoreMissingProperty: types.BoolValue(true),
				ResponseExportValues:  types.ListNull(types.StringType),
				Locks:                 types.ListNull(types.StringType),
				Output:                types.DynamicNull(),
				Timeouts:              nullTimeouts(),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			priorState := tfsdk.State{Schema: *test.upgrader.PriorSchema}
			if diagnostics := priorState.Set(ctx, test.state); diagnostics.HasError() {
				t.Fatalf("setting prior state returned diagnostics: %v", diagnostics)
			}

			response := resource.UpgradeStateResponse{
				State: tfsdk.State{Schema: currentSchema},
			}
			test.upgrader.StateUpgrader(ctx, resource.UpgradeStateRequest{State: &priorState}, &response)
			if response.Diagnostics.HasError() {
				t.Fatalf("upgrading state returned diagnostics: %v", response.Diagnostics)
			}

			var upgraded services.DataPlaneResourceModel
			if diagnostics := response.State.Get(ctx, &upgraded); diagnostics.HasError() {
				t.Fatalf("reading upgraded state returned diagnostics: %v", diagnostics)
			}
			if !upgraded.EvaluationID.IsNull() {
				t.Fatalf("expected evaluation_id to be null, got %s", upgraded.EvaluationID.String())
			}
			if got := upgraded.Name.ValueString(); got != "legacy-resource" {
				t.Fatalf("expected name to be preserved, got %q", got)
			}
		})
	}
}

func dataPlaneResourceSchema(t *testing.T, ctx context.Context) schema.Schema {
	t.Helper()

	dataPlaneResource := &services.DataPlaneResource{}
	var schemaResponse resource.SchemaResponse
	dataPlaneResource.Schema(ctx, resource.SchemaRequest{}, &schemaResponse)
	if schemaResponse.Diagnostics.HasError() {
		t.Fatalf("reading data plane resource schema returned diagnostics: %v", schemaResponse.Diagnostics)
	}
	return schemaResponse.Schema
}

type dataPlaneResourceV0State struct {
	ID                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	ParentID              types.String   `tfsdk:"parent_id"`
	Type                  types.String   `tfsdk:"type"`
	Body                  types.String   `tfsdk:"body"`
	IgnoreCasing          types.Bool     `tfsdk:"ignore_casing"`
	IgnoreMissingProperty types.Bool     `tfsdk:"ignore_missing_property"`
	ResponseExportValues  types.List     `tfsdk:"response_export_values"`
	Locks                 types.List     `tfsdk:"locks"`
	Output                types.String   `tfsdk:"output"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}

type dataPlaneResourceV1State struct {
	ID                    types.String   `tfsdk:"id"`
	Name                  types.String   `tfsdk:"name"`
	ParentID              types.String   `tfsdk:"parent_id"`
	Type                  types.String   `tfsdk:"type"`
	Body                  types.Dynamic  `tfsdk:"body"`
	IgnoreCasing          types.Bool     `tfsdk:"ignore_casing"`
	IgnoreMissingProperty types.Bool     `tfsdk:"ignore_missing_property"`
	ResponseExportValues  types.List     `tfsdk:"response_export_values"`
	Locks                 types.List     `tfsdk:"locks"`
	Output                types.Dynamic  `tfsdk:"output"`
	Timeouts              timeouts.Value `tfsdk:"timeouts"`
}
