package types

import (
	"encoding/json"
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
)

var _ TypeBase = &ObjectType{}

type ObjectType struct {
	Name                 string
	Properties           map[string]ObjectProperty
	AdditionalProperties *TypeReference
}

func (t *ObjectType) GetWriteOnly(body interface{}) interface{} {
	if t == nil || body == nil {
		return nil
	}
	// check body type
	bodyMap, ok := body.(map[string]interface{})
	if !ok {
		return nil
	}

	res := make(map[string]interface{})
	for key, def := range t.Properties {
		if _, ok := bodyMap[key]; ok {
			if (def.IsRequired() || (!def.IsReadOnly() && !def.IsDeployTimeConstant())) && def.Type != nil && def.Type.Type != nil {
				res[key] = (*def.Type.Type).GetWriteOnly(bodyMap[key])
			}
		}
	}

	if t.AdditionalProperties != nil && t.AdditionalProperties.Type != nil {
		for key, value := range bodyMap {
			if _, ok := t.Properties[key]; ok {
				continue
			}
			res[key] = (*t.AdditionalProperties.Type).GetWriteOnly(value)
		}
	}
	return res
}

func (t *ObjectType) Validate(body interface{}, path string) []error {
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
	// check properties defined in body, but not in schema
	for key, value := range bodyMap {
		if def, ok := t.Properties[key]; ok {
			if def.IsReadOnly() {
				errors = append(errors, utils.ErrorShouldNotDefineReadOnly(path+"."+key))
				continue
			}
			var valueDefType *TypeBase
			if def.Type != nil && def.Type.Type != nil {
				valueDefType = def.Type.Type
				errors = append(errors, (*valueDefType).Validate(value, path+"."+key)...)
			}
			continue
		}
		if t.AdditionalProperties != nil && t.AdditionalProperties.Type != nil {
			errors = append(errors, (*t.AdditionalProperties.Type).Validate(value, path+"."+key)...)
		} else {
			options := make([]string, 0)
			for key := range t.Properties {
				options = append(options, path+"."+key)
			}
			errors = append(errors, utils.ErrorShouldNotDefine(path+"."+key, options))
		}
	}

	// check properties required in schema, but not in body
	for key, value := range t.Properties {
		if value.IsRequired() && bodyMap[key] == nil {
			// skip name in body
			if path == "" && key == "name" {
				continue
			}
			errors = append(errors, utils.ErrorShouldDefine(path+"."+key))
		}
	}
	return errors
}

func (t *ObjectType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t *ObjectType) UnmarshalJSON(body []byte) error {
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
		case "Properties":
			if v != nil {
				var properties map[string]ObjectProperty
				err := json.Unmarshal(*v, &properties)
				if err != nil {
					return err
				}
				t.Properties = properties
			}
		case "AdditionalProperties":
			if v != nil {
				var index int
				err := json.Unmarshal(*v, &index)
				if err != nil {
					return err
				}
				t.AdditionalProperties = &TypeReference{TypeIndex: index}
			}
		default:
			return fmt.Errorf("unmarshalling object type, unrecognized key: %s", k)
		}
	}

	return nil
}

type ObjectProperty struct {
	Type        *TypeReference
	Flags       []ObjectPropertyFlag
	Description *string
}

func (o ObjectProperty) IsRequired() bool {
	for _, value := range o.Flags {
		if value == Required {
			return true
		}
	}
	return false
}

func (o ObjectProperty) IsReadOnly() bool {
	for _, value := range o.Flags {
		if value == ReadOnly {
			return true
		}
	}
	return false
}

func (o ObjectProperty) IsDeployTimeConstant() bool {
	for _, value := range o.Flags {
		if value == DeployTimeConstant {
			return true
		}
	}
	return false
}

func (o *ObjectProperty) UnmarshalJSON(body []byte) error {
	var m map[string]*json.RawMessage
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	for k, v := range m {
		switch k {
		case "Description":
			if v != nil {
				var description string
				err := json.Unmarshal(*v, &description)
				if err != nil {
					return err
				}
				o.Description = &description
			}
		case "Flags":
			if v != nil {
				var flag int
				err := json.Unmarshal(*v, &flag)
				if err != nil {
					return err
				}
				flags := make([]ObjectPropertyFlag, 0)
				for _, f := range PossibleObjectPropertyFlagValues() {
					if flag&int(f) != 0 {
						flags = append(flags, f)
					}
				}
				o.Flags = flags
			}
		case "Type":
			if v != nil {
				var index int
				err := json.Unmarshal(*v, &index)
				if err != nil {
					return err
				}
				o.Type = &TypeReference{TypeIndex: index}
			}
		default:
			return fmt.Errorf("unmarshalling object property, unrecognized key: %s", k)
		}
	}
	return nil
}

type ObjectPropertyFlag int

const (
	None ObjectPropertyFlag = 0

	Required ObjectPropertyFlag = 1 << 0

	ReadOnly ObjectPropertyFlag = 1 << 1

	WriteOnly ObjectPropertyFlag = 1 << 2

	DeployTimeConstant ObjectPropertyFlag = 1 << 3
)

func PossibleObjectPropertyFlagValues() []ObjectPropertyFlag {
	return []ObjectPropertyFlag{None, Required, ReadOnly, WriteOnly, DeployTimeConstant}
}
