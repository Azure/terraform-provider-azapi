package docstrings

const (
	typeStr = `In a format like %s<resource-type>@<api-version>%s. %s<resource-type>%s is the Azure resource type, for example, %sMicrosoft.Storage/storageAccounts%s. %s<api-version>%s is version of the API used to manage this azure resource.`
)

// Type returns the docstring for the type schema attribute.
func Type() string {
	return addBackquotes(typeStr)
}
