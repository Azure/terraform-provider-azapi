terraform {

  required_version = ">= 1.3.3"
  
  backend "azurerm" {
    resource_group_name  = "tfpipeline-rg"
    storage_account_name = "tfpipelinesa"
    container_name       = "terraform"
    key                  = "AMG/Grafana.tfstate"
  }

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.27"
    }
    azapi = {
      source = "Azure/azapi"
      version = "1.0.0"
    }
    
  }
}
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

provider "azapi" {
}

#####################
## Resource Group:-
#####################

resource "azurerm_resource_group" "azrg" {
  name        = "GrafanaRG"
  location    = "West Europe"
}

##############################
## Azure Managed Grafana:-
##############################
resource "azapi_resource" "azgrafana" {
  type        = "Microsoft.Dashboard/grafana@2022-08-01" 
  name        = "AMGrafanaTest"
  parent_id   = azurerm_resource_group.azrg.id
  location    = azurerm_resource_group.azrg.location
  
  identity {
    type      = "SystemAssigned"
  }

  body = jsonencode({
    sku = {
      name = "Standard"
    }
    properties = {
      publicNetworkAccess = "Enabled",
      zoneRedundancy = "Enabled",
      apiKey = "Enabled",
      deterministicOutboundIP = "Enabled"
    }
  })

}

