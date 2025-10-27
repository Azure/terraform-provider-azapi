---
subcategory: "Microsoft.Communication - Azure Communication Services"
page_title: "emailServices/domains"
description: |-
  Manages a Email Communication Service Domain.
---

# Microsoft.Communication/emailServices/domains - Email Communication Service Domain

This article demonstrates how to use `azapi` provider to manage the Email Communication Service Domain resource in Azure.



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

resource "azapi_resource" "emailService" {
  type      = "Microsoft.Communication/emailServices@2023-04-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      dataLocation = "United States"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "domain" {
  type      = "Microsoft.Communication/emailServices/domains@2023-04-01-preview"
  name      = "example.com"
  location  = "global"
  parent_id = azapi_resource.emailService.id
  tags = {
    env = "Test"
  }
  body = {
    properties = {
      domainManagement       = "CustomerManaged"
      userEngagementTracking = "Disabled"
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Communication/emailServices/domains@api-version`. The available api-versions for this resource are: [`2021-10-01-preview`, `2022-07-01-preview`, `2023-03-01-preview`, `2023-03-31`, `2023-04-01`, `2023-04-01-preview`, `2023-06-01-preview`, `2024-09-01-preview`, `2025-05-01`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Communication/emailServices/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Communication/emailServices/domains?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Communication/emailServices/{resourceName}/domains/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Communication/emailServices/{resourceName}/domains/{resourceName}?api-version=2025-05-01-preview
 ```
