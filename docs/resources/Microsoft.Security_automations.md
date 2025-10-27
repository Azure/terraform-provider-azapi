---
subcategory: "Microsoft.Security - Security Center"
page_title: "automations"
description: |-
  Manages a Security Center Automation and Continuous Export.
---

# Microsoft.Security/automations - Security Center Automation and Continuous Export

This article demonstrates how to use `azapi` provider to manage the Security Center Automation and Continuous Export resource in Azure.



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

resource "azapi_resource" "workspace" {
  type      = "Microsoft.OperationalInsights/workspaces@2022-10-01"
  parent_id = azapi_resource.resourceGroup.id
  name      = var.resource_name
  location  = var.location
  body = {
    properties = {
      features = {
        disableLocalAuth                            = false
        enableLogAccessUsingOnlyResourcePermissions = true
      }
      publicNetworkAccessForIngestion = "Enabled"
      publicNetworkAccessForQuery     = "Enabled"
      retentionInDays                 = 30
      sku = {
        name = "PerGB2018"
      }
      workspaceCapping = {
        dailyQuotaGb = -1
      }
    }
  }
  schema_validation_enabled = false
  response_export_values    = ["*"]
}

resource "azapi_resource" "automation" {
  type      = "Microsoft.Security/automations@2019-01-01-preview"
  name      = "ExportToWorkspace"
  parent_id = azapi_resource.resourceGroup.id
  location  = azapi_resource.resourceGroup.location
  body = {
    properties = {
      isEnabled = true,
      scopes = [
        {
          description = "Security Export for the subscription",
          scopePath   = azapi_resource.resourceGroup.id
        }
      ],
      sources = [
        {
          eventSource = "Assessments",
          ruleSets = [
            {
              rules = [
                {
                  propertyJPath = "type",
                  propertyType  = "String",
                  expectedValue = "Microsoft.Security/assessments",
                  operator      = "Contains"
                }
              ]
            }
          ]
        },
        {
          eventSource = "AssessmentsSnapshot",
          ruleSets = [
            {
              rules = [
                {
                  propertyJPath = "type",
                  propertyType  = "String",
                  expectedValue = "Microsoft.Security/assessments",
                  operator      = "Contains"
                }
              ]
            }
          ]
        },
        {
          eventSource = "SubAssessments"
        },
        {
          eventSource = "SubAssessmentsSnapshot"
        },
        {
          eventSource = "Alerts",
          ruleSets = [
            {
              rules = [
                {
                  propertyJPath = "Severity",
                  propertyType  = "String",
                  expectedValue = "low",
                  operator      = "Equals"
                }
              ]
            },
            {
              rules = [
                {
                  propertyJPath = "Severity",
                  propertyType  = "String",
                  expectedValue = "medium",
                  operator      = "Equals"
                }
              ]
            },
            {
              rules = [
                {
                  propertyJPath = "Severity",
                  propertyType  = "String",
                  expectedValue = "high",
                  operator      = "Equals"
                }
              ]
            },
            {
              rules = [
                {
                  propertyJPath = "Severity",
                  propertyType  = "String",
                  expectedValue = "informational",
                  operator      = "Equals"
                }
              ]
            }
          ]
        },
        {
          eventSource = "SecureScores"
        },
        {
          eventSource = "SecureScoresSnapshot"
        },
        {
          eventSource = "SecureScoreControls"
        },
        {
          eventSource = "SecureScoreControlsSnapshot"
        },
        {
          eventSource = "RegulatoryComplianceAssessment"
        },
        {
          eventSource = "RegulatoryComplianceAssessmentSnapshot"
        }
      ],
      actions = [
        {
          workspaceResourceId = azapi_resource.workspace.id
          actionType          = "Workspace"
        }
      ]
    }
  }
}

```



## Arguments Reference

The following arguments are supported:

* `type` - (Required) The type of the resource. This should be set to `Microsoft.Security/automations@api-version`. The available api-versions for this resource are: [`2019-01-01-preview`, `2023-12-01-preview`].

* `parent_id` - (Required) The ID of the azure resource in which this resource is created. The allowed values are:  
  `/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}`

* `name` - (Required) Specifies the name of the azure resource. Changing this forces a new resource to be created.

* `body` - (Required) Specifies the configuration of the resource. More information about the arguments in `body` can be found in the [Microsoft documentation](https://learn.microsoft.com/en-us/azure/templates/Microsoft.Security/automations?pivots=deployment-language-terraform).

For other arguments, please refer to the [azapi_resource](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/resource) documentation.

## Import

 ```shell
 # Azure resource can be imported using the resource id, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/automations/{resourceName}
 
 # It also supports specifying API version by using the resource id with api-version as a query parameter, e.g.
 terraform import azapi_resource.example /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/automations/{resourceName}?api-version=2023-12-01-preview
 ```
