package migration

import (
	"context"

	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func AzapiResourceMigrationV1ToV2(ctx context.Context) resource.StateUpgrader {
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

				"body": schema.DynamicAttribute{
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

				"output": schema.DynamicAttribute{
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
			type newModel struct {
				ID                            types.String        `tfsdk:"id"`
				Name                          types.String        `tfsdk:"name"`
				ParentID                      types.String        `tfsdk:"parent_id"`
				Type                          types.String        `tfsdk:"type"`
				Location                      types.String        `tfsdk:"location"`
				Identity                      types.List          `tfsdk:"identity"`
				Body                          types.Dynamic       `tfsdk:"body"`
				Locks                         types.List          `tfsdk:"locks"`
				SchemaValidationEnabled       types.Bool          `tfsdk:"schema_validation_enabled"`
				IgnoreCasing                  types.Bool          `tfsdk:"ignore_casing"`
				IgnoreMissingProperty         types.Bool          `tfsdk:"ignore_missing_property"`
				ReplaceTriggersExternalValues types.Dynamic       `tfsdk:"replace_triggers_external_values"`
				ReplaceTriggersRefs           types.List          `tfsdk:"replace_triggers_refs"`
				ResponseExportValues          types.Dynamic       `tfsdk:"response_export_values"`
				Retry                         retry.RetryValue    `tfsdk:"retry"`
				Output                        types.Dynamic       `tfsdk:"output"`
				Tags                          types.Map           `tfsdk:"tags"`
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
				Location:                      oldState.Location,
				Identity:                      oldState.Identity,
				Body:                          bodyVal,
				Locks:                         oldState.Locks,
				SchemaValidationEnabled:       oldState.SchemaValidationEnabled,
				IgnoreCasing:                  oldState.IgnoreCasing,
				IgnoreMissingProperty:         oldState.IgnoreMissingProperty,
				ReplaceTriggersExternalValues: types.DynamicNull(),
				ReplaceTriggersRefs:           types.ListNull(types.StringType),
				ResponseExportValues:          responseExportValues,
				Retry:                         retry.NewRetryValueNull(),
				Output:                        outputVal,
				Tags:                          oldState.Tags,
				Timeouts:                      oldState.Timeouts,
			}

			response.Diagnostics.Append(response.State.Set(ctx, newState)...)
		},
	}
}

func migrateToDynamicValue(input types.Dynamic) (types.Dynamic, error) {
	if input.IsNull() {
		return input, nil
	}
	if input.IsUnderlyingValueNull() {
		return input, nil
	}
	stringVal, ok := input.UnderlyingValue().(types.String)
	if !ok {
		return input, nil
	}
	dynamicVal, err := dynamic.FromJSONImplied([]byte(stringVal.ValueString()))
	if err != nil {
		return input, err
	}
	return dynamicVal, nil
}
