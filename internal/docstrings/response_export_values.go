package docstrings

const (
	responseExportValuesStr = `A list of path that needs to be exported from response body. Setting it to %s["*"]%s will export the full response body. Here's an example. If it sets to %s["properties.loginServer", "properties.policies.quarantinePolicy.status"]%s, it will set the following HCL object to computed property output.

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
`
)

// ResponseExportValues returns the docstring for the response_export_values schema attribute.
func ResponseExportValues() string {
	return addBackquotes(responseExportValuesStr)
}
