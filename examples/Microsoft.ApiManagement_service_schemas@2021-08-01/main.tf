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

resource "azapi_resource" "service" {
  type      = "Microsoft.ApiManagement/service@2021-08-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      certificates = [
      ]
      customProperties = {
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Ssl30" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls10" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Backend.Protocols.Tls11" = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls10"         = "false"
        "Microsoft.WindowsAzure.ApiManagement.Gateway.Security.Protocols.Tls11"         = "false"
      }
      disableGateway      = false
      publicNetworkAccess = "Enabled"
      publisherEmail      = "pub1@email.com"
      publisherName       = "pub1"
      virtualNetworkType  = "None"
    }
    sku = {
      capacity = 0
      name     = "Consumption"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "schema" {
  type      = "Microsoft.ApiManagement/service/schemas@2021-08-01"
  parent_id = azapi_resource.service.id
  name      = var.resource_name
  body = {
    properties = {
      description = ""
      schemaType  = "xml"
      value       = "    <xsd:schema xmlns:xsd=\"http://www.w3.org/2001/XMLSchema\"\n    xmlns:tns=\"http://tempuri.org/PurchaseOrderSchema.xsd\" targetNamespace=\"http://tempuri.org/PurchaseOrderSchema.xsd\" elementFormDefault=\"qualified\">\n    <xsd:element name=\"PurchaseOrder\" type=\"tns:PurchaseOrderType\"/>\n    <xsd:complexType name=\"PurchaseOrderType\">\n        <xsd:sequence>\n            <xsd:element name=\"ShipTo\" type=\"tns:USAddress\" maxOccurs=\"2\"/>\n            <xsd:element name=\"BillTo\" type=\"tns:USAddress\"/>\n        </xsd:sequence>\n        <xsd:attribute name=\"OrderDate\" type=\"xsd:date\"/>\n    </xsd:complexType>\n    <xsd:complexType name=\"USAddress\">\n        <xsd:sequence>\n            <xsd:element name=\"name\" type=\"xsd:string\"/>\n            <xsd:element name=\"street\" type=\"xsd:string\"/>\n            <xsd:element name=\"city\" type=\"xsd:string\"/>\n            <xsd:element name=\"state\" type=\"xsd:string\"/>\n            <xsd:element name=\"zip\" type=\"xsd:integer\"/>\n        </xsd:sequence>\n        <xsd:attribute name=\"country\" type=\"xsd:NMTOKEN\" fixed=\"US\"/>\n    </xsd:complexType>\n</xsd:schema>\n"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

