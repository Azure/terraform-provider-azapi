package docstrings

const (
	resourceActionStr = `The name of the resource action. It's also possible to make HTTP requests towards the resource ID if leave this field empty.`
)

// ResourceAction returns the docstring for the resourceAction schema attribute.
func ResourceAction() string {
	return resourceActionStr
}
