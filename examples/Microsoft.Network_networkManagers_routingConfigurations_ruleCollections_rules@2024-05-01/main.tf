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

variable "deploy_locations" {
  type    = list(string)
  default = ["westeurope"]
}

data "azapi_client_config" "current" {
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "networkManager" {
  type      = "Microsoft.Network/networkManagers@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      description = ""
      networkManagerScopeAccesses = [
        "Routing",
      ]
      networkManagerScopes = {
        managementGroups = [
        ]
        subscriptions = [
          "/subscriptions/${data.azapi_client_config.current.subscription_id}",
        ]
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "networkGroup" {
  type      = "Microsoft.Network/networkManagers/networkGroups@2024-05-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  body = {
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "virtualNetwork" {
  type      = "Microsoft.Network/virtualNetworks@2024-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addressSpace = {
        addressPrefixes = [
          "10.0.0.0/22",
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

resource "azapi_resource" "staticMember" {
  type      = "Microsoft.Network/networkManagers/networkGroups/staticMembers@2024-05-01"
  parent_id = azapi_resource.networkGroup.id
  name      = var.resource_name
  body = {
    properties = {
      resourceId = azapi_resource.virtualNetwork.id
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "routingConfiguration" {
  type      = "Microsoft.Network/networkManagers/routingConfigurations@2024-05-01"
  parent_id = azapi_resource.networkManager.id
  name      = var.resource_name
  body = {
    properties = {
      description = "example routing configuration"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "ruleCollection" {
  type      = "Microsoft.Network/networkManagers/routingConfigurations/ruleCollections@2024-05-01"
  parent_id = azapi_resource.routingConfiguration.id
  name      = var.resource_name
  body = {
    properties = {
      description = "example rule collection"
      appliesTo = [
        {
          networkGroupId = azapi_resource.networkGroup.id
        },
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "rule" {
  type      = "Microsoft.Network/networkManagers/routingConfigurations/ruleCollections/rules@2024-05-01"
  parent_id = azapi_resource.ruleCollection.id
  name      = var.resource_name
  body = {
    properties = {
      description = "example rule"
      destination = {
        type               = "AddressPrefix" # Required, possible values are "AddressPrefix", "ServiceTag"
        destinationAddress = "10.0.0.0/16"   # Required
      }
      nextHop = {
        nextHopType    = "VirtualNetworkGateway" # Required, possible values are "Internet", "NoNextHop", "VirtualAppliance", "VirtualNetworkGateway", "VnetLocal"
        nextHopAddress = ""                      # Required if nextHopType is "VirtualAppliance"
      }
    }
  }
  retry = {
    error_message_regex = [
      "Deleting failed(.+)resource has been deployed" # This retries deletion in case a rule might fail to delete immediately after undeploy during `terraform destroy`.
    ]
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource_action" "deploy" {
  type        = "Microsoft.Network/networkManagers@2024-05-01"
  action      = "commit"
  when        = "apply"
  resource_id = azapi_resource.networkManager.id
  body = {
    configurationIds = [azapi_resource.routingConfiguration.id] # Optional, to remove all configurations from the specified regions, leave the field as empty array
    targetLocations  = var.deploy_locations                     # Required
    commitType       = "Routing"                                # Required, possible values are "SecurityAdmin", "Connectivity", "Routing"
  }
  response_export_values = ["*"]
  depends_on             = [azapi_resource.rule]
  lifecycle {
    replace_triggered_by = [azapi_resource.rule.body.destination, azapi_resource.rule.body.nextHop] # trigger a new deployment when the rule is changed
  }
}

# this one is to remove the deployment when `terraform destroy` is called
resource "azapi_resource_action" "undeploy" {
  type        = "Microsoft.Network/networkManagers@2024-05-01"
  action      = "commit"
  when        = "destroy"
  resource_id = azapi_resource.networkManager.id
  body = {
    configurationIds = []                   # Optional, to remove all configurations from the specified regions, leave the field as empty array
    targetLocations  = var.deploy_locations # Required
    commitType       = "Routing"            # Required, possible values are "SecurityAdmin", "Connectivity", "Routing"
  }
  response_export_values = ["*"]
  depends_on             = [azapi_resource.rule]
}
