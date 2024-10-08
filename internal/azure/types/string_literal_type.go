package types

import (
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
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

func (t *StringLiteralType) Validate(body interface{}, path string) []error {
	if t == nil || body == nil {
		return []error{}
	}
	errors := make([]error, 0)
	if stringValue, ok := body.(string); ok {
		if stringValue != t.Value {
			errors = append(errors, utils.ErrorMismatch(path, t.Value, stringValue))
		}
	} else {
		errors = append(errors, utils.ErrorMismatch(path, "string", fmt.Sprintf("%T", body)))
	}
	return errors
}

func (t *StringLiteralType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}
