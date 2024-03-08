package services

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/resource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ResourceIdDataSourceModel struct {
	ID                types.String   `tfsdk:"id"`
	Type              types.String   `tfsdk:"type"`
	Name              types.String   `tfsdk:"name"`
	ParentID          types.String   `tfsdk:"parent_id"`
	ResourceID        types.String   `tfsdk:"resource_id"`
	ResourceGroupName types.String   `tfsdk:"resource_group_name"`
	SubscriptionID    types.String   `tfsdk:"subscription_id"`
	ProviderNamespace types.String   `tfsdk:"provider_namespace"`
	Parts             types.Map      `tfsdk:"parts"`
	Timeouts          timeouts.Value `tfsdk:"timeouts"`
}

type ResourceIdDataSource struct {
}

var _ datasource.DataSource = &ResourceIdDataSource{}
var _ datasource.DataSourceWithValidateConfig = &ResourceIdDataSource{}

func (r *ResourceIdDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource_id"
}

func (r *ResourceIdDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
			},

			"parent_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
			},

			"name": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},

			"resource_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
			},

			"resource_group_name": schema.StringAttribute{
				Computed: true,
			},

			"subscription_id": schema.StringAttribute{
				Computed: true,
			},

			"provider_namespace": schema.StringAttribute{
				Computed: true,
			},

			"parts": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},

		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx, timeouts.Opts{
				Read: true,
			}),
		},
	}
}

func (r *ResourceIdDataSource) ValidateConfig(ctx context.Context, request datasource.ValidateConfigRequest, response *datasource.ValidateConfigResponse) {
	var config *ResourceIdDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &config)...); response.Diagnostics.HasError() {
		return
	}

	if config == nil {
		return
	}

	if config.Name.IsNull() && !config.ParentID.IsNull() {
		response.Diagnostics.AddError("Invalid configuration", `The argument "name" is required when the argument "parent_id" is set`)
	}
	if !config.Name.IsNull() && config.ParentID.IsNull() {
		response.Diagnostics.AddError("Invalid configuration", `The argument "parent_id" is required when the argument "name" is set`)
	}
	if config.Name.IsNull() && config.ResourceID.IsNull() {
		response.Diagnostics.AddError("Invalid configuration", `One of the arguments "name" or "resource_id" must be set`)
	}
	if !config.Name.IsNull() && !config.ResourceID.IsNull() {
		response.Diagnostics.AddError("Invalid configuration", `Only one of the arguments "name" or "resource_id" can be set`)
	}
	if response.Diagnostics.HasError() {
		return
	}
}

func (r *ResourceIdDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model ResourceIdDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	readTimeout, diags := model.Timeouts.Read(ctx, 5*time.Minute)
	response.Diagnostics.Append(diags...)
	if response.Diagnostics.HasError() {
		return
	}

	ctx, cancel := context.WithTimeout(ctx, readTimeout)
	defer cancel()

	var id parse.ResourceId
	if name := model.Name.ValueString(); len(name) != 0 {
		buildId, err := parse.NewResourceID(model.Name.ValueString(), model.ParentID.ValueString(), model.Type.ValueString())
		if err != nil {
			response.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
		id = buildId
	} else {
		buildId, err := parse.ResourceIDWithResourceType(model.ResourceID.ValueString(), model.Type.ValueString())
		if err != nil {
			response.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
		id = buildId
	}

	model.ID = basetypes.NewStringValue(id.ID())
	model.Name = basetypes.NewStringValue(id.Name)
	model.ParentID = basetypes.NewStringValue(id.ParentId)
	model.ResourceID = basetypes.NewStringValue(id.AzureResourceId)

	armId, err := arm.ParseResourceID(id.AzureResourceId)
	if id.AzureResourceId == "/" {
		armId, err = &arm.ResourceID{
			ResourceType: arm.TenantResourceType,
		}, nil
	}
	if err != nil {
		response.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	model.ResourceGroupName = basetypes.NewStringValue(armId.ResourceGroupName)
	model.SubscriptionID = basetypes.NewStringValue(armId.SubscriptionID)
	model.ProviderNamespace = basetypes.NewStringValue(armId.ResourceType.Namespace)

	path := id.AzureResourceId
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	components := strings.Split(path, "/")
	parts := make(map[string]attr.Value)
	for i := 0; i < len(components)-1; i += 2 {
		parts[components[i]] = basetypes.NewStringValue(components[i+1])
	}
	model.Parts = basetypes.NewMapValueMust(types.StringType, parts)

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
