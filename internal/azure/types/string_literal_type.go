package types

import (
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ TypeBase = &StringLiteralType{}

type StringLiteralType struct {
	Type  string `json:"$type"`
	Value string `json:"value"`
}

func (t *StringLiteralType) GetReadOnly(i interface{}) interface{} {
	if t == nil || i == nil {
		return nil
	}
	return i
}

func (t *StringLiteralType) GetWriteOnly(i interface{}) interface{} {
	if t == nil || i == nil {
		return nil
	}
	return i
}

func (t *StringLiteralType) Validate(body attr.Value, path string) []error {
	if t == nil || body == nil || body.IsNull() || body.IsUnknown() {
		return []error{}
	}
	errors := make([]error, 0)

	switch v := body.(type) {
	case types.String:
		if v.ValueString() != t.Value {
			errors = append(errors, utils.ErrorMismatch(path, t.Value, v.ValueString()))
		}
	case types.Dynamic:
		return t.Validate(v.UnderlyingValue(), path)
	default:
		errors = append(errors, utils.ErrorMismatch(path, "string", fmt.Sprintf("%T", body)))
	}
	return errors
}

func (t *StringLiteralType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}
