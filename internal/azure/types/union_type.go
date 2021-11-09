package types

import (
	"encoding/json"

	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

var _ TypeBase = &UnionType{}

type UnionType struct {
	Elements []*TypeReference
}

func (t *UnionType) Validate(body interface{}, path string) []error {
	if t == nil || body == nil {
		return []error{}
	}
	errors := make([]error, 0)
	valid := false
	errs := make([]error, 0)
	for _, element := range t.Elements {
		if element.Type == nil {
			continue
		}
		temp := (*element.Type).Validate(body, path)
		if len(temp) < len(errs) || len(errs) == 0 {
			errs = temp
		}
		if len(temp) == 0 {
			valid = true
			errs = make([]error, 0)
			break
		}
	}
	errors = append(errors, errs...)
	if !valid {
		errors = append(errors, utils.ErrorNotMatchAny(path))
	}
	return errors
}

func (t *UnionType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t *UnionType) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		if k == "Elements" {
			if v != nil {
				var indexes []int
				err := json.Unmarshal(*v, &indexes)
				if err != nil {
					return err
				}
				elements := make([]*TypeReference, 0)
				for _, index := range indexes {
					elements = append(elements, &TypeReference{TypeIndex: index})
				}
				t.Elements = elements
			}
		}
	}
	return nil
}
