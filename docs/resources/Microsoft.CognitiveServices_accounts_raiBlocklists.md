---
subcategory: "Microsoft.CognitiveServices - Cognitive Services"
page_title: "accounts/raiBlocklists"
description: |-
  Manages a Cognitive Services Accounts Responsible AI Blocklists.
---

# Microsoft.CognitiveServices/accounts/raiBlocklists - Cognitive Services Accounts Responsible AI Blocklists

This article demonstrates how to use `azapi` provider to manage the Cognitive Services Accounts Responsible AI Blocklists resource in Azure.

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

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "account" {
  type      = "Microsoft.CognitiveServices/accounts@2024-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-ca"
  location  = var.location
  body = {
    kind = "OpenAI"
    properties = {
      allowedFqdnList               = []
      apiProperties                 = {}
      customSubDomainName           = ""
      disableLocalAuth              = false
      dynamicThrottlingEnabled      = false
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = false
    }
    sku = {
      name = "S0"
    }
  }
}

resource "azapi_resource" "raiBlocklist" {
  type      = "Microsoft.CognitiveServices/accounts/raiBlocklists@2024-10-01"
  parent_id = azapi_resource.account.id
  name      = "${var.resource_name}-crb"
  body = {
    properties = {
      description = "Acceptance test data new azurerm resource"
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.CognitiveServices/accounts/raiBlocklists@api-version`. The available api-versions for this resource are: [`2023-10-01-preview`, `2024-04-01-preview`, `2024-06-01-preview`, `2024-10-01`, `2025-04-01-preview`, `2025-06-01`, `2025-07-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CognitiveServices/accounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.CognitiveServices/accounts/raiBlocklists?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CognitiveServices/accounts/{resourceName}/raiBlocklists/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CognitiveServices/accounts/{resourceName}/raiBlocklists/{resourceName}?api-version=2025-07-01-preview
 ```
