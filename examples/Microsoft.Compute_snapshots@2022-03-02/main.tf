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
}

resource "azapi_resource" "disk" {
  type      = "Microsoft.Compute/disks@2023-04-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}disk"
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption    = "Empty"
        performancePlus = false
      }
      diskSizeGB = 10
      encryption = {
        type = "EncryptionAtRestWithPlatformKey"
      }
      networkAccessPolicy        = "AllowAll"
      optimizedForFrequentAttach = false
      publicNetworkAccess        = "Enabled"
    }
    sku = {
      name = "Standard_LRS"
    }
  }
}

resource "azapi_resource" "snapshot" {
  type      = "Microsoft.Compute/snapshots@2022-03-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}snapshot"
  location  = var.location
  body = {
    properties = {
      creationData = {
        createOption = "Copy"
        sourceUri    = azapi_resource.disk.id
      }
      diskSizeGB          = 20
      incremental         = false
      networkAccessPolicy = "AllowAll"
      publicNetworkAccess = "Enabled"
    }
  }
}
