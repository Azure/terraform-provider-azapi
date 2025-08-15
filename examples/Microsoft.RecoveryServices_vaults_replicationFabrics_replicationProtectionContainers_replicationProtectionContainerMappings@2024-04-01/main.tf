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

