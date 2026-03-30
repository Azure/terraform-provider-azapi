package metadata

import (
	"context"
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

type ResourceSchema struct {
	Description string
	Deprecation string

	Fields Fields

	// Including nested attribute object or block object.
	Nested NestedFields
}

func NewResourceSchema(ctx context.Context, sch schema.Schema) (schema ResourceSchema, diags diag.Diagnostics) {
	fields := Fields{}
	nested := NestedFields{}

	attrFields, attrNested, odiags := newResourceAttrFields(ctx, nil, sch.Attributes)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, attrFields)
	maps.Copy(nested, attrNested)

	blockFields, blockNested, odiags := newResourceBlockFields(ctx, nil, sch.Blocks)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	maps.Copy(fields, blockFields)
	maps.Copy(nested, blockNested)

	schema = ResourceSchema{
		Description: DescriptionOf(sch),
		Deprecation: sch.GetDeprecationMessage(),
		Fields:      fields,
		Nested:      nested,
	}
	return
}

func newResourceAttrFields(ctx context.Context, parents []string, attrs map[string]schema.Attribute) (fields Fields, nested NestedFields, diags diag.Diagnostics) {
	fields = Fields{}
	nested = NestedFields{}

	for name, attr := range attrs {
		var (
			objectNested NestedFields
			objectDiags  diag.Diagnostics
		)

		field := Field{
			parents:               parents,
			name:                  name,
			dataType:              DataType{inner: attr.GetType()},
			required:              attr.IsRequired(),
			optional:              attr.IsOptional(),
			computed:              attr.IsComputed(),
			sensitive:             attr.IsSensitive(),
			description:           DescriptionOf(attr),
			deprecation:           attr.GetDeprecationMessage(),
			customTypeDescription: PointerTo(MaybeDescriptionCtxOf(ctx, attr.GetType())),
			writeOnly:             attr.IsWriteOnly(),
		}

		switch attr := attr.(type) {
		case schema.BoolAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Bool) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Bool) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Bool) string { return DescriptionCtxOf(ctx, v) })

		case schema.Float32Attribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Float32) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Float32) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Float32) string { return DescriptionCtxOf(ctx, v) })
		case schema.Float64Attribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Float64) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Float64) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Float64) string { return DescriptionCtxOf(ctx, v) })
		case schema.Int32Attribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Int32) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Int32) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Int32) string { return DescriptionCtxOf(ctx, v) })
		case schema.Int64Attribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Int64) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Int64) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Int64) string { return DescriptionCtxOf(ctx, v) })
		case schema.NumberAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Number) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Number) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Number) string { return DescriptionCtxOf(ctx, v) })
		case schema.StringAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.String) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.String) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.String) string { return DescriptionCtxOf(ctx, v) })
		case schema.ListAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.List) string { return DescriptionCtxOf(ctx, v) })
		case schema.MapAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Map) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Map) string { return DescriptionCtxOf(ctx, v) })
		case schema.SetAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Set) string { return DescriptionCtxOf(ctx, v) })
		case schema.DynamicAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Dynamic) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Dynamic) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Dynamic) string { return DescriptionCtxOf(ctx, v) })
		case schema.ObjectAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Object) string { return DescriptionCtxOf(ctx, v) })

			field.isObject = true
			objects, objectDiags := newObjects(ctx, slices.Concat(parents, []string{name}), attr.AttributeTypes)
			if objectDiags.HasError() {
				return nil, nil, objectDiags
			}
			objectNested = objects.ToNestedFields(field)
		case schema.SingleNestedAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Object) string { return DescriptionCtxOf(ctx, v) })

			field.isObject = true

			// Nullify the validators of the single nested block since it is handled at the block level, which avoids to repeat in the nested schema level.
			nestedObj := attr.GetNestedObject().(schema.NestedAttributeObject)
			nestedObj.PlanModifiers = nil
			nestedObj.Validators = nil

			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), nestedObj)
		case schema.SetNestedAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Set) string { return DescriptionCtxOf(ctx, v) })

			field.isObject = true
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.MapNestedAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.Map) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.Map) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.Map) string { return DescriptionCtxOf(ctx, v) })

			field.isObject = true
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		case schema.ListNestedAttribute:
			field.planModifiers = MapSlice(attr.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(attr.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) })
			field.defaultDesc = MapOrZero(attr.Default, func(v defaults.List) string { return DescriptionCtxOf(ctx, v) })

			field.isObject = true
			objectNested, objectDiags = newResourceNestedAttrObjectFields(ctx, slices.Concat(parents, []string{name}), attr.GetNestedObject().(schema.NestedAttributeObject))
		default:
			diags.AddError("unknown schema type", fmt.Sprintf("%T", attr))
			return
		}

		fields[name] = field

		diags = append(diags, objectDiags...)
		if diags.HasError() {
			return
		}
		maps.Copy(nested, objectNested)
	}

	return
}

func newResourceNestedAttrObjectFields(ctx context.Context, parents []string, obj schema.NestedAttributeObject) (nested NestedFields, diags diag.Diagnostics) {
	nested = NestedFields{}

	attrFields, attrNested, attrDiags := newResourceAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	nested[strings.Join(parents, ".")] = NestedField{
		planModifiers: MapSlice(obj.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
		validators:    MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		fields:        attrFields,
	}
	maps.Copy(nested, attrNested)
	return
}

func newResourceBlockFields(ctx context.Context, parents []string, blks map[string]schema.Block) (fields Fields, nested NestedFields, diags diag.Diagnostics) {
	fields = Fields{}
	nested = NestedFields{}

	for name, blk := range blks {
		field := Field{
			parents:               parents,
			name:                  name,
			dataType:              DataType{isblk: true, inner: blk.Type()},
			optional:              true, // Always regard a block as optional.
			description:           DescriptionOf(blk),
			deprecation:           blk.GetDeprecationMessage(),
			customTypeDescription: PointerTo(MaybeDescriptionCtxOf(ctx, blk.Type())),
			isObject:              true,
		}

		switch blk := blk.(type) {
		case schema.SingleNestedBlock:
			field.planModifiers = MapSlice(blk.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(blk.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) })
		case schema.ListNestedBlock:
			field.planModifiers = MapSlice(blk.PlanModifiers, func(v planmodifier.List) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(blk.Validators, func(v validator.List) string { return DescriptionCtxOf(ctx, v) })
		case schema.SetNestedBlock:
			field.planModifiers = MapSlice(blk.PlanModifiers, func(v planmodifier.Set) string { return DescriptionCtxOf(ctx, v) })
			field.validators = MapSlice(blk.Validators, func(v validator.Set) string { return DescriptionCtxOf(ctx, v) })
		default:
			diags.AddError("unknown schema type", fmt.Sprintf("%T", blk))
			return
		}

		// Nullify the validators of the single nested block since it is handled at the block level, which avoids to repeat in the nested schema level.
		nestedObj := blk.GetNestedObject().(schema.NestedBlockObject)
		if _, ok := blk.(schema.SingleNestedBlock); ok {
			nestedObj.Validators = nil
			nestedObj.PlanModifiers = nil
		}
		objectNested, odiags := newResourceNestedBlkObjectFields(ctx, slices.Concat(parents, []string{name}), nestedObj)
		diags = append(diags, odiags...)
		if diags.HasError() {
			return
		}

		fields[name] = field
		maps.Copy(nested, objectNested)
	}

	return
}

func newResourceNestedBlkObjectFields(ctx context.Context, parents []string, obj schema.NestedBlockObject) (nested NestedFields, diags diag.Diagnostics) {
	attrFields, attrNested, attrDiags := newResourceAttrFields(ctx, parents, obj.Attributes)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	blkFields, blkNested, attrDiags := newResourceBlockFields(ctx, parents, obj.Blocks)
	diags.Append(attrDiags...)
	if diags.HasError() {
		return
	}

	fields := Fields{}
	maps.Copy(fields, attrFields)
	maps.Copy(fields, blkFields)

	nested = NestedFields{}
	maps.Copy(nested, attrNested)
	maps.Copy(nested, blkNested)

	nested[strings.Join(parents, ".")] = NestedField{
		planModifiers: MapSlice(obj.PlanModifiers, func(v planmodifier.Object) string { return DescriptionCtxOf(ctx, v) }),
		validators:    MapSlice(obj.Validators, func(v validator.Object) string { return DescriptionCtxOf(ctx, v) }),
		fields:        fields,
	}
	return
}
