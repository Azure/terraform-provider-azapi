package types

import (
	"encoding/json"
	"fmt"

	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

var _ TypeBase = &DiscriminatedObjectType{}

type DiscriminatedObjectType struct {
	Name           string
	Discriminator  string
	BaseProperties map[string]ObjectProperty
	Elements       map[string]*TypeReference
}

func (t *DiscriminatedObjectType) Validate(body interface{}, path string) []error {
	if t == nil || body == nil {
		return []error{}
	}
	errors := make([]error, 0)
	// check body type
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		errors = append(errors, utils.ErrorMismatch(path, "object", fmt.Sprintf("%T", body)))
		return errors
	}

	// check base properties
	otherProperties := make(map[string]interface{})
	for key, value := range bodyMap {
		if def, ok := t.BaseProperties[key]; ok {
			var valueDefType *TypeBase
			if def.Type != nil && def.Type.Type != nil {
				valueDefType = def.Type.Type
				errors = append(errors, (*valueDefType).Validate(value, path+"."+key)...)
			}
		} else {
			otherProperties[key] = value
		}
	}

	// check required base properties
	for key, value := range t.BaseProperties {
		if value.IsRequired() && bodyMap[key] == nil {
			errors = append(errors, utils.ErrorShouldDefine(path+"."+key))
		}
	}

	// check other properties which should be defined in discriminated objects
	valid := false
	errs := make([]error, 0)
	for _, element := range t.Elements {
		if element.Type == nil {
			continue
		}
		temp := (*element.Type).Validate(otherProperties, path)
		if len(temp) < len(errs) || len(errs) == 0 {
			errs = temp
		}
		if len(temp) == 0 {
			valid = true
			break
		}
	}
	errors = append(errors, errs...)
	if !valid {
		errors = append(errors, utils.ErrorNotMatchAny(path))
	}
	return errors
}

func (t *DiscriminatedObjectType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t *DiscriminatedObjectType) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "Name":
			if v != nil {
				var name string
				err := json.Unmarshal(*v, &name)
				if err != nil {
					return err
				}
				t.Name = name
			}
		case "Discriminator":
			if v != nil {
				var discriminator string
				err := json.Unmarshal(*v, &discriminator)
				if err != nil {
					return err
				}
				t.Discriminator = discriminator
			}
		case "BaseProperties":
			if v != nil {
				var baseProperties map[string]ObjectProperty
				err := json.Unmarshal(*v, &baseProperties)
				if err != nil {
					return err
				}
				t.BaseProperties = baseProperties
			}
		case "Elements":
			if v != nil {
				var elementIndexes map[string]int
				err := json.Unmarshal(*v, &elementIndexes)
				if err != nil {
					return err
				}
				elements := make(map[string]*TypeReference, 0)
				for key, index := range elementIndexes {
					elements[key] = &TypeReference{TypeIndex: index}
				}
				t.Elements = elements
			}
		}
	}
	return nil
}
