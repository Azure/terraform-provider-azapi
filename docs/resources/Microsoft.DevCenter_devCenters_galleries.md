---
subcategory: "Microsoft.DevCenter - Microsoft Dev Box"
page_title: "devCenters/galleries"
description: |-
  Manages a Dev Center Gallery.
---

# Microsoft.DevCenter/devCenters/galleries - Dev Center Gallery

This article demonstrates how to use `azapi` provider to manage the Dev Center Gallery resource in Azure.



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

data "azapi_client_config" "current" {}

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

resource "azapi_resource" "gallery" {
  type      = "Microsoft.Compute/galleries@2022-03-03"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}sig"
  location  = var.location
  body = {
    properties = {
      description = ""
    }
  }
}

resource "azapi_resource" "userAssignedIdentity" {
  type      = "Microsoft.ManagedIdentity/userAssignedIdentities@2023-01-31"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}uai"
  location  = var.location
}

resource "azapi_resource" "devCenter" {
  type      = "Microsoft.DevCenter/devCenters@2025-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}dc"
  location  = var.location
  identity {
    type         = "UserAssigned"
    identity_ids = [azapi_resource.userAssignedIdentity.id]
  }
  body = {
    properties = {}
  }
}

# Assign Reader role on the gallery to the Dev Center's user-assigned identity
locals {
  reader_role_id = "acdd72a7-3385-48ef-bd42-f606fba81ae7"
}

resource "azapi_resource" "roleAssignment" {
  type      = "Microsoft.Authorization/roleAssignments@2022-04-01"
  name      = uuidv5("url", "${azapi_resource.gallery.id}/roleAssignments/${azapi_resource.userAssignedIdentity.output.properties.principalId}")
  parent_id = azapi_resource.gallery.id
  body = {
    properties = {
      roleDefinitionId = "/subscriptions/${data.azapi_client_config.current.subscription_id}/providers/Microsoft.Authorization/roleDefinitions/${local.reader_role_id}"
      principalId      = azapi_resource.userAssignedIdentity.output.properties.principalId
      principalType    = "ServicePrincipal"
    }
  }
  depends_on = [azapi_resource.userAssignedIdentity, azapi_resource.gallery]
}

resource "azapi_resource" "gallery_1" {
  type      = "Microsoft.DevCenter/devCenters/galleries@2025-02-01"
  parent_id = azapi_resource.devCenter.id
  name      = "${var.resource_name}dcg"
  body = {
    properties = {
      galleryResourceId = azapi_resource.gallery.id
    }
  }
  depends_on = [azapi_resource.roleAssignment]
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.DevCenter/devCenters/galleries@api-version`. The available api-versions for this resource are: [`2022-08-01-preview`, `2022-09-01-preview`, `2022-10-12-preview`, `2022-11-11-preview`, `2023-01-01-preview`, `2023-04-01`, `2023-08-01-preview`, `2023-10-01-preview`, `2024-02-01`, `2024-05-01-preview`, `2024-06-01-preview`, `2024-07-01-preview`, `2024-08-01-preview`, `2024-10-01-preview`, `2025-02-01`, `2025-04-01-preview`, `2025-07-01-preview`, `2025-10-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.DevCenter/devCenters/galleries?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{resourceName}/galleries/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DevCenter/devCenters/{resourceName}/galleries/{resourceName}?api-version=2025-10-01-preview
 ```
