package services_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type IdDataSource struct{}

func TestAccIdDataSource_resourceGroup(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_id", "test")
	r := IdDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.resourceGroup(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_id").HasValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName"),
				check.That(data.ResourceName).Key("resource_group_name").HasValue("resourceGroupName"),
				check.That(data.ResourceName).Key("subscription_id").HasValue("00000000-0000-0000-0000-000000000000"),
				check.That(data.ResourceName).Key("provider_namespace").HasValue("Microsoft.Resources"),
			),
		},
	})
}

func TestAccIdDataSource_virtualNetworks(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_id", "test")
	r := IdDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.virtualNetworks(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_id").HasValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.Network/virtualNetworks/vnetName"),
				check.That(data.ResourceName).Key("resource_group_name").HasValue("resourceGroupName"),
				check.That(data.ResourceName).Key("subscription_id").HasValue("00000000-0000-0000-0000-000000000000"),
				check.That(data.ResourceName).Key("provider_namespace").HasValue("Microsoft.Network"),
				check.That(data.ResourceName).Key("parts.virtualNetworks").HasValue("vnetName"),
			),
		},
	})
}

func TestAccIdDataSource_subnets(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_id", "test")
	r := IdDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.subnets(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_id").HasValue("/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.Network/virtualNetworks/vnetName/subnets/subnetName"),
				check.That(data.ResourceName).Key("resource_group_name").HasValue("resourceGroupName"),
				check.That(data.ResourceName).Key("subscription_id").HasValue("00000000-0000-0000-0000-000000000000"),
				check.That(data.ResourceName).Key("provider_namespace").HasValue("Microsoft.Network"),
				check.That(data.ResourceName).Key("parts.virtualNetworks").HasValue("vnetName"),
				check.That(data.ResourceName).Key("parts.subnets").HasValue("subnetName"),
			),
		},
	})
}

func TestAccIdDataSource_tenants(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_id", "test")
	r := IdDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.tenants(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("resource_id").HasValue("/"),
				check.That(data.ResourceName).Key("resource_group_name").HasValue(""),
				check.That(data.ResourceName).Key("subscription_id").HasValue(""),
				check.That(data.ResourceName).Key("provider_namespace").HasValue("Microsoft.Resources"),
			),
		},
	})
}

func (r IdDataSource) resourceGroup() string {
	return `
data "azapi_resource_id" "test" {
  type      = "Microsoft.Resources/resourceGroups@2022-09-01"
  name      = "resourceGroupName"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000"
}
`
}

func (r IdDataSource) virtualNetworks() string {
	return `
data "azapi_resource_id" "test" {
  type      = "Microsoft.Network/virtualNetworks@2023-04-01"
  parent_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName"
  name      = "vnetName"
}
`
}

func (r IdDataSource) subnets() string {
	return `
data "azapi_resource_id" "test" {
  type        = "Microsoft.Network/virtualNetworks/subnets@2023-04-01"
  resource_id = "/subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resourceGroupName/providers/Microsoft.Network/virtualNetworks/vnetName/subnets/subnetName"
}
`
}

func (r IdDataSource) tenants() string {
	return `
data "azapi_resource_id" "test" {
  type        = "Microsoft.Resources/tenants@2021-04-01"
  resource_id = "/"
}


`
}
