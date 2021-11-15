package services_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/acceptance"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/acceptance/check"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/clients"
	"github.com/ms-henglu/terraform-provider-azurermg/internal/services/parse"
	"github.com/ms-henglu/terraform-provider-azurermg/utils"
)

type GenericResource struct{}

func ignoredProperties() []string {
	return []string{"body", "create_method", "update_method"}
}

func TestAccGenericResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccGenericResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_completeBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurermg_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.identityNone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.identityUserAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func (GenericResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	id, err := parse.ResourceID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, response, err := client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err != nil {
		if response.StatusCode == http.StatusNotFound {
			exist := false
			return &exist, nil
		}
		return nil, fmt.Errorf("retrieving %s: %+v", id, err)
	}
	exist := len(utils.GetId(resp)) != 0
	return &exist, nil
}

func (r GenericResource) basic(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false
}

resource "azurermg_resource" "test" {
  resource_id = "${azurerm_container_registry.test.id}/scopeMaps/acctest%[2]s"
  type        = "Microsoft.ContainerRegistry/registries/scopeMaps@2020-11-01-preview"
  body        = <<BODY
   {
      "properties": {
        "description": "Developer Scopes",
        "actions": [
          "repositories/testrepo/content/read"
        ]
      }
    }
  BODY
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurermg_resource" "import" {
  resource_id = azurermg_resource.test.resource_id
  type        = azurermg_resource.test.type
  body        = <<BODY
   {
      "properties": {
        "description": "Developer Scopes",
        "actions": [
          "repositories/testrepo/content/read"
        ]
      }
    }
  BODY
}
`, r.basic(data))
}

func (r GenericResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurermg_resource" "test" {
  resource_id = "${azurerm_resource_group.test.id}/providers/Microsoft.ContainerRegistry/registries/acctest%[2]s"
  type        = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  location    = "%[3]s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  body = <<BODY
    {
      "sku": {
        "name": "Standard"
      },
      "properties": {
        "adminUserEnabled": true
      }
    }
  BODY

  tags = {
    "Key" = "Value"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) completeBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurermg" {
  schema_validation_enabled = false
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurermg_resource" "test" {
  resource_id = "${azurerm_resource_group.test.id}/providers/Microsoft.ContainerRegistry/registries/acctest%[2]s"
  type        = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"

  body = <<BODY
    {
      "location": "${azurerm_resource_group.test.location}",
      "identity": {
		"type": "systemAssigned"
      },
      "sku": {
        "name": "Standard"
      },
      "properties": {
        "adminUserEnabled": true
      },
      "tags": {
        "key":"value"
      }
    }
  BODY
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identityNone(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurermg_resource" "test" {
  resource_id = "${azurerm_resource_group.test.id}/providers/Microsoft.ContainerRegistry/registries/acctest%[2]s"
  type        = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  location    = "%[3]s"
  body        = <<BODY
    {
      "sku": {
        "name": "Standard"
      },
      "properties": {
        "adminUserEnabled": true
      }
    }
  BODY
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurermg_resource" "test" {
  resource_id = "${azurerm_resource_group.test.id}/providers/Microsoft.ContainerRegistry/registries/acctest%[2]s"
  type        = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  location    = "%[3]s"
  identity {
    type = "SystemAssigned"
  }
  body = <<BODY
    {
      "sku": {
        "name": "Standard"
      },
      "properties": {
        "adminUserEnabled": true
      }
    }
  BODY
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identityUserAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurermg_resource" "test" {
  resource_id = "${azurerm_resource_group.test.id}/providers/Microsoft.ContainerRegistry/registries/acctest%[2]s"
  type        = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  location    = "%[3]s"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }
  body = <<BODY
    {
      "sku": {
        "name": "Standard"
      },
      "properties": {
        "adminUserEnabled": true
      }
    }
  BODY

  tags = {
    "Key" = "Value"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (GenericResource) template(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-%[1]d"
  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomStringOfLength(10))
}
