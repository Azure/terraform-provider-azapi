package metadata

import (
	"bytes"
	"fmt"
	"io"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

type Example struct {
	// The section header.
	Header string
	// The description of the example.
	Description string
	// The HCL code of the example.
	HCL string
}

func renderExamples(w io.Writer, examples []Example) error {
	if len(examples) != 0 {
		if _, err := io.WriteString(w, "## Example Usage\n"); err != nil {
			return err
		}
		for _, example := range examples {
			if example.Header != "" {
				if _, err := fmt.Fprintf(w, "\n### %s\n", example.Header); err != nil {
					return err
				}
			}
			if example.Description != "" {
				if _, err := fmt.Fprintf(w, "\n%s\n", example.Description); err != nil {
					return err
				}
			}
			if example.HCL != "" {
				if _, err := fmt.Fprintf(w, "\n```terraform\n%s\n```\n", bytes.TrimSpace(hclwrite.Format([]byte(example.HCL)))); err != nil {
					return err
				}
			}
		}
	}
	return nil
}
