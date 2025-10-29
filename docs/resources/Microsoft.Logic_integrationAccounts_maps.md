---
subcategory: "Microsoft.Logic - Logic Apps"
page_title: "integrationAccounts/maps"
description: |-
  Manages a Logic App Integration Account Map.
---

# Microsoft.Logic/integrationAccounts/maps - Logic App Integration Account Map

This article demonstrates how to use `azapi` provider to manage the Logic App Integration Account Map resource in Azure.



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

resource "azapi_resource" "integrationAccount" {
  type      = "Microsoft.Logic/integrationAccounts@2019-05-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
    }
    sku = {
      name = "Basic"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "map" {
  type      = "Microsoft.Logic/integrationAccounts/maps@2019-05-01"
  parent_id = azapi_resource.integrationAccount.id
  name      = var.resource_name
  body = {
    properties = {
      content     = "<xsl:stylesheet xmlns:xsl=\"http://www.w3.org/1999/XSL/Transform\"\n                xmlns:msxsl=\"urn:schemas-microsoft-com:xslt\"\n                xmlns:var=\"http://schemas.microsoft.com/BizTalk/2003/var\"\n                exclude-result-prefixes=\"msxsl var s0 userCSharp\"\n                version=\"1.0\"\n                xmlns:ns0=\"http://BizTalk_Server_Project4.StringFunctoidsDestinationSchema\"\n                xmlns:s0=\"http://BizTalk_Server_Project4.StringFunctoidsSourceSchema\"\n                xmlns:userCSharp=\"http://schemas.microsoft.com/BizTalk/2003/userCSharp\">\n<xsl:import href=\"http://btsfunctoids.blob.core.windows.net/functoids/functoids.xslt\" />\n<xsl:output omit-xml-declaration=\"yes\"\n            method=\"xml\"\n            version=\"1.0\" />\n<xsl:template match=\"/\">\n<xsl:apply-templates select=\"/s0:Root\" />\n</xsl:template>\n<xsl:template match=\"/s0:Root\">\n<xsl:variable name=\"var:v1\"\n              select=\"userCSharp:StringFind(string(StringFindSource/text()) , &quot;SearchString&quot;)\" />\n<xsl:variable name=\"var:v2\"\n              select=\"userCSharp:StringLeft(string(StringLeftSource/text()) , &quot;2&quot;)\" />\n<xsl:variable name=\"var:v3\"\n              select=\"userCSharp:StringRight(string(StringRightSource/text()) , &quot;2&quot;)\" />\n<xsl:variable name=\"var:v4\"\n              select=\"userCSharp:StringUpperCase(string(UppercaseSource/text()))\" />\n<xsl:variable name=\"var:v5\"\n              select=\"userCSharp:StringLowerCase(string(LowercaseSource/text()))\" />\n<xsl:variable name=\"var:v6\"\n              select=\"userCSharp:StringSize(string(SizeSource/text()))\" />\n<xsl:variable name=\"var:v7\"\n              select=\"userCSharp:StringSubstring(string(StringExtractSource/text()) , &quot;0&quot; , &quot;2&quot;)\" />\n<xsl:variable name=\"var:v8\"\n              select=\"userCSharp:StringConcat(string(StringConcatSource/text()))\" />\n<xsl:variable name=\"var:v9\"\n              select=\"userCSharp:StringTrimLeft(string(StringLeftTrimSource/text()))\" />\n<xsl:variable name=\"var:v10\"\n              select=\"userCSharp:StringTrimRight(string(StringRightTrimSource/text()))\" />\n<ns0:Root>\n<StringFindDestination>\n<xsl:value-of select=\"$var:v1\" />\n</StringFindDestination>\n<StringLeftDestination>\n<xsl:value-of select=\"$var:v2\" />\n</StringLeftDestination>\n<StringRightDestination>\n<xsl:value-of select=\"$var:v3\" />\n</StringRightDestination>\n<UppercaseDestination>\n<xsl:value-of select=\"$var:v4\" />\n</UppercaseDestination>\n<LowercaseDestination>\n<xsl:value-of select=\"$var:v5\" />\n</LowercaseDestination>\n<SizeDestination>\n<xsl:value-of select=\"$var:v6\" />\n</SizeDestination>\n<StringExtractDestination>\n<xsl:value-of select=\"$var:v7\" />\n</StringExtractDestination>\n<StringConcatDestination>\n<xsl:value-of select=\"$var:v8\" />\n</StringConcatDestination>\n<StringLeftTrimDestination>\n<xsl:value-of select=\"$var:v9\" />\n</StringLeftTrimDestination>\n<StringRightTrimDestination>\n<xsl:value-of select=\"$var:v10\" />\n</StringRightTrimDestination>\n</ns0:Root>\n</xsl:template>\n</xsl:stylesheet>\n"
      contentType = "application/xml"
      mapType     = "Xslt"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Logic/integrationAccounts/maps@api-version`. The available api-versions for this resource are: [`2015-08-01-preview`, `2016-06-01`, `2018-07-01-preview`, `2019-05-01`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Logic/integrationAccounts/maps?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{resourceName}/maps/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{resourceName}/maps/{resourceName}?api-version=2019-05-01
 ```
