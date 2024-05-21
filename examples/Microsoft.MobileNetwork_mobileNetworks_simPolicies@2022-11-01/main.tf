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
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "mobileNetwork" {
  type      = "Microsoft.MobileNetwork/mobileNetworks@2022-11-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      publicLandMobileNetworkIdentifier = {
        mcc = "001"
        mnc = "01"
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "service" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/services@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      pccRules = [
        {
          ruleName       = "default-rule"
          rulePrecedence = 1
          serviceDataFlowTemplates = [
            {
              direction = "Uplink"
              ports = [
              ]
              protocol = [
                "ip",
              ]
              remoteIpList = [
                "10.3.4.0/24",
              ]
              templateName = "IP-to-server"
            },
          ]
          trafficControl = "Enabled"
        },
      ]
      servicePrecedence = 0
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "dataNetwork" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/dataNetworks@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "slice" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/slices@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      snssai = {
        sst = 1
      }
    }

  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "simPolicy" {
  type      = "Microsoft.MobileNetwork/mobileNetworks/simPolicies@2022-11-01"
  parent_id = azapi_resource.mobileNetwork.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      defaultSlice = {
        id = azapi_resource.slice.id
      }
      registrationTimer = 3240
      sliceConfigurations = [
        {
          dataNetworkConfigurations = [
            {
              "5qi"                               = 9
              additionalAllowedSessionTypes       = null
              allocationAndRetentionPriorityLevel = 9
              allowedServices = [
                {
                  id = azapi_resource.service.id
                },
              ]
              dataNetwork = {
                id = azapi_resource.dataNetwork.id
              }
              defaultSessionType             = "IPv4"
              maximumNumberOfBufferedPackets = 10
              preemptionCapability           = "NotPreempt"
              preemptionVulnerability        = "Preemptable"
              sessionAmbr = {
                downlink = "1 Gbps"
                uplink   = "500 Mbps"
              }
            },
          ]
          defaultDataNetwork = {
            id = azapi_resource.dataNetwork.id
          }
          slice = {
            id = azapi_resource.slice.id
          }
        },
      ]
      ueAmbr = {
        downlink = "1 Gbps"
        uplink   = "500 Mbps"
      }
    }
    tags = {
      key = "value"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

