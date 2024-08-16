package docstrings

const (
	identityTenantID = `The Tenant ID for the Service Principal associated with the Managed Service Identity of this Azure resource.`
)

// IdentityTenantID returns the docstring for the identity principal ID attribute.
func IdentityTenantID() string {
	return identityTenantID
}
