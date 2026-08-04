package docstrings

const (
	ignoreMissingPropertyStr = `Whether ignore not returned properties like credentials in %sbody%s to suppress plan-diff. It's recommend to enable this option when some sensitive properties are not returned in response body, instead of setting them in %slifecycle.ignore_changes%s because it will make the sensitive fields unable to update.`

	ignoreMissingPropertyDeprecatedStr = ignoreMissingPropertyStr + ` **Deprecated**: to detect properties that exist in the remote resource but are missing from %sbody%s (for example, out-of-band changes made directly in Azure), use %scompute_complete_diff%s instead, which takes precedence over this property.`
)

// IgnoreMissingProperty returns the docstring for ignore_missing_property schema attribute.
func IgnoreMissingProperty() string {
	return addBackquotes(ignoreMissingPropertyStr)
}

// IgnoreMissingPropertyDeprecated returns the docstring for the ignore_missing_property schema
// attribute on resources that also expose compute_complete_diff, appending a deprecation notice
// that points users to compute_complete_diff.
func IgnoreMissingPropertyDeprecated() string {
	return addBackquotes(ignoreMissingPropertyDeprecatedStr)
}
