package metadata

import (
	"io"
)

type dataSourceRenderBuilder struct {
	ProviderName string
	ResourceType string

	Metadata DataSourceMetadata

	Subcategory string
	Examples    []Example
}

func (b dataSourceRenderBuilder) Category() Category {
	return CategoryDataSource
}

func (b dataSourceRenderBuilder) renderHeader(w io.Writer) error {
	return renderHeader(w, b.Category(), b.ProviderName, b.ResourceType, b.Subcategory, b.Metadata.Schema.Description)
}

func (b dataSourceRenderBuilder) renderDescription(w io.Writer) error {
	return renderDescription(w, b.Category(), b.ProviderName, b.ResourceType, b.Metadata.Schema.Deprecation, b.Metadata.Schema.Description)
}

func (b dataSourceRenderBuilder) renderExample(w io.Writer) error {
	return renderExamples(w, b.Examples)
}

func (b dataSourceRenderBuilder) renderSchema(w io.Writer) error {
	return renderSchema(w, b.Metadata.Schema.Fields, b.Metadata.Schema.Nested)
}
