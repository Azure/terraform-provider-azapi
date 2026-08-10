package docstrings

const (
	bodyResponseStr = `A dynamic attribute that contains the response body.`
)

// BodyResponse returns the docstring for the computed body schema attribute
// used by read-only data sources and ephemeral resources.
func BodyResponse() string {
	return bodyResponseStr
}
