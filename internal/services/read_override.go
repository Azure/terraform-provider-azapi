package services

import (
	"context"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/myvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

const readOverridePOSTVerb = "POST"

type ReadOverrideModel struct {
	Method types.String `tfsdk:"method"`
	Action types.String `tfsdk:"action"`
}

func readOverrideAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"method": types.StringType,
		"action": types.StringType,
	}
}

func readOverrideSchema() schema.Attribute {
	return schema.SingleNestedAttribute{
		Optional: true,
		Attributes: map[string]schema.Attribute{
			"method": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The HTTP method used to read the resource. The only supported value is `POST`.",
				Validators: []validator.String{
					stringvalidator.OneOf(readOverridePOSTVerb),
				},
			},
			"action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the action appended to the resource ID, for example `list`.",
				Validators: []validator.String{
					myvalidator.StringIsNotEmpty(),
				},
			},
		},
		MarkdownDescription: "Overrides the default `GET` request used to read the resource. When configured, the provider sends the specified action request instead of `GET` and uses its response for all read processing, including refreshing `body` and `output`. When omitted, the provider reads the resource with `GET`.",
	}
}

func readOverrideFromObject(ctx context.Context, value types.Object) (*ReadOverrideModel, diag.Diagnostics) {
	if value.IsNull() || value.IsUnknown() {
		return nil, nil
	}

	var model ReadOverrideModel
	diags := value.As(ctx, &model, basetypes.ObjectAsOptions{})
	return &model, diags
}

func readResource(ctx context.Context, client *clients.ResourceClient, azureResourceId, apiVersion string, readOverride *ReadOverrideModel, requestOptions clients.RequestOptions) (interface{}, error) {
	if readOverride == nil {
		return client.Get(ctx, azureResourceId, apiVersion, requestOptions)
	}

	action := strings.TrimSpace(readOverride.Action.ValueString())
	method := strings.TrimSpace(readOverride.Method.ValueString())
	return client.Action(ctx, azureResourceId, action, apiVersion, method, nil, requestOptions)
}
