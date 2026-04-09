package metadata

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/identityschema"
)

type ResourceIdentitySchema struct {
	Fields ResourceIdentityFields
}

func NewResourceIdentitySchema(ctx context.Context, sch identityschema.Schema) (schema ResourceIdentitySchema, diags diag.Diagnostics) {
	fields, odiags := newResourceIdentityAttrFields(ctx, sch.Attributes)
	diags.Append(odiags...)
	if diags.HasError() {
		return
	}
	schema = ResourceIdentitySchema{
		Fields: fields,
	}
	return
}

func newResourceIdentityAttrFields(ctx context.Context, attrs map[string]identityschema.Attribute) (fields ResourceIdentityFields, diags diag.Diagnostics) {
	fields = ResourceIdentityFields{}

	for name, attr := range attrs {
		field := ResourceIdentityField{
			Name:                  name,
			DataType:              DataType{inner: attr.GetType()},
			Required:              attr.IsRequiredForImport(),
			Optional:              attr.IsOptionalForImport(),
			Description:           DescriptionOf(attr),
			customTypeDescription: PointerTo(MaybeDescriptionCtxOf(ctx, attr.GetType())),
		}
		fields[name] = field
	}

	return
}
