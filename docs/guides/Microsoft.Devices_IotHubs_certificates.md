---
subcategory: "Microsoft.Devices - Azure IoT Hub, Azure IoT Hub Device Provisioning Service"
page_title: "IotHubs/certificates"
description: |-
  Manages a IoTHub Certificate.
---

# Microsoft.Devices/IotHubs/certificates - IoTHub Certificate

This article demonstrates how to use `azapi` provider to manage the IoTHub Certificate resource in Azure.

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

resource "azapi_resource" "IotHub" {
  type      = "Microsoft.Devices/IotHubs@2022-04-30-preview"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      cloudToDevice = {
      }
      enableFileUploadNotifications = false
      messagingEndpoints = {
      }
      routing = {
        fallbackRoute = {
          condition = "true"
          endpointNames = [
            "events",
          ]
          isEnabled = true
          source    = "DeviceMessages"
        }
      }
      storageEndpoints = {
      }
    }
    sku = {
      capacity = 1
      name     = "B1"
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "certificate" {
  type      = "Microsoft.Devices/IotHubs/certificates@2022-04-30-preview"
  parent_id = azapi_resource.IotHub.id
  name      = var.resource_name
  body = {
    properties = {
      certificate = "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tDQpNSUlEYnpDQ0FsZWdBd0lCQWdJSkFJempSRDM2c0liYk1BMEdDU3FHU0liM0RRRUJDd1VBTUUweEN6QUpCZ05WDQpCQVlUQWxWVE1STXdFUVlEVlFRSURBcFRiMjFsTFZOMFlYUmxNUkl3RUFZRFZRUUtEQWwwWlhKeVlXWnZjbTB4DQpGVEFUQmdOVkJBTU1ESFJsY25KaFptOXliUzVwYnpBZ0Z3MHhOekEwTWpFeU1EQTFNamRhR0E4eU1URTNNRE15DQpPREl3TURVeU4xb3dUVEVMTUFrR0ExVUVCaE1DVlZNeEV6QVJCZ05WQkFnTUNsTnZiV1V0VTNSaGRHVXhFakFRDQpCZ05WQkFvTUNYUmxjbkpoWm05eWJURVZNQk1HQTFVRUF3d01kR1Z5Y21GbWIzSnRMbWx2TUlJQklqQU5CZ2txDQpoa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDQVFFQTNMOUw1c3pUNCtGTHlrVEZOeXlQankvazNCUVRZQWZSDQpRelAyZGhuc3VVS20zY2RQQzBOeVord0VYSVVHaG9ETzJZRzZFWUNoT2w4ZnNEcURPamxvU1VHS3FZdysrbmxwDQpISXVVZ0p4OEl4eEcyWGtBTENqRlU3RW1GK3c3a243NmQwZXpwRUlZeG5MUCtLRzJEVm9ybm9FdDFhTGh2MU1MDQptcGdFWlpQaERiTVNMaFNZV2VUVlJNYXlYTHdxdGZnbkR1bVFTQis4ZC8xSnVKcXJTSTRwRDEySm96VlRoemI2DQpoc2pmYjZSTVg0ZXBQbXJHbjBQYlRQRUVBNmF3bXN4QkNYQjBzMTNuTlF0L08waExNMmFnd3ZBeW96aWxRVitzDQo2MTZDa2drNkRKb1VrcVpoRHk3dlBZTUlSU3I5OGZCd3M2emtyVjZ0VExqbUQ4eEF2b2JlUFFJREFRQUJvMUF3DQpUakFkQmdOVkhRNEVGZ1FVWElxTzQyMXpNTW1iY1JSWDl3Y3RaRkNRdVBJd0h3WURWUjBqQkJnd0ZvQVVYSXFPDQo0MjF6TU1tYmNSUlg5d2N0WkZDUXVQSXdEQVlEVlIwVEJBVXdBd0VCL3pBTkJna3Foa2lHOXcwQkFRc0ZBQU9DDQpBUUVBcjgyTmVUM0JZSk9LTGxVTDZPbTVMalVGNjZld2NKakc5bHRkdnlRd1ZuZU1jcTd0NVVBUHhnQ2h6cU5SDQpWazRkYThQemtYcGpCSnlXZXpIdXBkSk5YM1hxZVVrMmtTeHFRNi9nbWhxdmZJM3k3ZGpyd29PNmp2TUVZMjZXDQpxdGtUTk9SV0RQM1RISkpWaW1DM3pWK0tNVTVVQlZyRXpoT1ZoSFNVNzA5bEJQNzVvMEJCbjN4R3NQcVNxOWs4DQpJb3RJRmZ5QWM2YStYUDMrWk1wdmg3d3FBVW1sN3ZXYTV3bGNYRXhDeDM5aDFiYWxmRFNMR05DNHN3V1BDcDlBDQpNblFSMHArdk1heTloTlAxRWgrOVFZVWFpMTRkNUtTM2NGVitLeEUxY0pSNUhEL2lMbHRubk9FYnBNc0IwZVZPDQpaV2tGdkU3WTVsVzBvVlNBZmluNVR3VEpNUT09DQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0t"
      isVerified  = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}


```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Devices/IotHubs/certificates@api-version`. The available api-versions for this resource are: [`2017-07-01`, `2018-01-22`, `2018-04-01`, `2018-12-01-preview`, `2019-03-22`, `2019-03-22-preview`, `2019-07-01-preview`, `2019-11-04`, `2020-03-01`, `2020-04-01`, `2020-06-15`, `2020-07-10-preview`, `2020-08-01`, `2020-08-31`, `2020-08-31-preview`, `2021-02-01-preview`, `2021-03-03-preview`, `2021-03-31`, `2021-07-01`, `2021-07-01-preview`, `2021-07-02`, `2021-07-02-preview`, `2022-04-30-preview`, `2022-11-15-preview`, `2023-06-30`, `2023-06-30-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Devices/IotHubs/{resourceName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Devices/IotHubs/certificates?pivots=deployment-language-terraform).

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Devices/IotHubs/{resourceName}/certificates/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Devices/IotHubs/{resourceName}/certificates/{resourceName}?api-version=2023-06-30-preview
 ```
