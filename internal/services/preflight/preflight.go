package preflight

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func attrListToJSON(isPropertyValue bool, in []attr.Value) (bool, []json.RawMessage, error) {
	l := make([]json.RawMessage, 0)
	for _, v := range in {
		ok, vv, err := attrValueToJSON(isPropertyValue, v)
		if err != nil {
			return false, nil, err
		}

		if !ok {
			return false, nil, nil
		}

		l = append(l, json.RawMessage(vv))
	}

	return true, l, nil
}

func attrMapToJSON(isPropertyValue bool, in map[string]attr.Value) (bool, map[string]json.RawMessage, error) {
	m := map[string]json.RawMessage{}
	for k, v := range in {
		// preflight can only mock attributes under 'properties'.
		if k == "properties" {
			isPropertyValue = true
		}
		ok, vv, err := attrValueToJSON(isPropertyValue, v)
		if err != nil {
			return false, nil, err
		}

		if !ok {
			return false, nil, nil
		}

		m[k] = json.RawMessage(vv)
	}
	return true, m, nil
}

func attrValueToJSON(isPropertyValue bool, val attr.Value) (bool, []byte, error) {
	if val.IsNull() {
		val, err := json.Marshal(nil)
		return true, val, err
	}

	if val.IsUnknown() {
		if isPropertyValue {
			val, err := json.Marshal("[length('foo')]")
			return true, val, err
		} else {
			return false, nil, nil
		}

	}

	switch value := val.(type) {
	case types.Bool:
		val, err := json.Marshal(value.ValueBool())
		return true, val, err
	case types.String:
		val, err := json.Marshal(value.ValueString())
		return true, val, err
	case types.Int64:
		val, err := json.Marshal(value.ValueInt64())
		return true, val, err
	case types.Float64:
		val, err := json.Marshal(value.ValueFloat64())
		return true, val, err
	case types.Number:
		v, _ := value.ValueBigFloat().Float64()
		val, err := json.Marshal(v)
		return true, val, err
	case types.List:
		ok, l, err := attrListToJSON(isPropertyValue, value.Elements())
		if err != nil {
			return false, nil, err
		}

		if !ok {
			return false, nil, nil
		}

		val, err := json.Marshal(l)
		return true, val, err
	case types.Set:
		ok, l, err := attrListToJSON(isPropertyValue, value.Elements())
		if err != nil {
			return false, nil, err
		}

		if !ok {
			return false, nil, nil
		}

		val, err := json.Marshal(l)
		return true, val, err
	case types.Tuple:
		ok, l, err := attrListToJSON(isPropertyValue, value.Elements())
		if err != nil {
			return false, nil, err
		}

		if !ok {
			return false, nil, nil
		}

		val, err := json.Marshal(l)
		return true, val, err
	case types.Map:
		ok, m, err := attrMapToJSON(isPropertyValue, value.Elements())
		if err != nil {
			return false, nil, err
		}

		if !ok {
			return false, nil, nil
		}

		val, err := json.Marshal(m)
		return true, val, err
	case types.Object:
		ok, m, err := attrMapToJSON(isPropertyValue, value.Attributes())
		if err != nil {
			return false, nil, err
		}

		if !ok {
			return false, nil, nil
		}

		val, err := json.Marshal(m)
		return true, val, err
	default:
		return false, nil, fmt.Errorf("Unhandled type: %T", value)
	}
}

func ToJSON(d types.Dynamic) (bool, []byte, error) {
	return attrValueToJSON(false, d.UnderlyingValue())
}
