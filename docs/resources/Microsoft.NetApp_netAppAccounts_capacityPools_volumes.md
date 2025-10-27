---
subcategory: "Microsoft.NetApp - Azure NetApp Files"
page_title: "netAppAccounts/capacityPools/volumes"
description: |-
  Manages a NetApp Volume.
---

# Microsoft.NetApp/netAppAccounts/capacityPools/volumes - NetApp Volume

This article demonstrates how to use `azapi` provider to manage the NetApp Volume resource in Azure.



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
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "centralus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.6.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
    tags = {
      SkipASMAzSecPack = "true"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "netAppAccount" {
  type      = "Microsoft.NetApp/netAppAccounts@2022-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      activeDirectories = [
      ]
    }
    tags = {
      SkipASMAzSecPack = "true"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "GatewaySubnet"
  body = {
    properties = {
      addressPrefix = "10.6.1.0/24"
      delegations = [
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "subnet2" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.6.2.0/24"
      delegations = [
        {
          name = "testdelegation"
          properties = {
            serviceName = "Microsoft.Netapp/volumes"
          }
        },
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "capacityPool" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools@2022-05-01"
  parent_id = azapi_resource.netAppAccount.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      serviceLevel = "Standard"
      size         = 4.398046511104e+12
    }
    tags = {
      SkipASMAzSecPack = "true"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "volume" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes@2022-05-01"
  parent_id = azapi_resource.capacityPool.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      avsDataStore  = "Enabled"
      creationToken = "my-unique-file-path-230630034120103726"
      dataProtection = {
      }
      exportPolicy = {
        rules = [
          {
            allowedClients = "0.0.0.0/0"
            cifs           = false
            hasRootAccess  = true
            nfsv3          = true
            nfsv41         = false
            ruleIndex      = 1
            unixReadOnly   = false
            unixReadWrite  = true
          },
        ]
      }
      networkFeatures = "Basic"
      protocolTypes = [
        "NFSv3",
      ]
      serviceLevel             = "Standard"
      snapshotDirectoryVisible = true
      subnetId                 = azapi_resource.subnet2.id
      usageThreshold           = 1.073741824e+11
      volumeType               = ""
    }
    tags = {
      SkipASMAzSecPack = "true"
    }
    zones = [
    ]
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.NetApp/netAppAccounts/capacityPools/volumes@api-version`. The available api-versions for this resource are: [`2017-08-15`, `2019-05-01`, `2019-06-01`, `2019-07-01`, `2019-08-01`, `2019-10-01`, `2019-11-01`, `2020-02-01`, `2020-03-01`, `2020-05-01`, `2020-06-01`, `2020-07-01`, `2020-08-01`, `2020-09-01`, `2020-11-01`, `2020-12-01`, `2021-02-01`, `2021-04-01`, `2021-04-01-preview`, `2021-06-01`, `2021-08-01`, `2021-10-01`, `2022-01-01`, `2022-03-01`, `2022-05-01`, `2022-09-01`, `2022-11-01`, `2022-11-01-preview`, `2023-05-01`, `2023-05-01-preview`, `2023-07-01`, `2023-07-01-preview`, `2023-11-01`, `2023-11-01-preview`, `2024-01-01`, `2024-03-01`, `2024-03-01-preview`, `2024-05-01`, `2024-05-01-preview`, `2024-07-01`, `2024-07-01-preview`, `2024-09-01`, `2024-09-01-preview`, `2025-01-01`, `2025-01-01-preview`, `2025-03-01`, `2025-03-01-preview`, `2025-06-01`, `2025-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}/capacityPools/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.NetApp/netAppAccounts/capacityPools/volumes?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}/capacityPools/{resourceName}/volumes/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}/capacityPools/{resourceName}/volumes/{resourceName}?api-version=2025-07-01-preview
 ```
