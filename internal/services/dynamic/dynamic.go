package dynamic

import (
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type UnknownValueHandler func(val attr.Value) ([]byte, error)

func ToJSON(d types.Dynamic) ([]byte, error) {
	return attrValueToJSON(d.UnderlyingValue(), nil)
}

func ToJSONWithUnknownValueHandler(d types.Dynamic, handler UnknownValueHandler) ([]byte, error) {
	return attrValueToJSON(d.UnderlyingValue(), handler)
}

func attrListToJSON(in []attr.Value, handler UnknownValueHandler) ([]json.RawMessage, error) {
	l := make([]json.RawMessage, 0)
	for _, v := range in {
		vv, err := attrValueToJSON(v, handler)
		if err != nil {
			return nil, err
		}
		l = append(l, json.RawMessage(vv))
	}
	return l, nil
}

func attrMapToJSON(in map[string]attr.Value, handler UnknownValueHandler) (map[string]json.RawMessage, error) {
	m := map[string]json.RawMessage{}
	for k, v := range in {
		vv, err := attrValueToJSON(v, handler)
		if err != nil {
			return nil, err
		}
		m[k] = json.RawMessage(vv)
	}
	return m, nil
}

func attrValueToJSON(val attr.Value, handler UnknownValueHandler) ([]byte, error) {
	if val == nil || val.IsNull() {
		return json.Marshal(nil)
	}
	if val.IsUnknown() {
		if handler != nil {
			return handler(val)
		}
	}
	switch value := val.(type) {
	case types.Bool:
		return json.Marshal(value.ValueBool())
	case types.String:
		return json.Marshal(value.ValueString())
	case types.Int64:
		return json.Marshal(value.ValueInt64())
	case types.Float64:
		return json.Marshal(value.ValueFloat64())
	case types.Number:
		v, _ := value.ValueBigFloat().Float64()
		return json.Marshal(v)
	case types.List:
		l, err := attrListToJSON(value.Elements(), handler)
		if err != nil {
			return nil, err
		}
		return json.Marshal(l)
	case types.Set:
		l, err := attrListToJSON(value.Elements(), handler)
		if err != nil {
			return nil, err
		}
		return json.Marshal(l)
	case types.Tuple:
		l, err := attrListToJSON(value.Elements(), handler)
		if err != nil {
			return nil, err
		}
		return json.Marshal(l)
	case types.Map:
		m, err := attrMapToJSON(value.Elements(), handler)
		if err != nil {
			return nil, err
		}
		return json.Marshal(m)
	case types.Object:
		m, err := attrMapToJSON(value.Attributes(), handler)
		if err != nil {
			return nil, err
		}
		return json.Marshal(m)
	default:
		return nil, fmt.Errorf("Unhandled type: %T", value)
	}
}

func FromJSON(b []byte, typ attr.Type) (types.Dynamic, error) {
	v, err := attrValueFromJSON(b, typ)
	if err != nil {
		return types.Dynamic{}, err
	}
	return types.DynamicValue(v), nil
}

func attrListFromJSON(b []byte, etyp attr.Type) ([]attr.Value, error) {
	var l []json.RawMessage
	if err := json.Unmarshal(b, &l); err != nil {
		return nil, err
	}
	vals := make([]attr.Value, 0)
	for _, b := range l {
		val, err := attrValueFromJSON(b, etyp)
		if err != nil {
			return nil, err
		}
		vals = append(vals, val)
	}
	return vals, nil
}

func attrValueFromJSON(b []byte, typ attr.Type) (attr.Value, error) {
	switch typ := typ.(type) {
	case basetypes.BoolType:
		if b == nil || string(b) == "null" {
			return types.BoolNull(), nil
		}
		var v bool
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		return types.BoolValue(v), nil
	case basetypes.StringType:
		if b == nil || string(b) == "null" {
			return types.StringNull(), nil
		}
		var v string
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		return types.StringValue(v), nil
	case basetypes.Int64Type:
		if b == nil || string(b) == "null" {
			return types.Int64Null(), nil
		}
		var v int64
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		return types.Int64Value(v), nil
	case basetypes.Float64Type:
		if b == nil || string(b) == "null" {
			return types.Float64Null(), nil
		}
		var v float64
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		return types.Float64Value(v), nil
	case basetypes.NumberType:
		if b == nil || string(b) == "null" {
			return types.NumberNull(), nil
		}
		var v float64
		if err := json.Unmarshal(b, &v); err != nil {
			return nil, err
		}
		return types.NumberValue(big.NewFloat(v)), nil
	case basetypes.ListType:
		if b == nil || string(b) == "null" {
			return types.ListNull(typ.ElemType), nil
		}
		vals, err := attrListFromJSON(b, typ.ElemType)
		if err != nil {
			return nil, err
		}
		vv, diags := types.ListValue(typ.ElemType, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]
			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}
		return vv, nil
	case basetypes.SetType:
		if b == nil || string(b) == "null" {
			return types.SetNull(typ.ElemType), nil
		}
		vals, err := attrListFromJSON(b, typ.ElemType)
		if err != nil {
			return nil, err
		}
		vv, diags := types.SetValue(typ.ElemType, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]
			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}
		return vv, nil
	case basetypes.TupleType:
		if b == nil || string(b) == "null" {
			return types.TupleNull(typ.ElemTypes), nil
		}
		var l []json.RawMessage
		if err := json.Unmarshal(b, &l); err != nil {
			return nil, err
		}
		if len(l) != len(typ.ElemTypes) {
			return nil, fmt.Errorf("tuple element size not match: json=%d, type=%d", len(l), len(typ.ElemTypes))
		}
		vals := make([]attr.Value, 0)
		for i, b := range l {
			val, err := attrValueFromJSON(b, typ.ElemTypes[i])
			if err != nil {
				return nil, err
			}
			vals = append(vals, val)
		}
		vv, diags := types.TupleValue(typ.ElemTypes, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]
			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}
		return vv, nil
	case basetypes.MapType:
		if b == nil || string(b) == "null" {
			return types.MapNull(typ.ElemType), nil
		}
		var m map[string]json.RawMessage
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}
		vals := map[string]attr.Value{}
		for k, v := range m {
			val, err := attrValueFromJSON(v, typ.ElemType)
			if err != nil {
				return nil, err
			}
			vals[k] = val
		}
		vv, diags := types.MapValue(typ.ElemType, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]
			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}
		return vv, nil
	case basetypes.ObjectType:
		if b == nil || string(b) == "null" {
			return types.ObjectNull(typ.AttributeTypes()), nil
		}
		var m map[string]json.RawMessage
		if err := json.Unmarshal(b, &m); err != nil {
			return nil, err
		}
		vals := map[string]attr.Value{}
		attrTypes := typ.AttributeTypes()

		for k, attrType := range attrTypes {
			val, err := attrValueFromJSON(m[k], attrType)
			if err != nil {
				return nil, err
			}
			vals[k] = val
		}
		vv, diags := types.ObjectValue(attrTypes, vals)
		if diags.HasError() {
			diag := diags.Errors()[0]
			return nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}
		return vv, nil
	case basetypes.DynamicType:
		if b == nil || string(b) == "null" {
			return types.DynamicNull(), nil
		}
		_, vv, err := attrValueFromJSONImplied(b)
		return vv, err
	default:
		return nil, fmt.Errorf("Unhandled type: %T", typ)
	}
}

// FromJSONImplied is similar to FromJSON, while it is for typeless case.
// In which case, the following type conversion rules are applied (Go -> TF):
// - bool: bool
// - float64: number
// - string: string
// - []interface{}: tuple
// - map[string]interface{}: object
// - nil: null (dynamic)
func FromJSONImplied(b []byte) (types.Dynamic, error) {
	_, v, err := attrValueFromJSONImplied(b)
	if err != nil {
		return types.Dynamic{}, err
	}
	return types.DynamicValue(v), nil
}

func attrValueFromJSONImplied(b []byte) (attr.Type, attr.Value, error) {
	if string(b) == "null" {
		return types.DynamicType, types.DynamicNull(), nil
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal(b, &object); err == nil {
		attrTypes := map[string]attr.Type{}
		attrVals := map[string]attr.Value{}
		for k, v := range object {
			attrTypes[k], attrVals[k], err = attrValueFromJSONImplied(v)
			if err != nil {
				return nil, nil, err
			}
		}
		typ := types.ObjectType{AttrTypes: attrTypes}
		val, diags := types.ObjectValue(attrTypes, attrVals)
		if diags.HasError() {
			diag := diags.Errors()[0]
			return nil, nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}
		return typ, val, nil
	}

	var array []json.RawMessage
	if err := json.Unmarshal(b, &array); err == nil {
		eTypes := []attr.Type{}
		eVals := []attr.Value{}
		for _, e := range array {
			eType, eVal, err := attrValueFromJSONImplied(e)
			if err != nil {
				return nil, nil, err
			}
			eTypes = append(eTypes, eType)
			eVals = append(eVals, eVal)
		}
		typ := types.TupleType{ElemTypes: eTypes}
		val, diags := types.TupleValue(eTypes, eVals)
		if diags.HasError() {
			diag := diags.Errors()[0]
			return nil, nil, fmt.Errorf("%s: %s", diag.Summary(), diag.Detail())
		}
		return typ, val, nil
	}

	// Primitives
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal %s: %v", string(b), err)
	}

	switch v := v.(type) {
	case bool:
		return types.BoolType, types.BoolValue(v), nil
	case float64:
		return types.NumberType, types.NumberValue(big.NewFloat(v)), nil
	case string:
		return types.StringType, types.StringValue(v), nil
	case nil:
		return types.DynamicType, types.DynamicNull(), nil
	default:
		return nil, nil, fmt.Errorf("Unhandled type: %T", v)
	}
}

func SemanticallyEqual(a, b types.Dynamic) bool {
	aJson, err := ToJSON(a)
	if err != nil {
		return false
	}
	bJson, err := ToJSON(b)
	if err != nil {
		return false
	}
	return utils.NormalizeJson(string(aJson)) == utils.NormalizeJson(string(bJson))
}

// IsFullyKnown returns true if `val` is known. If `val` is an aggregate type,
// IsFullyKnown only returns true if all elements and attributes are known, as
// well.
func IsFullyKnown(val attr.Value) bool {
	if val == nil {
		return true
	}
	if val.IsUnknown() {
		return false
	}
	switch v := val.(type) {
	case types.Dynamic:
		return IsFullyKnown(v.UnderlyingValue())
	case types.List:
		for _, e := range v.Elements() {
			if !IsFullyKnown(e) {
				return false
			}
		}
		return true
	case types.Set:
		for _, e := range v.Elements() {
			if !IsFullyKnown(e) {
				return false
			}
		}
		return true
	case types.Tuple:
		for _, e := range v.Elements() {
			if !IsFullyKnown(e) {
				return false
			}
		}
		return true
	case types.Map:
		for _, e := range v.Elements() {
			if !IsFullyKnown(e) {
				return false
			}
		}
		return true
	case types.Object:
		for _, e := range v.Attributes() {
			if !IsFullyKnown(e) {
				return false
			}
		}
		return true
	default:
		return true
	}
}
