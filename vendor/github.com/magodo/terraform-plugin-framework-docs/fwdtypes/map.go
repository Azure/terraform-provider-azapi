package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.MapTypable = MapType{}
var _ attr.TypeWithMarkdownDescription = MapType{}

type MapType struct {
	description string
	basetypes.MapType
}

func NewMapType(description string, elemType attr.Type) MapType {
	return MapType{
		description: description,
		MapType:     basetypes.MapType{ElemType: elemType},
	}
}

func (s MapType) Equal(o attr.Type) bool {
	if s.MapType.ElementType() == nil {
		return false
	}

	other, ok := o.(MapType)

	if !ok {
		return false
	}

	return s.ElementType().Equal(other.ElementType())
}

func (s MapType) String() string {
	return "fwdtypes.MapType"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s MapType) MarkdownDescription(context.Context) string {
	return s.description
}
