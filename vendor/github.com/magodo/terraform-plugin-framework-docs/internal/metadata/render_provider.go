package metadata

import (
	"bytes"
	"io"
	"text/template"
)

type ProviderRenderOption struct {
	Examples []Example

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

type ProviderRender struct {
	Template    *template.Template
	Header      string
	Description string
	Example     string
	Schema      string
}

func (metadata Metadata) NewProviderRender(opt *ProviderRenderOption) (*ProviderRender, error) {
	src := providerRenderBuilder{
		ProviderName: metadata.ProviderName,
		Schema:       metadata.Provider.Schema,
	}

	if opt == nil {
		opt = metadata.Provider.RenderOption
	}

	var tpl *template.Template
	if opt != nil {
		tpl = opt.Template
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

	return &ProviderRender{
		Template:    tpl,
		Header:      headerBuf.String(),
		Description: descriptionBuf.String(),
		Example:     exampleBuf.String(),
		Schema:      schemaBuf.String(),
	}, nil
}

const providerTpl = `{{ .Header }}
{{ .Description }}
{{- with .Example }}
{{ . }}
{{- end }}
{{ .Schema }}`

func (render ProviderRender) Execute(w io.Writer) error {
	tpl := render.Template
	if tpl == nil {
		var err error
		tpl, err = template.New("provider").Parse(providerTpl)
		if err != nil {
			return err
		}
	}

	return tpl.Execute(w, render)
}
