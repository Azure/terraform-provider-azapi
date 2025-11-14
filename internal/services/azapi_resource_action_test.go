package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type AzapiResourceActionActionTest struct{}

func TestAccAzapiResourceActionAction_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := AzapiResourceActionActionTest{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config:            r.basic(data),
			ConfigStateChecks: []statecheck.StateCheck{},
		},
	})
}

func TestAccAzapiResourceActionAction_providerAction(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := AzapiResourceActionActionTest{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.providerAction(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"terraform_data.trigger",
					tfjsonpath.New("input"),
					knownvalue.StringExact("null"),
				),
			},
		},
	})
}

func TestAccAzapiResourceActionAction_withQueryParameters(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := AzapiResourceActionActionTest{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.withQueryParameters(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"terraform_data.trigger",
					tfjsonpath.New("input"),
					knownvalue.StringExact("null"),
				),
			},
		},
	})
}

func TestAccAzapiResourceActionAction_registerResourceProvider(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := AzapiResourceActionActionTest{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.registerResourceProvider(),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"terraform_data.trigger",
					tfjsonpath.New("input"),
					knownvalue.StringExact("null"),
				),
			},
		},
	})
}

func TestAccAzapiResourceActionAction_multipleTriggers(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource_action", "test")
	r := AzapiResourceActionActionTest{}

	data.DataSourceTest(t, []resource.TestStep{
		{
			Config: r.multipleTriggers(data),
			ConfigStateChecks: []statecheck.StateCheck{
				statecheck.ExpectKnownValue(
					"terraform_data.trigger1",
					tfjsonpath.New("input"),
					knownvalue.StringExact("null"),
				),
				statecheck.ExpectKnownValue(
					"terraform_data.trigger2",
					tfjsonpath.New("input"),
					knownvalue.StringExact("null"),
				),
			},
		},
	})
}

func (r AzapiResourceActionActionTest) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

action "azapi_resource_action" "test" {
  config {
    type        = "Microsoft.Automation/automationAccounts@2021-06-22"
    resource_id = azapi_resource.test.id
    action      = "agentRegistrationInformation/regenerateKey"
    body = {
      keyName = "primary"
    }
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = "acctest%[2]s"
  parent_id = azapi_resource.resourceGroup.id

  location = "%[3]s"

  body = {
    properties = {
      sku = {
        name = "Basic"
      }
    }
  }
}

resource "azapi_resource_action" "listKeys" {
  type        = "Microsoft.Automation/automationAccounts@2021-06-22"
  resource_id = azapi_resource.test.id
  action      = "listKeys"
}

resource "terraform_data" "test" {
  input = azapi_resource_action.listKeys.output
  lifecycle {
    action_trigger {
      events  = [after_create]
      actions = [action.azapi_resource_action.test]
    }
  }
}
`, GenericResource{}.template(data), data.RandomString, data.LocationPrimary)
}

func (r AzapiResourceActionActionTest) providerAction(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azapi_client_config" "current" {}

action "azapi_resource_action" "test" {
  config {
    type        = "Microsoft.Cache@2023-04-01"
    resource_id = "${data.azapi_client_config.current.subscription_resource_id}/providers/Microsoft.Cache"
    action      = "CheckNameAvailability"
    method      = "POST"
    headers = {
      "X-Custom-Header" = "test-value"
    }
    locks = [data.azapi_client_config.current.subscription_id]
    body = {
      type = "Microsoft.Cache/Redis"
      name = "test-%[1]s"
    }
  }
}

resource "terraform_data" "trigger" {
  input = "null"
  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.azapi_resource_action.test]
    }
  }
}
`, data.RandomString)
}

func (r AzapiResourceActionActionTest) withQueryParameters(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azapi_client_config" "current" {}

action "azapi_resource_action" "test" {
  config {
    type        = "Microsoft.Resources/subscriptions@2021-04-01"
    resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
    action      = "providers"
    method      = "GET"
    query_parameters = {
      "$expand" = ["metadata"]
    }
  }
}

resource "terraform_data" "trigger" {
  input = "null"
  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.azapi_resource_action.test]
    }
  }
}
`)
}

func (r AzapiResourceActionActionTest) registerResourceProvider() string {
	return `
data "azapi_client_config" "current" {}

action "azapi_resource_action" "test" {
  config {
    type        = "Microsoft.Resources/providers@2021-04-01"
    resource_id = "${data.azapi_client_config.current.subscription_resource_id}/providers/Microsoft.Compute"
    action      = "register"
    method      = "POST"
  }
}

resource "terraform_data" "trigger" {
  input = "null"
  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.azapi_resource_action.test]
    }
  }
}
`
}

func (r AzapiResourceActionActionTest) multipleTriggers(data acceptance.TestData) string {
	return fmt.Sprintf(`
data "azapi_client_config" "current" {}

action "azapi_resource_action" "test1" {
  config {
    type        = "Microsoft.Cache@2023-04-01"
    resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Cache"
    action      = "CheckNameAvailability"
    body = {
      type = "Microsoft.Cache/Redis"
      name = "test1-%[1]s"
    }
  }
}

action "azapi_resource_action" "test2" {
  config {
    type        = "Microsoft.Cache@2023-04-01"
    resource_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Cache"
    action      = "CheckNameAvailability"
    body = {
      type = "Microsoft.Cache/Redis"
      name = "test2-%[1]s"
    }
  }
}

resource "terraform_data" "trigger1" {
  input = "null"
  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.azapi_resource_action.test1]
    }
  }
}

resource "terraform_data" "trigger2" {
  input = "null"
  lifecycle {
    action_trigger {
      events  = [before_create]
      actions = [action.azapi_resource_action.test2]
    }
  }
}
`, data.RandomString)
}
