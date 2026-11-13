---
subcategory: "Microsoft.App - Azure Container Apps"
page_title: "agents"
description: |-
  Manages a SRE Agent.
---

# Microsoft.App/agents - SRE Agent

This article demonstrates how to use `azapi` provider to manage the SRE Agent resource in Azure.



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

data "azapi_client_config" "current" {}

data "azapi_resource_list" "roleDefinitions" {
  type                   = "Microsoft.Authorization/roleDefinitions@2022-04-01"
  parent_id              = data.azapi_client_config.current.subscription_resource_id
  response_export_values = ["value"]
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
  response_export_values    = ["properties.principalId"]
}

resource "azapi_resource" "roleAssignment" {
  type                      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name                      = uuidv5("url", "${azapi_resource.resourceGroup.id}/${azapi_resource.userAssignedIdentity.id}/Reader")
  parent_id                 = azapi_resource.resourceGroup.id
  schema_validation_enabled = false

  body = {
    properties = {
      principalId   = azapi_resource.userAssignedIdentity.output.properties.principalId
      principalType = "ServicePrincipal"
      roleDefinitionId = one([
        for role in data.azapi_resource_list.roleDefinitions.output.value :
        role.id if role.properties.roleName == "Reader" && role.properties.type == "BuiltInRole"
      ])
    }
  }
}

resource "azapi_resource" "workspace" {
  type                      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  name                      = "${var.resource_name}-log"
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
      Application_Type                = "web"
      WorkspaceResourceId             = azapi_resource.workspace.id
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
}

resource "azapi_resource" "virtualNetwork" {
  type                      = "Microsoft.Network/virtualNetworks@2022-07-01"
  name                      = "${var.resource_name}-vnet"
  parent_id                 = azapi_resource.resourceGroup.id
  location                  = var.location
  schema_validation_enabled = false

  body = {
    properties = {
      addressSpace = {
        addressPrefixes = ["10.30.0.0/24"]
      }
    }
  }

  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type                      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  name                      = "${var.resource_name}-subnet"
  parent_id                 = azapi_resource.virtualNetwork.id
  schema_validation_enabled = false

  body = {
    properties = {
      addressPrefix = "10.30.0.0/27"
      delegations = [
        {
          name = "sre-agent"
          properties = {
            serviceName = "Microsoft.App/environments"
          }
        }
      ]
    }
  }
}

resource "azapi_resource" "agent" {
  type                      = "Microsoft.App/agents@2026-07-01"
  name                      = var.resource_name
  parent_id                 = azapi_resource.resourceGroup.id
  location                  = var.location
  schema_validation_enabled = false

  tags = {
    purpose = "example"
  }

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
      incidentManagementConfiguration = {
        type           = "AzMonitor"
        connectionName = "azmonitor"
      }
      knowledgeGraphConfiguration = {
        identity         = azapi_resource.userAssignedIdentity.id
        managedResources = [azapi_resource.resourceGroup.id]
      }
      logConfiguration = {
        applicationInsightsConfiguration = {
          appId            = azapi_resource.component.output.properties.AppId
          connectionString = azapi_resource.component.output.properties.ConnectionString
        }
      }
      sandboxConfiguration = {
        egress = {
          mode                            = "AzureVNet"
          allowHttpMcpServerNetworkAccess = false
          allowedCodeRepositories         = ["Github"]
          allowedHosts                    = ["learn.microsoft.com"]
          allowedRegistries               = ["pypi"]
        }
        packages = [
          {
            name           = "numpy"
            packageManager = "pip"
            version        = "1.26.0"
          }
        ]
      }
      upgradeChannel = "Stable"
      vnetConfiguration = {
        subnetResourceId = azapi_resource.subnet.id
      }
    }
  }

  depends_on = [azapi_resource.roleAssignment]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.App/agents@api-version`. The available api-versions for this resource are: [`2026-01-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.App/agents?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/agents/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.App/agents/{resourceName}?api-version=2026-01-01
 ```
