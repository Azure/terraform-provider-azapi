---
subcategory: "Microsoft.CognitiveServices - Cognitive Services"
page_title: "accounts/raiPolicies"
description: |-
  Manages a Cognitive Services Accounts Responsible AI Policies.
---

# Microsoft.CognitiveServices/accounts/raiPolicies - Cognitive Services Accounts Responsible AI Policies

This article demonstrates how to use `azapi` provider to manage the Cognitive Services Accounts Responsible AI Policies resource in Azure.

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
  default = "acctest0003"
}

variable "location" {
  type    = string
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "account" {
  type      = "Microsoft.CognitiveServices/accounts@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location

  body = {
    kind = "OpenAI"
    properties = {
      disableLocalAuth              = false
      dynamicThrottlingEnabled      = false
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = false
    }
    sku = {
      name = "S0"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "raiPolicy" {
  type      = "Microsoft.CognitiveServices/accounts/raiPolicies@2024-10-01"
  name      = "NoModerationPolicy"
  parent_id = azapi_resource.account.id

  body = {
    properties = {
      basePolicyName = "Microsoft.Default"
      contentFilters = [
        {
          blocking          = true
          enabled           = true
          name              = "Hate"
          severityThreshold = "High"
          source            = "Prompt"
        }
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.CognitiveServices/accounts/raiPolicies@api-version`. The available api-versions for this resource are: [`2023-10-01-preview`, `2024-04-01-preview`, `2024-06-01-preview`, `2024-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CognitiveServices/accounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.CognitiveServices/accounts/raiPolicies?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CognitiveServices/accounts/{resourceName}/raiPolicies/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.CognitiveServices/accounts/{resourceName}/raiPolicies/{resourceName}?api-version=2024-10-01
 ```
