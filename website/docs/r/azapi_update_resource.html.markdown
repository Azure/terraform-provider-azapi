---
subcategory: ""
layout: "azapi"
page_title: "Generic Azure Update Resource: azapi_update_resource"
description: |-
  Manages a subset of an existing azure resource's properties
---

# azapi_update_resource

This resource can manage a subset of any existing Azure resource manager resource's properties.

-> **Note**: This resource is used to add or modify properties on an existing resource.
When delete `azapi_update_resource`, no operation will be performed, and these properties will stay unchanged.
If you want to restore the modified properties to some values, you must apply the restored properties before deleting.

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

resource "azurerm_public_ip" "example" {
  name                = "example-ip"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name
  allocation_method   = "Static"
}

resource "azurerm_lb" "example" {
  name                = "example-lb"
  location            = azurerm_resource_group.example.location
  resource_group_name = azurerm_resource_group.example.name

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.example.id
  }
}

resource "azurerm_lb_nat_rule" "example" {
  resource_group_name            = azurerm_resource_group.example.name
  loadbalancer_id                = azurerm_lb.example.id
  name                           = "RDPAccess"
  protocol                       = "Tcp"
  frontend_port                  = 3389
  backend_port                   = 3389
  frontend_ip_configuration_name = "PublicIPAddress"
}

resource "azapi_update_resource" "example" {
  type        = "Microsoft.Network/loadBalancers@2021-03-01"
  resource_id = azurerm_lb.example.id

  body = jsonencode({
    properties = {
      inboundNatRules = [
        {
          properties = {
            idleTimeoutInMinutes = 15
          }
        }
      ]
    }
  })

  depends_on = [
    azurerm_lb_nat_rule.example,
  ]
}

```

## Arguments Reference

The following arguments are supported:
* `name` - (Optional) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `parent_id` - (Optional) The ID of the azure resource in which this resource is created. Changing this forces a new resource to be created.

* `resource_id` - (Optional) The ID of an existing azure source. 
  Here're some examples 
  `Container Registry: /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1` and 
  `Virtual Machine: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1`.
  Changing this forces a new azure resource to be created.

~> **Note:** Configuring `name` and `parent_id` is an alternative way to configure `resource_id`.

* `type` - (Required) It is in a format like `<resource-type>@<api-version>`. `<resource-type>` is the Azure resource type, for example, `Microsoft.Storage/storageAccounts`.
  `<api-version>` is version of the API used to manage this azure resource.

* `body` - (Required) A JSON object that contains the request body used to add on an existing azure resource. 

---

* `response_export_values` - (Optional) A list of path that needs to be exported from response body. Here's an example. 
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

* `ignore_casing` - (Optional) Whether ignore incorrect casing returned in `body` to suppress plan-diff. Defaults to `false`.

* `ignore_missing_property` - (Optional) Whether ignore not returned properties like credentials in `body` to suppress plan-diff. Defaults to `false`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource.

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

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the azure resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the azure resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the azure resource.
