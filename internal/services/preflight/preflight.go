package preflight

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/terraform-provider-azapi/internal/azure/identity"
	aztypes "github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/dynamic"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
)

type RequestBodyModel struct {
	Provider  string                   `json:"provider"`
	Type      string                   `json:"type"`
	Location  string                   `json:"location"`
	Scope     string                   `json:"scope"`
	Resources []map[string]interface{} `json:"resources"`
}

// ParentIdPlaceholder generates a placeholder for the parentID based on the resource definition and subscription ID
func ParentIdPlaceholder(resourceDef *aztypes.ResourceType, subscriptionId string) (string, error) {
	// since the parentID is faked, there should exist only one scope type
	if resourceDef == nil || len(resourceDef.ScopeTypes) != 1 {
		return "", fmt.Errorf("failed to generate parentID placeholder because the resource definition is invalid")
	}

	scopeId := ""
	switch resourceDef.ScopeTypes[0] {
	case aztypes.Tenant:
		scopeId = "/"
	case aztypes.ManagementGroup:
		scopeId = "/providers/Microsoft.Management/managementGroups/" + NamePlaceholder()
	case aztypes.Subscription:
		scopeId = fmt.Sprintf("/subscriptions/%s", subscriptionId)
	case aztypes.ResourceGroup:
		scopeId = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionId, NamePlaceholder())
	default:
		return "", fmt.Errorf("failed to generate parentID placeholder because the scope type is not supported")
	}

	// for top-level resources, the parent ID is the same as the scope ID
	resourceType := strings.Split(resourceDef.Name, "@")[0]
	if utils.IsTopLevelResourceType(resourceType) {
		return scopeId, nil
	}

	parts := strings.Split(resourceType, "/")

	parentId := fmt.Sprintf("%s/providers/%s", scopeId, parts[0])
	for i := 1; i < len(parts)-1; i++ {
		parentId += fmt.Sprintf("/%s/%s", parts[i], NamePlaceholder())
	}

	return parentId, nil
}

// NamePlaceholder generates a random name placeholder
func NamePlaceholder() string {
	return acctest.RandStringFromCharSet(8, acctest.CharSetAlpha)
}

// Validate validates the resource using the preflight API
func Validate(ctx context.Context, client *clients.ResourceClient, resourceType string, parentId string, name string, location string, body types.Dynamic, identity types.List) error {
	azureResourceType, apiVersion, err := utils.GetAzureResourceTypeApiVersion(resourceType)
	if err != nil {
		return err
	}

	payload := RequestBodyModel{}
	payload.Provider, payload.Type, _ = strings.Cut(azureResourceType, "/")
	id, err := parse.NewResourceIDSkipScopeValidation(name, parentId, resourceType)
	if err != nil {
		return err
	}
	scopeId, err := ScopeID(id.ID())
	if err != nil {
		return err
	}
	payload.Scope = scopeId
	if location != "" {
		payload.Location = location
	}

	resource := make(map[string]interface{})
	err = unmarshalPreflightBody(body, identity, &resource)
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Skipping preflight validation for resource %s because the body is invalid: %v", resourceType, err))
		return nil
	}

	resource["name"] = name
	resource["apiVersion"] = apiVersion

	payload.Resources = []map[string]interface{}{resource}

	_, err = client.Action(ctx, "/providers/Microsoft.Resources", "validateResources", "2020-10-01", "POST", payload, clients.DefaultRequestOptions())
	return err
}

func ScopeID(resourceId string) (string, error) {
	armId, err := arm.ParseResourceID(resourceId)
	if err != nil {
		return "", err
	}
	if armId.Parent == nil {
		return "/", nil
	}
	scopeId := armId.Parent
	for scopeId.Parent != nil &&
		!strings.EqualFold(scopeId.ResourceType.String(), arm.SubscriptionResourceType.String()) &&
		!strings.EqualFold(scopeId.ResourceType.String(), arm.ResourceGroupResourceType.String()) &&
		!strings.EqualFold(scopeId.ResourceType.String(), arm.TenantResourceType.String()) &&
		!strings.EqualFold(scopeId.ResourceType.String(), "Microsoft.Management/managementGroups") {
		scopeId = scopeId.Parent
	}

	if scopeId == nil || scopeId.ResourceType.String() == arm.TenantResourceType.String() {
		return "/", nil
	}

	return scopeId.String(), nil
}

func unmarshalPreflightBody(input types.Dynamic, identityList types.List, out *map[string]interface{}) error {
	if input.IsNull() || input.IsUnknown() || input.IsUnderlyingValueUnknown() {
		return fmt.Errorf("input is null or unknown")
	}

	const unknownPlaceholder = "[length('foo')]"

	data, err := dynamic.ToJSONWithUnknownValueHandler(input, func(value attr.Value) ([]byte, error) {
		return json.Marshal(unknownPlaceholder)
	})

	if err != nil {
		return fmt.Errorf("marshaling failed: %v", err)
	}

	if err = json.Unmarshal(data, &out); err != nil {
		return fmt.Errorf(`unmarshaling failed: value: %s, err: %+v`, string(data), err)
	}

	if out == nil {
		out = &map[string]interface{}{}
	}
	// make sure that there's no unknown value outside the properties bag
	for k, v := range *out {
		if k == "properties" {
			continue
		}
		if searchForValue(v, unknownPlaceholder) {
			return fmt.Errorf("unknown value found outside the properties bag")
		}
	}

	if (*out)["identity"] == nil && !identityList.IsNull() && !identityList.IsUnknown() {
		identityModel := identity.FromList(identityList)
		expandedIdentity, err := identity.ExpandIdentity(identityModel)
		if err != nil {
			return err
		}
		(*out)["identity"] = expandedIdentity
	}
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
