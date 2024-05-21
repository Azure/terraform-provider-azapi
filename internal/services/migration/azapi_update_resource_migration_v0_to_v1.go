package migration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AzapiUpdateResourceMigrationV0ToV1(ctx context.Context) resource.StateUpgrader {
	return resource.StateUpgrader{
		PriorSchema: &schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed: true,
				},

				"name": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"parent_id": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"resource_id": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"type": schema.StringAttribute{
					Required: true,
				},

				"body": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"ignore_body_changes": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},

				"ignore_casing": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},

				"ignore_missing_property": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},

				"response_export_values": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},

				"locks": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},

				"output": schema.StringAttribute{
					Computed: true,
				},
			},

			Blocks: map[string]schema.Block{
				"timeouts": timeouts.Block(ctx, timeouts.Opts{
					Create: true,
					Update: true,
					Read:   true,
					Delete: true,
				}),
			},
			Version: 0,
		},
		StateUpgrader: func(ctx context.Context, request resource.UpgradeStateRequest, response *resource.UpgradeStateResponse) {
			type OldModel struct {
				ID                    types.String   `tfsdk:"id"`
				Name                  types.String   `tfsdk:"name"`
				ParentID              types.String   `tfsdk:"parent_id"`
				ResourceID            types.String   `tfsdk:"resource_id"`
				Type                  types.String   `tfsdk:"type"`
				Body                  types.String   `tfsdk:"body"`
				IgnoreCasing          types.Bool     `tfsdk:"ignore_casing"`
				IgnoreBodyChanges     types.List     `tfsdk:"ignore_body_changes"`
				IgnoreMissingProperty types.Bool     `tfsdk:"ignore_missing_property"`
				ResponseExportValues  types.List     `tfsdk:"response_export_values"`
				Locks                 types.List     `tfsdk:"locks"`
				Output                types.String   `tfsdk:"output"`
				Timeouts              timeouts.Value `tfsdk:"timeouts"`
			}
			type newModel struct {
				ID                    types.String   `tfsdk:"id"`
				Name                  types.String   `tfsdk:"name"`
				ParentID              types.String   `tfsdk:"parent_id"`
				ResourceID            types.String   `tfsdk:"resource_id"`
				Type                  types.String   `tfsdk:"type"`
				Body                  types.Dynamic  `tfsdk:"body"`
				IgnoreCasing          types.Bool     `tfsdk:"ignore_casing"`
				IgnoreBodyChanges     types.List     `tfsdk:"ignore_body_changes"`
				IgnoreMissingProperty types.Bool     `tfsdk:"ignore_missing_property"`
				ResponseExportValues  types.List     `tfsdk:"response_export_values"`
				Locks                 types.List     `tfsdk:"locks"`
				Output                types.Dynamic  `tfsdk:"output"`
				Timeouts              timeouts.Value `tfsdk:"timeouts"`
			}

			var oldState OldModel
			if response.Diagnostics.Append(request.State.Get(ctx, &oldState)...); response.Diagnostics.HasError() {
				return
			}

			newState := newModel{
				ID:                    oldState.ID,
				Name:                  oldState.Name,
				ParentID:              oldState.ParentID,
				ResourceID:            oldState.ResourceID,
				Type:                  oldState.Type,
				Body:                  types.DynamicValue(types.StringValue(oldState.Body.ValueString())),
				Locks:                 oldState.Locks,
				IgnoreBodyChanges:     oldState.IgnoreBodyChanges,
				IgnoreCasing:          oldState.IgnoreCasing,
				IgnoreMissingProperty: oldState.IgnoreMissingProperty,
				ResponseExportValues:  oldState.ResponseExportValues,
				Output:                types.DynamicValue(types.StringValue(oldState.Output.ValueString())),
				Timeouts:              oldState.Timeouts,
			}

			response.Diagnostics.Append(response.State.Set(ctx, newState)...)
		},
	}
}
