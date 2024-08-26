package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ActionDataSource struct{}

func TestAccActionDataSource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionDataSource_providerPermissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerPermissions(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionDataSource_providerAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionDataSource_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.headers(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccActionDataSource_queryParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "data.azapi_resource_action", "test")
	r := ActionDataSource{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.queryParameters(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r ActionDataSource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource_action" "test" {
  type                   = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id            = azapi_resource.test.id
  action                 = "listKeys"
  response_export_values = ["*"]
}
`, GenericResource{}.defaultTag(data))
}

func (r ActionDataSource) providerPermissions() string {
	return `

data "azapi_client_config" "current" {}

data "azapi_resource_action" "test" {
  type        = "Microsoft.Resources/providers@2021-04-01"
  resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Network"
  action      = "providerPermissions"
  method      = "GET"
}
`
}

func (r ActionDataSource) providerAction() string {
	return `
data "azapi_resource_action" "test" {
  type        = "Microsoft.ResourceGraph@2020-04-01-preview"
  resource_id = "/providers/Microsoft.ResourceGraph"
  action      = "resources"
  body = {
    query = "resources| limit 1"
  }
  response_export_values = ["*"]
}
`
}

func (r ActionDataSource) headers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

data "azapi_resource_action" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azapi_resource.test.id
  action      = "listKeys"
  headers = {
    "header1" = "value1"
  }
  response_export_values = ["*"]
}
`, GenericResource{}.defaultTag(data))
}

func (r ActionDataSource) queryParameters() string {
	return `
data "azapi_client_config" "current" {}

data "azapi_resource_action" "test" {
  type        = "Microsoft.Authorization@2021-06-01"
  resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization"
  action      = "policyDefinitions"
  method      = "GET"
  query_parameters = {
    "$filter" = ["policyType eq 'BuiltIn'"]
  }
  response_export_values = ["*"]
}`
}
