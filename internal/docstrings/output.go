package docstrings

import "strings"

const (
	outputStr = `The output HCL object containing the properties specified in %sresponse_export_values%s. Here are some examples to use the values.

	%s%s%sterraform
	// it will output "registry1.azurecr.io"
	output "login_server" {
		value = RESOURCE.example.output.properties.loginServer
	}

	// it will output "disabled"
	output "quarantine_policy" {
		value = RESOURCE.example.output.properties.policies.quarantinePolicy.status
	}
	%s%s%s
`
)

func Output(res string) string {
	return addBackquotes(strings.ReplaceAll(outputStr, "RESOURCE", res))
}
