package services_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/Azure/terraform-provider-azapi/internal/azure/location"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/Azure/terraform-provider-azapi/internal/services/parse"
	"github.com/Azure/terraform-provider-azapi/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

type GenericResource struct{}

func defaultIgnores() []string {
	return []string{"ignore_casing", "ignore_missing_property", "schema_validation_enabled", "body", "locks.#", "locks.0", "locks.1"}
}

func TestAccGenericResource_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.basic(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
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

func TestAccGenericResource_importWithApiVersion(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.importWithApiVersion(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		{
			ResourceName:            data.ResourceName,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateIdFunc:       r.ImportIdFunc,
			ImportStateVerifyIgnore: defaultIgnores(),
		},
	})
}

func TestAccGenericResource_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_completeBody(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.completeBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_identity(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.identityNone(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.identityUserAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.identitySystemAssigned(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.complete(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultTags(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultTagOverrideInBody(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultTag(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("default"),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultTagOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("tags.key").HasValue("override"),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultsNotApplicable(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
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
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_defaultLocation(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.defaultLocation(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStep(defaultIgnores()...),
		{
			Config: r.defaultLocationOverrideInHcl(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationSecondary)),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_subscriptionScope(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.subscriptionScope(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_extensionScope(t *testing.T) {
	t.Skip(`The service principle does not have authorization to perform action 'Microsoft.Authorization/locks/write'`)
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.extensionScope(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("location").HasValue(location.Normalize(data.LocationPrimary)),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_ignoreMissingProperty(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreMissingProperty(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_ignoreCasing(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.ignoreCasing(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_deleteLROEndsWithNotFoundError(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.deleteLROEndsWithNotFoundError(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func TestAccGenericResource_locks(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_resource", "test")
	r := GenericResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.locks(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(defaultIgnores()...),
	})
}

func (GenericResource) Exists(ctx context.Context, client *clients.Client, state *terraform.InstanceState) (*bool, error) {
	resourceType := state.Attributes["type"]
	id, err := parse.ResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return nil, err
	}

	_, err = client.ResourceClient.Get(ctx, id.AzureResourceId, id.ApiVersion)
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

func (GenericResource) ImportIdFunc(tfState *terraform.State) (string, error) {
	state := tfState.RootModule().Resources["azapi_resource.test"].Primary
	resourceType := state.Attributes["type"]
	id, err := parse.ResourceIDWithResourceType(state.ID, resourceType)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s?api-version=%s", id.AzureResourceId, id.ApiVersion), nil
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

resource "azapi_resource" "test" {
  type      = "Microsoft.ContainerRegistry/registries/scopeMaps@2022-02-01-preview"
  name      = "acctest%[2]s"
  parent_id = azurerm_container_registry.test.id

  body = jsonencode({
    properties = {
      description = "Developer Scopes"
      actions = [
        "repositories/testrepo/content/read"
      ]
    }
  })
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) requiresImport(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "import" {
  type      = azapi_resource.test.type
  name      = azapi_resource.test.name
  parent_id = azapi_resource.test.parent_id
  body = jsonencode({
    properties = {
      description = "Developer Scopes"
      actions = [
        "repositories/testrepo/content/read"
      ]
    }
  })
}
`, r.basic(data))
}

func (r GenericResource) importWithApiVersion(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_container_registry" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku                 = "Premium"
  admin_enabled       = false
}

resource "azapi_resource" "test" {
  type      = "Microsoft.ContainerRegistry/registries/scopeMaps@2020-11-01-preview"
  name      = "acctest%[2]s"
  parent_id = azurerm_container_registry.test.id

  body = jsonencode({
    properties = {
      description = "Developer Scopes"
      actions = [
        "repositories/testrepo/content/read"
      ]
    }
  })
}
`, r.template(data), data.RandomString)
}

func (r GenericResource) complete(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azapi_resource" "test" {
  type      = "Microsoft.ContainerRegistry/registries@2022-02-01-preview"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[3]s"
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  body = jsonencode({
    sku = {
      name = "Standard"
    }
    properties = {
      adminUserEnabled = true
    }
  })

  tags = {
    "Key" = "Value"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) completeBody(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

provider "azapi" {
}

resource "azurerm_user_assigned_identity" "test" {
  name                = "acctest%[2]s"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
}

resource "azapi_resource" "test" {
  name                      = "acctest%[2]s"
  parent_id                 = azurerm_resource_group.test.id
  type                      = "Microsoft.ContainerRegistry/registries@2022-02-01-preview"
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

resource "azapi_resource" "test" {
  type      = "Microsoft.ContainerRegistry/registries@2022-02-01-preview"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[3]s"

  body = jsonencode({
    sku = {
      name = "Standard"
    }
    properties = {
      adminUserEnabled = true
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) identitySystemAssigned(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s

resource "azapi_resource" "test" {
  type      = "Microsoft.ContainerRegistry/registries@2022-02-01-preview"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[3]s"
  identity {
    type = "SystemAssigned"
  }
  body = jsonencode({
    sku = {
      name = "Standard"
    }
    properties = {
      adminUserEnabled = true
    }
  })

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

resource "azapi_resource" "test" {
  type      = "Microsoft.ContainerRegistry/registries@2022-02-01-preview"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

  location = "%[3]s"
  identity {
    type         = "UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.test.id]
  }

  body = jsonencode({
    sku = {
      name = "Standard"
    }
    properties = {
      adminUserEnabled = true
    }
  })


  tags = {
    "Key" = "Value"
  }
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (r GenericResource) defaultTag(data acceptance.TestData) string {
	return fmt.Sprintf(`
%s
provider "azapi" {
  default_tags = {
    key = "default"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2022-08-08"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

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
provider "azapi" {
  default_tags = {
    key = "default"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2022-08-08"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

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
provider "azapi" {
  default_tags = {
    key = "default"
  }
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2022-08-08"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

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
provider "azapi" {
  default_location = "%[3]s"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2022-08-08"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

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
provider "azapi" {
  default_location = "%[3]s"
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Automation/automationAccounts@2022-08-08"
  name      = "acctest%[2]s"
  parent_id = azurerm_resource_group.test.id

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

provider "azapi" {
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

resource "azapi_resource" "test" {
  type      = "Microsoft.ContainerRegistry/registries/scopeMaps@2022-02-01-preview"
  name      = "acctest%[2]s"
  parent_id = azurerm_container_registry.test.id

  body = jsonencode({
    properties = {
      description = "Developer Scopes"
      actions = [
        "repositories/testrepo/content/read"
      ]
    }
  })
}
`, r.template(data), data.RandomString, data.LocationPrimary)
}

func (GenericResource) subscriptionScope(data acceptance.TestData) string {
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azapi_resource" "test" {
  type      = "Microsoft.Resources/resourceGroups@2022-09-01"
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

resource "azapi_resource" "locks" {
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

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring/storages@2022-11-01-preview"
  name      = "acctest-ss-%[2]d"
  parent_id = azurerm_spring_cloud_service.test.id

  body = jsonencode({
    properties = {
      accountKey  = azurerm_storage_account.test.primary_access_key
      accountName = azurerm_storage_account.test.name
      storageType = "StorageAccount"
    }
  })

  ignore_missing_property = true
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

resource "azapi_resource" "test" {
  type      = "Microsoft.AppPlatform/Spring/storages@2022-11-01-preview"
  name      = "acctest-ss-%[2]d"
  parent_id = azurerm_spring_cloud_service.test.id

  body = jsonencode({
    properties = {
      accountKey  = azurerm_storage_account.test.primary_access_key
      accountName = azurerm_storage_account.test.name
      storageType = "storageaccount"
    }
  })

  schema_validation_enabled = false
  ignore_casing             = true
  ignore_missing_property   = true
}
`, r.template(data), data.RandomInteger, data.RandomStringOfLength(10))
}

func (r GenericResource) deleteLROEndsWithNotFoundError(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_resource" "test" {
  type      = "Microsoft.ServiceBus/namespaces@2022-01-01-preview"
  name      = "acctest-sb-%[2]d"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location
  body = jsonencode({
    sku = {
      name = "Premium"
    }
  })
}

`, r.template(data), data.RandomInteger, data.RandomStringOfLength(10))
}

func (r GenericResource) locks(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s


resource "azurerm_route_table" "test" {
  name                = "acctestrt%[2]d"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Network/routeTables/routes@2022-07-01"
  name      = "first%[2]d"
  parent_id = azurerm_route_table.test.id
  body = jsonencode({
    properties = {
      nextHopType   = "VnetLocal"
      addressPrefix = "10.1.0.0/16"
    }
  })

  locks = [azurerm_route_table.test.id, azurerm_resource_group.test.id]
}

resource "azapi_resource" "test2" {
  type      = "Microsoft.Network/routeTables/routes@2022-07-01"
  name      = "second%[2]d"
  parent_id = azurerm_route_table.test.id
  body = jsonencode({
    properties = {
      nextHopType   = "VnetLocal"
      addressPrefix = "10.3.0.0/16"
    }
  })

  locks = [azurerm_route_table.test.id, azurerm_resource_group.test.id]
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
