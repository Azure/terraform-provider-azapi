package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.ListTypable = ListType{}
var _ attr.TypeWithMarkdownDescription = ListType{}

type ListType struct {
	description string
	basetypes.ListType
}

func NewListType(description string, elemType attr.Type) ListType {
	return ListType{
		description: description,
		ListType:    basetypes.ListType{ElemType: elemType},
	}
}

func (s ListType) Equal(o attr.Type) bool {
	if s.ListType.ElemType == nil {
		return false
	}

	other, ok := o.(ListType)

	if !ok {
		return false
	}

	return s.ElementType().Equal(other.ElementType())
}

func (s ListType) String() string {
	return "fwdtypes.ListType"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s ListType) MarkdownDescription(context.Context) string {
	return s.description
}
