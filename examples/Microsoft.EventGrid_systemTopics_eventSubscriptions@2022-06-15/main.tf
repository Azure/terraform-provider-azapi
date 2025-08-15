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
