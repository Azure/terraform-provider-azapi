package services

import (
	"encoding/json"
	"errors"
	"os"
	"time"

	"github.com/Azure/terraform-provider-azapi/internal/docstrings"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/myplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func Retry404MaxElapsedTime() time.Duration {
	if v := os.Getenv("AZAPI_RETRY_404_MAX_ELAPSED_TIME"); v != "" {
		timeout, err := time.ParseDuration(v)
		if err != nil {
			return timeout
		}
	}
	return 2 * time.Minute
}

func CommonAttributeResponseExportValues() schema.DynamicAttribute {
	return schema.DynamicAttribute{
		Optional: true,
		PlanModifiers: []planmodifier.Dynamic{
			myplanmodifier.DynamicUseStateWhen(dynamic.SemanticallyEqual),
		},
		MarkdownDescription: docstrings.ResponseExportValues(),
	}
}

func buildOutputFromBody(responseBody interface{}, modelResponseExportValues types.Dynamic, defaultResult interface{}) (types.Dynamic, error) {
	if modelResponseExportValues.IsNull() {
		if defaultResult == nil {
			return types.DynamicValue(types.ObjectValueMust(map[string]attr.Type{}, map[string]attr.Value{})), nil
		}
		data, err := json.Marshal(defaultResult)
		if err != nil {
			return types.DynamicNull(), err
		}
		return dynamic.FromJSONImplied(data)
	}

	data, err := dynamic.ToJSON(modelResponseExportValues)
	if err != nil {
		return types.DynamicNull(), err
	}

	switch modelResponseExportValues.UnderlyingValue().(type) {
	case types.List, types.Tuple, types.Set:
		var responseExportValues []string
		if err = json.Unmarshal(data, &responseExportValues); err != nil {
			return types.DynamicNull(), err
		}

		return types.DynamicValue(flattenOutput(responseBody, responseExportValues)), nil
	case types.Map, types.Object:
		var responseExportValues map[string]string
		if err = json.Unmarshal(data, &responseExportValues); err != nil {
			return types.DynamicNull(), err
		}

		return types.DynamicValue(flattenOutputJMES(responseBody, responseExportValues)), nil
	default:
		return types.DynamicNull(), errors.New("unsupported type for response_export_values, must be a list or map")
	}
}
