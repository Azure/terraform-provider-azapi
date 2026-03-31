package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.DynamicTypable = DynamicType{}
var _ attr.TypeWithMarkdownDescription = DynamicType{}

type DynamicType struct {
	description string
	basetypes.DynamicType
}

func NewDynamicType(description string) DynamicType {
	return DynamicType{
		description: description,
		DynamicType: basetypes.DynamicType{},
	}
}

func (s DynamicType) Equal(o attr.Type) bool {
	_, ok := o.(DynamicType)
	return ok
}

func (s DynamicType) String() string {
	return "fwdtypes.DynamicType"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s DynamicType) MarkdownDescription(context.Context) string {
	return s.description
}
