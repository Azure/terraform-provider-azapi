package types

import (
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

func (t *IntegerType) Validate(body attr.Value, path string) []error {
	if body == nil || body.IsNull() || body.IsUnknown() {
		return nil
	}
	var v int64
	switch input := body.(type) {
	case types.Int32:
		v = int64(input.ValueInt32())
	case types.Int64:
		v = input.ValueInt64()
	case types.Float64, types.Float32, types.Number:
		// TODO: skip validation for now because of the following issue:
		// the bicep-types-az parses float as integer type and it should be fixed: https://github.com/Azure/bicep-types-az/issues/1404
		return nil
	case types.Dynamic:
		return t.Validate(input.UnderlyingValue(), path)
	default:
		return []error{utils.ErrorMismatch(path, "integer", fmt.Sprintf("%T", body))}
	}
	if t.MinValue != nil && v < int64(*t.MinValue) {
		return []error{utils.ErrorCommon(path, fmt.Sprintf("value is less than %d", *t.MinValue))}
	}
	if t.MaxValue != nil && v > int64(*t.MaxValue) {
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
