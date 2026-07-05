package git

import (
	"fmt"
	"os"
	"os/exec"
)

func ShallowCloneAtCommit(repoURL, dest, commitSHA string) error {
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

	if err := Run("-C", dest, "init", "-q"); err != nil {
		return err
	}
	if err := Run("-C", dest, "remote", "add", "origin", repoURL); err != nil {
		return err
	}
	if err := Run("-C", dest, "fetch", "--depth", "1", "origin", commitSHA); err != nil {
		return err
	}
	if err := Run("-C", dest, "checkout", "FETCH_HEAD"); err != nil {
		return err
	}
	if err := Run("-C", dest, "submodule", "update", "--init", "--recursive", "--depth", "1"); err != nil {
		return err
	}

	return nil
}

func Run(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
