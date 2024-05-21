package docstrings

const (
	identityType = `The Type of Identity which should be used for this azure resource. Possible values are %sSystemAssigned%s, %sUserAssigned%s and %sSystemAssigned,UserAssigned%s`
)

// IdentityType returns the docstring for the identity type attribute.
func IdentityType() string {
	return addBackquotes(identityType)
}
