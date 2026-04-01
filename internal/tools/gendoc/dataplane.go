package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"

	tffwdocs "github.com/magodo/terraform-plugin-framework-docs"
)

type ApiPath struct {
	UrlFormat       string `json:"UrlFormat"`
	ResourceType    string `json:"ResourceType"`
	URL             string `json:"Url"`
	ParentIDExample string `json:"ParentIDExample"`
	ExamplePath     string `json:"ExamplePath"`
}

func genDataPlaneResource(ctx context.Context, gen *tffwdocs.Generator) error {
	apiData, err := os.ReadFile("internal/services/parse/data_plane_resources.json")
	if err != nil {
		return err
	}

	// Parse the JSON
	var apiPaths []ApiPath
	if err := json.Unmarshal(apiData, &apiPaths); err != nil {
		return fmt.Errorf("Error parsing JSON: %v", err)
	}

	// Sort by ResourceType for consistent output
	sort.Slice(apiPaths, func(i, j int) bool {
		return apiPaths[i].ResourceType < apiPaths[j].ResourceType
	})

	// Generate the markdown table
	table := generateMarkdownTable(apiPaths)

	// Generate example sections
	examples, err := generateExampleSections(apiPaths)
	if err != nil {
		return err
	}

	// #nosec G302
	f, err := os.OpenFile("./docs/resources/data_plane_resource.md", os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	opt := &tffwdocs.ResourceRenderOption{
		Examples: []tffwdocs.Example{
			{
				HCL: `
terraform {
  required_providers {
    azapi = {
      source = "Azure/azapi"
    }
  }
}

provider "azapi" {
}

provider "azurerm" {
  features {}
}

data "azurerm_synapse_workspace" "example" {
  name                = "example-workspace"
  resource_group_name = azurerm_resource_group.example.name
}

resource "azapi_data_plane_resource" "dataset" {
  type      = "Microsoft.Synapse/workspaces/datasets@2020-12-01"
  parent_id = trimprefix(data.azurerm_synapse_workspace.example.connectivity_endpoints.dev, "https://")
  name      = "example-dataset"
  body = {
    properties = {
      type = "AzureBlob",
      typeProperties = {
        folderPath = {
          value = "@dataset().MyFolderPath"
          type  = "Expression"
        }
        fileName = {
          value = "@dataset().MyFileName"
          type  = "Expression"
        }
        format = {
          type = "TextFormat"
        }
      }
      parameters = {
        MyFolderPath = {
          type = "String"
        }
        MyFileName = {
          type = "String"
        }
      }
    }
  }
}
`,
			},
		},
		ImportId: &tffwdocs.ImportId{
			Format:    "<resource_id>|<resource_type>",
			ExampleId: "exampleappconf.azconfig.io/kv/mykey|Microsoft.AppConfiguration/configurationStores/keyValues@1.0",
		},
		Template: template.Must(template.New("dataplane").Parse(fmt.Sprintf(`
{{ .Header }}
{{ .Description }}
{{- with .Example }}
{{ . }}
{{- end }}
{{ .Schema }}
## Available Resources

%s
## Resource Examples
%s


{{- with .Import }}
{{ . }}
{{- end }}
`, table, examples))),
	}

	return gen.RenderResource(ctx, f, "azapi_data_plane_resource", opt)
}

// generateMarkdownTable generates the markdown table from API paths
func generateMarkdownTable(apiPaths []ApiPath) string {
	var sb strings.Builder

	// Table header
	sb.WriteString("| Resource Type | URL | Parent ID Example                                                                           |\n")
	sb.WriteString("| --- | --- |---------------------------------------------------------------------------------------------|\n")

	// Table rows
	for _, path := range apiPaths {
		// Format the row with proper padding to align columns
		resourceType := path.ResourceType
		url := path.URL
		parentID := path.ParentIDExample

		// Write the row
		sb.WriteString(fmt.Sprintf("| %s | %s | %-91s |\n", resourceType, url, parentID))
	}

	return sb.String()
}

// generateExampleSections generates example usage sections for resources with ExamplePath
func generateExampleSections(apiPaths []ApiPath) (string, error) {
	var sb strings.Builder

	// Filter paths that have examples
	pathsWithExamples := []ApiPath{}
	for _, path := range apiPaths {
		if path.ExamplePath != "" {
			pathsWithExamples = append(pathsWithExamples, path)
		}
	}

	// If no examples, return empty string
	if len(pathsWithExamples) == 0 {
		return "", nil
	}

	// Generate sections for each example
	for _, path := range pathsWithExamples {
		// Create a section header based on the resource type
		sb.WriteString(fmt.Sprintf("\n### %s\n\n", path.ResourceType))

		b, err := os.ReadFile(path.ExamplePath)
		if err != nil {
			return "", fmt.Errorf("failed to read %s: %v", path.ExamplePath, err)
		}
		sb.WriteString(fmt.Sprintf("```terraform\n%s```\n", string(b)))
	}

	return sb.String(), nil
}
