package docstrings

const (
	listUniqueIdPropertyStr = `A mapping of list property paths to the field name used as a unique identifier when comparing and merging list items. When not set, list items are matched by a %sname%s property (if present) or by list ordering. To match using multiple fields, specify a comma-separated list of field names (e.g., %s"category, categoryGroup"%s).`
)

func ListUniqueIdProperty() string {
	return addBackquotes(listUniqueIdPropertyStr)
}
