# Terraform Provider for Azure Resource Manager Rest API

## How to use?

This provider hasn't been formly released, it can only be used by configuring provider override.
```
1. git clone https://github.com/Azure/terraform-provider-azapi.git
2. cd terraform-provider-azapi
3. go install
4. edit terraform.rc and add the following configuration, refs: https://www.terraform.io/docs/cli/config/config-file.html
  
  dev_overrides {
    "Azure/azapi" = "C:\\Users\\henglu\\go\\bin" #path to provider execute
  }
```

## Usage Example

The following example shows how to use `azapi_resource` to manage machine learning compute resource.

```hcl
terraform {
  required_providers {
    azapi = {
      source  = "Azure/azapi"
    }
  }
}

provider "azapi" {
  # More information on the authentication methods supported by
  # the AzApi Provider can be found here:
  # https://registry.terraform.io/providers/Azure/azapi/latest/docs

  # subscription_id = "..."
  # client_id       = "..."
  # client_secret   = "..."
  # tenant_id       = "..."
}

provider "azurerm" {
  features {}
}

data "azurerm_machine_learning_workspace" "existing" {
  name                = "example-workspace"
  resource_group_name = "example-resources"
}

resource "azapi_resource" "example" {
  name = "example"
  parent_id = data.azurerm_machine_learning_workspace.existing.id
  type = "Microsoft.MachineLearningServices/workspaces/computes@2021-07-01"
  
  location = "eastus"
  body = jsondecode({
    properties = {
      computeType      = "ComputeInstance"
      disableLocalAuth = true
      properties = {
        vmSize = "STANDARD_NC6"
      }
    }
  })
}

```

