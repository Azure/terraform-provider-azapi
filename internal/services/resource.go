package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	aztypes "github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func schemaValidate(config *AzapiResourceModel) error {
	if config == nil {
		return nil
	}

	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(config.Type.ValueString())
	if err != nil {
		return fmt.Errorf(`the argument "type" is invalid: %s`, err.Error())
	}
	resourceDef, _ := azure.GetResourceDefinition(azureResourceType, apiVersion)

	log.Printf("[INFO] prepare validation for resource type: %s, api-version: %s", azureResourceType, apiVersion)
	versions := azure.GetApiVersions(azureResourceType)
	if len(versions) == 0 {
		return schemaValidationError(fmt.Sprintf("the argument \"type\" is invalid.\n resource type %s can't be found.\n", azureResourceType))
	}
	isVersionValid := false
	for _, version := range versions {
		if version == apiVersion {
			isVersionValid = true
			break
		}
	}
	if !isVersionValid {
		return schemaValidationError(fmt.Sprintf("the argument \"type\"'s api-version is invalid.\n The supported versions are [%s].\n", strings.Join(versions, ", ")))
	}

	if resourceDef == nil {
		return nil
	}

	var bodyToValidate attr.Value
	if !config.Body.IsNull() && !config.Body.IsUnknown() && !config.Body.IsNull() && !config.Body.IsUnderlyingValueUnknown() {
		if v, ok := config.Body.UnderlyingValue().(types.Object); ok {
			attributes := v.Attributes()
			attributeTypes := v.AttributeTypes(context.Background())

			attributes["name"] = config.Name
			attributeTypes["name"] = types.StringType

			if !config.Location.IsNull() {
				attributes["location"] = config.Location
				attributeTypes["location"] = types.StringType
			}

			if !config.Tags.IsNull() {
				attributes["tags"] = config.Tags
				attributeTypes["tags"] = types.MapType{ElemType: types.StringType}
			}

			if !config.Identity.IsNull() {
				identityAttributeTypes := map[string]attr.Type{
					"type": types.StringType,
				}
				identityModel := identity.FromList(config.Identity)
				if len(identityModel.IdentityIDs.Elements()) != 0 {
					identityAttributeTypes["userAssignedIdentities"] = types.MapType{ElemType: types.DynamicType}
					elements := make(map[string]attr.Value)
					identityIds := identityModel.IdentityIDs.Elements()
					for _, identityId := range identityIds {
						elements[identityId.(types.String).ValueString()] = types.DynamicNull()
					}
					attributes["identity"] = types.ObjectValueMust(identityAttributeTypes, map[string]attr.Value{
						"type":                   identityModel.Type,
						"userAssignedIdentities": types.MapValueMust(types.DynamicType, elements),
					})
				} else {
					attributes["identity"] = types.ObjectValueMust(identityAttributeTypes, map[string]attr.Value{
						"type": identityModel.Type,
					})
				}
				attributeTypes["identity"] = types.ObjectType{AttrTypes: identityAttributeTypes}
			}

			bodyToValidate = types.ObjectValueMust(attributeTypes, attributes)
		}
	} else {
		bodyToValidate = config.Body
	}

	validateErrors := (*resourceDef).Validate(bodyToValidate, "")

	errors := make([]error, 0)
	// skip error that location is not in the body, because user might use the default location feature
	// same for tags
	for _, err := range validateErrors {
		if strings.Contains(err.Error(), "`location` is required") || strings.Contains(err.Error(), "`tags` is required") {
			continue
		}
		errors = append(errors, err)
	}

	if len(errors) != 0 {
		errorMsg := "the argument \"body\" is invalid:\n"
		for _, err := range errors {
			errorMsg += fmt.Sprintf("%s\n", err.Error())
		}
		return schemaValidationError(errorMsg)
	}

	return nil
}

func schemaValidationError(detail string) error {
	return fmt.Errorf("embedded schema validation failed: %s You can try to update `azapi` provider to "+
		"the latest version or disable the validation using the feature flag `schema_validation_enabled = false` "+
		"within the resource block", detail)
}

func canResourceHaveProperty(resourceDef *aztypes.ResourceType, property string) bool {
	if resourceDef == nil || resourceDef.Body == nil || resourceDef.Body.Type == nil {
		return false
	}
	objectType, ok := (*resourceDef.Body.Type).(*aztypes.ObjectType)
	if !ok {
		return false
	}
	if prop, ok := objectType.Properties[property]; ok {
		if !prop.IsReadOnly() {
			return true
		}
	}
	return false
}

func flattenBody(responseBody interface{}, resourceDef *aztypes.ResourceType) (types.Dynamic, error) {
	body := utils.NormalizeObject(responseBody)

	if resourceDef != nil {
		writeOnlyBody := (*resourceDef).GetWriteOnly(body)
		if bodyMap, ok := writeOnlyBody.(map[string]interface{}); ok {
			delete(bodyMap, "location")
			delete(bodyMap, "tags")
			delete(bodyMap, "name")
			delete(bodyMap, "identity")
			writeOnlyBody = bodyMap
		}
		body = writeOnlyBody
	}

	data, err := json.Marshal(body)
	if err != nil {
		return types.DynamicNull(), err
	}
	return dynamic.FromJSONImplied(data)
}

func flattenOutput(responseBody interface{}, paths []string) attr.Value {
	for _, path := range paths {
		if path == "*" {
			if v, ok := responseBody.(string); ok {
				return basetypes.NewStringValue(v)
			}
			data, err := json.Marshal(responseBody)
			if err != nil {
				return nil
			}
			out, err := dynamic.FromJSONImplied(data)
			if err != nil {
				return nil
			}
			return out
		}
	}

	var output interface{}
	output = make(map[string]interface{})
	for _, path := range paths {
		part := utils.ExtractObject(responseBody, path)
		if part == nil {
			continue
		}
		output = utils.MergeObject(output, part)
	}
	data, err := json.Marshal(output)
	if err != nil {
		return nil
	}
	out, err := dynamic.FromJSONImplied(data)
	if err != nil {
		return nil
	}
	return out
}

func flattenOutputJMES(responseBody interface{}, paths map[string]string) attr.Value {
	var output interface{}
	output = make(map[string]interface{})
	for pathKey, path := range paths {
		part := utils.ExtractObjectJMES(responseBody, pathKey, path)
		if part == nil {
			continue
		}
		output = utils.MergeObject(output, part)
	}
	data, err := json.Marshal(output)
	if err != nil {
		return nil
	}
	out, err := dynamic.FromJSONImplied(data)
	if err != nil {
		return nil
	}
	return out
}

func AsStringList(input types.List) []string {
	var result []string
	diags := input.ElementsAs(context.Background(), &result, false)
	if diags.HasError() {
		tflog.Warn(context.Background(), fmt.Sprintf("failed to convert list to string list: %s", diags))
	}
	return result
}

func AsMapOfString(input types.Map) map[string]string {
	result := make(map[string]string)
	diags := input.ElementsAs(context.Background(), &result, false)
	if diags.HasError() {
		tflog.Warn(context.Background(), fmt.Sprintf("failed to convert input to map of strings: %s", diags))
	}
	return result
}

func AsMapOfLists(input types.Map) map[string][]string {
	result := make(map[string][]string)
	diags := input.ElementsAs(context.Background(), &result, false)
	if diags.HasError() {
		tflog.Warn(context.Background(), fmt.Sprintf("failed to convert input to map of lists: %s", diags))
	}
	return result
}

func unmarshalBody(input types.Dynamic, out interface{}) error {
	if input.IsNull() || input.IsUnknown() || input.IsUnderlyingValueUnknown() {
		return nil
	}
	data, err := dynamic.ToJSON(input)
	if err != nil {
		return fmt.Errorf(`invalid dynamic value: value: %s, err: %+v`, input.String(), err)
	}
	if err = json.Unmarshal(data, &out); err != nil {
		return fmt.Errorf(`unmarshaling failed: value: %s, err: %+v`, string(data), err)
	}
	return nil
}
