terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "resource_group" {
  type     = "Microsoft.Resources/resourceGroups@2024-03-01"
  name     = "example-resource-group"
  location = "westus3"
}

resource "azapi_resource" "storage_account" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resource_group.id
  name      = "examplestorageaccount"
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
  name      = "example"
  body = {
    properties = {
      publicAccess = "None"
    }
  }
}

resource "azapi_resource" "blob_data_contributor" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name      = uuidv5("url", "${azapi_resource.storage_account.id}/${data.azapi_client_config.current.object_id}/blob-data-contributor")
  parent_id = azapi_resource.storage_account.id
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/ba92f5b4-2d11-453d-a403-e96b0029c9fe"
    }
  }
}

resource "azapi_data_plane_resource" "blob" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers/blobs@2023-11-03"
  parent_id = "${azapi_resource.storage_account.name}.blob.core.windows.net/${azapi_resource.container.name}"
  name      = "documents/hello.txt"

  body = {
    content_base64 = base64encode("Hello, world!")
    content_type   = "text/plain; charset=utf-8"
    metadata = {
      environment = "example"
    }
  }

  retry = {
    error_message_regex = ["AuthorizationPermissionMismatch", "Forbidden"]
  }

  depends_on = [azapi_resource.blob_data_contributor]
}
