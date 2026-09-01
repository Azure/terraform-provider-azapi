package services

import (
	"context"
	"fmt"
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
	"github.com/jmespath/go-jmespath"
)

const readActionPOSTVerb = "POST"

type ReadOverrideModel struct {
	Method       types.String `tfsdk:"method"`
	Action       types.String `tfsdk:"action"`
	ResponsePath types.String `tfsdk:"response_path"`
}

func readOverrideAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"method":        types.StringType,
		"action":        types.StringType,
		"response_path": types.StringType,
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
					stringvalidator.OneOf(readActionPOSTVerb),
				},
			},
			"action": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "The name of the action appended to the resource ID, for example `list`.",
				Validators: []validator.String{
					myvalidator.StringIsNotEmpty(),
				},
			},
			"response_path": schema.StringAttribute{
				Optional:            true,
				MarkdownDescription: "An optional [JMESPath](https://jmespath.org/) expression applied to the action response before it is merged into `body`.",
			},
		},
		MarkdownDescription: "Overrides the request used to read resource values when a plain `GET` does not return them. When unset, the provider uses the response from the default `GET` request.",
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

func resolveReadResponse(ctx context.Context, client *clients.ResourceClient, azureResourceId, apiVersion string, getResponse interface{}, readOverride *ReadOverrideModel, requestOptions clients.RequestOptions) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	if readOverride == nil {
		return getResponse, diags
	}

	action := strings.TrimSpace(readOverride.Action.ValueString())
	method := strings.TrimSpace(readOverride.Method.ValueString())

	response, err := client.Action(ctx, azureResourceId, action, apiVersion, method, nil, requestOptions)
	if err != nil {
		diags.AddError("Failed to read resource via action", fmt.Sprintf("%s %q on %q failed: %s", method, action, azureResourceId, err.Error()))
		return getResponse, diags
	}
	return applyResponsePath(response, readOverride.ResponsePath.ValueString(), &diags), diags
}

func applyResponsePath(response interface{}, path string, diags *diag.Diagnostics) interface{} {
	path = strings.TrimSpace(path)
	if path == "" {
		return response
	}
	result, err := jmespath.Search(path, response)
	if err != nil {
		diags.AddError("Invalid read_override response_path", fmt.Sprintf("failed to evaluate JMESPath %q against the read override response: %s", path, err.Error()))
		return response
	}
	return result
}
