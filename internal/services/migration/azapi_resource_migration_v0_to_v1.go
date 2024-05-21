package migration

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AzapiResourceMigrationV0ToV1(ctx context.Context) resource.StateUpgrader {
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

				"removing_special_chars": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},

				"parent_id": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"type": schema.StringAttribute{
					Required: true,
				},

				"location": schema.StringAttribute{
					Optional: true,
					Computed: true,
				},

				"body": schema.StringAttribute{
					Optional: true,
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

				"schema_validation_enabled": schema.BoolAttribute{
					Optional: true,
					Computed: true,
				},

				"output": schema.StringAttribute{
					Computed: true,
				},

				"tags": schema.MapAttribute{
					ElementType: types.StringType,
					Optional:    true,
					Computed:    true,
				},
			},
			Blocks: map[string]schema.Block{
				"identity": schema.ListNestedBlock{
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"type": schema.StringAttribute{
								Required: true,
							},

							"identity_ids": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
							},

							"principal_id": schema.StringAttribute{
								Computed: true,
							},

							"tenant_id": schema.StringAttribute{
								Computed: true,
							},
						},
					},
				},

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
				ID                      types.String   `tfsdk:"id"`
				Name                    types.String   `tfsdk:"name"`
				ParentID                types.String   `tfsdk:"parent_id"`
				Type                    types.String   `tfsdk:"type"`
				Location                types.String   `tfsdk:"location"`
				Identity                types.List     `tfsdk:"identity"`
				Body                    types.String   `tfsdk:"body"`
				Locks                   types.List     `tfsdk:"locks"`
				RemovingSpecialChars    types.Bool     `tfsdk:"removing_special_chars"`
				SchemaValidationEnabled types.Bool     `tfsdk:"schema_validation_enabled"`
				IgnoreBodyChanges       types.List     `tfsdk:"ignore_body_changes"`
				IgnoreCasing            types.Bool     `tfsdk:"ignore_casing"`
				IgnoreMissingProperty   types.Bool     `tfsdk:"ignore_missing_property"`
				ResponseExportValues    types.List     `tfsdk:"response_export_values"`
				Output                  types.String   `tfsdk:"output"`
				Tags                    types.Map      `tfsdk:"tags"`
				Timeouts                timeouts.Value `tfsdk:"timeouts"`
			}
			type newModel struct {
				ID                      types.String   `tfsdk:"id"`
				Name                    types.String   `tfsdk:"name"`
				ParentID                types.String   `tfsdk:"parent_id"`
				Type                    types.String   `tfsdk:"type"`
				Location                types.String   `tfsdk:"location"`
				Identity                types.List     `tfsdk:"identity"`
				Body                    types.Dynamic  `tfsdk:"body"`
				Locks                   types.List     `tfsdk:"locks"`
				RemovingSpecialChars    types.Bool     `tfsdk:"removing_special_chars"`
				SchemaValidationEnabled types.Bool     `tfsdk:"schema_validation_enabled"`
				IgnoreBodyChanges       types.List     `tfsdk:"ignore_body_changes"`
				IgnoreCasing            types.Bool     `tfsdk:"ignore_casing"`
				IgnoreMissingProperty   types.Bool     `tfsdk:"ignore_missing_property"`
				ResponseExportValues    types.List     `tfsdk:"response_export_values"`
				Output                  types.Dynamic  `tfsdk:"output"`
				Tags                    types.Map      `tfsdk:"tags"`
				Timeouts                timeouts.Value `tfsdk:"timeouts"`
			}

			var oldState OldModel
			if response.Diagnostics.Append(request.State.Get(ctx, &oldState)...); response.Diagnostics.HasError() {
				return
			}

			newState := newModel{
				ID:                      oldState.ID,
				Name:                    oldState.Name,
				ParentID:                oldState.ParentID,
				Type:                    oldState.Type,
				Location:                oldState.Location,
				Identity:                oldState.Identity,
				Body:                    types.DynamicValue(types.StringValue(oldState.Body.ValueString())),
				Locks:                   oldState.Locks,
				RemovingSpecialChars:    oldState.RemovingSpecialChars,
				SchemaValidationEnabled: oldState.SchemaValidationEnabled,
				IgnoreBodyChanges:       oldState.IgnoreBodyChanges,
				IgnoreCasing:            oldState.IgnoreCasing,
				IgnoreMissingProperty:   oldState.IgnoreMissingProperty,
				ResponseExportValues:    oldState.ResponseExportValues,
				Output:                  types.DynamicValue(types.StringValue(oldState.Output.ValueString())),
				Tags:                    oldState.Tags,
				Timeouts:                oldState.Timeouts,
			}

			response.Diagnostics.Append(response.State.Set(ctx, newState)...)
		},
	}
}
