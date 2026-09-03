package docstrings

const (
	dataPlaneTypeRefStr = ` For a list of supported data plane resource types, see the [Available Resources](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/data_plane_resource#available-resources) documentation.`
)

// DataPlaneType returns the docstring for the type schema attribute of data plane resources,
// including a reference to the list of supported data plane resource types.
func DataPlaneType() string {
	return addBackquotes(typeStr) + dataPlaneTypeRefStr
}
