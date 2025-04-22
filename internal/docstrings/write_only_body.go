package docstrings

func SensitiveBody() string {
	return `A dynamic attribute that contains the write-only properties of the request body. This will be merge-patched to the body to construct the actual request body.`
}
