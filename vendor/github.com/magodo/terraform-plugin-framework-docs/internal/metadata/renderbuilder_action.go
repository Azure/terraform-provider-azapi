package metadata

import (
	"io"
)

type actionRenderBuilder struct {
	ProviderName string
	ActionType   string

	Metadata ActionMetadata

	Subcategory string
	Examples    []Example
}

func (b actionRenderBuilder) Category() Category {
	return CategoryAction
}

func (b actionRenderBuilder) renderHeader(w io.Writer) error {
	return renderHeader(w, b.Category(), b.ProviderName, b.ActionType, b.Subcategory, b.Metadata.Schema.Description)
}

func (b actionRenderBuilder) renderDescription(w io.Writer) error {
	return renderDescription(w, b.Category(), b.ProviderName, b.ActionType, b.Metadata.Schema.Deprecation, b.Metadata.Schema.Description)
}

func (b actionRenderBuilder) renderExample(w io.Writer) error {
	return renderExamples(w, b.Examples)
}

func (b actionRenderBuilder) renderSchema(w io.Writer) error {
	return renderSchema(w, b.Metadata.Schema.Fields, b.Metadata.Schema.Nested)
}
