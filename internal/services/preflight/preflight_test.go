package preflight

import (
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestSuccessToJSON(t *testing.T) {
	input := types.DynamicValue(
		types.ObjectValueMust(
			map[string]attr.Type{
				"bool":         types.BoolType,
				"bool_null":    types.BoolType,
				"string":       types.StringType,
				"string_null":  types.StringType,
				"int64":        types.Int64Type,
				"int64_null":   types.Int64Type,
				"float64":      types.Float64Type,
				"float64_null": types.Float64Type,
				"number":       types.NumberType,
				"number_null":  types.NumberType,
				"list": types.ListType{
					ElemType: types.BoolType,
				},
				"list_empty": types.ListType{
					ElemType: types.BoolType,
				},
				"list_null": types.ListType{
					ElemType: types.BoolType,
				},
				"set": types.SetType{
					ElemType: types.BoolType,
				},
				"set_empty": types.SetType{
					ElemType: types.BoolType,
				},
				"set_null": types.SetType{
					ElemType: types.BoolType,
				},
				"tuple": types.TupleType{
					ElemTypes: []attr.Type{
						types.BoolType,
						types.StringType,
					},
				},
				"tuple_empty": types.TupleType{
					ElemTypes: []attr.Type{},
				},
				"tuple_null": types.TupleType{
					ElemTypes: []attr.Type{
						types.BoolType,
						types.StringType,
					},
				},
				"map": types.MapType{
					ElemType: types.BoolType,
				},
				"map_empty": types.MapType{
					ElemType: types.BoolType,
				},
				"map_null": types.MapType{
					ElemType: types.BoolType,
				},
				"object": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"bool":   types.BoolType,
						"string": types.StringType,
					},
				},
				"properties": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"bool_unknown":    types.BoolType,
						"string_unknown":  types.StringType,
						"int64_unknown":   types.Int64Type,
						"float64_unknown": types.Float64Type,
						"number_unknown":  types.NumberType,
						"list_unknown": types.ListType{
							ElemType: types.BoolType,
						},
						"set_unknown": types.SetType{
							ElemType: types.BoolType,
						},
						"tuple_unknown": types.TupleType{
							ElemTypes: []attr.Type{},
						},
						"map_unknown": types.MapType{
							ElemType: types.BoolType,
						},
						"object_unknown": types.ObjectType{
							AttrTypes: map[string]attr.Type{},
						},
					},
				},
				"object_empty": types.ObjectType{
					AttrTypes: map[string]attr.Type{},
				},
				"object_null": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"bool":   types.BoolType,
						"string": types.StringType,
					},
				},
			},
			map[string]attr.Value{
				"bool":         types.BoolValue(true),
				"bool_null":    types.BoolNull(),
				"string":       types.StringValue("a"),
				"string_null":  types.StringNull(),
				"int64":        types.Int64Value(123),
				"int64_null":   types.Int64Null(),
				"float64":      types.Float64Value(1.23),
				"float64_null": types.Float64Null(),
				"number":       types.NumberValue(big.NewFloat(1.23)),
				"number_null":  types.NumberNull(),
				"list": types.ListValueMust(
					types.BoolType,
					[]attr.Value{
						types.BoolValue(true),
						types.BoolValue(false),
					},
				),
				"list_empty": types.ListValueMust(types.BoolType, []attr.Value{}),
				"list_null":  types.ListNull(types.BoolType),
				"set": types.SetValueMust(
					types.BoolType,
					[]attr.Value{
						types.BoolValue(true),
						types.BoolValue(false),
					},
				),
				"set_empty": types.SetValueMust(types.BoolType, []attr.Value{}),
				"set_null":  types.SetNull(types.BoolType),
				"tuple": types.TupleValueMust(
					[]attr.Type{
						types.BoolType,
						types.StringType,
					},
					[]attr.Value{
						types.BoolValue(true),
						types.StringValue("a"),
					},
				),
				"tuple_empty": types.TupleValueMust(
					[]attr.Type{},
					[]attr.Value{},
				),
				"tuple_null": types.TupleNull(
					[]attr.Type{
						types.BoolType,
						types.StringType,
					},
				),
				"map": types.MapValueMust(
					types.BoolType,
					map[string]attr.Value{
						"a": types.BoolValue(true),
					},
				),
				"map_empty": types.MapValueMust(types.BoolType, map[string]attr.Value{}),
				"map_null":  types.MapNull(types.BoolType),
				"object": types.ObjectValueMust(
					map[string]attr.Type{
						"bool":   types.BoolType,
						"string": types.StringType,
					},
					map[string]attr.Value{
						"bool":   types.BoolValue(true),
						"string": types.StringValue("a"),
					},
				),
				"object_empty": types.ObjectValueMust(
					map[string]attr.Type{},
					map[string]attr.Value{},
				),
				"properties": types.ObjectValueMust(
					map[string]attr.Type{
						"bool_unknown":    types.BoolType,
						"string_unknown":  types.StringType,
						"int64_unknown":   types.Int64Type,
						"float64_unknown": types.Float64Type,
						"number_unknown":  types.NumberType,
						"list_unknown":    types.ListType{ElemType: types.BoolType},
						"set_unknown":     types.SetType{ElemType: types.BoolType},
						"tuple_unknown":   types.TupleType{ElemTypes: []attr.Type{}},
						"map_unknown":     types.MapType{ElemType: types.BoolType},
						"object_unknown":  types.ObjectType{AttrTypes: map[string]attr.Type{}},
					},
					map[string]attr.Value{
						"bool_unknown":    types.BoolUnknown(),
						"string_unknown":  types.StringUnknown(),
						"int64_unknown":   types.Int64Unknown(),
						"float64_unknown": types.Float64Unknown(),
						"number_unknown":  types.NumberUnknown(),
						"list_unknown":    types.ListUnknown(types.BoolType),
						"set_unknown":     types.SetUnknown(types.BoolType),
						"tuple_unknown":   types.TupleUnknown([]attr.Type{}),
						"map_unknown":     types.MapUnknown(types.BoolType),
						"object_unknown":  types.ObjectUnknown(map[string]attr.Type{}),
					},
				),
				"object_null": types.ObjectNull(
					map[string]attr.Type{
						"bool":   types.BoolType,
						"string": types.StringType,
					},
				),
			},
		),
	)

	expect := `
{
	"bool": true,
	"bool_null": null,
	"string": "a",
	"string_null": null,
	"int64": 123,
	"int64_null": null,
	"float64": 1.23,
	"float64_null": null,
	"number": 1.23,
	"number_null": null,
	"list": [true, false],
	"list_empty": [],
	"list_null": null,
	"set": [true, false],
	"set_empty": [],
	"set_null": null,
	"tuple": [true, "a"],
	"tuple_empty": [],
	"tuple_null": null,
	"map": {
		"a": true
	},
	"map_empty": {},
	"map_null": null,
	"object": {
		"bool": true,
		"string": "a"
	},
	"object_empty": {},
	"properties": {
		"bool_unknown": "[length('foo')]",
		"string_unknown": "[length('foo')]",
		"int64_unknown": "[length('foo')]",
		"float64_unknown": "[length('foo')]",
		"number_unknown": "[length('foo')]",
		"list_unknown": "[length('foo')]",
		"set_unknown": "[length('foo')]",
		"tuple_unknown": "[length('foo')]",
		"map_unknown": "[length('foo')]",
		"object_unknown": "[length('foo')]"
	},
	"object_null": null
}`

	mockResult, b, err := ToJSON(input)
	require.True(t, mockResult)
	require.NoError(t, err)
	require.JSONEq(t, expect, string(b))
}

func TestFailToJSON(t *testing.T) {
	input := types.DynamicValue(
		types.ObjectValueMust(
			map[string]attr.Type{
				"bool":         types.BoolType,
				"bool_null":    types.BoolType,
				"string":       types.StringType,
				"string_null":  types.StringType,
				"int64":        types.Int64Type,
				"int64_null":   types.Int64Type,
				"float64":      types.Float64Type,
				"float64_null": types.Float64Type,
				"number":       types.NumberType,
				"number_null":  types.NumberType,
				"list": types.ListType{
					ElemType: types.BoolType,
				},
				"list_empty": types.ListType{
					ElemType: types.BoolType,
				},
				"list_null": types.ListType{
					ElemType: types.BoolType,
				},
				"set": types.SetType{
					ElemType: types.BoolType,
				},
				"set_empty": types.SetType{
					ElemType: types.BoolType,
				},
				"set_null": types.SetType{
					ElemType: types.BoolType,
				},
				"tuple": types.TupleType{
					ElemTypes: []attr.Type{
						types.BoolType,
						types.StringType,
					},
				},
				"tuple_empty": types.TupleType{
					ElemTypes: []attr.Type{},
				},
				"tuple_null": types.TupleType{
					ElemTypes: []attr.Type{
						types.BoolType,
						types.StringType,
					},
				},
				"map": types.MapType{
					ElemType: types.BoolType,
				},
				"map_empty": types.MapType{
					ElemType: types.BoolType,
				},
				"map_null": types.MapType{
					ElemType: types.BoolType,
				},
				"object": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"bool":         types.BoolType,
						"bool_unknown": types.BoolType,
						"string":       types.StringType,
					},
				},
				"properties": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"bool_unknown":    types.BoolType,
						"string_unknown":  types.StringType,
						"int64_unknown":   types.Int64Type,
						"float64_unknown": types.Float64Type,
						"number_unknown":  types.NumberType,
						"list_unknown": types.ListType{
							ElemType: types.BoolType,
						},
						"set_unknown": types.SetType{
							ElemType: types.BoolType,
						},
						"tuple_unknown": types.TupleType{
							ElemTypes: []attr.Type{},
						},
						"map_unknown": types.MapType{
							ElemType: types.BoolType,
						},
						"object_unknown": types.ObjectType{
							AttrTypes: map[string]attr.Type{},
						},
					},
				},
				"object_empty": types.ObjectType{
					AttrTypes: map[string]attr.Type{},
				},
				"object_null": types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"bool":   types.BoolType,
						"string": types.StringType,
					},
				},
			},
			map[string]attr.Value{
				"bool":         types.BoolValue(true),
				"bool_null":    types.BoolNull(),
				"string":       types.StringValue("a"),
				"string_null":  types.StringNull(),
				"int64":        types.Int64Value(123),
				"int64_null":   types.Int64Null(),
				"float64":      types.Float64Value(1.23),
				"float64_null": types.Float64Null(),
				"number":       types.NumberValue(big.NewFloat(1.23)),
				"number_null":  types.NumberNull(),
				"list": types.ListValueMust(
					types.BoolType,
					[]attr.Value{
						types.BoolValue(true),
						types.BoolValue(false),
					},
				),
				"list_empty": types.ListValueMust(types.BoolType, []attr.Value{}),
				"list_null":  types.ListNull(types.BoolType),
				"set": types.SetValueMust(
					types.BoolType,
					[]attr.Value{
						types.BoolValue(true),
						types.BoolValue(false),
					},
				),
				"set_empty": types.SetValueMust(types.BoolType, []attr.Value{}),
				"set_null":  types.SetNull(types.BoolType),
				"tuple": types.TupleValueMust(
					[]attr.Type{
						types.BoolType,
						types.StringType,
					},
					[]attr.Value{
						types.BoolValue(true),
						types.StringValue("a"),
					},
				),
				"tuple_empty": types.TupleValueMust(
					[]attr.Type{},
					[]attr.Value{},
				),
				"tuple_null": types.TupleNull(
					[]attr.Type{
						types.BoolType,
						types.StringType,
					},
				),
				"map": types.MapValueMust(
					types.BoolType,
					map[string]attr.Value{
						"a": types.BoolValue(true),
					},
				),
				"map_empty": types.MapValueMust(types.BoolType, map[string]attr.Value{}),
				"map_null":  types.MapNull(types.BoolType),
				"object": types.ObjectValueMust(
					map[string]attr.Type{
						"bool":         types.BoolType,
						"bool_unknown": types.BoolType,
						"string":       types.StringType,
					},
					map[string]attr.Value{
						"bool":         types.BoolValue(true),
						"bool_unknown": types.BoolUnknown(),
						"string":       types.StringValue("a"),
					},
				),
				"object_empty": types.ObjectValueMust(
					map[string]attr.Type{},
					map[string]attr.Value{},
				),
				"properties": types.ObjectValueMust(
					map[string]attr.Type{
						"bool_unknown":    types.BoolType,
						"string_unknown":  types.StringType,
						"int64_unknown":   types.Int64Type,
						"float64_unknown": types.Float64Type,
						"number_unknown":  types.NumberType,
						"list_unknown":    types.ListType{ElemType: types.BoolType},
						"set_unknown":     types.SetType{ElemType: types.BoolType},
						"tuple_unknown":   types.TupleType{ElemTypes: []attr.Type{}},
						"map_unknown":     types.MapType{ElemType: types.BoolType},
						"object_unknown":  types.ObjectType{AttrTypes: map[string]attr.Type{}},
					},
					map[string]attr.Value{
						"bool_unknown":    types.BoolUnknown(),
						"string_unknown":  types.StringUnknown(),
						"int64_unknown":   types.Int64Unknown(),
						"float64_unknown": types.Float64Unknown(),
						"number_unknown":  types.NumberUnknown(),
						"list_unknown":    types.ListUnknown(types.BoolType),
						"set_unknown":     types.SetUnknown(types.BoolType),
						"tuple_unknown":   types.TupleUnknown([]attr.Type{}),
						"map_unknown":     types.MapUnknown(types.BoolType),
						"object_unknown":  types.ObjectUnknown(map[string]attr.Type{}),
					},
				),
				"object_null": types.ObjectNull(
					map[string]attr.Type{
						"bool":   types.BoolType,
						"string": types.StringType,
					},
				),
			},
		),
	)

	mockResult, _, err := ToJSON(input)
	require.False(t, mockResult)
	require.NoError(t, err)
}
