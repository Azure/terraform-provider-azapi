package metadata

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

type Metadata struct {
	ProviderName string
	Provider     ProviderMetadata
	Resources    ResourceMetadatas
	DataSources  DataSourceMetadatas
	Ephemerals   EphemeralMetadatas
	Actions      ActionMetadatas
	Lists        ListMetadatas
	Functions    FunctionMetadatas
}

type ProviderMetadata struct {
	Schema ProviderSchema

	RenderOption *ProviderRenderOption
}

type ResourceMetadatas map[string]ResourceMetadata

type ResourceMetadata struct {
	Schema   ResourceSchema
	Identity *ResourceIdentitySchema

	RenderOption *ResourceRenderOption
}

type DataSourceMetadatas map[string]DataSourceMetadata

type DataSourceMetadata struct {
	Schema DataSourceSchema

	RenderOption *DataSourceRenderOption
}

type EphemeralMetadatas map[string]EphemeralMetadata

type EphemeralMetadata struct {
	Schema EphemeralSchema

	RenderOption *EphemeralRenderOption
}

type ActionMetadatas map[string]ActionMetadata

type ActionMetadata struct {
	Schema ActionSchema

	RenderOption *ActionRenderOption
}

type ListMetadatas map[string]ListMetadata

type ListMetadata struct {
	Schema ListSchema

	RenderOption *ListRenderOption
}

type FunctionMetadatas map[string]FunctionMetadata

type FunctionMetadata struct {
	Schema FunctionSchema

	RenderOption *FunctionRenderOption
}

// Interface definitions mirrored from the interfaces.go at the root level.

type ProviderWithRenderOption interface {
	provider.Provider
	RenderOption() ProviderRenderOption
}

type ResourceWithRenderOption interface {
	resource.Resource
	RenderOption() ResourceRenderOption
}

type DataSourceWithRenderOption interface {
	datasource.DataSource
	RenderOption() DataSourceRenderOption
}

type EphemeralResourceWithRenderOption interface {
	ephemeral.EphemeralResource
	RenderOption() EphemeralRenderOption
}

type ListResourceWithRenderOption interface {
	list.ListResource
	RenderOption() ListRenderOption
}

type ActionWithRenderOption interface {
	action.Action
	RenderOption() ActionRenderOption
}

type FunctionWithRenderOption interface {
	function.Function
	RenderOption() FunctionRenderOption
}

func GetMetadata(ctx context.Context, p provider.Provider) (metadata Metadata, diags diag.Diagnostics) {
	var providerMetadataResp provider.MetadataResponse
	p.Metadata(ctx, provider.MetadataRequest{}, &providerMetadataResp)

	var schemaResp provider.SchemaResponse
	p.Schema(ctx, provider.SchemaRequest{}, &schemaResp)
	diags.Append(schemaResp.Diagnostics...)
	if diags.HasError() {
		diags = WrapErrDiags(diags, "provider")
		return
	}
	providerSchema, odiags := NewProviderSchema(ctx, schemaResp.Schema)
	diags.Append(odiags...)
	if diags.HasError() {
		diags = WrapErrDiags(diags, "provider")
		return
	}

	providerMetadata := ProviderMetadata{
		Schema: providerSchema,
	}
	if p, ok := p.(ProviderWithRenderOption); ok {
		providerMetadata.RenderOption = new(p.RenderOption())
	}

	metadata = Metadata{
		ProviderName: providerMetadataResp.TypeName,
		Provider:     providerMetadata,
		Resources:    ResourceMetadatas{},
		DataSources:  DataSourceMetadatas{},
		Ephemerals:   EphemeralMetadatas{},
		Actions:      ActionMetadatas{},
		Lists:        ListMetadatas{},
		Functions:    FunctionMetadatas{},
	}

	for _, builder := range p.Resources(ctx) {
		res := builder()

		// Get the resource type
		var metadataResp resource.MetadataResponse
		res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
		resourceType := metadataResp.TypeName

		var schemaResp resource.SchemaResponse
		res.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
		diags.Append(schemaResp.Diagnostics...)
		if diags.HasError() {
			diags = WrapErrDiags(diags, fmt.Sprintf("resource %q", resourceType))
			return
		}

		sch, odiags := NewResourceSchema(ctx, schemaResp.Schema)
		diags.Append(odiags...)
		if diags.HasError() {
			diags = WrapErrDiags(diags, fmt.Sprintf("resource %q", resourceType))
			return
		}

		resMetadata := ResourceMetadata{
			Schema: sch,
		}

		if res, ok := res.(ResourceWithRenderOption); ok {
			resMetadata.RenderOption = new(res.RenderOption())
		}

		if resourceWithIdentity, ok := res.(resource.ResourceWithIdentity); ok {
			var schemaResp resource.IdentitySchemaResponse
			resourceWithIdentity.IdentitySchema(ctx, resource.IdentitySchemaRequest{}, &schemaResp)
			diags.Append(schemaResp.Diagnostics...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("resource identity %q", resourceType))
				return
			}

			sch, odiags := NewResourceIdentitySchema(ctx, schemaResp.IdentitySchema)
			diags.Append(odiags...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("resource identity %q", resourceType))
				return
			}

			resMetadata.Identity = &sch
		}

		metadata.Resources[resourceType] = resMetadata
	}

	for _, builder := range p.DataSources(ctx) {
		res := builder()

		// Get the resource type
		var metadataResp datasource.MetadataResponse
		res.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
		resourceType := metadataResp.TypeName

		var schemaResp datasource.SchemaResponse
		res.Schema(ctx, datasource.SchemaRequest{}, &schemaResp)
		diags.Append(schemaResp.Diagnostics...)
		if diags.HasError() {
			diags = WrapErrDiags(diags, fmt.Sprintf("data source %q", resourceType))
			return
		}

		sch, odiags := NewDataSourceSchema(ctx, schemaResp.Schema)
		diags.Append(odiags...)
		if diags.HasError() {
			diags = WrapErrDiags(diags, fmt.Sprintf("data source %q", resourceType))
			return
		}

		resMetadata := DataSourceMetadata{
			Schema: sch,
		}

		if res, ok := res.(DataSourceWithRenderOption); ok {
			resMetadata.RenderOption = new(res.RenderOption())
		}

		metadata.DataSources[resourceType] = resMetadata
	}

	if p, ok := p.(provider.ProviderWithEphemeralResources); ok {
		for _, builder := range p.EphemeralResources(ctx) {
			res := builder()

			// Get the resource type
			var metadataResp ephemeral.MetadataResponse
			res.Metadata(ctx, ephemeral.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
			resourceType := metadataResp.TypeName

			var schemaResp ephemeral.SchemaResponse
			res.Schema(ctx, ephemeral.SchemaRequest{}, &schemaResp)
			diags.Append(schemaResp.Diagnostics...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("ephemeral resource %q", resourceType))
				return
			}

			sch, odiags := NewEphemeralSchema(ctx, schemaResp.Schema)
			diags.Append(odiags...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("ephemeral resource %q", resourceType))
				return
			}

			resMetadata := EphemeralMetadata{
				Schema: sch,
			}

			if res, ok := res.(EphemeralResourceWithRenderOption); ok {
				resMetadata.RenderOption = new(res.RenderOption())
			}

			metadata.Ephemerals[resourceType] = resMetadata
		}
	}

	if p, ok := p.(provider.ProviderWithActions); ok {
		for _, builder := range p.Actions(ctx) {
			res := builder()

			// Get the resource type
			var metadataResp action.MetadataResponse
			res.Metadata(ctx, action.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
			resourceType := metadataResp.TypeName

			var schemaResp action.SchemaResponse
			res.Schema(ctx, action.SchemaRequest{}, &schemaResp)
			diags.Append(schemaResp.Diagnostics...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("action %q", resourceType))
				return
			}

			sch, odiags := NewActionSchema(ctx, schemaResp.Schema)
			diags.Append(odiags...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("action %q", resourceType))
				return
			}

			resMetadata := ActionMetadata{
				Schema: sch,
			}

			if res, ok := res.(ActionWithRenderOption); ok {
				resMetadata.RenderOption = new(res.RenderOption())
			}

			metadata.Actions[resourceType] = resMetadata
		}
	}

	if p, ok := p.(provider.ProviderWithListResources); ok {
		for _, builder := range p.ListResources(ctx) {
			res := builder()

			// Get the resource type
			var metadataResp resource.MetadataResponse
			res.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: providerMetadataResp.TypeName}, &metadataResp)
			resourceType := metadataResp.TypeName

			var schemaResp list.ListResourceSchemaResponse
			res.ListResourceConfigSchema(ctx, list.ListResourceSchemaRequest{}, &schemaResp)
			diags.Append(schemaResp.Diagnostics...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("list resource %q", resourceType))
				return
			}

			sch, odiags := NewListSchema(ctx, schemaResp.Schema)
			diags.Append(odiags...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("list resource %q", resourceType))
				return
			}

			resMetadata := ListMetadata{
				Schema: sch,
			}

			if res, ok := res.(ListResourceWithRenderOption); ok {
				resMetadata.RenderOption = new(res.RenderOption())
			}

			metadata.Lists[resourceType] = resMetadata
		}
	}

	if p, ok := p.(provider.ProviderWithFunctions); ok {
		for _, builder := range p.Functions(ctx) {
			fun := builder()

			// Get the function name
			var metadataResp function.MetadataResponse
			fun.Metadata(ctx, function.MetadataRequest{}, &metadataResp)
			functionName := metadataResp.Name

			// Get the function definition
			var definitionResp function.DefinitionResponse
			fun.Definition(ctx, function.DefinitionRequest{}, &definitionResp)
			diags.Append(schemaResp.Diagnostics...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("function %q", functionName))
				return
			}

			sch, odiags := NewFunctionSchema(ctx, definitionResp.Definition)
			diags.Append(odiags...)
			if diags.HasError() {
				diags = WrapErrDiags(diags, fmt.Sprintf("function %q", functionName))
				return
			}

			funcMetadata := FunctionMetadata{
				Schema: sch,
			}

			if fun, ok := fun.(FunctionWithRenderOption); ok {
				funcMetadata.RenderOption = new(fun.RenderOption())
			}

			metadata.Functions[functionName] = funcMetadata
		}
	}

	return
}
