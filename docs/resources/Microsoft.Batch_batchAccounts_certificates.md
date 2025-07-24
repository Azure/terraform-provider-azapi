---
subcategory: "Microsoft.Batch - Batch"
page_title: "batchAccounts/certificates"
description: |-
  Manages a certificate in an Azure Batch account.
---

# Microsoft.Batch/batchAccounts/certificates - certificate in an Azure Batch account

This article demonstrates how to use `azapi` provider to manage the certificate in an Azure Batch account resource in Azure.

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

variable "certificate_data" {
  type        = string
  description = "The base64-encoded certificate data"
  sensitive   = true
}

variable "certificate_thumbprint" {
  type        = string
  description = "The thumbprint of the certificate"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "batchAccount" {
  type      = "Microsoft.Batch/batchAccounts@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.Batch"
      }
      poolAllocationMode  = "BatchService"
      publicNetworkAccess = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "certificate" {
  type      = "Microsoft.Batch/batchAccounts/certificates@2022-10-01"
  parent_id = azapi_resource.batchAccount.id
  name      = "SHA1-${var.certificate_thumbprint}"
  body = {
    properties = {
      data                = var.certificate_data
      format              = "Cer"
      thumbprint          = var.certificate_thumbprint
      thumbprintAlgorithm = "sha1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Batch/batchAccounts/certificates@api-version`. The available api-versions for this resource are: [`2017-09-01`, `2018-12-01`, `2019-04-01`, `2019-08-01`, `2020-03-01`, `2020-05-01`, `2020-09-01`, `2021-01-01`, `2021-06-01`, `2022-01-01`, `2022-06-01`, `2022-10-01`, `2023-05-01`, `2023-11-01`, `2024-02-01`, `2024-07-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Batch/batchAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Batch/batchAccounts/certificates?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Batch/batchAccounts/{resourceName}/certificates/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Batch/batchAccounts/{resourceName}/certificates/{resourceName}?api-version=2024-07-01
 ```
