package services_test

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/acceptance"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/acceptance/check"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/azure/location"
)

type GenericDataSource struct{}

func TestAccGenericDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_resource", "test")
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

data "azurermg_resource" "test" {
  resource_id = azurermg_resource.test.resource_id
  api_version = azurermg_resource.test.api_version
}
`, GenericResource{}.complete(data))
}
