package validation

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
)

var nodeVersionRe = regexp.MustCompile(`v?(\d+)`)

func CheckNodeVersion(minNodeMajor int) error {
	out, err := exec.Command("node", "--version").Output()
	if err != nil {
		return fmt.Errorf("node not found or not working: %w", err)
	}

	major, err := parseNodeVersion(out)
	if err != nil {
		return err
	}

	if major < minNodeMajor {
		return fmt.Errorf("node %d or later is required, found %d", minNodeMajor, major)
	}

	return nil
}

func parseNodeVersion(out []byte) (int, error) {
	match := nodeVersionRe.FindStringSubmatch(string(out))
	if match == nil {
		return 0, fmt.Errorf("unable to parse node version from %q", string(out))
	}

	major, err := strconv.Atoi(match[1])
	if err != nil {
		return 0, fmt.Errorf("unable to parse node major version from %q: %w", string(out), err)
	}

	return major, nil
}
