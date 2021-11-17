---
subcategory: ""
layout: "azurerm-restapi"
page_title: "Generic Azure Patch Resource: azurerm-restapi_patch_resource"
description: |-
  Manages a subset of an existing azure resource's properties
---

# azurerm-restapi_resource

This resource can manage a subset of any existing Azure resource manager resource's properties.

## Example Usage

```hcl
terraform {
  required_providers {
    azurerm-restapi = {
      source = "Azure/azurerm-restapi"
    }
  }
}

provider "azurerm-restapi" {
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

resource "azurerm-restapi_patch_resource" "example" {
  resource_id = azurerm_lb.example.id
  type        = "Microsoft.Network/loadBalancers@2021-03-01"
  body        = <<BODY
    {
      "properties": {
        "inboundNatRules": [
          {
            "properties": {
               "idleTimeoutInMinutes": 15
            }
          }
        ]
      }
    }
    BODY

  depends_on = [
    azurerm_lb_nat_rule.example,
  ]
}

```

## Arguments Reference

The following arguments are supported:
* `resource_id` - (Required) The ID of an existing azure source. 
  Here're some examples 
  `Container Registry: /subscriptions/00000000-0000-0000-0000-000000000000/resourcegroups/mygroup1/providers/Microsoft.ContainerRegistry/registries/myregistry1` and 
  `Virtual Machine: /subscriptions/00000000-0000-0000-0000-000000000000/resourceGroups/mygroup1/providers/Microsoft.Compute/virtualMachines/machine1`.
  Changing this forces a new azure resource to be created.

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


## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource.

* `output` - The output json containing the properties specified in `response_export_values`. Here're some examples to decode json and extract the value.
```
// it will output "registry1.azurecr.io"
output "login_server" {
  value = jsondecode(azurerm-restapi_resource.example.output).properties.loginServer
}

// it will output "disabled"
output "quarantine_policy" {
  value = jsondecode(azurerm-restapi_resource.example.output).properties.policies.quarantinePolicy.status
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the azure resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the azure resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the azure resource.
