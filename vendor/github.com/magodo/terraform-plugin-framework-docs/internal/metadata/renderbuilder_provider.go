package metadata

import (
	"io"
)

type providerRenderBuilder struct {
	ProviderName string
	Schema       ProviderSchema
	Examples     []Example
}

func (b providerRenderBuilder) Category() Category {
	return CategoryProvider
}

func (b providerRenderBuilder) renderHeader(w io.Writer) error {
	return renderHeader(w, b.Category(), b.ProviderName, "", "", b.Schema.Description)
}

func (b providerRenderBuilder) renderDescription(w io.Writer) error {
	return renderDescription(w, b.Category(), b.ProviderName, "", b.Schema.Deprecation, b.Schema.Description)
}

func (b providerRenderBuilder) renderExample(w io.Writer) error {
	return renderExamples(w, b.Examples)
}

func (b providerRenderBuilder) renderSchema(w io.Writer) error {
	return renderSchema(w, b.Schema.Fields, b.Schema.Nested)
}
