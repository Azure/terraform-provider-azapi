// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type ApiPath struct {
	UrlFormat       string `json:"UrlFormat"`
	ResourceType    string `json:"ResourceType"`
	URL             string `json:"Url"`
	ParentIDExample string `json:"ParentIDExample"`
	ExamplePath     string `json:"ExamplePath"`
}

func main() {
	// Default paths are relative to the repository root
	inputFile := flag.String("input", "data_plane_resources.json", "path to data_plane_resources.json")
	outputFile := flag.String("output", "", "path to template file")
	flag.Parse()

	// If output file is not specified, default to the repository root path
	if *outputFile == "" {
		*outputFile = "templates/resources/data_plane_resource.md.tmpl"
	}

	// If input file is relative, try to find it in the same directory as the script
	if *inputFile == "data_plane_resources.json" {
		*inputFile = "internal/services/parse/data_plane_resources.json"
	}

	// Read the JSON file
	content, err := os.ReadFile(*inputFile)
	if err != nil {
		log.Fatalf("Error reading input file: %v", err)
	}

	// Parse the JSON
	var apiPaths []ApiPath
	if err := json.Unmarshal(content, &apiPaths); err != nil {
		log.Fatalf("Error parsing JSON: %v", err)
	}

	// Sort by ResourceType for consistent output
	sort.Slice(apiPaths, func(i, j int) bool {
		return apiPaths[i].ResourceType < apiPaths[j].ResourceType
	})

	// Generate the markdown table
	table := generateMarkdownTable(apiPaths)

	// Generate example sections
	examples := generateExampleSections(apiPaths)

	// Read the template file
	templateContent, err := os.ReadFile(*outputFile)
	if err != nil {
		log.Fatalf("Error reading template file: %v", err)
	}

	// Replace the Available Resources section
	newContent := replaceAvailableResourcesSection(string(templateContent), table, examples)

	// Write the updated template
	if err := os.WriteFile(*outputFile, []byte(newContent), 0644); err != nil {
		log.Fatalf("Error writing template file: %v", err)
	}

	fmt.Printf("Successfully generated data plane resources table in %s\n", *outputFile)
	fmt.Printf("Total resources: %d\n", len(apiPaths))
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
func generateExampleSections(apiPaths []ApiPath) string {
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
		return ""
	}

	// Generate sections for each example
	for _, path := range pathsWithExamples {
		// Create a section header based on the resource type
		sb.WriteString(fmt.Sprintf("\n### %s\n\n", path.ResourceType))
		sb.WriteString(fmt.Sprintf("{{ tffile %q}}\n", path.ExamplePath))
	}

	return sb.String()
}

// replaceAvailableResourcesSection replaces the Available Resources section in the template
func replaceAvailableResourcesSection(templateContent, newTable, examples string) string {
	// Find the start of the Available Resources section
	sectionStart := "## Available Resources\n\n"
	startIdx := strings.Index(templateContent, sectionStart)
	if startIdx == -1 {
		log.Fatal("Could not find '## Available Resources' section in template")
	}

	// Position after the section header
	contentStart := startIdx + len(sectionStart)

	// The new content starts with the table, followed by examples
	newSection := newTable
	if examples != "" {
		newSection += "\n## Resource Examples\n" + examples
	}

	// Check if there's content after the table (we want to preserve it)
	// Find the end of the current table (next section or end of file)
	remainingContent := templateContent[contentStart:]

	// Find the start of the next section (if any) - starts with ## or end of file
	// Skip the "Resource Examples" section if it exists
	nextSectionIdx := strings.Index(remainingContent, "\n## ")

	// If there's a Resource Examples section, skip past it
	if nextSectionIdx != -1 {
		possibleExamplesSection := remainingContent[nextSectionIdx:]
		if strings.HasPrefix(possibleExamplesSection, "\n## Resource Examples") {
			// Find the section after Resource Examples
			afterExamples := possibleExamplesSection[len("\n## Resource Examples"):]
			nextNextSectionIdx := strings.Index(afterExamples, "\n## ")
			if nextNextSectionIdx != -1 {
				nextSectionIdx = nextSectionIdx + len("\n## Resource Examples") + nextNextSectionIdx
			} else {
				// No more sections after Resource Examples
				nextSectionIdx = -1
			}
		}
	}

	var afterTable string
	if nextSectionIdx != -1 {
		afterTable = remainingContent[nextSectionIdx:]
	}

	// Construct the new content
	return templateContent[:contentStart] + newSection + afterTable
}
