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

resource "azapi_resource" "lab" {
  type      = "Microsoft.LabServices/labs@2022-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoShutdownProfile = {
        shutdownOnDisconnect     = "Disabled"
        shutdownOnIdle           = "None"
        shutdownWhenNotConnected = "Disabled"
      }
      connectionProfile = {
        clientRdpAccess = "None"
        clientSshAccess = "None"
        webRdpAccess    = "None"
        webSshAccess    = "None"
      }
      securityProfile = {
        openAccess = "Disabled"
      }
      title = "Test Title"
      virtualMachineProfile = {
        additionalCapabilities = {
          installGpuDrivers = "Disabled"
        }
        adminUser = {
          password = "Password1234!"
          username = "testadmin"
        }
        createOption = "Image"
        imageReference = {
          offer     = "0001-com-ubuntu-server-focal"
          publisher = "canonical"
          sku       = "20_04-lts"
          version   = "latest"
        }
        sku = {
          capacity = 1
          name     = "Classic_Fsv2_2_4GB_128_S_SSD"
        }
        usageQuota        = "PT0S"
        useSharedPassword = "Disabled"
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "schedule" {
  type      = "Microsoft.LabServices/labs/schedules@2022-08-01"
  parent_id = azapi_resource.lab.id
  name      = var.resource_name
  body = {
    properties = {
      stopAt     = "2023-06-30T04:33:55Z"
      timeZoneId = "America/Los_Angeles"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

