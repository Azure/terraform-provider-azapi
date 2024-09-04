package azure

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
)

type Schema struct {
	Resources map[string]*Resource
	Functions map[string]*Function
}

type Resource struct {
	Definitions []*ResourceDefinition
}

type Function struct {
	Definitions []*FunctionDefinition
}

type ResourceDefinition struct {
	Definition *types.ResourceType
	Location   TypeLocation
	ApiVersion string
	mutex      sync.Mutex
}

type FunctionDefinition struct {
	Definition *types.ResourceFunctionType
	Location   TypeLocation
	ApiVersion string
	mutex      sync.Mutex
}

type TypeLocation struct {
	Location string
	Index    int
}

func (o *TypeLocation) UnmarshalJSON(body []byte) error {
	var m types.TypeReference
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	parts := strings.Split(m.Ref, "#/")
	if len(parts) != 2 {
		return fmt.Errorf("invalid ref: %s", m.Ref)
	}
	o.Location = parts[0]
	index, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid index, ref: %s: %s", m.Ref, err)
	}
	o.Index = int(index)
	return nil
}

func (o *TypeLocation) LoadResourceTypeDefinition() (*types.ResourceType, error) {
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

func (o *TypeLocation) LoadFunctionTypeDefinition() (*types.ResourceFunctionType, error) {
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
		if resourceType, ok := (*schema.Types[o.Index]).(*types.ResourceFunctionType); ok {
			return resourceType, nil
		}
	}
	return nil, fmt.Errorf("index invalid or the type is not a resource type")
}

type IndexRaw struct {
	Resources map[string]TypeLocation              `json:"resources"`
	Functions map[string]map[string][]TypeLocation `json:"resourceFunctions"`
}

func (o *Schema) UnmarshalJSON(body []byte) error {
	var m IndexRaw
	err := json.Unmarshal(body, &m)
	if err != nil {
		return err
	}
	if m.Resources == nil {
		return fmt.Errorf("resources block is nil")
	}
	if m.Functions == nil {
		return fmt.Errorf("functions block is nil")
	}
	o.Resources = make(map[string]*Resource)
	o.Functions = make(map[string]*Function)
	for k, v := range m.Resources {
		index := strings.Index(k, "@")
		if index == -1 {
			return fmt.Errorf("api-version is not specified, type: %s", k)
		}
		resourceType := k[0:index]
		resource := o.Resources[resourceType]
		if resource == nil {
			o.Resources[resourceType] = &Resource{
				Definitions: make([]*ResourceDefinition, 0),
			}
			resource = o.Resources[resourceType]
		}
		resource.Definitions = append(resource.Definitions, &ResourceDefinition{
			Definition: nil,
			Location:   v,
			ApiVersion: k[index+1:],
		})
	}
	for k, v := range m.Functions {
		for apiVersion, arr := range v {
			function := o.Functions[k]
			if function == nil {
				o.Functions[k] = &Function{
					Definitions: make([]*FunctionDefinition, 0),
				}
				function = o.Functions[k]
			}
			for _, item := range arr {
				function.Definitions = append(function.Definitions, &FunctionDefinition{
					Definition: nil,
					Location:   item,
					ApiVersion: apiVersion,
				})
			}
		}
	}

	return nil
}

func (o *ResourceDefinition) GetDefinition() (*types.ResourceType, error) {
	if o == nil {
		return nil, nil
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.Definition != nil {
		return o.Definition, nil
	}
	definition, err := o.Location.LoadResourceTypeDefinition()
	if err != nil {
		return nil, err
	}
	o.Definition = definition
	return o.Definition, nil
}

func (o *FunctionDefinition) GetDefinition() (*types.ResourceFunctionType, error) {
	if o == nil {
		return nil, nil
	}
	o.mutex.Lock()
	defer o.mutex.Unlock()
	if o.Definition != nil {
		return o.Definition, nil
	}
	definition, err := o.Location.LoadFunctionTypeDefinition()
	if err != nil {
		return nil, err
	}
	o.Definition = definition
	return o.Definition, nil
}
