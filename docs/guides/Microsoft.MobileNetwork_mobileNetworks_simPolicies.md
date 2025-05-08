---
subcategory: "Microsoft.MobileNetwork - Azure Private 5G Core"
page_title: "mobileNetworks/simPolicies"
description: |-
  Manages a Mobile Network Sim Policy.
---

# Microsoft.MobileNetwork/mobileNetworks/simPolicies - Mobile Network Sim Policy

This article demonstrates how to use `azapi` provider to manage the Mobile Network Sim Policy resource in Azure.

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
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "mobileNetwork" {
  type      = "Microsoft.MobileNetwork/mobileNetworks@2022-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicLandMobileNetworkIdentifier = {
        mcc = "001"
        mnc = "01"
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "service" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/services@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      pccRules = [
        {
          ruleName       = "default-rule"
          rulePrecedence = 1
          serviceDataFlowTemplates = [
            {
              direction = "Uplink"
              ports = [
              ]
              protocol = [
                "ip",
              ]
              remoteIpList = [
                "10.3.4.0/24",
              ]
              templateName = "IP-to-server"
            },
          ]
          trafficControl = "Enabled"
        },
      ]
      servicePrecedence = 0
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "dataNetwork" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/dataNetworks@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "slice" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/slices@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      snssai = {
        sst = 1
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "simPolicy" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/simPolicies@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      defaultSlice = {
        id = azapi_resource.slice.id
      }
      registrationTimer = 3240
      sliceConfigurations = [
        {
          dataNetworkConfigurations = [
            {
              "5qi"                               = 9
              additionalAllowedSessionTypes       = null
              allocationAndRetentionPriorityLevel = 9
              allowedServices = [
                {
                  id = azapi_resource.service.id
                },
              ]
              dataNetwork = {
                id = azapi_resource.dataNetwork.id
              }
              defaultSessionType             = "IPv4"
              maximumNumberOfBufferedPackets = 10
              preemptionCapability           = "NotPreempt"
              preemptionVulnerability        = "Preemptable"
              sessionAmbr = {
                downlink = "1 Gbps"
                uplink   = "500 Mbps"
              }
            },
          ]
          defaultDataNetwork = {
            id = azapi_resource.dataNetwork.id
          }
          slice = {
            id = azapi_resource.slice.id
          }
        },
      ]
      ueAmbr = {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
    tags = {
      key = "value"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.MobileNetwork/mobileNetworks/simPolicies@api-version`. The available api-versions for this resource are: [`2022-03-01-preview`, `2022-04-01-preview`, `2022-11-01`, `2023-06-01`, `2023-09-01`, `2024-02-01`, `2024-04-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/mobileNetworks/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.MobileNetwork/mobileNetworks/simPolicies?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/mobileNetworks/{resourceName}/simPolicies/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.MobileNetwork/mobileNetworks/{resourceName}/simPolicies/{resourceName}?api-version=2024-04-01
 ```
