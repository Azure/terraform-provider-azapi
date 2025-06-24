package docstrings

const (
	ignoreNullPropertyStr = `When set to %strue%s, the provider will ignore properties whose values are %snull%s in the %sbody%s.
These properties will not be included in the request body sent to the API, and the difference will not be shown in the plan output. Defaults to %sfalse%s.`
)

// IgnoreNullProperty returns the docstring for ignore_missing_property schema attribute.
func IgnoreNullProperty() string {
	return addBackquotes(ignoreNullPropertyStr)
}
