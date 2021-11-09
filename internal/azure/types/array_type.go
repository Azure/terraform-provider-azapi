package types

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

var _ TypeBase = &ArrayType{}

type ArrayType struct {
	ItemType *TypeReference
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

func (t *ArrayType) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		if k == "Type" {
			if v != nil {
				var index int
				err := json.Unmarshal(*v, &index)
				if err != nil {
					return err
				}
				t.ItemType = &TypeReference{TypeIndex: index}
			}
		}
	}
	return nil
}
