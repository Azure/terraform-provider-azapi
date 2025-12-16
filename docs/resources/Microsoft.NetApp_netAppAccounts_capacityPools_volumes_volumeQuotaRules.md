---
subcategory: "Microsoft.NetApp - Azure NetApp Files"
page_title: "netAppAccounts/capacityPools/volumes/volumeQuotaRules"
description: |-
  Manages a Volume Quota Rule.
---

# Microsoft.NetApp/netAppAccounts/capacityPools/volumes/volumeQuotaRules - Volume Quota Rule

This article demonstrates how to use `azapi` provider to manage the Volume Quota Rule resource in Azure.



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
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
    SkipNRMSNSG      = "true"
  }
}

resource "azapi_resource" "networkSecurityGroup" {
  type      = "Microsoft.Network/networkSecurityGroups@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nsg"
  location  = var.location
  body = {
    properties = {
      securityRules = []
    }
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.88.0.0/16"]
      }
      dhcpOptions = {
        dnsServers = []
      }
      subnets = []
    }
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"
  body = {
    properties = {
      addressPrefix         = "10.88.2.0/24"
      defaultOutboundAccess = true
      delegations = [{
        name = "netapp-delegation"
        properties = {
          serviceName = "Microsoft.NetApp/volumes"
        }
      }]
      networkSecurityGroup = {
        id = azapi_resource.networkSecurityGroup.id
      }
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "netAppAccount" {
  type      = "Microsoft.NetApp/netAppAccounts@2025-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-acct"
  location  = var.location
  body = {
    properties = {}
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "capacityPool" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools@2025-01-01"
  parent_id = azapi_resource.netAppAccount.id
  name      = "${var.resource_name}-pool"
  location  = var.location
  body = {
    properties = {
      coolAccess     = false
      encryptionType = "Single"
      qosType        = "Auto"
      serviceLevel   = "Standard"
      size           = 4398046511104
    }
  }
  tags = {
    CreatedOnDate    = "2023-08-17T08:01:00Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "volume" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes@2025-01-01"
  parent_id = azapi_resource.capacityPool.id
  name      = "${var.resource_name}-vol"
  location  = var.location
  body = {
    properties = {
      creationToken  = "${var.resource_name}-path"
      dataProtection = {}
      exportPolicy = {
        rules = []
      }
      protocolTypes  = ["NFSv3"]
      serviceLevel   = "Standard"
      subnetId       = azapi_resource.subnet.id
      usageThreshold = 107374182400
    }
  }
  tags = {
    CreatedOnDate    = "2022-07-08T23:50:21Z"
    SkipASMAzSecPack = "true"
  }
}

resource "azapi_resource" "volumeQuotaRule" {
  type      = "Microsoft.NetApp/netAppAccounts/capacityPools/volumes/volumeQuotaRules@2025-01-01"
  parent_id = azapi_resource.volume.id
  name      = "${var.resource_name}-quota"
  location  = var.location
  body = {
    properties = {
      quotaSizeInKiBs = 2048
      quotaType       = "DefaultGroupQuota"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.NetApp/netAppAccounts/capacityPools/volumes/volumeQuotaRules@api-version`. The available api-versions for this resource are: [`2022-01-01`, `2022-03-01`, `2022-05-01`, `2022-09-01`, `2022-11-01`, `2022-11-01-preview`, `2023-05-01`, `2023-05-01-preview`, `2023-07-01`, `2023-07-01-preview`, `2023-11-01`, `2023-11-01-preview`, `2024-01-01`, `2024-03-01`, `2024-03-01-preview`, `2024-05-01`, `2024-05-01-preview`, `2024-07-01`, `2024-07-01-preview`, `2024-09-01`, `2024-09-01-preview`, `2025-01-01`, `2025-01-01-preview`, `2025-03-01`, `2025-03-01-preview`, `2025-06-01`, `2025-07-01-preview`, `2025-08-01`, `2025-08-01-preview`, `2025-09-01`, `2025-09-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}/capacityPools/{resourceName}/volumes/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.NetApp/netAppAccounts/capacityPools/volumes/volumeQuotaRules?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}/capacityPools/{resourceName}/volumes/{resourceName}/volumeQuotaRules/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.NetApp/netAppAccounts/{resourceName}/capacityPools/{resourceName}/volumes/{resourceName}/volumeQuotaRules/{resourceName}?api-version=2025-09-01-preview
 ```
