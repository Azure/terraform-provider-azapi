package docstrings

const (
	bodyStr = `A dynamic attribute that contains the request body.`
)

// Body returns the docstring for the body schema attribute.
func Body() string {
	return bodyStr
}
