package preflight

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure"
	aztypes "github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RequestBodyModel struct {
	Provider  string                   `json:"provider"`
	Type      string                   `json:"type"`
	Location  string                   `json:"location"`
	Scope     string                   `json:"scope"`
	Resources []map[string]interface{} `json:"resources"`
}

func ParentIdPlaceholder(resourceDef *aztypes.ResourceType, subscriptionId string) string {
	// since the parentID is faked, there should exist only one scope type
	if resourceDef == nil || len(resourceDef.ScopeTypes) != 1 {
		return ""
	}

	parentId := ""
	switch resourceDef.ScopeTypes[0] {
	case aztypes.Tenant:
		parentId = "/"
	case aztypes.ManagementGroup:
		parentId = "/providers/Microsoft.Management/managementGroups/azapifakemg"
	case aztypes.Subscription:
		parentId = fmt.Sprintf("/subscriptions/%s", subscriptionId)
	case aztypes.ResourceGroup:
		parentId = fmt.Sprintf("/subscriptions/%s/resourceGroups/azapifakerg", subscriptionId)
	default:
	}
	return parentId
}

func NamePlaceholder() string {
	return "placeholder"
}

func IsSupported(parentId string, resourceType string) bool {
	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return false
	}

	if !utils.IsTopLevelResourceType(azureResourceType) {
		return false
	}

	// if the parentID is specified, it should be a resource group, subscription, tenant or management group
	if parentId != "" {
		parentResourceType := utils.GetResourceType(parentId)
		return strings.EqualFold(arm.ResourceGroupResourceType.String(), parentResourceType) ||
			strings.EqualFold(arm.SubscriptionResourceType.String(), parentResourceType) ||
			strings.EqualFold(arm.TenantResourceType.String(), parentResourceType) ||
			strings.EqualFold("Microsoft.Management/managementGroups", parentResourceType)
	}

	// if the parentID is not specified, the resource type should be able to deploy only at the tenant, management group, subscription or resource group level

	resourceDef, err := azure.GetResourceDefinition(azureResourceType, apiVersion)
	if err != nil || resourceDef == nil || len(resourceDef.ScopeTypes) != 1 {
		return false
	}

	deployedScope := resourceDef.ScopeTypes[0]

	return deployedScope == aztypes.Tenant || deployedScope == aztypes.ManagementGroup ||
		deployedScope == aztypes.Subscription || deployedScope == aztypes.ResourceGroup
}

func Validate(ctx context.Context, client *clients.ResourceClient, resourceType string, parentId string, name string, location string, body types.Dynamic) error {
	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return err
	}

	payload := RequestBodyModel{}
	payload.Provider, payload.Type, _ = strings.Cut(azureResourceType, "/")
	payload.Scope = parentId
	if location != "" {
		payload.Location = location
	}

	resource := make(map[string]interface{})
	err = unmarshalPreflightBody(body, &resource)
	if err != nil {
		return err
	}

	resource["name"] = name
	resource["apiVersion"] = apiVersion

	payload.Resources = []map[string]interface{}{resource}

	_, err = client.Action(ctx, "/providers/Microsoft.Resources", "validateResources", "2020-10-01", "POST", payload, clients.DefaultRequestOptions())
	if err != nil {
		return err
	}

	return nil
}

func unmarshalPreflightBody(input types.Dynamic, out interface{}) error {
	if input.IsNull() || input.IsUnknown() || input.IsUnderlyingValueUnknown() {
		return nil
	}

	const unknownPlaceholder = "[length('foo')]"

	data, err := dynamic.ToJSONWithUnknownValueHandler(input, func(value attr.Value) ([]byte, error) {
		return json.Marshal(unknownPlaceholder)
	})

	res := map[string]interface{}{}
	if err = json.Unmarshal(data, &res); err != nil {
		return fmt.Errorf(`unmarshaling failed: value: %s, err: %+v`, string(data), err)
	}

	// make sure that there's no unknown value outside the properties bag
	for k, v := range res {
		if k == "properties" {
			continue
		}
		if searchForValue(v, unknownPlaceholder) {
			return fmt.Errorf("unknown value found outside the properties bag")
		}
	}

	out = res
	return nil
}

func searchForValue(input interface{}, target string) bool {
	switch v := input.(type) {
	case map[string]interface{}:
		for _, value := range v {
			if searchForValue(value, target) {
				return true
			}
		}
	case []interface{}:
		for _, value := range v {
			if searchForValue(value, target) {
				return true
			}
		}
	case string:
		if v == target {
			return true
		}
	}

	return false
}
