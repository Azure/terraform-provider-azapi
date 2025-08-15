---
subcategory: "Microsoft.ServiceBus - Service Bus"
page_title: "namespaces/disasterRecoveryConfigs"
description: |-
  Manages a Disaster Recovery Config for a Service Bus Namespace.
---

# Microsoft.ServiceBus/namespaces/disasterRecoveryConfigs - Disaster Recovery Config for a Service Bus Namespace

This article demonstrates how to use `azapi` provider to manage the Disaster Recovery Config for a Service Bus Namespace resource in Azure.

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
  default = "westus"
}

variable "secondary_location" {
  type    = string
  default = "centralus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "resourceGroup_1" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = "${var.resource_name}rg2"
  location = var.secondary_location
}

resource "azapi_resource" "namespace" {
  type      = "Microsoft.ServiceBus/namespaces@2022-10-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}ns1"
  location  = var.location
  body = {
    properties = {
      disableLocalAuth           = false
      minimumTlsVersion          = "1.2"
      premiumMessagingPartitions = 1
      publicNetworkAccess        = "Enabled"
    }
    sku = {
      capacity = 1
      name     = "Premium"
      tier     = "Premium"
    }
  }
}

resource "azapi_resource" "namespace_1" {
  type      = "Microsoft.ServiceBus/namespaces@2022-10-01-preview"
  parent_id = azapi_resource.resourceGroup_1.id
  name      = "${var.resource_name}ns2"
  location  = var.secondary_location
  body = {
    properties = {
      disableLocalAuth           = false
      minimumTlsVersion          = "1.2"
      premiumMessagingPartitions = 1
      publicNetworkAccess        = "Enabled"
    }
    sku = {
      capacity = 1
      name     = "Premium"
      tier     = "Premium"
    }
  }
}

resource "azapi_resource" "disasterRecoveryConfig" {
  type      = "Microsoft.ServiceBus/namespaces/disasterRecoveryConfigs@2021-06-01-preview"
  parent_id = azapi_resource.namespace.id
  name      = "${var.resource_name}alias"
  body = {
    properties = {
      partnerNamespace = azapi_resource.namespace_1.id
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ServiceBus/namespaces/disasterRecoveryConfigs@api-version`. The available api-versions for this resource are: [`2017-04-01`, `2018-01-01-preview`, `2021-01-01-preview`, `2021-06-01-preview`, `2021-11-01`, `2022-01-01-preview`, `2022-10-01-preview`, `2023-01-01-preview`, `2024-01-01`, `2025-05-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ServiceBus/namespaces/disasterRecoveryConfigs?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}/disasterRecoveryConfigs/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceBus/namespaces/{resourceName}/disasterRecoveryConfigs/{resourceName}?api-version=2025-05-01-preview
 ```
