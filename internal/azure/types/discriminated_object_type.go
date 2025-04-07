package types

import (
	"context"
	"fmt"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ TypeBase = &DiscriminatedObjectType{}

type DiscriminatedObjectType struct {
	Type           string                    `json:"$type"`
	Name           string                    `json:"name"`
	Discriminator  string                    `json:"discriminator"`
	BaseProperties map[string]ObjectProperty `json:"baseProperties"`
	Elements       map[string]*TypeReference `json:"elements"`
}

func (t *DiscriminatedObjectType) GetReadOnly(body interface{}) interface{} {
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
			if def.Type != nil && def.Type.Type != nil {
				res[key] = (*def.Type.Type).GetReadOnly(bodyMap[key])
			} else {
				res[key] = bodyMap[key]
			}
		}
	}

	if _, ok := bodyMap[t.Discriminator]; !ok {
		return nil
	}

	if discriminator, ok := bodyMap[t.Discriminator].(string); ok {
		if t.Elements[discriminator] != nil && t.Elements[discriminator].Type != nil {
			if additionalProps := (*t.Elements[discriminator].Type).GetReadOnly(body); additionalProps != nil {
				if additionalMap, ok := additionalProps.(map[string]interface{}); ok {
					for key, value := range additionalMap {
						res[key] = value
					}
					return res
				}
			}
		}
	}

	// if the discriminator's type is not in the embedded schema, add unchecked properties to res
	for key, value := range bodyMap {
		if _, ok := t.BaseProperties[key]; ok {
			continue
		}
		res[key] = value
	}
	return res
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
		if t.Elements[discriminator] != nil && t.Elements[discriminator].Type != nil {
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

	// if the discriminator's type is not in the embedded schema, add unchecked properties to res
	for key, value := range bodyMap {
		if _, ok := t.BaseProperties[key]; ok {
			continue
		}
		res[key] = value
	}
	return res
}

func (t *DiscriminatedObjectType) Validate(body attr.Value, path string) []error {
	if t == nil || body == nil || body.IsNull() || body.IsUnknown() {
		return []error{}
	}

	errors := make([]error, 0)

	var bodyMap map[string]attr.Value
	switch v := body.(type) {
	case types.Object:
		bodyMap = v.Attributes()
	case types.Map:
		bodyMap = v.Elements()
	case types.Dynamic:
		return t.Validate(v.UnderlyingValue(), path)
	default:
		errors = append(errors, utils.ErrorMismatch(path, "object", fmt.Sprintf("%T", body)))
		return errors
	}

	// check base properties
	otherProperties := make(map[string]attr.Value)
	otherPropertyTypes := make(map[string]attr.Type)
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
			otherPropertyTypes[key] = value.Type(context.Background())
		}
	}

	// check required base properties
	for key, value := range t.BaseProperties {
		if !value.IsRequired() {
			continue
		}
		if _, ok := bodyMap[key]; !ok {
			errors = append(errors, utils.ErrorShouldDefine(path+"."+key))
		}
	}

	// check other properties which should be defined in discriminated objects
	if _, ok := otherProperties[t.Discriminator]; !ok {
		errors = append(errors, utils.ErrorShouldDefine(path+"."+t.Discriminator))
		return errors
	}

	if discriminatorStr, ok := otherProperties[t.Discriminator].(types.String); ok {
		discriminator := discriminatorStr.ValueString()
		switch {
		case t.Elements[discriminator] == nil:
			options := make([]string, 0)
			for key := range t.Elements {
				options = append(options, key)
			}
			errors = append(errors, utils.ErrorNotMatchAnyValues(path+"."+t.Discriminator, discriminator, options))
		case t.Elements[discriminator].Type != nil:
			errors = append(errors, (*t.Elements[discriminator].Type).Validate(types.ObjectValueMust(otherPropertyTypes, otherProperties), path)...)
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
