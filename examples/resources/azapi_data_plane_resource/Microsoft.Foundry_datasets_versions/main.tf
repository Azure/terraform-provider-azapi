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

variable "source_url" {
  type        = string
  description = "HTTPS raw GitHub or Azure Blob Storage URL for the dataset file."
}

variable "source_sha256" {
  type        = string
  default     = null
  description = "Optional SHA-256 checksum for source_url."
}

resource "azapi_data_plane_resource" "dataset_version" {
  type      = "Microsoft.Foundry/datasets/versions@v1"
  parent_id = "${var.project_endpoint}/datasets/example-dataset"

  body = {
    name          = "example-dataset"
    description   = "Dataset uploaded by the AzAPI provider"
    type          = "uri_file"
    version       = "1"
    format        = "jsonl"
    source_url    = var.source_url
    source_sha256 = var.source_sha256
  }
}
