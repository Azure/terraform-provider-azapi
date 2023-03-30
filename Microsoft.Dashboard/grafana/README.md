# AZ MANAGED GRAFANA WITH TERRAFORM AZAPI & DEVOPS

Greetings my fellow Technology Advocates and Specialists.

In this Session, I will demonstrate, __How to deploy Azure Managed Grafana with Terraform AzAPI and DevOps__

| __LIVE RECORDED SESSION:-__ |
| --------- |
| __LIVE DEMO__ was Recorded as part of my Presentation in __JOURNEY TO THE CLOUD: WRAP UP - 2022__ Forum/Platform |
| Duration of My Demo = __37 Mins 27 Secs__ |
| [![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/V8b4hg-pZ7s/0.jpg)](https://www.youtube.com/watch?v=V8b4hg-pZ7s) |

| __WHAT IS TERRAFORM AZAPI PROVIDER:-__ |
| --------- |
| As Per the official documentation: __The AzAPI provider is a very thin layer on top of the Azure ARM REST APIs. This provider compliments the AzureRM provider by enabling the management of Azure resources that are not yet or may never be supported in the AzureRM provider such as private/public preview services and features.__ |
| For more information, please refer [AzAPI Documentation](https://registry.terraform.io/providers/Azure/azapi/latest/docs) |

| __USE CASE:-__ |
| --------- |
| How to deploy __Azure Managed Grafana__ using Terraform when the required __AzureRM__ Provider is __NOT__ available ? |

| __REQUIREMENTS:-__ |
| --------- |

1. Azure Subscription.
2. Azure DevOps Organisation and Project.
3. Service Principal with Required RBAC ( __Contributor__) applied on Subscription or Resource Group(s).
4. Azure Resource Manager Service Connection in Azure DevOps.
5. Microsoft DevLabs Terraform Extension Installed in Azure DevOps and in Local System (VS Code Extension).

| __OUT OF SCOPE:-__ |
| --------- |
| __Azure DevOps Pipeline Code Snippet Explanation.__ |


| __HOW DOES MY CODE PLACEHOLDER LOOKS LIKE:-__ |
| --------- |
| ![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/bq6gnowefssegacz0rpz.png) |

| PIPELINE CODE SNIPPET:- | 
| --------- |

| AZURE DEVOPS YAML PIPELINE (azure-pipelines-azapi-az-managed-grafana-v1.0):- | 
| --------- |

```
trigger:
  none

######################
#DECLARE PARAMETERS:-
######################
parameters:
- name: SubscriptionID
  displayName: Subscription ID Details Follow Below:-
  default: 210e66cb-55cf-424e-8daa-6cad804ab604
  values:
  -  210e66cb-55cf-424e-8daa-6cad804ab604

- name: ServiceConnection
  displayName: Service Connection Name Follows Below:-
  default: amcloud-cicd-service-connection
  values:
  -  amcloud-cicd-service-connection

######################
#DECLARE VARIABLES:-
######################
variables:
  ResourceGroup: tfpipeline-rg
  StorageAccount: tfpipelinesa
  Container: terraform
  TfstateFile: AMG/Grafana.tfstate
  BuildAgent: windows-latest
  WorkingDir: $(System.DefaultWorkingDirectory)/AzAPI-Az-Managed-Grafana
  Target: $(build.artifactstagingdirectory)/AMTF
  Environment: NonProd
  Artifact: AM

#########################
# Declare Build Agents:-
#########################
pool:
  vmImage: $(BuildAgent)

###################
# Declare Stages:-
###################
stages:

- stage: PLAN
  jobs:
  - job: PLAN
    displayName: PLAN
    steps:
# Install Terraform Installer in the Build Agent:-
    - task: ms-devlabs.custom-terraform-tasks.custom-terraform-installer-task.TerraformInstaller@0
      displayName: INSTALL TERRAFORM VERSION - LATEST
      inputs:
        terraformVersion: 'latest'
# Terraform Init:-
    - task: TerraformTaskV2@2
      displayName: TERRAFORM INIT
      inputs:
        provider: 'azurerm'
        command: 'init'
        workingDirectory: '$(workingDir)' # Az DevOps can find the required Terraform code
        backendServiceArm: '${{ parameters.ServiceConnection }}' 
        backendAzureRmResourceGroupName: '$(ResourceGroup)' 
        backendAzureRmStorageAccountName: '$(StorageAccount)'
        backendAzureRmContainerName: '$(Container)'
        backendAzureRmKey: '$(TfstateFile)'
# Terraform Validate:-
    - task: TerraformTaskV2@2
      displayName: TERRAFORM VALIDATE
      inputs:
        provider: 'azurerm'
        command: 'validate'
        workingDirectory: '$(workingDir)'
        environmentServiceNameAzureRM: '${{ parameters.ServiceConnection }}'
# Terraform Plan:-
    - task: TerraformTaskV2@2
      displayName: TERRAFORM PLAN
      inputs:
        provider: 'azurerm'
        command: 'plan'
        workingDirectory: '$(workingDir)'
        commandOptions: "--var-file=az-managed-grafana.tfvars --out=tfplan"
        environmentServiceNameAzureRM: '${{ parameters.ServiceConnection }}'
    
# Copy Files to Artifacts Staging Directory:-
    - task: CopyFiles@2
      displayName: COPY FILES ARTIFACTS STAGING DIRECTORY
      inputs:
        SourceFolder: '$(workingDir)'
        Contents: |
          **/*.tf
          **/*.tfvars
          **/*tfplan*
        TargetFolder: '$(Target)'
# Publish Artifacts:-
    - task: PublishBuildArtifacts@1
      displayName: PUBLISH ARTIFACTS
      inputs:
        targetPath: '$(Target)'
        artifactName: '$(Artifact)' 

- stage: DEPLOY
  condition: succeeded()
  dependsOn: PLAN
  jobs:
  - deployment: 
    displayName: Deploy
    environment: $(Environment)
    pool:
      vmImage: '$(BuildAgent)'
    strategy:
      runOnce:
        deploy:
          steps:
# Download Artifacts:-
          - task: DownloadBuildArtifacts@0
            displayName: DOWNLOAD ARTIFACTS
            inputs:
              buildType: 'current'
              downloadType: 'single'
              artifactName: '$(Artifact)'
              downloadPath: '$(System.ArtifactsDirectory)' 
# Install Terraform Installer in the Build Agent:-
          - task: ms-devlabs.custom-terraform-tasks.custom-terraform-installer-task.TerraformInstaller@0
            displayName: INSTALL TERRAFORM VERSION - LATEST
            inputs:
              terraformVersion: 'latest'
# Terraform Init:-
          - task: TerraformTaskV2@2 
            displayName: TERRAFORM INIT
            inputs:
              provider: 'azurerm'
              command: 'init'
              workingDirectory: '$(System.ArtifactsDirectory)/$(Artifact)/AMTF/' # Az DevOps can find the required Terraform code
              backendServiceArm: '${{ parameters.ServiceConnection }}' 
              backendAzureRmResourceGroupName: '$(ResourceGroup)' 
              backendAzureRmStorageAccountName: '$(StorageAccount)'
              backendAzureRmContainerName: '$(Container)'
              backendAzureRmKey: '$(TfstateFile)'
# Terraform Apply:-
          - task: TerraformTaskV2@2
            displayName: TERRAFORM APPLY # The terraform Plan stored earlier is used here to apply only the changes.
            inputs:
              provider: 'azurerm'
              command: 'apply'
              workingDirectory: '$(System.ArtifactsDirectory)/$(Artifact)/AMTF'
              commandOptions: '--var-file=az-managed-grafana.tfvars' # The terraform Plan stored earlier is used here to apply. 
              environmentServiceNameAzureRM: '${{ parameters.ServiceConnection }}'

```

| __TERRAFORM CODE SNIPPET:-__ |
| --------- |


| __TERRAFORM (main.tf):-__ |
| --------- |

```
terraform {

  required_version = ">= 1.3.3"
  
  backend "azurerm" {
    resource_group_name  = "tfpipeline-rg"
    storage_account_name = "tfpipelinesa"
    container_name       = "terraform"
    key                  = "AMG/Grafana.tfstate"
  }

  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = "~> 3.27"
    }
    azapi = {
      source = "Azure/azapi"
      version = "1.0.0"
    }
    
  }
}
provider "azurerm" {
  features {}
  skip_provider_registration = true
}

provider "azapi" {
}

```

| __EXPLANATION:-__ |
| --------- |
| Browse to the provided [LINK](https://registry.terraform.io/providers/Azure/azapi/latest/docs) in order to know, as how to use the AzAPI Provider. |
| Below is how it looks:- |
| ![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/ihlbty5o22myvwijo443.jpg) |


| __TERRAFORM (az-managed-grafana.tf):-__ |
| --------- |

```
#####################
## Resource Group:-
#####################

resource "azurerm_resource_group" "azrg" {
  name        = var.rg-name
  location    = var.rg-location
}

##############################
## Azure Managed Grafana:-
##############################
resource "azapi_resource" "azgrafana" {
  type        = "Microsoft.Dashboard/grafana@2022-08-01" 
  name        = var.az-grafana-name
  parent_id   = azurerm_resource_group.azrg.id
  location    = azurerm_resource_group.azrg.location
  
  identity {
    type      = "SystemAssigned"
  }

  body = jsonencode({
    sku = {
      name = "Standard"
    }
    properties = {
      publicNetworkAccess = "Enabled",
      zoneRedundancy = "Enabled",
      apiKey = "Enabled",
      deterministicOutboundIP = "Enabled"
    }
  })

}

```

| __EXPLANATION:-__ |
| --------- |
| The Reference Example of Terraform Resource "__azapi_resource__" can be found [HERE](https://registry.terraform.io/providers/Azure/azapi/latest/docs/resources/azapi_resource) |
| The __JSONENCODE Body definition__ was build using [Azure REST API Reference](https://learn.microsoft.com/en-us/rest/api/azure/) |
| From the above Documentation, browse for "__Managed Grafana__" to find all details to set the "__azapi_resource__" correctly for Azure Managed Grafana. Below is how it looks:- |
| ![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/wrsybfr0kfufvdv9u85r.jpg) |
| ![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/kc0pqohvirb4xm0ngwhh.jpg) |
| ![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/wj7duy78g6tm9qv5j7eq.jpg) |


| __TERRAFORM (variables.tf):-__ |
| --------- |

```
variable "rg-name" {
  type        = string
  description = "Name of the Resource Group"
}
variable "rg-location" {
  type        = string
  description = "Location of the Resource Group"
}
variable "az-grafana-name" {
  type        = string
  description = "Name of the Azure Managed Grafana"
}

```

| __NOTE:-__ |
| --------- |
| This is self-explanatory. |



| __TERRAFORM (az-managed-grafana.tf):-__ |
| --------- |

```
rg-name              = "GrafanaRG"
rg-location          = "West Europe"
az-grafana-name      = "AMGrafanaTest"

```

| __NOTE:-__ |
| --------- |
| This is self-explanatory. |


| __TEST RESULTS:-__ |
| --------- |
| __Pipeline Execution:-__ |
| ![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/vdhsz1y8o4jwoeq6bvy5.jpg) |
| __Resource Deployment:-__ |
| ![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/isqo9so3z5vk1pd6zxdf.jpg) |
| __Important to Note:-__ |
| Please add the required RBAC "__Grafana Admin__" to user identity who will access Azure Managed Grafana URL. |
| Below is how Azure Managed Grafana looks like:- |
| ![Image description](https://dev-to-uploads.s3.amazonaws.com/uploads/articles/9ipnz6g31xzkpxvpxk5p.jpg) |

__Hope You Enjoyed the Session!!!__

__Stay Safe | Keep Learning | Spread Knowledge__
