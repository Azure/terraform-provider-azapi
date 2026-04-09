package metadata

import (
	"fmt"
	"io"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

type ImportId struct {
	// The Id format.
	Format string

	// The example id that will be displayed in an example `terraform import` command.
	// Note that the id is double quoted, ensure to escape any double quote included in the id.
	ExampleId string

	// The complete import by id block. If not specified, it will fill in the block with the `ExampleId`.
	ExampleBlk string
}

type resourceRenderBuilder struct {
	ProviderName string
	ResourceType string

	Metadata ResourceMetadata

	Subcategory string
	Examples    []Example

	// Import by id information.
	ImportId *ImportId

	// The identity import block examples.
	// The HCL shall only contain the config inside the `identity` object.
	IdentityExamples []Example
}

func (b resourceRenderBuilder) Category() Category {
	return CategoryResource
}

func (b resourceRenderBuilder) renderHeader(w io.Writer) error {
	return renderHeader(w, b.Category(), b.ProviderName, b.ResourceType, b.Subcategory, b.Metadata.Schema.Description)
}

func (b resourceRenderBuilder) renderDescription(w io.Writer) error {
	return renderDescription(w, b.Category(), b.ProviderName, b.ResourceType, b.Metadata.Schema.Deprecation, b.Metadata.Schema.Description)
}

func (b resourceRenderBuilder) renderExample(w io.Writer) error {
	return renderExamples(w, b.Examples)
}

func (b resourceRenderBuilder) renderSchema(w io.Writer) error {
	return renderSchema(w, b.Metadata.Schema.Fields, b.Metadata.Schema.Nested)
}

func (b resourceRenderBuilder) renderImport(w io.Writer) error {
	if identity := b.Metadata.Identity; b.ImportId != nil || identity != nil {
		io.WriteString(w, "## Import\n")

		if b.ImportId != nil {
			io.WriteString(w, "\n")
			if err := b.renderImportId(w, *b.ImportId); err != nil {
				return err
			}
		}

		if identity != nil {
			io.WriteString(w, "\n")
			if err := b.renderImportIdentity(w, *identity); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b resourceRenderBuilder) renderImportId(w io.Writer, importId ImportId) error {
	if importId.Format == "" {
		return fmt.Errorf("the `.Format` of the ImportId is not specified")
	}
	if importId.ExampleId == "" {
		return fmt.Errorf("the `.ExampleId` of the ImportId is not specified")
	}

	importBlk := importId.ExampleBlk
	if importBlk == "" {
		importBlk = fmt.Sprintf(`
import {
  to = %s.example
  id = "%s"
}
`, b.ResourceType, importId.ExampleId)
	}

	importBlk = string(hclwrite.Format([]byte(strings.TrimSpace(importBlk))))

	if _, err := fmt.Fprintf(w, `### Import ID

The [%[1]sterraform import%[1]s command](https://developer.hashicorp.com/terraform/cli/commands/import) can be used with the id format:

%[1]s%[1]s%[1]sshell
%[2]s
%[1]s%[1]s%[1]s

For example:

%[1]s%[1]s%[1]sshell
$ terraform import %[3]s.example "%[4]s"
%[1]s%[1]s%[1]s

In Terraform v1.5.0 and later, the [%[1]simport%[1]s block](https://developer.hashicorp.com/terraform/language/block/import) can be used with the %[1]sid%[1]s attribute, for example:

%[1]s%[1]s%[1]sterraform
%[5]s
%[1]s%[1]s%[1]s
`, "`", importId.Format, b.ResourceType, importId.ExampleId, importBlk); err != nil {
		return err
	}

	return nil
}

func (b resourceRenderBuilder) renderImportIdentity(w io.Writer, schema ResourceIdentitySchema) error {
	formatExample := func(example string) []byte {
		return hclwrite.Format([]byte(strings.TrimSpace(fmt.Sprintf(`
import {
	to = %s.example
	identity = {
		%s
	}
}
`, b.ResourceType, strings.TrimSpace(example)))))
	}

	if _, err := fmt.Fprintf(w, `### Import Identity

In Terraform v1.12.0 and later, the [%[1]simport%[1]s block](https://developer.hashicorp.com/terraform/language/block/import) can be used with the %[1]sidentity%[1]s attribute.
`, "`"); err != nil {
		return err
	}

	for _, example := range b.IdentityExamples {
		if example.Header != "" {
			if _, err := fmt.Fprintf(w, "\n#### Example: %s\n", example.Header); err != nil {
				return err
			}
		}
		if example.Description != "" {
			if _, err := fmt.Fprintf(w, "\n%s\n", example.Description); err != nil {
				return err
			}
		}
		if example.HCL != "" {
			if _, err := fmt.Fprintf(w, "\n```terraform\n%s\n```\n", formatExample(example.HCL)); err != nil {
				return err
			}
		}
	}

	sections := []struct {
		name   string
		fields []ResourceIdentityField
	}{
		{
			name:   "Required",
			fields: schema.Fields.RequiredFields(),
		},
		{
			name:   "Optional",
			fields: schema.Fields.OptionalFields(),
		},
	}

	for _, section := range sections {
		if len(section.fields) == 0 {
			continue
		}
		if _, err := fmt.Fprintf(w, `
%s:

`, section.name); err != nil {
			return err
		}

		for _, field := range section.fields {
			if err := b.renderIdentityField(w, field); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b resourceRenderBuilder) renderIdentityField(w io.Writer, field ResourceIdentityField) error {
	if _, err := fmt.Fprintf(w, "- `%s` (%s) %s\n", field.Name, field.Traits(), field.Description); err != nil {
		return err
	}
	if v := field.CustomTypeDescription(); v != "" {
		if _, err := fmt.Fprintf(w, "\n\t-> %s\n", v); err != nil {
			return err
		}
	}
	return nil
}
