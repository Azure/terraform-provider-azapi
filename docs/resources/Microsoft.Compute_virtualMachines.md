---
subcategory: "Microsoft.Compute - Virtual Machines, Virtual Machine Scale Sets"
page_title: "virtualMachines"
description: |-
  Manages a Virtual Machine.
---

# Microsoft.Compute/virtualMachines - Virtual Machine

This article demonstrates how to use `azapi` provider to manage the Virtual Machine resource in Azure.



## Example Usage

### attach_data_disk

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
  default = "westeurope"
}

variable "admin_username" {
  type        = string
  description = "The administrator username for the virtual machine"
}

variable "admin_password" {
  type        = string
  description = "The administrator password for the virtual machine"
  sensitive   = true
}

locals {
  os_disk_name            = "myosdisk1"
  data_disk_name          = "mydatadisk1"
  attached_data_disk_name = "myattacheddatadisk1"
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
          "10.0.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
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

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          name = "testconfiguration1"
          properties = {
            primary                   = true
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachedDisk" {
  type      = "Microsoft.Compute/disks@2022-03-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.attached_data_disk_name
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption = "Empty"
      }
      diskSizeGB = 1
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy = "AllowAll"
      osType              = "Linux"
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_F2"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.networkInterface.id
            properties = {
              primary = false
            }
          },
        ]
      }
      osProfile = {
        adminPassword = var.admin_password
        adminUsername = var.admin_username
        computerName  = "hostname230630032848831819"
        linuxConfiguration = {
          disablePasswordAuthentication = false
        }
      }
      storageProfile = {
        imageReference = {
          offer     = "UbuntuServer"
          publisher = "Canonical"
          sku       = "16.04-LTS"
          version   = "latest"
        }
        osDisk = {
          caching                 = "ReadWrite"
          createOption            = "FromImage"
          name                    = local.os_disk_name
          writeAcceleratorEnabled = false
        }
        dataDisks = [
          {
            caching      = "ReadWrite"
            createOption = "Empty"
            name         = local.data_disk_name
            diskSizeGB   = 1
            lun          = 1
            managedDisk = {
              storageAccountType = "Standard_LRS"
            }
          },
          {
            caching      = "ReadWrite"
            createOption = "Attach"
            name         = azapi_resource.attachedDisk.name
            lun          = 2
            managedDisk = {
              id = azapi_resource.attachedDisk.id
            }
          }
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```

### attach_os_disk

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
  default = "westeurope"
}

variable "admin_username" {
  type        = string
  description = "The administrator username for the virtual machine"
}

variable "admin_password" {
  type        = string
  description = "The administrator password for the virtual machine"
  sensitive   = true
}

variable "attached_resource_name" {
  type    = string
  default = "acctest0002"
}

locals {
  os_disk_name          = "myosdisk1"
  attached_os_disk_name = "myosdisk2"
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
          "10.0.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
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

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          name = "testconfiguration1"
          properties = {
            primary                   = true
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_F2"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.networkInterface.id
            properties = {
              primary = false
            }
          },
        ]
      }
      osProfile = {
        adminPassword = var.admin_password
        adminUsername = var.admin_username
        computerName  = "hostname230630032848831819"
        linuxConfiguration = {
          disablePasswordAuthentication = false
        }
      }
      storageProfile = {
        imageReference = {
          offer     = "UbuntuServer"
          publisher = "Canonical"
          sku       = "16.04-LTS"
          version   = "latest"
        }
        osDisk = {
          caching                 = "ReadWrite"
          createOption            = "FromImage"
          name                    = local.os_disk_name
          writeAcceleratorEnabled = false
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource" "managedDisk" {
  type      = "Microsoft.Compute/disks@2023-10-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.os_disk_name

  depends_on = [azapi_resource.virtualMachine]
}

resource "azapi_resource" "snapshot" {
  type      = "Microsoft.Compute/snapshots@2023-10-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    sku = {
      name = "Standard_ZRS"
    }
    properties = {
      creationData = {
        createOption     = "Copy"
        sourceResourceId = data.azapi_resource.managedDisk.id
      }
      diskSizeGB = 30
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy = "AllowAll"
      osType              = "Linux"
      hyperVGeneration    = "V1"
      incremental         = true
      publicNetworkAccess = "Enabled"
      supportedCapabilities = {
        architecture = "x64"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachedManagedDisk" {
  type      = "Microsoft.Compute/disks@2023-10-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.attached_os_disk_name
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption     = "Copy",
        sourceResourceId = azapi_resource.snapshot.id
      }

      diskSizeGB = 30
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy = "AllowAll"
      osType              = "Linux"
      hyperVGeneration    = "V1"
      publicNetworkAccess = "Enabled"
      supportedCapabilities = {
        architecture = "x64"
      }
    }
    sku = {
      name = "Standard_LRS"
    }
    zones = [
      "1"
    ]
  }

  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachedNetworkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.attached_resource_name
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          name = "testconfiguration2"
          properties = {
            primary                   = true
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        }
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "attachedVirtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.attached_resource_name
  location  = var.location
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_F2"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.attachedNetworkInterface.id
            properties = {
              primary = false
            }
          },
        ]
      }
      storageProfile = {
        osDisk = {
          caching                 = "ReadWrite"
          createOption            = "Attach"
          name                    = local.attached_os_disk_name
          osType                  = "Linux",
          writeAcceleratorEnabled = false
          managedDisk = {
            id = azapi_resource.attachedManagedDisk.id
          }
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```

### basic

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
  default = "westeurope"
}

variable "admin_username" {
  type        = string
  description = "The administrator username for the virtual machine"
}

variable "admin_password" {
  type        = string
  description = "The administrator password for the virtual machine"
  sensitive   = true
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
          "10.0.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
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

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          name = "testconfiguration1"
          properties = {
            primary                   = true
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_F2"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.networkInterface.id
            properties = {
              primary = false
            }
          },
        ]
      }
      osProfile = {
        adminPassword = var.admin_password
        adminUsername = var.admin_username
        computerName  = "hostname230630032848831819"
        linuxConfiguration = {
          disablePasswordAuthentication = false
        }
      }
      storageProfile = {
        imageReference = {
          offer     = "UbuntuServer"
          publisher = "Canonical"
          sku       = "16.04-LTS"
          version   = "latest"
        }
        osDisk = {
          caching                 = "ReadWrite"
          createOption            = "FromImage"
          name                    = "myosdisk1"
          writeAcceleratorEnabled = false
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```

### tag_os_disk

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
  default = "westeurope"
}

variable "admin_username" {
  type        = string
  description = "The administrator username for the virtual machine"
}

variable "admin_password" {
  type        = string
  description = "The administrator password for the virtual machine"
  sensitive   = true
}

locals {
  os_disk_name = "myosdisk1"

  tags = {
    environment = "accTest0001"
  }
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
          "10.0.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
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

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          name = "testconfiguration1"
          properties = {
            primary                   = true
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_F2"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.networkInterface.id
            properties = {
              primary = false
            }
          },
        ]
      }
      osProfile = {
        adminPassword = var.admin_password
        adminUsername = var.admin_username
        computerName  = "hostname230630032848831819"
        linuxConfiguration = {
          disablePasswordAuthentication = false
        }
      }
      storageProfile = {
        imageReference = {
          offer     = "UbuntuServer"
          publisher = "Canonical"
          sku       = "16.04-LTS"
          version   = "latest"
        }
        osDisk = {
          caching                 = "ReadWrite"
          createOption            = "FromImage"
          name                    = local.os_disk_name
          writeAcceleratorEnabled = false
        }
      }
    }
    tags = local.tags
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource" "managedDisk" {
  type      = "Microsoft.Compute/disks@2022-03-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.os_disk_name

  depends_on = [azapi_resource.virtualMachine]
}

resource "azapi_resource_action" "updateTags" {
  type        = "Microsoft.Compute/disks@2022-03-02"
  resource_id = data.azapi_resource.managedDisk.id
  method      = "PATCH"

  body = {
    tags = local.tags
  }

  depends_on = [azapi_resource.virtualMachine]
}
```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Compute/virtualMachines@api-version`. The available api-versions for this resource are: [`2015-06-15`, `2016-03-30`, `2016-04-30-preview`, `2017-03-30`, `2017-12-01`, `2018-04-01`, `2018-06-01`, `2018-10-01`, `2019-03-01`, `2019-07-01`, `2019-12-01`, `2020-06-01`, `2020-12-01`, `2021-03-01`, `2021-04-01`, `2021-07-01`, `2021-11-01`, `2022-03-01`, `2022-08-01`, `2022-11-01`, `2023-03-01`, `2023-07-01`, `2023-09-01`, `2024-03-01`, `2024-07-01`, `2024-11-01`, `2025-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Compute/virtualMachines?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/virtualMachines/{resourceName}?api-version=2025-04-01
 ```
