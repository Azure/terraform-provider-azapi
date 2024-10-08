package services_test

import (
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ListDataSource struct{}

func TestAccListDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_list", "test")
	r := ListDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.value.#").Exists(),
			),
		},
	})
}

func TestAccListDataSource_paging(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_list", "test")
	r := ListDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.paging(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccListDataSource_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_list", "test")
	r := ListDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.headers(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccListDataSource_queryParameter(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_list", "test")
	r := ListDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.queryParameter(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccListDataSource_defaultOutput(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_list", "test")
	r := ListDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.defaultOutput(),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.value.#").Exists(),
			),
		},
	})
}

func (r ListDataSource) basic() string {
	return `
data "azapi_client_config" "current" {}

data "azapi_resource_list" "test" {
  type                   = "Microsoft.Resources/resourceGroups@2024-03-01"
  parent_id              = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = ["*"]
}
`
}

func (r ListDataSource) paging() string {
	return `
data "azapi_client_config" "current" {}

data "azapi_resource_list" "test" {
  type                   = "Microsoft.Authorization/policyDefinitions@2021-06-01"
  parent_id              = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = ["*"]
}
`
}

func (r ListDataSource) headers() string {
	return `
data "azapi_client_config" "current" {}

data "azapi_resource_list" "test" {
  type                   = "Microsoft.Authorization/policyDefinitions@2021-06-01"
  parent_id              = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = ["*"]
  headers = {
    "header1" = "value1"
  }
}
`
}

func (r ListDataSource) queryParameter() string {
	return `
data "azapi_client_config" "current" {}

data "azapi_resource_list" "test" {
  type      = "Microsoft.Authorization/policyDefinitions@2021-06-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  query_parameters = {
    "$filter" = ["policyType eq 'BuiltIn'"]
  }
  response_export_values = ["*"]
}
`
}

func (r ListDataSource) defaultOutput() string {
	return `
data "azapi_client_config" "current" {}

data "azapi_resource_list" "test" {
  type      = "Microsoft.Resources/resourceGroups@2024-03-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
}
`
}
