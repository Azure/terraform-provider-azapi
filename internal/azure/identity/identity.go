package identity

import (
	"context"
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type IdentityType string

const (
	None                       IdentityType = "None"
	SystemAssigned             IdentityType = "SystemAssigned"
	UserAssigned               IdentityType = "UserAssigned"
	SystemAssignedUserAssigned IdentityType = "SystemAssigned, UserAssigned"
)

type Model struct {
	Type        types.String `tfsdk:"type"`
	IdentityIDs types.List   `tfsdk:"identity_ids"`
	PrincipalID types.String `tfsdk:"principal_id"`
	TenantID    types.String `tfsdk:"tenant_id"`
}

func (m Model) ModelType() attr.Type {
	return types.ObjectType{AttrTypes: m.AttrType()}
}

func (m Model) AttrType() map[string]attr.Type {
	return map[string]attr.Type{
		"type":         types.StringType,
		"identity_ids": types.ListType{ElemType: types.StringType},
		"principal_id": types.StringType,
		"tenant_id":    types.StringType,
	}
}

func (m Model) Value() attr.Value {
	var identityIDs attr.Value
	if m.IdentityIDs.IsNull() {
		identityIDs = basetypes.NewListNull(types.StringType)
	} else {
		identityIDs = basetypes.NewListValueMust(types.StringType, m.IdentityIDs.Elements())
	}

	return types.ObjectValueMust(m.AttrType(), map[string]attr.Value{
		"type":         types.StringValue(m.Type.ValueString()),
		"identity_ids": identityIDs,
		"principal_id": types.StringValue(m.PrincipalID.ValueString()),
		"tenant_id":    types.StringValue(m.TenantID.ValueString()),
	})
}

func (m Model) ObjectType() attr.Type {
	return types.ObjectType{AttrTypes: m.AttrType()}
}

func ExpandIdentity(input Model) (interface{}, error) {
	config := map[string]interface{}{}
	identityType := IdentityType(input.Type.ValueString())
	config["type"] = identityType
	identityIds := input.IdentityIDs.Elements()
	userAssignedIdentities := make(map[string]interface{}, len(identityIds))
	if len(identityIds) != 0 {
		if identityType != UserAssigned && identityType != SystemAssignedUserAssigned {
			return nil, fmt.Errorf("`identity_ids` can only be specified when `type` includes `UserAssigned`")
		}
		for _, id := range identityIds {
			userAssignedIdentities[id.(basetypes.StringValue).ValueString()] = make(map[string]interface{})
		}
		config["userAssignedIdentities"] = userAssignedIdentities
	}
	return config, nil
}

func FlattenIdentity(identity interface{}) *Model {
	if identity == nil {
		return nil
	}
	if identityMap, ok := identity.(map[string]interface{}); ok {
		identityIds := make([]attr.Value, 0)
		if identityMap["userAssignedIdentities"] != nil {
			userAssignedIdentities := identityMap["userAssignedIdentities"].(map[string]interface{})
			for key := range userAssignedIdentities {
				identityId, err := parse.UserAssignedIdentitiesID(key)
				if err == nil {
					identityIds = append(identityIds, basetypes.NewStringValue(identityId.ID()))
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

		principalId := ""
		if v := identityMap["principalId"]; v != nil {
			principalId = v.(string)
		}

		tenantId := ""
		if v := identityMap["tenantId"]; v != nil {
			tenantId = v.(string)
		}
		return &Model{
			Type:        basetypes.NewStringValue(identityType),
			IdentityIDs: basetypes.NewListValueMust(types.StringType, identityIds),
			PrincipalID: basetypes.NewStringValue(principalId),
			TenantID:    basetypes.NewStringValue(tenantId),
		}
	}
	return nil
}

func FromList(input types.List) Model {
	identityModel := Model{
		Type: types.StringValue(string(None)),
	}
	elements := input.Elements()
	if len(elements) == 0 {
		return identityModel
	}
	diags := elements[0].(basetypes.ObjectValue).As(context.Background(), &identityModel, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty:    false,
		UnhandledUnknownAsEmpty: false,
	})
	if diags.HasError() {
		tflog.Warn(context.Background(), fmt.Sprintf("failed to convert list to identity: %s", diags))
	}
	return identityModel
}

func ToList(input Model) types.List {
	return basetypes.NewListValueMust(input.ObjectType(), []attr.Value{input.Value()})
}
