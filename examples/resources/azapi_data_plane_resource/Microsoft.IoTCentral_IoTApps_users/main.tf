terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2021-04-01"
  name     = "example-resources"
  location = "westeurope"
}

resource "azapi_resource" "iotApp" {
  type      = "Microsoft.IoTCentral/iotApps@2021-11-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "exampleiotapp"
  location  = azapi_resource.resourceGroup.location
  body = {
    sku = {
      name = "ST2"
    }
    properties = {
      displayName = "Example IoT App"
      subdomain   = "exampleiotapp"
    }
  }
  response_export_values = {
    subdomain = "properties.subdomain"
  }
}

resource "azapi_data_plane_resource" "example" {
  type      = "Microsoft.IoTCentral/IoTApps/users@2022-07-31"
  parent_id = "${azapi_resource.iotApp.output.subdomain}.azureiotcentral.com"
  name      = "exampleuser"
  body = {
    type = "email"
    roles = [
      {
        role = "ae2c9854-393b-4f97-8c42-479d70ce626e"
      }
    ]
    email = "user@example.com"
  }
}
