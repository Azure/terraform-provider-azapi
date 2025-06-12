package docstrings

func SensitiveBody() string {
	return `A dynamic attribute that contains the write-only properties of the request body. This will be merge-patched to the body to construct the actual request body.`
}

func SensitiveBodyVersion() string {
	return "A map where the key is the path to the property in `sensitive_body` and the value is the version of the property. " +
		"The key is a string in the format of `path.to.property[index].subproperty`, where `index` is the index of the item in an array. " +
		"When the version is changed, the property will be included in the request body, otherwise it will be omitted from the request body. "
}
