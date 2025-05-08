---
subcategory: "Microsoft.Network - Various Networking Services"
page_title: "trafficManagerProfiles/NestedEndpoints"
description: |-
  Manages a Nested Endpoint within a Traffic Manager Profile.
---

# Microsoft.Network/trafficManagerProfiles/NestedEndpoints - Nested Endpoint within a Traffic Manager Profile

This article demonstrates how to use `azapi` provider to manage the Nested Endpoint within a Traffic Manager Profile resource in Azure.

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

resource "azapi_resource" "trafficManagerProfile" {
  type      = "Microsoft.Network/trafficManagerProfiles@2018-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      dnsConfig = {
        relativeName = "acctest-tmp-230630034107605443"
        ttl          = 30
      }
      monitorConfig = {
        expectedStatusCodeRanges = [
        ]
        intervalInSeconds         = 30
        path                      = "/"
        port                      = 443
        protocol                  = "HTTPS"
        timeoutInSeconds          = 10
        toleratedNumberOfFailures = 3
      }
      trafficRoutingMethod = "Weighted"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "trafficManagerProfile2" {
  type      = "Microsoft.Network/trafficManagerProfiles@2018-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = "global"
  body = {
    properties = {
      dnsConfig = {
        relativeName = "acctesttmpchild230630034107605443"
        ttl          = 30
      }
      monitorConfig = {
        expectedStatusCodeRanges = [
        ]
        intervalInSeconds         = 30
        path                      = "/"
        port                      = 443
        protocol                  = "HTTPS"
        timeoutInSeconds          = 10
        toleratedNumberOfFailures = 3
      }
      trafficRoutingMethod = "Priority"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "NestedEndpoint" {
  type      = "Microsoft.Network/trafficManagerProfiles/NestedEndpoints@2018-08-01"
  parent_id = azapi_resource.trafficManagerProfile.id
  name      = var.resource_name
  body = {
    properties = {
      customHeaders = [
      ]
      endpointStatus    = "Enabled"
      minChildEndpoints = 5
      subnets = [
      ]
      targetResourceId = azapi_resource.trafficManagerProfile2.id
      weight           = 3
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Network/trafficManagerProfiles/NestedEndpoints@api-version`. The available api-versions for this resource are: [`2018-08-01`, `2022-04-01`, `2022-04-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Network/trafficManagerProfiles/NestedEndpoints?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}/NestedEndpoints/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/trafficManagerProfiles/{resourceName}/NestedEndpoints/{resourceName}?api-version=2022-04-01-preview
 ```
