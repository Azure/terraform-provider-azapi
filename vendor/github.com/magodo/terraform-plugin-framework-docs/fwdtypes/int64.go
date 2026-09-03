package fwdtypes

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var _ basetypes.Int64Typable = Int64Type{}
var _ attr.TypeWithMarkdownDescription = Int64Type{}

type Int64Type struct {
	description string
	basetypes.Int64Type
}

func NewInt64Type(description string) Int64Type {
	return Int64Type{
		description: description,
		Int64Type:   basetypes.Int64Type{},
	}
}

func (s Int64Type) Equal(o attr.Type) bool {
	_, ok := o.(Int64Type)
	return ok
}

func (s Int64Type) String() string {
	return "fwdtypes.Int64Type"
}

// MarkdownDescription implements [attr.TypeWithMarkdownDescription].
func (s Int64Type) MarkdownDescription(context.Context) string {
	return s.description
}
