---
subcategory: "Microsoft.ApiManagement - API Management"
page_title: "service/schemas"
description: |-
  Manages a Global Schema within an API Management Service.
---

# Microsoft.ApiManagement/service/schemas - Global Schema within an API Management Service

This article demonstrates how to use `azapi` provider to manage the Global Schema within an API Management Service resource in Azure.



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


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ApiManagement/service/schemas@api-version`. The available api-versions for this resource are: [`2021-04-01-preview`, `2021-08-01`, `2021-12-01-preview`, `2022-04-01-preview`, `2022-08-01`, `2022-09-01-preview`, `2023-03-01-preview`, `2023-05-01-preview`, `2023-09-01-preview`, `2024-05-01`, `2024-06-01-preview`, `2024-10-01-preview`, `2025-03-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ApiManagement/service/schemas?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}/schemas/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}/schemas/{resourceName}?api-version=2025-03-01-preview
 ```
