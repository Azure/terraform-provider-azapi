package docstrings

const (
	locksStr = `A list of ARM resource IDs which are used to avoid create/modify/delete azapi resources at the same time.`
)

// Locks returns the docstring for the locks schema attribute.
func Locks() string {
	return locksStr
}
