---
subcategory: "Microsoft.EventGrid - Event Grid"
page_title: "systemTopics/eventSubscriptions"
description: |-
  Manages a EventGrid System Topic Event Subscription.
---

# Microsoft.EventGrid/systemTopics/eventSubscriptions - EventGrid System Topic Event Subscription

This article demonstrates how to use `azapi` provider to manage the EventGrid System Topic Event Subscription resource in Azure.



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

locals {
  system_topic_name        = "${var.resource_name}-st"
  storage_account_name     = "${var.resource_name}sa01"
  queue_name               = "${var.resource_name}queue"
  event_subscription1_name = "${var.resource_name}-es1"
  event_subscription2_name = "${var.resource_name}-es2"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "systemTopic" {
  type      = "Microsoft.EventGrid/systemTopics@2022-06-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.system_topic_name
  # For resource group source, system topic location must be global
  location = "global"
  body = {
    properties = {
      source    = azapi_resource.resourceGroup.id
      topicType = "Microsoft.Resources.ResourceGroups"
    }
  }
}

resource "azapi_resource" "storageAccount" {
  type      = "Microsoft.Storage/storageAccounts@2023-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = local.storage_account_name
  location  = var.location
  body = {
    kind = "StorageV2"
    properties = {
      accessTier                   = "Hot"
      allowBlobPublicAccess        = true
      allowCrossTenantReplication  = false
      allowSharedKeyAccess         = true
      defaultToOAuthAuthentication = false
      dnsEndpointType              = "Standard"
      encryption = {
        keySource = "Microsoft.Storage"
        services = {
          queue = { keyType = "Service" }
          table = { keyType = "Service" }
        }
      }
      isHnsEnabled       = false
      isLocalUserEnabled = true
      isNfsV3Enabled     = false
      isSftpEnabled      = false
      minimumTlsVersion  = "TLS1_2"
      networkAcls = {
        bypass              = "AzureServices"
        defaultAction       = "Allow"
        ipRules             = []
        resourceAccessRules = []
        virtualNetworkRules = []
      }
      publicNetworkAccess      = "Enabled"
      supportsHttpsTrafficOnly = true
    }
    sku = { name = "Standard_LRS" }
  }
  tags = { environment = "staging" }
}

# Create a queue in the storage account for the event subscription destination
locals { queue_service_id = "${azapi_resource.storageAccount.id}/queueServices/default" }

resource "azapi_resource" "queue" {
  type       = "Microsoft.Storage/storageAccounts/queueServices/queues@2023-05-01"
  parent_id  = local.queue_service_id
  name       = local.queue_name
  body       = {}
  depends_on = [azapi_resource.storageAccount]
}

resource "azapi_resource" "eventSubscription" {
  type      = "Microsoft.EventGrid/systemTopics/eventSubscriptions@2022-06-15"
  parent_id = azapi_resource.systemTopic.id
  name      = local.event_subscription1_name
  body = {
    properties = {
      deadLetterDestination = null
      destination = {
        endpointType = "StorageQueue"
        properties = {
          queueName  = local.queue_name
          resourceId = azapi_resource.storageAccount.id
        }
      }
      eventDeliverySchema = "EventGridSchema"
      filter = {
        advancedFilters = [{
          key          = "subject"
          operatorType = "StringBeginsWith"
          values       = ["foo"]
        }]
      }
      labels = []
    }
  }
  depends_on = [azapi_resource.queue]
}

resource "azapi_resource" "eventSubscription_1" {
  type      = "Microsoft.EventGrid/systemTopics/eventSubscriptions@2022-06-15"
  parent_id = azapi_resource.systemTopic.id
  name      = local.event_subscription2_name
  body = {
    properties = {
      deadLetterDestination = null
      destination = {
        endpointType = "StorageQueue"
        properties = {
          queueName  = local.queue_name
          resourceId = azapi_resource.storageAccount.id
        }
      }
      eventDeliverySchema = "EventGridSchema"
      filter = {
        advancedFilters = [{
          key          = "subject"
          operatorType = "StringEndsWith"
          values       = ["bar"]
        }]
      }
      labels = []
    }
  }
  depends_on = [azapi_resource.queue]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.EventGrid/systemTopics/eventSubscriptions@api-version`. The available api-versions for this resource are: [`2020-04-01-preview`, `2020-10-15-preview`, `2021-06-01-preview`, `2021-10-15-preview`, `2021-12-01`, `2022-06-15`, `2023-06-01-preview`, `2023-12-15-preview`, `2024-06-01-preview`, `2024-12-15-preview`, `2025-02-15`, `2025-04-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/systemTopics/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.EventGrid/systemTopics/eventSubscriptions?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/systemTopics/{resourceName}/eventSubscriptions/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventGrid/systemTopics/{resourceName}/eventSubscriptions/{resourceName}?api-version=2025-04-01-preview
 ```
