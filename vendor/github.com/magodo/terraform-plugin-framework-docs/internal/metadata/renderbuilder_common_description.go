package metadata

import (
	"fmt"
	"io"
)

func renderDescription(w io.Writer, category Category, providerName, resourceType, deprecation, description string) error {
	if category == CategoryProvider {
		if _, err := fmt.Fprintf(w, `# %s Provider
`, providerName); err != nil {
			return err
		}
	} else {
		if _, err := fmt.Fprintf(w, `# %s (%s)
`, resourceType, category); err != nil {
			return err
		}
	}

	if deprecation != "" {
		if _, err := fmt.Fprintf(w, `
!> %s
`, deprecation); err != nil {
			return nil
		}
	}

	if _, err := fmt.Fprintf(w, `
%s
`, description); err != nil {
		return err
	}

	return nil
}
