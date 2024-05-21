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
  type      = "Microsoft.Kusto/clusters@2023-05-02"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  identity {
    type         = "SystemAssigned"
    identity_ids = []
  }
  body = {
    properties = {
      enableAutoStop                = true
      enableDiskEncryption          = false
      enableDoubleEncryption        = false
      enablePurge                   = false
      enableStreamingIngest         = false
      engineType                    = "V2"
      publicIPType                  = "IPv4"
      publicNetworkAccess           = "Enabled"
      restrictOutboundNetworkAccess = "Disabled"
      trustedExternalTenants = [
      ]
    }
    sku = {
      capacity = 1
      name     = "Dev(No SLA)_Standard_D11_v2"
      tier     = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "database" {
  type      = "Microsoft.Kusto/clusters/databases@2023-05-02"
  parent_id = azapi_resource.cluster.id
  name      = var.resource_name
  location  = var.location
  body = {
    kind = "ReadWrite"
    properties = {
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "script" {
  type      = "Microsoft.Kusto/clusters/databases/scripts@2023-05-02"
  parent_id = azapi_resource.database.id
  name      = "create-table-script"
  body = {
    properties = {
      continueOnErrors = false
      forceUpdateTag   = "9e2e7874-aa37-7041-81b7-06397f03a37d"
      scriptContent    = ".create table TestTable(Id:string, Name:string, _ts:long, _timestamp:datetime)\n.create table TestTable ingestion json mapping \"TestMapping\"\n'['\n'    {\"column\":\"Id\",\"path\":\"$.id\"},'\n'    {\"column\":\"Name\",\"path\":\"$.name\"},'\n'    {\"column\":\"_ts\",\"path\":\"$._ts\"},'\n'    {\"column\":\"_timestamp\",\"path\":\"$._ts\", \"transform\":\"DateTimeFromUnixSeconds\"}'\n']'\n.alter table TestTable policy ingestionbatching \"{'MaximumBatchingTimeSpan': '0:0:10', 'MaximumNumberOfItems': 10000}\"\n"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

