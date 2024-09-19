package docstrings

const (
	responseExportValuesStr = `The attribute can accept either a list or a map.

- **List**: A list of paths that need to be exported from the response body. Setting it to %s["*"]%s will export the full response body. Here's an example. If it sets to %s["properties.loginServer", "properties.policies.quarantinePolicy.status"]%s, it will set the following HCL object to the computed property output.

	%s%s%stext
	{
		properties = {
			loginServer = "registry1.azurecr.io"
			policies = {
				quarantinePolicy = {
					status = "disabled"
				}
			}
		}
	}
	%s%s%s

- **Map**: A map where the key is the name for the result and the value is a JMESPath query string to filter the response. Here's an example. If it sets to %s{"login_server": "properties.loginServer", "quarantine_status": "properties.policies.quarantinePolicy.status"}%s, it will set the following HCL object to the computed property output.

	%s%s%stext
	{
		"login_server" = "registry1.azurecr.io"
		"quarantine_status" = "disabled"
	}
	%s%s%s

To learn more about JMESPath, visit [JMESPath](https://jmespath.org/).
`
)

// ResponseExportValues returns the docstring for the response_export_values schema attribute.
func ResponseExportValues() string {
	return addBackquotes(responseExportValuesStr)
}
