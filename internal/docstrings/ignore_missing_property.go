package docstrings

const (
	ignoreMissingPropertyStr = `Whether ignore not returned properties like credentials in %sbody%s to suppress plan-diff. Defaults to %strue%s. It's recommend to enable this option when some sensitive properties are not returned in response body, instead of setting them in %slifecycle.ignore_changes%s because it will make the sensitive fields unable to update.`
)

// IgnoreMissingProperty returns the docstring for ignore_missing_property schema attribute.
func IgnoreMissingProperty() string {
	return addBackquotes(ignoreMissingPropertyStr)
}
