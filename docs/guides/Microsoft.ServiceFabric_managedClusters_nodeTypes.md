---
subcategory: "Microsoft.ServiceFabric - Service Fabric"
page_title: "managedClusters/nodeTypes"
description: |-
  Manages a Service Fabric Managed Clusters Node Types.
---

# Microsoft.ServiceFabric/managedClusters/nodeTypes - Service Fabric Managed Clusters Node Types

This article demonstrates how to use `azapi` provider to manage the Service Fabric Managed Clusters Node Types resource in Azure.

## Example Usage

### default

```hcl
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

resource "azapi_resource" "managedCluster" {
  type      = "Microsoft.ServiceFabric/managedClusters@2021-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      addonFeatures = [
        "DnsService",
      ]
      adminPassword             = "NotV3ryS3cur3P@$$w0rd"
      adminUserName             = "testUser"
      clientConnectionPort      = 12345
      clusterUpgradeCadence     = "Wave0"
      dnsName                   = var.resource_name
      httpGatewayConnectionPort = 23456
      loadBalancingRules = [
        {
          backendPort      = 8000
          frontendPort     = 443
          probeProtocol    = "http"
          probeRequestPath = "/"
          protocol         = "tcp"
        },
      ]
      networkSecurityRules = [
        {
          access = "allow"
          destinationAddressPrefixes = [
            "0.0.0.0/0",
          ]
          destinationPortRanges = [
            "443",
          ]
          direction = "inbound"
          name      = "rule443-allow-fe"
          priority  = 1000
          protocol  = "tcp"
          sourceAddressPrefixes = [
            "0.0.0.0/0",
          ]
          sourcePortRanges = [
            "1-65535",
          ]
        },
      ]
    }
    sku = {
      name = "Standard"
    }
    tags = {
      Test = "value"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "nodeType" {
  type      = "Microsoft.ServiceFabric/managedClusters/nodeTypes@2021-05-01"
  parent_id = azapi_resource.managedCluster.id
  name      = var.resource_name
  body = {
    properties = {
      applicationPorts = {
        endPort   = 9000
        startPort = 7000
      }
      capacities = {
      }
      dataDiskSizeGB = 130
      dataDiskType   = "Standard_LRS"
      ephemeralPorts = {
        endPort   = 20000
        startPort = 10000
      }
      isPrimary               = true
      isStateless             = false
      multiplePlacementGroups = false
      placementProperties = {
      }
      vmImageOffer     = "WindowsServer"
      vmImagePublisher = "MicrosoftWindowsServer"
      vmImageSku       = "2016-Datacenter"
      vmImageVersion   = "latest"
      vmInstanceCount  = 5
      vmSecrets = [
      ]
      vmSize = "Standard_DS2_v2"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ServiceFabric/managedClusters/nodeTypes@api-version`. The available api-versions for this resource are: [`2020-01-01-preview`, `2021-01-01-preview`, `2021-05-01`, `2021-07-01-preview`, `2021-11-01-preview`, `2022-01-01`, `2022-02-01-preview`, `2022-06-01-preview`, `2022-08-01-preview`, `2022-10-01-preview`, `2023-02-01-preview`, `2023-03-01-preview`, `2023-07-01-preview`, `2023-09-01-preview`, `2023-11-01-preview`, `2023-12-01-preview`, `2024-02-01-preview`, `2024-04-01`, `2024-06-01-preview`, `2024-09-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/managedClusters/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ServiceFabric/managedClusters/nodeTypes?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/managedClusters/{resourceName}/nodeTypes/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ServiceFabric/managedClusters/{resourceName}/nodeTypes/{resourceName}?api-version=2024-09-01-preview
 ```
