---
subcategory: "Microsoft.SignalRService - Azure SignalR Service"
page_title: "webPubSub/sharedPrivateLinkResources"
description: |-
  Manages a Shared Private Link Resource for a Web Pubsub service.
---

# Microsoft.SignalRService/webPubSub/sharedPrivateLinkResources - Shared Private Link Resource for a Web Pubsub service

This article demonstrates how to use `azapi` provider to manage the Shared Private Link Resource for a Web Pubsub service resource in Azure.

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
  default = "westeurope"
}

data "azurerm_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "webPubSub" {
  type      = "Microsoft.SignalRService/webPubSub@2023-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      disableAadAuth      = false
      disableLocalAuth    = false
      publicNetworkAccess = "Enabled"
      tls = {
        clientCertEnabled = false
      }
    }
    sku = {
      capacity = 1
      name     = "Standard_S1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "vault" {
  type      = "Microsoft.KeyVault/vaults@2021-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      accessPolicies = [
        {
          objectId = data.azurerm_client_config.current.object_id
          permissions = {
            certificates = [
              "ManageContacts",
            ]
            keys = [
              "Create",
            ]
            secrets = [
              "Set",
            ]
            storage = [
            ]
          }
          tenantId = data.azurerm_client_config.current.tenant_id
        },
      ]
      createMode                   = "default"
      enableRbacAuthorization      = false
      enableSoftDelete             = true
      enabledForDeployment         = false
      enabledForDiskEncryption     = false
      enabledForTemplateDeployment = false
      publicNetworkAccess          = "Enabled"
      sku = {
        family = "A"
        name   = "standard"
      }
      softDeleteRetentionInDays = 7
      tenantId                  = data.azurerm_client_config.current.tenant_id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "sharedPrivateLinkResource" {
  type      = "Microsoft.SignalRService/webPubSub/sharedPrivateLinkResources@2023-02-01"
  parent_id = azapi_resource.webPubSub.id
  name      = var.resource_name
  body = {
    properties = {
      groupId               = "vault"
      privateLinkResourceId = azapi_resource.vault.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.SignalRService/webPubSub/sharedPrivateLinkResources@api-version`. The available api-versions for this resource are: [`2021-04-01-preview`, `2021-06-01-preview`, `2021-09-01-preview`, `2021-10-01`, `2022-08-01-preview`, `2023-02-01`, `2023-03-01-preview`, `2023-06-01-preview`, `2023-08-01-preview`, `2024-01-01-preview`, `2024-03-01`, `2024-04-01-preview`, `2024-08-01-preview`, `2024-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SignalRService/webPubSub/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.SignalRService/webPubSub/sharedPrivateLinkResources?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SignalRService/webPubSub/{resourceName}/sharedPrivateLinkResources/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SignalRService/webPubSub/{resourceName}/sharedPrivateLinkResources/{resourceName}?api-version=2024-10-01-preview
 ```
