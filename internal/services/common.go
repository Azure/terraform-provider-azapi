package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func CommonAttributeResponseExportValues() schema.DynamicAttribute {
	return schema.DynamicAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.Dynamic{
			myplanmodifier.DynamicUseStateWhen(dynamic.SemanticallyEqual),
		},
		MarkdownDescription: docstrings.ResponseExportValues(),
	}
}

func buildOutputFromBody(ctx context.Context, responseBody interface{}, modelResponseExportValues types.Dynamic) (types.Dynamic, error) {
	if modelResponseExportValues.IsNull() || modelResponseExportValues.IsUnknown() {
		return types.DynamicNull(), nil
	}

	data, err := dynamic.ToJSON(modelResponseExportValues)
	if err != nil {
		return types.DynamicNull(), err
	}

	tflog.Debug(ctx, fmt.Sprintf("data: %+v\n", data))
	tflog.Debug(ctx, fmt.Sprintf("modelResponseExportValues is: %+v\n", modelResponseExportValues))
	tflog.Debug(ctx, fmt.Sprintf("modelResponseExportValues underlying value is: %+v\n", modelResponseExportValues.UnderlyingValue()))
	tflog.Debug(ctx, fmt.Sprintf("modelResponseExportValues underlying value type is: %+v\n", modelResponseExportValues.UnderlyingValue().Type(ctx)))
	switch modelResponseExportValues.UnderlyingValue().(type) {
	case types.List, types.Tuple, types.Set:
		tflog.Debug(ctx, fmt.Sprintf("it is a list\n"))
		var responseExportValues []string
		if err = json.Unmarshal(data, &responseExportValues); err != nil {
			tflog.Debug(ctx, fmt.Sprintf("list error: %+v\n", err))
			return types.DynamicNull(), err
		}

		tflog.Debug(ctx, fmt.Sprintf("responseExportValues is: %+v\n", responseExportValues))
		return types.DynamicValue(flattenOutput(responseBody, responseExportValues)), nil
	case types.Map, types.Object:
		tflog.Debug(ctx, fmt.Sprintf("it is a map\n"))
		var responseExportValues map[string]string
		if err = json.Unmarshal(data, &responseExportValues); err != nil {
			fmt.Printf("map error: %+v\n", err)
			return types.DynamicNull(), err
		}

		tflog.Debug(ctx, fmt.Sprintf("responseExportValues is: %+v\n", responseExportValues))

		return types.DynamicValue(flattenOutputJMES(responseBody, responseExportValues)), nil
	default:
		return types.DynamicNull(), errors.New("unsupported type for response_export_values, must be a list or map")
	}
}
