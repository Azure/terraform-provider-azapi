package services

import (
	"fmt"
	"slices"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/services/common"
	"github.com/Azure/terraform-provider-azapi/internal/services/customization"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// validateDataPlaneResourceWritable returns an error if the resource type has a registered
// customization but exposes no write operations (nil CreateFunc). Such types are read-only by
// design and should only be used with the data source, not the managed resource.
func validateDataPlaneResourceWritable(resourceType string) error {
	customizedResource := customization.GetCustomization(resourceType)
	if customizedResource == nil {
		return nil
	}
	if (*customizedResource).CreateFunc() == nil {
		typeBase := strings.Split(resourceType, "@")[0]
		return fmt.Errorf(`resource type %q does not support create/update/delete; use data.azapi_data_plane_resource to read it instead`, typeBase)
	}
	return nil
}

func validateDataPlaneResourceAddress(config *DataPlaneResourceModel) error {
	if config == nil || config.Type.IsNull() || config.Type.IsUnknown() {
		return nil
	}

	customizedResource := customization.GetCustomization(config.Type.ValueString())
	hasCreateResult := false
	if customizedResource != nil {
		if v, ok := (*customizedResource).(customization.DataPlaneResourceWithCreateResult); ok && v.CreateResultFunc() != nil {
			hasCreateResult = true
		}
	}

	placeholderKeys, err := parse.DataPlaneResourcePlaceholderKeys(config.Type.ValueString())
	if err != nil {
		return err
	}

	requiresName := slices.Contains(placeholderKeys, "name")
	requiredIdentifiers := make([]string, 0)
	for _, key := range placeholderKeys {
		if key != "name" {
			requiredIdentifiers = append(requiredIdentifiers, key)
		}
	}

	if requiresName && hasCreateResult {
		if !config.Name.IsNull() && !config.Name.IsUnknown() && strings.TrimSpace(config.Name.ValueString()) != "" {
			return fmt.Errorf(`the argument "name" should not be set for resource type %q because the service generates the identifier`, strings.Split(config.Type.ValueString(), "@")[0])
		}
	} else if requiresName {
		if config.Name.IsUnknown() {
			return nil
		}
		if config.Name.IsNull() || strings.TrimSpace(config.Name.ValueString()) == "" {
			return fmt.Errorf(`the argument "name" must be set for resource type %q`, strings.Split(config.Type.ValueString(), "@")[0])
		}
	} else if !config.Name.IsNull() && !config.Name.IsUnknown() && strings.TrimSpace(config.Name.ValueString()) != "" {
		return fmt.Errorf(`the argument "name" is not used for resource type %q`, strings.Split(config.Type.ValueString(), "@")[0])
	}

	if len(requiredIdentifiers) == 0 {
		return nil
	}
	if config.Identifiers.IsUnknown() {
		return nil
	}
	// ValidateConfig is called once per block before for_each is expanded.
	// At that point each.value.pk / each.value.rk are not yet resolved and
	// appear as unknown attribute values inside an otherwise-known map.
	// Skip validation now; the plan walk re-invokes ValidateConfig with
	// concrete values for each expanded instance.
	for _, v := range config.Identifiers.Elements() {
		if v.IsUnknown() {
			return nil
		}
	}

	identifiers := common.AsMapOfString(config.Identifiers)
	missing := make([]string, 0)
	for _, key := range requiredIdentifiers {
		if strings.TrimSpace(identifiers[key]) == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) != 0 {
		return fmt.Errorf(`the argument "identifiers" must include non-empty values for [%s] for resource type %q`, strings.Join(missing, ", "), strings.Split(config.Type.ValueString(), "@")[0])
	}
	return nil
}

func stringMapToTypesMap(input map[string]string) types.Map {
	if len(input) == 0 {
		return types.MapNull(types.StringType)
	}
	values := make(map[string]attr.Value, len(input))
	for key, value := range input {
		values[key] = types.StringValue(value)
	}
	return types.MapValueMust(types.StringType, values)
}
