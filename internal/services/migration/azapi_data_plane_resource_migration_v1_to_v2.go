package migration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AzapiDataPlaneResourceMigrationV1ToV2(ctx context.Context) resource.StateUpgrader {
	return resource.StateUpgrader{
		PriorSchema: &schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed: true,
				},

				"name": schema.StringAttribute{
					Required: true,
				},

				"parent_id": schema.StringAttribute{
					Required: true,
				},

				"type": schema.StringAttribute{
					Required: true,
				},

				"body": schema.DynamicAttribute{
					Optional: true,
					Computed: true,
				},

				"ignore_casing": schema.BoolAttribute{
					Optional: true,
					Computed: true},

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

				"output": schema.DynamicAttribute{
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
				Type                  types.String   `tfsdk:"type"`
				Body                  types.Dynamic  `tfsdk:"body"`
				IgnoreCasing          types.Bool     `tfsdk:"ignore_casing"`
				IgnoreMissingProperty types.Bool     `tfsdk:"ignore_missing_property"`
				ResponseExportValues  types.List     `tfsdk:"response_export_values"`
				Locks                 types.List     `tfsdk:"locks"`
				Output                types.Dynamic  `tfsdk:"output"`
				Timeouts              timeouts.Value `tfsdk:"timeouts"`
			}
			type newModel struct {
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

			var oldState OldModel
			if response.Diagnostics.Append(request.State.Get(ctx, &oldState)...); response.Diagnostics.HasError() {
				return
			}

			bodyVal, err := migrateToDynamicValue(oldState.Body)
			if err != nil {
				response.Diagnostics.AddError("failed to migrate body to dynamic value", err.Error())
				return
			}

			outputVal, err := migrateToDynamicValue(oldState.Output)
			if err != nil {
				response.Diagnostics.AddError("failed to migrate output to dynamic value", err.Error())
				return
			}
			newState := newModel{
				ID:                    oldState.ID,
				Name:                  oldState.Name,
				ParentID:              oldState.ParentID,
				Type:                  oldState.Type,
				Body:                  bodyVal,
				Locks:                 oldState.Locks,
				IgnoreCasing:          oldState.IgnoreCasing,
				IgnoreMissingProperty: oldState.IgnoreMissingProperty,
				ResponseExportValues:  oldState.ResponseExportValues,
				Output:                outputVal,
				Timeouts:              oldState.Timeouts,
			}

			response.Diagnostics.Append(response.State.Set(ctx, newState)...)
		},
	}
}
