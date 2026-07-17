package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.ObjectTypable = ObjectType{}
var _ attr.TypeWithMarkdownDescription = ObjectType{}

type ObjectType struct {
	description string
	basetypes.ObjectType
}

func NewObjectType(description string, attrTypes map[string]attr.Type) ObjectType {
	return ObjectType{
		description: description,
		ObjectType:  basetypes.ObjectType{AttrTypes: attrTypes},
	}
}

func (s ObjectType) Equal(candidate attr.Type) bool {
	other, ok := candidate.(ObjectType)
	if !ok {
		return false
	}
	if len(other.AttrTypes) != len(s.AttrTypes) {
		return false
	}
	for k, v := range s.AttrTypes {
		attr, ok := other.AttrTypes[k]
		if !ok {
			return false
		}
		if !v.Equal(attr) {
			return false
		}
	}
	return true
}

func (s ObjectType) String() string {
	return "fwdtypes.ObjectType"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s ObjectType) MarkdownDescription(context.Context) string {
	return s.description
}
