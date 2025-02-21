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

variable "resource_site_name" {
  type    = string
  default = "acctestsite0001"
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

resource "azapi_resource" "serverfarm" {
  type      = "Microsoft.Web/serverfarms@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hyperV         = false
      perSiteScaling = false
      reserved       = false
      zoneRedundant  = false
    }
    sku = {
      name = "S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "site" {
  type      = "Microsoft.Web/sites@2022-09-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_site_name
  location  = var.location
  body = {
    properties = {
      clientAffinityEnabled = false
      clientCertEnabled     = false
      clientCertMode        = "Required"
      enabled               = true
      httpsOnly             = false
      publicNetworkAccess   = "Enabled"
      serverFarmId          = azapi_resource.serverfarm.id
      siteConfig = {
        acrUseManagedIdentityCreds       = false
        alwaysOn                         = true
        autoHealEnabled                  = false
        ftpsState                        = "Disabled"
        http20Enabled                    = false
        loadBalancing                    = "LeastRequests"
        localMySqlEnabled                = false
        managedPipelineMode              = "Integrated"
        minTlsVersion                    = "1.2"
        publicNetworkAccess              = "Enabled"
        remoteDebuggingEnabled           = false
        scmIpSecurityRestrictionsUseMain = false
        scmMinTlsVersion                 = "1.2"
        use32BitWorkerProcess            = true
        vnetRouteAllEnabled              = false
        webSocketsEnabled                = false
        windowsFxVersion                 = ""
      }
      vnetRouteAllEnabled = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "dynatraceSiteExtension" {
  type                      = "Microsoft.Web/sites/siteextensions@2022-09-01"
  parent_id                 = azapi_resource.site.id
  name                      = "Dynatrace"
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["*"]
}
