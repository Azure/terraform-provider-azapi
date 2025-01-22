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

	sensitiveResponseExportValuesStr = `The attribute can accept either a list or a map.

- **List**: A list of paths that need to be exported from the response body. Setting it to %s["*"]%s will export the full response body. Here's an example. If it sets to %s["properties.loginServer", "properties.policies.quarantinePolicy.status"]%s, it will set the following HCL object to the computed property sensitive_output.

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

- **Map**: A map where the key is the name for the result and the value is a JMESPath query string to filter the response. Here's an example. If it sets to %s{"login_server": "properties.loginServer", "quarantine_status": "properties.policies.quarantinePolicy.status"}%s, it will set the following HCL object to the computed property sensitive_output.

	%s%s%stext
	{
		"login_server" = "registry1.azurecr.io"
		"quarantine_status" = "disabled"
	}
	%s%s%s

To learn more about JMESPath, visit [JMESPath](https://jmespath.org/).
`

	responseExportValuesForResourceListStr = `The attribute can accept either a list or a map.

- **List**: A list of paths that need to be exported from the response body. Setting it to %s["*"]%s will export the full response body. Here's an example. If it sets to %s["value"]%s, it will set the following HCL object to the computed property output.

	%s%s%stext
	{
	  "value" = [
		{
		  "id" = "/subscriptions/000000/resourceGroups/demo-rg/providers/Microsoft.Automation/automationAccounts/example"
		  "location" = "eastus2"
		  "name" = "example"
		  "properties" = {
			"creationTime" = "2024-10-11T08:18:38.737+00:00"
			"disableLocalAuth" = false
			"lastModifiedTime" = "2024-10-11T08:18:38.737+00:00"
			"publicNetworkAccess" = true
		  }
		  "tags" = {}
		  "type" = "Microsoft.Automation/AutomationAccounts"
		}
	  ]
	}
	%s%s%s

- **Map**: A map where the key is the name for the result and the value is a JMESPath query string to filter the response. Here's an example. If it sets to %s{"values": "value[].{name: name, publicNetworkAccess: properties.publicNetworkAccess}", "names": "value[].name"}%s, it will set the following HCL object to the computed property output.

	%s%s%stext
	{
		"names" = [
			"example",
			"fredaccount01",
		]
		"values" = [
			{
			  "name" = "example"
			  "publicNetworkAccess" = true
			},
			{
			  "name" = "fredaccount01"
			  "publicNetworkAccess" = null
			},
		]
	}
	%s%s%s

To learn more about JMESPath, visit [JMESPath](https://jmespath.org/).
`
)

// ResponseExportValues returns the docstring for the response_export_values schema attribute.
func ResponseExportValues() string {
	return addBackquotes(responseExportValuesStr)
}

// ResponseExportValuesForResourceList returns the docstring for the response_export_values schema attribute for the resource list data source.
func ResponseExportValuesForResourceList() string {
	return addBackquotes(responseExportValuesForResourceListStr)
}

// SensitiveResponseExportValues returns the docstring for the response_export_values schema attribute.
func SensitiveResponseExportValues() string {
	return addBackquotes(sensitiveResponseExportValuesStr)
}
