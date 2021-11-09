package types

import (
	"fmt"

	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

type BuiltInTypeKind int

const (
	Any         BuiltInTypeKind = 1
	NULL        BuiltInTypeKind = 2
	Bool        BuiltInTypeKind = 3
	Int         BuiltInTypeKind = 4
	String      BuiltInTypeKind = 5
	Object      BuiltInTypeKind = 6
	Array       BuiltInTypeKind = 7
	ResourceRef BuiltInTypeKind = 8
)

func (kind BuiltInTypeKind) String() string {
	switch kind {
	case Any:
		return "any"
	case NULL:
		return "null"
	case Bool:
		return "bool"
	case Int:
		return "int"
	case String:
		return "string"
	case Object:
		return "object"
	case Array:
		return "array"
	case ResourceRef:
		return "resource reference"
	}
	return ""
}

type BuiltInType struct {
	Kind BuiltInTypeKind `json:"Kind"`
}

func (t *BuiltInType) Validate(body interface{}, path string) []error {
	if t == nil || body == nil {
		return []error{}
	}
	errors := make([]error, 0)
	kind := t.Kind
	switch body.(type) {
	case string:
		if kind != String {
			errors = append(errors, utils.ErrorMismatch(path, fmt.Sprint(kind), "string"))
		}
	case int, float32, float64, int32, int64:
		if kind != Int {
			errors = append(errors, utils.ErrorMismatch(path, fmt.Sprint(kind), "number"))
		}
	case bool:
		if kind != Bool {
			errors = append(errors, utils.ErrorMismatch(path, fmt.Sprint(kind), "bool"))
		}
	case map[string]interface{}:
		if kind != Object {
			errors = append(errors, utils.ErrorMismatch(path, fmt.Sprint(kind), "object"))
		}
	case []interface{}:
		if kind != String {
			errors = append(errors, utils.ErrorMismatch(path, fmt.Sprint(kind), "array"))
		}
	default:
		// check other cases like Any/Null/ResourceRef
	}
	return errors
}

func (t *BuiltInType) AsTypeBase() *TypeBase {
	typeBase := TypeBase(t)
	return &typeBase
}

var _ TypeBase = &BuiltInType{}
