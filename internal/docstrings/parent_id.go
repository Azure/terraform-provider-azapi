package docstrings

const (
	parentIDStr = `The ID of the azure resource in which this resource is created. It supports different kinds of deployment scope for **top level** resources:

  - resource group scope: %sparent_id%s should be the ID of a resource group, it's recommended to manage a resource group by azurerm_resource_group.
	- management group scope: %sparent_id%s should be the ID of a management group, it's recommended to manage a management group by azurerm_management_group.
	- extension scope: %sparent_id%s should be the ID of the resource you're adding the extension to.
	- subscription scope: %sparent_id%s should be like \x60/subscriptions/00000000-0000-0000-0000-000000000000\x60
	- tenant scope: %sparent_id%s should be /

  For child level resources, the %sparent_id%s should be the ID of its parent resource, for example, subnet resource's %sparent_id%s is the ID of the vnet.

  For type %sMicrosoft.Resources/resourceGroups%s, the %sparent_id%s could be omitted, it defaults to subscription ID specified in provider or the default subscription (You could check the default subscription by azure cli command: %saz account show%s).`
)

// ParentID returns the docstring for the parent_id schema attribute.
func ParentID() string {
	return addBackquotes(parentIDStr)
}
