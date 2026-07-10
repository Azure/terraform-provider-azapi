package tffwdocs

import (
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/list"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

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
	RenderOption() EphemeralResourceRenderOption
}

type ListResourceWithRenderOption interface {
	list.ListResource
	RenderOption() ListResourceRenderOption
}

type ActionWithRenderOption interface {
	action.Action
	RenderOption() ActionRenderOption
}

type FunctionWithRenderOption interface {
	function.Function
	RenderOption() FunctionRenderOption
}
