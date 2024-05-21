package docstrings

const (
	identityPrincipalID = `The Principal ID for the Service Principal associated with the Managed Service Identity of this Azure resource.`
)

// IdentityPrincipalID returns the docstring for the identity principal ID attribute.
func IdentityPrincipalID() string {
	return identityPrincipalID
}
