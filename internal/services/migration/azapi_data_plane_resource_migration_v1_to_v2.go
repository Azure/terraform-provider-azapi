package migration

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/retry"
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
				ID                            types.String        `tfsdk:"id"`
				Name                          types.String        `tfsdk:"name"`
				ParentID                      types.String        `tfsdk:"parent_id"`
				Type                          types.String        `tfsdk:"type"`
				Body                          types.Dynamic       `tfsdk:"body"`
				IgnoreCasing                  types.Bool          `tfsdk:"ignore_casing"`
				IgnoreMissingProperty         types.Bool          `tfsdk:"ignore_missing_property"`
				ResponseExportValues          types.Dynamic       `tfsdk:"response_export_values"`
				ReplaceTriggersExternalValues types.Dynamic       `tfsdk:"replace_triggers_external_values"`
				ReplaceTriggersRefs           types.List          `tfsdk:"replace_triggers_refs"`
				Retry                         retry.RetryValue    `tfsdk:"retry"`
				Locks                         types.List          `tfsdk:"locks"`
				Output                        types.Dynamic       `tfsdk:"output"`
				Timeouts                      timeouts.Value      `tfsdk:"timeouts"`
				CreateHeaders                 map[string]string   `tfsdk:"create_headers"`
				CreateQueryParameters         map[string][]string `tfsdk:"create_query_parameters"`
				UpdateHeaders                 map[string]string   `tfsdk:"update_headers"`
				UpdateQueryParameters         map[string][]string `tfsdk:"update_query_parameters"`
				DeleteHeaders                 map[string]string   `tfsdk:"delete_headers"`
				DeleteQueryParameters         map[string][]string `tfsdk:"delete_query_parameters"`
				ReadHeaders                   map[string]string   `tfsdk:"read_headers"`
				ReadQueryParameters           map[string][]string `tfsdk:"read_query_parameters"`
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

			responseExportValues := types.DynamicNull()
			if !oldState.ResponseExportValues.IsNull() {
				responseExportValues = types.DynamicValue(oldState.ResponseExportValues)
			}

			newState := newModel{
				ID:                            oldState.ID,
				Name:                          oldState.Name,
				ParentID:                      oldState.ParentID,
				Type:                          oldState.Type,
				Body:                          bodyVal,
				Locks:                         oldState.Locks,
				IgnoreCasing:                  oldState.IgnoreCasing,
				IgnoreMissingProperty:         oldState.IgnoreMissingProperty,
				ResponseExportValues:          responseExportValues,
				ReplaceTriggersRefs:           types.ListNull(types.StringType),
				ReplaceTriggersExternalValues: types.DynamicNull(),
				Retry:                         retry.NewRetryValueNull(),
				Output:                        outputVal,
				Timeouts:                      oldState.Timeouts,
			}

			response.Diagnostics.Append(response.State.Set(ctx, newState)...)
		},
	}
}
