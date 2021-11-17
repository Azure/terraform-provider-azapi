package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/location"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type GenericDataSource struct{}

func TestAccGenericDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
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
				check.That(data.ResourceName).Key("location").HasValue(location.LocationNormalize(data.LocationPrimary)),
				check.That(data.ResourceName).Key("tags.%").HasValue("1"),
			),
		},
	})
}

func (r GenericDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azurerm-restapi_resource" "test" {
  resource_id = azurerm-restapi_resource.test.resource_id
  type        = azurerm-restapi_resource.test.type
}
`, GenericResource{}.complete(data))
}
