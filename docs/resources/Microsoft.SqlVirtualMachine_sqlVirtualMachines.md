---
subcategory: "Microsoft.SqlVirtualMachine - SQL Server on Azure Virtual Machines"
page_title: "sqlVirtualMachines"
description: |-
  Manages a Microsoft SQL Virtual Machine.
---

# Microsoft.SqlVirtualMachine/sqlVirtualMachines - Microsoft SQL Virtual Machine

This article demonstrates how to use `azapi` provider to manage the Microsoft SQL Virtual Machine resource in Azure.

## Example Usage

### SQL_best_practices_assessment

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
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "networkSecurityGroup" {
  type      = "Microsoft.Network/networkSecurityGroups@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      securityRules = [{
        name = "MSSQLRule"
        properties = {
          access                     = "Allow"
          destinationAddressPrefix   = "*"
          destinationAddressPrefixes = []
          destinationPortRange       = "1433"
          destinationPortRanges      = []
          direction                  = "Inbound"
          priority                   = 1001
          protocol                   = "Tcp"
          sourceAddressPrefix        = "167.220.255.0/25"
          sourceAddressPrefixes      = []
          sourcePortRange            = "*"
          sourcePortRanges           = []
        }
      }]
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.0.0/24"
      networkSecurityGroup = {
        id = azapi_resource.networkSecurityGroup.id
      }
    }
  }
}

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 4
      ipTags                   = []
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Dynamic"
    }
    sku = {
      name = "Basic"
      tier = "Regional"
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      disableTcpStateTracking = false
      dnsSettings = {
        dnsServers = []
      }
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          type = "Microsoft.Network/networkInterfaces/ipConfigurations"
          name = "testconfiguration1"
          properties = {
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            publicIPAddress = {
              id = azapi_resource.publicIPAddress.id
            }
            subnet = {
              id = azapi_resource.subnet.id
            }
            primary          = true
            privateIPAddress = "10.0.0.4"
          }
        }
      ]
      nicType       = "Standard"
      auxiliaryMode = "None"
      auxiliarySku  = "None"
    }
  }
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2024-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      osProfile = {
        adminUsername            = "testadmin"
        adminPassword            = var.admin_password
        allowExtensionOperations = true
        computerName             = "winhost01"
        secrets                  = []
        windowsConfiguration = {
          timeZone               = "Pacific Standard Time"
          enableAutomaticUpdates = true
          patchSettings = {
            patchMode      = "AutomaticByOS"
            assessmentMode = "ImageDefault"
          }
          provisionVMAgent = true
        }
      }
      storageProfile = {
        dataDisks = []
        imageReference = {
          offer     = "SQL2017-WS2016"
          publisher = "MicrosoftSQLServer"
          sku       = "SQLDEV"
          version   = "latest"
        }
        osDisk = {
          diskSizeGB = 127
          managedDisk = {
            storageAccountType = "Premium_LRS"
          }
          name                    = "acctvm-250116171212663925OSDisk"
          osType                  = "Windows"
          writeAcceleratorEnabled = false
          caching                 = "ReadOnly"
          createOption            = "FromImage"
          deleteOption            = "Detach"
        }
      }
      hardwareProfile = {
        vmSize = "Standard_F2s"
      }
      networkProfile = {
        networkInterfaces = [
          {
            properties = {
              primary = false
            }
            id = azapi_resource.networkInterface.id
          }
        ]
      }
    }
  }
}


resource "azapi_resource" "extension" {
  type      = "Microsoft.Compute/virtualMachines/extensions@2024-07-01"
  parent_id = azapi_resource.virtualMachine.id
  name      = "AzureMonitorWindowsAgent"
  location  = "westeurope"
  body = {
    properties = {
      autoUpgradeMinorVersion = true
      enableAutomaticUpgrade  = true
      publisher               = "Microsoft.Azure.Monitor"
      suppressFailures        = false
      type                    = "AzureMonitorWindowsAgent"
      typeHandlerVersion      = "1.0"
    }
  }
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2020-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
        legacy                                      = 0
        searchVersion                               = 1
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
}

resource "azapi_resource" "table" {
  type      = "Microsoft.OperationalInsights/workspaces/tables@2023-09-01"
  parent_id = azapi_resource.workspace.id
  name      = "SqlAssessment_CL"
  body = {
    properties = {
      schema = {
        name = "SqlAssessment_CL"
        columns = [
          {
            name = "TimeGenerated"
            type = "datetime"
          },
          {
            type = "string"
            name = "RawData"
          }
        ]
      }
    }
  }
}

resource "azapi_resource" "dataCollectionEndpoint" {
  type      = "Microsoft.Insights/dataCollectionEndpoints@2022-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.location}-DCE-1"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      networkAcls = {
        publicNetworkAccess = "Enabled"
      }
    }
  }
}

resource "azapi_resource" "dataCollectionRule" {
  type      = "Microsoft.Insights/dataCollectionRules@2022-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${azapi_resource.workspace.output.properties.customerId}_${azapi_resource.resourceGroup.location}_DCR_1"
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      streamDeclarations = {
        Custom-SqlAssessment_CL = {
          columns = [
            {
              name = "TimeGenerated"
              type = "datetime"
            },
            {
              type = "string"
              name = "RawData"
            }
          ]
        }
      }
      dataCollectionEndpointId = azapi_resource.dataCollectionEndpoint.id
      dataFlows = [
        {
          outputStream = "Custom-SqlAssessment_CL"
          streams = [
            "Custom-SqlAssessment_CL"
          ]
          transformKql = "source"
          destinations = [
            azapi_resource.workspace.name
          ]
        }
      ]
      dataSources = {
        logFiles = [
          {
            filePatterns = [
              "C:\\Windows\\System32\\config\\systemprofile\\AppData\\Local\\Microsoft SQL Server IaaS Agent\\Assessment\\*.csv"
            ]
            format = "text"
            name   = "Custom-SqlAssessment_CL"
            settings = {
              text = {
                recordStartTimestampFormat = "ISO 8601"
              }
            }
            streams = [
              "Custom-SqlAssessment_CL"
            ]
          }
        ]
      }
      description = ""
      destinations = {
        logAnalytics = [
          {
            name                = azapi_resource.workspace.name
            workspaceResourceId = azapi_resource.workspace.id
          }
        ]
      }
    }
  }
  depends_on = [azapi_resource.table]
}

resource "azapi_resource" "dataCollectionRuleAssociation" {
  type      = "Microsoft.Insights/dataCollectionRuleAssociations@2022-06-01"
  parent_id = azapi_resource.virtualMachine.id
  name      = "${azapi_resource.workspace.output.properties.customerId}_${azapi_resource.resourceGroup.location}_DCRA_1"
  body = {
    properties = {
      dataCollectionRuleId = azapi_resource.dataCollectionRule.id
    }
  }
}

resource "azapi_resource" "sqlvirtualMachine" {
  type      = "Microsoft.SqlVirtualMachine/sqlVirtualMachines@2023-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.virtualMachine.name
  location  = azapi_resource.virtualMachine.location
  body = {
    properties = {
      sqlServerLicenseType     = "PAYG"
      virtualMachineResourceId = azapi_resource.virtualMachine.id
      enableAutomaticUpgrade   = true
      leastPrivilegeMode       = "Enabled"
      sqlImageOffer            = "SQL2017-WS2016"
      sqlImageSku              = "Developer"
      sqlManagement            = "Full"
      assessmentSettings = {
        enable         = true
        runImmediately = false
        schedule = {
          dayOfWeek      = "Monday"
          enable         = true
          startTime      = "00:00"
          weeklyInterval = 1
        }
      }
    }
  }
  depends_on = [
    azapi_resource.dataCollectionRuleAssociation,
    azapi_resource.extension,
  ]
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

variable "vm_admin_password" {
  type        = string
  description = "The administrator password for the SQL virtual machine"
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
    }
  }
  schema_validation_enabled = false
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "networkSecurityGroup" {
  type      = "Microsoft.Network/networkSecurityGroups@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      securityRules = [{
        name = "MSSQLRule"
        properties = {
          access                     = "Allow"
          destinationAddressPrefix   = "*"
          destinationAddressPrefixes = []
          destinationPortRange       = "1433"
          destinationPortRanges      = []
          direction                  = "Inbound"
          priority                   = 1001
          protocol                   = "Tcp"
          sourceAddressPrefix        = "167.220.255.0/25"
          sourceAddressPrefixes      = []
          sourcePortRange            = "*"
          sourcePortRanges           = []
        }
      }]
    }
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = var.resource_name
  body = {
    properties = {
      addressPrefix = "10.0.0.0/24"
      networkSecurityGroup = {
        id = azapi_resource.networkSecurityGroup.id
      }
    }
  }
}

resource "azapi_resource" "publicIPAddress" {
  type      = "Microsoft.Network/publicIPAddresses@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      ddosSettings = {
        protectionMode = "VirtualNetworkInherited"
      }
      idleTimeoutInMinutes     = 4
      ipTags                   = []
      publicIPAddressVersion   = "IPv4"
      publicIPAllocationMethod = "Dynamic"
    }
    sku = {
      name = "Basic"
      tier = "Regional"
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      disableTcpStateTracking = false
      dnsSettings = {
        dnsServers = []
      }
      enableAcceleratedNetworking = false
      enableIPForwarding          = false
      ipConfigurations = [
        {
          type = "Microsoft.Network/networkInterfaces/ipConfigurations"
          name = "testconfiguration1"
          properties = {
            privateIPAddressVersion   = "IPv4"
            privateIPAllocationMethod = "Dynamic"
            publicIPAddress = {
              id = azapi_resource.publicIPAddress.id
            }
            subnet = {
              id = azapi_resource.subnet.id
            }
            primary          = true
            privateIPAddress = "10.0.0.4"
          }
        }
      ]
      nicType       = "Standard"
      auxiliaryMode = "None"
      auxiliarySku  = "None"
    }
  }
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2024-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      osProfile = {
        adminUsername            = "testadmin"
        adminPassword            = var.vm_admin_password
        allowExtensionOperations = true
        computerName             = "winhost01"
        secrets                  = []
        windowsConfiguration = {
          timeZone               = "Pacific Standard Time"
          enableAutomaticUpdates = true
          patchSettings = {
            patchMode      = "AutomaticByOS"
            assessmentMode = "ImageDefault"
          }
          provisionVMAgent = true
        }
      }
      storageProfile = {
        dataDisks = []
        imageReference = {
          offer     = "SQL2017-WS2016"
          publisher = "MicrosoftSQLServer"
          sku       = "SQLDEV"
          version   = "latest"
        }
        osDisk = {
          diskSizeGB = 127
          managedDisk = {
            storageAccountType = "Premium_LRS"
          }
          name                    = "acctvm-250116171212663925OSDisk"
          osType                  = "Windows"
          writeAcceleratorEnabled = false
          caching                 = "ReadOnly"
          createOption            = "FromImage"
          deleteOption            = "Detach"
        }
      }
      hardwareProfile = {
        vmSize = "Standard_F2s"
      }
      networkProfile = {
        networkInterfaces = [
          {
            properties = {
              primary = false
            }
            id = azapi_resource.networkInterface.id
          }
        ]
      }
    }
  }
}


resource "azapi_resource" "sqlvirtualMachine" {
  type      = "Microsoft.SqlVirtualMachine/sqlVirtualMachines@2023-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = azapi_resource.virtualMachine.name
  location  = azapi_resource.virtualMachine.location
  body = {
    properties = {
      sqlServerLicenseType     = "PAYG"
      virtualMachineResourceId = azapi_resource.virtualMachine.id
      enableAutomaticUpgrade   = true
      leastPrivilegeMode       = "Enabled"
      sqlImageOffer            = "SQL2017-WS2016"
      sqlImageSku              = "Developer"
      sqlManagement            = "Full"
    }
  }
}
```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.SqlVirtualMachine/sqlVirtualMachines@api-version`. The available api-versions for this resource are: [`2017-03-01-preview`, `2021-11-01-preview`, `2022-02-01`, `2022-02-01-preview`, `2022-07-01-preview`, `2022-08-01-preview`, `2023-01-01-preview`, `2023-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.SqlVirtualMachine/sqlVirtualMachines?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachines/{resourceName}?api-version=2023-10-01
 ```
