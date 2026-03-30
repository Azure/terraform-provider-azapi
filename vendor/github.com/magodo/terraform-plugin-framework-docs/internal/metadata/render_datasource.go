package metadata

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

type DataSourceRenderOption struct {
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

type DataSourceRender struct {
	Template    *template.Template
	Header      string
	Description string
	Example     string
	Schema      string
}

func (metadata Metadata) NewDataSourceRender(dataSourceType string, opt *DataSourceRenderOption) (*DataSourceRender, error) {
	resmetadata, ok := metadata.DataSources[dataSourceType]
	if !ok {
		return nil, fmt.Errorf("unknown data source type %q", dataSourceType)
	}

	src := dataSourceRenderBuilder{
		ProviderName: metadata.ProviderName,
		ResourceType: dataSourceType,
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

	return &DataSourceRender{
		Template:    tpl,
		Header:      headerBuf.String(),
		Description: descriptionBuf.String(),
		Example:     exampleBuf.String(),
		Schema:      schemaBuf.String(),
	}, nil
}

const dataSourceTpl = `{{ .Header }}
{{ .Description }}
{{- with .Example }}
{{ . }}
{{- end }}
{{ .Schema }}`

func (render DataSourceRender) Execute(w io.Writer) error {
	tpl := render.Template
	if tpl == nil {
		var err error
		tpl, err = template.New("dataSource").Parse(dataSourceTpl)
		if err != nil {
			return err
		}
	}

	return tpl.Execute(w, render)
}
