# Terraform Provider for Azure Resource Manager Rest API

## How to use?
This project is used in API test only.

This provider hasn't been formly released, it can only be used by configuring provider override.
```
1. git clone https://github.com/Azure/terraform-provider-azurerm-restapi.git
2. cd terraform-provider-azurerm-restapi
3. go install
4. edit terraform.rc and add the following configuration, refs: https://www.terraform.io/docs/cli/config/config-file.html
  
  dev_overrides {
    "Azure/azurerm-restapi" = "C:\\Users\\henglu\\go\\bin" #path to provider execute
  }
```

## Usage Example

The following example shows how to use `azurerm-restapi_resource` to manage machine learning compute resource.

```hcl
terraform {
  required_providers {
    azurerm-restapi = {
      source  = "Azure/azurerm-restapi"
    }
  }
}

provider "azurerm-restapi" {
  # More information on the authentication methods supported by
  # the AzureRM Provider can be found here:
  # https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs

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

resource "azurerm-restapi_resource" "example" {
  resource_id = "${data.azurerm_machine_learning_workspace.existing.id}/computes/example"
  type = "Microsoft.MachineLearningServices/workspaces/computes@2021-07-01"
  body = <<BODY
    {
      "location": "eastus",
      "properties": {
        "computeType": "ComputeInstance",
        "disableLocalAuth": true,
        "properties": {
          "vmSize": "STANDARD_NC6"
        }
      }
    }
  BODY
}

```

