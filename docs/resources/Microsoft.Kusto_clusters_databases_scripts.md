---
subcategory: "Microsoft.Kusto - Azure Data Explorer"
page_title: "clusters/databases/scripts"
description: |-
  Manages a Kusto Script.
---

# Microsoft.Kusto/clusters/databases/scripts - Kusto Script

This article demonstrates how to use `azapi` provider to manage the Kusto Script resource in Azure.



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
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "cluster" {
  type      = "Microsoft.Kusto/clusters@2023-05-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      enableAutoStop                = true
      enableDiskEncryption          = false
      enableDoubleEncryption        = false
      enablePurge                   = false
      enableStreamingIngest         = false
      engineType                    = "V2"
      publicIPType                  = "IPv4"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      trustedExternalTenants = [
      ]
    }
    sku = {
      capacity = 1
      name     = "Dev(No SLA)_Standard_D11_v2"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Kusto/clusters/databases@2023-05-02"
  parent_id = azapi_resource.cluster.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "ReadWrite"
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "script" {
  type      = "Microsoft.Kusto/clusters/databases/scripts@2023-05-02"
  parent_id = azapi_resource.database.id
  name      = "create-table-script"
  body = {
    properties = {
      continueOnErrors = false
      forceUpdateTag   = "9e2e7874-aa37-7041-81b7-06397f03a37d"
      scriptContent    = ".create table TestTable(Id:string, Name:string, _ts:long, _timestamp:datetime)\n.create table TestTable ingestion json mapping \"TestMapping\"\n'['\n'    {\"column\":\"Id\",\"path\":\"$.id\"},'\n'    {\"column\":\"Name\",\"path\":\"$.name\"},'\n'    {\"column\":\"_ts\",\"path\":\"$._ts\"},'\n'    {\"column\":\"_timestamp\",\"path\":\"$._ts\", \"transform\":\"DateTimeFromUnixSeconds\"}'\n']'\n.alter table TestTable policy ingestionbatching \"{'MaximumBatchingTimeSpan': '0:0:10', 'MaximumNumberOfItems': 10000}\"\n"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Kusto/clusters/databases/scripts@api-version`. The available api-versions for this resource are: [`2021-01-01`, `2021-08-27`, `2022-02-01`, `2022-07-07`, `2022-11-11`, `2022-12-29`, `2023-05-02`, `2023-08-15`, `2024-04-13`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{resourceName}/databases/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Kusto/clusters/databases/scripts?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{resourceName}/databases/{resourceName}/scripts/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Kusto/clusters/{resourceName}/databases/{resourceName}/scripts/{resourceName}?api-version=2024-04-13
 ```
