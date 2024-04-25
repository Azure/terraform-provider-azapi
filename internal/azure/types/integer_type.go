package types

import (
	"fmt"
	"math"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
)

var _ TypeBase = &IntegerType{}

type IntegerType struct {
	Type     string `json:"$type"`
	MinValue *int   `json:"minValue"`
	MaxValue *int   `json:"maxValue"`
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
	case float64:
		if math.Round(input) != input {
			return []error{utils.ErrorMismatch(path, "integer", fmt.Sprintf("%T", body))}
		}
		v = int(input)
	case float32:
		if math.Round(float64(input)) != float64(input) {
			return []error{utils.ErrorMismatch(path, "integer", fmt.Sprintf("%T", body))}
		}
		v = int(input)
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
