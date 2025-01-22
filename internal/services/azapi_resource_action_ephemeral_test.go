package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

type ActionEphemeral struct{}

func TestAccEphemeral_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azapi_resource_action", "test")
	r := ActionEphemeral{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.basic(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccEphemeral_providerPermissions(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azapi_resource_action", "test")
	r := ActionEphemeral{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerPermissions(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccEphemeral_providerAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azapi_resource_action", "test")
	r := ActionEphemeral{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccEphemeral_headers(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azapi_resource_action", "test")
	r := ActionEphemeral{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.headers(data),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func TestAccEphemeral_queryParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "ephemeral.azapi_resource_action", "test")
	r := ActionEphemeral{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.queryParameters(),
			Check:  resource.ComposeTestCheckFunc(),
		},
	})
}

func (r ActionEphemeral) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azapi_resource_action" "test" {
  type                   = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id            = azapi_resource.test.id
  action                 = "listKeys"
  response_export_values = ["*"]
}
`, GenericResource{}.defaultTag(data))
}

func (r ActionEphemeral) providerPermissions() string {
	return `

data "azapi_client_config" "current" {}

ephemeral "azapi_resource_action" "test" {
  type        = "Microsoft.Resources/providers@2021-04-01"
  resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Network"
  action      = "providerPermissions"
  method      = "GET"
}
`
}

func (r ActionEphemeral) providerAction() string {
	return `
ephemeral "azapi_resource_action" "test" {
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

func (r ActionEphemeral) headers(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

ephemeral "azapi_resource_action" "test" {
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

func (r ActionEphemeral) queryParameters() string {
	return `
data "azapi_client_config" "current" {}

ephemeral "azapi_resource_action" "test" {
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
