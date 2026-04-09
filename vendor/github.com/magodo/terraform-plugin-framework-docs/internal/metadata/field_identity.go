package metadata

import (
	"errors"
	"fmt"
	"sort"
)

type ResourceIdentityFields map[string]ResourceIdentityField

func (fields ResourceIdentityFields) RequiredFields() []ResourceIdentityField {
	var out []ResourceIdentityField
	for _, info := range fields {
		if !info.Required {
			continue
		}
		out = append(out, info)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func (fields ResourceIdentityFields) OptionalFields() []ResourceIdentityField {
	var out []ResourceIdentityField
	for _, info := range fields {
		if !info.Optional {
			continue
		}
		out = append(out, info)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Name < out[j].Name
	})
	return out
}

func (fields ResourceIdentityFields) Lint() error {
	var errs []error
	for _, field := range fields {
		if field.Description == "" {
			return fmt.Errorf("no description specified for identity field: %s", field.Name)
		}
	}
	return errors.Join(errs...)
}

type ResourceIdentityField struct {
	Name     string
	DataType DataType

	Required bool
	Optional bool

	Description           string
	customTypeDescription string
}

func (field ResourceIdentityField) Traits() string {
	return field.DataType.String()
}

func (field ResourceIdentityField) CustomTypeDescription() string {
	return Sentencefy(field.customTypeDescription)
}
