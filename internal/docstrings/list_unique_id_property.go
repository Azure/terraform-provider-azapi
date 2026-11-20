package docstrings

const (
	listUniqueIdPropertyStr = `A mapping of list property paths to the string field name or dot-separated nested path used as a unique identifier within each list item when comparing and merging list items (e.g., %s"properties.virtualNetworkRules" = "subnet.id"%s). An exact field name takes precedence over nested path traversal. When not set, list items are matched by a %sname%s property (if present) or by list ordering. To match using multiple fields, specify comma-separated field names or nested paths (e.g., %s"subnet.id, category"%s). Missing, null, or non-string identifier values are treated as unresolved. If no component resolves to a non-empty string, list ordering is used. Identifier values are compared case-insensitively when %signore_casing%s is enabled, field names remain case-sensitive.`
)

func ListUniqueIdProperty() string {
	return addBackquotes(listUniqueIdPropertyStr)
}
