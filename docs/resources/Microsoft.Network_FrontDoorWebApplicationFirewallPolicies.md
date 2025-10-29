---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "FrontDoorWebApplicationFirewallPolicies"
description: |-
  Manages a Front Door Firewall Policy instance.
---

# Microsoft.Network/FrontDoorWebApplicationFirewallPolicies - Front Door Firewall Policy instance

This article demonstrates how to use `azapi` provider to manage the Front Door Firewall Policy instance resource in Azure.



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

resource "azapi_resource" "FrontDoorWebApplicationFirewallPolicy" {
  type      = "Microsoft.Network/FrontDoorWebApplicationFirewallPolicies@2020-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      customRules = {
        rules = [
          {
            action       = "Block"
            enabledState = "Enabled"
            matchConditions = [
              {
                matchValue = [
                  "192.168.1.0/24",
                  "10.0.0.0/24",
                ]
                matchVariable   = "RemoteAddr"
                negateCondition = false
                operator        = "IPMatch"
              },
            ]
            name                       = "Rule1"
            priority                   = 1
            rateLimitDurationInMinutes = 1
            rateLimitThreshold         = 10
            ruleType                   = "MatchRule"
          },
        ]
      }
      managedRules = {
        managedRuleSets = [
          {
            ruleGroupOverrides = [
              {
                ruleGroupName = "PHP"
                rules = [
                  {
                    action       = "Block"
                    enabledState = "Disabled"
                    ruleId       = "933111"
                  },
                ]
              },
            ]
            ruleSetAction  = "Block"
            ruleSetType    = "DefaultRuleSet"
            ruleSetVersion = "preview-0.1"
          },
          {
            ruleSetAction  = "Block"
            ruleSetType    = "BotProtection"
            ruleSetVersion = "preview-0.1"
          },
        ]
      }
      policySettings = {
        customBlockResponseBody       = "PGh0bWw+CjxoZWFkZXI+PHRpdGxlPkhlbGxvPC90aXRsZT48L2hlYWRlcj4KPGJvZHk+CkhlbGxvIHdvcmxkCjwvYm9keT4KPC9odG1sPg=="
        customBlockResponseStatusCode = 403
        enabledState                  = "Enabled"
        mode                          = "Prevention"
        redirectUrl                   = "https://www.fabrikam.com"
      }
    }
    sku = {
      name = "Premium_AzureFrontDoor"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/FrontDoorWebApplicationFirewallPolicies@api-version`. The available api-versions for this resource are: [`2018-08-01`, `2019-03-01`, `2019-10-01`, `2020-04-01`, `2020-11-01`, `2021-06-01`, `2022-05-01`, `2024-02-01`, `2025-03-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/FrontDoorWebApplicationFirewallPolicies?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/FrontDoorWebApplicationFirewallPolicies/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/FrontDoorWebApplicationFirewallPolicies/{resourceName}?api-version=2025-03-01
 ```
