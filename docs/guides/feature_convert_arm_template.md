---
layout: "azapi"
page_title: "Feature: Convert ARM Template/Resource JSON to AzAPI Configuration"
description: |-
  This guide will cover how to use AzAPI VSCode extension to convert ARM template/resource JSON to AzAPI configuration.

---

The AzAPI VSCode extension provides a feature to convert ARM template/resource JSON to AzAPI configuration. This feature is useful when you have an existing ARM template or resource JSON and want to convert it to AzAPI configuration.

This article covers how to use the feature to convert ARM template/resource JSON to AzAPI configuration.

## Prerequisites

- [AzAPI VSCode extension](https://marketplace.visualstudio.com/items?itemName=azapi-vscode.azapi) version 2.0.1 or later
- [Terraform AzAPI provider](https://registry.terraform.io/providers/azure/azapi) version 2.0.1 or later

## Convert ARM template to AzAPI configuration

This feature only supports ARM template exported from the Azure portal. Other ARM templates may not be supported.

1. Export the ARM template from the Azure portal.  
A step-by-step guide on how to export the ARM template from the Azure portal can be found [here](https://learn.microsoft.com/en-us/azure/azure-resource-manager/templates/export-template-portal).

2. Copy the entire ARM template JSON. For example, the following is an ARM template JSON:

```json
{
    "$schema": "https://schema.management.azure.com/schemas/2019-04-01/deploymentTemplate.json#",
    "contentVersion": "1.0.0.0",
    "parameters": {
        "Clusters_tfmonitoring_name": {
            "defaultValue": "tfmonitoring",
            "type": "String"
        },
        "accounts_tfvoice_name": {
            "defaultValue": "tfvoice",
            "type": "String"
        },
        "userAssignedIdentities_monitoringtoolsuai_name": {
            "defaultValue": "monitoringtoolsuai",
            "type": "String"
        }
    },
    "variables": {},
    "resources": [
        {
            "type": "Microsoft.CognitiveServices/accounts",
            "apiVersion": "2024-06-01-preview",
            "name": "[parameters('accounts_tfvoice_name')]",
            "location": "eastus",
            "sku": {
                "name": "F0"
            },
            "kind": "SpeechServices",
            "identity": {
                "type": "None"
            },
            "properties": {
                "networkAcls": {
                    "defaultAction": "Allow",
                    "virtualNetworkRules": [],
                    "ipRules": []
                },
                "publicNetworkAccess": "Enabled",
                "disableLocalAuth": true
            }
        },
        {
            "type": "Microsoft.Kusto/Clusters",
            "apiVersion": "2023-08-15",
            "name": "[parameters('Clusters_tfmonitoring_name')]",
            "location": "East US",
            "tags": {
                "NRMS.KustoPlatform.Classification.1P": "Corp"
            },
            "sku": {
                "name": "Dev(No SLA)_Standard_E2a_v4",
                "tier": "Basic",
                "capacity": 1
            },
            "zones": [
                "1",
                "3",
                "2"
            ],
            "identity": {
                "type": "SystemAssigned"
            },
            "properties": {
                "trustedExternalTenants": [],
                "enableDiskEncryption": false,
                "enableStreamingIngest": false,
                "languageExtensions": {
                    "value": []
                },
                "enablePurge": false,
                "enableDoubleEncryption": false,
                "engineType": "V3",
                "acceptedAudiences": [],
                "restrictOutboundNetworkAccess": "Disabled",
                "allowedFqdnList": [],
                "publicNetworkAccess": "Enabled",
                "allowedIpRangeList": [],
                "enableAutoStop": true,
                "publicIPType": "IPv4"
            }
        }
    ]
}
```

3. Paste the ARM template JSON in a `*.tf` file in the VSCode editor.  
The AzAPI VSCode extension will automatically detect the ARM template JSON and convert it to AzAPI configuration.

```hcl
variable "subscriptionId" {
  type        = string
  description = "The subscription id"
}

variable "resourceGroupName" {
  type        = string
  description = "The resource group name"
}

variable "Clusters_tfmonitoring_name" {
  type    = string
  default = "tfmonitoring"
}

variable "accounts_tfvoice_name" {
  type    = string
  default = "tfvoice"
}

variable "userAssignedIdentities_monitoringtoolsuai_name" {
  type    = string
  default = "monitoringtoolsuai"
}

resource "azapi_resource" "account" {
  type      = "Microsoft.CognitiveServices/accounts@2024-10-01"
  parent_id = "/subscriptions/${var.subscriptionId}/resourceGroups/${var.resourceGroupName}"
  name      = "${var.accounts_tfvoice_name}"
  location  = "eastus"
  body = {
    kind = "SpeechServices"
    properties = {
      disableLocalAuth = true
      networkAcls = {
        defaultAction       = "Allow"
        ipRules             = []
        virtualNetworkRules = []
      }
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "F0"
    }
  }
}

resource "azapi_resource" "Cluster" {
  type      = "Microsoft.Kusto/Clusters@2024-04-13"
  parent_id = "/subscriptions/${var.subscriptionId}/resourceGroups/${var.resourceGroupName}"
  name      = "${var.Clusters_tfmonitoring_name}"
  location  = "East US"
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      acceptedAudiences      = []
      allowedFqdnList        = []
      allowedIpRangeList     = []
      enableAutoStop         = true
      enableDiskEncryption   = false
      enableDoubleEncryption = false
      enablePurge            = false
      enableStreamingIngest  = false
      engineType             = "V3"
      languageExtensions = {
        value = []
      }
      publicIPType                  = "IPv4"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      trustedExternalTenants        = []
    }
    sku = {
      capacity = 1
      name     = "Dev(No SLA)_Standard_E2a_v4"
      tier     = "Basic"
    }
    zones = ["1", "3", "2"]
  }
  tags = {
    "NRMS.KustoPlatform.Classification.1P" = "Corp"
  }
}
```

## Convert resource JSON to AzAPI configuration

This feature supports converting resource JSON to AzAPI configuration.

There are different ways to get the resource JSON:

1. From the Azure Portal:
   1. Navigate to the target resource in the Azure portal.
   2. On the `Overview` page, in the `Essentials` section, select the `JSON View` link.
   3. The JSON view will display the resource JSON. And it's recommended to select the latest API version in the JSON view.  


2. From the Azure CLI:  
  More details could be found [here](https://docs.microsoft.com/en-us/cli/azure/resource?view=azure-cli-latest#az-resource-show).  
  For example, to get the resource JSON of a Cognitive Services account, you can run the following command:

```bash
az resource show --ids "/subscriptions/<subscription-id>/resourceGroups/<resource-group-name>/providers/Microsoft.CognitiveServices/accounts/<account-name>"
``` 

For example, the following is a resource JSON:

```json
{
    "apiVersion": "2024-10-01",
    "id": "/subscriptions/000000/resourceGroups/example-rg/providers/Microsoft.CognitiveServices/accounts/terraform",
    "name": "terraform",
    "type": "microsoft.cognitiveservices/accounts",
    "sku": {
        "name": "S0"
    },
    "kind": "OpenAI",
    "location": "eastus",
    "tags": {},
    "properties": {
        "endpoint": "https://terraform.openai.azure.com/",
        "provisioningState": "Succeeded",
        "dateCreated": "2023-03-29T06:22:27.5216767Z",
        "isMigrated": false,
        "customSubDomainName": "terraform",
        "privateEndpointConnections": [],
        "publicNetworkAccess": "Enabled",
        "disableLocalAuth": true,
        "endpoints": {
            "OpenAI Language Model Instance API": "https://terraform.openai.azure.com/",
            "OpenAI Dall-E API": "https://terraform.openai.azure.com/",
            "OpenAI Whisper API": "https://terraform.openai.azure.com/",
            "OpenAI Model Scaleset API": "https://terraform.openai.azure.com/",
            "OpenAI Realtime API": "https://terraform.openai.azure.com/",
            "Token Service API": "https://terraform.openai.azure.com/"
        }
    },
    "systemData": {
    },
    "etag": "\"4207bbe6-0000-0100-0000-674b228f0000\""
}
```

Copy the resource JSON and paste it in a `*.tf` file in the VSCode editor. The AzAPI VSCode extension will automatically detect the resource JSON and convert it to AzAPI configuration, it will remove the read-only fields and only keep the fields that can be modified.

Here is the converted AzAPI configuration:
```hcl
/*
Note: This is a generated HCL content from the JSON input which is based on the latest API version available.
To import the resource, please run the following command:
terraform import azapi_resource.account /subscriptions/000000/resourceGroups/example-rg/providers/Microsoft.CognitiveServices/accounts/terraform?api-version=2024-10-01

Or add the below config:
import {
  id = "/subscriptions/000000/resourceGroups/example-rg/providers/Microsoft.CognitiveServices/accounts/terraform?api-version=2024-10-01"
  to = azapi_resource.account
}
*/

resource "azapi_resource" "account" {
  type      = "Microsoft.CognitiveServices/accounts@2024-10-01"
  parent_id = "/subscriptions/000000/resourceGroups/example-rg"
  name      = "terraform"
  location  = "eastus"
  body = {
    kind = "OpenAI"
    properties = {
      customSubDomainName = "terraform"
      disableLocalAuth    = true
      publicNetworkAccess = "Enabled"
    }
    sku = {
      name = "S0"
    }
  }
}

```

To import the resource, you can run the following command:
```bash
terraform import azapi_resource.account /subscriptions/000000/resourceGroups/example-rg/providers/Microsoft.CognitiveServices/accounts/terraform?api-version=2024-10-01
```

Or add the below config:
```hcl
import {
  id = "/subscriptions/000000/resourceGroups/example-rg/providers/Microsoft.CognitiveServices/accounts/terraform?api-version=2024-10-01"
  to = azapi_resource.account
}
```
