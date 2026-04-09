package metadata

import (
	"fmt"
	"io"
	"maps"
	"slices"
)

func renderSchema(w io.Writer, fields Fields, nested NestedFields) error {
	if _, err := fmt.Fprintf(w, `## Schema
`); err != nil {
		return err
	}

	sections := []struct {
		name   string
		fields []Field
	}{
		{
			name:   "Required",
			fields: fields.RequiredFields(),
		},
		{
			name:   "Optional",
			fields: fields.OptionalFields(),
		},
		{
			name:   "Read-Only",
			fields: fields.ComputedFields(),
		},
	}

	for _, section := range sections {
		if len(section.fields) == 0 {
			continue
		}
		if _, err := fmt.Fprintf(w, `
### %s

`, section.name); err != nil {
			return err
		}

		for _, field := range section.fields {
			if err := renderField(w, field); err != nil {
				return err
			}
		}
	}

	if nested := nested; len(nested) != 0 {
		io.WriteString(w, "\n")
		if err := renderNestedFields(w, nested); err != nil {
			return err
		}
	}

	return nil
}

func renderNestedFields(w io.Writer, fields NestedFields) error {
	keys := slices.Collect(maps.Keys(fields))
	slices.Sort(keys)
	for _, key := range keys {
		field := fields[key]

		if _, err := fmt.Fprintf(w, `<a id="nested--%[1]s"></a>
### Nested Schema for %[2]s%[1]s%[2]s
`, key, "`"); err != nil {
			return err
		}

		if l := field.Validators(); len(l) != 0 {
			for _, e := range l {
				if e == "" {
					continue
				}
				if _, err := fmt.Fprintf(w, "\n-> %s\n", e); err != nil {
					return err
				}
			}
		}

		sections := []struct {
			name   string
			fields []Field
		}{
			{
				name:   "Required",
				fields: field.Fields().RequiredFields(),
			},
			{
				name:   "Optional",
				fields: field.Fields().OptionalFields(),
			},
			{
				name:   "Read-Only",
				fields: field.Fields().ComputedFields(),
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
				if err := renderField(w, field); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func renderField(w io.Writer, field Field) error {
	if _, err := fmt.Fprintf(w, "- `%s` (%s)", field.Name(), field.Traits()); err != nil {
		return err
	}

	if v := field.Description(); v != "" {
		fmt.Fprintf(w, " %s", v)
	}

	if v := field.Default(); v != "" {
		fmt.Fprintf(w, " %s", v)
	}
	if v := field.NestedLink(); v != "" {
		fmt.Fprintf(w, " %s", v)
	}
	io.WriteString(w, "\n")

	if l := field.PlanModifiers(); len(l) != 0 {
		for _, e := range l {
			if e == "" {
				continue
			}
			if _, err := fmt.Fprintf(w, "\n\t~> %s\n", e); err != nil {
				return err
			}
		}
	}

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

	if v := field.Deprecation(); v != "" {
		if _, err := fmt.Fprintf(w, "\n\t!> %s\n", v); err != nil {
			return err
		}
	}

	return nil
}
