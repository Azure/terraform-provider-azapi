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

provider "azurerm" {
  features {
  }
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

data "azuread_service_principal" "test" {
  display_name = "Windows Virtual Desktop"
}

resource "azurerm_role_assignment" "test" {
  name                             = "4a23d649-5bf7-41b8-9701-c17f811f68da"
  scope                            = azapi_resource.resourceGroup.id
  role_definition_name             = "Desktop Virtualization Power On Off Contributor"
  principal_id                     = data.azuread_service_principal.test.object_id
  skip_service_principal_aad_check = true
}

resource "azapi_resource" "hostPool" {
  type      = "Microsoft.DesktopVirtualization/hostPools@2023-09-05"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      hostPoolType          = "Personal"
      loadBalancerType      = "Persistent"
      maxSessionLimit       = 999999
      preferredAppGroupType = "Desktop"
      publicNetworkAccess   = "Enabled"
      startVMOnConnect      = true
      validationEnvironment = false
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "scalingPlan" {
  depends_on = [azurerm_role_assignment.test]
  type       = "Microsoft.DesktopVirtualization/scalingPlans@2023-11-01-preview"
  name       = var.resource_name
  location   = var.location
  parent_id  = azapi_resource.resourceGroup.id
  body = {
    properties = {
      timeZone     = "W. Europe Standard Time"
      hostPoolType = "Personal"
      exclusionTag = "no-schedule"
      schedules    = []
      hostPoolReferences = [
        {
          hostPoolArmPath    = azapi_resource.hostPool.id,
          scalingPlanEnabled = true
        }
      ]
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "personalSchedule" {
  type      = "Microsoft.DesktopVirtualization/scalingPlans/personalSchedules@2023-11-01-preview"
  name      = "Weekdays"
  parent_id = azapi_resource.scalingPlan.id
  body = {
    properties = {
      daysOfWeek = [
        "Monday",
        "Tuesday",
        "Wednesday",
        "Thursday",
        "Friday"
      ]

      rampUpStartTime = {
        hour   = 7
        minute = 0
      }

      rampUpAutoStartHosts            = "None"
      rampUpStartVMOnConnect          = "Enable"
      rampUpMinutesToWaitOnDisconnect = 45
      rampUpActionOnDisconnect        = "Hibernate"
      rampUpMinutesToWaitOnLogoff     = 30
      rampUpActionOnLogoff            = "Hibernate"

      peakStartTime = {
        hour   = 8
        minute = 0
      }

      peakStartVMOnConnect          = "Enable"
      peakMinutesToWaitOnDisconnect = 60
      peakActionOnDisconnect        = "Hibernate"
      peakMinutesToWaitOnLogoff     = 60
      peakActionOnLogoff            = "Hibernate"

      rampDownStartTime = {
        hour   = 16
        minute = 30
      }

      rampDownStartVMOnConnect          = "Enable"
      rampDownMinutesToWaitOnDisconnect = 45
      rampDownActionOnDisconnect        = "Hibernate"
      rampDownMinutesToWaitOnLogoff     = 30
      rampDownActionOnLogoff            = "Hibernate"

      offPeakStartTime = {
        hour   = 17
        minute = 30
      }

      offPeakStartVMOnConnect          = "Enable"
      offPeakMinutesToWaitOnDisconnect = 20
      offPeakActionOnDisconnect        = "Hibernate"
      offPeakMinutesToWaitOnLogoff     = 15
      offPeakActionOnLogoff            = "Hibernate"
    }
  }
}
