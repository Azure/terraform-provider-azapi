package docstrings

const (
	ignoreCasingStr = `Whether ignore the casing of the property names in the response body. Defaults to %sfalse%s.`
)

// IgnoreCasing returns the docstring for the ignore_casing schema attribute.
func IgnoreCasing() string {
	return addBackquotes(ignoreCasingStr)
}
