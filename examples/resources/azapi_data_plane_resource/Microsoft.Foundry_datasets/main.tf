terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {}

variable "project_endpoint" {
  type        = string
  description = "Foundry project endpoint, for example contoso.services.ai.azure.com/api/projects/example."
}

variable "dataset_name" {
  description = "The name of the dataset."
  type        = string
}

variable "dataset_version" {
  description = "The version of the dataset."
  type        = string
}

variable "dataset_description" {
  description = "The description of the dataset."
  type        = string
}

variable "dataset_type" {
  description = "The dataset type. The current upload workflow supports uri_file."
  type        = string
  default     = "uri_file"

  validation {
    condition     = var.dataset_type == "uri_file"
    error_message = "dataset_type must be uri_file because folder uploads are not supported by this workflow."
  }
}

variable "dataset_format" {
  description = "The format of the dataset, e.g., jsonl."
  type        = string
  default     = "jsonl"
}

variable "source_url" {
  type        = string
  description = "URL for the dataset file. Add module-specific host allowlist validation as appropriate."

  validation {
    condition     = can(regex("^https://", var.source_url))
    error_message = "source_url must be an HTTPS URL."
  }
}

variable "source_sha256" {
  type        = string
  default     = null
  description = "Optional SHA-256 checksum for source_url."
}

resource "azapi_data_plane_resource" "dataset" {
  type = "Microsoft.Foundry/datasets/versions@2025-05-01"

  parent_id = "${var.project_endpoint}/datasets/${var.dataset_name}"

  # The generic AzAPI resource calls this attribute "name".
  # For Foundry dataset versions, it is the version in the REST path:
  #
  # {project-endpoint}/datasets/{dataset_name}/versions/{dataset_version}
  name = var.dataset_version

  body = {
    name          = var.dataset_name
    version       = var.dataset_version
    description   = var.dataset_description
    type          = var.dataset_type
    format        = var.dataset_format
    source_url    = var.source_url
    source_sha256 = var.source_sha256
  }
}

output "dataset_id" {
  value = azapi_data_plane_resource.dataset.id
}

output "computed_sha256" {
  value = azapi_data_plane_resource.dataset.output.computed_sha256
}
