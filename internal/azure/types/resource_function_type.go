package types

import (
	"encoding/json"
	"fmt"
)

var _ TypeBase = &ResourceFunctionType{}

type ResourceFunctionType struct {
	Name         string
	ResourceType string
	ApiVersion   string
	Input        *TypeReference
	Output       *TypeReference
}

func (t ResourceFunctionType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

func (t ResourceFunctionType) Validate(body interface{}, path string) []error {
	return []error{}
}

func (t ResourceFunctionType) GetWriteOnly(body interface{}) interface{} {
	return body
}

func (t *ResourceFunctionType) UnmarshalJSON(body []byte) error {
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
		case "ResourceType":
			if v != nil {
				var resourceType string
				err := json.Unmarshal(*v, &resourceType)
				if err != nil {
					return err
				}
				t.ResourceType = resourceType
			}
		case "ApiVersion":
			if v != nil {
				var apiVersion string
				err := json.Unmarshal(*v, &apiVersion)
				if err != nil {
					return err
				}
				t.ApiVersion = apiVersion
			}
		case "Input":
			if v != nil {
				var index int
				err := json.Unmarshal(*v, &index)
				if err != nil {
					return err
				}
				t.Input = &TypeReference{TypeIndex: index}
			}
		case "Output":
			if v != nil {
				var index int
				err := json.Unmarshal(*v, &index)
				if err != nil {
					return err
				}
				t.Output = &TypeReference{TypeIndex: index}
			}
		default:
			return fmt.Errorf("unmarshalling resource function type, unrecognized key: %s", k)
		}
	}

	return nil
}
