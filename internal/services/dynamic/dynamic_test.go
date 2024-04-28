package dynamic

import (
	"context"
	"math/big"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/require"
)

func TestToJSON(t *testing.T) {
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
	"object_null": null
}`

	b, err := ToJSON(input)
	require.NoError(t, err)
	require.JSONEq(t, expect, string(b))
}

func TestFromJSON(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect types.Dynamic
	}{
		{
			name: "basic",
			input: `
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
	"number_null":    null,
	"list": [true, false],
	"list_empty": [],
	"list_null": null,
	"set": [true, false],
	"set_empty": [],
	"set_null": null,
	"tuple": [true, "a"],
	"tuple_empty": [],
	"tuple_null":  null,
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
	"object_null": null
}`,
			expect: types.DynamicValue(
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
						"object_null": types.ObjectNull(
							map[string]attr.Type{
								"bool":   types.BoolType,
								"string": types.StringType,
							},
						),
					},
				),
			),
		},
		{
			name: "fields not defined in type is ignored",
			input: `
{
	"str1": "a",
	"str2": "b"
}
`,
			expect: types.DynamicValue(
				types.ObjectValueMust(
					map[string]attr.Type{
						"str1": types.StringType,
					},
					map[string]attr.Value{
						"str1": types.StringValue("a"),
					},
				),
			),
		},
		{
			name: "fields defined in type not in JSON, set it as null",
			input: `
{
	"str1": "a"
}
`,
			expect: types.DynamicValue(
				types.ObjectValueMust(
					map[string]attr.Type{
						"str1": types.StringType,
						"str2": types.StringType,
					},
					map[string]attr.Value{
						"str1": types.StringValue("a"),
						"str2": types.StringNull(),
					},
				),
			),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := FromJSON([]byte(tt.input), tt.expect.UnderlyingValue().Type(context.TODO()))
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
		})
	}
}

func TestFromJSONImplied(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		expect types.Dynamic
	}{
		{
			name: "basic",
			input: `
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
	"number_null":    null,
	"list": [true, false],
	"list_null": null,
	"set": [true, false],
	"set_null": null,
	"tuple": [true, "a"],
	"tuple_null":  null,
	"map": {
		"a": true
	},
	"map_null": null,
	"object": {
		"bool": true,
		"string": "a"
	},
	"object_null": null
}`,
			expect: types.DynamicValue(
				types.ObjectValueMust(
					map[string]attr.Type{
						"bool":         types.BoolType,
						"bool_null":    types.DynamicType,
						"string":       types.StringType,
						"string_null":  types.DynamicType,
						"int64":        types.NumberType,
						"int64_null":   types.DynamicType,
						"float64":      types.NumberType,
						"float64_null": types.DynamicType,
						"number":       types.NumberType,
						"number_null":  types.DynamicType,
						"list": types.TupleType{
							ElemTypes: []attr.Type{
								types.BoolType,
								types.BoolType,
							},
						},
						"list_null": types.DynamicType,
						"set": types.TupleType{
							ElemTypes: []attr.Type{
								types.BoolType,
								types.BoolType,
							},
						},
						"set_null": types.DynamicType,
						"tuple": types.TupleType{
							ElemTypes: []attr.Type{
								types.BoolType,
								types.StringType,
							},
						},
						"tuple_null": types.DynamicType,
						"map": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"a": types.BoolType,
							},
						},
						"map_null": types.DynamicType,
						"object": types.ObjectType{
							AttrTypes: map[string]attr.Type{
								"bool":   types.BoolType,
								"string": types.StringType,
							},
						},
						"object_null": types.DynamicType,
					},
					map[string]attr.Value{
						"bool":         types.BoolValue(true),
						"bool_null":    types.DynamicNull(),
						"string":       types.StringValue("a"),
						"string_null":  types.DynamicNull(),
						"int64":        types.NumberValue(big.NewFloat(123)),
						"int64_null":   types.DynamicNull(),
						"float64":      types.NumberValue(big.NewFloat(1.23)),
						"float64_null": types.DynamicNull(),
						"number":       types.NumberValue(big.NewFloat(1.23)),
						"number_null":  types.DynamicNull(),
						"list": types.TupleValueMust(
							[]attr.Type{
								types.BoolType,
								types.BoolType,
							},
							[]attr.Value{
								types.BoolValue(true),
								types.BoolValue(false),
							},
						),
						"list_null": types.DynamicNull(),
						"set": types.TupleValueMust(
							[]attr.Type{
								types.BoolType,
								types.BoolType,
							},
							[]attr.Value{
								types.BoolValue(true),
								types.BoolValue(false),
							},
						),
						"set_null": types.DynamicNull(),
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
						"tuple_null": types.DynamicNull(),
						"map": types.ObjectValueMust(
							map[string]attr.Type{
								"a": types.BoolType,
							},
							map[string]attr.Value{
								"a": types.BoolValue(true),
							},
						),
						"map_null": types.DynamicNull(),
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
						"object_null": types.DynamicNull(),
					},
				),
			),
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := FromJSONImplied([]byte(tt.input))
			require.NoError(t, err)
			require.Equal(t, tt.expect, actual)
		})
	}
}
