package types

import (
	"fmt"
	"log"
	"regexp"

	"github.com/Azure/terraform-provider-azapi/internal/azure/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ TypeBase = &StringType{}

type StringType struct {
	Type      string `json:"$type"`
	MinLength *int   `json:"minLength"`
	MaxLength *int   `json:"maxLength"`
	Sensitive bool   `json:"sensitive"`
	Pattern   string `json:"pattern"`
}

func (s *StringType) GetReadOnly(i interface{}) interface{} {
	if s == nil || i == nil {
		return nil
	}
	return i
}

func (s *StringType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(s)
	return &typeBase
}

func (s *StringType) Validate(body attr.Value, path string) []error {
	if s == nil || body == nil || body.IsNull() || body.IsUnknown() {
		return nil
	}

	var v string
	switch input := body.(type) {
	case types.String:
		v = input.ValueString()
	case types.Dynamic:
		return s.Validate(input.UnderlyingValue(), path)
	default:
		return []error{utils.ErrorMismatch(path, "string", fmt.Sprintf("%T", body))}
	}
	if s.MinLength != nil && len(v) < *s.MinLength {
		return []error{utils.ErrorCommon(path, fmt.Sprintf("string length is less than %d", *s.MinLength))}
	}
	if s.MaxLength != nil && len(v) > *s.MaxLength {
		return []error{utils.ErrorCommon(path, fmt.Sprintf("string length is greater than %d", *s.MaxLength))}
	}
	if s.Pattern != "" {
		isMatch, err := regexp.Match(s.Pattern, []byte(v))
		if err != nil {
			log.Printf("[WARN] failed to match pattern %s: %s", s.Pattern, err)
			return nil
		}
		if !isMatch {
			return []error{utils.ErrorCommon(path, fmt.Sprintf("string does not match pattern %s", s.Pattern))}
		}
	}
	return nil
}

func (s *StringType) GetWriteOnly(i interface{}) interface{} {
	if s == nil || i == nil {
		return nil
	}
	return i
}
