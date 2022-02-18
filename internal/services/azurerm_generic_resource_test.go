package services_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/clients"
	"github.com/Azure/terraform-provider-azurerm-restapi/internal/services/parse"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type GenericResource struct{}

func TestAccGenericResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
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
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
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
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
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
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
		{
			Config: r.identityUserAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
		{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
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
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
		{
			Config: r.defaultTagOverrideInBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
		{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
		{
			Config: r.defaultTagOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
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
				check.That(data.ResourceName).Key("tags").DoesNotExist(),
				check.That(data.ResourceName).Key("location").IsEmpty(),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
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
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
		{
			Config: r.defaultLocationOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationSecondary)),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
	})
}

func TestAccGenericResource_subscriptionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subscriptionScope(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
	})
}

func TestAccGenericResource_extensionScope(t *testing.T) {
	t.Skip(`The service principle does not have authorization to perform action 'Microsoft.Authorization/locks/write'`)
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionScope(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
	})
}

func TestAccGenericResource_ignoreMissingProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreMissingProperty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
	})
}

func TestAccGenericResource_ignoreCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm-restapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreCasing(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(r.ImportIdFunc, r.importStateCheckFunc),
	})
}

func (GenericResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.NewResourceID(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	_, _, err = client.NewResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion)
	if err == nil {
		b := true
		return &b, nil
	}
	var responseErr *azcore.ResponseError
	if errors.As(err, &responseErr) && responseErr.StatusCode == http.StatusNotFound {
		b := false
		return &b, nil
	}
	return nil, fmt.Errorf("checking for presence of existing %s: %+v", id, err)
}

func (GenericResource) ImportIdFunc(tfState *terraform.State) (string, error) {
	state := tfState.RootModule().Resources["azurerm-restapi_resource.test"].Primary
	resourceType := state.Attributes["type"]
	id, err := parse.NewResourceID(state.ID, resourceType)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?api-version=%s", id.AzureResourceId, id.ApiVersion), nil
}

func (GenericResource) importStateCheckFunc(states []*terraform.InstanceState) error {
	if len(states) != 1 {
		return fmt.Errorf("expect states length is 1, but got %d", len(states))
	}
	state := states[0]
	props := []string{"name", "parent_id", "type", "id", "body"}
	for _, prop := range props {
		if len(state.Attributes[prop]) == 0 {
			return fmt.Errorf("expect `%s` is not empty", prop)
		}
	}
	return nil
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
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azurerm-restapi_resource" "test" {
  name                      = "acctest%[2]s"
  parent_id                 = azurerm_resource_group.test.id
  type                      = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  schema_validation_enabled = false
  body                      = <<BODY
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

func (GenericResource) subscriptionScope(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azurerm-restapi_resource" "test" {
  type      = "Microsoft.Resources/resourceGroups@2021-04-01"
  name      = "acctestRG-%[1]d"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"

  location = "%[2]s"
}
`, data.RandomInteger, data.LocationPrimary, data.RandomStringOfLength(10))
}

// nolint staticcheck
func (r GenericResource) extensionScope(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm-restapi_resource" "locks" {
  type      = "Microsoft.Authorization/locks@2015-01-01"
  name      = "acctest-%[2]d"
  parent_id = azurerm_resource_group.test.id

  body = jsonencode({
    properties = {
      level = "CanNotDelete"
    }
  })
}
`, r.template(data), data.RandomInteger)
}

func (r GenericResource) ignoreMissingProperty(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm-restapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring/storages@2021-09-01-preview"
  name      = "acctest-ss-%[2]d"
  parent_id = azurerm_spring_cloud_service.test.id

  body = jsonencode({
    properties = {
      accountKey  = azurerm_storage_account.test.primary_access_key
      accountName = azurerm_storage_account.test.name
      storageType = "StorageAccount"
    }
  })

  ignore_missing_property_enabled = true
}
`, r.template(data), data.RandomInteger, data.RandomStringOfLength(10))
}

func (r GenericResource) ignoreCasing(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azurerm_storage_account" "test" {
  name                     = "acctestsa%[3]s"
  resource_group_name      = azurerm_resource_group.test.name
  location                 = azurerm_resource_group.test.location
  account_tier             = "Standard"
  account_replication_type = "GRS"
}
resource "azurerm_spring_cloud_service" "test" {
  name                = "acctest-sc-%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azurerm-restapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring/storages@2021-09-01-preview"
  name      = "acctest-ss-%[2]d"
  parent_id = azurerm_spring_cloud_service.test.id

  body = jsonencode({
    properties = {
      accountKey  = azurerm_storage_account.test.primary_access_key
      accountName = azurerm_storage_account.test.name
      storageType = "storageaccount"
    }
  })

  schema_validation_enabled       = false
  ignore_casing_enabled           = true
  ignore_missing_property_enabled = true
}
`, r.template(data), data.RandomInteger, data.RandomStringOfLength(10))
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
