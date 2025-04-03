package types

import (
	"fmt"
	"strconv"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ TypeBase = &ArrayType{}

type ArrayType struct {
	Type      string         `json:"$type"`
	ItemType  *TypeReference `json:"itemType"`
	MinLength *int           `json:"minLength"`
	MaxLength *int           `json:"maxLength"`
}

func (t *ArrayType) GetReadOnly(i interface{}) interface{} {
	if t == nil || i == nil {
		return nil
	}
	var itemType *TypeBase
	if t.ItemType != nil {
		itemType = t.ItemType.Type
	}
	// check body type
	bodyArray, ok := i.([]interface{})
	if !ok {
		return nil
	}
	if itemType == nil {
		return bodyArray
	}

	res := make([]interface{}, 0)
	for _, value := range bodyArray {
		res = append(res, (*itemType).GetReadOnly(value))
	}
	return res
}

func (t *ArrayType) GetWriteOnly(body interface{}) interface{} {
	if t == nil || body == nil {
		return nil
	}
	var itemType *TypeBase
	if t.ItemType != nil {
		itemType = t.ItemType.Type
	}
	// check body type
	bodyArray, ok := body.([]interface{})
	if !ok {
		return nil
	}
	res := make([]interface{}, 0)
	for _, value := range bodyArray {
		if itemType != nil {
			res = append(res, (*itemType).GetWriteOnly(value))
		}
	}
	return res
}

func (t *ArrayType) Validate(body attr.Value, path string) []error {
	if t == nil || body == nil || body.IsNull() || body.IsUnknown() {
		return []error{}
	}
	errors := make([]error, 0)
	var itemType *TypeBase
	if t.ItemType != nil {
		itemType = t.ItemType.Type
	}
	// check body type
	var items []attr.Value
	switch v := body.(type) {
	case types.List:
		items = v.Elements()
	case types.Tuple:
		items = v.Elements()
	case types.Set:
		items = v.Elements()
	case types.Dynamic:
		return t.Validate(v.UnderlyingValue(), path)
	default:
		errors = append(errors, utils.ErrorMismatch(path, "array", fmt.Sprintf("%T", body)))
	}

	// check the length
	if t.MinLength != nil && len(items) < *t.MinLength {
		errors = append(errors, utils.ErrorCommon(path, fmt.Sprintf("array length is less than %d", *t.MinLength)))
	}

	if t.MaxLength != nil && len(items) > *t.MaxLength {
		errors = append(errors, utils.ErrorCommon(path, fmt.Sprintf("array length is greater than %d", *t.MaxLength)))
	}

	for index, value := range items {
		if itemType != nil {
			errors = append(errors, (*itemType).Validate(value, path+"."+strconv.Itoa(index))...)
		}
	}
	return errors
}

func (t *ArrayType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}
