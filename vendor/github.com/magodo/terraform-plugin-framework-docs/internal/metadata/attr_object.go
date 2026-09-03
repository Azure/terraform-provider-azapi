package metadata

import (
	"context"
	"maps"
	"slices"
	"strings"

	fwattr "github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type Objects map[string]Object

type Object struct {
	key    string
	fields map[string]ObjectField
}

type ObjectField struct {
	parents     []string
	name        string
	dataType    DataType
	description string
	isObject    bool
}

func newObjects(ctx context.Context, parents []string, attrs map[string]fwattr.Type) (objects Objects, diags diag.Diagnostics) {
	objects = Objects{}
	fields := map[string]ObjectField{}
	for name, attr := range attrs {
		field := ObjectField{
			parents:     parents,
			name:        name,
			dataType:    DataType{inner: attr},
			description: PointerTo(MaybeDescriptionCtxOf(ctx, attr)),
		}

		// If this is an object nested in an object, recursively process it.
		if obj, ok := attr.(basetypes.ObjectTypable); ok {
			field.isObject = true

			if obj, ok := obj.(fwattr.TypeWithAttributeTypes); ok {
				nestedObjects, odiags := newObjects(ctx, slices.Concat(parents, []string{name}), obj.AttributeTypes())
				diags = append(diags, odiags...)
				if diags.HasError() {
					return nil, diags
				}
				maps.Copy(objects, nestedObjects)
			}
		}

		fields[name] = field
	}
	key := strings.Join(parents, ".")
	objects[key] = Object{
		key:    key,
		fields: fields,
	}

	return objects, diags
}

func (objs Objects) ToFunctionObjects() FunctionObjects {
	out := FunctionObjects{}
	for key, obj := range objs {
		fields := map[string]FunctionField{}
		for key, field := range obj.fields {
			fields[key] = field.ToFunctionField()
		}
		out[key] = FunctionObject{
			functionKey: obj.key,
			fields:      fields,
		}
	}
	return out
}

func (obj ObjectField) ToFunctionField() FunctionField {
	return FunctionField{
		parents:     obj.parents,
		name:        obj.name,
		dataType:    obj.dataType,
		description: obj.description,
		isObject:    obj.isObject,
	}
}

func (objs Objects) ToNestedFields(rootField Field) NestedFields {
	out := NestedFields{}
	for key, obj := range objs {
		fields := Fields{}
		for key, field := range obj.fields {
			fields[key] = field.ToField(rootField)
		}
		out[key] = NestedField{
			fields: fields,
		}
	}
	return out
}

func (obj ObjectField) ToField(rootField Field) Field {
	return Field{
		parents:     obj.parents,
		name:        obj.name,
		dataType:    obj.dataType,
		required:    rootField.required,
		optional:    rootField.optional,
		computed:    rootField.computed,
		description: obj.description,
		isObject:    obj.isObject,
	}
}
