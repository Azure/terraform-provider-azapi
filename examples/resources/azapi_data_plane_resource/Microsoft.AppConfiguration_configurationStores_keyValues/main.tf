terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "example-resources"
  location = "westeurope"
}

resource "azapi_resource" "appconf" {
  type      = "Microsoft.AppConfiguration/configurationStores@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "exampleappconf"
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "standard"
    }
  }
  response_export_values = {
    endpoint = "properties.endpoint"
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    appConfigDataOwnerRoleId = "value[?properties.roleName == 'App Configuration Data Owner'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.appconf.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.appConfigDataOwnerRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.AppConfiguration/configurationStores/keyValues@1.0"
  parent_id = replace(azapi_resource.appconf.output.endpoint, "https://", "")
  name      = "mykey"
  body = {
    content_type = ""
    value        = "myvalue"
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
