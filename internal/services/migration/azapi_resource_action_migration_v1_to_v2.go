package migration

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AzapiResourceActionMigrationV1ToV2(ctx context.Context) resource.StateUpgrader {
	return resource.StateUpgrader{
		PriorSchema: &schema.Schema{
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed: true,
				},

				"type": schema.StringAttribute{
					Required: true,
				},

				"resource_id": schema.StringAttribute{
					Required: true,
				},

				"action": schema.StringAttribute{
					Optional: true,
				},

				"method": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"body": schema.DynamicAttribute{
					Optional: true,
				},

				"when": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"locks": schema.ListAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},

				"response_export_values": schema.ListAttribute{
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
			type newModel struct {
				ID                            types.String        `tfsdk:"id"`
				Type                          types.String        `tfsdk:"type"`
				ResourceId                    types.String        `tfsdk:"resource_id"`
				Action                        types.String        `tfsdk:"action"`
				Method                        types.String        `tfsdk:"method"`
				Body                          types.Dynamic       `tfsdk:"body"`
				When                          types.String        `tfsdk:"when"`
				Locks                         types.List          `tfsdk:"locks"`
				ResponseExportValues          types.Dynamic       `tfsdk:"response_export_values"`
				SensitiveResponseExportValues types.Dynamic       `tfsdk:"sensitive_response_export_values"`
				Output                        types.Dynamic       `tfsdk:"output"`
				SensitiveOutput               types.Dynamic       `tfsdk:"sensitive_output"`
				Timeouts                      timeouts.Value      `tfsdk:"timeouts"`
				Retry                         retry.RetryValue    `tfsdk:"retry"`
				Headers                       map[string]string   `tfsdk:"headers"`
				QueryParameters               map[string][]string `tfsdk:"query_parameters"`
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
				Type:                          oldState.Type,
				ResourceId:                    oldState.ResourceId,
				Action:                        oldState.Action,
				Method:                        oldState.Method,
				Body:                          bodyVal,
				When:                          oldState.When,
				Locks:                         oldState.Locks,
				ResponseExportValues:          responseExportValues,
				SensitiveResponseExportValues: types.DynamicNull(),
				Output:                        outputVal,
				SensitiveOutput:               types.DynamicNull(),
				Timeouts:                      oldState.Timeouts,
				Retry:                         retry.NewRetryValueNull(),
			}

			response.Diagnostics.Append(response.State.Set(ctx, newState)...)
		},
	}
}
