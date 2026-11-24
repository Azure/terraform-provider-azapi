package docstrings

const (
	listUniqueIdPropertyStr = `A mapping of list property paths to the field name used as a unique identifier when comparing and merging list items. When not set, list items are matched by a %sname%s property (if present) or by list ordering. To match using multiple fields, specify a comma-separated list of field names (e.g., %s"category, categoryGroup"%s).`
)

func ListUniqueIdProperty() string {
	return addBackquotes(listUniqueIdPropertyStr)
}

func ListUniqueIdPropertyDeprecationMessage() string {
	return `list_unique_id_property is deprecated. To ignore ordering differences when managing a complete collection, use Terraform's [toset()](https://developer.hashicorp.com/terraform/language/functions/toset) function. toset() does not replace matching or merging list items by custom identifiers. Configuration that manages only part of the collection should be updated to manage the complete collection. It is against AzAPI design principles to workaround ARM API bug / design deficiencies.`
}
