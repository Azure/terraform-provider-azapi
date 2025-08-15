terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
  skip_provider_registration = false
}

data "azapi_client_config" "current" {}

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "applicationDefinition" {
  type      = "Microsoft.Solutions/applicationDefinitions@2021-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-appdef"
  location  = var.location
  body = {
    properties = {
      authorizations = [{
        principalId      = data.azapi_client_config.current.object_id
        roleDefinitionId = "b24988ac-6180-42a0-ab88-20f7382dd24c"
      }]
      createUiDefinition = "    {\n      \"$schema\": \"https://schema.management.azure.com/schemas/0.1.2-preview/CreateUIDefinition.MultiVm.json#\",\n      \"handler\": \"Microsoft.Azure.CreateUIDef\",\n      \"version\": \"0.1.2-preview\",\n      \"parameters\": {\n         \"basics\": [],\n         \"steps\": [],\n         \"outputs\": {}\n      }\n    }\n"
      description        = "Test Managed App Definition"
      displayName        = "TestManagedAppDefinition"
      isEnabled          = true
      lockLevel          = "ReadOnly"
      mainTemplate       = "    {\n      \"$schema\": \"https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#\",\n      \"contentVersion\": \"1.0.0.0\",\n      \"parameters\": {\n\n         \"boolParameter\": {\n            \"type\": \"bool\"\n         },\n         \"intParameter\": {\n            \"type\": \"int\"\n         },\n         \"stringParameter\": {\n            \"type\": \"string\"\n         },\n         \"secureStringParameter\": {\n            \"type\": \"secureString\"\n         },\n         \"objectParameter\": {\n            \"type\": \"object\"\n         },\n         \"arrayParameter\": {\n            \"type\": \"array\"\n         }\n\n      },\n      \"variables\": {},\n      \"resources\": [],\n      \"outputs\": {\n        \"boolOutput\": {\n          \"type\": \"bool\",\n          \"value\": true\n        },\n        \"intOutput\": {\n          \"type\": \"int\",\n          \"value\": 100\n        },\n        \"stringOutput\": {\n          \"type\": \"string\",\n          \"value\": \"stringOutputValue\"\n        },\n        \"objectOutput\": {\n          \"type\": \"object\",\n          \"value\": {\n            \"nested_bool\": true,\n            \"nested_array\": [\"value_1\", \"value_2\"],\n            \"nested_object\": {\n              \"key_0\": 0\n            }\n          }\n        },\n        \"arrayOutput\": {\n          \"type\": \"array\",\n          \"value\": [\"value_1\", \"value_2\"]\n        }\n      }\n    }\n"
    }
  }
}

resource "azapi_resource" "application" {
  type      = "Microsoft.Solutions/applications@2021-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-app"
  location  = var.location
  body = {
    kind = "ServiceCatalog"
    properties = {
      applicationDefinitionId = azapi_resource.applicationDefinition.id
      managedResourceGroupId  = "/subscriptions/${data.azapi_client_config.current.subscription_id}/resourceGroups/${var.resource_name}-infragroup"
      parameters = {
        arrayParameter = {
          value = ["value_1", "value_2"]
        }
        boolParameter = {
          value = true
        }
        intParameter = {
          value = 100
        }
        objectParameter = {
          value = {
            nested_array = ["value_1", "value_2"]
            nested_bool  = true
            nested_object = {
              key_0 = 0
            }
          }
        }
        secureStringParameter = {
          value = ""
        }
        stringParameter = {
          value = "value_1"
        }
      }
    }
  }
}

