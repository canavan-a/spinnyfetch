package modules

import (
	"fmt"
	"os/exec"
	"strings"
)

type packagesModule struct{}

func (packagesModule) Key() string { return "Packages" }

func (packagesModule) Fetch() (string, error) {
	sys := countLines("nix-store", "-q", "--requisites", "/run/current-system/sw")
	user := countLines("nix", "profile", "list")

	var parts []string
	if sys >= 0 {
		parts = append(parts, fmt.Sprintf("%d (nix-system)", sys))
	}
	if user >= 0 {
		parts = append(parts, fmt.Sprintf("%d (nix-user)", user))
	}
	if len(parts) == 0 {
		return "", fmt.Errorf("no package managers found")
	}
	return strings.Join(parts, ", "), nil
}

// countLines runs a command and counts non-empty output lines, or returns
// -1 if the command isn't available or fails.
func countLines(name string, args ...string) int {
	out, err := exec.Command(name, args...).Output()
	if err != nil {
		return -1
	}
	n := 0
	for _, l := range strings.Split(string(out), "\n") {
		if strings.TrimSpace(l) != "" {
			n++
		}
	}
	return n
}
