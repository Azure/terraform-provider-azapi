package metadata

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

type ListRenderOption struct {
	// The subcategory of the document.
	Subcategory string
	Examples    []Example

	// A custom template that overrides the default template:
	//
	// {{ .Header }}
	// {{ .Description }}
	// {{- with .Example }}
	// {{ . }}
	// {{- end }}
	// {{ .Schema }}
	Template *template.Template
}

type ListRender struct {
	Template    *template.Template
	Header      string
	Description string
	Example     string
	Schema      string
}

func (metadata Metadata) NewListRender(listType string, opt *ListRenderOption) (*ListRender, error) {
	resmetadata, ok := metadata.Lists[listType]
	if !ok {
		return nil, fmt.Errorf("unknown list resource type %q", listType)
	}

	src := listRenderBuilder{
		ProviderName: metadata.ProviderName,
		ListType:     listType,
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

	return &ListRender{
		Template:    tpl,
		Header:      headerBuf.String(),
		Description: descriptionBuf.String(),
		Example:     exampleBuf.String(),
		Schema:      schemaBuf.String(),
	}, nil
}

const listTpl = `{{ .Header }}
{{ .Description }}
{{- with .Example }}
{{ . }}
{{- end }}
{{ .Schema }}`

func (render ListRender) Execute(w io.Writer) error {
	tpl := render.Template
	if tpl == nil {
		var err error
		tpl, err = template.New("list").Parse(listTpl)
		if err != nil {
			return err
		}
	}

	return tpl.Execute(w, render)
}
