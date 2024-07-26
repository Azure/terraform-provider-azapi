package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type GenericUpdateResource struct{}

func TestAccGenericUpdateResource_automationAccount(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.automationAccount(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_dynamicSchema(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dynamicSchema(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_withNameParentId(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.automationAccountWithNameParentId(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_siteConfigSlotConfigNames(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.siteConfigSlotConfigNames(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_locks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.locks(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func TestAccGenericUpdateResource_ignoreChanges(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreChanges(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_ignoreChangesArray(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreChangesArray(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_ignoreOrderInArray(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreOrderInArray(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_timeouts(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.timeouts(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccGenericUpdateResource_ignoreIDCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreIDCasing(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			Config: r.ignoreIDCasingUpdated(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (r GenericUpdateResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.ResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	resp, err := client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if utils.ResponseErrorWasNotFound(err) {
			exist := false
			return &exist, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	exist := utils.GetId(resp) != nil
	return &exist, nil
}

func (r GenericUpdateResource) automationAccount(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azurerm_automation_account.test.id
  body = jsonencode({
    properties = {
      publicNetworkAccess = true
    }
  })
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) dynamicSchema(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azurerm_automation_account.test.id
  body = {
    properties = {
      publicNetworkAccess = true
    }
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) automationAccountWithNameParentId(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2023-11-01"
  name      = azurerm_automation_account.test.name
  parent_id = azurerm_resource_group.test.id
  body = jsonencode({
    properties = {
      publicNetworkAccess = true
    }
  })
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) siteConfigSlotConfigNames(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_app_service_plan" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku {
    tier = "Standard"
    size = "S1"
  }
}

resource "azurerm_app_service" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  app_service_plan_id = azurerm_app_service_plan.test.id

  site_config {
    dotnet_framework_version = "v4.0"
    scm_type                 = "LocalGit"
  }

  app_settings = {
    "SOME_KEY" = "some-value"
  }

  connection_string {
    name  = "Database"
    type  = "SQLServer"
    value = "Server=some-server.mydomain.com;Integrated Security=SSPI"
  }
}

resource "azapi_update_resource" "test" {
  type      = "Microsoft.Web/sites/config@2021-03-01"
  name      = "slotConfigNames"
  parent_id = azurerm_app_service.test.id
  body = jsonencode({
    properties = {
      connectionStringNames   = ["test1", "test2"]
      appSettingNames         = ["test3"]
      azureStorageConfigNames = ["test4"]
    }
  })
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) locks(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test1" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azurerm_automation_account.test.id
  body = jsonencode({
    properties = {
      publicNetworkAccess = true
    }
  })
  locks = [azurerm_automation_account.test.id]
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azurerm_automation_account.test.id
  body = jsonencode({
    properties = {
      publicNetworkAccess = true
    }
  })
  locks = [azurerm_automation_account.test.id]
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) ignoreChanges(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azurerm_automation_account.test.id
  body = jsonencode({
    properties = {
      sku = {
        name = "Free"
      }
    }
  })

  ignore_body_changes = ["properties.sku.name"]
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) ignoreChangesArray(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_virtual_network" "test" {
  name                = "acctest%[2]d"
  address_space       = ["10.0.0.0/16"]
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm_subnet" "test" {
  name                 = "default"
  resource_group_name  = azurerm_resource_group.test.name
  virtual_network_name = azurerm_virtual_network.test.name
  address_prefixes     = ["10.0.1.0/24"]
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Network/virtualNetworks@2022-07-01"
  resource_id = azurerm_virtual_network.test.id
  body = jsonencode({
    properties = {
      subnets = [
        {
          name = "second"
          properties = {
            addressPrefix = "10.0.2.0/24"
          }
        }
      ]
    }
  })

  ignore_body_changes = ["properties.subnets"]
  depends_on          = [azurerm_subnet.test]
}
`, r.template(data), data.RandomInteger)
}

func (r GenericUpdateResource) ignoreOrderInArray(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "vnet" {
  type      = "Microsoft.Network/virtualNetworks@2023-09-01"
  parent_id = azurerm_resource_group.test.id
  name      = "acctest%[2]d"
  location  = azurerm_resource_group.test.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
        {
          name = "first"
          properties = {
            addressPrefix = "10.0.3.0/24"
          }
        },
        {
          name = "second"
          properties = {
            addressPrefix = "10.0.4.0/24"
          }
        }
      ]
    }
  }
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Network/virtualNetworks@2022-07-01"
  resource_id = azapi_resource.vnet.id
  body = {
    properties = {
      subnets = [
        {
          name = "second"
          properties = {
            addressPrefix = "10.0.4.0/24"
          }
        },
        {
          name = "first"
          properties = {
            addressPrefix = "10.0.3.0/24"
          }
        }
      ]
    }
  }
}
`, r.template(data), data.RandomInteger)
}

func (r GenericUpdateResource) ignoreIDCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_update_resource" "test" {
  type        = "Microsoft.SignalRService/WebPubSub@2024-04-01-preview"
  resource_id = azurerm_web_pubsub.test.id

  body = {
    properties = {
      networkACLs = {
        defaultAction = "Deny"
        publicNetwork = {
          allow = ["ClientConnection"]
        }
        ipRules = [{
          value  = "0.0.0.0/0"
          action = "Allow"
        }]
      }
    }
  }
}
`, r.templateForIDCasing(data))
}

func (r GenericUpdateResource) ignoreIDCasingUpdated(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_update_resource" "test" {
  type        = "Microsoft.SignalRService/WebPubSub@2024-04-01-preview"
  resource_id = azurerm_web_pubsub.test.id

  body = {
    properties = {
      networkACLs = {
        defaultAction = "Deny"
        publicNetwork = {
          allow = ["ClientConnection", "RESTAPI"]
        }
        ipRules = [{
          value  = "0.0.0.0/0"
          action = "Allow"
        }]
      }
    }
  }
}
`, r.templateForIDCasing(data))
}

func (GenericUpdateResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomString)
}

func (r GenericUpdateResource) templateForIDCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_web_pubsub" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name

  sku      = "Standard_S1"
  capacity = 1

  public_network_access_enabled = false

  live_trace {
    enabled                   = true
    messaging_logs_enabled    = true
    connectivity_logs_enabled = false
  }

  identity {
    type = "SystemAssigned"
  }
}
`, r.template(data), data.RandomString)
}

func (r GenericUpdateResource) timeouts(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_automation_account" "test" {
  name                = "acctest-%[2]s"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku_name            = "Basic"
}

resource "azapi_update_resource" "test" {
  type        = "Microsoft.Automation/automationAccounts@2023-11-01"
  resource_id = azurerm_automation_account.test.id
  body = jsonencode({
    properties = {
      publicNetworkAccess = true
    }
  })
  timeouts {
    create = "10m"
    update = "10m"
    delete = "10m"
    read   = "10m"
  }
}
`, r.template(data), data.RandomString)
}
