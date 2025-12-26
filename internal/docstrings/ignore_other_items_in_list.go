package docstrings

const (
	ignoreOtherItemsInListStr = `A list of list property paths where items not specified in configuration should be ignored.

This is intended for partial list management when combined with %slist_unique_id_property%s (for example, to avoid perpetual drift from server-side ordering).`
)

func IgnoreOtherItemsInList() string {
	return addBackquotes(ignoreOtherItemsInListStr)
}
