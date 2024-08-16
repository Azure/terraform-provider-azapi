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

resource "azapi_resource" "component" {
  type      = "Microsoft.Insights/components@2020-02-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "web"
    properties = {
      Application_Type                = "web"
      DisableIpMasking                = false
      DisableLocalAuth                = false
      ForceCustomerStorageForProfiler = false
      RetentionInDays                 = 90
      SamplingPercentage              = 100
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "webTest" {
  type      = "Microsoft.Insights/webTests@2022-06-15"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "standard"
    properties = {
      Description = ""
      Enabled     = false
      Frequency   = 300
      Kind        = "standard"
      Locations = [
        {
          Id = "us-tx-sn1-azr"
        },
      ]
      Name = var.resource_name
      Request = {
        FollowRedirects = false
        Headers = [
          {
            key   = "x-header"
            value = "testheader"
          },
          {
            key   = "x-header-2"
            value = "testheader2"
          },
        ]
        HttpVerb               = "GET"
        ParseDependentRequests = false
        RequestUrl             = "http://microsoft.com"
      }
      RetryEnabled       = false
      SyntheticMonitorId = var.resource_name
      Timeout            = 30
      ValidationRules = {
        ExpectedHttpStatusCode = 200
        SSLCheck               = false
      }
    }
    tags = {
      "hidden-link:${azapi_resource.component.id}" = "Resource"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

