package metadata

import (
	"io"
)

type ephemeralRenderBuilder struct {
	ProviderName string
	ResourceType string

	Metadata EphemeralMetadata

	Subcategory string
	Examples    []Example
}

func (b ephemeralRenderBuilder) Category() Category {
	return CategoryEphemeral
}

func (b ephemeralRenderBuilder) renderHeader(w io.Writer) error {
	return renderHeader(w, b.Category(), b.ProviderName, b.ResourceType, b.Subcategory, b.Metadata.Schema.Description)
}

func (b ephemeralRenderBuilder) renderDescription(w io.Writer) error {
	return renderDescription(w, b.Category(), b.ProviderName, b.ResourceType, b.Metadata.Schema.Deprecation, b.Metadata.Schema.Description)
}

func (b ephemeralRenderBuilder) renderExample(w io.Writer) error {
	return renderExamples(w, b.Examples)
}

func (b ephemeralRenderBuilder) renderSchema(w io.Writer) error {
	return renderSchema(w, b.Metadata.Schema.Fields, b.Metadata.Schema.Nested)
}
