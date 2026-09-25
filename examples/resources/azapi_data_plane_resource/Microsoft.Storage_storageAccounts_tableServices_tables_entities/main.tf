terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

locals {
  storage_table_headers = {
    "x-ms-version" = "2026-04-06"
  }
  storage_table_entity_headers = {
    "x-ms-version"        = "2026-04-06"
    Accept                = "application/json;odata=nometadata"
    DataServiceVersion    = "3.0;NetFx"
    MaxDataServiceVersion = "3.0;NetFx"
  }
}

data "azapi_client_config" "current" {}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier = "Hot"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_data_plane_resource" "table" {
  type           = "Microsoft.Storage/storageAccounts/tableServices/tables@2026-04-06"
  parent_id      = "${azapi_resource.storageAccount.name}.table.core.windows.net"
  name           = var.resource_name
  create_headers = local.storage_table_headers
  read_headers   = local.storage_table_headers
  delete_headers = local.storage_table_headers
  body = {
    TableName = var.resource_name
  }

  retry = {
    error_message_regex  = ["AuthorizationPermissionMismatch", "AuthorizationFailure", "Forbidden", "Unauthorized", "authorization"]
    interval_seconds     = 20
    max_interval_seconds = 120
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    storageTableDataContributorRoleId = "value[?properties.roleName == 'Storage Table Data Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.storageAccount.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.storageTableDataContributorRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

# Table entity writes use insert-or-replace semantics. Each create or update request
# must include the complete desired entity content. Properties omitted from a write,
# including properties added outside Terraform, are removed by that write.
#
# With sensitive_body_version, unchanged sensitive properties are omitted from update
# requests and can therefore be removed by this replacement operation.
resource "azapi_data_plane_resource" "entity" {
  type           = "Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06"
  parent_id      = "${azapi_resource.storageAccount.name}.table.core.windows.net/${azapi_data_plane_resource.table.name}(PartitionKey='example',RowKey='state')"
  create_headers = local.storage_table_entity_headers
  read_headers   = local.storage_table_entity_headers
  update_headers = local.storage_table_entity_headers
  delete_headers = merge(local.storage_table_entity_headers, {
    "If-Match" = "*"
  })
  body = {
    outputs = jsonencode({
      status = "ok"
    })
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}

data "azapi_data_plane_resource" "entity" {
  type      = "Microsoft.Storage/storageAccounts/tableServices/tables/entities@2026-04-06"
  parent_id = azapi_data_plane_resource.entity.parent_id
  headers   = local.storage_table_entity_headers

  depends_on = [
    azapi_data_plane_resource.entity,
  ]
}

data "azapi_data_plane_resource" "entities" {
  type      = "Microsoft.Storage/storageAccounts/tableServices/tables/entitiesCollection@2026-04-06"
  parent_id = "${azapi_resource.storageAccount.name}.table.core.windows.net/${azapi_data_plane_resource.table.name}"
  headers   = local.storage_table_entity_headers
  query_parameters = {
    "$filter" = ["PartitionKey eq 'example'"]
  }

  depends_on = [
    azapi_data_plane_resource.entity,
  ]
}
