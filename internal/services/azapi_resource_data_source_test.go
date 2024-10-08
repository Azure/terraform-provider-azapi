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
				check.That(data.ResourceName).Key("output.properties.%").Exists(),
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
				check.That(data.ResourceName).Key("output.properties.%").Exists(),
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

func TestAccGenericDataSource_withRetry(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource", "test")
	r := GenericDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.withRetry(data),
			ExternalProviders: map[string]resource.ExternalProvider{
				"time": {
					Source:            "hashicorp/time",
					VersionConstraint: "0.12.0",
				},
			},
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("location").Exists(),
			),
		},
	})
}

func TestAccGenericDataSource_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource", "test")
	r := GenericDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.headers(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccGenericDataSource_queryParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource", "test")
	r := GenericDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.queryParameter(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccGenericDataSource_defaultOutput(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource", "test")
	r := GenericDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.defaultOutput(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.properties.automationHybridServiceUrl").Exists(),
			),
		},
	})
}

func (r GenericDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  name                   = azapi_resource.test.name
  parent_id              = azapi_resource.test.parent_id
  type                   = azapi_resource.test.type
  response_export_values = ["*"]
}
`, GenericResource{}.complete(data))
}

func (r GenericDataSource) defaultParentId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  type                   = "Microsoft.Resources/resourceGroups@2024-03-01"
  name                   = azapi_resource.test.name
  response_export_values = ["*"]
}
`, GenericResource{}.defaultParentId(data))
}

func (r GenericDataSource) withResourceId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  type                   = azapi_resource.test.type
  resource_id            = azapi_resource.test.id
  response_export_values = ["*"]
}
`, GenericResource{}.complete(data))
}

func (r GenericDataSource) withRetry(data acceptance.TestData) string {
	return fmt.Sprintf(`
resource "time_sleep" "wait_30_seconds" {
  create_duration = "30s"
}

data "azapi_client_config" "this" {}

resource "azapi_resource" "test" {
  name       = "acctestRG-%[1]d"
  type       = "Microsoft.Resources/resourceGroups@2024-03-01"
  parent_id  = "/subscriptions/${data.azapi_client_config.this.subscription_id}"
  location   = "%[2]s"
  depends_on = [time_sleep.wait_30_seconds]
}

resource "terraform_data" "read_data_source_during_apply" {
  input = "acctestRG-%[1]d"
}

data "azapi_resource" "test" {
  type      = "Microsoft.Resources/resourceGroups@2024-03-01"
  name      = terraform_data.read_data_source_during_apply.output
  parent_id = "/subscriptions/${data.azapi_client_config.this.subscription_id}"

  retry = {
    error_message_regex = ["ResourceGroupNotFound"]
  }

  timeouts {
    read = "2m"
  }
}
`, data.RandomInteger, data.LocationPrimary, data.RandomInteger)
}

func (r GenericDataSource) headers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  name                   = azapi_resource.test.name
  parent_id              = azapi_resource.test.parent_id
  type                   = azapi_resource.test.type
  response_export_values = ["*"]
  headers = {
    "header1" = "value1"
  }
}
`, GenericResource{}.complete(data))
}

func (r GenericDataSource) queryParameter(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  name                   = azapi_resource.test.name
  parent_id              = azapi_resource.test.parent_id
  type                   = azapi_resource.test.type
  response_export_values = ["*"]
  query_parameters = {
    "query1" = ["value1"]
  }
}
`, GenericResource{}.complete(data))
}

func (r GenericDataSource) defaultOutput(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource" "test" {
  name      = azapi_resource.test.name
  parent_id = azapi_resource.test.parent_id
  type      = azapi_resource.test.type
}
`, GenericResource{}.defaultOutput(data))
}
