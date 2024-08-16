package docstrings

const (
	schemaValidationEnabledStr = `Whether enabled the validation on %stype%s and %sbody%s with embedded schema. Defaults to %strue%s.`
)

// SchemaValidationEnabled returns the docstring for the schema_validation_enabled schema attribute.
func SchemaValidationEnabled() string {
	return addBackquotes(schemaValidationEnabledStr)
}
