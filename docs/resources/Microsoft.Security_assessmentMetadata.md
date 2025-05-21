---
subcategory: "Microsoft.Security - Security Center"
page_title: "assessmentMetadata"
description: |-
  Manages a Security Center Assessment Metadata for Azure Security Center.
---

# Microsoft.Security/assessmentMetadata - Security Center Assessment Metadata for Azure Security Center

This article demonstrates how to use `azapi` provider to manage the Security Center Assessment Metadata for Azure Security Center resource in Azure.

## Example Usage

### default

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
    azurerm = {
      source = "hashicorp/azurerm"
    }
  }
}

provider "azurerm" {
  features {
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
  default = "eastus"
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "assessmentMetadatum" {
  type      = "Microsoft.Security/assessmentMetadata@2020-01-01"
  parent_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}"
  name      = "95c7a001-d595-43af-9754-1310c740d34c"
  body = {
    properties = {
      assessmentType = "CustomerManaged"
      description    = "Test Description"
      displayName    = "Test Display Name"
      severity       = "Medium"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Security/assessmentMetadata@api-version`. The available api-versions for this resource are: [`2019-01-01-preview`, `2020-01-01`, `2021-06-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/`  
  `/subscriptions/{subscriptionId}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Security/assessmentMetadata?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Security/assessmentMetadata/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example //providers/Microsoft.Security/assessmentMetadata/{resourceName}?api-version=2021-06-01
 ```
