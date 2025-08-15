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

resource "azapi_resource" "hostPool" {
  type      = "Microsoft.DesktopVirtualization/hostPools@2024-04-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-hp"
  location  = var.location
  body = {
    properties = {
      customRdpProperty             = ""
      description                   = ""
      friendlyName                  = ""
      hostPoolType                  = "Pooled"
      loadBalancerType              = "BreadthFirst"
      maxSessionLimit               = 999999
      personalDesktopAssignmentType = ""
      preferredAppGroupType         = "Desktop"
      publicNetworkAccess           = "Enabled"
      startVMOnConnect              = false
      validationEnvironment         = false
      vmTemplate                    = ""
    }
  }
}

resource "azapi_resource" "applicationGroup" {
  type      = "Microsoft.DesktopVirtualization/applicationGroups@2024-04-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-ag"
  location  = var.location
  body = {
    properties = {
      applicationGroupType = "Desktop"
      description          = ""
      friendlyName         = ""
      hostPoolArmPath      = azapi_resource.hostPool.id
    }
  }
}
