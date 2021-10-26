---
subcategory: ""
layout: "azurermg"
page_title: "Generic Azure Resource: azurermg_resource"
description: |-
  Manages a Azure resource
---

# azurermg_resource

This resource can manage any Azure resource manager resource.

## Example Usage

```hcl
terraform {
  required_providers {
    azurermg = {
      source = "ms-henglu/azurermg"
    }
  }
}

provider "azurermg" {
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
resource "azurermg_resource" "example" {
  url         = "${azurerm_resource_group.test.id}/providers/Microsoft.ContainerRegistry/registries/registry1"
  api_version = "2020-11-01-preview"
  location    = azurerm_resource_group.example.location
  identity {
    type         = "SystemAssigned, UserAssigned"
    identity_ids = [azurerm_user_assigned_identity.example.id]
  }

  body = <<BODY
    {
      "sku": {
        "name": "Standard"
      },
      "properties": {
        "adminUserEnabled": true
      }
    }
  BODY

  tags = {
    "Key" = "Value"
  }

  paths = ["properties.loginServer", "properties.policies.quarantinePolicy.status"]
}

// it will output "registry1.azurecr.io"
output "login_server" {
  value = jsondecode(azurermg_resource.test.output).properties.loginServer
}

// it will output "disabled"
output "quarantine_policy" {
  value = jsondecode(azurermg_resource.test.output).properties.policies.quarantinePolicy.status
}
```

## Arguments Reference

The following arguments are supported:
* `url` - (Required) The url which should be used to manage this azure resource. It's also same with the ID of an azure source. 
  Here're some examples 
  `Container Registry: /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1` and 
  `Virtual Machine: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1`.
  Changing this forces a new azure resource to be created.

* `api_version` - (Required) The version of the API used to manage this azure resource. Changing this forces a new azure resource to be created.

* `body` - (Required) A JSON object that contains the request body used to create and update azure resource. 

---

* `create_method` - (Optional) The HTTP method used to create this azure resource. Possible values are `PUT` and `POST`. Defaults to `PUT`.

* `update_method` - (Optional) The HTTP method used to create this azure resource. Possible values are `PUT` and `POST`. Defaults to `PUT`.
  
* `location` - (Optional) The Azure Region where the azure resource should exist. 
  
* `identity` - (Optional) A `identity` block as defined below. 

* `tags` - (Optional) A mapping of tags which should be assigned to the azure resource. 

* `paths` - (Optional) A list of path that needs to be exported from response body. Here's an example. 
  If it sets to `["properties.loginServer", "properties.policies.quarantinePolicy.status"]`, it will set the following json to computed property `output`.
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
---

A `identity` block supports the following:

* `type` - (Required) The Type of Identity which should be used for this azure resource. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned,UserAssigned`. 

* `identity_ids` - (Optional) A list of User Managed Identity ID's which should be assigned to the azure resource. 


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this azure resource.

* `output` - The output json containing the properties specified in `paths`. Here're some examples to decode json and extract the value.
```
// it will output "registry1.azurecr.io"
output "login_server" {
  value = jsondecode(azurermg_resource.test.output).properties.loginServer
}

// it will output "disabled"
output "quarantine_policy" {
  value = jsondecode(azurermg_resource.test.output).properties.policies.quarantinePolicy.status
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
terraform import azurermg_resource.example /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/resGroup1/providers/Microsoft.MachineLearningServices/workspaces/workspace1/computes/cluster1?api-version=2021-07-01
```
