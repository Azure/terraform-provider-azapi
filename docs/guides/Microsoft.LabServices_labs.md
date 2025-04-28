---
subcategory: "Microsoft.LabServices - Azure Lab Services"
page_title: "labs"
description: |-
  Manages a Lab Services Labs.
---

# Microsoft.LabServices/labs - Lab Services Labs

This article demonstrates how to use `azapi` provider to manage the Lab Services Labs resource in Azure.

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

resource "azapi_resource" "lab" {
  type      = "Microsoft.LabServices/labs@2022-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      autoShutdownProfile = {
        shutdownOnDisconnect     = "Disabled"
        shutdownOnIdle           = "None"
        shutdownWhenNotConnected = "Disabled"
      }
      connectionProfile = {
        clientRdpAccess = "None"
        clientSshAccess = "None"
        webRdpAccess    = "None"
        webSshAccess    = "None"
      }
      securityProfile = {
        openAccess = "Disabled"
      }
      title = "Test Title"
      virtualMachineProfile = {
        additionalCapabilities = {
          installGpuDrivers = "Disabled"
        }
        adminUser = {
          password = "Password1234!"
          username = "testadmin"
        }
        createOption = "Image"
        imageReference = {
          offer     = "0001-com-ubuntu-server-focal"
          publisher = "canonical"
          sku       = "20_04-lts"
          version   = "latest"
        }
        sku = {
          capacity = 1
          name     = "Classic_Fsv2_2_4GB_128_S_SSD"
        }
        usageQuota        = "PT0S"
        useSharedPassword = "Disabled"
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.LabServices/labs@api-version`. The available api-versions for this resource are: [`2021-10-01-preview`, `2021-11-15-preview`, `2022-08-01`, `2023-06-07`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.LabServices/labs?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.LabServices/labs/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.LabServices/labs/{resourceName}?api-version=2023-06-07
 ```
