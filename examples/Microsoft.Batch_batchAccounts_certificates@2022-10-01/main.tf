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

variable "certificate_data" {
  type        = string
  description = "The base64-encoded certificate data"
  sensitive   = true
}

variable "certificate_thumbprint" {
  type        = string
  description = "The thumbprint of the certificate"
}

resource "azapi_resource" "resourceGroup" {
  type     = "Microsoft.Resources/resourceGroups@2020-06-01"
  name     = var.resource_name
  location = var.location
}

resource "azapi_resource" "batchAccount" {
  type      = "Microsoft.Batch/batchAccounts@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      encryption = {
        keySource = "Microsoft.Batch"
      }
      poolAllocationMode  = "BatchService"
      publicNetworkAccess = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "certificate" {
  type      = "Microsoft.Batch/batchAccounts/certificates@2022-10-01"
  parent_id = azapi_resource.batchAccount.id
  name      = "SHA1-${var.certificate_thumbprint}"
  body = {
    properties = {
      data                = var.certificate_data
      format              = "Cer"
      thumbprint          = var.certificate_thumbprint
      thumbprintAlgorithm = "sha1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

