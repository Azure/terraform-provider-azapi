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

type DataPlaneResource struct{}

func TestAccDataPlaneResource_appConfigKeyValues(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.appConfigKeyValues(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_dynamicSchema(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.dynamicSchema(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_purviewClassification(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purviewClassification(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_purviewCollection(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.purviewCollection(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_keyVaultIssuer(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.keyVaultIssuer(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func TestAccDataPlaneResource_iotAppsUser(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.iotAppsUser(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
	})
}

func (DataPlaneResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.DataPlaneResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	_, err = client.DataPlaneClient.Get(ctx, id)
	if err == nil {
		b := true
		return &b, nil
	}
	if utils.ResponseErrorWasNotFound(err) {
		b := false
		return &b, nil
	}
	return nil, fmt.Errorf("checking for presence of existing %s: %+v", id, err)
}

func (r DataPlaneResource) appConfigKeyValues(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azurerm_app_configuration" "appconf" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

data "azurerm_client_config" "current" {}

resource "azurerm_role_assignment" "test" {
  scope                = azurerm_app_configuration.appconf.id
  role_definition_name = "App Configuration Data Owner"
  principal_id         = data.azurerm_client_config.current.object_id
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
  parent_id = replace(azurerm_app_configuration.appconf.endpoint, "https://", "")
  name      = "mykey"
  body = jsonencode(
    {
      content_type = ""
      value        = "myvalue"
    }
  )

  depends_on = [
    azurerm_role_assignment.test,
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) dynamicSchema(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azurerm_purview_account" "example" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Purview/accounts/Scanning/classificationrules@2022-07-01-preview"
  parent_id = replace(azurerm_purview_account.example.scan_endpoint, "https://", "")
  name      = "acctest%[2]s"
  payload = {
    kind = "Custom"
    properties = {
      description        = "Let's put a cool desc here"
      classificationName = "MICROSOFT.FINANCIAL.AUSTRALIA.BANK_ACCOUNT_NUMBER"
      columnPatterns = [
        {
          pattern = "^data$"
          kind    = "Regex"
        }
      ]
      dataPatterns = [
        {
          pattern = "^[0-9]{2}-[0-9]{4}-[0-9]{6}-[0-9]{3}$"
          kind    = "Regex"
        }
      ]
      minimumPercentageMatch = 60
      ruleStatus             = "Enabled"
    }
  }
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) purviewClassification(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azurerm_purview_account" "example" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Purview/accounts/Scanning/classificationrules@2022-07-01-preview"
  parent_id = replace(azurerm_purview_account.example.scan_endpoint, "https://", "")
  name      = "acctest%[2]s"
  body = jsonencode({
    kind = "Custom"
    properties = {
      description        = "Let's put a cool desc here"
      classificationName = "MICROSOFT.FINANCIAL.AUSTRALIA.BANK_ACCOUNT_NUMBER"
      columnPatterns = [
        {
          pattern = "^data$"
          kind    = "Regex"
        }
      ]
      dataPatterns = [
        {
          pattern = "^[0-9]{2}-[0-9]{4}-[0-9]{6}-[0-9]{3}$"
          kind    = "Regex"
        }
      ]
      minimumPercentageMatch = 60
      ruleStatus             = "Enabled"
    }
  })
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) purviewCollection(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azurerm_purview_account" "example" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location

  identity {
    type = "SystemAssigned"
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Purview/accounts/Account/collections@2019-11-01-preview"
  parent_id = "${azurerm_purview_account.example.name}.purview.azure.com"
  name      = "defaultResourceSetRuleConfig"
  body = jsonencode({
    friendlyName = "Finance"
  })
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) keyVaultIssuer(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm_resource_group" "example" {
  name     = "acctest%[2]s"
  location = "%[1]s"
}

resource "azurerm_key_vault" "example" {
  name                        = "acctest%[2]s"
  location                    = azurerm_resource_group.example.location
  resource_group_name         = azurerm_resource_group.example.name
  enabled_for_disk_encryption = true
  tenant_id                   = data.azurerm_client_config.current.tenant_id
  soft_delete_retention_days  = 7
  purge_protection_enabled    = false

  sku_name = "standard"
}

resource "azapi_resource_action" "add_accesspolicy" {
  type        = "Microsoft.KeyVault/vaults/accessPolicies@2023-02-01"
  resource_id = "${azurerm_key_vault.example.id}/accessPolicies/add"
  method      = "PUT"
  body = jsonencode({
    properties = {
      accessPolicies = [{
        tenantId = data.azurerm_client_config.current.tenant_id
        objectId = data.azurerm_client_config.current.object_id
        permissions = {
          certificates = ["managecontacts", "getissuers", "setissuers", "deleteissuers"]
        }
      }]
    }
  })
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.KeyVault/vaults/certificates/issuers@7.4"
  parent_id = replace(azurerm_key_vault.example.vault_uri, "https://", "")
  name      = "acctest%[2]s"
  body = jsonencode({
    provider = "Test"
    credentials = {
      account_id = "keyvaultuser"
    }
    org_details = {
      admin_details = [
        {
          first_name = "John"
          last_name  = "Doe"
          email      = "admin@microsoft.com"
          phone      = "4255555555"
        }
      ]
    }
  })
  depends_on = [
    azapi_resource_action.add_accesspolicy
  ]
}
`, data.LocationPrimary, data.RandomString)
}

func (r DataPlaneResource) iotAppsUser(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "acctest%[2]s"
  location = "%[1]s"
}
resource "azurerm_iotcentral_application" "example" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sub_domain          = "acctest%[2]s"
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.IoTCentral/IoTApps/users@2022-07-31"
  parent_id = "${azurerm_iotcentral_application.example.sub_domain}.azureiotcentral.com"
  name      = "acctest%[2]s"
  body = jsonencode({
    type = "email"
    roles = [
      {
        role = "ae2c9854-393b-4f97-8c42-479d70ce626e"
      }
    ]
    email = "user5@contoso.com"
  })
}

`, data.LocationPrimary, data.RandomString)
}
