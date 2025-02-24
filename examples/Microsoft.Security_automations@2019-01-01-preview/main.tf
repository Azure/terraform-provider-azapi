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
