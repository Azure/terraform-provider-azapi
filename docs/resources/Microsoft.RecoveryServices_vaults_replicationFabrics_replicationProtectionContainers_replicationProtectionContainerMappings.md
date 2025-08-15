---
subcategory: "Microsoft.RecoveryServices - Azure Site Recovery"
page_title: "vaults/replicationFabrics/replicationProtectionContainers/replicationProtectionContainerMappings"
description: |-
  Manages a Site Recovery protection container mapping on Azure.
---

# Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers/replicationProtectionContainerMappings - Site Recovery protection container mapping on Azure

This article demonstrates how to use `azapi` provider to manage the Site Recovery protection container mapping on Azure resource in Azure.

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
  name     = "${var.resource_name}-rg"
  location = var.location
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.RecoveryServices/vaults@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "acctest-vault-250703130022502990"
  location  = var.location
  body = {
    properties = {
      publicNetworkAccess = "Enabled"
      redundancySettings = {
        crossRegionRestore            = "Disabled"
        standardTierStorageRedundancy = "GeoRedundant"
      }
    }
    sku = {
      name = "Standard"
    }
  }
}

resource "azapi_resource" "replicationFabric" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics@2024-04-01"
  parent_id = azapi_resource.vault.id
  name      = "acctest-fabric1-250703130022502990"
  body = {
    properties = {
      customDetails = {
        instanceType = "Azure"
        location     = "westeurope"
      }
    }
  }
}

resource "azapi_resource" "replicationFabric_1" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics@2024-04-01"
  parent_id = azapi_resource.vault.id
  name      = "acctest-fabric2b-250703130022502990"
  body = {
    properties = {
      customDetails = {
        instanceType = "Azure"
        location     = "westus2"
      }
    }
  }
}

resource "azapi_resource" "replicationPolicy" {
  type      = "Microsoft.RecoveryServices/vaults/replicationPolicies@2024-04-01"
  parent_id = azapi_resource.vault.id
  name      = "acctest-policy-250703130022502990"
  body = {
    properties = {
      providerSpecificInput = {
        appConsistentFrequencyInMinutes = 240
        instanceType                    = "A2A"
        multiVmSyncStatus               = "Enable"
        recoveryPointHistory            = 1440
      }
    }
  }
}

resource "azapi_resource" "replicationProtectionContainer" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers@2024-04-01"
  parent_id = azapi_resource.replicationFabric.id
  name      = "acctest-protection-cont1-250703130022502990"
  body = {
    properties = {}
  }
}

resource "azapi_resource" "replicationProtectionContainer_1" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers@2024-04-01"
  parent_id = azapi_resource.replicationFabric_1.id
  name      = "acctest-protection-cont2-250703130022502990"
  body = {
    properties = {}
  }
}

resource "azapi_resource" "replicationProtectionContainerMapping" {
  type      = "Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers/replicationProtectionContainerMappings@2024-04-01"
  parent_id = azapi_resource.replicationProtectionContainer.id
  name      = "mapping-250703130022502990"
  body = {
    properties = {
      policyId = azapi_resource.replicationPolicy.id
      providerSpecificInput = {
        instanceType = "A2A"
      }
      targetProtectionContainerId = azapi_resource.replicationProtectionContainer_1.id
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers/replicationProtectionContainerMappings@api-version`. The available api-versions for this resource are: [`2016-08-10`, `2018-01-10`, `2018-07-10`, `2021-02-10`, `2021-03-01`, `2021-04-01`, `2021-06-01`, `2021-07-01`, `2021-08-01`, `2021-10-01`, `2021-11-01`, `2021-12-01`, `2022-01-01`, `2022-02-01`, `2022-03-01`, `2022-04-01`, `2022-05-01`, `2022-08-01`, `2022-09-10`, `2022-10-01`, `2023-01-01`, `2023-02-01`, `2023-04-01`, `2023-06-01`, `2023-08-01`, `2024-01-01`, `2024-02-01`, `2024-04-01`, `2024-10-01`, `2025-01-01`, `2025-02-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{resourceName}/replicationProtectionContainers/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.RecoveryServices/vaults/replicationFabrics/replicationProtectionContainers/replicationProtectionContainerMappings?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{resourceName}/replicationProtectionContainers/{resourceName}/replicationProtectionContainerMappings/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{resourceName}/replicationFabrics/{resourceName}/replicationProtectionContainers/{resourceName}/replicationProtectionContainerMappings/{resourceName}?api-version=2025-02-01
 ```
