---
subcategory: ""
layout: "azapi"
page_title: "Azure Resource Action: azapi_resource_action"
description: |-
  Perform resource action which changes an existing resource's state
---

# azapi_resource_action

This resource can perform any Azure resource manager resource action.
It's recommended to use `azapi_resource_action` resource to perform actions which change a resource's state, please use `azapi_resource_action` data source,
if user wants to perform readonly action.

-> **Note** The action can be performed on either apply or destroy. The default is apply, see `when` argument for more details.

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

variable "enabled" {
  type        = bool
  default     = false
  description = "whether start the spring service"
}

resource "azurerm_resource_group" "example" {
  name     = "example-rg"
  location = "west europe"
}

resource "azurerm_spring_cloud_service" "test" {
  name                = "example-spring"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  sku_name            = "S0"
}

resource "azapi_resource_action" "start" {
  type                   = "Microsoft.AppPlatform/Spring@2022-05-01-preview"
  resource_id            = azurerm_spring_cloud_service.test.id
  action                 = "start"
  response_export_values = ["*"]

  count = var.enabled ? 1 : 0
}

resource "azapi_resource_action" "stop" {
  type                   = "Microsoft.AppPlatform/Spring@2022-05-01-preview"
  resource_id            = azurerm_spring_cloud_service.test.id
  action                 = "stop"
  response_export_values = ["*"]

  count = var.enabled ? 0 : 1
}
```

Here's an example to use the `azapi_resource_action` resource to register a provider.

```hcl
provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

resource "azapi_resource_action" "test" {
  type        = "Microsoft.Resources/providers@2021-04-01"
  resource_id = "/subscriptions/${data.azurerm_client_config.current.subscription_id}/providers/Microsoft.Compute"
  action      = "register"
  method      = "POST"
}
```

Here's an example to use the `azapi_resource_action` resource to perform a provider action.

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

data "azapi_resource_action" "test" {
  type        = "Microsoft.ResourceGraph@2020-04-01-preview"
  resource_id = "/providers/Microsoft.ResourceGraph"
  action      = "resources"
  body = {
    query = "resources| where name contains \"test\""
  }
  response_export_values = ["*"]
}
```

## Arguments Reference

The following arguments are supported:

* `type` - (Required) It is in a format like `<resource-type>@<api-version>`. `<resource-type>` is the Azure resource type, for example, `Microsoft.Storage/storageAccounts`.
  `<api-version>` is version of the API used to manage this azure resource.

* `resource_id` - (Required) The ID of an existing azure source.

* `action` - (Optional) The name of the resource action. It's also possible to make Http requests towards the resource ID if leave this field empty.

---

* `body` - (Optional) A dynamic attribute that contains the request body.

* `locks` - (Optional) A list of ARM resource IDs which are used to avoid modify azapi resources at the same time.

* `method` - (Optional) Specifies the Http method of the azure resource action. Allowed values are `POST`, `PATCH`, `PUT` and `DELETE`. Defaults to `POST`.

* `response_export_values` - (Optional) A list of path that needs to be exported from response body.
  Setting it to `["*"]` will export the full response body.
  Here's an example. If it sets to `["keys"]`, it will set the following HCL object to computed property `output`.

```
{
  keys = [
    {
      KeyName = "Primary"
      Permissions = "Full"
      Value = "nHGYNd******i4wdug=="
    },
    {
      KeyName = "Secondary"
      Permissions = "Full"
      Value = "6yoCad******SLzKzg=="
    }
  ]
}
```

* `when` - (Optional) When to perform the action, value must be one of: `apply`, `destroy`. Default is `apply`.

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource action.

* `output` - The HCL object containing the properties specified in `response_export_values`. Here are some examples to use the values.

```hcl
// it will output "nHGYNd******i4wdug=="
output "primary_key" {
  value = azapi_resource_action.test.output.keys.0.Value
}

// it will output "6yoCad******SLzKzg=="
output "secondary_key" {
  value = azapi_resource_action.test.output.keys.1.Value
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the azure resource action.
* `update` - (Defaults to 30 minutes) Used when updating the azure resource action.
* `read` - (Defaults to 5 minutes) Used when retrieving the azure resource action.
* `delete` - (Defaults to 30 minutes) Used when deleting the azure resource action.
