package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccGenericUpdateResource_readOverrideAppSettings(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_update_resource", "test")
	r := GenericUpdateResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.readOverrideAppSettings(data),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
				check.That(data.ResourceName).Key("resource_id").Exists(),
				check.That(data.ResourceName).Key("parent_id").Exists(),
				check.That(data.ResourceName).Key("name").Exists(),
			),
		},
	})
}

func readOverrideAppSettingsTemplate(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s


resource "azapi_resource" "serverfarm" {
  type      = "Microsoft.Web/serverfarms@2023-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "functionapp"
    sku = {
      name = "FC1"
      tier = "FlexConsumption"
    }
    properties = {
      reserved = true
    }
  }
}

resource "azapi_resource" "storage" {
  type      = "Microsoft.Storage/storageAccounts@2023-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest%[2]s"
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "StorageV2"
    sku = {
      name = "Standard_LRS"
    }
    properties = {
      minimumTlsVersion = "TLS1_2"
    }
  }
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-01-01"
  parent_id = "${azapi_resource.storage.id}/blobServices/default"
  name      = "deployments"
  body = {
    properties = {}
  }
}

resource "azapi_resource" "site" {
  type      = "Microsoft.Web/sites@2023-12-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-%[2]s"
  location  = azapi_resource.resourceGroup.location
  identity {
    type = "SystemAssigned"
  }
  body = {
    kind = "functionapp,linux"
    properties = {
      reserved     = true
      serverFarmId = azapi_resource.serverfarm.id
      functionAppConfig = {
        deployment = {
          storage = {
            type  = "blobContainer"
            value = "https://acctest%[2]s.blob.core.windows.net/deployments"
            authentication = {
              type = "SystemAssignedIdentity"
            }
          }
        }
        runtime = {
          name    = "dotnet-isolated"
          version = "8.0"
        }
        scaleAndConcurrency = {
          maximumInstanceCount = 40
          instanceMemoryMB     = 2048
        }
      }
    }
  }
  depends_on = [azapi_resource.container]
}
`, GenericUpdateResource{}.template(data), data.RandomString)
}

func (r GenericUpdateResource) readOverrideAppSettings(data acceptance.TestData) string {
	return fmt.Sprintf(`
%[1]s

resource "azapi_update_resource" "test" {
  type      = "Microsoft.Web/sites/config@2023-12-01"
  name      = "appsettings"
  parent_id = azapi_resource.site.id
  body = {
    properties = {
      KEY1 = "value1"
      KEY2 = "value2"
      KEY3 = "value3"
    }
  }
  ignore_casing           = false
  ignore_missing_property = false
  read_override = {
    method = "POST"
    action = "list"
  }
}
`, readOverrideAppSettingsTemplate(data))
}
