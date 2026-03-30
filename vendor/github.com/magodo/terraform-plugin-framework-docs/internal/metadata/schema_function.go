package metadata

import (
	"context"
	"fmt"
	"maps"
	"slices"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

type FunctionSchema struct {
	Description string
	Summary     string
	Deprecation string

	Parameters FunctionFields
	Objects    FunctionObjects

	Return        FunctionField
	ReturnObjects FunctionObjects
}

func NewFunctionSchema(ctx context.Context, sch function.Definition) (schema FunctionSchema, diags diag.Diagnostics) {
	nested := FunctionObjects{}

	paramFields, paramNested, odiags := newFunctionParameters(ctx, nil, sch.Parameters)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(nested, paramNested)

	if sch.VariadicParameter != nil {
		varFields, paramNested, odiags := newFunctionParameters(ctx, nil, []function.Parameter{sch.VariadicParameter})
		diags.Append(odiags...)
		if diags.HasError() {
			return
		}
		field := varFields[0]
		field.isVariadic = true

		paramFields = append(paramFields, field)
		maps.Copy(nested, paramNested)
	}

	retField, retNested, odiags := newFunctionReturn(ctx, sch.Return)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}

	schema = FunctionSchema{
		Summary: sch.Summary,
		Description: func() string {
			v := sch.MarkdownDescription
			if v != "" {
				return v
			}
			return sch.Description
		}(),
		Deprecation:   sch.DeprecationMessage,
		Parameters:    paramFields,
		Objects:       nested,
		Return:        retField,
		ReturnObjects: retNested,
	}
	return
}

func newFunctionReturn(ctx context.Context, ret function.Return) (field FunctionField, objects FunctionObjects, diags diag.Diagnostics) {
	field = FunctionField{
		dataType:              DataType{inner: ret.GetType()},
		customTypeDescription: PointerTo(MaybeDescriptionCtxOf(ctx, ret.GetType())),
	}

	if obj, ok := ret.(function.ObjectReturn); ok {
		field.isObject = true
		plainObjs, odiags := newObjects(ctx, nil, obj.AttributeTypes)
		diags = append(diags, odiags...)
		if diags.HasError() {
			return
		}
		objects = plainObjs.ToFunctionObjects()
	}

	return field, objects, diags
}

func newFunctionParameters(ctx context.Context, parents []string, params []function.Parameter) (fields FunctionFields, nested FunctionObjects, diags diag.Diagnostics) {
	fields = FunctionFields{}
	nested = FunctionObjects{}

	for _, attr := range params {
		field := FunctionField{
			parents:               parents,
			name:                  attr.GetName(),
			dataType:              DataType{inner: attr.GetType()},
			description:           DescriptionOf(attr),
			customTypeDescription: PointerTo(MaybeDescriptionCtxOf(ctx, attr.GetType())),
			allowNull:             attr.GetAllowNullValue(),
			allowUnknown:          attr.GetAllowUnknownValues(),
		}

		switch attr := attr.(type) {
		case function.BoolParameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.BoolParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.Float32Parameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.Float32ParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.Float64Parameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.Float64ParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.Int32Parameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.Int32ParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.Int64Parameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.Int64ParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.NumberParameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.NumberParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.StringParameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.StringParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.ListParameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.ListParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.MapParameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.MapParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.SetParameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.SetParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.DynamicParameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.DynamicParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
		case function.ObjectParameter:
			field.validators = MapSliceSome(attr.Validators,
				func(v function.ObjectParameterValidator) *string {
					return MaybeDescriptionCtxOf(ctx, v)
				})
			field.isObject = true
			nestedObjects, odiags := newObjects(ctx, slices.Concat(parents, []string{attr.GetName()}), attr.AttributeTypes)
			diags = append(diags, odiags...)
			if diags.HasError() {
				return
			}
			maps.Copy(nested, nestedObjects.ToFunctionObjects())

		default:
			diags.AddError("unknown schema type", fmt.Sprintf("%T", attr))
			return
		}

		fields = append(fields, field)
	}

	return
}
