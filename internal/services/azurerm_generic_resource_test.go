package services_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/terraform-provider-azurerm-restapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/clients"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azurerm-restapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type GenericResource struct{}

func ignoredProperties() []string {
	return []string{"body"}
}

func TestAccGenericResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestAccGenericResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_completeBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.identityNone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.identityUserAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_defaultTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.defaultTagOverrideInBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.defaultTagOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_defaultsNotApplicable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultsNotApplicable(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("tags").DoesNotExist(),
				check.That(data.ResourceName).Key("location").DoesNotExist(),
			),
		},
		data.ImportStep(ignoredProperties()...),
	})
}

func TestAccGenericResource_defaultLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultLocation(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStep(ignoredProperties()...),
		{
			Config: r.defaultLocationOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationSecondary)),
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

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_container_registry.test.id
  type      = "Microsoft.ContainerRegistry/registries/scopeMaps@2020-11-01-preview"
  body      = <<BODY
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

resource "azurerm-restapi_resource" "import" {
  name      = azurerm-restapi_resource.test.name
  parent_id = azurerm-restapi_resource.test.parent_id
  type      = azurerm-restapi_resource.test.type
  body      = <<BODY
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

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  location  = "%[3]s"
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

provider "azurerm-restapi" {
  schema_validation_enabled = false
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"

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

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  location  = "%[3]s"
  body      = <<BODY
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

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  location  = "%[3]s"
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

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  location  = "%[3]s"
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

func (r GenericResource) defaultTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm-restapi" {
  default_tags = {
    key = "default"
  }
}

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultTagOverrideInBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm-restapi" {
  default_tags = {
    key = "default"
  }
}

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
    tags = {
      key = "override"
    }
  })

}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultTagOverrideInHcl(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm-restapi" {
  default_tags = {
    key = "default"
  }
}

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = azurerm_resource_group.test.location
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })

  tags = {
    key = "override"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultLocation(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm-restapi" {
  default_location = "%[3]s"
}

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultLocationOverrideInHcl(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azurerm-restapi" {
  default_location = "%[3]s"
}

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id
  type      = "Microsoft.Automation/automationAccounts@2020-01-13-preview"

  location = "%[4]s"
  identity {
    type = "SystemAssigned"
  }

  body = jsonencode({
    properties = {
      sku = {
        name = "Basic"
      }
    }
    tags = {
      key = "override"
    }
  })

}
`, r.template(data), data.RandomString, data.LocationPrimary, data.LocationSecondary)
}

func (r GenericResource) defaultsNotApplicable(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azurerm-restapi" {
  default_tags = {
    key = "default"
  }
  default_location = "%[3]s"
}

resource "azurerm_container_registry" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false
}

resource "azurerm-restapi_resource" "test" {
  name      = "acctest%[2]s"
  parent_id = azurerm_container_registry.test.id
  type      = "Microsoft.ContainerRegistry/registries/scopeMaps@2020-11-01-preview"
  body      = <<BODY
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
