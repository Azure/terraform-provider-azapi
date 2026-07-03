package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/Azure/terraform-provider-azapi/tools/bicep-types-update/config"
	"github.com/Azure/terraform-provider-azapi/tools/bicep-types-update/validation"
)

func main() {
	if err := validation.CheckGitVersion(2, 10); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if err := validation.CheckNodeVersion(24); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	bicepTypesAzURL := "https://github.com/Azure/bicep-types-az.git"

	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "failed to determine source file location")
		os.Exit(1)
	}

	dir := filepath.Dir(currentFile)

	configPath := filepath.Join(dir, "config.yml")
	cfg, err := config.ReadConfig(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config failed: %v\n", err)
		os.Exit(1)
	}

	restAPISpecsDir := filepath.Join(dir, "azure-rest-api-specs")
	azureRestAPISpecsURL := "https://github.com/Azure/azure-rest-api-specs.git"
	if os.Getenv("SKIP_CLONE_AND_PATCH") != "true" {
		fmt.Printf("➡️  Cloning %s at %s into %s\n", azureRestAPISpecsURL, cfg.AzureRestAPISpecsCommitSHA, restAPISpecsDir)
		if err := shallowCloneAtCommit(azureRestAPISpecsURL, restAPISpecsDir, cfg.AzureRestAPISpecsCommitSHA); err != nil {
			fmt.Fprintf(os.Stderr, "clone failed: %v\n", err)
			os.Exit(1)
		}
	}

	bicepTypesAzDir := filepath.Join(dir, "bicep-types-az")
	if os.Getenv("SKIP_CLONE_AND_PATCH") != "true" {
		fmt.Printf("➡️  Cloning %s at %s into %s\n", bicepTypesAzURL, cfg.BicepTypesAzCommitSHA, bicepTypesAzDir)
		if err := shallowCloneAtCommit(bicepTypesAzURL, bicepTypesAzDir, cfg.BicepTypesAzCommitSHA); err != nil {
			fmt.Fprintf(os.Stderr, "clone failed: %v\n", err)
			os.Exit(1)
		}
	}

	bicepTypesPatch := filepath.Join(dir, "bicep-types-support-all-actions.patch")
	bicepTypesDir := filepath.Join(bicepTypesAzDir, "bicep-types")
	if os.Getenv("SKIP_CLONE_AND_PATCH") != "true" {
		fmt.Printf("➡️  Applying patch %s to %s\n", bicepTypesPatch, bicepTypesDir)
		if err := runGit("-C", bicepTypesDir, "apply", bicepTypesPatch); err != nil {
			fmt.Fprintf(os.Stderr, "apply patch failed: %v\n", err)
			os.Exit(1)
		}
	}

	bicepTypesAzPatch := filepath.Join(dir, "bicep-types-az-support-all-actions.patch")
	if os.Getenv("SKIP_CLONE_AND_PATCH") != "true" {
		fmt.Printf("➡️  Applying patch %s to %s\n", bicepTypesAzPatch, bicepTypesAzDir)
		if err := runGit("-C", bicepTypesAzDir, "apply", bicepTypesAzPatch); err != nil {
			fmt.Fprintf(os.Stderr, "apply patch failed: %v\n", err)
			os.Exit(1)
		}
	}

	bicepTypesSrcBicepTypesDir := filepath.Join(bicepTypesDir, "src", "bicep-types")

	fmt.Printf("➡️  Running npm ci in %s\n", bicepTypesSrcBicepTypesDir)
	if err := runCommandInDir(bicepTypesSrcBicepTypesDir, npmCommand(), "ci"); err != nil {
		fmt.Fprintf(os.Stderr, "npm ci failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("➡️  Running npm run build in %s\n", bicepTypesSrcBicepTypesDir)
	if err := runCommandInDir(bicepTypesSrcBicepTypesDir, npmCommand(), "run", "build"); err != nil {
		fmt.Fprintf(os.Stderr, "npm run build failed: %v\n", err)
		os.Exit(1)
	}

	autorestBicepDir := filepath.Join(bicepTypesAzDir, "src", "autorest.bicep")

	fmt.Printf("➡️  Running npm ci in %s\n", autorestBicepDir)
	if err := runCommandInDir(autorestBicepDir, npmCommand(), "ci"); err != nil {
		fmt.Fprintf(os.Stderr, "npm ci failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("➡️  Running npm run build in %s\n", autorestBicepDir)
	if err := runCommandInDir(autorestBicepDir, npmCommand(), "run", "build"); err != nil {
		fmt.Fprintf(os.Stderr, "npm run build failed: %v\n", err)
		os.Exit(1)
	}

	generatorDir := filepath.Join(bicepTypesAzDir, "src", "generator")

	fmt.Printf("➡️  Running npm ci in %s\n", generatorDir)
	if err := runCommandInDir(generatorDir, npmCommand(), "ci"); err != nil {
		fmt.Fprintf(os.Stderr, "npm ci failed: %v\n", err)
		os.Exit(1)
	}
}

func shallowCloneAtCommit(repoURL, dest, commitSHA string) error {
	if dest == "" {
		return fmt.Errorf("destination cannot be empty")
	}
	if commitSHA == "" {
		return fmt.Errorf("commit SHA cannot be empty")
	}

	if err := os.RemoveAll(dest); err != nil {
		return fmt.Errorf("failed to remove existing destination: %w", err)
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return fmt.Errorf("failed to create destination: %w", err)
	}

	if err := runGit("-C", dest, "init", "-q"); err != nil {
		return err
	}
	if err := runGit("-C", dest, "remote", "add", "origin", repoURL); err != nil {
		return err
	}
	if err := runGit("-C", dest, "fetch", "--depth", "1", "origin", commitSHA); err != nil {
		return err
	}
	if err := runGit("-C", dest, "checkout", "FETCH_HEAD"); err != nil {
		return err
	}
	if err := runGit("-C", dest, "submodule", "update", "--init", "--recursive", "--depth", "1"); err != nil {
		return err
	}

	return nil
}

func runGit(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runCommandInDir(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func npmCommand() string {
	if runtime.GOOS == "windows" {
		return "npm.cmd"
	}
	return "npm"
}
