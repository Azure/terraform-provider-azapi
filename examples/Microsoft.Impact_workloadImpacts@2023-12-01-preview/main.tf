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

data "azapi_client_config" "current" {}

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
        adminPassword = "Password1234!"
        adminUsername = "testadmin"
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
}



resource "azapi_resource" "workloadImpact" {
  type      = "Microsoft.Impact/workloadImpacts@2023-12-01-preview"
  parent_id = "/subscriptions/${data.azapi_client_config.current.subscription_id}"
  name      = var.resource_name
  body = {
    properties = {
      additionalProperties = {
        CollectTelemetry = true
        Location         = "DataCenter1"
        LogUrl           = "http://example.com/log"
        Manufacturer     = "ManufacturerName"
        ModelNumber      = "Model123"
        NodeId           = "node-123"
        PhysicalHostName = "host123"
        SerialNumber     = "SN123456"
        VmUniqueId       = "vm-unique-id"
      }
      armCorrelationIds = [
        "id1",
        "id2",
      ]
      clientIncidentDetails = {
        clientIncidentId     = "id"
        clientIncidentSource = "AzureDevops"
      }
      confidenceLevel = "High"
      connectivity = {
        port     = 1443
        protocol = "TCP"
        source = {
          azureResourceId = azapi_resource.virtualMachine.id
        }
        target = {
          azureResourceId = azapi_resource.virtualMachine.id
        }
      }
      endDateTime = "2024-12-04T01:15:00Z"
      errorDetails = {
        errorCode    = "code"
        errorMessage = "errorMessage"
      }
      impactCategory     = "Resource.Availability"
      impactDescription  = "impact description"
      impactGroupId      = "impact groupid"
      impactedResourceId = azapi_resource.virtualMachine.id
      performance = [
        {
          actual   = 2
          expected = 2
          expectedValueRange = {
            max = 5
            min = 1
          }
          metricName = "example"
          unit       = "ByteSeconds"
        },
      ]
      startDateTime = "2024-12-03T01:15:00Z"
      workload = {
        context = "context"
        toolset = "Ansible"
      }
    }
  }

  schema_validation_enabled = false
}
