package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"gopkg.in/yaml.v3"
)

func main() {
	bicepTypesAzURL := "https://github.com/Azure/bicep-types-az.git"
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Fprintln(os.Stderr, "failed to determine source file location")
		os.Exit(1)
	}

	dir := filepath.Dir(currentFile)
	dest := filepath.Join(dir, "bicep-types-az")
	configPath := filepath.Join(dir, "config.yml")
	commitSHA, err := readCommitSHA(configPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read config failed: %v\n", err)
		os.Exit(1)
	}

	if err := cloneAtCommit(bicepTypesAzURL, dest, commitSHA); err != nil {
		fmt.Fprintf(os.Stderr, "clone failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Cloned %s at %s into %s\n", bicepTypesAzURL, commitSHA, dest)
}

func readCommitSHA(configPath string) (string, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf("read config: %w", err)
	}

	var config struct {
		CommitSHA string `yaml:"bicep-types-az-commit-sha"`
	}
	if err := yaml.Unmarshal(data, &config); err != nil {
		return "", fmt.Errorf("parse config: %w", err)
	}
	if config.CommitSHA == "" {
		return "", fmt.Errorf("commit SHA not found in %s", configPath)
	}

	return config.CommitSHA, nil
}

func cloneAtCommit(repoURL, dest, commitSHA string) error {
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
