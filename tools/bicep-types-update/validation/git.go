package validation

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

var gitVersionRe = regexp.MustCompile(`(\d+)\.(\d+)`)

func CheckGitVersion(minGitMajor, minGitMinor int) error {
	out, err := exec.Command("git", "--version").Output()
	if err != nil {
		return fmt.Errorf("git not found or not working: %w", err)
	}

	match := gitVersionRe.FindStringSubmatch(string(out))
	if match == nil {
		return fmt.Errorf("unable to parse git version from %q", string(out))
	}

	major, err := strconv.Atoi(match[1])
	if err != nil {
		return fmt.Errorf("unable to parse git major version from %q: %w", string(out), err)
	}
	minor, err := strconv.Atoi(match[2])
	if err != nil {
		return fmt.Errorf("unable to parse git minor version from %q: %w", string(out), err)
	}

	if major < minGitMajor || (major == minGitMajor && minor < minGitMinor) {
		return fmt.Errorf("git %d.%d or later is required, found %d.%d", minGitMajor, minGitMinor, major, minor)
	}

	return nil
}
