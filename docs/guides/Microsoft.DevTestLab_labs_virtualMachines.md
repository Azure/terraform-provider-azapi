---
subcategory: "Microsoft.DevTestLab - Azure Lab Services"
page_title: "labs/virtualMachines"
description: |-
  Manages a Virtual Machine within a Dev Test Lab.
---

# Microsoft.DevTestLab/labs/virtualMachines - Virtual Machine within a Dev Test Lab

This article demonstrates how to use `azapi` provider to manage the Virtual Machine within a Dev Test Lab resource in Azure.

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
  type      = "Microsoft.DevTestLab/labs@2018-09-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      labStorageType = "Premium"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

data "azapi_resource_id" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2023-04-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
}

data "azapi_resource_id" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2023-04-01"
  parent_id = data.azapi_resource_id.virtualNetwork.id
  name      = "${var.resource_name}Subnet"
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.DevTestLab/labs/virtualNetworks@2018-09-15"
  parent_id = azapi_resource.lab.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      subnetOverrides = [
        {
          labSubnetName                = data.azapi_resource_id.subnet.name
          resourceId                   = data.azapi_resource_id.subnet.id
          useInVmCreationPermission    = "Allow"
          usePublicIpAddressPermission = "Allow"
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualMachine" {
  type      = "Microsoft.DevTestLab/labs/virtualMachines@2018-09-15"
  parent_id = azapi_resource.lab.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      allowClaim              = true
      disallowPublicIpAddress = false
      galleryImageReference = {
        offer     = "WindowsServer"
        osType    = "Windows"
        publisher = "MicrosoftWindowsServer"
        sku       = "2012-Datacenter"
        version   = "latest"
      }
      isAuthenticationWithSshKey = false
      labSubnetName              = data.azapi_resource_id.subnet.name
      labVirtualNetworkId        = azapi_resource.virtualNetwork.id
      networkInterface = {
      }
      notes       = ""
      osType      = "Windows"
      password    = "Pa$w0rd1234!"
      size        = "Standard_F2"
      storageType = "Standard"
      userName    = "acct5stU5er"
    }
  }
  ignore_casing             = true
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DevTestLab/labs/virtualMachines@api-version`. The available api-versions for this resource are: [`2015-05-21-preview`, `2016-05-15`, `2018-09-15`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DevTestLab/labs/virtualMachines?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{resourceName}/virtualMachines/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevTestLab/labs/{resourceName}/virtualMachines/{resourceName}?api-version=2018-09-15
 ```
