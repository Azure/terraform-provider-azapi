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
  default = "eastus"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = "${var.resource_name}-rg"
  location = var.location
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2022-07-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-vnet"
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/16",
        ]
      }
      dhcpOptions = {
        dnsServers = [
        ]
      }
      subnets = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
  lifecycle {
    ignore_changes = [body.properties.subnets]
  }
}

resource "azapi_resource" "subnet" {
  type      = "Microsoft.Network/virtualNetworks/subnets@2022-07-01"
  parent_id = azapi_resource.virtualNetwork.id
  name      = "${var.resource_name}-subnet"
  body = {
    properties = {
      addressPrefix = "10.0.2.0/24"
      delegations = [
      ]
      privateEndpointNetworkPolicies    = "Enabled"
      privateLinkServiceNetworkPolicies = "Enabled"
      serviceEndpointPolicies = [
      ]
      serviceEndpoints = [
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "containerGroupProfile" {
  type      = "Microsoft.ContainerInstance/containerGroupProfiles@2024-05-01-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-contianerGroup"
  location  = var.location

  body = {
    properties = {
      "containers" : [
        {
          "name" : "mycontainergroupprofile",
          "properties" : {
            "command" : [],
            "environmentVariables" : [],
            "image" : "mcr.microsoft.com/azuredocs/aci-helloworld:latest",
            "ports" : [
              {
                "port" : 8000
              }
            ],
            "resources" : {
              "requests" : {
                "cpu" : 1,
                "memoryInGB" : 1.5
              }
            }
          }
        }
      ],
      "imageRegistryCredentials" : [],
      "ipAddress" : {
        "ports" : [
          {
            "protocol" : "TCP",
            "port" : 8000
          }
        ],
        "type" : "Public"
      },
      "osType" : "Linux",
      "sku" : "Standard"
    }
  }

  schema_validation_enabled = false
}



resource "azapi_resource" "standbyContainerGroupPool" {
  type      = "Microsoft.StandbyPool/standbyContainerGroupPools@2025-03-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = "${var.resource_name}-CGPool"
  body = {
    properties = {
      containerGroupProperties = {
        containerGroupProfile = {
          id       = azapi_resource.containerGroupProfile.id
          revision = 1
        }
        subnetIds = [
          {
            id = azapi_resource.subnet.id
          },
        ]
      }
      elasticityProfile = {
        maxReadyCapacity = 5
        refillPolicy     = "always"
      }
      zones = [
        "1",
        "2",
        "3",
      ]
    },
    location = "${var.location}"
  }
  schema_validation_enabled = false
  ignore_casing             = false
  ignore_missing_property   = false
}
