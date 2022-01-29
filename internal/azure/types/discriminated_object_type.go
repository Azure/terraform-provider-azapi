package types

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/utils"
)

var _ TypeBase = &DiscriminatedObjectType{}

type DiscriminatedObjectType struct {
	Name           string
	Discriminator  string
	BaseProperties map[string]ObjectProperty
	Elements       map[string]*TypeReference
}

func (t *DiscriminatedObjectType) GetWriteOnly(body interface{}) interface{} {
	if t == nil || body == nil {
		return []error{}
	}
	// check body type
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil
	}

	res := make(map[string]interface{})
	for key, def := range t.BaseProperties {
		if _, ok := bodyMap[key]; ok {
			if !def.IsReadOnly() && def.Type != nil && def.Type.Type != nil {
				res[key] = (*def.Type.Type).GetWriteOnly(bodyMap[key])
			}
		}
	}

	if _, ok := bodyMap[t.Discriminator]; !ok {
		return nil
	}

	if discriminator, ok := bodyMap[t.Discriminator].(string); ok {
		if t.Elements[discriminator].Type != nil {
			if additionalProps := (*t.Elements[discriminator].Type).GetWriteOnly(body); additionalProps != nil {
				if additionalMap, ok := additionalProps.(map[string]interface{}); ok {
					for key, value := range additionalMap {
						res[key] = value
					}
					return res
				}
			}
		}
	}
	return nil
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
			if def.IsReadOnly() {
				errors = append(errors, utils.ErrorShouldNotDefineReadOnly(path+"."+key))
				continue
			}
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
	if _, ok := otherProperties[t.Discriminator]; !ok {
		errors = append(errors, utils.ErrorShouldDefine(path+"."+t.Discriminator))
		return errors
	}

	if discriminator, ok := otherProperties[t.Discriminator].(string); ok {
		switch {
		case t.Elements[discriminator] == nil:
			options := make([]string, 0)
			for key := range t.Elements {
				options = append(options, key)
			}
			errors = append(errors, utils.ErrorNotMatchAnyValues(path+"."+t.Discriminator, discriminator, options))
		case t.Elements[discriminator].Type != nil:
			errors = append(errors, (*t.Elements[discriminator].Type).Validate(otherProperties, path)...)
		}
	} else {
		errors = append(errors, utils.ErrorMismatch(path+"."+t.Discriminator, "string", fmt.Sprintf("%T", otherProperties[t.Discriminator])))
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
				elements := make(map[string]*TypeReference)
				for key, index := range elementIndexes {
					elements[key] = &TypeReference{TypeIndex: index}
				}
				t.Elements = elements
			}
		}
	}
	return nil
}
