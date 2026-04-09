package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.Float64Typable = Float64Type{}
var _ attr.TypeWithMarkdownDescription = Float64Type{}

type Float64Type struct {
	description string
	basetypes.Float64Type
}

func NewFloat64Type(description string) Float64Type {
	return Float64Type{
		description: description,
		Float64Type: basetypes.Float64Type{},
	}
}

func (s Float64Type) Equal(o attr.Type) bool {
	_, ok := o.(Float64Type)
	return ok
}

func (s Float64Type) String() string {
	return "fwdtypes.Float64Type"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s Float64Type) MarkdownDescription(context.Context) string {
	return s.description
}
