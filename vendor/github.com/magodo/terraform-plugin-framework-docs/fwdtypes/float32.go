package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.Float32Typable = Float32Type{}
var _ attr.TypeWithMarkdownDescription = Float32Type{}

type Float32Type struct {
	description string
	basetypes.Float32Type
}

func NewFloat32Type(description string) Float32Type {
	return Float32Type{
		description: description,
		Float32Type: basetypes.Float32Type{},
	}
}

func (s Float32Type) Equal(o attr.Type) bool {
	_, ok := o.(Float32Type)
	return ok
}

func (s Float32Type) String() string {
	return "fwdtypes.Float32Type"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s Float32Type) MarkdownDescription(context.Context) string {
	return s.description
}
