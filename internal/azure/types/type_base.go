package types

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
)

type TypeBase interface {
	AsTypeBase() *TypeBase
	Validate(attr.Value, string) []error
	GetWriteOnly(interface{}) interface{}
	GetReadOnly(interface{}) interface{}
}
