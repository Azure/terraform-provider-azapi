---
subcategory: ""
layout: "azapi"
page_title: "Azure Resource: azapi_resource"
description: |-
  Manages a Azure resource
---

# azapi_resource

This resource can manage any Azure resource manager resource.

## Example Usage

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_user_assigned_identity" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
}

// manage a container registry resource
resource "azapi_resource" "example" {
  type      = "Microsoft.ContainerRegistry/registries@2020-11-01-preview"
  name      = "registry1"
  parent_id = azurerm_resource_group.example.id

  location = azurerm_resource_group.example.location
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  body = jsonencode({
    sku = {
      name = "Standard"
    }
    properties = {
      adminUserEnabled = true
    }
  })

  tags = {
    "Key" = "Value"
  }

  response_export_values = ["properties.loginServer", "properties.policies.quarantinePolicy.status"]
}

// it will output "registry1.azurecr.io"
output "login_server" {
  value = jsondecode(azapi_resource.example.output).properties.loginServer
}

// it will output "disabled"
output "quarantine_policy" {
  value = jsondecode(azapi_resource.example.output).properties.policies.quarantinePolicy.status
}
```

## Arguments Reference

The following arguments are supported:
* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created. 
* `parent_id` - (Required) The ID of the azure resource in which this resource is created. Changing this forces a new resource to be created. It supports different kinds of deployment scope for **top level** resources: 
    - resource group scope: `parent_id` should be the ID of a resource group, it's recommended to manage a resource group by [azurerm_resource_group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/resource_group).
    - management group scope: `parent_id` should be the ID of a management group, it's recommended to manage a management group by [azurerm_management_group](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/resources/management_group).
    - extension scope: `parent_id` should be the ID of the resource you're adding the extension to.
    - subscription scope: `parent_id` should be like `/subscriptions/00000000-0000-0000-0000-000000000000`
    - tenant scope: `parent_id` should be `/`

  For child level resources, the `parent_id` should be the ID of its parent resource, for example, subnet resource's `parent_id` is the ID of the vnet.

* `type` - (Required) It is in a format like `<resource-type>@<api-version>`. `<resource-type>` is the Azure resource type, for example, `Microsoft.Storage/storageAccounts`.
  `<api-version>` is version of the API used to manage this azure resource.

* `body` - (Required) A JSON object that contains the request body used to create and update azure resource. 

---
  
* `location` - (Optional) The Azure Region where the azure resource should exist. 
  
* `identity` - (Optional) A `identity` block as defined below. 

* `tags` - (Optional) A mapping of tags which should be assigned to the azure resource. 

* `response_export_values` - (Optional) A list of path that needs to be exported from response body.
  Setting it to `["*"]` will export the full response body.
  Here's an example. If it sets to `["properties.loginServer", "properties.policies.quarantinePolicy.status"]`, it will set the following json to computed property `output`.
```
{
  "properties" : {
    "loginServer" : "registry1.azurecr.io"
    "policies" : {
      "quarantinePolicy" = {
        "status" = "disabled"
      }
    }
  }
}
```

* `locks` - (Optional) A list of ARM resource IDs which are used to avoid create/modify/delete azapi resources at the same time.

* `ignore_casing` - (Optional) Whether ignore incorrect casing returned in `body` to suppress plan-diff. Defaults to `false`.

* `ignore_missing_property` - (Optional) Whether ignore not returned properties like credentials in `body` to suppress plan-diff. Defaults to `false`.

* `schema_validation_enabled` - (Optional) Whether enabled the validation on `type` and `body` with embedded schema. Defaults to `true`.

---

A `identity` block supports the following:

* `type` - (Required) The Type of Identity which should be used for this azure resource. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned,UserAssigned`. 

* `identity_ids` - (Optional) A list of User Managed Identity ID's which should be assigned to the azure resource. 


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this azure resource.

* `output` - The output json containing the properties specified in `response_export_values`. Here're some examples to decode json and extract the value.
```
// it will output "registry1.azurecr.io"
output "login_server" {
  value = jsondecode(azapi_resource.example.output).properties.loginServer
}

// it will output "disabled"
output "quarantine_policy" {
  value = jsondecode(azapi_resource.example.output).properties.policies.quarantinePolicy.status
}
```


---

A `identity` block exports the following:

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this azure resource.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this azure resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the azure resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the azure resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the azure resource.

## Import

Azure resource can be imported using the `resource id`, e.g.

```shell
terraform import azapi_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1
```

It also supports specifying API version by using the `resource id` with `api-version` as a query parameter, e.g.

```shell
terraform import azapi_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1?api-version=2021-07-01
```
