terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    tls = {
      source  = "hashicorp/tls"
      version = "~> 4.1"
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
  default = "centraluseuap"
}

resource "tls_private_key" "ssh" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "neuroShield" {
  type                      = "Microsoft.NetworkSecurity/neuroShields@2026-10-01-preview"
  parent_id                 = azapi_resource.resourceGroup.id
  name                      = var.resource_name
  location                  = var.location
  schema_validation_enabled = false

  body = {
    properties = {
      responseMode      = "Manual"
      severityThreshold = "Medium"
      description       = "Adaptive virtual patching example"
    }
  }
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
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

  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"

  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
    }
  }
}

resource "azapi_resource" "networkInterface" {
  type      = "Microsoft.Network/networkInterfaces@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-nic"
  location  = var.location

  body = {
    properties = {
      ipConfigurations = [
        {
          name = "internal"
          properties = {
            privateIPAllocationMethod = "Dynamic"
            subnet = {
              id = azapi_resource.subnet.id
            }
          }
        }
      ]
    }
  }
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.Compute/virtualMachines@2023-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location

  identity {
    type = "SystemAssigned"
  }

  body = {
    properties = {
      hardwareProfile = {
        vmSize = "Standard_D2s_v6"
      }
      networkProfile = {
        networkInterfaces = [
          {
            id = azapi_resource.networkInterface.id
          }
        ]
      }
      osProfile = {
        adminUsername = "azureuser"
        computerName  = var.resource_name
        linuxConfiguration = {
          disablePasswordAuthentication = true
          ssh = {
            publicKeys = [
              {
                path    = "/home/azureuser/.ssh/authorized_keys"
                keyData = tls_private_key.ssh.public_key_openssh
              }
            ]
          }
        }
      }
      storageProfile = {
        imageReference = {
          publisher = "Canonical"
          offer     = "ubuntu-26_04-lts"
          sku       = "server"
          version   = "latest"
        }
        osDisk = {
          caching      = "ReadWrite"
          createOption = "FromImage"
          managedDisk = {
            storageAccountType = "Standard_LRS"
          }
        }
      }
    }
  }
}

resource "azapi_resource" "protectedResource" {
  type                      = "Microsoft.NetworkSecurity/neuroShields/protectedResources@2026-10-01-preview"
  parent_id                 = azapi_resource.neuroShield.id
  name                      = var.resource_name
  schema_validation_enabled = false

  body = {
    properties = {
      targetResourceId = azapi_resource.virtualMachine.id
      targetIdentity = {
        type        = "SystemAssigned"
        principalId = azapi_resource.virtualMachine.identity[0].principal_id
        tenantId    = azapi_resource.virtualMachine.identity[0].tenant_id
      }
    }
  }

  # The target resource and identity are create-only; changes require replacing the registration.
  # https://github.com/Azure/azure-rest-api-specs/blob/e0f840aaf5c7c2daeec7da8abcd310dd7f47f56b/specification/networksecurity/resource-manager/Microsoft.NetworkSecurity/neuroShield/models.tsp#L278-L311
  replace_triggers_refs = [
    "properties.targetResourceId",
    "properties.targetIdentity",
  ]
}
