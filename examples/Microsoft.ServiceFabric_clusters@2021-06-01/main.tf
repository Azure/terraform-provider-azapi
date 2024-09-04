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

variable "resource_name" {
  type    = string
  default = "acctest0001"
}

variable "location" {
  type    = string
  default = "westeurope"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "cluster" {
  type      = "Microsoft.ServiceFabric/clusters@2021-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addOnFeatures = [
      ]
      fabricSettings = [
      ]
      managementEndpoint = "http://example:80"
      nodeTypes = [
        {
          capacities = {
          }
          clientConnectionEndpointPort = 2020
          durabilityLevel              = "Bronze"
          httpGatewayEndpointPort      = 80
          isPrimary                    = true
          isStateless                  = false
          multipleAvailabilityZones    = false
          name                         = "first"
          placementProperties = {
          }
          vmInstanceCount = 3
        },
      ]
      reliabilityLevel = "Bronze"
      upgradeMode      = "Automatic"
      vmImage          = "Windows"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

