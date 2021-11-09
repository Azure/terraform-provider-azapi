package identity

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/parse"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/validate"
)

type IdentityType string

const (
	None                       IdentityType = "None"
	SystemAssigned             IdentityType = "SystemAssigned"
	UserAssigned               IdentityType = "UserAssigned"
	SystemAssignedUserAssigned IdentityType = "SystemAssigned, UserAssigned"
)

func SchemaIdentity() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Required: true,
					ValidateFunc: validation.StringInSlice([]string{
						string(None),
						string(UserAssigned),
						string(SystemAssigned),
						string(SystemAssignedUserAssigned),
					}, false),
				},

				"identity_ids": {
					Type:     schema.TypeList,
					Optional: true,
					Elem: &schema.Schema{
						Type:         schema.TypeString,
						ValidateFunc: validate.UserAssignedIdentityID,
					},
				},

				"principal_id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func SchemaIdentityDataSource() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Computed: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"type": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"identity_ids": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},

				"principal_id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"tenant_id": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

func ExpandIdentity(input []interface{}) (interface{}, error) {
	body := make(map[string]interface{}, 0)
	if len(input) == 0 || input[0] == nil {
		return body, nil
	}

	v := input[0].(map[string]interface{})

	config := map[string]interface{}{}
	identityType := IdentityType(v["type"].(string))
	config["type"] = identityType
	identityIds := v["identity_ids"].([]interface{})
	userAssignedIdentities := make(map[string]interface{}, len(identityIds))
	if len(identityIds) != 0 {
		if identityType != UserAssigned && identityType != SystemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}
		for _, id := range identityIds {
			userAssignedIdentities[id.(string)] = make(map[string]interface{})
		}
		config["userAssignedIdentities"] = userAssignedIdentities
	}
	body["identity"] = config
	return body, nil
}

func FlattenIdentity(body interface{}) []interface{} {
	if body != nil {
		if bodyMap, ok := body.(map[string]interface{}); ok && bodyMap["identity"] != nil {
			if identityMap, ok := bodyMap["identity"].(map[string]interface{}); ok {
				identityIds := make([]string, 0)
				if identityMap["userAssignedIdentities"] != nil {
					userAssignedIdentities := identityMap["userAssignedIdentities"].(map[string]interface{})
					for key := range userAssignedIdentities {
						identityId, err := parse.UserAssignedIdentitiesID(key)
						if err == nil {
							identityIds = append(identityIds, identityId.ID())
						}
					}
				}

				identityType := identityMap["type"].(string)
				switch {
				case strings.Contains(identityType, ","):
					identityType = string(SystemAssignedUserAssigned)
				case strings.EqualFold(identityType, string(UserAssigned)):
					identityType = string(UserAssigned)
				case strings.EqualFold(identityType, string(SystemAssigned)):
					identityType = string(SystemAssigned)
				default:
					identityType = string(None)
				}

				return []interface{}{
					map[string]interface{}{
						"type":         identityType,
						"identity_ids": identityIds,
						"principal_id": identityMap["principalId"],
						"tenant_id":    identityMap["tenantId"],
					},
				}
			}
		}
	}
	return nil
}
