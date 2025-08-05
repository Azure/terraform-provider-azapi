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

resource "azapi_resource" "gallery" {
  type      = "Microsoft.Compute/galleries@2022-03-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}sig"
  location  = var.location
  body = {
    properties = {
      description = ""
    }
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}acc"
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

resource "azapi_resource" "application" {
  type      = "Microsoft.Compute/galleries/applications@2022-03-03"
  parent_id = azapi_resource.gallery.id
  name      = "${var.resource_name}-app"
  location  = var.location
  body = {
    properties = {
      supportedOSType = "Linux"
    }
  }
}

resource "azapi_resource" "container" {
  type      = "Microsoft.Storage/storageAccounts/blobServices/containers@2023-05-01"
  parent_id = "${azapi_resource.storageAccount.id}/blobServices/default"
  name      = "mycontainer"
  body = {
    properties = {
      publicAccess = "Blob"
    }
  }
}

resource "azapi_resource" "version" {
  type      = "Microsoft.Compute/galleries/applications/versions@2022-03-03"
  parent_id = azapi_resource.application.id
  name      = "0.0.1"
  location  = var.location
  body = {
    properties = {
      publishingProfile = {
        enableHealthCheck = false
        excludeFromLatest = false
        manageActions = {
          install = "[install command]"
          remove  = "[remove command]"
          update  = ""
        }
        source = {
          defaultConfigurationLink = ""
          mediaLink                = "https://${azapi_resource.storageAccount.name}.blob.core.windows.net/mycontainer/myblob"
        }
        targetRegions = [{
          name                 = var.location
          regionalReplicaCount = 1
          storageAccountType   = "Standard_LRS"
        }]
      }
      safetyProfile = {
        allowDeletionOfReplicatedLocations = true
      }
    }
  }
  depends_on = [azapi_resource.container]
}

