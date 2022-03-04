package azure

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
)

type Schema struct {
	Resources map[string]*Resource
}

type Resource struct {
	Definitions []ResourceDefinition
}

type ResourceDefinition struct {
	Definition *types.ResourceType
	Location   TypeLocation
	ApiVersion string
}

type TypeLocation struct {
	Location string `json:"RelativePath"`
	Index    int    `json:"Index"`
}

func (o *Schema) UnmarshalJSON(body []byte) error {
	var m map[string]map[string]TypeLocation
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	if m["Resources"] == nil {
		return nil
	}
	o.Resources = make(map[string]*Resource)
	for k, v := range m["Resources"] {
		index := strings.Index(k, "@")
		if index == -1 {
			return fmt.Errorf("api-version is not specified, type: %s", k)
		}
		resourceType := k[0:index]
		resource := o.Resources[resourceType]
		if resource == nil {
			o.Resources[resourceType] = &Resource{
				Definitions: make([]ResourceDefinition, 0),
			}
			resource = o.Resources[resourceType]
		}
		resource.Definitions = append(resource.Definitions, ResourceDefinition{
			Definition: nil,
			Location:   v,
			ApiVersion: k[index+1:],
		})
	}

	return nil
}

func (o *ResourceDefinition) GetDefinition() (*types.ResourceType, error) {
	if o == nil {
		return nil, nil
	}
	if o.Definition != nil {
		return o.Definition, nil
	}
	definition, err := o.Location.LoadDefinition()
	if err != nil {
		return nil, err
	}
	o.Definition = definition
	return o.Definition, nil
}

func (o *TypeLocation) LoadDefinition() (*types.ResourceType, error) {
	if o == nil {
		return nil, nil
	}
	data, err := StaticFiles.ReadFile("generated/" + o.Location)
	if err != nil {
		return nil, err
	}
	var schema types.Schema
	err = json.Unmarshal(data, &schema)
	if err != nil {
		return nil, err
	}
	if o.Index < len(schema.Types) && schema.Types[o.Index] != nil {
		if resourceType, ok := (*schema.Types[o.Index]).(*types.ResourceType); ok {
			return resourceType, nil
		}
	}
	return nil, fmt.Errorf("index invalid or the type is not a resource type")
}
