package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.NumberTypable = NumberType{}
var _ attr.TypeWithMarkdownDescription = NumberType{}

type NumberType struct {
	description string
	basetypes.NumberType
}

func NewNumberType(description string) NumberType {
	return NumberType{
		description: description,
		NumberType:  basetypes.NumberType{},
	}
}

func (s NumberType) Equal(o attr.Type) bool {
	_, ok := o.(NumberType)

	return ok
}

func (s NumberType) String() string {
	return "fwdtypes.NumberType"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s NumberType) MarkdownDescription(context.Context) string {
	return s.description
}
