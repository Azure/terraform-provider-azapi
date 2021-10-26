---
subcategory: ""
layout: "azurermg"
page_title: "Generic Azure Data Source: azurermg_resource"
description: |-
  Gets information from an existing azure resource
---

# azurermg_resource

This resource can access any existing Azure resource manager resource.

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

resource "azurerm_container_registry" "example" {
  name                = "example"
  resource_group_name = azurerm_resource_group.example.name
  location            = azurerm_resource_group.example.location
  sku                 = "Premium"
  admin_enabled       = false
}

data "azurermg_resource" "test" {
  url         = azurerm_container_registry.example.id
  api_version = "2020-11-01-preview"

  paths = ["properties.loginServer", "properties.policies.quarantinePolicy.status"]
}

// it will output "registry1.azurecr.io"
output "login_server" {
  value = jsondecode(data.azurermg_resource.test.output).properties.loginServer
}

// it will output "disabled"
output "quarantine_policy" {
  value = jsondecode(data.azurermg_resource.test.output).properties.policies.quarantinePolicy.status
}
```

## Arguments Reference

The following arguments are supported:
* `url` - (Required) The url which should be used to access this azure resource. It's also same with the ID of an azure source.
  Here're some examples
  `Container Registry: /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1` and
  `Virtual Machine: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1`.
  Changing this forces a new azure resource to be created.

* `api_version` - (Required) The version of the API used to access this azure resource. Changing this forces a new azure resource to be created.

---

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

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource.

* `identity` - An `identity` block as defined below, which contains the Managed Service Identity information for this azure resource.

* `location` - The Azure Region where the azure resource should exist.

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

* `tags` - A mapping of tags which should be assigned to the azure resource.

---

A `identity` block exports the following:

* `type` - The Type of Identity which should be used for this azure resource. Possible values are `SystemAssigned`, `UserAssigned` and `SystemAssigned,UserAssigned`.

* `identity_ids` - A list of User Managed Identity ID's which should be assigned to the azure resource.

* `principal_id` - The Principal ID for the Service Principal associated with the Managed Service Identity of this azure resource.

* `tenant_id` - The Tenant ID for the Service Principal associated with the Managed Service Identity of this azure resource.

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `read` - (Defaults to 5 minutes) Used when retrieving the azure resource.
