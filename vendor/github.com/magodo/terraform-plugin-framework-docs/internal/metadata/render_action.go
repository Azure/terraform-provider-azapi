package metadata

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

type ActionRenderOption struct {
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

type ActionRender struct {
	Template    *template.Template
	Header      string
	Description string
	Example     string
	Schema      string
	Import      string
}

func (metadata Metadata) NewActionRender(actionType string, opt *ActionRenderOption) (*ActionRender, error) {
	resmetadata, ok := metadata.Actions[actionType]
	if !ok {
		return nil, fmt.Errorf("unknown action type %q", actionType)
	}

	src := actionRenderBuilder{
		ProviderName: metadata.ProviderName,
		ActionType:   actionType,
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

	return &ActionRender{
		Template:    tpl,
		Header:      headerBuf.String(),
		Description: descriptionBuf.String(),
		Example:     exampleBuf.String(),
		Schema:      schemaBuf.String(),
	}, nil
}

const actionTpl = `{{ .Header }}
{{ .Description }}
{{- with .Example }}
{{ . }}
{{- end }}
{{ .Schema }}`

func (render ActionRender) Execute(w io.Writer) error {
	tpl := render.Template
	if tpl == nil {
		var err error
		tpl, err = template.New("action").Parse(actionTpl)
		if err != nil {
			return err
		}
	}

	return tpl.Execute(w, render)
}
