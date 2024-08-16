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

resource "azapi_resource" "deployment" {
  type      = "Microsoft.Resources/deployments@2020-06-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  body = {
    properties = {
      mode = "Complete"
      template = {
        "$schema"      = "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#"
        contentVersion = "1.0.0.0"
        parameters = {
          storageAccountType = {
            allowedValues = [
              "Standard_LRS",
              "Standard_GRS",
              "Standard_ZRS",
            ]
            defaultValue = "Standard_LRS"
            metadata = {
              description = "Storage Account type"
            }
            type = "string"
          }
        }
        resources = [
          {
            apiVersion = "[variables('apiVersion')]"
            location   = "[variables('location')]"
            name       = "[variables('storageAccountName')]"
            properties = {
              accountType = "[parameters('storageAccountType')]"
            }
            type = "Microsoft.Storage/storageAccounts"
          },
          {
            apiVersion = "[variables('apiVersion')]"
            location   = "[variables('location')]"
            name       = "[variables('publicIPAddressName')]"
            properties = {
              dnsSettings = {
                domainNameLabel = "[variables('dnsLabelPrefix')]"
              }
              publicIPAllocationMethod = "[variables('publicIPAddressType')]"
            }
            type = "Microsoft.Network/publicIPAddresses"
          },
        ]
        variables = {
          apiVersion          = "2015-06-15"
          dnsLabelPrefix      = "[concat('terraform-tdacctest', uniquestring(resourceGroup().id))]"
          location            = "[resourceGroup().location]"
          publicIPAddressName = "[concat('myPublicIp', uniquestring(resourceGroup().id))]"
          publicIPAddressType = "Dynamic"
          storageAccountName  = "[concat(uniquestring(resourceGroup().id), 'storage')]"
        }
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

