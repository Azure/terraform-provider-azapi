---
subcategory: "Microsoft.ApiManagement - API Management"
page_title: "service/apis/schemas"
description: |-
  Manages a API Schema within an API Management Service.
---

# Microsoft.ApiManagement/service/apis/schemas - API Schema within an API Management Service

This article demonstrates how to use `azapi` provider to manage the API Schema within an API Management Service resource in Azure.

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

resource "azapi_resource" "api" {
  type      = "Microsoft.ApiManagement/service/apis@2021-08-01"
  parent_id = azapi_resource.service.id
  name      = "${var.resource_name};rev=1"
  body = {
    properties = {
      apiRevisionDescription = ""
      apiType                = "http"
      apiVersion             = ""
      apiVersionDescription  = ""
      authenticationSettings = {
      }
      description = ""
      displayName = "api1"
      path        = "api1"
      protocols = [
        "https",
      ]
      serviceUrl           = ""
      subscriptionRequired = true
      type                 = "http"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "schema" {
  type      = "Microsoft.ApiManagement/service/apis/schemas@2021-08-01"
  parent_id = azapi_resource.api.id
  name      = var.resource_name
  body = {
    properties = {
      contentType = "application/vnd.ms-azure-apim.xsd+xml"
      document = {
        value = "<!--\n Copyright (c) HashiCorp, Inc.\n SPDX-License-Identifier: MPL-2.0\n-->\n\n<s:schema elementFormDefault=\"qualified\" targetNamespace=\"http://ws.cdyne.com/WeatherWS/\" xmlns:tns=\"http://ws.cdyne.com/WeatherWS/\" xmlns:s=\"http://www.w3.org/2001/XMLSchema\" xmlns:soap12=\"http://schemas.xmlsoap.org/wsdl/soap12/\" xmlns:mime=\"http://schemas.xmlsoap.org/wsdl/mime/\" xmlns:soap=\"http://schemas.xmlsoap.org/wsdl/soap/\" xmlns:tm=\"http://microsoft.com/wsdl/mime/textMatching/\" xmlns:http=\"http://schemas.xmlsoap.org/wsdl/http/\" xmlns:soapenc=\"http://schemas.xmlsoap.org/soap/encoding/\" xmlns:wsdl=\"http://schemas.xmlsoap.org/wsdl/\" xmlns:apim-wsdltns=\"http://ws.cdyne.com/WeatherWS/\"> <s:element name=\"GetWeatherInformation\"> <s:complexType /> </s:element> <s:element name=\"GetWeatherInformationResponse\"> <s:complexType> <s:sequence> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"GetWeatherInformationResult\" type=\"tns:ArrayOfWeatherDescription\" /> </s:sequence> </s:complexType> </s:element> <s:complexType name=\"ArrayOfWeatherDescription\"> <s:sequence> <s:element minOccurs=\"0\" maxOccurs=\"unbounded\" name=\"WeatherDescription\" type=\"tns:WeatherDescription\" /> </s:sequence> </s:complexType> <s:complexType name=\"WeatherDescription\"> <s:sequence> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"WeatherID\" type=\"s:short\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Description\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"PictureURL\" type=\"s:string\" /> </s:sequence> </s:complexType> <s:element name=\"GetCityForecastByZIP\"> <s:complexType> <s:sequence> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ZIP\" type=\"s:string\" /> </s:sequence> </s:complexType> </s:element> <s:element name=\"GetCityForecastByZIPResponse\"> <s:complexType> <s:sequence> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"GetCityForecastByZIPResult\" type=\"tns:ForecastReturn\" /> </s:sequence> </s:complexType> </s:element> <s:complexType name=\"ForecastReturn\"> <s:sequence> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"Success\" type=\"s:boolean\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ResponseText\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"State\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"City\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"WeatherStationCity\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ForecastResult\" type=\"tns:ArrayOfForecast\" /> </s:sequence> </s:complexType> <s:complexType name=\"ArrayOfForecast\"> <s:sequence> <s:element minOccurs=\"0\" maxOccurs=\"unbounded\" name=\"Forecast\" nillable=\"true\" type=\"tns:Forecast\" /> </s:sequence> </s:complexType> <s:complexType name=\"Forecast\"> <s:sequence> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"Date\" type=\"s:dateTime\" /> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"WeatherID\" type=\"s:short\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Desciption\" type=\"s:string\" /> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"Temperatures\" type=\"tns:temp\" /> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"ProbabilityOfPrecipiation\" type=\"tns:POP\" /> </s:sequence> </s:complexType> <s:complexType name=\"temp\"> <s:sequence> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"MorningLow\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"DaytimeHigh\" type=\"s:string\" /> </s:sequence> </s:complexType> <s:complexType name=\"POP\"> <s:sequence> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Nighttime\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Daytime\" type=\"s:string\" /> </s:sequence> </s:complexType> <s:element name=\"GetCityWeatherByZIP\"> <s:complexType> <s:sequence> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ZIP\" type=\"s:string\" /> </s:sequence> </s:complexType> </s:element> <s:element name=\"GetCityWeatherByZIPResponse\"> <s:complexType> <s:sequence> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"GetCityWeatherByZIPResult\" type=\"tns:WeatherReturn\" /> </s:sequence> </s:complexType> </s:element> <s:complexType name=\"WeatherReturn\"> <s:sequence> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"Success\" type=\"s:boolean\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"ResponseText\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"State\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"City\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"WeatherStationCity\" type=\"s:string\" /> <s:element minOccurs=\"1\" maxOccurs=\"1\" name=\"WeatherID\" type=\"s:short\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Description\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Temperature\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"RelativeHumidity\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Wind\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Pressure\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Visibility\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"WindChill\" type=\"s:string\" /> <s:element minOccurs=\"0\" maxOccurs=\"1\" name=\"Remarks\" type=\"s:string\" /> </s:sequence> </s:complexType> <s:element name=\"ArrayOfWeatherDescription\" nillable=\"true\" type=\"tns:ArrayOfWeatherDescription\" /> <s:element name=\"ForecastReturn\" nillable=\"true\" type=\"tns:ForecastReturn\" /> <s:element name=\"WeatherReturn\" type=\"tns:WeatherReturn\" /> </s:schema>"
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.ApiManagement/service/apis/schemas@api-version`. The available api-versions for this resource are: [`2017-03-01`, `2018-01-01`, `2018-06-01-preview`, `2019-01-01`, `2019-12-01`, `2019-12-01-preview`, `2020-06-01-preview`, `2020-12-01`, `2021-01-01-preview`, `2021-04-01-preview`, `2021-08-01`, `2021-12-01-preview`, `2022-04-01-preview`, `2022-08-01`, `2022-09-01-preview`, `2023-03-01-preview`, `2023-05-01-preview`, `2023-09-01-preview`, `2024-05-01`, `2024-06-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}/apis/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.ApiManagement/service/apis/schemas?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}/apis/{resourceName}/schemas/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{resourceName}/apis/{resourceName}/schemas/{resourceName}?api-version=2024-06-01-preview
 ```
