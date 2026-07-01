package migration

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AzapiResourceActionMigrationV2ToV3(ctx context.Context) resource.StateUpgrader {
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
				"ignore_not_found": schema.BoolAttribute{
					Optional: true,
				},
				"exist": schema.BoolAttribute{
					Computed: true,
				},
				"response_export_values": schema.DynamicAttribute{
					Optional: true,
				},
				"sensitive_response_export_values": schema.DynamicAttribute{
					Optional: true,
				},
				"output": schema.DynamicAttribute{
					Computed: true,
				},
				"sensitive_output": schema.DynamicAttribute{
					Computed:  true,
					Sensitive: true,
				},
				"retry": retry.RetrySchema(ctx),
				"headers": schema.MapAttribute{
					ElementType: types.StringType,
					Optional:    true,
				},
				"query_parameters": schema.MapAttribute{
					ElementType: types.ListType{
						ElemType: types.StringType,
					},
					Optional: true,
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
			Version: 2,
		},
		StateUpgrader: func(ctx context.Context, request resource.UpgradeStateRequest, response *resource.UpgradeStateResponse) {
			type oldModel struct {
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
			type newModel struct {
				ID                            types.String     `tfsdk:"id"`
				Type                          types.String     `tfsdk:"type"`
				ResourceId                    types.String     `tfsdk:"resource_id"`
				Action                        types.String     `tfsdk:"action"`
				Method                        types.String     `tfsdk:"method"`
				Body                          types.Dynamic    `tfsdk:"body"`
				SensitiveBodyVersion          types.Map        `tfsdk:"sensitive_body_version"`
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

			var oldState oldModel
			if response.Diagnostics.Append(request.State.Get(ctx, &oldState)...); response.Diagnostics.HasError() {
				return
			}

			newState := newModel{
				ID:                            oldState.ID,
				Type:                          oldState.Type,
				ResourceId:                    oldState.ResourceId,
				Action:                        oldState.Action,
				Method:                        oldState.Method,
				Body:                          oldState.Body,
				SensitiveBodyVersion:          types.MapNull(types.StringType),
				When:                          oldState.When,
				Locks:                         oldState.Locks,
				IgnoreNotFound:                oldState.IgnoreNotFound,
				Exist:                         oldState.Exist,
				ResponseExportValues:          oldState.ResponseExportValues,
				SensitiveResponseExportValues: oldState.SensitiveResponseExportValues,
				Output:                        oldState.Output,
				SensitiveOutput:               oldState.SensitiveOutput,
				Timeouts:                      oldState.Timeouts,
				Retry:                         oldState.Retry,
				Headers:                       oldState.Headers,
				QueryParameters:               oldState.QueryParameters,
			}

			response.Diagnostics.Append(response.State.Set(ctx, newState)...)
		},
	}
}
