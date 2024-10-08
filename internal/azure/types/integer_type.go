package types

import (
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
)

var _ TypeBase = &IntegerType{}

type IntegerType struct {
	Type     string `json:"$type"`
	MinValue *int   `json:"minValue"`
	MaxValue *int   `json:"maxValue"`
}

func (t *IntegerType) GetReadOnly(i interface{}) interface{} {
	if t == nil || i == nil {
		return nil
	}
	return i
}

func (t *IntegerType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t *IntegerType) Validate(body interface{}, path string) []error {
	if body == nil {
		return nil
	}
	var v int
	switch input := body.(type) {
	case float64, float32:
		// TODO: skip validation for now because of the following issue:
		// the bicep-types-az parses float as integer type and it should be fixed: https://github.com/Azure/bicep-types-az/issues/1404
		return nil
	case int64:
		v = int(input)
	case int32:
		v = int(input)
	case int:
		v = input
	default:
		return []error{utils.ErrorMismatch(path, "integer", fmt.Sprintf("%T", body))}
	}
	if t.MinValue != nil && v < *t.MinValue {
		return []error{utils.ErrorCommon(path, fmt.Sprintf("value is less than %d", *t.MinValue))}
	}
	if t.MaxValue != nil && v > *t.MaxValue {
		return []error{utils.ErrorCommon(path, fmt.Sprintf("value is greater than %d", *t.MaxValue))}
	}
	return nil
}

func (t *IntegerType) GetWriteOnly(i interface{}) interface{} {
	if t == nil || i == nil {
		return nil
	}
	return i
}
