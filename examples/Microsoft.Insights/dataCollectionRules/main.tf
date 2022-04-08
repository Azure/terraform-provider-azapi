terraform {
  required_providers {
    azapi = {
      source = "azure/azapi"
    }
  }
}

provider "azurerm" {
  features {}
}

provider "azapi" {
}

resource "azurerm_resource_group" "test" {
  name     = "myResourceGroup"
  location = "westus"
}

resource "azurerm_log_analytics_workspace" "test" {
  name                = "myLogWorkspace"
  location            = azurerm_resource_group.test.location
  resource_group_name = azurerm_resource_group.test.name
  sku                 = "PerGB2018"
  retention_in_days   = 30
}

resource "azapi_resource" "test" {
  type      = "Microsoft.Insights/dataCollectionRules@2021-04-01"
  name      = "myDataCollectionRules"
  parent_id = azurerm_resource_group.test.id
  location  = azurerm_resource_group.test.location

  body = jsonencode({
    kind = "Windows"
    properties = {
      dataSources = {
        performanceCounters = [
          {
            name = "cloudTeamCoreCounters"
            streams = [
              "Microsoft-Perf"
            ]
            samplingFrequencyInSeconds = 15
            counterSpecifiers = [
              "\\Processor(_Total)\\% Processor Time",
              "\\Memory\\Committed Bytes",
              "\\LogicalDisk(_Total)\\Free Megabytes",
              "\\PhysicalDisk(_Total)\\Avg. Disk Queue Length",
            ]
          },
          {
            name = "appTeamExtraCounters"
            streams = [
              "Microsoft-Perf"
            ]
            samplingFrequencyInSeconds = 15
            counterSpecifiers = [
              "\\Process(_Total)\\Thread Count"
            ]
          }
        ]
        windowsEventLogs = [
          {
            name = "cloudSecurityTeamEvents"
            streams = [
              "Microsoft-WindowsEvent"
            ]
            xPathQueries = [
              "Security!"
            ]
          },
          {
            name = "appTeam1AppEvents"
            streams = [
              "Microsoft-WindowsEvent"
            ]
            xPathQueries = [
              "System![System[(Level = 1 or Level = 2 or Level = 3)]]",
              "Application!*[System[(Level = 1 or Level = 2 or Level = 3)]]"
            ]
          }
        ]
      }

      dataFlows = [
        {
          streams = [
            "Microsoft-Perf",
            "Microsoft-Syslog",
            "Microsoft-WindowsEvent",
          ]
          destinations = [
            "centralWorkspace"
          ]
        }
      ]

      destinations = {
        logAnalytics = [
          {
            name                = "centralWorkspace"
            workspaceResourceId = azurerm_log_analytics_workspace.test.id
          }
        ]
      }
    }
  })
}