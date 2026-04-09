package metadata

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

type ResourceRenderOption struct {
	// The subcategory of the document.
	Subcategory string
	Examples    []Example

	// The information about import by id (including via command and via import block).
	ImportId *ImportId
	// The examples for importing by identity via import block.
	IdentityExamples []Example

	// A custom template that overrides the default template:
	//
	// {{ .Header }}
	// {{ .Description }}
	// {{- with .Example }}
	// {{ . }}
	// {{- end }}
	// {{ .Schema }}
	// {{- with .Import }}
	// {{ . }}
	// {{- end }}
	Template *template.Template
}

type ResourceRender struct {
	Template    *template.Template
	Header      string
	Description string
	Example     string
	Schema      string
	Import      string
}

func (metadata Metadata) NewResourceRender(resourceType string, opt *ResourceRenderOption) (*ResourceRender, error) {
	resmetadata, ok := metadata.Resources[resourceType]
	if !ok {
		return nil, fmt.Errorf("unknown resource type %q", resourceType)
	}

	src := resourceRenderBuilder{
		ProviderName: metadata.ProviderName,
		ResourceType: resourceType,
		Metadata:     resmetadata,
	}

	if opt == nil {
		opt = resmetadata.RenderOption
	}

	var tpl *template.Template
	if opt != nil {
		tpl = opt.Template
		src.Subcategory = opt.Subcategory
		src.Examples = opt.Examples
		src.ImportId = opt.ImportId
		src.IdentityExamples = opt.IdentityExamples
	}

	headerBuf := bytes.NewBuffer(nil)
	if err := src.renderHeader(headerBuf); err != nil {
		return nil, err
	}
	descriptionBuf := bytes.NewBuffer(nil)
	if err := src.renderDescription(descriptionBuf); err != nil {
		return nil, err
	}
	exampleBuf := bytes.NewBuffer(nil)
	if err := src.renderExample(exampleBuf); err != nil {
		return nil, err
	}
	schemaBuf := bytes.NewBuffer(nil)
	if err := src.renderSchema(schemaBuf); err != nil {
		return nil, err
	}
	importBuf := bytes.NewBuffer(nil)
	if err := src.renderImport(importBuf); err != nil {
		return nil, err
	}

	return &ResourceRender{
		Template:    tpl,
		Header:      headerBuf.String(),
		Description: descriptionBuf.String(),
		Example:     exampleBuf.String(),
		Schema:      schemaBuf.String(),
		Import:      importBuf.String(),
	}, nil
}

const resourceTpl = `{{ .Header }}
{{ .Description }}
{{- with .Example }}
{{ . }}
{{- end }}
{{ .Schema }}
{{- with .Import }}
{{ . }}
{{- end }}`

func (render ResourceRender) Execute(w io.Writer) error {
	tpl := render.Template
	if tpl == nil {
		var err error
		tpl, err = template.New("resource").Parse(resourceTpl)
		if err != nil {
			return err
		}
	}

	return tpl.Execute(w, render)
}
