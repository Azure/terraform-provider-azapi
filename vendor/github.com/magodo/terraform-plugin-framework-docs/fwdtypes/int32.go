package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.Int32Typable = Int32Type{}
var _ attr.TypeWithMarkdownDescription = Int32Type{}

type Int32Type struct {
	description string
	basetypes.Int32Type
}

func NewInt32Type(description string) Int32Type {
	return Int32Type{
		description: description,
		Int32Type:   basetypes.Int32Type{},
	}
}

func (s Int32Type) Equal(o attr.Type) bool {
	_, ok := o.(Int32Type)
	return ok
}

func (s Int32Type) String() string {
	return "fwdtypes.Int32Type"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s Int32Type) MarkdownDescription(context.Context) string {
	return s.description
}
