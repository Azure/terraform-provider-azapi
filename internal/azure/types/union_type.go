package types

import (
	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
)

var _ TypeBase = &UnionType{}

type UnionType struct {
	Type     string           `json:"$type"`
	Elements []*TypeReference `json:"elements"`
}

func (t *UnionType) GetReadOnly(i interface{}) interface{} {
	if t == nil || i == nil {
		return nil
	}
	// TODO: improve this
	return i
}

func (t *UnionType) GetWriteOnly(body interface{}) interface{} {
	if t == nil || body == nil {
		return nil
	}
	// TODO: improve this
	return body
}

func (t *UnionType) Validate(body interface{}, path string) []error {
	if t == nil || body == nil {
		return []error{}
	}
	errors := make([]error, 0)
	valid := false
	for _, element := range t.Elements {
		if element.Type == nil {
			continue
		}
		temp := (*element.Type).Validate(body, path)
		if len(temp) == 0 {
			valid = true
			break
		}
	}
	if !valid {
		options := make([]string, 0)
		for _, element := range t.Elements {
			if element.Type != nil {
				if stringLiteralType, ok := (*element.Type).(*StringLiteralType); ok {
					options = append(options, stringLiteralType.Value)
				}
			}
		}
		if len(options) == 0 {
			errors = append(errors, utils.ErrorNotMatchAny(path))
		} else {
			value := ""
			if current, ok := body.(string); ok {
				value = current
			}
			errors = append(errors, utils.ErrorNotMatchAnyValues(path, value, options))
		}
	}
	return errors
}

func (t *UnionType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}
