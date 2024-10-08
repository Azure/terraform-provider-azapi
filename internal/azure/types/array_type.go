package types

import (
	"fmt"
	"strconv"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
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

func (t *ArrayType) Validate(body interface{}, path string) []error {
	if t == nil || body == nil {
		return []error{}
	}
	errors := make([]error, 0)
	var itemType *TypeBase
	if t.ItemType != nil {
		itemType = t.ItemType.Type
	}
	// check body type
	bodyArray, ok := body.([]interface{})
	if !ok {
		errors = append(errors, utils.ErrorMismatch(path, "array", fmt.Sprintf("%T", body)))
		return errors
	}

	// check the length
	if t.MinLength != nil && len(bodyArray) < *t.MinLength {
		errors = append(errors, utils.ErrorCommon(path, fmt.Sprintf("array length is less than %d", *t.MinLength)))
	}

	if t.MaxLength != nil && len(bodyArray) > *t.MaxLength {
		errors = append(errors, utils.ErrorCommon(path, fmt.Sprintf("array length is greater than %d", *t.MaxLength)))
	}

	for index, value := range bodyArray {
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
