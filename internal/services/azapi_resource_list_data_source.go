package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/retry"
	"github.com/Azure/terraform-provider-azapi/internal/services/common"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework-timeouts/datasource/timeouts"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	tffwdocs "github.com/magodo/terraform-plugin-framework-docs"
)

type ResourceListDataSourceModel struct {
	ID                   types.String     `tfsdk:"id"`
	Type                 types.String     `tfsdk:"type"`
	ParentID             types.String     `tfsdk:"parent_id"`
	ResponseExportValues types.Dynamic    `tfsdk:"response_export_values"`
	Output               types.Dynamic    `tfsdk:"output"`
	Timeouts             timeouts.Value   `tfsdk:"timeouts"`
	Retry                retry.RetryValue `tfsdk:"retry"`
	Headers              types.Map        `tfsdk:"headers"`
	QueryParameters      types.Map        `tfsdk:"query_parameters"`
}

type ResourceListDataSource struct {
	ProviderData *clients.Client
}

var _ datasource.DataSource = &ResourceListDataSource{}
var _ datasource.DataSourceWithConfigure = &ResourceListDataSource{}
var _ tffwdocs.DataSourceWithRenderOption = &ResourceListDataSource{}

func (r *ResourceListDataSource) Configure(ctx context.Context, request datasource.ConfigureRequest, response *datasource.ConfigureResponse) {
	if v, ok := request.ProviderData.(*clients.Client); ok {
		r.ProviderData = v
	}
}

func (r *ResourceListDataSource) Metadata(ctx context.Context, request datasource.MetadataRequest, response *datasource.MetadataResponse) {
	response.TypeName = request.ProviderTypeName + "_resource_list"
}

func (r *ResourceListDataSource) Schema(ctx context.Context, request datasource.SchemaRequest, response *datasource.SchemaResponse) {
	response.Schema = schema.Schema{
		MarkdownDescription: "This data source allows you to list resources of a specific type under a given scope (e.g., subscription, resource group).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.ID(),
			},

			"type": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceType(),
				},
				MarkdownDescription: docstrings.Type(),
			},

			"parent_id": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					myvalidator.StringIsResourceID(),
				},
				MarkdownDescription: docstrings.ParentID(),
			},

			"response_export_values": schema.DynamicAttribute{
				Optional:            true,
				MarkdownDescription: docstrings.ResponseExportValuesForResourceList(),
			},

			"output": schema.DynamicAttribute{
				Computed:            true,
				MarkdownDescription: docstrings.Output("data.azapi_resource_list"),
			},

			"retry": retry.RetryDsSchema(ctx),

			"headers": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "A map of headers to include in the request.",
			},

			"query_parameters": schema.MapAttribute{
				ElementType: types.ListType{
					ElemType: types.StringType,
				},
				Optional:            true,
				MarkdownDescription: "A map of query parameters to include in the request.",
			},
		},

		Blocks: map[string]schema.Block{
			"timeouts": timeouts.Block(ctx),
		},
	}
}

func (r *ResourceListDataSource) Read(ctx context.Context, request datasource.ReadRequest, response *datasource.ReadResponse) {
	var model ResourceListDataSourceModel
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

	id, err := parse.NewResourceIDSkipScopeValidation("", model.ParentID.ValueString(), model.Type.ValueString())
	if err != nil {
		response.Diagnostics.AddError("Invalid configuration", err.Error())
		return
	}

	ctx = tflog.SetField(ctx, "resource_id", id.ID())

	listUrl := strings.TrimSuffix(id.AzureResourceId, "/")

	client := r.ProviderData.ResourceClient
	requestOptions := clients.RequestOptions{
		Headers:         common.AsMapOfString(model.Headers),
		QueryParameters: clients.NewQueryParameters(common.AsMapOfLists(model.QueryParameters)),
	}
	requestOptions.RetryOptions, requestOptions.LastRetryError = clients.NewRetryOptions(model.Retry)
	responseBody, err := client.List(ctx, listUrl, id.ApiVersion, requestOptions)
	if err != nil {
		response.Diagnostics.AddError("Failed to list resources", fmt.Sprintf("Failed to list resources, url: %s, error: %s", listUrl, err.Error()))
		return
	}

	model.ID = basetypes.NewStringValue(listUrl)
	var defaultOutput interface{}
	if !r.ProviderData.Features.DisableDefaultOutput {
		defaultOutput = responseBody
		defaultOutput = utils.RemoveFields(defaultOutput, volatileFieldList())
	}
	output, err := buildOutputFromBody(responseBody, model.ResponseExportValues, defaultOutput)
	if err != nil {
		response.Diagnostics.AddError("Failed to build output", err.Error())
		return
	}
	model.Output = output

	response.Diagnostics.Append(response.State.Set(ctx, &model)...)
}

func (r *ResourceListDataSource) RenderOption() tffwdocs.DataSourceRenderOption {
	return tffwdocs.DataSourceRenderOption{
		Examples: []tffwdocs.Example{
			{
				HCL: `
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

data "azapi_resource_list" "listBySubscription" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
  response_export_values = {
    "values" = "value[].{name: name, publicNetworkAccess: properties.publicNetworkAccess}"
    "names"  = "value[].name"
  }
}

data "azapi_resource_list" "listByResourceGroup" {
  type      = "Microsoft.Automation/automationAccounts@2021-06-22"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1"
  response_export_values = {
    "names" = "value[].name"
  }
}

data "azapi_resource_list" "listSubnetsByVnet" {
  type                   = "Microsoft.Network/virtualNetworks/subnets@2021-02-01"
  parent_id              = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/rg1/providers/Microsoft.Network/virtualNetworks/vnet1"
  response_export_values = ["*"]
}
`,
			},
		},
	}
}
