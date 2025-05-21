---
subcategory: "Microsoft.OperationalInsights - Azure Monitor"
page_title: "workspaces/tables"
description: |-
  Manages a Operational Insights Workspaces Tables.
---

# Microsoft.OperationalInsights/workspaces/tables - Operational Insights Workspaces Tables

This article demonstrates how to use `azapi` provider to manage the Operational Insights Workspaces Tables resource in Azure.

## Example Usage

### audit_log

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
  default = "westeurope"
}

locals {
  audit_log_table_name = "AuditLog_CL"
  audit_log_columns = [
    {
      "name" : "appId",
      "type" : "string"
    },
    {
      "name" : "correlationId",
      "type" : "string"
    },
    {
      "name" : "TimeGenerated",
      "type" : "datetime"
    }
  ]
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "table" {
  type      = "Microsoft.OperationalInsights/workspaces/tables@2022-10-01"
  parent_id = azapi_resource.workspace.id
  name      = local.audit_log_table_name
  body = {
    properties = {
      schema = {
        name    = local.audit_log_table_name
        columns = local.audit_log_columns
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```

### basic

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
  default = "westeurope"
}

locals {
  sentinel_ti_alerts_table_name = "SentinelTIAlerts_CL"
  sentinel_ti_alerts_columns = [
    {
      "name" : "ConfidenceScore",
      "type" : "int"
    },
    {
      "name" : "ExternalIndicatorId",
      "type" : "string"
    },
    {
      "name" : "IndicatorType",
      "type" : "string"
    },
    {
      "name" : "Indicator",
      "type" : "string"
    },
    {
      "name" : "TimeGenerated",
      "type" : "datetime"
    },
    {
      "name" : "MatchType",
      "type" : "string"
    },
    {
      "name" : "OriginTimestamp",
      "type" : "datetime"
    },
    {
      "name" : "Details",
      "type" : "dynamic"
    }
  ]
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "table" {
  type      = "Microsoft.OperationalInsights/workspaces/tables@2022-10-01"
  parent_id = azapi_resource.workspace.id
  name      = local.sentinel_ti_alerts_table_name
  body = {
    properties = {
      schema = {
        name    = local.sentinel_ti_alerts_table_name
        columns = local.sentinel_ti_alerts_columns
      }
      retentionInDays      = 30
      totalRetentionInDays = 30
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```

### data_collection_logs

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
  default = "westeurope"
}

locals {
  data_collection_log_table_name = "DataCollectionLog_CL"
  data_collection_log_columns = [
    {
      "name" : "RawData",
      "type" : "string"
    },
    {
      "name" : "FilePath",
      "type" : "string"
    },
    {
      "name" : "TimeGenerated",
      "type" : "datetime"
    }
  ]
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "table" {
  type      = "Microsoft.OperationalInsights/workspaces/tables@2022-10-01"
  parent_id = azapi_resource.workspace.id
  name      = local.data_collection_log_table_name
  body = {
    properties = {
      schema = {
        name    = local.data_collection_log_table_name
        columns = local.data_collection_log_columns
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.OperationalInsights/workspaces/tables@api-version`. The available api-versions for this resource are: [`2020-03-01-preview`, `2020-08-01`, `2020-10-01`, `2021-12-01-preview`, `2022-10-01`, `2023-09-01`, `2025-02-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.OperationalInsights/workspaces/tables?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{resourceName}/tables/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.OperationalInsights/workspaces/{resourceName}/tables/{resourceName}?api-version=2025-02-01
 ```
