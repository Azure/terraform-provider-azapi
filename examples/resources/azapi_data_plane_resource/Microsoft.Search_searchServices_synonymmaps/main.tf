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

resource "azapi_resource" "searchService" {
  type      = "Microsoft.Search/searchServices@2023-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "examplesearch"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      replicaCount   = 1
      partitionCount = 1
      hostingMode    = "default"
      authOptions = {
        aadOrApiKey = {
          aadAuthFailureMode = "http401WithBearerChallenge"
        }
      }
    }
    sku = {
      name = "basic"
    }
  }
}

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type      = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  response_export_values = {
    searchIndexDataContributorRoleId = "value[?properties.roleName == 'Search Index Data Contributor'].id | [0]"
  }
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  parent_id = azapi_resource.searchService.id
  name      = uuid()
  body = {
    properties = {
      principalId      = data.azapi_client_config.current.object_id
      roleDefinitionId = data.azapi_resource_list.roleDefinitions.output.searchIndexDataContributorRoleId
    }
  }
  lifecycle {
    ignore_changes = [name]
  }
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.Search/searchServices/synonymmaps@2024-07-01"
  parent_id = "${azapi_resource.searchService.name}.search.windows.net"
  name      = "mysynonymmap"
  body = {
    format   = "solr"
    synonyms = "hotel, motel\nairport, aerodrome"
  }

  depends_on = [
    azapi_resource.roleAssignment,
  ]
}
