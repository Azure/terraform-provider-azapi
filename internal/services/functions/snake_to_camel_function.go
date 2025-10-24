package functions

import (
	"context"
	"strings"
	"unicode"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Snake2CamelFunction builds resource IDs for tenant scope
type Snake2CamelFunction struct{}

func (f *Snake2CamelFunction) Metadata(ctx context.Context, request function.MetadataRequest, response *function.MetadataResponse) {
	response.Name = "snake2camel"
}

func (f *Snake2CamelFunction) Definition(ctx context.Context, request function.DefinitionRequest, response *function.DefinitionResponse) {
	response.Definition = function.Definition{
		Parameters: []function.Parameter{
			function.DynamicParameter{
				AllowNullValue:      true,
				AllowUnknownValues:  true,
				Name:                "input",
				Description:         "The input object to convert from snake_case to camelCase.",
				MarkdownDescription: "The input object to convert from snake_case to camelCase.",
			},
		},
		Return:              function.DynamicReturn{},
		Summary:             "The camelCase returned object.",
		Description:         "Converts all keys in the input from snake_case to camelCase. Retains the original structure and values.",
		MarkdownDescription: "Converts all keys in the input from snake_case to camelCase. Retains the original structure and values.",
	}
}

func (f *Snake2CamelFunction) Run(ctx context.Context, request function.RunRequest, response *function.RunResponse) {
	var inputObj types.Dynamic
	if response.Error = request.Arguments.Get(ctx, &inputObj); response.Error != nil {
		return
	}
	if inputObj.IsNull() {
		response.Result = function.NewResultData(types.DynamicNull())
		return
	}

	if inputObj.IsUnknown() {
		response.Result = function.NewResultData(types.DynamicUnknown())
		return
	}

	result := dynamicSnake2Camel(ctx, inputObj)
	response.Error = response.Result.Set(ctx, result)
}

var _ function.Function = &Snake2CamelFunction{}

func dynamicSnake2Camel(ctx context.Context, input attr.Value) attr.Value {
	switch t := input.(type) {
	case types.Dynamic:
		val := t.UnderlyingValue()
		newVal := dynamicSnake2Camel(ctx, val)
		return types.DynamicValue(newVal) // to ensure val is attr.Value
	case types.Object:
		newObjectType := convertAttrType(t.Type(ctx)).(types.ObjectType)

		if t.IsUnknown() {
			return types.ObjectUnknown(newObjectType.AttrTypes)
		}

		if t.IsNull() {
			return types.ObjectNull(newObjectType.AttrTypes)
		}

		newAttrValues := make(map[string]attr.Value)
		for key, attrValue := range t.Attributes() {
			newKey := snake2Camel(key)
			newAttrValues[newKey] = dynamicSnake2Camel(ctx, attrValue)
		}
		return types.ObjectValueMust(newObjectType.AttrTypes, newAttrValues)
	case types.Map:
		newMapType := convertAttrType(t.Type(ctx)).(types.MapType)

		if t.IsUnknown() {
			return types.MapUnknown(newMapType.ElemType)
		}

		if t.IsNull() {
			return types.MapNull(newMapType.ElemType)
		}

		newAttrValues := make(map[string]attr.Value)
		for key, attrValue := range t.Elements() {
			newKey := snake2Camel(key)
			newAttrValues[newKey] = dynamicSnake2Camel(ctx, attrValue)
		}
		return types.MapValueMust(newMapType.ElemType, newAttrValues)
	case types.List:
		newElemType := convertAttrType(t.Type(ctx)).(types.ListType)

		if t.IsUnknown() {
			return types.ListUnknown(newElemType.ElemType)
		}

		if t.IsNull() {
			return types.ListNull(newElemType.ElemType)
		}
		newElemValues := make([]attr.Value, len(t.Elements()))
		for i, elemValue := range t.Elements() {
			newElemValues[i] = dynamicSnake2Camel(ctx, elemValue)
		}
		return types.ListValueMust(newElemType.ElemType, newElemValues)
	case types.Set:
		newElemType := convertAttrType(t.Type(ctx)).(types.SetType)

		if t.IsUnknown() {
			return types.SetUnknown(newElemType.ElemType)
		}

		if t.IsNull() {
			return types.SetNull(newElemType.ElemType)
		}
		newElemValues := make([]attr.Value, len(t.Elements()))
		for i, elemValue := range t.Elements() {
			newElemValues[i] = dynamicSnake2Camel(ctx, elemValue)
		}
		return types.SetValueMust(newElemType.ElemType, newElemValues)
	case types.Tuple:
		newType := convertAttrType(t.Type(ctx)).(types.TupleType)

		if t.IsUnknown() {
			return types.TupleUnknown(newType.ElementTypes())
		}

		if t.IsNull() {
			return types.TupleNull(newType.ElementTypes())
		}

		newElemValues := make([]attr.Value, len(t.Elements()))
		for i, elemValue := range t.Elements() {
			newElemValues[i] = dynamicSnake2Camel(ctx, elemValue)
		}
		return types.TupleValueMust(newType.ElementTypes(), newElemValues)
	}
	return input
}

func convertAttrType(attrType attr.Type) attr.Type {
	switch t := attrType.(type) {
	case types.ObjectType:
		newAttrTypes := make(map[string]attr.Type)
		for key, innerType := range t.AttrTypes {
			newKey := snake2Camel(key)
			newAttrTypes[newKey] = convertAttrType(innerType)
		}
		return types.ObjectType{AttrTypes: newAttrTypes}
	case types.MapType:
		newElemType := convertAttrType(t.ElemType)
		return types.MapType{ElemType: newElemType}
	case types.ListType:
		newElemType := convertAttrType(t.ElemType)
		return types.ListType{ElemType: newElemType}
	case types.SetType:
		newElemType := convertAttrType(t.ElemType)
		return types.SetType{ElemType: newElemType}
	case types.TupleType:
		newElemTypes := make([]attr.Type, len(t.ElementTypes()))
		for i, elemType := range t.ElementTypes() {
			newElemTypes[i] = convertAttrType(elemType)
		}
		return types.TupleType{ElemTypes: newElemTypes}
	}
	return attrType
}

// snake2Camel converts snake_case to camelCase
func snake2Camel(input string) string {
	sb := strings.Builder{}
	upperNext := false
	for _, r := range input {
		if r == '_' {
			upperNext = true
			continue
		}
		if upperNext {
			sb.WriteRune(unicode.ToUpper(r))
			upperNext = false
			continue
		}
		sb.WriteRune(unicode.ToLower(r))
	}
	return sb.String()
}
