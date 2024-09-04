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

resource "azapi_resource" "serverfarm" {
  type      = "Microsoft.Web/serverfarms@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "Windows"
    properties = {
      isXenon = false
    }
    sku = {
      capacity = 1
      name     = "S1"
      size     = "S1"
      tier     = "Standard"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "site" {
  type      = "Microsoft.Web/sites@2021-02-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      clientAffinityEnabled = false
      clientCertEnabled     = false
      enabled               = true
      httpsOnly             = false
      serverFarmId          = azapi_resource.serverfarm.id
      siteConfig = {
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "publicCertificate" {
  type      = "Microsoft.Web/sites/publicCertificates@2022-09-01"
  parent_id = azapi_resource.site.id
  name      = var.resource_name
  body = {
    properties = {
      blob                      = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNzakNDQVpvQ0NRQ01kdDdEdnlnUHREQU5CZ2txaGtpRzl3MEJBUXNGQURBYk1Sa3dGd1lEVlFRRERCQmgKY0drdWRHVnljbUZtYjNKdExtbHZNQjRYRFRFNE1EY3dOVEV3TXpNek1Gb1hEVEk0TURjd01qRXdNek16TUZvdwpHekVaTUJjR0ExVUVBd3dRWVhCcExuUmxjbkpoWm05eWJTNXBiekNDQVNJd0RRWUpLb1pJaHZjTkFRRUJCUUFECmdnRVBBRENDQVFvQ2dnRUJBS1FXMzMyT2wyOENzaWRBaGVEMWFMOVVsOEpXbktMZGFWeEtaM3NzbDVDWGpQRE8KbU03SVhrMFNnYlFuVUM4bElsUEZaaURHYlExc0I2T1RNdW42Wlo0aXBMcDgwZHRsMHJvQ0x0Q25EUU9CR3pDTgpBckNZQW9YUnVyamtYRVk3dHBEMHd3dFU3MiszN2gzSFE0ZzBWUzZWSXRKQ3FKOVFBRFYrSE8yWld1WlRlejcwCk1ob0w2T0xmWlA3SEdZZEpES2dmRVZORjVYbGJWek5BR2tESUpGZGhqTnh5R0d1NU5mc20xcGZRaEF5dW5razcKSlZhbWpVZzVJb2pSZG82M0lTOXd3ek1PZGVHU0FiQmNzSmZZZUNmVmcya3VwUjhxMFRtWit4OTNSbW1PbGJTaQo2NmtFWXhSelo5WUNRZUhKbW4xWWZKOTJCcENVaXk5QTZaMWlhS1VDQXdFQUFUQU5CZ2txaGtpRzl3MEJBUXNGCkFBT0NBUUVBSjdKaGxlY1A3SjQ4d0kyUUhUTWJBTWtrV0J2L2lXcTEvUUlGNHVnSDNaYjVQb3JPditOZmhRMEwKbFdpdy9Tek44QWU5NXZVaXhBR1lITVNhMjhvdW1NNUsxT3NxS0VrVklvMUFvQkg4bkJ6K1ZjVHBSRC9tSFhvdApBSFBBWnQ5ajVMcWVIWCtlblI2UmJJTkFmM2puK1lVM01kVmUwTXNBRGRGQVNWRGZqbVFQMlI3bzlhSmIvUXFPCmczYlpCV3NpQkRFSVNmeWFIMitwZ1VNN3d0d0VvRldtRU1sZ2pMSzFNUkJzMWNEWlhxbkhhQ2QvcnMrTm1XVjkKbmFFdTd4NWZ5UU9rNEhvemtwd2VSK0p4MXNCbFRSc2E0OS9xU0h0LzZVTEtmTzAxL2NUczRpRjcxeWtYUGJoMwpLajljSTJ1bzlhWXRYa3hraEtyR3lVcEE3RkpxV3c9PQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCg=="
      publicCertificateLocation = "Unknown"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

