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
        adminPassword            = "Password1234!"
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