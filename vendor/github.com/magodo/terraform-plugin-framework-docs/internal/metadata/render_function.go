package metadata

import (
	"bytes"
	"fmt"
	"io"
	"text/template"
)

type FunctionRenderOption struct {
	// The subcategory of the document.
	Subcategory string
	Examples    []Example

	// Description of the return value.
	ReturnDescription *string

	// A custom template that overrides the default template:
	//
	// {{ .Header }}
	// {{ .Description }}
	// {{- with .Example }}
	// {{ . }}
	// {{- end }}
	// {{ .Signature }}
	// {{- with .Arguments }}
	// {{ . }}
	// {{- end }}
	// {{- with .Return }}
	// {{ . }}
	// {{- end }}
	Template *template.Template
}

type FunctionRender struct {
	Template    *template.Template
	Header      string
	Description string
	Example     string
	Signature   string
	Arguments   string
	Return      string
}

func (metadata Metadata) NewFunctionRender(functionName string, opt *FunctionRenderOption) (*FunctionRender, error) {
	fmetadata, ok := metadata.Functions[functionName]
	if !ok {
		return nil, fmt.Errorf("unknown function name %q", functionName)
	}

	src := functionRenderBuilder{
		ProviderName: metadata.ProviderName,
		FunctionName: functionName,
		Metadata:     fmetadata,
	}

	if opt == nil {
		opt = fmetadata.RenderOption
	}

	var tpl *template.Template
	if opt != nil {
		tpl = opt.Template
		src.Subcategory = opt.Subcategory
		src.Examples = opt.Examples
		src.ReturnDescription = opt.ReturnDescription
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
	signatureBuf := bytes.NewBuffer(nil)
	if err := src.renderSignature(signatureBuf); err != nil {
		return nil, err
	}
	argumentsBuf := bytes.NewBuffer(nil)
	if err := src.renderArguments(argumentsBuf); err != nil {
		return nil, err
	}
	returnBuf := bytes.NewBuffer(nil)
	if err := src.renderReturn(returnBuf); err != nil {
		return nil, err
	}

	return &FunctionRender{
		Template:    tpl,
		Header:      headerBuf.String(),
		Description: descriptionBuf.String(),
		Example:     exampleBuf.String(),
		Signature:   signatureBuf.String(),
		Arguments:   argumentsBuf.String(),
		Return:      returnBuf.String(),
	}, nil
}

const functionTpl = `{{ .Header }}
{{ .Description }}
{{- with .Example }}
{{ . }}
{{- end }}
{{ .Signature }}
{{- with .Arguments }}
{{ . }}
{{- end }}
{{- with .Return }}
{{ . }}
{{- end }}`

func (render FunctionRender) Execute(w io.Writer) error {
	tpl := render.Template
	if tpl == nil {
		var err error
		tpl, err = template.New("function").Parse(functionTpl)
		if err != nil {
			return err
		}
	}

	return tpl.Execute(w, render)
}
