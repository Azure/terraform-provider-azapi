package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/azure/tags"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type AzapiResourceDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	ParentID             types.String `tfsdk:"parent_id"`
	ResourceID           types.String `tfsdk:"resource_id"`
	Type                 types.String `tfsdk:"type"`
	ResponseExportValues types.List   `tfsdk:"response_export_values"`
	Location             types.String `tfsdk:"location"`
	Identity             types.List   `tfsdk:"identity"`
	Output               types.String `tfsdk:"output"`
	Tags                 types.Map    `tfsdk:"tags"`
}

type AzapiResourceDataSource struct {
	ProviderData *clients.Client
}

var _ datasource.DataSource = &AzapiResourceDataSource{}
var _ datasource.DataSourceWithConfigure = &AzapiResourceDataSource{}

func (r *AzapiResourceDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *AzapiResourceDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource"
}

func (r *AzapiResourceDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
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

			`name`: schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					myvalidator.StringIsNotEmpty(),
				},
			},

			"parent_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
			},

			"resource_id": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
			},

			"location": schema.StringAttribute{
				Computed: true,
			},

			"identity": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"type": schema.StringAttribute{
							Computed: true,
						},

						"principal_id": schema.StringAttribute{
							Computed: true,
						},

						"tenant_id": schema.StringAttribute{
							Computed: true,
						},

						"identity_ids": schema.ListAttribute{
							Computed:    true,
							ElementType: types.StringType,
						},
					},
				},
			},

			"response_export_values": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.ValueStringsAre(myvalidator.StringIsNotEmpty()),
				},
			},

			"output": schema.StringAttribute{
				Computed: true,
			},

			"tags": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *AzapiResourceDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model AzapiResourceDataSourceModel
	if response.Diagnostics.Append(request.Config.Get(ctx, &model)...); response.Diagnostics.HasError() {
		return
	}

	var id parse.ResourceId
	resourceType := model.Type.ValueString()
	azureResourceType, _, _ := utils.GetAzureResourceTypeApiVersion(resourceType)
	if name := model.Name.ValueString(); len(name) != 0 {
		parentId := model.ParentID.ValueString()
		if parentId == "" && strings.EqualFold(azureResourceType, arm.ResourceGroupResourceType.String()) {
			parentId = fmt.Sprintf("/subscriptions/%s", r.ProviderData.Account.GetSubscriptionId())
		}
		buildId, err := parse.NewResourceID(name, parentId, resourceType)
		if err != nil {
			response.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
		id = buildId
	} else {
		resourceId := model.ResourceID.ValueString()
		if resourceId == "" && strings.EqualFold(azureResourceType, arm.SubscriptionResourceType.String()) {
			resourceId = fmt.Sprintf("/subscriptions/%s", r.ProviderData.Account.GetSubscriptionId())
		}
		buildId, err := parse.ResourceIDWithResourceType(resourceId, resourceType)
		if err != nil {
			response.Diagnostics.AddError("Invalid configuration", err.Error())
			return
		}
		id = buildId
	}

	client := r.ProviderData.ResourceClient
	responseBody, err := client.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			response.Diagnostics.AddError("Resource not found", fmt.Errorf("resource %q not found", id).Error())
			return
		}
		response.Diagnostics.AddError("Failed to retrieve resource", fmt.Errorf("retrieving resource %q: %+v", id, err).Error())
		return
	}

	model.ID = basetypes.NewStringValue(id.ID())
	model.Name = basetypes.NewStringValue(id.Name)
	model.ParentID = basetypes.NewStringValue(id.ParentId)
	model.ResourceID = basetypes.NewStringValue(id.AzureResourceId)
	if bodyMap, ok := responseBody.(map[string]interface{}); ok {
		model.Tags = tags.FlattenTags(bodyMap["tags"])
		model.Location = basetypes.NewStringNull()
		if v := bodyMap["location"]; v != nil {
			model.Location = basetypes.NewStringValue(location.Normalize(v.(string)))
		}
		model.Identity = basetypes.NewListNull(identity.Model{}.ModelType())
		if v := identity.FlattenIdentity(bodyMap["identity"]); v != nil {
			model.Identity = identity.ToList(*v)
		}
	}
	model.Output = basetypes.NewStringValue(flattenOutput(responseBody, AsStringList(model.ResponseExportValues)))
	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}
