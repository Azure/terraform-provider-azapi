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

-> **Note** When delete `azapi_resource_action`, no operation will be performed.

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

Here's an example to use `azapi_resource_action` resource to perform a provider action.

```hcl
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

data "azapi_resource_action" "test" {
  type        = "Microsoft.Resources/providers@2020-04-01-preview"
  resource_id = "/providers/Microsoft.ResourceGraph"
  action      = "resources"
  body = jsonencode({
    query = "resources| where name contains \"test\""
  })
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

* `body` - (Optional) A JSON object that contains the request body.

* `locks` - (Optional) A list of ARM resource IDs which are used to avoid modify azapi resources at the same time.

* `method` - (Optional) Specifies the Http method of the azure resource action. Allowed values are `POST`, `PATCH`, `PUT` and `DELETE`. Defaults to `POST`.

* `response_export_values` - (Optional) A list of path that needs to be exported from response body.
  Setting it to `["*"]` will export the full response body.
  Here's an example. If it sets to `["keys"]`, it will set the following json to computed property `output`.
```
{
  "keys": [
    {
      "KeyName": "Primary",
      "Permissions": "Full",
      "Value": "nHGYNd******i4wdug=="
    },
    {
      "KeyName": "Secondary",
      "Permissions": "Full",
      "Value": "6yoCad******SLzKzg=="
    }
  ]
}
```

## Attributes Reference

In addition to the Arguments listed above - the following Attributes are exported:

* `id` - The ID of the azure resource action.

* `output` - The output json containing the properties specified in `response_export_values`. Here are some examples to decode json and extract the value.

```hcl
// it will output "nHGYNd******i4wdug=="
output "primary_key" {
  value = jsondecode(azapi_resource_action.test.output).keys.0.Value
}

// it will output "6yoCad******SLzKzg=="
output "secondary_key" {
  value = jsondecode(azapi_resource_action.test.output).keys.1.Value
}
```

## Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 30 minutes) Used when creating the azure resource.
* `read` - (Defaults to 5 minutes) Used when retrieving the azure resource.
* `delete` - (Defaults to 30 minutes) Used when deleting the azure resource.
