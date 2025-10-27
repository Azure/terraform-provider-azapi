---
subcategory: "Microsoft.RecoveryServices - Azure Site Recovery"
page_title: "vaults/backupFabrics/protectionContainers/protectedItems"
description: |-
  Manages a Azure Backup Protected Item.
---

# Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems - Azure Backup Protected Item

This article demonstrates how to use `azapi` provider to manage the Azure Backup Protected Item resource in Azure.



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

locals {
  sa_base   = substr(lower(join("", split("-", var.resource_name))), 0, 24)
  sa_name   = local.sa_base
  dns_label = substr(lower(var.resource_name), 0, 63)
  comp_name = substr(var.resource_name, 0, 15)
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-pip"
  location  = var.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      dnsSettings = {
        domainNameLabel = local.dns_label
      }
      idleTimeoutInMinutes     = 4
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Dynamic"
    }
    sku = {
      name = "Basic"
      tier = "Regional"
    }
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
        addressPrefixes = ["10.0.0.0/16"]
      }
      dhcpOptions = {
        dnsServers = []
      }
      privateEndpointVNetPolicies = "Disabled"
      subnets                     = []
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2024-05-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"
  body = {
    properties = {
      addressPrefix                     = "10.0.10.0/24"
      defaultOutboundAccess             = true
      delegations                       = []
      privateEndpointNetworkPolicies    = "Disabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies           = []
      serviceEndpoints                  = []
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nic"
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [{
        name = "acctestipconfig"
        properties = {
          primary                   = true
          privateIPAddressVersion   = "IPv4"
          privateIPAllocationMethod = "Dynamic"
          publicIPAddress = {
            id = azapi_resource.publicIPAddress.id
          }
          subnet = {
            id = azapi_resource.subnet.id
          }
        }
      }]
    }
  }
}

resource "azapi_resource" "disk" {
  type      = "Microsoft.Compute/disks@2023-04-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-datadisk"
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption = "Empty"
      }
      diskSizeGB = 1023
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy        = "AllowAll"
      optimizedForFrequentAttach = false
      osType                     = null
      publicNetworkAccess        = "Enabled"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.sa_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = {
            keyType = "Service"
          }
          table = {
            keyType = "Service"
          }
        }
      }
      isHnsEnabled       = false
      isLocalUserEnabled = true
      isNfsV3Enabled     = false
      isSftpEnabled      = false
      minimumTlsVersion  = "TLS1_2"
      networkAcls = {
        bypass              = "AzureServices"
        defaultAction       = "Allow"
        ipRules             = []
        resourceAccessRules = []
        virtualNetworkRules = []
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.RecoveryServices/vaults@2024-01-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-rsv"
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

resource "azapi_resource" "backupPolicy" {
  type      = "Microsoft.RecoveryServices/vaults/backupPolicies@2024-10-01"
  parent_id = azapi_resource.vault.id
  name      = "${var.resource_name}-policy"
  body = {
    properties = {
      backupManagementType = "AzureIaasVM"
      policyType           = "V1"
      retentionPolicy = {
        dailySchedule = {
          retentionDuration = {
            count        = 10
            durationType = "Days"
          }
          retentionTimes = ["2025-07-03T23:00:00Z"]
        }
        retentionPolicyType = "LongTermRetentionPolicy"
      }
      schedulePolicy = {
        schedulePolicyType   = "SimpleSchedulePolicy"
        scheduleRunDays      = []
        scheduleRunFrequency = "Daily"
        scheduleRunTimes     = ["2025-07-03T23:00:00Z"]
      }
      tieringPolicy = {
        ArchivedRP = {
          duration     = 0
          durationType = "Invalid"
          tieringMode  = "DoNotTier"
        }
      }
      timeZone = "UTC"
    }
  }
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2024-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vm"
  location  = var.location
  body = {
    properties = {
      diagnosticsProfile = {
        bootDiagnostics = {
          enabled    = true
          storageUri = "https://${local.sa_name}.blob.core.windows.net/"
        }
      }
      hardwareProfile = {
        vmSize = "Standard_B1s"
      }
      networkProfile = {
        networkInterfaces = [{
          id = azapi_resource.networkInterface.id
          properties = {
            primary = true
          }
        }]
      }
      osProfile = {
        adminPassword = "Password123!@#"
        adminUsername = "vmadmin"
        computerName  = local.comp_name
        linuxConfiguration = {
          disablePasswordAuthentication = false
        }
      }
      storageProfile = {
        dataDisks = [{
          createOption = "Attach"
          diskSizeGB   = 1023
          lun          = 0
          managedDisk = {
            id                 = azapi_resource.disk.id
            storageAccountType = "Standard_LRS"
          }
          name                    = "${var.resource_name}-datadisk"
          writeAcceleratorEnabled = false
          }, {
          createOption = "Empty"
          diskSizeGB   = 4
          lun          = 1
          managedDisk = {
            storageAccountType = "Standard_LRS"
          }
          name                    = "${var.resource_name}-datadisk2"
          writeAcceleratorEnabled = false
        }]
        imageReference = {
          offer     = "0001-com-ubuntu-server-jammy"
          publisher = "Canonical"
          sku       = "22_04-lts"
          version   = "latest"
        }
        osDisk = {
          caching      = "ReadWrite"
          createOption = "FromImage"
          managedDisk = {
            storageAccountType = "Standard_LRS"
          }
          name                    = "${var.resource_name}-osdisk"
          writeAcceleratorEnabled = false
        }
      }
    }
  }
}

resource "azapi_resource" "protectedItem" {
  type      = "Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems@2023-02-01"
  parent_id = "${azapi_resource.vault.id}/backupFabrics/Azure/protectionContainers/iaasvmcontainer;iaasvmcontainerv2;${azapi_resource.resourceGroup.name};${azapi_resource.virtualMachine.name}"
  name      = "VM;iaasvmcontainerv2;${azapi_resource.resourceGroup.name};${azapi_resource.virtualMachine.name}"
  body = {
    properties = {
      extendedProperties = {
        diskExclusionProperties = {
          diskLunList     = [0]
          isInclusionList = true
        }
      }
      policyId          = azapi_resource.backupPolicy.id
      protectedItemType = "Microsoft.Compute/virtualMachines"
      sourceResourceId  = azapi_resource.virtualMachine.id
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems@api-version`. The available api-versions for this resource are: [`2016-06-01`, `2019-05-13`, `2019-06-15`, `2020-10-01`, `2020-12-01`, `2021-01-01`, `2021-02-01`, `2021-02-01-preview`, `2021-02-10`, `2021-03-01`, `2021-04-01`, `2021-06-01`, `2021-07-01`, `2021-08-01`, `2021-10-01`, `2021-12-01`, `2022-01-01`, `2022-02-01`, `2022-03-01`, `2022-04-01`, `2022-06-01-preview`, `2022-09-01-preview`, `2022-09-30-preview`, `2022-10-01`, `2023-01-01`, `2023-02-01`, `2023-04-01`, `2023-06-01`, `2023-08-01`, `2024-01-01`, `2024-02-01`, `2024-04-01`, `2024-04-30-preview`, `2024-07-30-preview`, `2024-10-01`, `2024-11-01-preview`, `2025-01-01`, `2025-02-01`, `2025-02-28-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `{any azure resource id}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.RecoveryServices/vaults/backupFabrics/protectionContainers/protectedItems?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example {any azure resource id}/protectedItems/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example {any azure resource id}/protectedItems/{resourceName}?api-version=2025-02-28-preview
 ```
