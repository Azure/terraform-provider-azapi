package metadata

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type DataType struct {
	isblk bool
	inner attr.Type
}

func (dt DataType) String() string {
	switch typ := dt.inner.(type) {
	case basetypes.BoolTypable:
		return "Boolean"
	case basetypes.Float32Typable:
		return "Float32"
	case basetypes.Float64Typable:
		return "Float64"
	case basetypes.Int32Typable:
		return "Int32"
	case basetypes.Int64Typable:
		return "Int64"
	case basetypes.NumberTypable:
		return "Number"
	case basetypes.StringTypable:
		return "String"
	case basetypes.DynamicTypable:
		return "Dynamic"
	case basetypes.ListTypable:
		if etype := maybeElementTypeString(dt.isblk, typ); etype != "" {
			return "List of " + etype
		} else {
			return "List"
		}
	case basetypes.SetTypable:
		if etype := maybeElementTypeString(dt.isblk, typ); etype != "" {
			return "Set of " + etype
		} else {
			return "Set"
		}
	case basetypes.MapTypable:
		if etype := maybeElementTypeString(dt.isblk, typ); etype != "" {
			return "Map of " + etype
		} else {
			return "Map"
		}
	case basetypes.ObjectTypable:
		if dt.isblk {
			return "Block"
		} else {
			return "Object"
		}
	case basetypes.TupleType:
		// Note there is no TupleTypable type in the fw.
		// TODO consider enrich this.
		return "Tuple"
	default:
		panic(fmt.Sprintf("unhandled inner data type: %t", dt.inner))
	}
}

func maybeElementTypeString(isblk bool, typ attr.Type) string {
	if typeWithElement, ok := typ.(attr.TypeWithElementType); ok {
		return DataType{
			isblk: isblk,
			inner: typeWithElement.ElementType(),
		}.String()
	}
	return ""
}
