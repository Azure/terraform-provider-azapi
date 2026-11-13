---
subcategory: "Microsoft.App - Azure Container Apps"
page_title: "agents/connectors"
description: |-
  Manages a SRE Agent Connector.
---

# Microsoft.App/agents/connectors - SRE Agent Connector

This article demonstrates how to use `azapi` provider to manage the SRE Agent Connector resource in Azure.



## Example Usage

### default

```hcl
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
  default = "eastus2"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "userAssignedIdentity" {
  type                      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  name                      = "${var.resource_name}-id"
  parent_id                 = azapi_resource.resourceGroup.id
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["properties.clientId", "properties.tenantId"]
}

resource "azapi_resource" "workspace" {
  type                      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  name                      = "${var.resource_name}-logs"
  parent_id                 = azapi_resource.resourceGroup.id
  location                  = var.location
  schema_validation_enabled = false

  body = {
    properties = {
      sku = {
        name = "PerGB2018"
      }
      retentionInDays = 30
    }
  }
}

resource "azapi_resource" "component" {
  type                      = "Microsoft.Insights/components@2020-02-02"
  name                      = "${var.resource_name}-appi"
  parent_id                 = azapi_resource.resourceGroup.id
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["properties.AppId", "properties.ConnectionString"]

  body = {
    kind = "web"
    properties = {
      Application_Type    = "web"
      WorkspaceResourceId = azapi_resource.workspace.id
    }
  }
}

resource "azapi_resource" "agent" {
  type                      = "Microsoft.App/agents@2026-07-01"
  name                      = "${var.resource_name}-agent"
  parent_id                 = azapi_resource.resourceGroup.id
  location                  = var.location
  schema_validation_enabled = false

  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }

  body = {
    properties = {
      actionConfiguration = {
        accessLevel = "Low"
        identity    = azapi_resource.userAssignedIdentity.id
        mode        = "Autonomous"
      }
      defaultModel = {
        name     = "Automatic"
        provider = "MicrosoftFoundry"
      }
      logConfiguration = {
        applicationInsightsConfiguration = {
          appId            = azapi_resource.component.output.properties.AppId
          connectionString = azapi_resource.component.output.properties.ConnectionString
        }
      }
      upgradeChannel = "Stable"
    }
  }
}

resource "azapi_resource" "cluster" {
  type                      = "Microsoft.Kusto/clusters@2023-05-02"
  name                      = "adx${var.resource_name}"
  parent_id                 = azapi_resource.resourceGroup.id
  location                  = var.location
  schema_validation_enabled = false
  response_export_values    = ["properties.uri"]

  body = {
    sku = {
      name     = "Dev(No SLA)_Standard_D11_v2"
      tier     = "Basic"
      capacity = 1
    }
    properties = {
      enableAutoStop      = true
      engineType          = "V2"
      publicNetworkAccess = "Enabled"
    }
  }
}

resource "azapi_resource" "database" {
  type                      = "Microsoft.Kusto/clusters/databases@2023-05-02"
  name                      = "${var.resource_name}-telemetry"
  parent_id                 = azapi_resource.cluster.id
  location                  = var.location
  schema_validation_enabled = false

  body = {
    kind       = "ReadWrite"
    properties = {}
  }
}

resource "azapi_resource" "principalAssignment" {
  type                      = "Microsoft.Kusto/clusters/databases/principalAssignments@2023-05-02"
  name                      = "${var.resource_name}-viewer"
  parent_id                 = azapi_resource.database.id
  schema_validation_enabled = false

  body = {
    properties = {
      principalId   = azapi_resource.userAssignedIdentity.output.properties.clientId
      principalType = "App"
      role          = "Viewer"
      tenantId      = azapi_resource.userAssignedIdentity.output.properties.tenantId
    }
  }
}

resource "azapi_resource" "connector" {
  type                      = "Microsoft.App/agents/connectors@2026-07-01"
  name                      = "${var.resource_name}-kusto"
  parent_id                 = azapi_resource.agent.id
  schema_validation_enabled = false

  body = {
    properties = {
      dataConnectorType = "Kusto"
      dataSource        = "${azapi_resource.cluster.output.properties.uri}/${azapi_resource.database.name}"
      endpoint          = azapi_resource.cluster.output.properties.uri
      identity          = azapi_resource.userAssignedIdentity.id
      extendedProperties = {
        armResourceId = azapi_resource.database.id
        resource = {
          name = azapi_resource.database.name
        }
      }
    }
  }

  depends_on = [azapi_resource.principalAssignment]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.App/agents/connectors@api-version`. The available api-versions for this resource are: [`2026-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/agents/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.App/agents/connectors?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/agents/{resourceName}/connectors/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/agents/{resourceName}/connectors/{resourceName}?api-version=2026-01-01
 ```
