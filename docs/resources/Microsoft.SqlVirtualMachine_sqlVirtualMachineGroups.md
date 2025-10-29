---
subcategory: "Microsoft.SqlVirtualMachine - SQL Server on Azure Virtual Machines"
page_title: "sqlVirtualMachineGroups"
description: |-
  Manages a Microsoft SQL Virtual Machine Group.
---

# Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups - Microsoft SQL Virtual Machine Group

This article demonstrates how to use `azapi` provider to manage the Microsoft SQL Virtual Machine Group resource in Azure.



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

resource "azapi_resource" "sqlVirtualMachineGroup" {
  type      = "Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups@2023-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      sqlImageOffer = "SQL2017-WS2016"
      sqlImageSku   = "Developer"
      wsfcDomainProfile = {
        clusterBootstrapAccount  = ""
        clusterOperatorAccount   = ""
        clusterSubnetType        = "SingleSubnet"
        domainFqdn               = "testdomain.com"
        ouPath                   = ""
        sqlServiceAccount        = ""
        storageAccountPrimaryKey = ""
        storageAccountUrl        = ""
      }
    }
  }
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups@api-version`. The available api-versions for this resource are: [`2017-03-01-preview`, `2021-11-01-preview`, `2022-02-01`, `2022-02-01-preview`, `2022-07-01-preview`, `2022-08-01-preview`, `2023-01-01-preview`, `2023-10-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/{resourceName}?api-version=2023-10-01
 ```
