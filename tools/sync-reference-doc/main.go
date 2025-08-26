// sync-reference-doc
//
// Purpose
//   Transforms Terraform example folders from an input directory (inDir)
//   into a docs-friendly output structure (outDir) and generates a remarks.json
//   index per namespace that lists the sample files.
//
// Key behaviors
//   - Input folder layout: examples are grouped by resource-type-and-version folder names
//       <Namespace>_<type_path_tokens_separated_by_"_">@<apiVersion>/[<scenario>]/main.tf
//       e.g. Microsoft.AppPlatform_Spring_apps_deployments@2024-01-01/basic/main.tf
//   - Multiple scenarios: each resource folder may contain one or more scenario subfolders;
//       each scenario folder must contain a main.tf. We also accept main.tf at the resource
//       folder root when no scenario subfolders exist.
//   - Scenario normalization: scenario names are sanitized by removing underscores and
//       trimming whitespace; output paths and descriptions use the sanitized name.
//   - Output folder layout: per-namespace folder under outDir, with nested sample folders
//       by all type path tokens (lowercased) and optional scenario folder:
//         <outDir>/<namespace_lower>/samples/<token1>/<token2>/.../<scenario>/main.tf
//       e.g. .../microsoft.appplatform/samples/spring/apps/deployments/basic/main.tf
//   - remarks.json location: written at the same level as the namespace's samples folder:
//         <outDir>/<namespace_lower>/remarks.json
//       Only the TerraformSamples array is updated; other top-level fields are preserved.
//       If missing or unparsable, a new file is created with a default $schema value.
//   - Friendly names: if tools/generator-example-doc/resource_types.json is provided (or
//       -resourceTypes path is specified), descriptions use human-friendly resource names.
//       The mapping is case-insensitive on the resource type key.
//   - Logging: prints progress for discovery, copying, and remarks file generation.
//
// CLI
//   -inDir         Input root directory for examples (default depends on current working dir)
//   -outDir        Output root directory for docs structure (default depends on current dir)
//   -resourceTypes Optional path to resource_types.json to enable friendly names
//
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// This tool syncs Terraform example files from inDir to outDir in a docs-friendly structure
// and generates/updates a remarks.json file per namespace next to the samples folder.

// Default CLI values are determined dynamically based on where the tool is executed from.
// - If run from repo root: inDir = "examples", outDir = "../bicep-refdocs-generator/settings/remarks"
// - If run from tools/sync-reference-doc: inDir = "../../examples", outDir = "../../../bicep-refdocs-generator/settings/remarks"

type sampleInfo struct {
	Namespace       string   // e.g. Microsoft.AlertsManagement
	TypePathTokens  []string // e.g. [actionRules] or [service apis diagnostics]
	APIVersion      string   // e.g. 2021-08-08
	InputDir        string   // absolute or relative path to the input sample dir
	InputMainTfPath string
	ScenarioName    string // e.g. basic, advanced; derived from subfolder name under the resource folder
}

type remarksDoc struct {
	Schema           string            `json:"$schema"`
	TerraformSamples []terraformSample `json:"TerraformSamples"`
}

type terraformSample struct {
	ResourceType string `json:"ResourceType"`
	Path         string `json:"Path"`
	Description  string `json:"Description"`
}

// loaded at startup from tools/generator-example-doc/resource_types.json if present
// (can be overridden via -resourceTypes). Keys are normalized to lowercase for case-insensitive lookups.
var friendlyNames map[string]string

func main() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}

	var defIn, defOut, resourceTypeJsonFile string
	switch filepath.Base(wd) {
	case "sync-reference-doc":
		defIn = filepath.Join(wd, "..", "..", "examples")
		defOut = filepath.Join(wd, "..", "..", "..", "bicep-refdocs-generator", "settings", "remarks")
		resourceTypeJsonFile = filepath.Join(wd, "..", "generator-example-doc", "resource_types.json")
	case "terraform-provider-azapi":
		defIn = filepath.Join(wd, "examples")
		defOut = filepath.Join(wd, "..", "bicep-refdocs-generator", "settings", "remarks")
		resourceTypeJsonFile = filepath.Join(wd, "tools", "generator-example-doc", "resource_types.json")
	default:
	}

	inDirFlag := flag.String("inDir", defIn, "input directory containing sample folders (e.g. examples)")
	outDirFlag := flag.String("outDir", defOut, "output root directory for remarks structure")
	resourceTypeJsonFileFlag := flag.String("resourceTypes", resourceTypeJsonFile, "path to resource_types.json for friendly names (optional)")
	flag.Parse()

	inDir := *inDirFlag
	outDir := *outDirFlag
	resourceTypeJsonFile = *resourceTypeJsonFileFlag

	log.Printf("sync-reference-doc: inDir=%s outDir=%s resourceTypes=%s", inDir, outDir, resourceTypeJsonFile)

	// Load friendly names for resource types (best-effort)
	if resourceTypeJsonFile != "" {
		if err := initFriendlyNames(resourceTypeJsonFile); err != nil {
			fmt.Fprintf(os.Stderr, "warn: could not load friendly names: %v\n", err)
		} else {
			log.Printf("friendly names loaded: %d", len(friendlyNames))
		}
	} else {
		log.Printf("friendly names disabled: no resourceTypes file provided")
	}

	log.Printf("scanning resource folders in: %s", inDir)
	entries, err := os.ReadDir(inDir)
	if err != nil {
		log.Fatalf("failed to read input directory %s: %v", inDir, err)
	}
	log.Printf("found %d resource folders", len(entries))

	var samples []sampleInfo
	for _, e := range entries {
		if !e.IsDir() {
			// skip files
			continue
		}

		// e is the resource folder like Microsoft.SqlVirtualMachine_sqlVirtualMachines@2023-10-01
		log.Printf("resource folder: %s", e.Name())
		parsed, err := parseSampleDirName(e.Name())
		if err != nil {
			// Ignore unrecognized dirs; keep going
			fmt.Fprintf(os.Stderr, "warn: skipping %s: %v\n", e.Name(), err)
			continue
		}
		resDir := filepath.Join(inDir, e.Name())
		scenEntries, err := os.ReadDir(resDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warn: skipping %s: cannot read scenarios: %v\n", e.Name(), err)
			continue
		}
		foundAny := false
		for _, se := range scenEntries {
			if !se.IsDir() {
				continue
			}
			scenario := se.Name()
			mainTf := filepath.Join(resDir, scenario, "main.tf")
			if _, err := os.Stat(mainTf); err != nil {
				fmt.Fprintf(os.Stderr, "warn: skipping %s/%s: no main.tf found\n", e.Name(), scenario)
				continue
			}
			s := parsed
			s.InputDir = filepath.Join(resDir, scenario)
			s.InputMainTfPath = mainTf
			s.ScenarioName = sanitizeScenarioName(scenario)
			samples = append(samples, s)
			log.Printf("  scenario: %s (main.tf)", s.ScenarioName)
			foundAny = true
		}
		if !foundAny {
			// As a last resort, allow main.tf directly under the resource folder
			mainTf := filepath.Join(resDir, "main.tf")
			if _, err := os.Stat(mainTf); err == nil {
				s := parsed
				s.InputDir = resDir
				s.InputMainTfPath = mainTf
				s.ScenarioName = ""
				samples = append(samples, s)
				log.Printf("  scenario: <root> (main.tf)")
			} else {
				fmt.Fprintf(os.Stderr, "warn: skipping %s: no main.tf found\n", e.Name())
			}
		}
	}

	// summary logged after building namespace map

	// Group by namespace for remarks generation
	byNS := make(map[string][]sampleInfo)
	for _, s := range samples {
		byNS[s.Namespace] = append(byNS[s.Namespace], s)
	}
	log.Printf("total samples discovered: %d (namespaces=%d)", len(samples), len(byNS))

	// Copy files and collect remarks per namespace
	for ns, list := range byNS {
		nsLower := strings.ToLower(ns)
		nsDir := filepath.Join(outDir, nsLower)
		samplesDir := filepath.Join(nsDir, "samples")
		if err := os.MkdirAll(samplesDir, 0o755); err != nil {
			log.Fatalf("failed to create samples dir %s: %v", samplesDir, err)
		}

		var remarks []terraformSample
		for _, s := range list {
			// Sample directory path is all type tokens nested, lower-cased
			lowerTokens := make([]string, len(s.TypePathTokens))
			for i, t := range s.TypePathTokens {
				lowerTokens[i] = strings.ToLower(t)
			}
			destSampleDir := filepath.Join(append([]string{samplesDir}, lowerTokens...)...)
			if s.ScenarioName != "" {
				destSampleDir = filepath.Join(destSampleDir, strings.ToLower(s.ScenarioName))
			}
			if err := os.MkdirAll(destSampleDir, 0o755); err != nil {
				log.Fatalf("failed to create sample dir %s: %v", destSampleDir, err)
			}
			// Copy main.tf
			destMain := filepath.Join(destSampleDir, "main.tf")
			log.Printf("copy: %s -> %s", s.InputMainTfPath, destMain)
			if err := copyFile(s.InputMainTfPath, destMain); err != nil {
				log.Fatalf("failed to copy %s to %s: %v", s.InputMainTfPath, destMain, err)
			}

			// Build remarks entry
			resourceType := ns + "/" + strings.Join(s.TypePathTokens, "/")
			rel := []string{"samples"}
			rel = append(rel, lowerTokens...)
			if s.ScenarioName != "" {
				rel = append(rel, strings.ToLower(s.ScenarioName))
			}
			relPath := filepath.ToSlash(filepath.Join(append(rel, "main.tf")...))
			desc := defaultDescription(resourceType, s.APIVersion, s.ScenarioName)
			remarks = append(remarks, terraformSample{
				ResourceType: resourceType,
				Path:         relPath,
				Description:  desc,
			})
		}

		// Sort remarks entries for stable output
		sort.Slice(remarks, func(i, j int) bool {
			if remarks[i].ResourceType == remarks[j].ResourceType {
				return remarks[i].Path < remarks[j].Path
			}
			return remarks[i].ResourceType < remarks[j].ResourceType
		})

		// Write or merge remarks.json at the same level as the samples folder
		remarksPath := filepath.Join(nsDir, "remarks.json")
		log.Printf("write remarks: %s (entries=%d)", remarksPath, len(remarks))
		if err := writeOrMergeRemarks(remarksPath, remarks); err != nil {
			log.Fatalf("failed to write remarks.json for %s: %v", ns, err)
		}
	}

	log.Printf("sync complete: namespaces=%d samples=%d", len(byNS), len(samples))
}

// initFriendlyNames initializes the friendlyNames map by reading the JSON file
func initFriendlyNames(resourceTypeJsonFile string) error {
	friendlyNames = make(map[string]string)
	b, err := os.ReadFile(resourceTypeJsonFile)
	if err != nil {
		return err
	}
	var entries []struct {
		ResourceType string `json:"resourceType"`
		FriendlyName string `json:"friendlyName"`
	}
	if err := json.Unmarshal(b, &entries); err != nil {
		return err
	}
	for _, e := range entries {
		if e.ResourceType == "" || e.FriendlyName == "" {
			continue
		}
		friendlyNames[strings.ToLower(e.ResourceType)] = e.FriendlyName
	}
	return nil
}

func parseSampleDirName(name string) (sampleInfo, error) {
	// Expected format: <Namespace>_<typePathTokens joined by '_'>@<apiVersion>
	// Examples:
	//   Microsoft.AlertsManagement_actionRules@2021-08-08
	//   Microsoft.ApiManagement_service_apis_diagnostics@2021-08-01
	var out sampleInfo
	if name == "" {
		return out, errors.New("empty name")
	}
	at := strings.LastIndex(name, "@")
	if at <= 0 || at == len(name)-1 {
		return out, fmt.Errorf("missing or malformed @version in %q", name)
	}
	prefix := name[:at]
	version := name[at+1:]
	// Split prefix by '_' -> first part is namespace, rest are type tokens
	parts := strings.Split(prefix, "_")
	if len(parts) < 2 {
		return out, fmt.Errorf("expected <Namespace>_<TypePath>@<Version>, got %q", name)
	}
	ns := parts[0]
	typeTokens := parts[1:]
	out.Namespace = ns
	out.TypePathTokens = typeTokens
	out.APIVersion = version
	return out, nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	tmp := dst + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	cerr := out.Close()
	if err == nil {
		err = cerr
	}
	if err != nil {
		_ = os.Remove(tmp)
		return err
	}
	// Best-effort to copy file mode from src
	if fi, err := os.Stat(src); err == nil {
		_ = os.Chmod(tmp, fi.Mode())
	}
	return os.Rename(tmp, dst)
}

func writeJSON(path string, v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	b = append(b, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, b, 0o644)
}

// writeOrMergeRemarks writes a remarks.json file at path. If a file already exists,
// it preserves existing top-level fields and only updates the TerraformSamples array.
// If the file doesn't exist or cannot be parsed, it creates a new document with the
// default $schema value and the provided samples.
func writeOrMergeRemarks(path string, samples []terraformSample) error {
	var existing map[string]any
	b, err := os.ReadFile(path)
	if err == nil {
		if json.Unmarshal(b, &existing) != nil {
			existing = nil // treat as not existing / unparsable
		}
	}
	if existing == nil {
		// Create a new document
		doc := remarksDoc{
			Schema:           "../../schemas/remarks.schema.json",
			TerraformSamples: samples,
		}
		return writeJSON(path, doc)
	}
	// Merge: set TerraformSamples, keep others as-is
	existing["TerraformSamples"] = samples
	// Ensure $schema exists; if not, add default
	if _, ok := existing["$schema"]; !ok {
		existing["$schema"] = "../../schemas/remarks.schema.json"
	}
	return writeJSON(path, existing)
}

// sanitizeScenarioName removes underscores and trims whitespace from scenario names.
func sanitizeScenarioName(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	return strings.ReplaceAll(s, "_", "")
}

func defaultDescription(resourceType, apiVersion, scenario string) string {
	// Normalize a human-friendly scenario label
	scen := strings.TrimSpace(strings.ToLower(scenario))
	if strings.Contains(scen, "_") {
		scen = strings.ReplaceAll(scen, "_", "")
	}
	var scenPrefix string
	switch scen {
	case "", "default":
		scenPrefix = "A basic example of"
	case "advanced":
		scenPrefix = "An advanced example of"
	case "complete", "full":
		scenPrefix = "A complete example of"
	case "minimal", "basic":
		scenPrefix = "A basic example of"
	default:
		// Use the scenario name verbatim when not one of the known ones
		scenPrefix = fmt.Sprintf("A %s example of", scen)
	}

	if friendlyNames != nil {
		if fn, ok := friendlyNames[strings.ToLower(resourceType)]; ok && fn != "" {
			return fmt.Sprintf("%s deploying %s.", scenPrefix, fn)
		}
	}
	return fmt.Sprintf("%s deploying %s (%s).", scenPrefix, resourceType, apiVersion)
}
