---
subcategory: "Microsoft.App - Azure Container Apps"
page_title: "managedEnvironments/privateEndpointConnections"
description: |-
  Manages a Container App Environment Private Endpoint Connection.
---

# Microsoft.App/managedEnvironments/privateEndpointConnections - Container App Environment Private Endpoint Connection

This article demonstrates how to use `azapi` provider to manage the Container App Environment Private Endpoint Connection resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

variable "resource_name" {
  type    = string
  default = "acctest5925"
}

variable "location" {
  type    = string
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_action" "sharedKeys" {
  type                   = "Microsoft.OperationalInsights/workspaces@2020-08-01"
  resource_id            = azapi_resource.workspace.id
  action                 = "sharedKeys"
  response_export_values = ["*"]
}

# Create Virtual Network for private endpoint
resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.0.0.0/16"]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

# Create Subnet for private endpoint (minimum /27 required for workload profiles)
resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2023-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"
  body = {
    properties = {
      addressPrefix                  = "10.0.0.0/21"
      privateEndpointNetworkPolicies = "Disabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

# Create Container Apps Environment with workload profiles (required for private endpoints)
resource "azapi_resource" "managedEnvironment" {
  type      = "Microsoft.App/managedEnvironments@2024-10-02-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = azapi_resource.workspace.output.properties.customerId
          sharedKey  = data.azapi_resource_action.sharedKeys.output.primarySharedKey
        }
      }
      # Enable workload profiles (default for new environments, but explicitly set)
      workloadProfiles = [
        {
          name                = "Consumption"
          workloadProfileType = "Consumption"
        }
      ]
      # Disable public network access to enable private endpoints
      publicNetworkAccess = "Disabled"
      vnetConfiguration = {
        # Note: For private endpoints, we don't inject into VNet but create separate private endpoint
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

# Create Private Endpoint for the Container Apps Environment
resource "azapi_resource" "privateEndpoint" {
  type      = "Microsoft.Network/privateEndpoints@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-pe"
  location  = var.location
  body = {
    properties = {
      subnet = {
        id = azapi_resource.subnet.id
      }
      privateLinkServiceConnections = [
        {
          name = "${var.resource_name}-connection"
          properties = {
            privateLinkServiceId = azapi_resource.managedEnvironment.id
            groupIds             = ["managedEnvironments"]
          }
        }
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

# Create Private DNS Zone for Container Apps
resource "azapi_resource" "privateDnsZone" {
  type      = "Microsoft.Network/privateDnsZones@2020-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctestzone.azurecontainerapps.dev"
  location  = "global"
  body = {
    properties = {}
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

# Link VNet to Private DNS Zone
resource "azapi_resource" "vnetLink" {
  type      = "Microsoft.Network/privateDnsZones/virtualNetworkLinks@2020-06-01"
  parent_id = azapi_resource.privateDnsZone.id
  name      = "${var.resource_name}-vnet-link"
  location  = "global"
  body = {
    properties = {
      virtualNetwork = {
        id = azapi_resource.virtualNetwork.id
      }
      registrationEnabled = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

# Create DNS Zone Group for automatic DNS record management
resource "azapi_resource" "privateDnsZoneGroup" {
  type      = "Microsoft.Network/privateEndpoints/privateDnsZoneGroups@2023-05-01"
  parent_id = azapi_resource.privateEndpoint.id
  name      = "default"
  body = {
    properties = {
      privateDnsZoneConfigs = [
        {
          name = "config"
          properties = {
            privateDnsZoneId = azapi_resource.privateDnsZone.id
          }
        }
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_list" "privateEndpointConnections" {
  type                   = "Microsoft.App/managedEnvironments/privateEndpointConnections@2024-10-02-preview"
  parent_id              = azapi_resource.managedEnvironment.id
  response_export_values = ["*"]
}

output "privateEndpointConnections" {
  value = data.azapi_resource_list.privateEndpointConnections.output
}

locals {
  privateEndpointConnectionIds = [
    for pe in data.azapi_resource_list.privateEndpointConnections.output.value : pe.id
  ]
}

# Note: The private endpoint connection is automatically created when the private endpoint is created
# This resource manages the approval state of the connection
resource "azapi_update_resource" "privateEndpointConnection" {
  for_each    = toset(local.privateEndpointConnectionIds)
  type        = "Microsoft.App/managedEnvironments/privateEndpointConnections@2024-10-02-preview"
  resource_id = each.value
  body = {
    properties = {
      privateLinkServiceConnectionState = {
        status          = "Approved"
        description     = "Auto-approved"
        actionsRequired = "None"
      }
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.App/managedEnvironments/privateEndpointConnections@api-version`. The available api-versions for this resource are: [`2024-02-02-preview`, `2024-08-02-preview`, `2024-10-02-preview`, `2025-02-02-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/managedEnvironments/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.App/managedEnvironments/privateEndpointConnections?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/managedEnvironments/{resourceName}/privateEndpointConnections/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/managedEnvironments/{resourceName}/privateEndpointConnections/{resourceName}?api-version=2025-02-02-preview
 ```
