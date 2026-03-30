package metadata

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strings"
)

type functionRenderBuilder struct {
	ProviderName string
	FunctionName string

	Subcategory string
	Examples    []Example

	Metadata FunctionMetadata

	ReturnDescription *string
}

func (b functionRenderBuilder) Category() Category {
	return CategoryFunction
}

func (b functionRenderBuilder) renderHeader(w io.Writer) error {
	return renderHeader(w, b.Category(), b.ProviderName, b.FunctionName, b.Subcategory, b.Metadata.Schema.Summary)
}

func (b functionRenderBuilder) renderDescription(w io.Writer) error {
	return renderDescription(w, b.Category(), b.ProviderName, b.FunctionName, b.Metadata.Schema.Deprecation, b.Metadata.Schema.Description)
}

func (b functionRenderBuilder) renderExample(w io.Writer) error {
	return renderExamples(w, b.Examples)
}

func (b functionRenderBuilder) renderSignature(w io.Writer) error {
	if _, err := fmt.Fprintf(w, `## Signature

%[1]s%[1]s%[1]stext
%[2]s(`, "`", b.FunctionName); err != nil {
		return err
	}
	schema := b.Metadata.Schema

	var params []string
	for _, param := range schema.Parameters {
		if param.isVariadic {
			params = append(params, fmt.Sprintf("%s %s...", param.name, param.dataType))
		} else {
			params = append(params, fmt.Sprintf("%s %s", param.name, param.dataType))
		}
	}
	if _, err := io.WriteString(w, strings.Join(params, ",")); err != nil {
		return err
	}
	if _, err := fmt.Fprintf(w, `) %[1]s
%[2]s%[2]s%[2]s
`, schema.Return.dataType, "`"); err != nil {
		return err
	}

	return nil
}

func (b functionRenderBuilder) renderArguments(w io.Writer) error {
	if len(b.Metadata.Schema.Parameters) == 0 {
		return nil
	}
	if _, err := fmt.Fprintf(w, `## Arguments

`); err != nil {
		return err
	}

	for _, param := range b.Metadata.Schema.Parameters {
		if err := b.renderArgument(w, param); err != nil {
			return err
		}
	}

	return b.renderObjects(w, b.Metadata.Schema.Objects)
}

func (b functionRenderBuilder) renderReturn(w io.Writer) error {
	if !(b.Metadata.Schema.ReturnObjects != nil || b.ReturnDescription != nil) {
		return nil
	}

	if _, err := fmt.Fprintf(w, `## Return
`); err != nil {
		return err
	}

	if v := b.Metadata.Schema.Return.CustomTypeDescription(); v != "" {
		if _, err := fmt.Fprintf(w, "\n-> %s\n", v); err != nil {
			return err
		}
	}

	if desc := b.ReturnDescription; desc != nil {
		if _, err := fmt.Fprintf(w, "\n%s\n", *desc); err != nil {
			return err
		}
	}

	if objs := b.Metadata.Schema.ReturnObjects; objs != nil {
		if _, err := fmt.Fprintf(w, "\nThe `object` returned from `%s` has the following attributes:\n\n", b.FunctionName); err != nil {
			return err
		}

		objs := maps.Clone(objs)

		if err := b.renderObject(w, objs[""]); err != nil {
			return err
		}

		delete(objs, "")
		if err := b.renderObjects(w, objs); err != nil {
			return err
		}
	}

	return nil
}

func (b functionRenderBuilder) renderArgument(w io.Writer, field FunctionField) error {
	if _, err := fmt.Fprintf(w, "1. `%s` (%s) %s", field.Name(), field.Traits(), field.Description()); err != nil {
		return err
	}
	if v := field.NestedLink(); v != "" {
		fmt.Fprintf(w, " %s", v)
	}
	io.WriteString(w, "\n")

	if l := field.Validators(); len(l) != 0 {
		for _, e := range l {
			if e == "" {
				continue
			}
			if _, err := fmt.Fprintf(w, "\n\t-> %s\n", e); err != nil {
				return err
			}
		}
	}

	if v := field.CustomTypeDescription(); v != "" {
		if _, err := fmt.Fprintf(w, "\n\t-> %s\n", v); err != nil {
			return err
		}
	}

	return nil
}

func (b functionRenderBuilder) renderObjects(w io.Writer, objs FunctionObjects) error {
	keys := slices.Collect(maps.Keys(objs))
	slices.Sort(keys)
	for _, key := range keys {
		obj := objs[key]

		if _, err := fmt.Fprintf(w, `
<a id="nested--%[1]s"></a>
### Fields of %[2]s%[1]s%[2]s

`, key, "`"); err != nil {
			return err
		}

		if err := b.renderObject(w, obj); err != nil {
			return err
		}
	}
	return nil
}

func (b functionRenderBuilder) renderObject(w io.Writer, object FunctionObject) error {
	keys := slices.Collect(maps.Keys(object.fields))
	slices.Sort(keys)
	for _, key := range keys {
		field := object.fields[key]
		if _, err := fmt.Fprintf(w, "- `%s` (%s)", key, field.dataType); err != nil {
			return err
		}
		if v := field.Description(); v != "" {
			if _, err := fmt.Fprintf(w, " %s", v); err != nil {
				return err
			}
		}
		if v := field.NestedLink(); v != "" {
			fmt.Fprintf(w, " %s", v)
		}
		io.WriteString(w, "\n")
	}

	if v := object.CustomTypeDescription(); v != "" {
		if _, err := fmt.Fprintf(w, "\n-> %s", v); err != nil {
			return err
		}
	}

	return nil
}
