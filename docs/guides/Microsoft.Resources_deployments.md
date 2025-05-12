---
subcategory: "Microsoft.Resources - Azure Resource Manager"
page_title: "deployments"
description: |-
  Manages a Template Deployment.
---

# Microsoft.Resources/deployments - Template Deployment

This article demonstrates how to use `azapi` provider to manage the Template Deployment resource in Azure.

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

resource "azapi_resource" "deployment" {
  type      = "Microsoft.Resources/deployments@2020-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  body = {
    properties = {
      mode = "Complete"
      template = {
        "$schema"      = "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#"
        contentVersion = "1.0.0.0"
        parameters = {
          storageAccountType = {
            allowedValues = [
              "Standard_LRS",
              "Standard_GRS",
              "Standard_ZRS",
            ]
            defaultValue = "Standard_LRS"
            metadata = {
              description = "Storage Account type"
            }
            type = "string"
          }
        }
        resources = [
          {
            apiVersion = "[variables('apiVersion')]"
            location   = "[variables('location')]"
            name       = "[variables('storageAccountName')]"
            properties = {
              accountType = "[parameters('storageAccountType')]"
            }
            type = "Microsoft.Storage/storageAccounts"
          },
          {
            apiVersion = "[variables('apiVersion')]"
            location   = "[variables('location')]"
            name       = "[variables('publicIPAddressName')]"
            properties = {
              dnsSettings = {
                domainNameLabel = "[variables('dnsLabelPrefix')]"
              }
              publicIPAllocationMethod = "[variables('publicIPAddressType')]"
            }
            type = "Microsoft.Network/publicIPAddresses"
          },
        ]
        variables = {
          apiVersion          = "2015-06-15"
          dnsLabelPrefix      = "[concat('terraform-tdacctest', uniquestring(resourceGroup().id))]"
          location            = "[resourceGroup().location]"
          publicIPAddressName = "[concat('myPublicIp', uniquestring(resourceGroup().id))]"
          publicIPAddressType = "Dynamic"
          storageAccountName  = "[concat(uniquestring(resourceGroup().id), 'storage')]"
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Resources/deployments@api-version`. The available api-versions for this resource are: [`2015-11-01`, `2016-02-01`, `2016-07-01`, `2016-09-01`, `2017-05-10`, `2018-02-01`, `2018-05-01`, `2019-03-01`, `2019-05-01`, `2019-05-10`, `2019-07-01`, `2019-08-01`, `2019-10-01`, `2020-06-01`, `2020-08-01`, `2020-10-01`, `2021-01-01`, `2021-04-01`, `2022-09-01`, `2023-07-01`, `2024-03-01`, `2024-07-01`, `2024-11-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Resources/deployments?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Resources/deployments/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Resources/deployments/{resourceName}?api-version=2024-11-01
 ```
