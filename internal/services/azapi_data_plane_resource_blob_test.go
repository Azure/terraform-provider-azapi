package services_test

import (
	"fmt"
	"testing"

	"github.com/Azure/terraform-provider-azapi/internal/acceptance"
	"github.com/Azure/terraform-provider-azapi/internal/acceptance/check"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccDataPlaneResource_storageBlob(t *testing.T) {
	data := acceptance.BuildTestData(t, "azapi_data_plane_resource", "test")
	r := DataPlaneResource{}

	data.ResourceTest(t, r, []resource.TestStep{
		{
			Config: r.storageBlob(data, "Hello from AzAPI!", "text/plain", "create"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("name").HasValue("documents/test.txt"),
				check.That(data.ResourceName).Key("output.content_type").HasValue("text/plain"),
				check.That(data.ResourceName).Key("output.metadata_environment").HasValue("create"),
				check.That(data.ResourceName).Key("output.content_length").HasValue("17"),
				check.That(data.ResourceName).Key("output.blob_type").HasValue("BlockBlob"),
			),
		},
		{
			Config: r.storageBlob(data, "Updated from AzAPI!", "text/markdown", "update"),
			Check: resource.ComposeTestCheckFunc(
				check.That(data.ResourceName).Key("output.content_type").HasValue("text/markdown"),
				check.That(data.ResourceName).Key("output.metadata_environment").HasValue("update"),
				check.That(data.ResourceName).Key("output.content_length").HasValue("19"),
				check.That(data.ResourceName).Key("output.blob_type").HasValue("BlockBlob"),
			),
		},
		{
			ResourceName:      data.ResourceName,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateIdFunc: storageBlobImportID,
			ImportStateVerifyIgnore: []string{
				"body",
				"ignore_casing",
				"ignore_missing_property",
				"locks",
				"output",
				"response_export_values",
				"retry",
			},
		},
	})
}

func storageBlobImportID(tfState *terraform.State) (string, error) {
	state := tfState.RootModule().Resources["azapi_data_plane_resource.test"].Primary
	return fmt.Sprintf("%s|%s", state.ID, state.Attributes["type"]), nil
}

func (r DataPlaneResource) storageBlob(data acceptance.TestData, content, contentType, environment string) string {
	return fmt.Sprintf(`
data "azapi_client_config" "current" {}

resource "azapi_resource" "resource_group" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "acctest%[1]s"
  location = "%[2]s"
}

resource "azapi_resource" "storage_account" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resource_group.id
  name      = "acctest%[1]s"
  location  = azapi_resource.resource_group.location
  body = {
    kind = "StorageV2"
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = "${azapi_resource.storage_account.id}/blobServices/default"
  name      = "content"
  body = {
    properties = {
      publicAccess = "None"
    }
  }
}

resource "azapi_resource" "blob_data_contributor" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.storage_account.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/ba92f5b4-2d11-453d-a403-e96b0029c9fe"
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "test" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers/blobs@2023-11-03"
  parent_id = "${azapi_resource.storage_account.name}.blob.core.windows.net/${azapi_resource.container.name}"
  name      = "documents/test.txt"

  body = {
    content_base64 = base64encode(%[3]q)
    content_type   = %[4]q
    metadata = {
      environment = %[5]q
    }
  }

  response_export_values = {
    content_type         = "content_type"
    metadata_environment = "metadata.environment"
    content_length       = "content_length"
    blob_type            = "blob_type"
  }

  retry = {
    error_message_regex  = ["AuthorizationPermissionMismatch", "Forbidden"]
    interval_seconds     = 10
    max_interval_seconds = 60
  }

  depends_on = [
    azapi_resource.blob_data_contributor,
  ]
}
`, data.RandomString, data.LocationPrimary, content, contentType, environment)
}
