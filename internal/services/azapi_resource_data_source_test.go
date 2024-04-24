package services_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type GenericDataSource struct{}

func TestAccGenericDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource", "test")
	r := GenericDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccGenericDataSource_withResourceId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource", "test")
	r := GenericDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.withResourceId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("identity.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.type").HasValue("SystemAssigned, UserAssigned"),
				check.That(data.ResourceName).Key("identity.0.identity_ids.#").HasValue("1"),
				check.That(data.ResourceName).Key("identity.0.principal_id").Exists(),
				check.That(data.ResourceName).Key("identity.0.tenant_id").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func TestAccGenericDataSource_defaultParentId(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource", "test")
	r := GenericDataSource{}

	subscriptionId := os.Getenv("ARM_SUBSCRIPTION_ID")
	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.defaultParentId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("parent_id").HasValue(fmt.Sprintf("/subscriptions/%s", subscriptionId)),
			),
		},
	})
}

func (r GenericDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  name      = azapi_resource.test.name
  parent_id = azapi_resource.test.parent_id
  type      = azapi_resource.test.type
}
`, GenericResource{}.complete(data))
}

func (r GenericDataSource) defaultParentId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  type = "Microsoft.Resources/resourceGroups@2024-03-01"
  name = azapi_resource.test.name
}
`, GenericResource{}.defaultParentId(data))
}

func (r GenericDataSource) withResourceId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  type        = azapi_resource.test.type
  resource_id = azapi_resource.test.id
}
`, GenericResource{}.complete(data))
}
