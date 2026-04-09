package metadata

import (
	"io"
)

type listRenderBuilder struct {
	ProviderName string
	ListType     string

	Metadata ListMetadata

	Subcategory string
	Examples    []Example
}

func (b listRenderBuilder) Category() Category {
	return CategoryList
}

func (b listRenderBuilder) renderHeader(w io.Writer) error {
	return renderHeader(w, b.Category(), b.ProviderName, b.ListType, b.Subcategory, b.Metadata.Schema.Description)
}

func (b listRenderBuilder) renderDescription(w io.Writer) error {
	return renderDescription(w, b.Category(), b.ProviderName, b.ListType, b.Metadata.Schema.Deprecation, b.Metadata.Schema.Description)
}

func (b listRenderBuilder) renderExample(w io.Writer) error {
	return renderExamples(w, b.Examples)
}

func (b listRenderBuilder) renderSchema(w io.Writer) error {
	return renderSchema(w, b.Metadata.Schema.Fields, b.Metadata.Schema.Nested)
}
