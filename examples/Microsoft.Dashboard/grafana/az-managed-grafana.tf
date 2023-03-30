#####################
## Resource Group:-
#####################

resource "azurerm_resource_group" "azrg" {
  name        = var.rg-name
  location    = var.rg-location
}

##############################
## Azure Managed Grafana:-
##############################
resource "azapi_resource" "azgrafana" {
  type        = "Microsoft.Dashboard/grafana@2022-08-01" 
  name        = var.az-grafana-name
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

